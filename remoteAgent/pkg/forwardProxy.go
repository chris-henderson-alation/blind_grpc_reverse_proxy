package grpcinverter

import (
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc/peer"

	"go.uber.org/zap"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter/ioc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
)

type ForwardProxy struct {
	listeners *ListenerMap
	pending   *PendingJobs

	internalAddr string
	internalPort uint16
	externalAddr string
	externalPort uint16

	internal *grpc.Server
	external *grpc.Server
	ioc.UnimplementedGrpcInverterServer
}

func NewForwardProxy(internalAddr string, internalPort uint16, externalAddr string, externalPort uint16) *ForwardProxy {
	forwardProxy := &ForwardProxy{
		internalAddr: internalAddr,
		internalPort: internalPort,
		externalAddr: externalAddr,
		externalPort: externalPort,
		listeners:    NewListenerMap(),
		pending:      NewPendingJobs(),
	}
	internalServer := grpc.NewServer(
		grpc.KeepaliveParams(KEEPALIVE_SERVER_PARAMETERS),
		grpc.KeepaliveEnforcementPolicy(KEEPALIVE_ENFORCEMENT_POLICY),
		grpc.ForceServerCodec(NoopCodec{}),
		grpc.UnknownServiceHandler(forwardProxy.GenericStreamHandler()))
	externalServer := grpc.NewServer(
		grpc.KeepaliveParams(KEEPALIVE_SERVER_PARAMETERS),
		grpc.KeepaliveEnforcementPolicy(KEEPALIVE_ENFORCEMENT_POLICY))
	ioc.RegisterGrpcInverterServer(externalServer, forwardProxy)
	forwardProxy.internal = internalServer
	forwardProxy.external = externalServer
	return forwardProxy
}

func (proxy *ForwardProxy) JobTunnel(_ *empty.Empty, agent ioc.GrpcInverter_JobTunnelServer) error {
	p, _ := peer.FromContext(agent.Context())
	agentId, err := ExtractAgentId(agent.Context())
	if err != nil {
		// @TODO Don't try too hard here, this will get replaced by looking up the contents of the cert.
		LOGGER.Error("An agent failed to identify itself")
		return err
	}
	jobs, alreadyRegistered := proxy.listeners.Register(agentId)
	if alreadyRegistered {
		LOGGER.Warn("An agent attempted to connect to the job tunnel, however an agent with its ID is already registered",
			zap.Uint64("agent", agentId),
			zap.String("agentIP", p.Addr.String()))
		return fmt.Errorf("already registered")
	}
	LOGGER.Info("An agent has connected to its job tunnel.",
		zap.Uint64("agent", agentId),
		zap.String("agentIP", p.Addr.String()))
	for {
		select {
		case <-agent.Context().Done():
			LOGGER.Info("An agent has disconnected from its job tunnel",
				zap.Uint64("agent", agentId),
				zap.String("agentIP", p.Addr.String()))
			proxy.listeners.Unregister(agentId)
			return nil
		case job := <-jobs:
			LOGGER.Info("Sending job to agent",
				zap.String("method", job.Method),
				zap.Uint64("job", job.JobId),
				zap.Uint64("connector", job.Connector),
				zap.String("agentIP", p.Addr.String()))
			err := agent.Send(&ioc.Job{
				JobID:     job.JobId,
				Method:    job.Method,
				Connector: job.Connector,
			})
			if err != nil {
				LOGGER.Error("An error occured while sending a job to an agent",
					zap.String("method", job.Method),
					zap.Uint64("job", job.JobId),
					zap.Uint64("connector", job.Connector),
					zap.String("agentIP", p.Addr.String()),
					zap.Error(err))
				s, ok := status.FromError(err)
				if !ok {
					job.Alation.SendError(status.New(codes.Unavailable, err.Error()))
				} else {
					job.Alation.SendError(s)
				}
				continue
			}
			LOGGER.Info("Successfully sent job to agent and awaiting callback",
				zap.String("method", job.Method),
				zap.Uint64("job", job.JobId),
				zap.Uint64("connector", job.Connector),
				zap.String("agentIP", p.Addr.String()))
			proxy.pending.Submit(job)
		}
	}
}

// Pipe is the endpoint that an agent calls back to to begin piping job results
// from the connector and back up to Alation.
func (proxy *ForwardProxy) Pipe(agent ioc.GrpcInverter_PipeServer) error {
	agentId, err := ExtractAgentId(agent.Context())
	if err != nil {
		return err
	}
	jobId, err := ExtractJobId(agent.Context())
	if err != nil {
		return err
	}
	job := proxy.pending.Retrieve(agentId, jobId)
	if job == nil {
		return fmt.Errorf("nice catch blanco nino, but too bad your ass got saaaaaacked")
	}
	p, _ := peer.FromContext(agent.Context())
	LOGGER.Info("Agent established callback stream",
		zap.String("method", job.Method),
		zap.Uint64("job", job.JobId),
		zap.Uint64("connector", job.Connector),
		zap.String("agentIP", p.Addr.String()))
	ioc.ForwardProxy(job.Alation, agent)
	return nil
}

// GenericStreamHandler is grpc.StreamHandler that gets registered within the Alation facing
// gRPC server as the fallback "UnknownServiceHandler". This means that any request that comes
// in that cannot be handled by any other registered gRPC server is instead forwarded to this stream handler.
//
// In our particular case, that is all of them! The ForwardProxy does not register any other gRPC service, meaning
// that this handler receives any-and-all incoming gRPC requests regardless of their original protobuf definitions.
//
// This is critical to our scheme as we need to be able to accept ANY incoming gRPC request without having to maintain
// symmetric protobuf definitions (which would be an egregious maintenance burden for us all).
func (proxy *ForwardProxy) GenericStreamHandler() grpc.StreamHandler {
	return func(_ interface{}, upstream grpc.ServerStream) error {
		// @TODO check return just in case
		p, _ := peer.FromContext(upstream.Context())
		fullMethodName, ok := grpc.MethodFromServerStream(upstream)
		if !ok {
			LOGGER.Error("No method was attached to the incoming gRPC call",
				zap.String("alationIP", p.Addr.String()))
			return status.Error(codes.Internal, "@TODO this is comically fatal")
		}
		agentId, err := ExtractAgentId(upstream.Context())
		if err != nil {
			LOGGER.Error("No agent header was attached to the incoming request",
				zap.String("method", fullMethodName),
				zap.String("wanted", AgentIdHeader),
				zap.String("alationIP", p.Addr.String()))
			return err
		}
		connectorId, err := ExtractConnectorId(upstream.Context())
		if err != nil {
			LOGGER.Error("No connector header was attached to the incoming request",
				zap.String("method", fullMethodName),
				zap.Uint64("agent", agentId),
				zap.String("wanted", ConnectorIdHeader),
				zap.String("alationIP", p.Addr.String()))
			return err
		}
		jobId, err := ExtractJobId(upstream.Context())
		if err != nil {
			LOGGER.Error("No job header was attached to the incoming request",
				zap.String("method", fullMethodName),
				zap.Uint64("agent", agentId),
				zap.Uint64("connector", connectorId),
				zap.String("wanted", JobIdHeader),
				zap.String("alationIP", p.Addr.String()))
			return err
		}
		LOGGER.Info("Received a job from Alation",
			zap.String("method", fullMethodName),
			zap.Uint64("connector", connectorId),
			zap.Uint64("agent", agentId),
			zap.Uint64("job", jobId),
			zap.String("alationIP", p.Addr.String()))
		errors := make(chan error)
		call := &GrpcCall{
			Method:    fullMethodName,
			Agent:     agentId,
			Connector: connectorId,
			JobId:     jobId,
			Alation: ioc.ServerStreamWithError{
				ServerStream: upstream,
				Error:        errors,
			},
		}
		err = proxy.Submit(call)
		if err != nil {
			// This is a warning because it is likely simply just a disconnected agent.
			LOGGER.Warn("Failed to submit job to connector",
				zap.String("method", fullMethodName),
				zap.Uint64("connector", connectorId),
				zap.Uint64("agent", agentId),
				zap.Uint64("job", jobId),
				zap.String("alationIP", p.Addr.String()),
				zap.Error(err))
			return err
		}
		return <-errors
	}
}

// Submit submits a job to the listening agent. If the agent is not currently connected
// then callers will receive an error immediately informing them as such.
func (proxy *ForwardProxy) Submit(call *GrpcCall) error {
	return proxy.listeners.Submit(call)
}

// Start starts both the Alation facing and internet facing servers for the forward proxy.
//
// This method BLOCKS until both servers are shutdown (likely via the Stop method). if you wish
// to run this server concurrently (say, for example, in unit tests) then simply run `go proxy.Start()`.
func (proxy *ForwardProxy) Start() error {
	internalListener, err := net.Listen("tcp", proxy.InternalAddress())
	if err != nil {
		return err
	}
	LOGGER.Info("Began listening on internal address",
		zap.String("address", proxy.internalAddr),
		zap.Uint16("port", proxy.internalPort),
		zap.String("protocol", "tcp"))
	externalListener, err := net.Listen("tcp", proxy.ExternalAddress())
	if err != nil {
		return err
	}
	LOGGER.Info("Began listening on external address",
		zap.String("address", proxy.externalAddr),
		zap.Uint16("port", proxy.externalPort),
		zap.String("protocol", "tcp"))
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := proxy.internal.Serve(internalListener)
		LOGGER.Warn("Internal server stopped",
			zap.String("address", proxy.internalAddr),
			zap.Uint16("port", proxy.internalPort),
			zap.String("protocol", "tcp"),
			zap.Error(err))
	}()
	go func() {
		defer wg.Done()
		err := proxy.external.Serve(externalListener)
		LOGGER.Warn("External server stopped",
			zap.String("address", proxy.externalAddr),
			zap.Uint16("port", proxy.externalPort),
			zap.String("protocol", "tcp"),
			zap.Error(err))
	}()
	wg.Wait()
	return nil
}

// Stop stops both the Alation facing and internet facing gRPC servers.
func (proxy *ForwardProxy) Stop() {
	proxy.internal.Stop()
	proxy.external.Stop()
}

// InternalAddress returns the "<addr>:<port>" string for the Alation facing server.
//
// This is mostly useful for unit testing.
func (proxy ForwardProxy) InternalAddress() string {
	return fmt.Sprintf("%s:%d", proxy.internalAddr, proxy.internalPort)
}

// ExternalAddress returns the "<addr>:<port>" string for the internet facing server.
//
// This is mostly useful for unit testing.
func (proxy ForwardProxy) ExternalAddress() string {
	return fmt.Sprintf("%s:%d", proxy.externalAddr, proxy.externalPort)
}

type GrpcCall struct {
	Alation   ioc.ServerStreamWithError
	Method    string
	Connector uint64
	Agent     uint64
	JobId     uint64
}

package forward

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/shared"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/logging"
	"github.com/Alation/alation_connector_manager/docker/remoteAgent/protocol"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Proxy struct {
	Listeners *ListenerMap
	pending   *PendingJobs

	internalAddr string
	internalPort uint16
	ExternalAddr string
	ExternalPort uint16

	internal *grpc.Server
	external *grpc.Server
	protocol.UnimplementedGrpcInverterServer

	ctx  context.Context
	stop context.CancelFunc
}

func NewProxy(internalAddr string, internalPort uint16, externalAddr string, externalPort uint16) *Proxy {
	ctx, stop := context.WithCancel(context.Background())
	forwardProxy := &Proxy{
		internalAddr: internalAddr,
		internalPort: internalPort,
		ExternalAddr: externalAddr,
		ExternalPort: externalPort,
		Listeners:    NewListenerMap(),
		pending:      NewPendingJobs(),
		ctx:          ctx,
		stop:         stop,
	}
	internalServer := grpc.NewServer(
		grpc.KeepaliveParams(shared.KEEPALIVE_SERVER_PARAMETERS),
		grpc.KeepaliveEnforcementPolicy(shared.KEEPALIVE_ENFORCEMENT_POLICY),
		grpc.ForceServerCodec(shared.NoopCodec{}),
		grpc.UnknownServiceHandler(forwardProxy.GenericStreamHandler()))
	externalServer := grpc.NewServer(
		grpc.KeepaliveParams(shared.KEEPALIVE_SERVER_PARAMETERS),
		grpc.KeepaliveEnforcementPolicy(shared.KEEPALIVE_ENFORCEMENT_POLICY))
	protocol.RegisterGrpcInverterServer(externalServer, forwardProxy)
	forwardProxy.internal = internalServer
	forwardProxy.external = externalServer
	return forwardProxy
}

func (proxy *Proxy) JobTunnel(_ *empty.Empty, agent protocol.GrpcInverter_JobTunnelServer) error {
	agentId, err := shared.ExtractAgentId(agent.Context())
	if err != nil {
		// @TODO Don't try too hard here, this will get replaced by looking up the contents of the cert.
		logging.LOGGER.Error("An agent failed to identify itself", logging.PeerAgent(agent.Context()))
		return err
	}
	jobs := proxy.Listeners.Register(agentId)
	if jobs == nil {
		logging.LOGGER.Warn("An agent attempted to connect to the job tunnel, however an agent with its ID is already registered",
			logging.Agent(agentId),
			logging.PeerAgent(agent.Context()))
		return fmt.Errorf("already registered")
	}
	logger := logging.LOGGER.With(
		logging.Agent(agentId),
		logging.PeerAgent(agent.Context()))
	for {
		select {
		case <-agent.Context().Done():
			logger.Info("An agent has disconnected from its job tunnel")
			proxy.Listeners.Unregister(agentId)
			return nil
		case job := <-jobs:
			logger.Info("Sending job to agent",
				logging.Method(job.Method),
				logging.Connector(job.Connector),
				logging.Job(job.JobId))
			err := agent.Send(&protocol.Job{
				JobID:     job.JobId,
				Method:    job.Method,
				Connector: job.Connector,
			})
			if err != nil {
				logger.Error("An error occurred while sending a job to an agent",
					logging.Method(job.Method),
					logging.Connector(job.Connector),
					logging.Job(job.JobId),
					logging.Error(err))
				s, ok := status.FromError(err)
				if !ok {
					job.Alation.SendError(status.New(codes.Unavailable, err.Error()))
				} else {
					job.Alation.SendError(s)
				}
				continue
			}
			logger.Info("Successfully sent job to agent and awaiting callback",
				logging.Method(job.Method),
				logging.Connector(job.Connector),
				logging.Job(job.JobId))
			proxy.pending.Submit(job)
		}
	}
}

// Pipe is the endpoint that an agent calls back to to begin piping job results
// from the connector and back up to Alation.
func (proxy *Proxy) Pipe(agent protocol.GrpcInverter_PipeServer) error {
	agentId, err := shared.ExtractAgentId(agent.Context())
	if err != nil {
		return err
	}
	jobId, err := shared.ExtractJobId(agent.Context())
	if err != nil {
		return err
	}
	job := proxy.pending.Retrieve(agentId, jobId)
	if job == nil {
		return fmt.Errorf("nice catch blanco nino, but too bad your ass got saaaaaacked")
	}
	jobLogger := logging.LOGGER.With(
		logging.Method(job.Method),
		logging.Agent(agentId),
		logging.Connector(job.Connector),
		logging.Job(job.JobId),
		logging.PeerAlation(job.Alation.Context()),
		logging.PeerAgent(agent.Context()))
	jobLogger.Info("Agent established callback stream")
	protocol.ForwardProxy(job.Alation, agent, jobLogger)
	return nil
}

// GenericStreamHandler is grpc.StreamHandler that gets registered within the Alation facing
// gRPC server as the fallback "UnknownServiceHandler". This means that any request that comes
// in that cannot be handled by any other registered gRPC server is instead forwarded to this stream handler.
//
// In our particular case, that is all of them! The Proxy does not register any other gRPC service, meaning
// that this handler receives any-and-all incoming gRPC requests regardless of their original protobuf definitions.
//
// This is critical to our scheme as we need to be able to accept ANY incoming gRPC request without having to maintain
// symmetric protobuf definitions (which would be an egregious maintenance burden for us all).
func (proxy *Proxy) GenericStreamHandler() grpc.StreamHandler {
	return func(_ interface{}, upstream grpc.ServerStream) error {
		grpcMethod, ok := grpc.MethodFromServerStream(upstream)
		if !ok {
			logging.LOGGER.Error("No method was attached to the incoming gRPC call",
				logging.PeerAlation(upstream.Context()))
			return status.Error(codes.Internal, "@TODO this is comically fatal")
		}
		agentId, err := shared.ExtractAgentId(upstream.Context())
		if err != nil {
			logging.LOGGER.Error("No agent header was attached to the incoming request",
				logging.Method(grpcMethod),
				logging.Wanted(shared.AgentIdHeader),
				logging.PeerAlation(upstream.Context()))
			return err
		}
		connectorId, err := shared.ExtractConnectorId(upstream.Context())
		if err != nil {
			logging.LOGGER.Error("No connector header was attached to the incoming request",
				logging.Method(grpcMethod),
				logging.Agent(agentId),
				logging.Wanted(shared.ConnectorIdHeader),
				logging.PeerAlation(upstream.Context()))
			return err
		}
		jobId, err := shared.ExtractJobId(upstream.Context())
		if err != nil {
			logging.LOGGER.Error("No job header was attached to the incoming request",
				logging.Method(grpcMethod),
				logging.Agent(agentId),
				logging.Connector(connectorId),
				logging.Wanted(shared.JobIdHeader),
				logging.PeerAlation(upstream.Context()))
			return err
		}
		logger := logging.LOGGER.With(logging.Method(grpcMethod),
			logging.Agent(agentId),
			logging.Connector(connectorId),
			logging.Job(jobId),
			logging.PeerAlation(upstream.Context()))
		logger.Info("Received a job from Alation")
		errors := make(chan error)
		call := &GrpcCall{
			Method:    grpcMethod,
			Agent:     agentId,
			Connector: connectorId,
			JobId:     jobId,
			Alation: protocol.ServerStreamWithError{
				ServerStream: upstream,
				Error:        errors,
			},
		}
		err = proxy.Submit(call)
		if err != nil {
			// This is a warning because it is likely simply just a disconnected agent.
			logging.LOGGER.Warn("Failed to submit job to connector",
				logging.Method(grpcMethod),
				logging.Agent(agentId),
				logging.Connector(connectorId),
				logging.Job(jobId),
				logging.PeerAlation(upstream.Context()),
				logging.Error(err))
			return err
		}
		return <-errors
	}
}

// Submit submits a job to the listening agent. If the agent is not currently connected
// then callers will receive an error immediately informing them as such.
func (proxy *Proxy) Submit(call *GrpcCall) error {
	return proxy.Listeners.Submit(call)
}

const tcp = "tcp"

// Start starts both the Alation facing and internet facing servers for the forward proxy.
//
// This method BLOCKS until both servers are shutdown (likely via the Stop method). if you wish
// to run this server concurrently (say, for example, in unit tests) then simply run `go proxy.Start()`.
func (proxy *Proxy) Start() error {
	internalListener, err := net.Listen(tcp, proxy.InternalAddress())
	if err != nil {
		return err
	}
	logging.LOGGER.Info("Began listening on internal address",
		logging.Address(proxy.internalAddr),
		logging.Port(proxy.internalPort),
		logging.Protocol(tcp))
	externalListener, err := net.Listen(tcp, proxy.ExternalAddress())
	if err != nil {
		return err
	}
	logging.LOGGER.Info("Began listening on external address",
		logging.Address(proxy.ExternalAddr),
		logging.Port(proxy.ExternalPort),
		logging.Protocol(tcp))
	var internalServerError error
	var externalServerError error
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		internalServerError = proxy.internal.Serve(internalListener)
		logging.LOGGER.Warn("Internal server stopped",
			logging.Address(proxy.internalAddr),
			logging.Port(proxy.internalPort),
			logging.Protocol(tcp),
			logging.Error(err))
		proxy.stop()
	}()
	go func() {
		defer wg.Done()
		externalServerError = proxy.external.Serve(externalListener)
		logging.LOGGER.Warn("External server stopped",
			logging.Address(proxy.ExternalAddr),
			logging.Port(proxy.ExternalPort),
			logging.Protocol(tcp),
			logging.Error(err))
		proxy.stop()
	}()
	<-proxy.ctx.Done()
	wg.Wait()
	logging.LOGGER.Info("Forward proxy shutting down.")
	if internalServerError == nil && externalServerError == nil {
		return nil
	} else if internalServerError != nil {
		return internalServerError
	} else if externalServerError != nil {
		return externalServerError
	} else {
		return errors.Wrapf(internalServerError, externalServerError.Error())
	}
}

// Stop stops both the Alation facing and internet facing gRPC servers.
func (proxy *Proxy) Stop() {
	proxy.internal.Stop()
	proxy.external.Stop()
	proxy.stop()
}

// InternalAddress returns the "<addr>:<port>" string for the Alation facing server.
//
// This is mostly useful for unit testing.
func (proxy Proxy) InternalAddress() string {
	return fmt.Sprintf("%s:%d", proxy.internalAddr, proxy.internalPort)
}

// ExternalAddress returns the "<addr>:<port>" string for the internet facing server.
//
// This is mostly useful for unit testing.
func (proxy Proxy) ExternalAddress() string {
	return fmt.Sprintf("%s:%d", proxy.ExternalAddr, proxy.ExternalPort)
}

type GrpcCall struct {
	Alation   protocol.ServerStreamWithError
	Method    string
	Connector uint64
	Agent     uint64
	JobId     uint64
}

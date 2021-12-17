package reverse

import (
	"context"
	"fmt"
	"sync"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/shared"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/logging"
	"github.com/Alation/alation_connector_manager/docker/remoteAgent/protocol"
	"github.com/cenkalti/backoff/v4"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const ConnectorBasePort = 11000

type Agent struct {
	host   string
	port   uint16
	Id     uint64
	ctx    context.Context
	cancel context.CancelFunc
	client protocol.GrpcInverterClient

	lock     sync.Mutex
	connPool map[uint64]*grpc.ClientConn
}

func NewAgent(id uint64, host string, port uint16) *Agent {
	ctx, cancel := context.WithCancel(context.Background())
	agent := &Agent{
		host:     host,
		port:     port,
		Id:       id,
		ctx:      ctx,
		cancel:   cancel,
		lock:     sync.Mutex{},
		connPool: make(map[uint64]*grpc.ClientConn),
	}
	return agent
}

func (a *Agent) EventLoop() {
	logging.LOGGER.Info("Beginning agent event loop")
Init:
	a.NewUpstreamClient()
	jobs, err := a.client.JobTunnel(shared.NewHeaderBuilder().SetAgentId(a.Id).Build(a.ctx), &empty.Empty{})
	if err != nil {
		logging.LOGGER.Error("Client connection was established, however connecting to the job tunnel failed. Reconnected will be attempted.", zap.Error(err))
		goto Init
	}
	for {
		job, err := jobs.Recv()
		if err != nil {
			if a.Cancelled() {
				logging.LOGGER.Info("Agent shutting down")
				return
			}
			logging.LOGGER.Error("The job tunnel appears to have been shutdown. Reconnects will be attempted.", zap.Error(err))
			goto Init
		}
		go a.Dispatch(job)
	}
}

func (a *Agent) Dispatch(job *protocol.Job) {
	//////////////////////////////////////////////
	// Establish the callback stream.
	logging.LOGGER.Info("Received job, attempting to establish callback stream",
		logging.Method(job.Method),
		logging.Connector(job.Connector),
		logging.Job(job.JobID))
	upstream, err := a.client.Pipe(shared.NewHeaderBuilder().
		SetConnectorId(job.Connector).
		SetAgentId(a.Id).
		SetJobId(job.JobID).
		Build(context.Background()))
	if err != nil {
		panic(err)
	}
	logging.LOGGER.Info("Established callback stream",
		logging.Method(job.Method),
		logging.Connector(job.Connector),
		logging.Job(job.JobID),
		logging.PeerProxy(upstream.Context()))
	//////////////////////////////////////////////
	//////////////////////////////////////////////
	// Retrieve or create the connection to the connector.
	logging.LOGGER.Info("Attempting to establish connection to downstream connector",
		logging.Method(job.Method),
		logging.Connector(job.Connector),
		logging.Job(job.JobID))
	downstream, err := a.NewDownStreamClient(job.Connector, job.Method)
	if err != nil {
		e2 := upstream.Send(&protocol.Message{Error: ConnectorDown.Fmt(err, a.Id, job.Connector)})
		if e2 != nil {
			// @TODO
			logging.LOGGER.Error("bad and worse")
		}
		return
	}
	logging.LOGGER.Info("Connection established to downstream connector",
		logging.Method(job.Method),
		logging.Connector(job.Connector),
		logging.Job(job.JobID),
		logging.PeerConnector(downstream.Context()))
	//////////////////////////////////////////////
	//////////////////////////////////////////////
	// Begin proxying.
	logger := logging.LOGGER.With(logging.Method(job.Method),
		logging.Connector(job.Connector),
		logging.Job(job.JobID),
		logging.PeerProxy(upstream.Context()),
		logging.PeerConnector(downstream.Context()))
	protocol.ReverseProxy(upstream, downstream, logger)
}

func (a *Agent) NewDownStreamClient(connectorId uint64, method string) (grpc.ClientStream, error) {
	conn, err := func() (*grpc.ClientConn, error) {
		a.lock.Lock()
		defer a.lock.Unlock()
		conn, ok := a.connPool[connectorId]
		if ok {
			return conn, nil
		}
		conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%d", connectorId+ConnectorBasePort),
			grpc.WithKeepaliveParams(shared.KEEPALIVE_CLIENT_PARAMETERS),
			grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		a.connPool[connectorId] = conn
		return conn, nil
	}()
	if err != nil {
		return nil, err
	}
	return conn.NewStream(context.Background(), shared.BIDIRECTIONAL_STREAM_DESC, method, grpc.ForceCodec(shared.NoopCodec{}))
}

func (a *Agent) NewUpstreamClient() {
	a.lock.Lock()
	var conn *grpc.ClientConn
	target := fmt.Sprintf("%s:%d", a.host, a.port)
	// Returning a *backoff.PermanentError is the only way that backoff.Retry
	// will ever return an error. Alternatively, if the provided backoff
	// configuration has a non-zero MaxElapsedTime then backoff.Retry will
	// eventually return the error if it takes too long.
	//
	// We do neither of the above here because we have no interesting in giving up
	// the operation. This agent is entirely useless unless it can reach Alation.
	//
	// If you ever change this function to return a *backoff.PermanentError, or to
	// have a maximum elapsed time, then you will have to handle the possible error return.
	_ = backoff.Retry(func() error {
		var err error
		conn, err = grpc.Dial(target,
			grpc.WithInsecure(), // @TODO NOT INSECURE
			grpc.WithBlock(),
			grpc.WithKeepaliveParams(shared.KEEPALIVE_CLIENT_PARAMETERS))
		if err != nil {
			return err
		}
		return nil
	}, shared.RECONNECT_EXP_BACKOFF_CONFIG)
	a.client = protocol.NewGrpcInverterClient(conn)
	a.lock.Unlock()
}

// Stop stops the agent and halts the event loop. This used primarily for testing purposes.
func (a *Agent) Stop() {
	a.cancel()
}

// Cancelled returns whether-or-not the Agent has received a shutdown signal.
func (a *Agent) Cancelled() bool {
	return a.ctx.Err() != nil
}

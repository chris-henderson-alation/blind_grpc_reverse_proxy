package grpcinverter

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc/codes"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter/ioc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/cenkalti/backoff/v4"
)

const ConnectorBasePort = 11000

type Agent struct {
	host   string
	port   uint64
	id     uint64
	ctx    context.Context
	cancel context.CancelFunc
	client ioc.GrpcInverterClient

	lock     sync.Mutex
	connPool map[uint64]*grpc.ClientConn
}

func NewAgent(id uint64, host string, port uint64) *Agent {
	ctx, cancel := context.WithCancel(context.Background())
	agent := &Agent{
		host:     host,
		port:     port,
		id:       id,
		ctx:      ctx,
		cancel:   cancel,
		lock:     sync.Mutex{},
		connPool: make(map[uint64]*grpc.ClientConn),
	}
	return agent
}

func (a *Agent) EventLoop() {
Init:
	a.NewUpstreamClient()
	jobs, err := a.client.JobTunnel(NewHeaderBuilder().SetAgentId(a.id).Build(a.ctx), &empty.Empty{})
	if err != nil {
		logrus.Warn(err)
		goto Init
	}
	for {
		job, err := jobs.Recv()
		if err != nil {
			if a.Cancelled() {
				logrus.Info("Agent shutting down.")
				return
			}
			logrus.Warnf("Job tunnel shutdown with error %v, attempting reconnects.", err)
			goto Init
		}
		go a.Dispatch(job)
	}
}

func (a *Agent) Dispatch(job *ioc.Job) {
	logrus.Infof("Received job %v", job)
	//////////////////////////////////////////////
	// Establish the callback name.
	upstream, err := a.client.Pipe(NewHeaderBuilder().
		SetConnectorId(job.Connector).
		SetAgentId(a.id).
		SetJobId(job.JobID).
		Build(context.Background()))
	if err != nil {
		panic(err)
	}
	//////////////////////////////////////////////
	//////////////////////////////////////////////
	// Retrieve or create the connection to the connector.
	downstream, err := a.NewDownStreamClient(job.Connector, job.Method)
	if err != nil {
		e2 := upstream.Send(&ioc.Message{Error: a.connectorUnavailable(err)})
		if e2 != nil {
			logrus.Errorf("Original %v second %v", err, err)
		}
		return
	}
	//////////////////////////////////////////////
	//////////////////////////////////////////////
	// Begin proxying.
	ioc.ReverseProxy(upstream, downstream)
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
			grpc.WithKeepaliveParams(KEEPALIVE_CLIENT_PARAMETERS),
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
	return conn.NewStream(context.Background(), BIDIRECTIONAL_STREAM_DESC, method, grpc.ForceCodec(NoopCodec{}))
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
		logrus.Warnf("Attempting connection to %s", target)
		conn, err = grpc.Dial(target,
			grpc.WithInsecure(), // @TODO NOT INSECURE
			grpc.WithBlock(),
			grpc.WithKeepaliveParams(KEEPALIVE_CLIENT_PARAMETERS))
		if err != nil {
			logrus.Warnf("Connection to %s failed. This will be reattempted.", target)
			return err
		}
		return nil
	}, RECONNECT_EXP_BACKOFF_CONFIG)
	a.client = ioc.NewGrpcInverterClient(conn)
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

func (a *Agent) connectorUnavailable(err error) *ioc.Error {
	e := ioc.ErrorFromGoError(err)
	if e.Code == int32(codes.Unknown) {
		e.Code = int32(codes.Unavailable)
	}
	return e
}

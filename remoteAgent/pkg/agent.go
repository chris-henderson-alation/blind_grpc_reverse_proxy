package grpcinverter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter/ioc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

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
	//conn, err := grpc.Dial(fmt.Sprintf("fresh-crane.alation-test.com:%s", port),
	//	grpc.WithInsecure(),
	//	grpc.WithKeepaliveParams(shared.KEEPALIVE_CLIENT_PARAMETERS))
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(KEEPALIVE_CLIENT_PARAMETERS))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Agent{
		host:   host,
		port:   port,
		id:     id,
		ctx:    ctx,
		cancel: cancel,
		client: ioc.NewGrpcInverterClient(conn),

		lock:     sync.Mutex{},
		connPool: make(map[uint64]*grpc.ClientConn),
	}
}

func (a *Agent) Stop() {
	a.cancel()
}

func (a *Agent) EventLoop() {
	jobs, err := a.client.JobTunnel(NewHeaderBuilder().SetAgentId(a.id).Build(a.ctx), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			job, err := jobs.Recv()
			if err != nil {
				logrus.Warn(err)
				return
			}
			go a.Dispatch(job)
		}
	}()
}

func (a *Agent) Dispatch(job *ioc.Job) {
	logrus.Infof("Received job %v", job)
	headers := NewHeaderBuilder().
		SetConnectorId(job.Connector).
		SetAgentId(a.id).
		SetJobId(job.JobID).
		Build(context.Background())
	upstream, err := a.client.Pipe(headers)
	if err != nil {
		panic(err)
	}
	conn, err := a.fetchConnection(job.Connector)
	if err != nil {
		fmt.Println("booo tits")
		upstream.SendMsg(err)
		return
	}
	downstream, err := conn.NewStream(context.Background(), &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}, job.Method, grpc.ForceCodec(NoopCodec{}))
	if err != nil {
		fmt.Println("booo")
		upstream.SendMsg(err)
		return
	}
	var bookmark *ioc.Message
	var shouldRetry bool
	for {
		bookmark, shouldRetry = ioc.ReverseProxy(upstream, downstream, bookmark)
		if shouldRetry {
			logrus.Error("Failed! Trying again in a minute!")
			time.Sleep(time.Minute)
			panic("I dunno figure this out again in a bit")
			//c = NewClient()
			//upstream, err = c.Pipe(headers)
			if err != nil {
				panic("dunno, loops on loops")
			}
			logrus.Info("yoooo reconnected!")
		} else {
			return
		}
	}
}

func (a *Agent) fetchConnection(connectorId uint64) (*grpc.ClientConn, error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	conn, ok := a.connPool[connectorId]
	if ok {
		return conn, nil
	}
	var err error
	conn, err = grpc.Dial(fmt.Sprintf("0.0.0.0:%d", connectorId+10000),
		grpc.WithKeepaliveParams(KEEPALIVE_CLIENT_PARAMETERS),
		grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	a.connPool[connectorId] = conn
	return conn, nil
}

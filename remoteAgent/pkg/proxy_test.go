package grpcinverter

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter/ioc"

	"google.golang.org/grpc/keepalive"

	"google.golang.org/grpc"
)

var alationFacingPort uint64 = 20000
var agentFacingPort uint64 = 30000
var connectorId uint64 = 0
var agentId uint64 = 0

var host = "0.0.0.0"

const connectorBasePort = 10000

func TestPerfStream(t *testing.T) {
	forward := NewForwardProxyFacade()
	agent := NewAgent(atomic.AddUint64(&agentId, 1), host, forward.external)
	go agent.EventLoop()
	defer agent.Stop()
	defer forward.Stop()
	for {
		if forward.agents.Listening(agentId) {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
	connector := NewConnector(&PerformanceConnector{})
	connector.Start()
	defer connector.Stop()
	alation := NewAlationClient(forward.internal)
	perf, err := alation.Performance(NewHeaderBuilder().
		SetJobId(1).
		SetAgentId(agent.id).
		SetConnectorId(connector.id).
		Build(context.Background()),
		&String{Message: "don't really care"})
	if err != nil {
		t.Fatal(err)
	}
	for _, err := perf.Recv(); err == nil; _, err = perf.Recv() {
	}
}

func TestUnary(t *testing.T) {
	forward := NewForwardProxyFacade()
	agent := NewAgent(atomic.AddUint64(&agentId, 1), host, forward.external)
	go agent.EventLoop()
	defer agent.Stop()
	defer forward.Stop()
	for {
		if forward.agents.Listening(agentId) {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
	input := randomString()
	mutation := randomString()
	want := fmt.Sprintf("%s%s", input, mutation)
	connector := NewConnector(&UnaryConnector{mutation: mutation})
	connector.Start()
	defer connector.Stop()
	alation := NewAlationClient(forward.internal)
	unary, err := alation.Unary(
		NewHeaderBuilder().
			SetJobId(1).
			SetAgentId(agent.id).
			SetConnectorId(connector.id).
			Build(context.Background()),
		&String{Message: input})
	if err != nil {
		t.Fatal(err)
	}
	if unary.Message != want {
		t.Fatalf("wanted '%s' got '%s'", want, unary.Message)
	}

	input = randomString()
	want = fmt.Sprintf("%s%s", input, mutation)
	unary, err = alation.Unary(
		NewHeaderBuilder().
			SetJobId(1).
			SetAgentId(agent.id).
			SetConnectorId(connector.id).
			Build(context.Background()),
		&String{Message: input})
	if err != nil {
		t.Fatal(err)
	}
	if unary.Message != want {
		t.Fatalf("wanted '%s' got '%s'", want, unary.Message)
	}
}

type UnaryConnector struct {
	mutation string
	t        *testing.T
	UnimplementedTestServer
}

func (u *UnaryConnector) Unary(ctx context.Context, s *String) (*String, error) {
	return &String{Message: fmt.Sprintf("%s%s", s.Message, u.mutation)}, nil
}

type PerformanceConnector struct {
	UnimplementedTestServer
}

func (p *PerformanceConnector) Performance(s *String, server Test_PerformanceServer) error {
	log.Info("Alright, you asked for it")
	b := strings.Builder{}
	for i := b.Len(); i < 1024*50; i = b.Len() {
		b.WriteString(randomString())
	}
	length := b.Len()
	str := &String{Message: b.String()}
	gigabyte := 1024 * 1024 * 1024
	for i := gigabyte * 2; i > 0; i -= length {
		err := server.Send(str)
		if err != nil {
			fmt.Println("wut")
			return err
		}
	}
	return nil
}

type ForwardFacade struct {
	internal uint64
	external uint64
	i        *grpc.Server
	e        *grpc.Server
	agents   *Agents
}

func NewForwardProxyFacade() *ForwardFacade {
	forwardProxy := NewForwardProxy()
	internalServer := grpc.NewServer(
		grpc.KeepaliveParams(KEEPALIVE_SERVER_PARAMETERS),
		grpc.KeepaliveEnforcementPolicy(KEEPALIVE_ENFORCEMENT_POLICY),
		grpc.ForceServerCodec(NoopCodec{}),
		grpc.UnknownServiceHandler(forwardProxy.StreamHandler()))
	externalServer := grpc.NewServer(
		grpc.KeepaliveParams(KEEPALIVE_SERVER_PARAMETERS),
		grpc.KeepaliveEnforcementPolicy(KEEPALIVE_ENFORCEMENT_POLICY))
	ioc.RegisterGrpcInverterServer(externalServer, forwardProxy)
	internal := atomic.AddUint64(&alationFacingPort, 1)
	internalListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", internal))
	if err != nil {
		panic(err)
	}
	external := atomic.AddUint64(&agentFacingPort, 1)
	externalListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", external))
	if err != nil {
		panic(err)
	}
	go func() {
		internalServer.Serve(internalListener)
	}()
	go func() {
		externalServer.Serve(externalListener)
	}()
	return &ForwardFacade{
		internal: internal,
		external: external,
		i:        internalServer,
		e:        externalServer,
		agents:   forwardProxy.agents,
	}
}

func (f *ForwardFacade) Stop() {
	f.i.Stop()
	f.e.Stop()
}

func NewAlationClient(target uint64) TestClient {
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%d", target),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 30,
			PermitWithoutStream: true,
		}))
	if err != nil {
		panic(err)
	}
	return NewTestClient(conn)
}

func NewConnector(connector TestServer) *MockConnector {
	id := atomic.AddUint64(&connectorId, 1)
	return &MockConnector{id: id, server: nil, TestServer: connector}
}

type MockConnector struct {
	id     uint64
	server *grpc.Server
	TestServer
}

func (c *MockConnector) Stop() {
	c.server.Stop()
}

func (c *MockConnector) Start() {
	s := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		Time: time.Second * 30,
	}), grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             time.Second * 10,
		PermitWithoutStream: true,
	}))
	RegisterTestServer(s, c)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", connectorBasePort+c.id))
	if err != nil {
		panic(err)
	}
	c.server = s
	go c.server.Serve(listener)
}

//func (c *MockConnector) Unary(ctx context.Context, s *String) (*String, error) {
//	response := randomString()
//	log.Infof("Got Unary request: '%s' adding '%s'", s.Message, response)
//	str := &String{Message: fmt.Sprintf("%s %s", s.Message, response)}
//	return str, nil
//}
//
//func (c *MockConnector) ServerStream(s *String, server Test_ServerStreamServer) error {
//	log.Infof("Got ServerStream request: '%s'", s.Message)
//	for i := 0; i < 5; i++ {
//		time.Sleep(time.Millisecond * time.Duration(rand.Int63n(2000)))
//		str := &String{Message: fmt.Sprintf("%s %s", s.Message, randomString())}
//		log.Infof("Sending '%s'", str.Message)
//		err := server.Send(str)
//		if err != nil {
//			log.Errorf("ServerStream error: %v", err)
//			return err
//		}
//	}
//	return nil
//}
//
//func (c *MockConnector) ClientStream(server Test_ClientStreamServer) error {
//	log.Info("Got ClientStream request")
//	for {
//		str := &String{}
//		err := server.RecvMsg(str)
//		switch err {
//		case nil:
//			log.Infof("Got response '%s'", str.Message)
//		case io.EOF:
//			log.Info("Shutdown from client")
//			return nil
//		case io.ErrClosedPipe:
//			log.Info("Closed pipe from client")
//			return nil
//		default:
//			log.Infof("Some other error from client: %v", err)
//			return err
//		}
//	}
//}
//
//func (c *MockConnector) Bidirectional(server Test_BidirectionalServer) error {
//	for i := 1; i <= 2; i++ {
//		message, err := server.Recv()
//		if err != nil {
//			log.Errorf("BiDirection recv %d err %v", i, err)
//			return err
//		}
//		response := randomString()
//		log.Infof("Got Bidirectional request: '%s' adding '%s'", message.Message, response)
//		str := &String{Message: fmt.Sprintf("%s %s", message.Message, response)}
//		err = server.Send(str)
//		if err != nil {
//			log.Errorf("BiDirection send %d err %v", i, err)
//			return err
//		}
//	}
//
//	return nil
//}
//
//func (c *MockConnector) IntentionalError(s *String, server Test_IntentionalErrorServer) error {
//	return status.New(codes.DataLoss, "THIS IS INTENTIONALLY BAD").Err()
//}
//
//func (c *MockConnector) Performance(s *String, server Test_PerformanceServer) error {
//	log.Info("Alright, you asked for it")
//	b := strings.Builder{}
//	for i := b.Len(); i < 1024*50; i = b.Len() {
//		b.WriteString(randomString())
//	}
//	length := b.Len()
//	str := &String{Message: b.String()}
//	gigabyte := 1024 * 1024 * 1024
//	for i := gigabyte * 2; i > 0; i -= length {
//		err := server.Send(str)
//		if err != nil {
//			return err
//		}
//	}
//	log.Info("I mean...I'm done...")
//	return nil
//}

func randomString() string {
	length := rand.Intn(16)
	s := strings.Builder{}
	for i := 0; i < length; i++ {
		s.WriteByte(byte(rand.Intn(42) + 48))
	}
	return s.String()
}

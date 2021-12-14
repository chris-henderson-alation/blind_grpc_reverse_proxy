package grpcinverter

import (
	"context"
	// We use both math and crypto rand
	crand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/Alation/alation_connector_manager/docker/remoteAgent/grpcinverter/ioc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var alationFacingPort uint64 = 20000
var agentFacingPort uint64 = 30000
var connectorId uint64 = 0
var agentId uint64 = 0

var host = "0.0.0.0"

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type UnaryConnector struct {
	mutation string
	t        *testing.T
	UnimplementedTestServer
}

func (u *UnaryConnector) Unary(ctx context.Context, s *String) (*String, error) {
	return &String{Message: fmt.Sprintf("%s%s", s.Message, u.mutation)}, nil
}

func TestUnary(t *testing.T) {
	input := randomString()
	mutation := randomString()
	want := fmt.Sprintf("%s%s", input, mutation)
	stack := newStack(&UnaryConnector{mutation: mutation})
	defer stack.Stop()
	alation, agent, connector := stack.alation, stack.agent, stack.connector

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
}

func TestParallelUnary(t *testing.T) {
	t.Parallel()
	type data struct {
		alationGives  string
		connectorAdds string
		want          string
	}
	mutation := randomString()
	testData := make([]data, 100)
	for i := 0; i < 100; i++ {
		alationGives := randomString()
		testData[i] = data{
			alationGives:  alationGives,
			connectorAdds: mutation,
			want:          fmt.Sprintf("%s%s", alationGives, mutation),
		}
	}
	stack := newStack(&UnaryConnector{mutation: mutation})
	defer stack.Stop()
	for _, test := range testData {
		t.Run(test.alationGives, func(t2 *testing.T) {
			got, err := stack.alation.Unary(stack.Headers(context.Background()), &String{Message: test.alationGives})
			if err != nil {
				t2.Fatal(err)
			}
			if got.Message != test.want {
				t2.Errorf("got '%s' want '%s'", got.Message, test.want)
			}
		})
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ComplexTypeConnector struct {
	want   *ComplexPbType
	output *ComplexPbType
	UnimplementedTestServer
}

func (c *ComplexTypeConnector) ComplexType(ctx context.Context, input *ComplexPbType) (*ComplexPbType, error) {
	if !proto.Equal(c.want, input) {
		return nil, fmt.Errorf("wanted %v got %v", c.want, input)
	}
	return c.output, nil
}

func TestComplexPbType(t *testing.T) {
	_struct, err := structpb.NewStruct(map[string]interface{}{
		"hello": "world",
		"neat":  1,
		"cool":  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	nested := &ComplexPbType_Nested{
		Value: 42,
	}
	input := &ComplexPbType{
		Scalar:      42,
		Array:       []uint64{5, 6, 7, 8},
		Enumeration: ComplexPbType_VariantC,
		Struct:      _struct,
		Nested:      nested,
		Oneof:       &ComplexPbType_OneOfVariantB{},
		Map: map[int32]*ComplexPbType_Nested{
			1:  {Value: 1},
			2:  {Value: 2},
			3:  {Value: 3},
			42: nested,
		},
	}
	_struct2, err := structpb.NewStruct(map[string]interface{}{
		"the best":            "song in the world",
		"different data type": true,
		"cool but different":  2,
	})
	nested2 := &ComplexPbType_Nested{
		Value: 86,
	}
	output := &ComplexPbType{
		Scalar:      42,
		Array:       []uint64{5, 6, 7, 8},
		Enumeration: ComplexPbType_VariantC,
		Struct:      _struct2,
		Nested:      nested2,
		Oneof:       &ComplexPbType_OneOfVariantB{},
		Map: map[int32]*ComplexPbType_Nested{
			4:  {Value: 4},
			5:  {Value: 5},
			6:  {Value: 6},
			86: nested2,
		},
	}
	stack := newStack(&ComplexTypeConnector{want: input, output: output})
	defer stack.Stop()
	got, err := stack.alation.ComplexType(stack.Headers(context.Background()), input)
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(output, got) {
		t.Fatalf("got %v got %v", got, output)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ClientStreamConnector struct {
	want []string
	UnimplementedTestServer
}

func (c *ClientStreamConnector) ClientStream(server Test_ClientStreamServer) error {
	found := 0
	for {
		got, err := server.Recv()
		switch err {
		case nil:
			want := c.want[found]
			if got.Message != want {
				return fmt.Errorf("wanted '%s' got '%s'", want, got.Message)
			}
			found += 1
		case io.EOF:
			if found != len(c.want) {
				return fmt.Errorf("wanted %d string but got %d", len(c.want), found)
			}
		default:
			return err
		}
	}
}

func TestClientStream(t *testing.T) {
	input := make([]string, 100)
	for i := 0; i < 100; i++ {
		input[i] = randomString()
	}
	stack := newStack(&ClientStreamConnector{want: input})
	defer stack.Stop()
	stream, err := stack.alation.ClientStream(stack.Headers(context.Background()))
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range input {
		err := stream.Send(&String{Message: i})
		if err != nil {
			t.Fatal(err)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type IntentionalErrorConnector struct {
	Error *status.Status
	UnimplementedTestServer
}

func (i *IntentionalErrorConnector) IntentionalError(s *String, server Test_IntentionalErrorServer) error {
	return i.Error.Err()
}

func TestIntentionalError(t *testing.T) {
	want, err := status.New(codes.Aborted, "this isn't the greatest song in the world").WithDetails(&String{Message: "here is some extra context"})
	if err != nil {
		t.Fatal(err)
	}
	stack := newStack(&IntentionalErrorConnector{Error: want})
	defer stack.Stop()
	stream, err := stack.alation.IntentionalError(
		stack.Headers(context.Background()),
		&String{Message: "boy, I sure do hope that the connector doesn't crash or nothing 😇"})
	if err != nil {
		// this is not the error we are looking for.
		t.Fatal(err)
	}
	s, got := stream.Recv()
	if got == nil {
		t.Fatalf("wanted %v as an error but got nil, also got %v as an actual return", want, s)
	}
	protoGot, _ := status.FromError(got)
	if !proto.Equal(protoGot.Proto(), want.Proto()) {
		t.Fatalf("wanted '%v' got '%v'", want, protoGot)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type BidirectionalStreamConnectorProtocolDriven struct {
	UnimplementedTestServer
}

func (b BidirectionalStreamConnectorProtocolDriven) Bidirectional(server Test_BidirectionalServer) error {
	builder := strings.Builder{}
	for {
		msg, err := server.Recv()
		switch err {
		case nil:
			builder.WriteString(msg.Message)
		case io.EOF:
			return nil
		default:
			panic(err)
		}
		msg, err = server.Recv()
		switch err {
		case nil:
			builder.WriteString(msg.Message)
		case io.EOF:
			return nil
		default:
			return err
		}
		err = server.Send(&String{Message: builder.String()})
		if err != nil {
			panic(err)
		}
		builder.Reset()
	}
}

func TestBidirectionalStream(t *testing.T) {
	stack := newStack(&BidirectionalStreamConnectorProtocolDriven{})
	defer stack.Stop()
	type data struct {
		one  *String
		two  *String
		want *String
	}
	testData := make([]data, 100)
	for i := 0; i < 100; i++ {
		one := randomString()
		two := randomString()
		combined := fmt.Sprintf("%s%s", one, two)
		testData[i] = data{
			one:  &String{Message: one},
			two:  &String{Message: two},
			want: &String{Message: combined},
		}
	}
	stream, err := stack.alation.Bidirectional(stack.Headers(context.Background()))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range testData {
		err = stream.Send(test.one)
		if err != nil {
			t.Fatal(err)
		}
		err = stream.Send(test.two)
		if err != nil {
			t.Fatal(err)
		}
		got, err := stream.Recv()
		if err != nil {
			t.Fatal(err)
		}
		if !proto.Equal(got, test.want) {
			t.Fatalf("wanted %v got %v", test.want, got)
		}
	}
	stream.CloseSend()
	// Receive in order to get the EOF from the server or else everything shuts down too fast
	// and errors happen as transports begin to shut down.
	if _, err = stream.Recv(); err != io.EOF {
		t.Fatal(err)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Test Header Passthrough

type HeaderConnector struct {
	want metadata.MD
	UnimplementedTestServer
}

func (h *HeaderConnector) Unary(ctx context.Context, s *String) (*String, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no headers")
	}
	if !reflect.DeepEqual(md, h.want) {
		return nil, fmt.Errorf("wanted %v got %v", h.want, md)
	}
	return nil, nil
}

func TestHeaders(t *testing.T) {
	t.SkipNow()
	connector := &HeaderConnector{}
	stack := newStack(connector)
	defer stack.Stop()
	builder := NewHeaderBuilder().
		SetJobId(1).
		SetAgentId(stack.agent.id).
		SetConnectorId(stack.connector.id)
	md := metadata.New(builder.headers)
	connector.want = md
	_, err := stack.alation.Unary(metadata.NewOutgoingContext(context.Background(), md), &String{Message: "cool headers"})
	if err != nil {
		t.Fatal(err)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Test Reconnect

type ReconnectServerStream struct {
	UnimplementedTestServer
}

func (r ReconnectServerStream) ServerStream(s *String, server Test_ServerStreamServer) error {
	for i := 0; i < 20; i++ {
		sent := fmt.Sprintf("%d", i)
		err := server.Send(&String{Message: sent})
		if err != nil {
			return err
		}
		fmt.Println(sent)
		time.Sleep(time.Second)
	}
	return nil
}

func TestReconnect(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Log("This test requires escalated privileges in order to manipulate the network")
		t.SkipNow()
	}
	link := NewLocalLink(t.Name())
	defer link.Destroy()
	stack := newCustomNetworkingStack(&ReconnectServerStream{}, nil, &link.addr, nil)
	defer stack.Stop()
	stream, err := stack.alation.ServerStream(stack.Headers(context.Background()), &String{Message: ""})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		want := fmt.Sprintf("%d", i)
		s, err := stream.Recv()
		if err != nil {
			t.Fatal(err)
		}
		if s.Message != want {
			t.Fatalf("wanted '%s' got '%s'", want, s.Message)
		}
	}
	link.Down()
	t.Log("sleeping!")
	time.Sleep(time.Second * 30)
	t.Log("finihsings")
	link.Up()
	for i := 10; i < 20; i++ {
		want := fmt.Sprintf("%d", i)
		s, err := stream.Recv()
		if err != nil {
			t.Fatal(err)
		}
		if s.Message != want {
			t.Fatalf("wanted '%s' got '%s'", want, s.Message)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Test Reconnect

type BrokenTunnelTestConnector struct {
	UnimplementedTestServer
}

func (r *BrokenTunnelTestConnector) ServerStream(s *String, server Test_ServerStreamServer) error {
	for i := 0; i < 20; i++ {
		err := server.Send(&String{Message: fmt.Sprintf("%d", i)})
		if err != nil {
			return err
		}
	}
	return nil
}

func TestBrokenTunnel(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Log("This test requires escalated privileges in order to manipulate the network")
		t.SkipNow()
	}
	t.Log("what is all this then come now")
	link := NewLocalLink(t.Name())
	defer link.Destroy()
	stack := newCustomNetworkingStack(&BrokenTunnelTestConnector{}, nil, &link.addr, nil)
	defer stack.Stop()
	link.Down()
	time.Sleep(time.Second * 60)
	link.Up()
	iters := 0
	for {
		if iters > 100 {
			t.Fatal("not in time")
		}
		iters += 1
		if stack.forward.agents.Listening(stack.agent.id) {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	t.Log("yaaaas")
	stream, err := stack.alation.ServerStream(stack.Headers(context.Background()), &String{Message: ""})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("yus")
	for i := 0; i < 20; i++ {
		want := fmt.Sprintf("%d", i)
		s, err := stream.Recv()
		if err != nil {
			t.Fatal(err)
		}
		if s.Message != want {
			t.Fatalf("wanted '%s' got '%s'", want, s.Message)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Test utility functions.

type TechStack struct {
	alation   TestClient
	forward   *ForwardFacade
	agent     *Agent
	connector *MockConnector
}

func (t *TechStack) Stop() {
	t.connector.Stop()
	t.agent.Stop()
	t.forward.Stop()
}

func (t *TechStack) Headers(ctx context.Context) context.Context {
	return NewHeaderBuilder().
		SetJobId(1).
		SetAgentId(t.agent.id).
		SetConnectorId(t.connector.id).
		Build(ctx)
}

func newStack(connector TestServer) *TechStack {
	return newCustomNetworkingStack(connector, nil, nil, nil)
}
func newCustomNetworkingStack(connector TestServer, alationToForwardAddr, forwardToAgentAddr, agentToConnectorAddr *string) *TechStack {
	if alationToForwardAddr == nil {
		alationToForwardAddr = &host
	}
	if forwardToAgentAddr == nil {
		forwardToAgentAddr = &host
	}
	if agentToConnectorAddr == nil {
		agentToConnectorAddr = &host
	}
	forward := NewForwardProxyFacade(*alationToForwardAddr, *forwardToAgentAddr)
	agent := NewAgent(atomic.AddUint64(&agentId, 1), *forwardToAgentAddr, forward.external)
	go agent.EventLoop()
	for {
		if forward.agents.Listening(agentId) {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	connectorServer := NewConnector(connector, *agentToConnectorAddr)
	connectorServer.Start()
	alation := NewAlationClient(forward.internal)
	return &TechStack{
		alation:   alation,
		forward:   forward,
		agent:     agent,
		connector: connectorServer,
	}
}

type ForwardFacade struct {
	internal uint64
	external uint64
	i        *grpc.Server
	e        *grpc.Server
	agents   *Agents
}

func NewForwardProxyFacade(internalAddr, externalAddr string) *ForwardFacade {
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
	internalListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", internalAddr, internal))
	if err != nil {
		panic(err)
	}
	external := atomic.AddUint64(&agentFacingPort, 1)
	externalListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", externalAddr, external))
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

func NewConnector(connector TestServer, addr string) *MockConnector {
	id := atomic.AddUint64(&connectorId, 1)
	s := grpc.NewServer(
		grpc.KeepaliveParams(KEEPALIVE_SERVER_PARAMETERS),
		grpc.KeepaliveEnforcementPolicy(KEEPALIVE_ENFORCEMENT_POLICY))
	RegisterTestServer(s, connector)
	return &MockConnector{id: id, server: s, TestServer: connector, addr: addr}
}

type MockConnector struct {
	id     uint64
	server *grpc.Server
	addr   string
	TestServer
}

func (c *MockConnector) Stop() {
	c.server.Stop()
}

func (c *MockConnector) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.addr, ConnectorBasePort+c.id))
	if err != nil {
		panic(err)
	}
	go c.server.Serve(listener)
}

func randomString() string {
	length := rand.Intn(16)
	s := strings.Builder{}
	for i := 0; i < length; i++ {
		s.WriteByte(byte(rand.Intn(42) + 48))
	}
	return s.String()
}

// Shameless copying in order to figure out how to have mulitple dummy links in order
// to simulate down networks.
//
// https://linuxconfig.org/configuring-virtual-network-interfaces-in-linux
type LocalLink struct {
	addr     string
	addrCIDR string
	mac      string
	name     string
	label    string
}

var nextAddrSpace uint32 = 0

func NewLocalLink(name string) LocalLink {
	next := atomic.AddUint32(&nextAddrSpace, 1)
	name = fmt.Sprintf("ocf%d", next)
	addr := fmt.Sprintf("10.0.%d.0", next)
	addrCIDR := fmt.Sprintf("%s/24", addr)
	macBytes := make([]byte, 6)
	n, err := crand.Read(macBytes)
	if err != nil {
		panic(err)
	}
	if n != len(macBytes) {
		panic("wanted 6 bytes read into random MAC address")
	}
	mac := fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X",
		macBytes[0],
		macBytes[1],
		macBytes[2],
		macBytes[3],
		macBytes[4],
		macBytes[5])
	link := LocalLink{
		addr:     addr,
		addrCIDR: addrCIDR,
		mac:      mac,
		name:     name,
		label:    fmt.Sprintf("%s:0", name),
	}
	link.Destroy()
	if err := exec.Command("modprobe", "dummy").Run(); err != nil {
		panic(err)
	}
	if err := exec.Command("ip", "link", "add", link.name, "type", "dummy").Run(); err != nil {
		panic(err)
	}
	if err := exec.Command("ifconfig", link.name, "hw", "ether", link.mac).Run(); err != nil {
		panic(err)
	}
	if err := exec.Command("ip", "addr", "add", link.addrCIDR, "brd", "+", "dev", link.name, "label", link.label).Run(); err != nil {
		panic(err)
	}
	link.Up()
	return link
}

func (l LocalLink) Up() {
	if err := exec.Command("ip", "link", "set", "dev", l.name, "up").Run(); err != nil {
		panic(err)
	}
}

func (l LocalLink) Down() {
	if err := exec.Command("ip", "link", "set", "dev", l.name, "down").Run(); err != nil {
		panic(err)
	}
}

func (l LocalLink) Destroy() {
	exec.Command("ip", "addr", "del", l.addrCIDR, "brd", "+", "dev", l.name, "label", l.label).Run()
	exec.Command("ip", "link", "delete", l.name, "type", "dummy").Run()
}

package remoteAgent

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Alation/alation_connector_manager/remoteAgent/forward"
	"github.com/Alation/alation_connector_manager/remoteAgent/reverse"
	"github.com/Alation/alation_connector_manager/remoteAgent/shared"
	"github.com/cenkalti/backoff/v4"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

var alationFacingPort uint32 = 20000
var agentFacingPort uint32 = 30000
var connectorId uint64 = 0
var agentId uint64 = 0

var loopback = "0.0.0.0"

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
		shared.NewHeaderBuilder().
			SetJobId(1).
			SetAgentId(agent.Id).
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
	for i, test := range testData {
		t.Run(test.alationGives, func(t2 *testing.T) {
			got, err := stack.alation.Unary(stack.HeadersWithJobId(context.Background(), uint64(i+1)), &String{Message: test.alationGives})
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
		&String{Message: "boy, I sure do hope that the connector doesn't crash or nothing ????"})
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
// Test multiple agents.

func TestMultipleAgents(t *testing.T) {
	mutation := randomString()
	stack := newStack(&UnaryConnector{mutation: mutation})
	defer stack.Stop()
	type data struct {
		alationGives  string
		connectorAdds string
		want          string
	}
	numAgents := 50
	// Test for 100 agents connected at once
	agents := make([]*reverse.Agent, numAgents)
	for i := 0; i < numAgents; i++ {
		agent := newAgent(stack.forward)
		defer agent.Stop()
		agents[i] = agent
	}
	var jobId uint64 = 0
	for i := 0; i < numAgents; i++ {
		agent := agents[i]
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			for j := 0; j < 100; j++ {
				alationGives := randomString()
				test := data{
					alationGives:  alationGives,
					connectorAdds: mutation,
					want:          fmt.Sprintf("%s%s", alationGives, mutation),
				}
				got, err := stack.alation.Unary(shared.NewHeaderBuilder().
					SetJobId(atomic.AddUint64(&jobId, 1)).
					SetAgentId(agent.Id).
					SetConnectorId(stack.connector.id).
					Build(context.Background()), &String{Message: test.alationGives})
				if err != nil {
					t.Fatal(err)
				}
				if got.Message != test.want {
					t.Errorf("wanted '%s' got '%s'", test.want, got.Message)
				}
			}
		})
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Test when the connector is down.

func TestConnectorDown(t *testing.T) {
	stack := newStack(&UnimplementedTestServer{})
	defer stack.Stop()
	stack.connector.Stop()
	resp, err := stack.alation.Unary(stack.Headers(context.Background()), &String{Message: "Nothing, please."})
	t.Log(resp)
	t.Log(err)
	s, _ := status.FromError(err)
	t.Log(s.Details())
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
	builder := shared.NewHeaderBuilder().
		SetJobId(1).
		SetAgentId(stack.agent.Id).
		SetConnectorId(stack.connector.id)
	md := metadata.New(builder.Headers)
	connector.want = md
	_, err := stack.alation.Unary(metadata.NewOutgoingContext(context.Background(), md), &String{Message: "cool headers"})
	if err != nil {
		t.Fatal(err)
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
	forward   *forward.Proxy
	agent     *reverse.Agent
	connector *MockConnector
}

func (t *TechStack) Stop() {
	t.connector.Stop()
	t.agent.Stop()
	t.forward.Stop()
}

func (t *TechStack) Headers(ctx context.Context) context.Context {
	return shared.NewHeaderBuilder().
		SetJobId(1).
		SetAgentId(t.agent.Id).
		SetConnectorId(t.connector.id).
		Build(ctx)
}

func (t *TechStack) HeadersWithJobId(ctx context.Context, jobId uint64) context.Context {
	return shared.NewHeaderBuilder().
		SetJobId(jobId).
		SetAgentId(t.agent.Id).
		SetConnectorId(t.connector.id).
		Build(ctx)
}

func newStack(connector TestServer) *TechStack {
	return newCustomNetworkingStack(connector, nil, nil, nil)
}

func newCustomNetworkingStack(connector TestServer, alationToForwardAddr, forwardToAgentAddr, agentToConnectorAddr *string) *TechStack {
	if alationToForwardAddr == nil {
		alationToForwardAddr = &loopback
	}
	if forwardToAgentAddr == nil {
		forwardToAgentAddr = &loopback
	}
	if agentToConnectorAddr == nil {
		agentToConnectorAddr = &loopback
	}
	forward := forward.NewProxy(
		*alationToForwardAddr, uint16(atomic.AddUint32(&alationFacingPort, 1)),
		*forwardToAgentAddr, uint16(atomic.AddUint32(&agentFacingPort, 1)))
	go func() {
		if err := forward.Start(); err != nil {
			panic(err)
		}
	}()
	agent := newAgent(forward)
	connectorServer := NewConnector(connector, *agentToConnectorAddr)
	connectorServer.Start()
	alation := NewAlationClient(forward.InternalAddress())
	return &TechStack{
		alation:   alation,
		forward:   forward,
		agent:     agent,
		connector: connectorServer,
	}
}

func newAgent(target *forward.Proxy) *reverse.Agent {
	agent := reverse.NewAgent(atomic.AddUint64(&agentId, 1), target.ExternalAddr, target.ExternalPort)
	go func() {
		agent.EventLoop()
	}()
	err := backoff.Retry(func() error {
		if target.Listeners.Listening(agentId) {
			return nil
		}
		return fmt.Errorf("not yet connected")
	}, &backoff.ExponentialBackOff{
		InitialInterval:     time.Millisecond * 100,
		RandomizationFactor: 0.2,
		Multiplier:          1.6,
		MaxInterval:         time.Second,
		MaxElapsedTime:      time.Second * 2,
		Clock:               backoff.SystemClock,
	})
	if err != nil {
		panic(err)
	}
	return agent
}

func NewAlationClient(target string) TestClient {
	conn, err := grpc.Dial(target,
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
		grpc.KeepaliveParams(shared.KEEPALIVE_SERVER_PARAMETERS),
		grpc.KeepaliveEnforcementPolicy(shared.KEEPALIVE_ENFORCEMENT_POLICY))
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
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.addr, reverse.ConnectorBasePort+c.id))
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

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpcinverter

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TestClient is the client API for Test service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TestClient interface {
	Unary(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	ServerStream(ctx context.Context, in *String, opts ...grpc.CallOption) (Test_ServerStreamClient, error)
	ClientStream(ctx context.Context, opts ...grpc.CallOption) (Test_ClientStreamClient, error)
	Bidirectional(ctx context.Context, opts ...grpc.CallOption) (Test_BidirectionalClient, error)
	IntentionalError(ctx context.Context, in *String, opts ...grpc.CallOption) (Test_IntentionalErrorClient, error)
	Performance(ctx context.Context, in *String, opts ...grpc.CallOption) (Test_PerformanceClient, error)
	PerformanceBytes(ctx context.Context, in *TestBytes, opts ...grpc.CallOption) (Test_PerformanceBytesClient, error)
}

type testClient struct {
	cc grpc.ClientConnInterface
}

func NewTestClient(cc grpc.ClientConnInterface) TestClient {
	return &testClient{cc}
}

func (c *testClient) Unary(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := c.cc.Invoke(ctx, "/test.Test/Unary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testClient) ServerStream(ctx context.Context, in *String, opts ...grpc.CallOption) (Test_ServerStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Test_ServiceDesc.Streams[0], "/test.Test/ServerStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &testServerStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Test_ServerStreamClient interface {
	Recv() (*String, error)
	grpc.ClientStream
}

type testServerStreamClient struct {
	grpc.ClientStream
}

func (x *testServerStreamClient) Recv() (*String, error) {
	m := new(String)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *testClient) ClientStream(ctx context.Context, opts ...grpc.CallOption) (Test_ClientStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Test_ServiceDesc.Streams[1], "/test.Test/ClientStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &testClientStreamClient{stream}
	return x, nil
}

type Test_ClientStreamClient interface {
	Send(*String) error
	CloseAndRecv() (*String, error)
	grpc.ClientStream
}

type testClientStreamClient struct {
	grpc.ClientStream
}

func (x *testClientStreamClient) Send(m *String) error {
	return x.ClientStream.SendMsg(m)
}

func (x *testClientStreamClient) CloseAndRecv() (*String, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(String)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *testClient) Bidirectional(ctx context.Context, opts ...grpc.CallOption) (Test_BidirectionalClient, error) {
	stream, err := c.cc.NewStream(ctx, &Test_ServiceDesc.Streams[2], "/test.Test/Bidirectional", opts...)
	if err != nil {
		return nil, err
	}
	x := &testBidirectionalClient{stream}
	return x, nil
}

type Test_BidirectionalClient interface {
	Send(*String) error
	Recv() (*String, error)
	grpc.ClientStream
}

type testBidirectionalClient struct {
	grpc.ClientStream
}

func (x *testBidirectionalClient) Send(m *String) error {
	return x.ClientStream.SendMsg(m)
}

func (x *testBidirectionalClient) Recv() (*String, error) {
	m := new(String)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *testClient) IntentionalError(ctx context.Context, in *String, opts ...grpc.CallOption) (Test_IntentionalErrorClient, error) {
	stream, err := c.cc.NewStream(ctx, &Test_ServiceDesc.Streams[3], "/test.Test/IntentionalError", opts...)
	if err != nil {
		return nil, err
	}
	x := &testIntentionalErrorClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Test_IntentionalErrorClient interface {
	Recv() (*String, error)
	grpc.ClientStream
}

type testIntentionalErrorClient struct {
	grpc.ClientStream
}

func (x *testIntentionalErrorClient) Recv() (*String, error) {
	m := new(String)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *testClient) Performance(ctx context.Context, in *String, opts ...grpc.CallOption) (Test_PerformanceClient, error) {
	stream, err := c.cc.NewStream(ctx, &Test_ServiceDesc.Streams[4], "/test.Test/Performance", opts...)
	if err != nil {
		return nil, err
	}
	x := &testPerformanceClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Test_PerformanceClient interface {
	Recv() (*String, error)
	grpc.ClientStream
}

type testPerformanceClient struct {
	grpc.ClientStream
}

func (x *testPerformanceClient) Recv() (*String, error) {
	m := new(String)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *testClient) PerformanceBytes(ctx context.Context, in *TestBytes, opts ...grpc.CallOption) (Test_PerformanceBytesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Test_ServiceDesc.Streams[5], "/test.Test/PerformanceBytes", opts...)
	if err != nil {
		return nil, err
	}
	x := &testPerformanceBytesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Test_PerformanceBytesClient interface {
	Recv() (*TestBytes, error)
	grpc.ClientStream
}

type testPerformanceBytesClient struct {
	grpc.ClientStream
}

func (x *testPerformanceBytesClient) Recv() (*TestBytes, error) {
	m := new(TestBytes)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TestServer is the server API for Test service.
// All implementations must embed UnimplementedTestServer
// for forward compatibility
type TestServer interface {
	Unary(context.Context, *String) (*String, error)
	ServerStream(*String, Test_ServerStreamServer) error
	ClientStream(Test_ClientStreamServer) error
	Bidirectional(Test_BidirectionalServer) error
	IntentionalError(*String, Test_IntentionalErrorServer) error
	Performance(*String, Test_PerformanceServer) error
	PerformanceBytes(*TestBytes, Test_PerformanceBytesServer) error
	mustEmbedUnimplementedTestServer()
}

// UnimplementedTestServer must be embedded to have forward compatible implementations.
type UnimplementedTestServer struct {
}

func (UnimplementedTestServer) Unary(context.Context, *String) (*String, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unary not implemented")
}
func (UnimplementedTestServer) ServerStream(*String, Test_ServerStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ServerStream not implemented")
}
func (UnimplementedTestServer) ClientStream(Test_ClientStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ClientStream not implemented")
}
func (UnimplementedTestServer) Bidirectional(Test_BidirectionalServer) error {
	return status.Errorf(codes.Unimplemented, "method Bidirectional not implemented")
}
func (UnimplementedTestServer) IntentionalError(*String, Test_IntentionalErrorServer) error {
	return status.Errorf(codes.Unimplemented, "method IntentionalError not implemented")
}
func (UnimplementedTestServer) Performance(*String, Test_PerformanceServer) error {
	return status.Errorf(codes.Unimplemented, "method Performance not implemented")
}
func (UnimplementedTestServer) PerformanceBytes(*TestBytes, Test_PerformanceBytesServer) error {
	return status.Errorf(codes.Unimplemented, "method PerformanceBytes not implemented")
}
func (UnimplementedTestServer) mustEmbedUnimplementedTestServer() {}

// UnsafeTestServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TestServer will
// result in compilation errors.
type UnsafeTestServer interface {
	mustEmbedUnimplementedTestServer()
}

func RegisterTestServer(s grpc.ServiceRegistrar, srv TestServer) {
	s.RegisterService(&Test_ServiceDesc, srv)
}

func _Test_Unary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServer).Unary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/test.Test/Unary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServer).Unary(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _Test_ServerStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(String)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TestServer).ServerStream(m, &testServerStreamServer{stream})
}

type Test_ServerStreamServer interface {
	Send(*String) error
	grpc.ServerStream
}

type testServerStreamServer struct {
	grpc.ServerStream
}

func (x *testServerStreamServer) Send(m *String) error {
	return x.ServerStream.SendMsg(m)
}

func _Test_ClientStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TestServer).ClientStream(&testClientStreamServer{stream})
}

type Test_ClientStreamServer interface {
	SendAndClose(*String) error
	Recv() (*String, error)
	grpc.ServerStream
}

type testClientStreamServer struct {
	grpc.ServerStream
}

func (x *testClientStreamServer) SendAndClose(m *String) error {
	return x.ServerStream.SendMsg(m)
}

func (x *testClientStreamServer) Recv() (*String, error) {
	m := new(String)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Test_Bidirectional_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TestServer).Bidirectional(&testBidirectionalServer{stream})
}

type Test_BidirectionalServer interface {
	Send(*String) error
	Recv() (*String, error)
	grpc.ServerStream
}

type testBidirectionalServer struct {
	grpc.ServerStream
}

func (x *testBidirectionalServer) Send(m *String) error {
	return x.ServerStream.SendMsg(m)
}

func (x *testBidirectionalServer) Recv() (*String, error) {
	m := new(String)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Test_IntentionalError_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(String)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TestServer).IntentionalError(m, &testIntentionalErrorServer{stream})
}

type Test_IntentionalErrorServer interface {
	Send(*String) error
	grpc.ServerStream
}

type testIntentionalErrorServer struct {
	grpc.ServerStream
}

func (x *testIntentionalErrorServer) Send(m *String) error {
	return x.ServerStream.SendMsg(m)
}

func _Test_Performance_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(String)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TestServer).Performance(m, &testPerformanceServer{stream})
}

type Test_PerformanceServer interface {
	Send(*String) error
	grpc.ServerStream
}

type testPerformanceServer struct {
	grpc.ServerStream
}

func (x *testPerformanceServer) Send(m *String) error {
	return x.ServerStream.SendMsg(m)
}

func _Test_PerformanceBytes_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TestBytes)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TestServer).PerformanceBytes(m, &testPerformanceBytesServer{stream})
}

type Test_PerformanceBytesServer interface {
	Send(*TestBytes) error
	grpc.ServerStream
}

type testPerformanceBytesServer struct {
	grpc.ServerStream
}

func (x *testPerformanceBytesServer) Send(m *TestBytes) error {
	return x.ServerStream.SendMsg(m)
}

// Test_ServiceDesc is the grpc.ServiceDesc for Test service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Test_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "test.Test",
	HandlerType: (*TestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Unary",
			Handler:    _Test_Unary_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ServerStream",
			Handler:       _Test_ServerStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ClientStream",
			Handler:       _Test_ClientStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Bidirectional",
			Handler:       _Test_Bidirectional_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "IntentionalError",
			Handler:       _Test_IntentionalError_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Performance",
			Handler:       _Test_Performance_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PerformanceBytes",
			Handler:       _Test_PerformanceBytes_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/protobuf/test.proto",
}

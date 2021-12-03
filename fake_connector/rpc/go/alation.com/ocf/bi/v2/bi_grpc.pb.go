// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v2

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BiClient is the client API for Bi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BiClient interface {
	// Verify that the connector has been configured correctly.
	ConfigurationVerification(ctx context.Context, in *Request, opts ...grpc.CallOption) (*empty.Empty, error)
	// Return the BI Server folder structure.
	ListFileSystem(ctx context.Context, in *ListFileSystemRequest, opts ...grpc.CallOption) (Bi_ListFileSystemClient, error)
	// Run full metadata extraction.
	RunFullExtraction(ctx context.Context, in *RunFullExtractionRequest, opts ...grpc.CallOption) (Bi_RunFullExtractionClient, error)
	// Certify or de-certify a BI object.
	PropagateCertification(ctx context.Context, in *PropagateCertificationRequest, opts ...grpc.CallOption) (Bi_PropagateCertificationClient, error)
}

type biClient struct {
	cc grpc.ClientConnInterface
}

func NewBiClient(cc grpc.ClientConnInterface) BiClient {
	return &biClient{cc}
}

func (c *biClient) ConfigurationVerification(ctx context.Context, in *Request, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/bi.Bi/ConfigurationVerification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *biClient) ListFileSystem(ctx context.Context, in *ListFileSystemRequest, opts ...grpc.CallOption) (Bi_ListFileSystemClient, error) {
	stream, err := c.cc.NewStream(ctx, &Bi_ServiceDesc.Streams[0], "/bi.Bi/ListFileSystem", opts...)
	if err != nil {
		return nil, err
	}
	x := &biListFileSystemClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Bi_ListFileSystemClient interface {
	Recv() (*BIObjectResponse, error)
	grpc.ClientStream
}

type biListFileSystemClient struct {
	grpc.ClientStream
}

func (x *biListFileSystemClient) Recv() (*BIObjectResponse, error) {
	m := new(BIObjectResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *biClient) RunFullExtraction(ctx context.Context, in *RunFullExtractionRequest, opts ...grpc.CallOption) (Bi_RunFullExtractionClient, error) {
	stream, err := c.cc.NewStream(ctx, &Bi_ServiceDesc.Streams[1], "/bi.Bi/RunFullExtraction", opts...)
	if err != nil {
		return nil, err
	}
	x := &biRunFullExtractionClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Bi_RunFullExtractionClient interface {
	Recv() (*BIObjectResponse, error)
	grpc.ClientStream
}

type biRunFullExtractionClient struct {
	grpc.ClientStream
}

func (x *biRunFullExtractionClient) Recv() (*BIObjectResponse, error) {
	m := new(BIObjectResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *biClient) PropagateCertification(ctx context.Context, in *PropagateCertificationRequest, opts ...grpc.CallOption) (Bi_PropagateCertificationClient, error) {
	stream, err := c.cc.NewStream(ctx, &Bi_ServiceDesc.Streams[2], "/bi.Bi/propagateCertification", opts...)
	if err != nil {
		return nil, err
	}
	x := &biPropagateCertificationClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Bi_PropagateCertificationClient interface {
	Recv() (*BIObjectResponse, error)
	grpc.ClientStream
}

type biPropagateCertificationClient struct {
	grpc.ClientStream
}

func (x *biPropagateCertificationClient) Recv() (*BIObjectResponse, error) {
	m := new(BIObjectResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BiServer is the server API for Bi service.
// All implementations must embed UnimplementedBiServer
// for forward compatibility
type BiServer interface {
	// Verify that the connector has been configured correctly.
	ConfigurationVerification(context.Context, *Request) (*empty.Empty, error)
	// Return the BI Server folder structure.
	ListFileSystem(*ListFileSystemRequest, Bi_ListFileSystemServer) error
	// Run full metadata extraction.
	RunFullExtraction(*RunFullExtractionRequest, Bi_RunFullExtractionServer) error
	// Certify or de-certify a BI object.
	PropagateCertification(*PropagateCertificationRequest, Bi_PropagateCertificationServer) error
	mustEmbedUnimplementedBiServer()
}

// UnimplementedBiServer must be embedded to have forward compatible implementations.
type UnimplementedBiServer struct {
}

func (UnimplementedBiServer) ConfigurationVerification(context.Context, *Request) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigurationVerification not implemented")
}
func (UnimplementedBiServer) ListFileSystem(*ListFileSystemRequest, Bi_ListFileSystemServer) error {
	return status.Errorf(codes.Unimplemented, "method ListFileSystem not implemented")
}
func (UnimplementedBiServer) RunFullExtraction(*RunFullExtractionRequest, Bi_RunFullExtractionServer) error {
	return status.Errorf(codes.Unimplemented, "method RunFullExtraction not implemented")
}
func (UnimplementedBiServer) PropagateCertification(*PropagateCertificationRequest, Bi_PropagateCertificationServer) error {
	return status.Errorf(codes.Unimplemented, "method PropagateCertification not implemented")
}
func (UnimplementedBiServer) mustEmbedUnimplementedBiServer() {}

// UnsafeBiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BiServer will
// result in compilation errors.
type UnsafeBiServer interface {
	mustEmbedUnimplementedBiServer()
}

func RegisterBiServer(s grpc.ServiceRegistrar, srv BiServer) {
	s.RegisterService(&Bi_ServiceDesc, srv)
}

func _Bi_ConfigurationVerification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BiServer).ConfigurationVerification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bi.Bi/ConfigurationVerification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BiServer).ConfigurationVerification(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bi_ListFileSystem_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListFileSystemRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BiServer).ListFileSystem(m, &biListFileSystemServer{stream})
}

type Bi_ListFileSystemServer interface {
	Send(*BIObjectResponse) error
	grpc.ServerStream
}

type biListFileSystemServer struct {
	grpc.ServerStream
}

func (x *biListFileSystemServer) Send(m *BIObjectResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Bi_RunFullExtraction_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RunFullExtractionRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BiServer).RunFullExtraction(m, &biRunFullExtractionServer{stream})
}

type Bi_RunFullExtractionServer interface {
	Send(*BIObjectResponse) error
	grpc.ServerStream
}

type biRunFullExtractionServer struct {
	grpc.ServerStream
}

func (x *biRunFullExtractionServer) Send(m *BIObjectResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Bi_PropagateCertification_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(PropagateCertificationRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BiServer).PropagateCertification(m, &biPropagateCertificationServer{stream})
}

type Bi_PropagateCertificationServer interface {
	Send(*BIObjectResponse) error
	grpc.ServerStream
}

type biPropagateCertificationServer struct {
	grpc.ServerStream
}

func (x *biPropagateCertificationServer) Send(m *BIObjectResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Bi_ServiceDesc is the grpc.ServiceDesc for Bi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Bi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bi.Bi",
	HandlerType: (*BiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConfigurationVerification",
			Handler:    _Bi_ConfigurationVerification_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListFileSystem",
			Handler:       _Bi_ListFileSystem_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "RunFullExtraction",
			Handler:       _Bi_RunFullExtraction_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "propagateCertification",
			Handler:       _Bi_PropagateCertification_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "bi.proto",
}
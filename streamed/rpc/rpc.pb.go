// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpc.proto

package rpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GrpcCallType int32

const (
	GrpcCallType_Unary               GrpcCallType = 0
	GrpcCallType_ServerStream        GrpcCallType = 1
	GrpcCallType_ClientStream        GrpcCallType = 2
	GrpcCallType_BiDirectionalStream GrpcCallType = 3
)

var GrpcCallType_name = map[int32]string{
	0: "Unary",
	1: "ServerStream",
	2: "ClientStream",
	3: "BiDirectionalStream",
}

var GrpcCallType_value = map[string]int32{
	"Unary":               0,
	"ServerStream":        1,
	"ClientStream":        2,
	"BiDirectionalStream": 3,
}

func (x GrpcCallType) String() string {
	return proto.EnumName(GrpcCallType_name, int32(x))
}

func (GrpcCallType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{0}
}

type PushStreamItem struct {
	Headers              *Headers `protobuf:"bytes,1,opt,name=headers,proto3" json:"headers,omitempty"`
	Sequence             uint64   `protobuf:"varint,2,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Body                 [][]byte `protobuf:"bytes,3,rep,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PushStreamItem) Reset()         { *m = PushStreamItem{} }
func (m *PushStreamItem) String() string { return proto.CompactTextString(m) }
func (*PushStreamItem) ProtoMessage()    {}
func (*PushStreamItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{0}
}

func (m *PushStreamItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PushStreamItem.Unmarshal(m, b)
}
func (m *PushStreamItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PushStreamItem.Marshal(b, m, deterministic)
}
func (m *PushStreamItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PushStreamItem.Merge(m, src)
}
func (m *PushStreamItem) XXX_Size() int {
	return xxx_messageInfo_PushStreamItem.Size(m)
}
func (m *PushStreamItem) XXX_DiscardUnknown() {
	xxx_messageInfo_PushStreamItem.DiscardUnknown(m)
}

var xxx_messageInfo_PushStreamItem proto.InternalMessageInfo

func (m *PushStreamItem) GetHeaders() *Headers {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *PushStreamItem) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *PushStreamItem) GetBody() [][]byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type JobOptional struct {
	// Types that are valid to be assigned to Options:
	//	*JobOptional_Job
	//	*JobOptional_Nothing
	Options              isJobOptional_Options `protobuf_oneof:"options"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *JobOptional) Reset()         { *m = JobOptional{} }
func (m *JobOptional) String() string { return proto.CompactTextString(m) }
func (*JobOptional) ProtoMessage()    {}
func (*JobOptional) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{1}
}

func (m *JobOptional) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JobOptional.Unmarshal(m, b)
}
func (m *JobOptional) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JobOptional.Marshal(b, m, deterministic)
}
func (m *JobOptional) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JobOptional.Merge(m, src)
}
func (m *JobOptional) XXX_Size() int {
	return xxx_messageInfo_JobOptional.Size(m)
}
func (m *JobOptional) XXX_DiscardUnknown() {
	xxx_messageInfo_JobOptional.DiscardUnknown(m)
}

var xxx_messageInfo_JobOptional proto.InternalMessageInfo

type isJobOptional_Options interface {
	isJobOptional_Options()
}

type JobOptional_Job struct {
	Job *Job `protobuf:"bytes,1,opt,name=job,proto3,oneof"`
}

type JobOptional_Nothing struct {
	Nothing *empty.Empty `protobuf:"bytes,2,opt,name=nothing,proto3,oneof"`
}

func (*JobOptional_Job) isJobOptional_Options() {}

func (*JobOptional_Nothing) isJobOptional_Options() {}

func (m *JobOptional) GetOptions() isJobOptional_Options {
	if m != nil {
		return m.Options
	}
	return nil
}

func (m *JobOptional) GetJob() *Job {
	if x, ok := m.GetOptions().(*JobOptional_Job); ok {
		return x.Job
	}
	return nil
}

func (m *JobOptional) GetNothing() *empty.Empty {
	if x, ok := m.GetOptions().(*JobOptional_Nothing); ok {
		return x.Nothing
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*JobOptional) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*JobOptional_Job)(nil),
		(*JobOptional_Nothing)(nil),
	}
}

type Job struct {
	Headers              *Headers `protobuf:"bytes,1,opt,name=headers,proto3" json:"headers,omitempty"`
	Body                 []byte   `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Job) Reset()         { *m = Job{} }
func (m *Job) String() string { return proto.CompactTextString(m) }
func (*Job) ProtoMessage()    {}
func (*Job) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{2}
}

func (m *Job) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Job.Unmarshal(m, b)
}
func (m *Job) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Job.Marshal(b, m, deterministic)
}
func (m *Job) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Job.Merge(m, src)
}
func (m *Job) XXX_Size() int {
	return xxx_messageInfo_Job.Size(m)
}
func (m *Job) XXX_DiscardUnknown() {
	xxx_messageInfo_Job.DiscardUnknown(m)
}

var xxx_messageInfo_Job proto.InternalMessageInfo

func (m *Job) GetHeaders() *Headers {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *Job) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type Headers struct {
	Method               string       `protobuf:"bytes,1,opt,name=Method,proto3" json:"Method,omitempty"`
	Type                 GrpcCallType `protobuf:"varint,2,opt,name=type,proto3,enum=rpc.GrpcCallType" json:"type,omitempty"`
	JobID                uint64       `protobuf:"varint,3,opt,name=jobID,proto3" json:"jobID,omitempty"`
	Connector            uint64       `protobuf:"varint,4,opt,name=connector,proto3" json:"connector,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Headers) Reset()         { *m = Headers{} }
func (m *Headers) String() string { return proto.CompactTextString(m) }
func (*Headers) ProtoMessage()    {}
func (*Headers) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{3}
}

func (m *Headers) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Headers.Unmarshal(m, b)
}
func (m *Headers) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Headers.Marshal(b, m, deterministic)
}
func (m *Headers) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Headers.Merge(m, src)
}
func (m *Headers) XXX_Size() int {
	return xxx_messageInfo_Headers.Size(m)
}
func (m *Headers) XXX_DiscardUnknown() {
	xxx_messageInfo_Headers.DiscardUnknown(m)
}

var xxx_messageInfo_Headers proto.InternalMessageInfo

func (m *Headers) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *Headers) GetType() GrpcCallType {
	if m != nil {
		return m.Type
	}
	return GrpcCallType_Unary
}

func (m *Headers) GetJobID() uint64 {
	if m != nil {
		return m.JobID
	}
	return 0
}

func (m *Headers) GetConnector() uint64 {
	if m != nil {
		return m.Connector
	}
	return 0
}

func init() {
	proto.RegisterEnum("rpc.GrpcCallType", GrpcCallType_name, GrpcCallType_value)
	proto.RegisterType((*PushStreamItem)(nil), "rpc.PushStreamItem")
	proto.RegisterType((*JobOptional)(nil), "rpc.JobOptional")
	proto.RegisterType((*Job)(nil), "rpc.Job")
	proto.RegisterType((*Headers)(nil), "rpc.Headers")
}

func init() { proto.RegisterFile("rpc.proto", fileDescriptor_77a6da22d6a3feb1) }

var fileDescriptor_77a6da22d6a3feb1 = []byte{
	// 413 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0xb5, 0xeb, 0xa4, 0x8e, 0x27, 0x51, 0x15, 0xa6, 0xa8, 0x58, 0xa1, 0x87, 0xc8, 0x12, 0x28,
	0xe2, 0xe0, 0x4a, 0x46, 0x1c, 0x39, 0x90, 0x16, 0x91, 0x46, 0x42, 0x54, 0x2e, 0x70, 0xe0, 0x66,
	0x3b, 0x43, 0xec, 0xc8, 0xf1, 0x2c, 0xeb, 0x0d, 0x92, 0x25, 0x0e, 0x7c, 0x3a, 0xda, 0xb5, 0x4d,
	0xe1, 0x10, 0xa9, 0x37, 0xcf, 0x7b, 0xe3, 0x37, 0xfb, 0xde, 0x0c, 0x78, 0x52, 0x64, 0xa1, 0x90,
	0xac, 0x18, 0x1d, 0x29, 0xb2, 0xd9, 0xf3, 0x2d, 0xf3, 0xb6, 0xa4, 0x2b, 0x03, 0xa5, 0x87, 0xef,
	0x57, 0xb4, 0x17, 0xaa, 0x69, 0x3b, 0x82, 0x1c, 0xce, 0xee, 0x0e, 0x75, 0x7e, 0xaf, 0x24, 0x25,
	0xfb, 0x5b, 0x45, 0x7b, 0x7c, 0x09, 0x6e, 0x4e, 0xc9, 0x86, 0x64, 0xed, 0xdb, 0x73, 0x7b, 0x31,
	0x8e, 0x26, 0xa1, 0x16, 0x5c, 0xb5, 0x58, 0xdc, 0x93, 0x38, 0x83, 0x51, 0x4d, 0x3f, 0x0e, 0x54,
	0x65, 0xe4, 0x9f, 0xcc, 0xed, 0xc5, 0x20, 0xfe, 0x5b, 0x23, 0xc2, 0x20, 0xe5, 0x4d, 0xe3, 0x3b,
	0x73, 0x67, 0x31, 0x89, 0xcd, 0x77, 0xb0, 0x83, 0xf1, 0x9a, 0xd3, 0x4f, 0x42, 0x15, 0x5c, 0x25,
	0x25, 0x5e, 0x82, 0xb3, 0xe3, 0xb4, 0x1b, 0x31, 0x32, 0x23, 0xd6, 0x9c, 0xae, 0xac, 0x58, 0xc3,
	0x18, 0x81, 0x5b, 0xb1, 0xca, 0x8b, 0x6a, 0x6b, 0xb4, 0xc7, 0xd1, 0x45, 0xd8, 0xba, 0x08, 0x7b,
	0x17, 0xe1, 0x7b, 0xed, 0x62, 0x65, 0xc5, 0x7d, 0xe3, 0xd2, 0x03, 0x97, 0x8d, 0x7a, 0x1d, 0xbc,
	0x03, 0x67, 0xcd, 0xe9, 0xa3, 0xad, 0xf4, 0xcf, 0xd5, 0xa3, 0xfa, 0xe7, 0xfe, 0x02, 0xb7, 0xeb,
	0xc3, 0x0b, 0x38, 0xfd, 0x48, 0x2a, 0xe7, 0x8d, 0x51, 0xf1, 0xe2, 0xae, 0xc2, 0x17, 0x30, 0x50,
	0x8d, 0x68, 0xdd, 0x9f, 0x45, 0x4f, 0x8c, 0xf6, 0x07, 0x29, 0xb2, 0xeb, 0xa4, 0x2c, 0x3f, 0x37,
	0x82, 0x62, 0x43, 0xe3, 0x53, 0x18, 0xee, 0x38, 0xbd, 0xbd, 0xf1, 0x1d, 0x93, 0x52, 0x5b, 0xe0,
	0x25, 0x78, 0x19, 0x57, 0x15, 0x65, 0x8a, 0xa5, 0x3f, 0x30, 0xcc, 0x03, 0xf0, 0xea, 0x2b, 0x4c,
	0xfe, 0x55, 0x42, 0x0f, 0x86, 0x5f, 0xaa, 0x44, 0x36, 0x53, 0x0b, 0xa7, 0x30, 0xb9, 0x27, 0xf9,
	0x93, 0x64, 0xbb, 0xb3, 0xa9, 0xad, 0x91, 0xeb, 0xb2, 0xa0, 0x4a, 0x75, 0xc8, 0x09, 0x3e, 0x83,
	0xf3, 0x65, 0x71, 0x53, 0x48, 0xca, 0xda, 0xb4, 0x3b, 0xc2, 0x89, 0x7e, 0xdb, 0x30, 0xd2, 0x7f,
	0x17, 0x19, 0x49, 0x7c, 0x03, 0xee, 0x1d, 0x97, 0xa5, 0x4e, 0xea, 0x48, 0xbc, 0xb3, 0x69, 0xbf,
	0x98, 0x7e, 0x6f, 0x81, 0x85, 0x6f, 0x01, 0x1e, 0x4e, 0x06, 0xcf, 0x4d, 0xc7, 0xff, 0x37, 0x34,
	0x3b, 0x22, 0x17, 0x58, 0x0b, 0x7b, 0x39, 0xfc, 0xa6, 0xaf, 0x32, 0x3d, 0x35, 0xd4, 0xeb, 0x3f,
	0x01, 0x00, 0x00, 0xff, 0xff, 0xff, 0xdb, 0x28, 0x80, 0xae, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ServicerClient is the client API for Servicer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServicerClient interface {
	PollJob(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*JobOptional, error)
	PushStream(ctx context.Context, opts ...grpc.CallOption) (Servicer_PushStreamClient, error)
}

type servicerClient struct {
	cc *grpc.ClientConn
}

func NewServicerClient(cc *grpc.ClientConn) ServicerClient {
	return &servicerClient{cc}
}

func (c *servicerClient) PollJob(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*JobOptional, error) {
	out := new(JobOptional)
	err := c.cc.Invoke(ctx, "/rpc.Servicer/PollJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicerClient) PushStream(ctx context.Context, opts ...grpc.CallOption) (Servicer_PushStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Servicer_serviceDesc.Streams[0], "/rpc.Servicer/PushStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &servicerPushStreamClient{stream}
	return x, nil
}

type Servicer_PushStreamClient interface {
	Send(*PushStreamItem) error
	CloseAndRecv() (*empty.Empty, error)
	grpc.ClientStream
}

type servicerPushStreamClient struct {
	grpc.ClientStream
}

func (x *servicerPushStreamClient) Send(m *PushStreamItem) error {
	return x.ClientStream.SendMsg(m)
}

func (x *servicerPushStreamClient) CloseAndRecv() (*empty.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(empty.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ServicerServer is the server API for Servicer service.
type ServicerServer interface {
	PollJob(context.Context, *empty.Empty) (*JobOptional, error)
	PushStream(Servicer_PushStreamServer) error
}

// UnimplementedServicerServer can be embedded to have forward compatible implementations.
type UnimplementedServicerServer struct {
}

func (*UnimplementedServicerServer) PollJob(ctx context.Context, req *empty.Empty) (*JobOptional, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PollJob not implemented")
}
func (*UnimplementedServicerServer) PushStream(srv Servicer_PushStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method PushStream not implemented")
}

func RegisterServicerServer(s *grpc.Server, srv ServicerServer) {
	s.RegisterService(&_Servicer_serviceDesc, srv)
}

func _Servicer_PollJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicerServer).PollJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Servicer/PollJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicerServer).PollJob(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Servicer_PushStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServicerServer).PushStream(&servicerPushStreamServer{stream})
}

type Servicer_PushStreamServer interface {
	SendAndClose(*empty.Empty) error
	Recv() (*PushStreamItem, error)
	grpc.ServerStream
}

type servicerPushStreamServer struct {
	grpc.ServerStream
}

func (x *servicerPushStreamServer) SendAndClose(m *empty.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *servicerPushStreamServer) Recv() (*PushStreamItem, error) {
	m := new(PushStreamItem)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Servicer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Servicer",
	HandlerType: (*ServicerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PollJob",
			Handler:    _Servicer_PollJob_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PushStream",
			Handler:       _Servicer_PushStream_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "rpc.proto",
}

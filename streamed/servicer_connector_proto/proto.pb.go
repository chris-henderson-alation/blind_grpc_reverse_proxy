// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto.proto

package servicer_connector_proto

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type Complex_Corpus int32

const (
	Complex_UNIVERSAL Complex_Corpus = 0
	Complex_WEB       Complex_Corpus = 1
	Complex_IMAGES    Complex_Corpus = 2
	Complex_LOCAL     Complex_Corpus = 3
	Complex_NEWS      Complex_Corpus = 4
	Complex_PRODUCTS  Complex_Corpus = 5
	Complex_VIDEO     Complex_Corpus = 6
)

var Complex_Corpus_name = map[int32]string{
	0: "UNIVERSAL",
	1: "WEB",
	2: "IMAGES",
	3: "LOCAL",
	4: "NEWS",
	5: "PRODUCTS",
	6: "VIDEO",
}

var Complex_Corpus_value = map[string]int32{
	"UNIVERSAL": 0,
	"WEB":       1,
	"IMAGES":    2,
	"LOCAL":     3,
	"NEWS":      4,
	"PRODUCTS":  5,
	"VIDEO":     6,
}

func (x Complex_Corpus) String() string {
	return proto.EnumName(Complex_Corpus_name, int32(x))
}

func (Complex_Corpus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{2, 0}
}

type Bytes struct {
	Inner                []byte   `protobuf:"bytes,1,opt,name=inner,proto3" json:"inner,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Bytes) Reset()         { *m = Bytes{} }
func (m *Bytes) String() string { return proto.CompactTextString(m) }
func (*Bytes) ProtoMessage()    {}
func (*Bytes) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{0}
}

func (m *Bytes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Bytes.Unmarshal(m, b)
}
func (m *Bytes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Bytes.Marshal(b, m, deterministic)
}
func (m *Bytes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bytes.Merge(m, src)
}
func (m *Bytes) XXX_Size() int {
	return xxx_messageInfo_Bytes.Size(m)
}
func (m *Bytes) XXX_DiscardUnknown() {
	xxx_messageInfo_Bytes.DiscardUnknown(m)
}

var xxx_messageInfo_Bytes proto.InternalMessageInfo

func (m *Bytes) GetInner() []byte {
	if m != nil {
		return m.Inner
	}
	return nil
}

type Ball struct {
	Ball                 string   `protobuf:"bytes,1,opt,name=ball,proto3" json:"ball,omitempty"`
	Complex              *Complex `protobuf:"bytes,2,opt,name=complex,proto3" json:"complex,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ball) Reset()         { *m = Ball{} }
func (m *Ball) String() string { return proto.CompactTextString(m) }
func (*Ball) ProtoMessage()    {}
func (*Ball) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{1}
}

func (m *Ball) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ball.Unmarshal(m, b)
}
func (m *Ball) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ball.Marshal(b, m, deterministic)
}
func (m *Ball) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ball.Merge(m, src)
}
func (m *Ball) XXX_Size() int {
	return xxx_messageInfo_Ball.Size(m)
}
func (m *Ball) XXX_DiscardUnknown() {
	xxx_messageInfo_Ball.DiscardUnknown(m)
}

var xxx_messageInfo_Ball proto.InternalMessageInfo

func (m *Ball) GetBall() string {
	if m != nil {
		return m.Ball
	}
	return ""
}

func (m *Ball) GetComplex() *Complex {
	if m != nil {
		return m.Complex
	}
	return nil
}

type Complex struct {
	Array                []uint32        `protobuf:"varint,1,rep,packed,name=array,proto3" json:"array,omitempty"`
	Corpus               Complex_Corpus  `protobuf:"varint,4,opt,name=corpus,proto3,enum=ocf.Complex_Corpus" json:"corpus,omitempty"`
	Result               *Complex_Result `protobuf:"bytes,5,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Complex) Reset()         { *m = Complex{} }
func (m *Complex) String() string { return proto.CompactTextString(m) }
func (*Complex) ProtoMessage()    {}
func (*Complex) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{2}
}

func (m *Complex) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Complex.Unmarshal(m, b)
}
func (m *Complex) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Complex.Marshal(b, m, deterministic)
}
func (m *Complex) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Complex.Merge(m, src)
}
func (m *Complex) XXX_Size() int {
	return xxx_messageInfo_Complex.Size(m)
}
func (m *Complex) XXX_DiscardUnknown() {
	xxx_messageInfo_Complex.DiscardUnknown(m)
}

var xxx_messageInfo_Complex proto.InternalMessageInfo

func (m *Complex) GetArray() []uint32 {
	if m != nil {
		return m.Array
	}
	return nil
}

func (m *Complex) GetCorpus() Complex_Corpus {
	if m != nil {
		return m.Corpus
	}
	return Complex_UNIVERSAL
}

func (m *Complex) GetResult() *Complex_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

type Complex_Result struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Title                string   `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Snippets             []string `protobuf:"bytes,3,rep,name=snippets,proto3" json:"snippets,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Complex_Result) Reset()         { *m = Complex_Result{} }
func (m *Complex_Result) String() string { return proto.CompactTextString(m) }
func (*Complex_Result) ProtoMessage()    {}
func (*Complex_Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{2, 0}
}

func (m *Complex_Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Complex_Result.Unmarshal(m, b)
}
func (m *Complex_Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Complex_Result.Marshal(b, m, deterministic)
}
func (m *Complex_Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Complex_Result.Merge(m, src)
}
func (m *Complex_Result) XXX_Size() int {
	return xxx_messageInfo_Complex_Result.Size(m)
}
func (m *Complex_Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Complex_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Complex_Result proto.InternalMessageInfo

func (m *Complex_Result) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Complex_Result) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Complex_Result) GetSnippets() []string {
	if m != nil {
		return m.Snippets
	}
	return nil
}

type Incoming struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Incoming) Reset()         { *m = Incoming{} }
func (m *Incoming) String() string { return proto.CompactTextString(m) }
func (*Incoming) ProtoMessage()    {}
func (*Incoming) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{3}
}

func (m *Incoming) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Incoming.Unmarshal(m, b)
}
func (m *Incoming) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Incoming.Marshal(b, m, deterministic)
}
func (m *Incoming) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Incoming.Merge(m, src)
}
func (m *Incoming) XXX_Size() int {
	return xxx_messageInfo_Incoming.Size(m)
}
func (m *Incoming) XXX_DiscardUnknown() {
	xxx_messageInfo_Incoming.DiscardUnknown(m)
}

var xxx_messageInfo_Incoming proto.InternalMessageInfo

func (m *Incoming) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Status struct {
	Fname                string   `protobuf:"bytes,1,opt,name=fname,proto3" json:"fname,omitempty"`
	Ok                   bool     `protobuf:"varint,2,opt,name=ok,proto3" json:"ok,omitempty"`
	Error                string   `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_2fcc84b9998d60d8, []int{4}
}

func (m *Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Status.Unmarshal(m, b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Status.Marshal(b, m, deterministic)
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return xxx_messageInfo_Status.Size(m)
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetFname() string {
	if m != nil {
		return m.Fname
	}
	return ""
}

func (m *Status) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

func (m *Status) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterEnum("ocf.Complex_Corpus", Complex_Corpus_name, Complex_Corpus_value)
	proto.RegisterType((*Bytes)(nil), "ocf.Bytes")
	proto.RegisterType((*Ball)(nil), "ocf.Ball")
	proto.RegisterType((*Complex)(nil), "ocf.Complex")
	proto.RegisterType((*Complex_Result)(nil), "ocf.Complex.Result")
	proto.RegisterType((*Incoming)(nil), "ocf.Incoming")
	proto.RegisterType((*Status)(nil), "ocf.Status")
}

func init() { proto.RegisterFile("proto.proto", fileDescriptor_2fcc84b9998d60d8) }

var fileDescriptor_2fcc84b9998d60d8 = []byte{
	// 468 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xdb, 0x6a, 0xdb, 0x40,
	0x10, 0x86, 0x23, 0xeb, 0x60, 0x6b, 0xe2, 0x04, 0x31, 0x2d, 0x45, 0xa8, 0xb4, 0xb8, 0xa2, 0x14,
	0x43, 0xa9, 0x52, 0xdc, 0x27, 0xb0, 0x6c, 0x51, 0x0c, 0x6e, 0x1c, 0xd6, 0x39, 0x40, 0x6e, 0x82,
	0xa2, 0xae, 0x8d, 0x88, 0xb4, 0x2b, 0x56, 0xab, 0x52, 0xf7, 0x75, 0xfa, 0xa2, 0x65, 0x77, 0xad,
	0x1e, 0x2e, 0x7a, 0x23, 0xcd, 0x3f, 0xf3, 0x49, 0x33, 0xbb, 0xf3, 0xc3, 0x69, 0x23, 0xb8, 0xe4,
	0x89, 0x7e, 0xa2, 0xcd, 0x8b, 0x5d, 0xf4, 0x72, 0xcf, 0xf9, 0xbe, 0xa2, 0x17, 0x3a, 0xf5, 0xd8,
	0xed, 0x2e, 0x68, 0xdd, 0xc8, 0x83, 0x21, 0xe2, 0x57, 0xe0, 0xa6, 0x07, 0x49, 0x5b, 0x7c, 0x0e,
	0x6e, 0xc9, 0x18, 0x15, 0xa1, 0x35, 0xb1, 0xa6, 0x63, 0x62, 0x44, 0x9c, 0x82, 0x93, 0xe6, 0x55,
	0x85, 0x08, 0xce, 0x63, 0x5e, 0x55, 0xba, 0xe8, 0x13, 0x1d, 0xe3, 0x3b, 0x18, 0x16, 0xbc, 0x6e,
	0x2a, 0xfa, 0x3d, 0x1c, 0x4c, 0xac, 0xe9, 0xe9, 0x6c, 0x9c, 0xf0, 0x62, 0x97, 0x2c, 0x4c, 0x8e,
	0xf4, 0xc5, 0xf8, 0xe7, 0x00, 0x86, 0xc7, 0xa4, 0xea, 0x92, 0x0b, 0x91, 0x1f, 0x42, 0x6b, 0x62,
	0x4f, 0xcf, 0x88, 0x11, 0xf8, 0x1e, 0xbc, 0x82, 0x8b, 0xa6, 0x6b, 0x43, 0x67, 0x62, 0x4d, 0xcf,
	0x67, 0xcf, 0xfe, 0xfe, 0x51, 0xb2, 0xd0, 0x25, 0x72, 0x44, 0x14, 0x2c, 0x68, 0xdb, 0x55, 0x32,
	0x74, 0x75, 0xd7, 0x7f, 0x61, 0xa2, 0x4b, 0xe4, 0x88, 0x44, 0x6b, 0xf0, 0x4c, 0x06, 0x03, 0xb0,
	0x3b, 0xd1, 0x1f, 0x40, 0x85, 0x6a, 0x16, 0x59, 0xca, 0x8a, 0xea, 0xe9, 0x7d, 0x62, 0x04, 0x46,
	0x30, 0x6a, 0x59, 0xd9, 0x34, 0x54, 0xb6, 0xa1, 0x3d, 0xb1, 0xa7, 0x3e, 0xf9, 0xad, 0xe3, 0x7b,
	0xf0, 0xcc, 0x30, 0x78, 0x06, 0xfe, 0xcd, 0xe5, 0xea, 0x36, 0x23, 0xdb, 0xf9, 0x3a, 0x38, 0xc1,
	0x21, 0xd8, 0x77, 0x59, 0x1a, 0x58, 0x08, 0xe0, 0xad, 0xbe, 0xcc, 0x3f, 0x67, 0xdb, 0x60, 0x80,
	0x3e, 0xb8, 0xeb, 0xcd, 0x62, 0xbe, 0x0e, 0x6c, 0x1c, 0x81, 0x73, 0x99, 0xdd, 0x6d, 0x03, 0x07,
	0xc7, 0x30, 0xba, 0x22, 0x9b, 0xe5, 0xcd, 0xe2, 0x7a, 0x1b, 0xb8, 0x0a, 0xb9, 0x5d, 0x2d, 0xb3,
	0x4d, 0xe0, 0xc5, 0x6f, 0x61, 0xb4, 0x62, 0x05, 0xaf, 0x4b, 0xb6, 0xc7, 0x10, 0x86, 0x35, 0x6d,
	0xdb, 0x7c, 0x4f, 0x8f, 0xf3, 0xf6, 0x32, 0x5e, 0x82, 0xb7, 0x95, 0xb9, 0xec, 0xf4, 0xbe, 0x76,
	0x2c, 0xaf, 0x7b, 0xc2, 0x08, 0x3c, 0x87, 0x01, 0x7f, 0xd2, 0x07, 0x1a, 0x91, 0x01, 0x7f, 0x52,
	0x14, 0x15, 0x82, 0x8b, 0xd0, 0x36, 0x94, 0x16, 0xb3, 0x1f, 0xe0, 0x5d, 0x53, 0xc6, 0xca, 0x16,
	0x5f, 0x83, 0x73, 0xa5, 0x3a, 0xfa, 0xfa, 0x12, 0xd5, 0xaa, 0xa3, 0x3f, 0x61, 0x7c, 0x82, 0x6f,
	0xc0, 0x25, 0x79, 0x53, 0x7e, 0xfd, 0x1f, 0xf0, 0xd1, 0xc2, 0x0f, 0x60, 0xa7, 0xe5, 0x1e, 0x5f,
	0x24, 0xc6, 0x66, 0x49, 0x6f, 0xb3, 0x24, 0x53, 0x36, 0x8b, 0xc0, 0xd0, 0xca, 0x63, 0x0a, 0x4f,
	0xa3, 0xfb, 0xb0, 0xa5, 0xe2, 0x5b, 0x59, 0x50, 0xf1, 0x50, 0x70, 0xc6, 0x68, 0x21, 0xb9, 0x78,
	0x30, 0x1f, 0x79, 0xfa, 0xf5, 0xe9, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf4, 0x99, 0xa6, 0xbc,
	0xc4, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TennisClient is the client API for Tennis service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TennisClient interface {
	Ping(ctx context.Context, in *Ball, opts ...grpc.CallOption) (*Ball, error)
	Rapid(ctx context.Context, in *Ball, opts ...grpc.CallOption) (Tennis_RapidClient, error)
	Big(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (Tennis_BigClient, error)
}

type tennisClient struct {
	cc *grpc.ClientConn
}

func NewTennisClient(cc *grpc.ClientConn) TennisClient {
	return &tennisClient{cc}
}

func (c *tennisClient) Ping(ctx context.Context, in *Ball, opts ...grpc.CallOption) (*Ball, error) {
	out := new(Ball)
	err := c.cc.Invoke(ctx, "/ocf.Tennis/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tennisClient) Rapid(ctx context.Context, in *Ball, opts ...grpc.CallOption) (Tennis_RapidClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Tennis_serviceDesc.Streams[0], "/ocf.Tennis/Rapid", opts...)
	if err != nil {
		return nil, err
	}
	x := &tennisRapidClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Tennis_RapidClient interface {
	Recv() (*Ball, error)
	grpc.ClientStream
}

type tennisRapidClient struct {
	grpc.ClientStream
}

func (x *tennisRapidClient) Recv() (*Ball, error) {
	m := new(Ball)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *tennisClient) Big(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (Tennis_BigClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Tennis_serviceDesc.Streams[1], "/ocf.Tennis/Big", opts...)
	if err != nil {
		return nil, err
	}
	x := &tennisBigClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Tennis_BigClient interface {
	Recv() (*Bytes, error)
	grpc.ClientStream
}

type tennisBigClient struct {
	grpc.ClientStream
}

func (x *tennisBigClient) Recv() (*Bytes, error) {
	m := new(Bytes)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TennisServer is the server API for Tennis service.
type TennisServer interface {
	Ping(context.Context, *Ball) (*Ball, error)
	Rapid(*Ball, Tennis_RapidServer) error
	Big(*empty.Empty, Tennis_BigServer) error
}

// UnimplementedTennisServer can be embedded to have forward compatible implementations.
type UnimplementedTennisServer struct {
}

func (*UnimplementedTennisServer) Ping(ctx context.Context, req *Ball) (*Ball, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedTennisServer) Rapid(req *Ball, srv Tennis_RapidServer) error {
	return status.Errorf(codes.Unimplemented, "method Rapid not implemented")
}
func (*UnimplementedTennisServer) Big(req *empty.Empty, srv Tennis_BigServer) error {
	return status.Errorf(codes.Unimplemented, "method Big not implemented")
}

func RegisterTennisServer(s *grpc.Server, srv TennisServer) {
	s.RegisterService(&_Tennis_serviceDesc, srv)
}

func _Tennis_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ball)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TennisServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ocf.Tennis/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TennisServer).Ping(ctx, req.(*Ball))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tennis_Rapid_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Ball)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TennisServer).Rapid(m, &tennisRapidServer{stream})
}

type Tennis_RapidServer interface {
	Send(*Ball) error
	grpc.ServerStream
}

type tennisRapidServer struct {
	grpc.ServerStream
}

func (x *tennisRapidServer) Send(m *Ball) error {
	return x.ServerStream.SendMsg(m)
}

func _Tennis_Big_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TennisServer).Big(m, &tennisBigServer{stream})
}

type Tennis_BigServer interface {
	Send(*Bytes) error
	grpc.ServerStream
}

type tennisBigServer struct {
	grpc.ServerStream
}

func (x *tennisBigServer) Send(m *Bytes) error {
	return x.ServerStream.SendMsg(m)
}

var _Tennis_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ocf.Tennis",
	HandlerType: (*TennisServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Tennis_Ping_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Rapid",
			Handler:       _Tennis_Rapid_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Big",
			Handler:       _Tennis_Big_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto.proto",
}

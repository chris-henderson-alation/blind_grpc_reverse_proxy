// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/protobuf/test.proto

package grpcinverter

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type String struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *String) Reset()         { *m = String{} }
func (m *String) String() string { return proto.CompactTextString(m) }
func (*String) ProtoMessage()    {}
func (*String) Descriptor() ([]byte, []int) {
	return fileDescriptor_6641c415498ccafb, []int{0}
}

func (m *String) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_String.Unmarshal(m, b)
}
func (m *String) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_String.Marshal(b, m, deterministic)
}
func (m *String) XXX_Merge(src proto.Message) {
	xxx_messageInfo_String.Merge(m, src)
}
func (m *String) XXX_Size() int {
	return xxx_messageInfo_String.Size(m)
}
func (m *String) XXX_DiscardUnknown() {
	xxx_messageInfo_String.DiscardUnknown(m)
}

var xxx_messageInfo_String proto.InternalMessageInfo

func (m *String) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type TestBytes struct {
	Body                 []byte   `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TestBytes) Reset()         { *m = TestBytes{} }
func (m *TestBytes) String() string { return proto.CompactTextString(m) }
func (*TestBytes) ProtoMessage()    {}
func (*TestBytes) Descriptor() ([]byte, []int) {
	return fileDescriptor_6641c415498ccafb, []int{1}
}

func (m *TestBytes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TestBytes.Unmarshal(m, b)
}
func (m *TestBytes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TestBytes.Marshal(b, m, deterministic)
}
func (m *TestBytes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TestBytes.Merge(m, src)
}
func (m *TestBytes) XXX_Size() int {
	return xxx_messageInfo_TestBytes.Size(m)
}
func (m *TestBytes) XXX_DiscardUnknown() {
	xxx_messageInfo_TestBytes.DiscardUnknown(m)
}

var xxx_messageInfo_TestBytes proto.InternalMessageInfo

func (m *TestBytes) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func init() {
	proto.RegisterType((*String)(nil), "test.String")
	proto.RegisterType((*TestBytes)(nil), "test.TestBytes")
}

func init() { proto.RegisterFile("pkg/protobuf/test.proto", fileDescriptor_6641c415498ccafb) }

var fileDescriptor_6641c415498ccafb = []byte{
	// 257 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0x5d, 0x89, 0x95, 0x8e, 0x51, 0xc3, 0x5e, 0x2c, 0xbd, 0x28, 0x01, 0xa1, 0x17, 0xd3,
	0x5a, 0x2f, 0x82, 0xb7, 0x88, 0x07, 0x6f, 0xd2, 0xe8, 0xc5, 0xdb, 0x26, 0x9d, 0x86, 0xc5, 0x66,
	0x37, 0xcc, 0x8e, 0x81, 0xfc, 0x58, 0xff, 0x8b, 0x64, 0x83, 0xa2, 0x1e, 0x4c, 0x6f, 0xef, 0x2d,
	0xdf, 0xdb, 0x19, 0xde, 0xc0, 0x59, 0xfd, 0x56, 0xce, 0x6b, 0xb2, 0x6c, 0xf3, 0xf7, 0xcd, 0x9c,
	0xd1, 0x71, 0xe2, 0x9d, 0x0c, 0x3a, 0x1d, 0xc7, 0x30, 0xca, 0x98, 0xb4, 0x29, 0xe5, 0x04, 0x0e,
	0x2b, 0x74, 0x4e, 0x95, 0x38, 0x11, 0x17, 0x62, 0x36, 0x5e, 0x7d, 0xd9, 0xf8, 0x1c, 0xc6, 0xcf,
	0xe8, 0x38, 0x6d, 0x19, 0x9d, 0x94, 0x10, 0xe4, 0x76, 0xdd, 0x7a, 0x26, 0x5c, 0x79, 0xbd, 0xfc,
	0xd8, 0x87, 0xa0, 0x23, 0xe4, 0x25, 0x1c, 0xbc, 0x18, 0x45, 0xad, 0x0c, 0x13, 0x3f, 0xa9, 0xff,
	0x7a, 0xfa, 0xcb, 0xc5, 0x7b, 0x32, 0x81, 0x30, 0x43, 0x6a, 0x90, 0x32, 0x26, 0x54, 0xd5, 0xff,
	0xf4, 0x42, 0x74, 0xfc, 0xfd, 0x56, 0xa3, 0xe1, 0x5d, 0xf8, 0x99, 0x90, 0xd7, 0x70, 0x9c, 0xea,
	0xb5, 0x26, 0x2c, 0x58, 0x5b, 0xa3, 0xb6, 0x43, 0x81, 0x85, 0x90, 0x4b, 0x88, 0x1e, 0x0d, 0xa3,
	0xe9, 0x03, 0x0f, 0x44, 0x96, 0x06, 0xd7, 0xba, 0x82, 0xa3, 0x27, 0xa4, 0x8d, 0xa5, 0x4a, 0x99,
	0x02, 0x07, 0xf1, 0x5b, 0x88, 0x7e, 0xe0, 0x7d, 0x9b, 0xa7, 0x3d, 0xf5, 0x5d, 0xef, 0xf4, 0xef,
	0x43, 0x97, 0x4c, 0xa3, 0xd7, 0x93, 0xe4, 0xae, 0xa4, 0xba, 0xd0, 0xa6, 0x41, 0x62, 0xa4, 0x7c,
	0xe4, 0x6f, 0x78, 0xf3, 0x19, 0x00, 0x00, 0xff, 0xff, 0x9d, 0xa3, 0x92, 0xd4, 0xde, 0x01, 0x00,
	0x00,
}
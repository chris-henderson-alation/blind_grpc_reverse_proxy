// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: bi_certify.proto

package v2

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// A PropagateCertificationRequest is made whenever an object get an endorsement, warning or deprecation.
type PropagateCertificationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Configuration *ProtobufConfiguration `protobuf:"bytes,1,opt,name=configuration,proto3" json:"configuration,omitempty"`
	BiObject      *GBMBIObject           `protobuf:"bytes,2,opt,name=biObject,proto3" json:"biObject,omitempty"`
	// If certify is true, then it contains the concatenation of all the endorsements.
	// If certify is false, then it contains the deprecation or warning message from the flag.
	Note string `protobuf:"bytes,3,opt,name=note,proto3" json:"note,omitempty"`
	// Set to true, if the object has only endorsements without any deprecations or warnings.
	Certify bool `protobuf:"varint,4,opt,name=certify,proto3" json:"certify,omitempty"`
}

func (x *PropagateCertificationRequest) Reset() {
	*x = PropagateCertificationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bi_certify_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PropagateCertificationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PropagateCertificationRequest) ProtoMessage() {}

func (x *PropagateCertificationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_bi_certify_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PropagateCertificationRequest.ProtoReflect.Descriptor instead.
func (*PropagateCertificationRequest) Descriptor() ([]byte, []int) {
	return file_bi_certify_proto_rawDescGZIP(), []int{0}
}

func (x *PropagateCertificationRequest) GetConfiguration() *ProtobufConfiguration {
	if x != nil {
		return x.Configuration
	}
	return nil
}

func (x *PropagateCertificationRequest) GetBiObject() *GBMBIObject {
	if x != nil {
		return x.BiObject
	}
	return nil
}

func (x *PropagateCertificationRequest) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

func (x *PropagateCertificationRequest) GetCertify() bool {
	if x != nil {
		return x.Certify
	}
	return false
}

var File_bi_certify_proto protoreflect.FileDescriptor

var file_bi_certify_proto_rawDesc = []byte{
	0x0a, 0x10, 0x62, 0x69, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x62, 0x69, 0x2e, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x79, 0x1a, 0x16,
	0x62, 0x69, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x62, 0x69, 0x5f, 0x6d, 0x64, 0x65, 0x5f, 0x6f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcb, 0x01, 0x0a,
	0x1d, 0x50, 0x72, 0x6f, 0x70, 0x61, 0x67, 0x61, 0x74, 0x65, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x43,
	0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x37, 0x0a, 0x08, 0x62, 0x69, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x62, 0x69, 0x2e, 0x6d, 0x64, 0x65, 0x5f, 0x6f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2e, 0x47, 0x42, 0x4d, 0x42, 0x49, 0x4f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x52, 0x08, 0x62, 0x69, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x6f, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x79, 0x42, 0x4e, 0x0a, 0x1b, 0x61, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x62, 0x69, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x79, 0x42, 0x16, 0x50, 0x72, 0x6f, 0x70, 0x61,
	0x67, 0x61, 0x74, 0x65, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x50, 0x01, 0x5a, 0x15, 0x61, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6f, 0x63, 0x66, 0x2f, 0x62, 0x69, 0x2f, 0x76, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_bi_certify_proto_rawDescOnce sync.Once
	file_bi_certify_proto_rawDescData = file_bi_certify_proto_rawDesc
)

func file_bi_certify_proto_rawDescGZIP() []byte {
	file_bi_certify_proto_rawDescOnce.Do(func() {
		file_bi_certify_proto_rawDescData = protoimpl.X.CompressGZIP(file_bi_certify_proto_rawDescData)
	})
	return file_bi_certify_proto_rawDescData
}

var file_bi_certify_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_bi_certify_proto_goTypes = []interface{}{
	(*PropagateCertificationRequest)(nil), // 0: bi.certify.PropagateCertificationRequest
	(*ProtobufConfiguration)(nil),         // 1: common.ProtobufConfiguration
	(*GBMBIObject)(nil),                   // 2: bi.mde_objects.GBMBIObject
}
var file_bi_certify_proto_depIdxs = []int32{
	1, // 0: bi.certify.PropagateCertificationRequest.configuration:type_name -> common.ProtobufConfiguration
	2, // 1: bi.certify.PropagateCertificationRequest.biObject:type_name -> bi.mde_objects.GBMBIObject
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_bi_certify_proto_init() }
func file_bi_certify_proto_init() {
	if File_bi_certify_proto != nil {
		return
	}
	file_bi_configuration_proto_init()
	file_bi_mde_objects_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_bi_certify_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PropagateCertificationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_bi_certify_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_bi_certify_proto_goTypes,
		DependencyIndexes: file_bi_certify_proto_depIdxs,
		MessageInfos:      file_bi_certify_proto_msgTypes,
	}.Build()
	File_bi_certify_proto = out.File
	file_bi_certify_proto_rawDesc = nil
	file_bi_certify_proto_goTypes = nil
	file_bi_certify_proto_depIdxs = nil
}

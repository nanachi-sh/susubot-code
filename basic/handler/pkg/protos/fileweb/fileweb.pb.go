// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.28.2
// source: pkg/protos/fileweb/fileweb.proto

package fileweb

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

type Errors int32

const (
	Errors_EMPTY     Errors = 0
	Errors_Undefined Errors = 1
)

// Enum value maps for Errors.
var (
	Errors_name = map[int32]string{
		0: "EMPTY",
		1: "Undefined",
	}
	Errors_value = map[string]int32{
		"EMPTY":     0,
		"Undefined": 1,
	}
)

func (x Errors) Enum() *Errors {
	p := new(Errors)
	*p = x
	return p
}

func (x Errors) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Errors) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_protos_fileweb_fileweb_proto_enumTypes[0].Descriptor()
}

func (Errors) Type() protoreflect.EnumType {
	return &file_pkg_protos_fileweb_fileweb_proto_enumTypes[0]
}

func (x Errors) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Errors.Descriptor instead.
func (Errors) EnumDescriptor() ([]byte, []int) {
	return file_pkg_protos_fileweb_fileweb_proto_rawDescGZIP(), []int{0}
}

type UploadRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Buf           []byte                 `protobuf:"bytes,1,opt,name=Buf,proto3" json:"Buf,omitempty"`
	ValidTime     *uint32                `protobuf:"varint,2,opt,name=ValidTime,proto3,oneof" json:"ValidTime,omitempty"` //过期时间(ms)
	AutoRefresh   bool                   `protobuf:"varint,3,opt,name=AutoRefresh,proto3" json:"AutoRefresh,omitempty"`   //资源被请求后自动重置过期时间
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadRequest) Reset() {
	*x = UploadRequest{}
	mi := &file_pkg_protos_fileweb_fileweb_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadRequest) ProtoMessage() {}

func (x *UploadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_fileweb_fileweb_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadRequest.ProtoReflect.Descriptor instead.
func (*UploadRequest) Descriptor() ([]byte, []int) {
	return file_pkg_protos_fileweb_fileweb_proto_rawDescGZIP(), []int{0}
}

func (x *UploadRequest) GetBuf() []byte {
	if x != nil {
		return x.Buf
	}
	return nil
}

func (x *UploadRequest) GetValidTime() uint32 {
	if x != nil && x.ValidTime != nil {
		return *x.ValidTime
	}
	return 0
}

func (x *UploadRequest) GetAutoRefresh() bool {
	if x != nil {
		return x.AutoRefresh
	}
	return false
}

type UploadResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Body:
	//
	//	*UploadResponse_Hash
	//	*UploadResponse_Err
	Body          isUploadResponse_Body `protobuf_oneof:"Body"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadResponse) Reset() {
	*x = UploadResponse{}
	mi := &file_pkg_protos_fileweb_fileweb_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadResponse) ProtoMessage() {}

func (x *UploadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_fileweb_fileweb_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadResponse.ProtoReflect.Descriptor instead.
func (*UploadResponse) Descriptor() ([]byte, []int) {
	return file_pkg_protos_fileweb_fileweb_proto_rawDescGZIP(), []int{1}
}

func (x *UploadResponse) GetBody() isUploadResponse_Body {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *UploadResponse) GetHash() string {
	if x != nil {
		if x, ok := x.Body.(*UploadResponse_Hash); ok {
			return x.Hash
		}
	}
	return ""
}

func (x *UploadResponse) GetErr() Errors {
	if x != nil {
		if x, ok := x.Body.(*UploadResponse_Err); ok {
			return x.Err
		}
	}
	return Errors_EMPTY
}

type isUploadResponse_Body interface {
	isUploadResponse_Body()
}

type UploadResponse_Hash struct {
	Hash string `protobuf:"bytes,1,opt,name=Hash,proto3,oneof"`
}

type UploadResponse_Err struct {
	Err Errors `protobuf:"varint,2,opt,name=err,proto3,enum=susubot.basic.fileweb.Errors,oneof"`
}

func (*UploadResponse_Hash) isUploadResponse_Body() {}

func (*UploadResponse_Err) isUploadResponse_Body() {}

var File_pkg_protos_fileweb_fileweb_proto protoreflect.FileDescriptor

var file_pkg_protos_fileweb_fileweb_proto_rawDesc = []byte{
	0x0a, 0x20, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x66, 0x69, 0x6c,
	0x65, 0x77, 0x65, 0x62, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x77, 0x65, 0x62, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x15, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x62, 0x61, 0x73, 0x69,
	0x63, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x77, 0x65, 0x62, 0x22, 0x74, 0x0a, 0x0d, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x42, 0x75,
	0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x42, 0x75, 0x66, 0x12, 0x21, 0x0a, 0x09,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x48,
	0x00, 0x52, 0x09, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12,
	0x20, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x6f, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x41, 0x75, 0x74, 0x6f, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73,
	0x68, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22,
	0x61, 0x0a, 0x0e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x14, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x12, 0x31, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x62,
	0x61, 0x73, 0x69, 0x63, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x77, 0x65, 0x62, 0x2e, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x73, 0x48, 0x00, 0x52, 0x03, 0x65, 0x72, 0x72, 0x42, 0x06, 0x0a, 0x04, 0x42, 0x6f,
	0x64, 0x79, 0x2a, 0x22, 0x0a, 0x06, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x12, 0x09, 0x0a, 0x05,
	0x45, 0x4d, 0x50, 0x54, 0x59, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x6e, 0x64, 0x65, 0x66,
	0x69, 0x6e, 0x65, 0x64, 0x10, 0x01, 0x32, 0x60, 0x0a, 0x07, 0x46, 0x69, 0x6c, 0x65, 0x57, 0x65,
	0x62, 0x12, 0x55, 0x0a, 0x06, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x24, 0x2e, 0x73, 0x75,
	0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2e, 0x66, 0x69, 0x6c, 0x65,
	0x77, 0x65, 0x62, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x25, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x62, 0x61, 0x73, 0x69,
	0x63, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x77, 0x65, 0x62, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x10, 0x5a, 0x0e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x77, 0x65, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_pkg_protos_fileweb_fileweb_proto_rawDescOnce sync.Once
	file_pkg_protos_fileweb_fileweb_proto_rawDescData = file_pkg_protos_fileweb_fileweb_proto_rawDesc
)

func file_pkg_protos_fileweb_fileweb_proto_rawDescGZIP() []byte {
	file_pkg_protos_fileweb_fileweb_proto_rawDescOnce.Do(func() {
		file_pkg_protos_fileweb_fileweb_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_protos_fileweb_fileweb_proto_rawDescData)
	})
	return file_pkg_protos_fileweb_fileweb_proto_rawDescData
}

var file_pkg_protos_fileweb_fileweb_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_protos_fileweb_fileweb_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_protos_fileweb_fileweb_proto_goTypes = []any{
	(Errors)(0),            // 0: susubot.basic.fileweb.Errors
	(*UploadRequest)(nil),  // 1: susubot.basic.fileweb.UploadRequest
	(*UploadResponse)(nil), // 2: susubot.basic.fileweb.UploadResponse
}
var file_pkg_protos_fileweb_fileweb_proto_depIdxs = []int32{
	0, // 0: susubot.basic.fileweb.UploadResponse.err:type_name -> susubot.basic.fileweb.Errors
	1, // 1: susubot.basic.fileweb.FileWeb.Upload:input_type -> susubot.basic.fileweb.UploadRequest
	2, // 2: susubot.basic.fileweb.FileWeb.Upload:output_type -> susubot.basic.fileweb.UploadResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_protos_fileweb_fileweb_proto_init() }
func file_pkg_protos_fileweb_fileweb_proto_init() {
	if File_pkg_protos_fileweb_fileweb_proto != nil {
		return
	}
	file_pkg_protos_fileweb_fileweb_proto_msgTypes[0].OneofWrappers = []any{}
	file_pkg_protos_fileweb_fileweb_proto_msgTypes[1].OneofWrappers = []any{
		(*UploadResponse_Hash)(nil),
		(*UploadResponse_Err)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_protos_fileweb_fileweb_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_protos_fileweb_fileweb_proto_goTypes,
		DependencyIndexes: file_pkg_protos_fileweb_fileweb_proto_depIdxs,
		EnumInfos:         file_pkg_protos_fileweb_fileweb_proto_enumTypes,
		MessageInfos:      file_pkg_protos_fileweb_fileweb_proto_msgTypes,
	}.Build()
	File_pkg_protos_fileweb_fileweb_proto = out.File
	file_pkg_protos_fileweb_fileweb_proto_rawDesc = nil
	file_pkg_protos_fileweb_fileweb_proto_goTypes = nil
	file_pkg_protos_fileweb_fileweb_proto_depIdxs = nil
}

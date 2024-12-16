// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: plugin/randomfortune/LLOneBot/protos/randomfortune/randomfortune.proto

package randomfortune

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

type BasicRequest_ReturnMethods int32

const (
	BasicRequest_Hash BasicRequest_ReturnMethods = 0
	BasicRequest_Raw  BasicRequest_ReturnMethods = 1
)

// Enum value maps for BasicRequest_ReturnMethods.
var (
	BasicRequest_ReturnMethods_name = map[int32]string{
		0: "Hash",
		1: "Raw",
	}
	BasicRequest_ReturnMethods_value = map[string]int32{
		"Hash": 0,
		"Raw":  1,
	}
)

func (x BasicRequest_ReturnMethods) Enum() *BasicRequest_ReturnMethods {
	p := new(BasicRequest_ReturnMethods)
	*p = x
	return p
}

func (x BasicRequest_ReturnMethods) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BasicRequest_ReturnMethods) Descriptor() protoreflect.EnumDescriptor {
	return file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_enumTypes[0].Descriptor()
}

func (BasicRequest_ReturnMethods) Type() protoreflect.EnumType {
	return &file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_enumTypes[0]
}

func (x BasicRequest_ReturnMethods) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BasicRequest_ReturnMethods.Descriptor instead.
func (BasicRequest_ReturnMethods) EnumDescriptor() ([]byte, []int) {
	return file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDescGZIP(), []int{0, 0}
}

type BasicRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReturnMethod BasicRequest_ReturnMethods `protobuf:"varint,1,opt,name=ReturnMethod,proto3,enum=susubot.plugin.randomfortune.BasicRequest_ReturnMethods" json:"ReturnMethod,omitempty"`
}

func (x *BasicRequest) Reset() {
	*x = BasicRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicRequest) ProtoMessage() {}

func (x *BasicRequest) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicRequest.ProtoReflect.Descriptor instead.
func (*BasicRequest) Descriptor() ([]byte, []int) {
	return file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDescGZIP(), []int{0}
}

func (x *BasicRequest) GetReturnMethod() BasicRequest_ReturnMethods {
	if x != nil {
		return x.ReturnMethod
	}
	return BasicRequest_Hash
}

type BasicResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response *BasicResponse_UploadResponse `protobuf:"bytes,1,opt,name=Response,proto3,oneof" json:"Response,omitempty"`
	Buf      []byte                        `protobuf:"bytes,2,opt,name=Buf,proto3,oneof" json:"Buf,omitempty"`
}

func (x *BasicResponse) Reset() {
	*x = BasicResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicResponse) ProtoMessage() {}

func (x *BasicResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicResponse.ProtoReflect.Descriptor instead.
func (*BasicResponse) Descriptor() ([]byte, []int) {
	return file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDescGZIP(), []int{1}
}

func (x *BasicResponse) GetResponse() *BasicResponse_UploadResponse {
	if x != nil {
		return x.Response
	}
	return nil
}

func (x *BasicResponse) GetBuf() []byte {
	if x != nil {
		return x.Buf
	}
	return nil
}

type BasicResponse_UploadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash    string `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	URLPath string `protobuf:"bytes,2,opt,name=URLPath,proto3" json:"URLPath,omitempty"`
}

func (x *BasicResponse_UploadResponse) Reset() {
	*x = BasicResponse_UploadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicResponse_UploadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicResponse_UploadResponse) ProtoMessage() {}

func (x *BasicResponse_UploadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BasicResponse_UploadResponse.ProtoReflect.Descriptor instead.
func (*BasicResponse_UploadResponse) Descriptor() ([]byte, []int) {
	return file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDescGZIP(), []int{1, 0}
}

func (x *BasicResponse_UploadResponse) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *BasicResponse_UploadResponse) GetURLPath() string {
	if x != nil {
		return x.URLPath
	}
	return ""
}

var File_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto protoreflect.FileDescriptor

var file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDesc = []byte{
	0x0a, 0x46, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x66,
	0x6f, 0x72, 0x74, 0x75, 0x6e, 0x65, 0x2f, 0x4c, 0x4c, 0x4f, 0x6e, 0x65, 0x42, 0x6f, 0x74, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x66, 0x6f, 0x72,
	0x74, 0x75, 0x6e, 0x65, 0x2f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x66, 0x6f, 0x72, 0x74, 0x75,
	0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1c, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f,
	0x74, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x66,
	0x6f, 0x72, 0x74, 0x75, 0x6e, 0x65, 0x22, 0x90, 0x01, 0x0a, 0x0c, 0x42, 0x61, 0x73, 0x69, 0x63,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x5c, 0x0a, 0x0c, 0x52, 0x65, 0x74, 0x75, 0x72,
	0x6e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x38, 0x2e,
	0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72,
	0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x66, 0x6f, 0x72, 0x74, 0x75, 0x6e, 0x65, 0x2e, 0x42, 0x61, 0x73,
	0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e,
	0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x52, 0x0c, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x4d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x22, 0x22, 0x0a, 0x0d, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x4d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x12, 0x08, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x10, 0x00,
	0x12, 0x07, 0x0a, 0x03, 0x52, 0x61, 0x77, 0x10, 0x01, 0x22, 0xd8, 0x01, 0x0a, 0x0d, 0x42, 0x61,
	0x73, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5b, 0x0a, 0x08, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x3a, 0x2e,
	0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72,
	0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x66, 0x6f, 0x72, 0x74, 0x75, 0x6e, 0x65, 0x2e, 0x42, 0x61, 0x73,
	0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x08, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x88, 0x01, 0x01, 0x12, 0x15, 0x0a, 0x03, 0x42, 0x75, 0x66, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x01, 0x52, 0x03, 0x42, 0x75, 0x66, 0x88, 0x01, 0x01, 0x1a,
	0x3e, 0x0a, 0x0e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x48, 0x61, 0x73, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x55, 0x52, 0x4c, 0x50, 0x61, 0x74, 0x68,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x55, 0x52, 0x4c, 0x50, 0x61, 0x74, 0x68, 0x42,
	0x0b, 0x0a, 0x09, 0x5f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x06, 0x0a, 0x04,
	0x5f, 0x42, 0x75, 0x66, 0x32, 0x76, 0x0a, 0x0d, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x46, 0x6f,
	0x72, 0x74, 0x75, 0x6e, 0x65, 0x12, 0x65, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x72, 0x74,
	0x75, 0x6e, 0x65, 0x12, 0x2a, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x66, 0x6f, 0x72, 0x74, 0x75,
	0x6e, 0x65, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2b, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x66, 0x6f, 0x72, 0x74, 0x75, 0x6e, 0x65, 0x2e, 0x42,
	0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x16, 0x5a, 0x14,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x66, 0x6f, 0x72,
	0x74, 0x75, 0x6e, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDescOnce sync.Once
	file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDescData = file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDesc
)

func file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDescGZIP() []byte {
	file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDescOnce.Do(func() {
		file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDescData = protoimpl.X.CompressGZIP(file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDescData)
	})
	return file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDescData
}

var file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_goTypes = []any{
	(BasicRequest_ReturnMethods)(0),      // 0: susubot.plugin.randomfortune.BasicRequest.ReturnMethods
	(*BasicRequest)(nil),                 // 1: susubot.plugin.randomfortune.BasicRequest
	(*BasicResponse)(nil),                // 2: susubot.plugin.randomfortune.BasicResponse
	(*BasicResponse_UploadResponse)(nil), // 3: susubot.plugin.randomfortune.BasicResponse.UploadResponse
}
var file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_depIdxs = []int32{
	0, // 0: susubot.plugin.randomfortune.BasicRequest.ReturnMethod:type_name -> susubot.plugin.randomfortune.BasicRequest.ReturnMethods
	3, // 1: susubot.plugin.randomfortune.BasicResponse.Response:type_name -> susubot.plugin.randomfortune.BasicResponse.UploadResponse
	1, // 2: susubot.plugin.randomfortune.randomFortune.GetFortune:input_type -> susubot.plugin.randomfortune.BasicRequest
	2, // 3: susubot.plugin.randomfortune.randomFortune.GetFortune:output_type -> susubot.plugin.randomfortune.BasicResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_init() }
func file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_init() {
	if File_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*BasicRequest); i {
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
		file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*BasicResponse); i {
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
		file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*BasicResponse_UploadResponse); i {
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
	file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_goTypes,
		DependencyIndexes: file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_depIdxs,
		EnumInfos:         file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_enumTypes,
		MessageInfos:      file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_msgTypes,
	}.Build()
	File_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto = out.File
	file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_rawDesc = nil
	file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_goTypes = nil
	file_plugin_randomfortune_LLOneBot_protos_randomfortune_randomfortune_proto_depIdxs = nil
}
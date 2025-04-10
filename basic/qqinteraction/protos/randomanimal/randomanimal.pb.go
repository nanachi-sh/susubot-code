// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: plugin/randomanimal/protos/randomanimal/randomanimal.proto

package randomanimal

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

type Type int32

const (
	Type_Image Type = 0
	Type_Video Type = 1
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "Image",
		1: "Video",
	}
	Type_value = map[string]int32{
		"Image": 0,
		"Video": 1,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDescGZIP(), []int{0}
}

type BasicRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AutoUpload bool `protobuf:"varint,1,opt,name=AutoUpload,proto3" json:"AutoUpload,omitempty"`
}

func (x *BasicRequest) Reset() {
	*x = BasicRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicRequest) ProtoMessage() {}

func (x *BasicRequest) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_msgTypes[0]
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
	return file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDescGZIP(), []int{0}
}

func (x *BasicRequest) GetAutoUpload() bool {
	if x != nil {
		return x.AutoUpload
	}
	return false
}

type BasicResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type     Type                          `protobuf:"varint,1,opt,name=Type,proto3,enum=susubot.plugin.randomanimal.Type" json:"Type,omitempty"`
	Buf      []byte                        `protobuf:"bytes,2,opt,name=Buf,proto3,oneof" json:"Buf,omitempty"`
	Response *BasicResponse_UploadResponse `protobuf:"bytes,3,opt,name=Response,proto3,oneof" json:"Response,omitempty"`
}

func (x *BasicResponse) Reset() {
	*x = BasicResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicResponse) ProtoMessage() {}

func (x *BasicResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_msgTypes[1]
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
	return file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDescGZIP(), []int{1}
}

func (x *BasicResponse) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_Image
}

func (x *BasicResponse) GetBuf() []byte {
	if x != nil {
		return x.Buf
	}
	return nil
}

func (x *BasicResponse) GetResponse() *BasicResponse_UploadResponse {
	if x != nil {
		return x.Response
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
		mi := &file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BasicResponse_UploadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BasicResponse_UploadResponse) ProtoMessage() {}

func (x *BasicResponse_UploadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_msgTypes[2]
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
	return file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDescGZIP(), []int{1, 0}
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

var File_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto protoreflect.FileDescriptor

var file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDesc = []byte{
	0x0a, 0x43, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x61,
	0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2f, 0x4c, 0x4c, 0x4f, 0x6e, 0x65, 0x42, 0x6f, 0x74, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x61, 0x6e, 0x69, 0x6d,
	0x61, 0x6c, 0x2f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1b, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x61, 0x6e, 0x69, 0x6d,
	0x61, 0x6c, 0x22, 0x2e, 0x0a, 0x0c, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x41, 0x75, 0x74, 0x6f, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x41, 0x75, 0x74, 0x6f, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x22, 0x8e, 0x02, 0x0a, 0x0d, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x21, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c,
	0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x15, 0x0a, 0x03, 0x42,
	0x75, 0x66, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x03, 0x42, 0x75, 0x66, 0x88,
	0x01, 0x01, 0x12, 0x5a, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x39, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x61, 0x6e, 0x69, 0x6d,
	0x61, 0x6c, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48,
	0x01, 0x52, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x88, 0x01, 0x01, 0x1a, 0x3e,
	0x0a, 0x0e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x55, 0x52, 0x4c, 0x50, 0x61, 0x74, 0x68, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x55, 0x52, 0x4c, 0x50, 0x61, 0x74, 0x68, 0x42, 0x06,
	0x0a, 0x04, 0x5f, 0x42, 0x75, 0x66, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x2a, 0x1c, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x10,
	0x01, 0x32, 0xfb, 0x03, 0x0a, 0x0c, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x41, 0x6e, 0x69, 0x6d,
	0x61, 0x6c, 0x12, 0x5f, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x44, 0x6f, 0x67, 0x12, 0x29, 0x2e, 0x73,
	0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72, 0x61,
	0x6e, 0x64, 0x6f, 0x6d, 0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f,
	0x74, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x61,
	0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x5f, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x78, 0x12, 0x29, 0x2e,
	0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72,
	0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x42, 0x61, 0x73, 0x69,
	0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62,
	0x6f, 0x74, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d,
	0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x60, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x44, 0x75, 0x63, 0x6b, 0x12,
	0x29, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x42, 0x61,
	0x73, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x73, 0x75, 0x73,
	0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72, 0x61, 0x6e, 0x64,
	0x6f, 0x6d, 0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5f, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x43, 0x61, 0x74,
	0x12, 0x29, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x42,
	0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x73, 0x75,
	0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72, 0x61, 0x6e,
	0x64, 0x6f, 0x6d, 0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x66, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x43, 0x68,
	0x69, 0x6b, 0x65, 0x6e, 0x5f, 0x43, 0x58, 0x4b, 0x12, 0x29, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62,
	0x6f, 0x74, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d,
	0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x61, 0x6e, 0x69, 0x6d, 0x61,
	0x6c, 0x2e, 0x42, 0x61, 0x73, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x15, 0x5a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d,
	0x61, 0x6e, 0x69, 0x6d, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDescOnce sync.Once
	file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDescData = file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDesc
)

func file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDescGZIP() []byte {
	file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDescOnce.Do(func() {
		file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDescData = protoimpl.X.CompressGZIP(file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDescData)
	})
	return file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDescData
}

var file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_goTypes = []any{
	(Type)(0),                            // 0: susubot.plugin.randomanimal.Type
	(*BasicRequest)(nil),                 // 1: susubot.plugin.randomanimal.BasicRequest
	(*BasicResponse)(nil),                // 2: susubot.plugin.randomanimal.BasicResponse
	(*BasicResponse_UploadResponse)(nil), // 3: susubot.plugin.randomanimal.BasicResponse.UploadResponse
}
var file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_depIdxs = []int32{
	0, // 0: susubot.plugin.randomanimal.BasicResponse.Type:type_name -> susubot.plugin.randomanimal.Type
	3, // 1: susubot.plugin.randomanimal.BasicResponse.Response:type_name -> susubot.plugin.randomanimal.BasicResponse.UploadResponse
	1, // 2: susubot.plugin.randomanimal.randomAnimal.GetDog:input_type -> susubot.plugin.randomanimal.BasicRequest
	1, // 3: susubot.plugin.randomanimal.randomAnimal.GetFox:input_type -> susubot.plugin.randomanimal.BasicRequest
	1, // 4: susubot.plugin.randomanimal.randomAnimal.GetDuck:input_type -> susubot.plugin.randomanimal.BasicRequest
	1, // 5: susubot.plugin.randomanimal.randomAnimal.GetCat:input_type -> susubot.plugin.randomanimal.BasicRequest
	1, // 6: susubot.plugin.randomanimal.randomAnimal.GetChiken_CXK:input_type -> susubot.plugin.randomanimal.BasicRequest
	2, // 7: susubot.plugin.randomanimal.randomAnimal.GetDog:output_type -> susubot.plugin.randomanimal.BasicResponse
	2, // 8: susubot.plugin.randomanimal.randomAnimal.GetFox:output_type -> susubot.plugin.randomanimal.BasicResponse
	2, // 9: susubot.plugin.randomanimal.randomAnimal.GetDuck:output_type -> susubot.plugin.randomanimal.BasicResponse
	2, // 10: susubot.plugin.randomanimal.randomAnimal.GetCat:output_type -> susubot.plugin.randomanimal.BasicResponse
	2, // 11: susubot.plugin.randomanimal.randomAnimal.GetChiken_CXK:output_type -> susubot.plugin.randomanimal.BasicResponse
	7, // [7:12] is the sub-list for method output_type
	2, // [2:7] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_init() }
func file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_init() {
	if File_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_msgTypes[0].Exporter = func(v any, i int) any {
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
		file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_msgTypes[1].Exporter = func(v any, i int) any {
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
		file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_msgTypes[2].Exporter = func(v any, i int) any {
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
	file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_goTypes,
		DependencyIndexes: file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_depIdxs,
		EnumInfos:         file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_enumTypes,
		MessageInfos:      file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_msgTypes,
	}.Build()
	File_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto = out.File
	file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_rawDesc = nil
	file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_goTypes = nil
	file_plugin_randomanimal_LLOneBot_protos_randomanimal_randomanimal_proto_depIdxs = nil
}

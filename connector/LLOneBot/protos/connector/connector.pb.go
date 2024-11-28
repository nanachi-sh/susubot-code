// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: connector/LLOneBot/protos/connector/connector.proto

package connector

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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_connector_LLOneBot_protos_connector_connector_proto_rawDescGZIP(), []int{0}
}

type ConnectRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr  string  `protobuf:"bytes,1,opt,name=Addr,proto3" json:"Addr,omitempty"`
	Port  int32   `protobuf:"varint,2,opt,name=Port,proto3" json:"Port,omitempty"`
	Token *string `protobuf:"bytes,3,opt,name=Token,proto3,oneof" json:"Token,omitempty"`
}

func (x *ConnectRequest) Reset() {
	*x = ConnectRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectRequest) ProtoMessage() {}

func (x *ConnectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectRequest.ProtoReflect.Descriptor instead.
func (*ConnectRequest) Descriptor() ([]byte, []int) {
	return file_connector_LLOneBot_protos_connector_connector_proto_rawDescGZIP(), []int{1}
}

func (x *ConnectRequest) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *ConnectRequest) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *ConnectRequest) GetToken() string {
	if x != nil && x.Token != nil {
		return *x.Token
	}
	return ""
}

type ConnectResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OK bool `protobuf:"varint,1,opt,name=OK,proto3" json:"OK,omitempty"`
}

func (x *ConnectResponse) Reset() {
	*x = ConnectResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectResponse) ProtoMessage() {}

func (x *ConnectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectResponse.ProtoReflect.Descriptor instead.
func (*ConnectResponse) Descriptor() ([]byte, []int) {
	return file_connector_LLOneBot_protos_connector_connector_proto_rawDescGZIP(), []int{2}
}

func (x *ConnectResponse) GetOK() bool {
	if x != nil {
		return x.OK
	}
	return false
}

type ReadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsClose bool   `protobuf:"varint,1,opt,name=isClose,proto3" json:"isClose,omitempty"`
	Buf     []byte `protobuf:"bytes,2,opt,name=buf,proto3" json:"buf,omitempty"`
}

func (x *ReadResponse) Reset() {
	*x = ReadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadResponse) ProtoMessage() {}

func (x *ReadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadResponse.ProtoReflect.Descriptor instead.
func (*ReadResponse) Descriptor() ([]byte, []int) {
	return file_connector_LLOneBot_protos_connector_connector_proto_rawDescGZIP(), []int{3}
}

func (x *ReadResponse) GetIsClose() bool {
	if x != nil {
		return x.IsClose
	}
	return false
}

func (x *ReadResponse) GetBuf() []byte {
	if x != nil {
		return x.Buf
	}
	return nil
}

type WriteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Buf []byte `protobuf:"bytes,1,opt,name=buf,proto3" json:"buf,omitempty"`
}

func (x *WriteRequest) Reset() {
	*x = WriteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteRequest) ProtoMessage() {}

func (x *WriteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteRequest.ProtoReflect.Descriptor instead.
func (*WriteRequest) Descriptor() ([]byte, []int) {
	return file_connector_LLOneBot_protos_connector_connector_proto_rawDescGZIP(), []int{4}
}

func (x *WriteRequest) GetBuf() []byte {
	if x != nil {
		return x.Buf
	}
	return nil
}

type CloseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OK bool `protobuf:"varint,1,opt,name=OK,proto3" json:"OK,omitempty"`
}

func (x *CloseResponse) Reset() {
	*x = CloseResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CloseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloseResponse) ProtoMessage() {}

func (x *CloseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloseResponse.ProtoReflect.Descriptor instead.
func (*CloseResponse) Descriptor() ([]byte, []int) {
	return file_connector_LLOneBot_protos_connector_connector_proto_rawDescGZIP(), []int{5}
}

func (x *CloseResponse) GetOK() bool {
	if x != nil {
		return x.OK
	}
	return false
}

var File_connector_LLOneBot_protos_connector_connector_proto protoreflect.FileDescriptor

var file_connector_LLOneBot_protos_connector_connector_proto_rawDesc = []byte{
	0x0a, 0x33, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x4c, 0x4c, 0x4f, 0x6e,
	0x65, 0x42, 0x6f, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x63, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x5d,
	0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x41, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x41, 0x64, 0x64, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x50, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x19, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x21, 0x0a,
	0x0f, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x4f, 0x4b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x4f, 0x4b,
	0x22, 0x3a, 0x0a, 0x0c, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x69, 0x73, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x69, 0x73, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x62, 0x75,
	0x66, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x62, 0x75, 0x66, 0x22, 0x20, 0x0a, 0x0c,
	0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x62, 0x75, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x62, 0x75, 0x66, 0x22, 0x1f,
	0x0a, 0x0d, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x4f, 0x4b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x4f, 0x4b, 0x32,
	0x9b, 0x01, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x2c, 0x0a,
	0x07, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x0f, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x43, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x04, 0x52,
	0x65, 0x61, 0x64, 0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x52, 0x65,
	0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x1e, 0x0a, 0x05,
	0x57, 0x72, 0x69, 0x74, 0x65, 0x12, 0x0d, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x1f, 0x0a, 0x05,
	0x43, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0e, 0x2e,
	0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x48, 0x5a,
	0x46, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x61, 0x6e, 0x61,
	0x63, 0x68, 0x69, 0x2d, 0x73, 0x68, 0x2f, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2d, 0x63,
	0x6f, 0x64, 0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2f, 0x4c, 0x4c,
	0x4f, 0x6e, 0x65, 0x42, 0x6f, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x63, 0x6f,
	0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_connector_LLOneBot_protos_connector_connector_proto_rawDescOnce sync.Once
	file_connector_LLOneBot_protos_connector_connector_proto_rawDescData = file_connector_LLOneBot_protos_connector_connector_proto_rawDesc
)

func file_connector_LLOneBot_protos_connector_connector_proto_rawDescGZIP() []byte {
	file_connector_LLOneBot_protos_connector_connector_proto_rawDescOnce.Do(func() {
		file_connector_LLOneBot_protos_connector_connector_proto_rawDescData = protoimpl.X.CompressGZIP(file_connector_LLOneBot_protos_connector_connector_proto_rawDescData)
	})
	return file_connector_LLOneBot_protos_connector_connector_proto_rawDescData
}

var file_connector_LLOneBot_protos_connector_connector_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_connector_LLOneBot_protos_connector_connector_proto_goTypes = []any{
	(*Empty)(nil),           // 0: Empty
	(*ConnectRequest)(nil),  // 1: ConnectRequest
	(*ConnectResponse)(nil), // 2: ConnectResponse
	(*ReadResponse)(nil),    // 3: ReadResponse
	(*WriteRequest)(nil),    // 4: WriteRequest
	(*CloseResponse)(nil),   // 5: CloseResponse
}
var file_connector_LLOneBot_protos_connector_connector_proto_depIdxs = []int32{
	1, // 0: Connector.Connect:input_type -> ConnectRequest
	0, // 1: Connector.Read:input_type -> Empty
	4, // 2: Connector.Write:input_type -> WriteRequest
	0, // 3: Connector.Close:input_type -> Empty
	2, // 4: Connector.Connect:output_type -> ConnectResponse
	3, // 5: Connector.Read:output_type -> ReadResponse
	0, // 6: Connector.Write:output_type -> Empty
	5, // 7: Connector.Close:output_type -> CloseResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_connector_LLOneBot_protos_connector_connector_proto_init() }
func file_connector_LLOneBot_protos_connector_connector_proto_init() {
	if File_connector_LLOneBot_protos_connector_connector_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Empty); i {
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
		file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ConnectRequest); i {
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
		file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*ConnectResponse); i {
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
		file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ReadResponse); i {
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
		file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*WriteRequest); i {
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
		file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*CloseResponse); i {
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
	file_connector_LLOneBot_protos_connector_connector_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_connector_LLOneBot_protos_connector_connector_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_connector_LLOneBot_protos_connector_connector_proto_goTypes,
		DependencyIndexes: file_connector_LLOneBot_protos_connector_connector_proto_depIdxs,
		MessageInfos:      file_connector_LLOneBot_protos_connector_connector_proto_msgTypes,
	}.Build()
	File_connector_LLOneBot_protos_connector_connector_proto = out.File
	file_connector_LLOneBot_protos_connector_connector_proto_rawDesc = nil
	file_connector_LLOneBot_protos_connector_connector_proto_goTypes = nil
	file_connector_LLOneBot_protos_connector_connector_proto_depIdxs = nil
}
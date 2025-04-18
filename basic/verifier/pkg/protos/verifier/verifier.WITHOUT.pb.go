// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.28.2
// source: pkg/protos/verifier/verifier.WITHOUT.proto

package verifier

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
	Errors_EMPTY         Errors = 0
	Errors_Intervaling   Errors = 1 //正处于间隔时间内
	Errors_NoFriend      Errors = 2 //请求QQ号非机器人好友
	Errors_ErrVerified   Errors = 3 //已验证(重复验证)
	Errors_VerifyNoFound Errors = 4 //未找到验证请求
	Errors_Expired       Errors = 5 //已过期
	Errors_CodeWrong     Errors = 6 //验证码错误
	Errors_UnVerified    Errors = 7 //还未验证
	Errors_Undefined     Errors = 8
)

// Enum value maps for Errors.
var (
	Errors_name = map[int32]string{
		0: "EMPTY",
		1: "Intervaling",
		2: "NoFriend",
		3: "ErrVerified",
		4: "VerifyNoFound",
		5: "Expired",
		6: "CodeWrong",
		7: "UnVerified",
		8: "Undefined",
	}
	Errors_value = map[string]int32{
		"EMPTY":         0,
		"Intervaling":   1,
		"NoFriend":      2,
		"ErrVerified":   3,
		"VerifyNoFound": 4,
		"Expired":       5,
		"CodeWrong":     6,
		"UnVerified":    7,
		"Undefined":     8,
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
	return file_pkg_protos_verifier_verifier_WITHOUT_proto_enumTypes[0].Descriptor()
}

func (Errors) Type() protoreflect.EnumType {
	return &file_pkg_protos_verifier_verifier_WITHOUT_proto_enumTypes[0]
}

func (x Errors) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Errors.Descriptor instead.
func (Errors) EnumDescriptor() ([]byte, []int) {
	return file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescGZIP(), []int{0}
}

type Result int32

const (
	Result_Verified Result = 0 //已验证
)

// Enum value maps for Result.
var (
	Result_name = map[int32]string{
		0: "Verified",
	}
	Result_value = map[string]int32{
		"Verified": 0,
	}
)

func (x Result) Enum() *Result {
	p := new(Result)
	*p = x
	return p
}

func (x Result) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Result) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_protos_verifier_verifier_WITHOUT_proto_enumTypes[1].Descriptor()
}

func (Result) Type() protoreflect.EnumType {
	return &file_pkg_protos_verifier_verifier_WITHOUT_proto_enumTypes[1]
}

func (x Result) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Result.Descriptor instead.
func (Result) EnumDescriptor() ([]byte, []int) {
	return file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescGZIP(), []int{1}
}

type QQ_NewVerifyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Interval      int32                  `protobuf:"varint,1,opt,name=interval,proto3" json:"interval,omitempty"` //请求间隔时间(ms)
	Expires       int32                  `protobuf:"varint,2,opt,name=expires,proto3" json:"expires,omitempty"`   //验证码过期时间
	QQID          string                 `protobuf:"bytes,3,opt,name=QQID,proto3" json:"QQID,omitempty"`          //QQ号
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QQ_NewVerifyRequest) Reset() {
	*x = QQ_NewVerifyRequest{}
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QQ_NewVerifyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QQ_NewVerifyRequest) ProtoMessage() {}

func (x *QQ_NewVerifyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QQ_NewVerifyRequest.ProtoReflect.Descriptor instead.
func (*QQ_NewVerifyRequest) Descriptor() ([]byte, []int) {
	return file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescGZIP(), []int{0}
}

func (x *QQ_NewVerifyRequest) GetInterval() int32 {
	if x != nil {
		return x.Interval
	}
	return 0
}

func (x *QQ_NewVerifyRequest) GetExpires() int32 {
	if x != nil {
		return x.Expires
	}
	return 0
}

func (x *QQ_NewVerifyRequest) GetQQID() string {
	if x != nil {
		return x.QQID
	}
	return ""
}

type QQ_NewVerifyResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Body:
	//
	//	*QQ_NewVerifyResponse_Err
	//	*QQ_NewVerifyResponse_VerifyHash
	Body          isQQ_NewVerifyResponse_Body `protobuf_oneof:"Body"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QQ_NewVerifyResponse) Reset() {
	*x = QQ_NewVerifyResponse{}
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QQ_NewVerifyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QQ_NewVerifyResponse) ProtoMessage() {}

func (x *QQ_NewVerifyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QQ_NewVerifyResponse.ProtoReflect.Descriptor instead.
func (*QQ_NewVerifyResponse) Descriptor() ([]byte, []int) {
	return file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescGZIP(), []int{1}
}

func (x *QQ_NewVerifyResponse) GetBody() isQQ_NewVerifyResponse_Body {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *QQ_NewVerifyResponse) GetErr() Errors {
	if x != nil {
		if x, ok := x.Body.(*QQ_NewVerifyResponse_Err); ok {
			return x.Err
		}
	}
	return Errors_EMPTY
}

func (x *QQ_NewVerifyResponse) GetVerifyHash() string {
	if x != nil {
		if x, ok := x.Body.(*QQ_NewVerifyResponse_VerifyHash); ok {
			return x.VerifyHash
		}
	}
	return ""
}

type isQQ_NewVerifyResponse_Body interface {
	isQQ_NewVerifyResponse_Body()
}

type QQ_NewVerifyResponse_Err struct {
	Err Errors `protobuf:"varint,1,opt,name=err,proto3,enum=susubot.basic.verifier.Errors,oneof"`
}

type QQ_NewVerifyResponse_VerifyHash struct {
	VerifyHash string `protobuf:"bytes,2,opt,name=VerifyHash,proto3,oneof"`
}

func (*QQ_NewVerifyResponse_Err) isQQ_NewVerifyResponse_Body() {}

func (*QQ_NewVerifyResponse_VerifyHash) isQQ_NewVerifyResponse_Body() {}

type QQ_VerifyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	VerifyHash    string                 `protobuf:"bytes,1,opt,name=VerifyHash,proto3" json:"VerifyHash,omitempty"`
	VerifyCode    string                 `protobuf:"bytes,2,opt,name=VerifyCode,proto3" json:"VerifyCode,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QQ_VerifyRequest) Reset() {
	*x = QQ_VerifyRequest{}
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QQ_VerifyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QQ_VerifyRequest) ProtoMessage() {}

func (x *QQ_VerifyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QQ_VerifyRequest.ProtoReflect.Descriptor instead.
func (*QQ_VerifyRequest) Descriptor() ([]byte, []int) {
	return file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescGZIP(), []int{2}
}

func (x *QQ_VerifyRequest) GetVerifyHash() string {
	if x != nil {
		return x.VerifyHash
	}
	return ""
}

func (x *QQ_VerifyRequest) GetVerifyCode() string {
	if x != nil {
		return x.VerifyCode
	}
	return ""
}

type QQ_VerifyResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Body:
	//
	//	*QQ_VerifyResponse_Err
	//	*QQ_VerifyResponse_Resp
	Body          isQQ_VerifyResponse_Body `protobuf_oneof:"Body"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QQ_VerifyResponse) Reset() {
	*x = QQ_VerifyResponse{}
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QQ_VerifyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QQ_VerifyResponse) ProtoMessage() {}

func (x *QQ_VerifyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QQ_VerifyResponse.ProtoReflect.Descriptor instead.
func (*QQ_VerifyResponse) Descriptor() ([]byte, []int) {
	return file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescGZIP(), []int{3}
}

func (x *QQ_VerifyResponse) GetBody() isQQ_VerifyResponse_Body {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *QQ_VerifyResponse) GetErr() Errors {
	if x != nil {
		if x, ok := x.Body.(*QQ_VerifyResponse_Err); ok {
			return x.Err
		}
	}
	return Errors_EMPTY
}

func (x *QQ_VerifyResponse) GetResp() *QQ_VerifyResponse_Response {
	if x != nil {
		if x, ok := x.Body.(*QQ_VerifyResponse_Resp); ok {
			return x.Resp
		}
	}
	return nil
}

type isQQ_VerifyResponse_Body interface {
	isQQ_VerifyResponse_Body()
}

type QQ_VerifyResponse_Err struct {
	Err Errors `protobuf:"varint,1,opt,name=err,proto3,enum=susubot.basic.verifier.Errors,oneof"`
}

type QQ_VerifyResponse_Resp struct {
	Resp *QQ_VerifyResponse_Response `protobuf:"bytes,2,opt,name=resp,proto3,oneof"`
}

func (*QQ_VerifyResponse_Err) isQQ_VerifyResponse_Body() {}

func (*QQ_VerifyResponse_Resp) isQQ_VerifyResponse_Body() {}

type QQ_VerifiedRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	VerifyHash    string                 `protobuf:"bytes,1,opt,name=VerifyHash,proto3" json:"VerifyHash,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QQ_VerifiedRequest) Reset() {
	*x = QQ_VerifiedRequest{}
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QQ_VerifiedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QQ_VerifiedRequest) ProtoMessage() {}

func (x *QQ_VerifiedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QQ_VerifiedRequest.ProtoReflect.Descriptor instead.
func (*QQ_VerifiedRequest) Descriptor() ([]byte, []int) {
	return file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescGZIP(), []int{4}
}

func (x *QQ_VerifiedRequest) GetVerifyHash() string {
	if x != nil {
		return x.VerifyHash
	}
	return ""
}

type QQ_VerifiedResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Body:
	//
	//	*QQ_VerifiedResponse_Err
	//	*QQ_VerifiedResponse_Resp
	Body          isQQ_VerifiedResponse_Body `protobuf_oneof:"Body"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QQ_VerifiedResponse) Reset() {
	*x = QQ_VerifiedResponse{}
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QQ_VerifiedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QQ_VerifiedResponse) ProtoMessage() {}

func (x *QQ_VerifiedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QQ_VerifiedResponse.ProtoReflect.Descriptor instead.
func (*QQ_VerifiedResponse) Descriptor() ([]byte, []int) {
	return file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescGZIP(), []int{5}
}

func (x *QQ_VerifiedResponse) GetBody() isQQ_VerifiedResponse_Body {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *QQ_VerifiedResponse) GetErr() Errors {
	if x != nil {
		if x, ok := x.Body.(*QQ_VerifiedResponse_Err); ok {
			return x.Err
		}
	}
	return Errors_EMPTY
}

func (x *QQ_VerifiedResponse) GetResp() *QQ_VerifiedResponse_Response {
	if x != nil {
		if x, ok := x.Body.(*QQ_VerifiedResponse_Resp); ok {
			return x.Resp
		}
	}
	return nil
}

type isQQ_VerifiedResponse_Body interface {
	isQQ_VerifiedResponse_Body()
}

type QQ_VerifiedResponse_Err struct {
	Err Errors `protobuf:"varint,1,opt,name=err,proto3,enum=susubot.basic.verifier.Errors,oneof"`
}

type QQ_VerifiedResponse_Resp struct {
	Resp *QQ_VerifiedResponse_Response `protobuf:"bytes,2,opt,name=resp,proto3,oneof"`
}

func (*QQ_VerifiedResponse_Err) isQQ_VerifiedResponse_Body() {}

func (*QQ_VerifiedResponse_Resp) isQQ_VerifiedResponse_Body() {}

type QQ_VerifyResponse_Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Result        Result                 `protobuf:"varint,2,opt,name=result,proto3,enum=susubot.basic.verifier.Result" json:"result,omitempty"`
	VarifyId      string                 `protobuf:"bytes,3,opt,name=VarifyId,proto3" json:"VarifyId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QQ_VerifyResponse_Response) Reset() {
	*x = QQ_VerifyResponse_Response{}
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QQ_VerifyResponse_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QQ_VerifyResponse_Response) ProtoMessage() {}

func (x *QQ_VerifyResponse_Response) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QQ_VerifyResponse_Response.ProtoReflect.Descriptor instead.
func (*QQ_VerifyResponse_Response) Descriptor() ([]byte, []int) {
	return file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescGZIP(), []int{3, 0}
}

func (x *QQ_VerifyResponse_Response) GetResult() Result {
	if x != nil {
		return x.Result
	}
	return Result_Verified
}

func (x *QQ_VerifyResponse_Response) GetVarifyId() string {
	if x != nil {
		return x.VarifyId
	}
	return ""
}

type QQ_VerifiedResponse_Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Result        Result                 `protobuf:"varint,2,opt,name=result,proto3,enum=susubot.basic.verifier.Result" json:"result,omitempty"`
	VarifyId      string                 `protobuf:"bytes,3,opt,name=VarifyId,proto3" json:"VarifyId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QQ_VerifiedResponse_Response) Reset() {
	*x = QQ_VerifiedResponse_Response{}
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QQ_VerifiedResponse_Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QQ_VerifiedResponse_Response) ProtoMessage() {}

func (x *QQ_VerifiedResponse_Response) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QQ_VerifiedResponse_Response.ProtoReflect.Descriptor instead.
func (*QQ_VerifiedResponse_Response) Descriptor() ([]byte, []int) {
	return file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescGZIP(), []int{5, 0}
}

func (x *QQ_VerifiedResponse_Response) GetResult() Result {
	if x != nil {
		return x.Result
	}
	return Result_Verified
}

func (x *QQ_VerifiedResponse_Response) GetVarifyId() string {
	if x != nil {
		return x.VarifyId
	}
	return ""
}

var File_pkg_protos_verifier_verifier_WITHOUT_proto protoreflect.FileDescriptor

var file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x76, 0x65, 0x72,
	0x69, 0x66, 0x69, 0x65, 0x72, 0x2f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e, 0x57,
	0x49, 0x54, 0x48, 0x4f, 0x55, 0x54, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x73, 0x75,
	0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2e, 0x76, 0x65, 0x72, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x22, 0x5f, 0x0a, 0x13, 0x51, 0x51, 0x5f, 0x4e, 0x65, 0x77, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x78, 0x70, 0x69, 0x72,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65,
	0x73, 0x12, 0x12, 0x0a, 0x04, 0x51, 0x51, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x51, 0x51, 0x49, 0x44, 0x22, 0x74, 0x0a, 0x14, 0x51, 0x51, 0x5f, 0x4e, 0x65, 0x77, 0x56,
	0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a,
	0x03, 0x65, 0x72, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x73, 0x75, 0x73,
	0x75, 0x62, 0x6f, 0x74, 0x2e, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x48, 0x00, 0x52, 0x03, 0x65, 0x72,
	0x72, 0x12, 0x20, 0x0a, 0x0a, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x48, 0x61, 0x73, 0x68, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0a, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x48,
	0x61, 0x73, 0x68, 0x42, 0x06, 0x0a, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x22, 0x52, 0x0a, 0x10, 0x51,
	0x51, 0x5f, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1e, 0x0a, 0x0a, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x1e, 0x0a, 0x0a, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x43, 0x6f, 0x64, 0x65, 0x22,
	0xf9, 0x01, 0x0a, 0x11, 0x51, 0x51, 0x5f, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x62, 0x61, 0x73,
	0x69, 0x63, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x73, 0x48, 0x00, 0x52, 0x03, 0x65, 0x72, 0x72, 0x12, 0x48, 0x0a, 0x04, 0x72, 0x65, 0x73,
	0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f,
	0x74, 0x2e, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72,
	0x2e, 0x51, 0x51, 0x5f, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x04, 0x72,
	0x65, 0x73, 0x70, 0x1a, 0x5e, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x36, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1e, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2e,
	0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52,
	0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x56, 0x61, 0x72, 0x69, 0x66,
	0x79, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x56, 0x61, 0x72, 0x69, 0x66,
	0x79, 0x49, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x22, 0x34, 0x0a, 0x12, 0x51,
	0x51, 0x5f, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x48, 0x61, 0x73, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x48, 0x61, 0x73,
	0x68, 0x22, 0xfd, 0x01, 0x0a, 0x13, 0x51, 0x51, 0x5f, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x03, 0x65, 0x72, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74,
	0x2e, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x48, 0x00, 0x52, 0x03, 0x65, 0x72, 0x72, 0x12, 0x4a, 0x0a,
	0x04, 0x72, 0x65, 0x73, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x73, 0x75,
	0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2e, 0x76, 0x65, 0x72, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x2e, 0x51, 0x51, 0x5f, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x48, 0x00, 0x52, 0x04, 0x72, 0x65, 0x73, 0x70, 0x1a, 0x5e, 0x0a, 0x08, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e,
	0x62, 0x61, 0x73, 0x69, 0x63, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x1a, 0x0a,
	0x08, 0x56, 0x61, 0x72, 0x69, 0x66, 0x79, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x56, 0x61, 0x72, 0x69, 0x66, 0x79, 0x49, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x42, 0x6f, 0x64,
	0x79, 0x2a, 0x91, 0x01, 0x0a, 0x06, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x12, 0x09, 0x0a, 0x05,
	0x45, 0x4d, 0x50, 0x54, 0x59, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x76, 0x61, 0x6c, 0x69, 0x6e, 0x67, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x6f, 0x46, 0x72,
	0x69, 0x65, 0x6e, 0x64, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x45, 0x72, 0x72, 0x56, 0x65, 0x72,
	0x69, 0x66, 0x69, 0x65, 0x64, 0x10, 0x03, 0x12, 0x11, 0x0a, 0x0d, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x79, 0x4e, 0x6f, 0x46, 0x6f, 0x75, 0x6e, 0x64, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07, 0x45, 0x78,
	0x70, 0x69, 0x72, 0x65, 0x64, 0x10, 0x05, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x6f, 0x64, 0x65, 0x57,
	0x72, 0x6f, 0x6e, 0x67, 0x10, 0x06, 0x12, 0x0e, 0x0a, 0x0a, 0x55, 0x6e, 0x56, 0x65, 0x72, 0x69,
	0x66, 0x69, 0x65, 0x64, 0x10, 0x07, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x6e, 0x64, 0x65, 0x66, 0x69,
	0x6e, 0x65, 0x64, 0x10, 0x08, 0x2a, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12,
	0x0c, 0x0a, 0x08, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x10, 0x00, 0x32, 0xc5, 0x02,
	0x0a, 0x08, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x6b, 0x0a, 0x0c, 0x51, 0x51,
	0x5f, 0x4e, 0x65, 0x77, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x12, 0x2b, 0x2e, 0x73, 0x75, 0x73,
	0x75, 0x62, 0x6f, 0x74, 0x2e, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x2e, 0x51, 0x51, 0x5f, 0x4e, 0x65, 0x77, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f,
	0x74, 0x2e, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72,
	0x2e, 0x51, 0x51, 0x5f, 0x4e, 0x65, 0x77, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x62, 0x0a, 0x09, 0x51, 0x51, 0x5f, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x12, 0x28, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x62,
	0x61, 0x73, 0x69, 0x63, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e, 0x51, 0x51,
	0x5f, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29,
	0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74, 0x2e, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2e, 0x76,
	0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e, 0x51, 0x51, 0x5f, 0x56, 0x65, 0x72, 0x69, 0x66,
	0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x68, 0x0a, 0x0b, 0x51,
	0x51, 0x5f, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x12, 0x2a, 0x2e, 0x73, 0x75, 0x73,
	0x75, 0x62, 0x6f, 0x74, 0x2e, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x2e, 0x51, 0x51, 0x5f, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x73, 0x75, 0x73, 0x75, 0x62, 0x6f, 0x74,
	0x2e, 0x62, 0x61, 0x73, 0x69, 0x63, 0x2e, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x2e,
	0x51, 0x51, 0x5f, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x11, 0x5a, 0x0f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f,
	0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescOnce sync.Once
	file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescData = file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDesc
)

func file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescGZIP() []byte {
	file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescOnce.Do(func() {
		file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescData)
	})
	return file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDescData
}

var file_pkg_protos_verifier_verifier_WITHOUT_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_pkg_protos_verifier_verifier_WITHOUT_proto_goTypes = []any{
	(Errors)(0),                          // 0: susubot.basic.verifier.Errors
	(Result)(0),                          // 1: susubot.basic.verifier.Result
	(*QQ_NewVerifyRequest)(nil),          // 2: susubot.basic.verifier.QQ_NewVerifyRequest
	(*QQ_NewVerifyResponse)(nil),         // 3: susubot.basic.verifier.QQ_NewVerifyResponse
	(*QQ_VerifyRequest)(nil),             // 4: susubot.basic.verifier.QQ_VerifyRequest
	(*QQ_VerifyResponse)(nil),            // 5: susubot.basic.verifier.QQ_VerifyResponse
	(*QQ_VerifiedRequest)(nil),           // 6: susubot.basic.verifier.QQ_VerifiedRequest
	(*QQ_VerifiedResponse)(nil),          // 7: susubot.basic.verifier.QQ_VerifiedResponse
	(*QQ_VerifyResponse_Response)(nil),   // 8: susubot.basic.verifier.QQ_VerifyResponse.Response
	(*QQ_VerifiedResponse_Response)(nil), // 9: susubot.basic.verifier.QQ_VerifiedResponse.Response
}
var file_pkg_protos_verifier_verifier_WITHOUT_proto_depIdxs = []int32{
	0,  // 0: susubot.basic.verifier.QQ_NewVerifyResponse.err:type_name -> susubot.basic.verifier.Errors
	0,  // 1: susubot.basic.verifier.QQ_VerifyResponse.err:type_name -> susubot.basic.verifier.Errors
	8,  // 2: susubot.basic.verifier.QQ_VerifyResponse.resp:type_name -> susubot.basic.verifier.QQ_VerifyResponse.Response
	0,  // 3: susubot.basic.verifier.QQ_VerifiedResponse.err:type_name -> susubot.basic.verifier.Errors
	9,  // 4: susubot.basic.verifier.QQ_VerifiedResponse.resp:type_name -> susubot.basic.verifier.QQ_VerifiedResponse.Response
	1,  // 5: susubot.basic.verifier.QQ_VerifyResponse.Response.result:type_name -> susubot.basic.verifier.Result
	1,  // 6: susubot.basic.verifier.QQ_VerifiedResponse.Response.result:type_name -> susubot.basic.verifier.Result
	2,  // 7: susubot.basic.verifier.verifier.QQ_NewVerify:input_type -> susubot.basic.verifier.QQ_NewVerifyRequest
	4,  // 8: susubot.basic.verifier.verifier.QQ_Verify:input_type -> susubot.basic.verifier.QQ_VerifyRequest
	6,  // 9: susubot.basic.verifier.verifier.QQ_Verified:input_type -> susubot.basic.verifier.QQ_VerifiedRequest
	3,  // 10: susubot.basic.verifier.verifier.QQ_NewVerify:output_type -> susubot.basic.verifier.QQ_NewVerifyResponse
	5,  // 11: susubot.basic.verifier.verifier.QQ_Verify:output_type -> susubot.basic.verifier.QQ_VerifyResponse
	7,  // 12: susubot.basic.verifier.verifier.QQ_Verified:output_type -> susubot.basic.verifier.QQ_VerifiedResponse
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_pkg_protos_verifier_verifier_WITHOUT_proto_init() }
func file_pkg_protos_verifier_verifier_WITHOUT_proto_init() {
	if File_pkg_protos_verifier_verifier_WITHOUT_proto != nil {
		return
	}
	file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[1].OneofWrappers = []any{
		(*QQ_NewVerifyResponse_Err)(nil),
		(*QQ_NewVerifyResponse_VerifyHash)(nil),
	}
	file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[3].OneofWrappers = []any{
		(*QQ_VerifyResponse_Err)(nil),
		(*QQ_VerifyResponse_Resp)(nil),
	}
	file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes[5].OneofWrappers = []any{
		(*QQ_VerifiedResponse_Err)(nil),
		(*QQ_VerifiedResponse_Resp)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_protos_verifier_verifier_WITHOUT_proto_goTypes,
		DependencyIndexes: file_pkg_protos_verifier_verifier_WITHOUT_proto_depIdxs,
		EnumInfos:         file_pkg_protos_verifier_verifier_WITHOUT_proto_enumTypes,
		MessageInfos:      file_pkg_protos_verifier_verifier_WITHOUT_proto_msgTypes,
	}.Build()
	File_pkg_protos_verifier_verifier_WITHOUT_proto = out.File
	file_pkg_protos_verifier_verifier_WITHOUT_proto_rawDesc = nil
	file_pkg_protos_verifier_verifier_WITHOUT_proto_goTypes = nil
	file_pkg_protos_verifier_verifier_WITHOUT_proto_depIdxs = nil
}

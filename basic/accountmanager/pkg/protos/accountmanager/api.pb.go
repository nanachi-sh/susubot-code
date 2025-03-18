// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.28.2
// source: pkg/protos/accountmanager/api.proto

package accountmanager

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

type Error int32

const (
	Error_ERROR_UNKNOWN                      Error = 0
	Error_ERROR_UNDEFINED                    Error = 1
	Error_ERROR_INVALID_ARGUMENT             Error = 2
	Error_ERROR_NO_VERIFYCODE_AUTH           Error = 1001
	Error_ERROR_EMAIL_EXISTED                Error = 1002 //邮箱已注册
	Error_ERROR_VERIFYCODE_ANSWER_FAIL       Error = 1003 //验证码答案错误
	Error_ERROR_EMAIL_VERIFYCODE_ANSWER_FAIL Error = 1004 //邮箱验证码答案错误
)

// Enum value maps for Error.
var (
	Error_name = map[int32]string{
		0:    "ERROR_UNKNOWN",
		1:    "ERROR_UNDEFINED",
		2:    "ERROR_INVALID_ARGUMENT",
		1001: "ERROR_NO_VERIFYCODE_AUTH",
		1002: "ERROR_EMAIL_EXISTED",
		1003: "ERROR_VERIFYCODE_ANSWER_FAIL",
		1004: "ERROR_EMAIL_VERIFYCODE_ANSWER_FAIL",
	}
	Error_value = map[string]int32{
		"ERROR_UNKNOWN":                      0,
		"ERROR_UNDEFINED":                    1,
		"ERROR_INVALID_ARGUMENT":             2,
		"ERROR_NO_VERIFYCODE_AUTH":           1001,
		"ERROR_EMAIL_EXISTED":                1002,
		"ERROR_VERIFYCODE_ANSWER_FAIL":       1003,
		"ERROR_EMAIL_VERIFYCODE_ANSWER_FAIL": 1004,
	}
)

func (x Error) Enum() *Error {
	p := new(Error)
	*p = x
	return p
}

func (x Error) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Error) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_protos_accountmanager_api_proto_enumTypes[0].Descriptor()
}

func (Error) Type() protoreflect.EnumType {
	return &file_pkg_protos_accountmanager_api_proto_enumTypes[0]
}

func (x Error) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Error.Descriptor instead.
func (Error) EnumDescriptor() ([]byte, []int) {
	return file_pkg_protos_accountmanager_api_proto_rawDescGZIP(), []int{0}
}

var File_pkg_protos_accountmanager_api_proto protoreflect.FileDescriptor

var file_pkg_protos_accountmanager_api_proto_rawDesc = []byte{
	0x0a, 0x23, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2a, 0xd0, 0x01, 0x0a, 0x05, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x12, 0x11, 0x0a, 0x0d, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x55, 0x4e, 0x4b,
	0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f,
	0x55, 0x4e, 0x44, 0x45, 0x46, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x41, 0x52, 0x47,
	0x55, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x02, 0x12, 0x1d, 0x0a, 0x18, 0x45, 0x52, 0x52, 0x4f, 0x52,
	0x5f, 0x4e, 0x4f, 0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x59, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x41,
	0x55, 0x54, 0x48, 0x10, 0xe9, 0x07, 0x12, 0x18, 0x0a, 0x13, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f,
	0x45, 0x4d, 0x41, 0x49, 0x4c, 0x5f, 0x45, 0x58, 0x49, 0x53, 0x54, 0x45, 0x44, 0x10, 0xea, 0x07,
	0x12, 0x21, 0x0a, 0x1c, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x59,
	0x43, 0x4f, 0x44, 0x45, 0x5f, 0x41, 0x4e, 0x53, 0x57, 0x45, 0x52, 0x5f, 0x46, 0x41, 0x49, 0x4c,
	0x10, 0xeb, 0x07, 0x12, 0x27, 0x0a, 0x22, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x45, 0x4d, 0x41,
	0x49, 0x4c, 0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x59, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x41, 0x4e,
	0x53, 0x57, 0x45, 0x52, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x10, 0xec, 0x07, 0x42, 0x17, 0x5a, 0x15,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_protos_accountmanager_api_proto_rawDescOnce sync.Once
	file_pkg_protos_accountmanager_api_proto_rawDescData = file_pkg_protos_accountmanager_api_proto_rawDesc
)

func file_pkg_protos_accountmanager_api_proto_rawDescGZIP() []byte {
	file_pkg_protos_accountmanager_api_proto_rawDescOnce.Do(func() {
		file_pkg_protos_accountmanager_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_protos_accountmanager_api_proto_rawDescData)
	})
	return file_pkg_protos_accountmanager_api_proto_rawDescData
}

var file_pkg_protos_accountmanager_api_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_protos_accountmanager_api_proto_goTypes = []any{
	(Error)(0), // 0: accountmanager.api.Error
}
var file_pkg_protos_accountmanager_api_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_protos_accountmanager_api_proto_init() }
func file_pkg_protos_accountmanager_api_proto_init() {
	if File_pkg_protos_accountmanager_api_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_protos_accountmanager_api_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_protos_accountmanager_api_proto_goTypes,
		DependencyIndexes: file_pkg_protos_accountmanager_api_proto_depIdxs,
		EnumInfos:         file_pkg_protos_accountmanager_api_proto_enumTypes,
	}.Build()
	File_pkg_protos_accountmanager_api_proto = out.File
	file_pkg_protos_accountmanager_api_proto_rawDesc = nil
	file_pkg_protos_accountmanager_api_proto_goTypes = nil
	file_pkg_protos_accountmanager_api_proto_depIdxs = nil
}

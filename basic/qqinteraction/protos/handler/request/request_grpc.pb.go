// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: basic/handler/protos/handler/request/request.proto

package request

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RequestHandler_SendGroupMessage_FullMethodName   = "/susubot.basic.handler.request.RequestHandler/SendGroupMessage"
	RequestHandler_SendFriendMessage_FullMethodName  = "/susubot.basic.handler.request.RequestHandler/SendFriendMessage"
	RequestHandler_MessageRecall_FullMethodName      = "/susubot.basic.handler.request.RequestHandler/MessageRecall"
	RequestHandler_GetMessage_FullMethodName         = "/susubot.basic.handler.request.RequestHandler/GetMessage"
	RequestHandler_GetGroupInfo_FullMethodName       = "/susubot.basic.handler.request.RequestHandler/GetGroupInfo"
	RequestHandler_GetGroupMemberInfo_FullMethodName = "/susubot.basic.handler.request.RequestHandler/GetGroupMemberInfo"
	RequestHandler_GetFriendList_FullMethodName      = "/susubot.basic.handler.request.RequestHandler/GetFriendList"
	RequestHandler_GetFriendInfo_FullMethodName      = "/susubot.basic.handler.request.RequestHandler/GetFriendInfo"
)

// RequestHandlerClient is the client API for RequestHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RequestHandlerClient interface {
	SendGroupMessage(ctx context.Context, in *SendGroupMessageRequest, opts ...grpc.CallOption) (*BasicResponse, error)
	SendFriendMessage(ctx context.Context, in *SendFriendMessageRequest, opts ...grpc.CallOption) (*BasicResponse, error)
	MessageRecall(ctx context.Context, in *MessageRecallRequest, opts ...grpc.CallOption) (*BasicResponse, error)
	GetMessage(ctx context.Context, in *GetMessageRequest, opts ...grpc.CallOption) (*BasicResponse, error)
	GetGroupInfo(ctx context.Context, in *GetGroupInfoRequest, opts ...grpc.CallOption) (*BasicResponse, error)
	GetGroupMemberInfo(ctx context.Context, in *GetGroupMemberInfoRequest, opts ...grpc.CallOption) (*BasicResponse, error)
	GetFriendList(ctx context.Context, in *BasicRequest, opts ...grpc.CallOption) (*BasicResponse, error)
	GetFriendInfo(ctx context.Context, in *GetFriendInfoRequest, opts ...grpc.CallOption) (*BasicResponse, error)
}

type requestHandlerClient struct {
	cc grpc.ClientConnInterface
}

func NewRequestHandlerClient(cc grpc.ClientConnInterface) RequestHandlerClient {
	return &requestHandlerClient{cc}
}

func (c *requestHandlerClient) SendGroupMessage(ctx context.Context, in *SendGroupMessageRequest, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, RequestHandler_SendGroupMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requestHandlerClient) SendFriendMessage(ctx context.Context, in *SendFriendMessageRequest, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, RequestHandler_SendFriendMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requestHandlerClient) MessageRecall(ctx context.Context, in *MessageRecallRequest, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, RequestHandler_MessageRecall_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requestHandlerClient) GetMessage(ctx context.Context, in *GetMessageRequest, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, RequestHandler_GetMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requestHandlerClient) GetGroupInfo(ctx context.Context, in *GetGroupInfoRequest, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, RequestHandler_GetGroupInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requestHandlerClient) GetGroupMemberInfo(ctx context.Context, in *GetGroupMemberInfoRequest, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, RequestHandler_GetGroupMemberInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requestHandlerClient) GetFriendList(ctx context.Context, in *BasicRequest, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, RequestHandler_GetFriendList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *requestHandlerClient) GetFriendInfo(ctx context.Context, in *GetFriendInfoRequest, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, RequestHandler_GetFriendInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RequestHandlerServer is the server API for RequestHandler service.
// All implementations must embed UnimplementedRequestHandlerServer
// for forward compatibility.
type RequestHandlerServer interface {
	SendGroupMessage(context.Context, *SendGroupMessageRequest) (*BasicResponse, error)
	SendFriendMessage(context.Context, *SendFriendMessageRequest) (*BasicResponse, error)
	MessageRecall(context.Context, *MessageRecallRequest) (*BasicResponse, error)
	GetMessage(context.Context, *GetMessageRequest) (*BasicResponse, error)
	GetGroupInfo(context.Context, *GetGroupInfoRequest) (*BasicResponse, error)
	GetGroupMemberInfo(context.Context, *GetGroupMemberInfoRequest) (*BasicResponse, error)
	GetFriendList(context.Context, *BasicRequest) (*BasicResponse, error)
	GetFriendInfo(context.Context, *GetFriendInfoRequest) (*BasicResponse, error)
	mustEmbedUnimplementedRequestHandlerServer()
}

// UnimplementedRequestHandlerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRequestHandlerServer struct{}

func (UnimplementedRequestHandlerServer) SendGroupMessage(context.Context, *SendGroupMessageRequest) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendGroupMessage not implemented")
}
func (UnimplementedRequestHandlerServer) SendFriendMessage(context.Context, *SendFriendMessageRequest) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendFriendMessage not implemented")
}
func (UnimplementedRequestHandlerServer) MessageRecall(context.Context, *MessageRecallRequest) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MessageRecall not implemented")
}
func (UnimplementedRequestHandlerServer) GetMessage(context.Context, *GetMessageRequest) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessage not implemented")
}
func (UnimplementedRequestHandlerServer) GetGroupInfo(context.Context, *GetGroupInfoRequest) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupInfo not implemented")
}
func (UnimplementedRequestHandlerServer) GetGroupMemberInfo(context.Context, *GetGroupMemberInfoRequest) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupMemberInfo not implemented")
}
func (UnimplementedRequestHandlerServer) GetFriendList(context.Context, *BasicRequest) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFriendList not implemented")
}
func (UnimplementedRequestHandlerServer) GetFriendInfo(context.Context, *GetFriendInfoRequest) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFriendInfo not implemented")
}
func (UnimplementedRequestHandlerServer) mustEmbedUnimplementedRequestHandlerServer() {}
func (UnimplementedRequestHandlerServer) testEmbeddedByValue()                        {}

// UnsafeRequestHandlerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RequestHandlerServer will
// result in compilation errors.
type UnsafeRequestHandlerServer interface {
	mustEmbedUnimplementedRequestHandlerServer()
}

func RegisterRequestHandlerServer(s grpc.ServiceRegistrar, srv RequestHandlerServer) {
	// If the following call pancis, it indicates UnimplementedRequestHandlerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RequestHandler_ServiceDesc, srv)
}

func _RequestHandler_SendGroupMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendGroupMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestHandlerServer).SendGroupMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequestHandler_SendGroupMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestHandlerServer).SendGroupMessage(ctx, req.(*SendGroupMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequestHandler_SendFriendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendFriendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestHandlerServer).SendFriendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequestHandler_SendFriendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestHandlerServer).SendFriendMessage(ctx, req.(*SendFriendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequestHandler_MessageRecall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageRecallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestHandlerServer).MessageRecall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequestHandler_MessageRecall_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestHandlerServer).MessageRecall(ctx, req.(*MessageRecallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequestHandler_GetMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestHandlerServer).GetMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequestHandler_GetMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestHandlerServer).GetMessage(ctx, req.(*GetMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequestHandler_GetGroupInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestHandlerServer).GetGroupInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequestHandler_GetGroupInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestHandlerServer).GetGroupInfo(ctx, req.(*GetGroupInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequestHandler_GetGroupMemberInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupMemberInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestHandlerServer).GetGroupMemberInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequestHandler_GetGroupMemberInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestHandlerServer).GetGroupMemberInfo(ctx, req.(*GetGroupMemberInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequestHandler_GetFriendList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BasicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestHandlerServer).GetFriendList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequestHandler_GetFriendList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestHandlerServer).GetFriendList(ctx, req.(*BasicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RequestHandler_GetFriendInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFriendInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestHandlerServer).GetFriendInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RequestHandler_GetFriendInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestHandlerServer).GetFriendInfo(ctx, req.(*GetFriendInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RequestHandler_ServiceDesc is the grpc.ServiceDesc for RequestHandler service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RequestHandler_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "susubot.basic.handler.request.RequestHandler",
	HandlerType: (*RequestHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendGroupMessage",
			Handler:    _RequestHandler_SendGroupMessage_Handler,
		},
		{
			MethodName: "SendFriendMessage",
			Handler:    _RequestHandler_SendFriendMessage_Handler,
		},
		{
			MethodName: "MessageRecall",
			Handler:    _RequestHandler_MessageRecall_Handler,
		},
		{
			MethodName: "GetMessage",
			Handler:    _RequestHandler_GetMessage_Handler,
		},
		{
			MethodName: "GetGroupInfo",
			Handler:    _RequestHandler_GetGroupInfo_Handler,
		},
		{
			MethodName: "GetGroupMemberInfo",
			Handler:    _RequestHandler_GetGroupMemberInfo_Handler,
		},
		{
			MethodName: "GetFriendList",
			Handler:    _RequestHandler_GetFriendList_Handler,
		},
		{
			MethodName: "GetFriendInfo",
			Handler:    _RequestHandler_GetFriendInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "basic/handler/protos/handler/request/request.proto",
}
// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: plugin/uno/protos/uno/uno.proto

package uno

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
	Uno_CreateRoom_FullMethodName     = "/susubot.plugin.uno.uno/CreateRoom"
	Uno_GetRooms_FullMethodName       = "/susubot.plugin.uno.uno/GetRooms"
	Uno_GetRoom_FullMethodName        = "/susubot.plugin.uno.uno/GetRoom"
	Uno_JoinRoom_FullMethodName       = "/susubot.plugin.uno.uno/JoinRoom"
	Uno_ExitRoom_FullMethodName       = "/susubot.plugin.uno.uno/ExitRoom"
	Uno_StartRoom_FullMethodName      = "/susubot.plugin.uno.uno/StartRoom"
	Uno_DrawCard_FullMethodName       = "/susubot.plugin.uno.uno/DrawCard"
	Uno_SendCardAction_FullMethodName = "/susubot.plugin.uno.uno/SendCardAction"
)

// UnoClient is the client API for Uno service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UnoClient interface {
	CreateRoom(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*CreateRoomResponse, error)
	GetRooms(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetRoomsResponse, error)
	GetRoom(ctx context.Context, in *GetRoomRequest, opts ...grpc.CallOption) (*GetRoomResponse, error)
	JoinRoom(ctx context.Context, in *JoinRoomRequest, opts ...grpc.CallOption) (*JoinRoomResponse, error)
	ExitRoom(ctx context.Context, in *ExitRoomRequest, opts ...grpc.CallOption) (*ExitRoomResponse, error)
	StartRoom(ctx context.Context, in *StartRoomRequest, opts ...grpc.CallOption) (*BasicResponse, error)
	DrawCard(ctx context.Context, in *DrawCardRequest, opts ...grpc.CallOption) (*DrawCardResponse, error)
	SendCardAction(ctx context.Context, in *SendCardActionRequest, opts ...grpc.CallOption) (*SendCardActionResponse, error)
}

type unoClient struct {
	cc grpc.ClientConnInterface
}

func NewUnoClient(cc grpc.ClientConnInterface) UnoClient {
	return &unoClient{cc}
}

func (c *unoClient) CreateRoom(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*CreateRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateRoomResponse)
	err := c.cc.Invoke(ctx, Uno_CreateRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *unoClient) GetRooms(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetRoomsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomsResponse)
	err := c.cc.Invoke(ctx, Uno_GetRooms_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *unoClient) GetRoom(ctx context.Context, in *GetRoomRequest, opts ...grpc.CallOption) (*GetRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomResponse)
	err := c.cc.Invoke(ctx, Uno_GetRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *unoClient) JoinRoom(ctx context.Context, in *JoinRoomRequest, opts ...grpc.CallOption) (*JoinRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(JoinRoomResponse)
	err := c.cc.Invoke(ctx, Uno_JoinRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *unoClient) ExitRoom(ctx context.Context, in *ExitRoomRequest, opts ...grpc.CallOption) (*ExitRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExitRoomResponse)
	err := c.cc.Invoke(ctx, Uno_ExitRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *unoClient) StartRoom(ctx context.Context, in *StartRoomRequest, opts ...grpc.CallOption) (*BasicResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BasicResponse)
	err := c.cc.Invoke(ctx, Uno_StartRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *unoClient) DrawCard(ctx context.Context, in *DrawCardRequest, opts ...grpc.CallOption) (*DrawCardResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DrawCardResponse)
	err := c.cc.Invoke(ctx, Uno_DrawCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *unoClient) SendCardAction(ctx context.Context, in *SendCardActionRequest, opts ...grpc.CallOption) (*SendCardActionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendCardActionResponse)
	err := c.cc.Invoke(ctx, Uno_SendCardAction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UnoServer is the server API for Uno service.
// All implementations must embed UnimplementedUnoServer
// for forward compatibility.
type UnoServer interface {
	CreateRoom(context.Context, *Empty) (*CreateRoomResponse, error)
	GetRooms(context.Context, *Empty) (*GetRoomsResponse, error)
	GetRoom(context.Context, *GetRoomRequest) (*GetRoomResponse, error)
	JoinRoom(context.Context, *JoinRoomRequest) (*JoinRoomResponse, error)
	ExitRoom(context.Context, *ExitRoomRequest) (*ExitRoomResponse, error)
	StartRoom(context.Context, *StartRoomRequest) (*BasicResponse, error)
	DrawCard(context.Context, *DrawCardRequest) (*DrawCardResponse, error)
	SendCardAction(context.Context, *SendCardActionRequest) (*SendCardActionResponse, error)
	mustEmbedUnimplementedUnoServer()
}

// UnimplementedUnoServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUnoServer struct{}

func (UnimplementedUnoServer) CreateRoom(context.Context, *Empty) (*CreateRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRoom not implemented")
}
func (UnimplementedUnoServer) GetRooms(context.Context, *Empty) (*GetRoomsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRooms not implemented")
}
func (UnimplementedUnoServer) GetRoom(context.Context, *GetRoomRequest) (*GetRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoom not implemented")
}
func (UnimplementedUnoServer) JoinRoom(context.Context, *JoinRoomRequest) (*JoinRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinRoom not implemented")
}
func (UnimplementedUnoServer) ExitRoom(context.Context, *ExitRoomRequest) (*ExitRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExitRoom not implemented")
}
func (UnimplementedUnoServer) StartRoom(context.Context, *StartRoomRequest) (*BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartRoom not implemented")
}
func (UnimplementedUnoServer) DrawCard(context.Context, *DrawCardRequest) (*DrawCardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DrawCard not implemented")
}
func (UnimplementedUnoServer) SendCardAction(context.Context, *SendCardActionRequest) (*SendCardActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCardAction not implemented")
}
func (UnimplementedUnoServer) mustEmbedUnimplementedUnoServer() {}
func (UnimplementedUnoServer) testEmbeddedByValue()             {}

// UnsafeUnoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UnoServer will
// result in compilation errors.
type UnsafeUnoServer interface {
	mustEmbedUnimplementedUnoServer()
}

func RegisterUnoServer(s grpc.ServiceRegistrar, srv UnoServer) {
	// If the following call pancis, it indicates UnimplementedUnoServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Uno_ServiceDesc, srv)
}

func _Uno_CreateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnoServer).CreateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Uno_CreateRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnoServer).CreateRoom(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Uno_GetRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnoServer).GetRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Uno_GetRooms_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnoServer).GetRooms(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Uno_GetRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnoServer).GetRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Uno_GetRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnoServer).GetRoom(ctx, req.(*GetRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Uno_JoinRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnoServer).JoinRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Uno_JoinRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnoServer).JoinRoom(ctx, req.(*JoinRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Uno_ExitRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExitRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnoServer).ExitRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Uno_ExitRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnoServer).ExitRoom(ctx, req.(*ExitRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Uno_StartRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnoServer).StartRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Uno_StartRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnoServer).StartRoom(ctx, req.(*StartRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Uno_DrawCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DrawCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnoServer).DrawCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Uno_DrawCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnoServer).DrawCard(ctx, req.(*DrawCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Uno_SendCardAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendCardActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnoServer).SendCardAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Uno_SendCardAction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnoServer).SendCardAction(ctx, req.(*SendCardActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Uno_ServiceDesc is the grpc.ServiceDesc for Uno service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Uno_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "susubot.plugin.uno.uno",
	HandlerType: (*UnoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRoom",
			Handler:    _Uno_CreateRoom_Handler,
		},
		{
			MethodName: "GetRooms",
			Handler:    _Uno_GetRooms_Handler,
		},
		{
			MethodName: "GetRoom",
			Handler:    _Uno_GetRoom_Handler,
		},
		{
			MethodName: "JoinRoom",
			Handler:    _Uno_JoinRoom_Handler,
		},
		{
			MethodName: "ExitRoom",
			Handler:    _Uno_ExitRoom_Handler,
		},
		{
			MethodName: "StartRoom",
			Handler:    _Uno_StartRoom_Handler,
		},
		{
			MethodName: "DrawCard",
			Handler:    _Uno_DrawCard_Handler,
		},
		{
			MethodName: "SendCardAction",
			Handler:    _Uno_SendCardAction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugin/uno/protos/uno/uno.proto",
}

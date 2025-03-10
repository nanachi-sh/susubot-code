package api

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/nanachi-sh/susubot-code/plugin/uno/define"
	uno_pb "github.com/nanachi-sh/susubot-code/plugin/uno/protos/uno"
	"github.com/nanachi-sh/susubot-code/plugin/uno/uno"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type (
	unoService struct {
		uno_pb.UnoServer
	}
)

func GetCookies(md metadata.MD) ([]*http.Cookie, error) {
	csStr := []string{}
	if s := md.Get("grpcgateway-cookie"); len(s) > 0 {
		csStr = s
	} else if s := md.Get("cookie"); len(s) > 0 {
		csStr = s
	} else {
		return nil, nil
	}
	cs := []*http.Cookie{}
	for _, v := range csStr {
		c, err := http.ParseCookie(v)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c...)
	}
	return cs, nil
}

func GRPCServe() error {
	portStr := os.Getenv("GRPC_LISTEN_PORT")
	if portStr == "" {
		return errors.New("gRPC服务监听端口未设置")
	}
	port, err := strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		return err
	}
	if port <= 0 || port > 65535 {
		return errors.New("gRPC服务监听端口范围不正确")
	}
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		return err
	}
	opts := []grpc.ServerOption{}
	if define.EnableTLS {
		cert, err := tls.LoadX509KeyPair(fmt.Sprintf("%v/tls.pem", define.CertsDir), fmt.Sprintf("%v/tls.key", define.CertsDir))
		if err != nil {
			return err
		}
		cred := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
		})
		opts = append(opts, grpc.Creds(cred))
	}
	gs := grpc.NewServer(opts...)
	uno_pb.RegisterUnoServer(gs, new(unoService))
	return gs.Serve(l)
}

func (*unoService) GetPlayer(ctx context.Context, req *uno_pb.GetPlayerRequest) (*uno_pb.GetPlayerResponse, error) {
	type d struct {
		data *uno_pb.GetPlayerResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp := uno.GetPlayer(cookies, req)
		if resp == nil {
			ret.data = &uno_pb.GetPlayerResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) CreateUser(ctx context.Context, req *uno_pb.CreateUserRequest) (*uno_pb.BasicResponse, error) {
	type d struct {
		data *uno_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp, err := uno.CreateUser(cookies, req)
		if err != nil {
			ret.err = err
			return
		}
		if resp == nil {
			ret.data = &uno_pb.BasicResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) GetUser(ctx context.Context, req *uno_pb.GetUserRequest) (*uno_pb.GetUserResponse, error) {
	type d struct {
		data *uno_pb.GetUserResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp, err := uno.GetUser(cookies, req)
		if err != nil {
			ret.err = err
			return
		}
		if resp == nil {
			ret.data = &uno_pb.GetUserResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) RoomEvent(req *uno_pb.RoomEventRequest, stream grpc.ServerStreamingServer[uno_pb.RoomEventResponse]) error {
	resp, err := uno.RoomEvent(req, stream)
	if err != nil {
		return err
	}
	if resp != nil && resp.Err != nil {
		stream.Send(&uno_pb.RoomEventResponse{Err: resp.Err})
	}
	return nil
}

// 仅允许正式玩家创建
func (*unoService) CreateRoom(ctx context.Context, _ *uno_pb.Empty) (*uno_pb.CreateRoomResponse, error) {
	type d struct {
		data *uno_pb.CreateRoomResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp, err := uno.CreateRoom(cookies)
		if err != nil {
			ret.err = err
			return
		}
		if resp == nil {
			ret.data = &uno_pb.CreateRoomResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) GetRooms(ctx context.Context, _ *uno_pb.Empty) (*uno_pb.GetRoomsResponse, error) {
	type d struct {
		data *uno_pb.GetRoomsResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		resp := uno.GetRooms()
		if resp == nil {
			ret.data = &uno_pb.GetRoomsResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) GetRoom(ctx context.Context, req *uno_pb.GetRoomRequest) (*uno_pb.GetRoomResponse, error) {
	type d struct {
		data *uno_pb.GetRoomResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		//普通或临时玩家则返回桌基本信息
		resp := uno.GetRoom(cookies, req)
		if resp == nil {
			ret.data = &uno_pb.GetRoomResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) JoinRoom(ctx context.Context, req *uno_pb.JoinRoomRequest) (*uno_pb.JoinRoomResponse, error) {
	type d struct {
		data *uno_pb.JoinRoomResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp, err := uno.JoinRoom(cookies, req)
		if err != nil {
			ret.err = err
			return
		}
		if resp == nil {
			ret.data = &uno_pb.JoinRoomResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) ExitRoom(ctx context.Context, req *uno_pb.ExitRoomRequest) (*uno_pb.ExitRoomResponse, error) {
	type d struct {
		data *uno_pb.ExitRoomResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp := uno.ExitRoom(cookies, req)
		if resp == nil {
			ret.data = &uno_pb.ExitRoomResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) StartRoom(ctx context.Context, req *uno_pb.StartRoomRequest) (*uno_pb.BasicResponse, error) {
	type d struct {
		data *uno_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp := uno.StartRoom(cookies, req)
		if resp == nil {
			ret.data = &uno_pb.BasicResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) DrawCard(ctx context.Context, req *uno_pb.DrawCardRequest) (*uno_pb.DrawCardResponse, error) {
	type d struct {
		data *uno_pb.DrawCardResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp := uno.DrawCard(cookies, req)
		if resp == nil {
			ret.data = &uno_pb.DrawCardResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) SendCard(ctx context.Context, req *uno_pb.SendCardRequest) (*uno_pb.SendCardResponse, error) {
	type d struct {
		data *uno_pb.SendCardResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp := uno.SendCard(cookies, req)
		if resp == nil {
			ret.data = &uno_pb.SendCardResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) NoSendCard(ctx context.Context, req *uno_pb.NoSendCardRequest) (*uno_pb.NoSendCardResponse, error) {
	type d struct {
		data *uno_pb.NoSendCardResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp := uno.NoSendCard(cookies, req)
		if resp == nil {
			ret.data = &uno_pb.NoSendCardResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) CallUNO(ctx context.Context, req *uno_pb.CallUNORequest) (*uno_pb.CallUNOResponse, error) {
	type d struct {
		data *uno_pb.CallUNOResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp := uno.CallUNO(cookies, req)
		if resp == nil {
			ret.data = &uno_pb.CallUNOResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) Challenge(ctx context.Context, req *uno_pb.ChallengeRequest) (*uno_pb.ChallengeResponse, error) {
	type d struct {
		data *uno_pb.ChallengeResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp := uno.Challenge(cookies, req)
		if resp == nil {
			ret.data = &uno_pb.ChallengeResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) IndicateUNO(ctx context.Context, req *uno_pb.IndicateUNORequest) (*uno_pb.IndicateUNOResponse, error) {
	type d struct {
		data *uno_pb.IndicateUNOResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp := uno.IndicateUNO(cookies, req)
		if resp == nil {
			ret.data = &uno_pb.IndicateUNOResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*unoService) TEST_SetPlayerCard(ctx context.Context, req *uno_pb.TEST_SetPlayerCardRequest) (*uno_pb.BasicResponse, error) {
	type d struct {
		data *uno_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			ret.err = errors.New("从context获取metadata失败")
		}
		cookies, err := GetCookies(md)
		if err != nil {
			ret.err = err
		}
		resp := uno.TEST_SetPlayerCard(cookies, req)
		if resp == nil {
			ret.data = &uno_pb.BasicResponse{}
		} else {
			ret.data = resp
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

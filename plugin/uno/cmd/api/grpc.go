package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	uno_pb "github.com/nanachi-sh/susubot-code/plugin/uno/protos/uno"
	"github.com/nanachi-sh/susubot-code/plugin/uno/uno"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type (
	unoService struct {
		uno_pb.UnoServer
	}
)

func GetCookies(md metadata.MD) []*http.Cookie {
	csStr := []string{}
	if s := md.Get("grpcgateway-cookie"); len(s) > 0 {
		csStr = s
	} else if s := md.Get("cookie"); len(s) > 0 {
		csStr = s
	} else {
		return nil
	}
	cs := []*http.Cookie{}
	for _, v := range csStr {
		vS := strings.Split(v, ";")
		//仅一个Cookie
		if len(vS) == 1 {
			vCS := strings.Split(v, "=")
			//异常Cookie结构
			if len(vCS) != 2 {
				return nil
			}
			key := strings.TrimSpace(vCS[0])
			value := strings.TrimSpace(vCS[1])
			cs = append(cs, &http.Cookie{
				Name:  key,
				Value: value,
			})
		} else { //多个Cookie
			for _, vSv := range vS {
				vSvCS := strings.Split(vSv, "=")
				//异常Cookie结构
				if len(vSvCS) != 2 {
					return nil
				}
				key := strings.TrimSpace(vSvCS[0])
				value := strings.TrimSpace(vSvCS[1])
				cs = append(cs, &http.Cookie{
					Name:  key,
					Value: value,
				})
			}
		}
	}
	return cs
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
	gs := grpc.NewServer()
	uno_pb.RegisterUnoServer(gs, new(unoService))
	return gs.Serve(l)
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
		cookies := GetCookies(md)
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
		cookies := GetCookies(md)
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
		cookies := GetCookies(md)
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
		cookies := GetCookies(md)
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
		cookies := GetCookies(md)
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
		cookies := GetCookies(md)
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
		cookies := GetCookies(md)
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
		cookies := GetCookies(md)
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
		cookies := GetCookies(md)
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
		cookies := GetCookies(md)
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
		cookies := GetCookies(md)
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
		cookies := GetCookies(md)
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

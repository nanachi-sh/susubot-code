package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	uno_pb "github.com/nanachi-sh/susubot-code/plugin/uno/protos/uno"
	"github.com/nanachi-sh/susubot-code/plugin/uno/uno"
	"google.golang.org/grpc"
)

type (
	unoService struct {
		uno_pb.UnoServer
	}
)

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

func (*unoService) CreateRoom(ctx context.Context, _ *uno_pb.Empty) (*uno_pb.CreateRoomResponse, error) {
	type d struct {
		data *uno_pb.CreateRoomResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		resp := uno.CreateRoom()
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
		resp := uno.GetRoom(req)
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
		resp := uno.JoinRoom(req)
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
		resp := uno.ExitRoom(req)
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
		resp := uno.StartRoom(req)
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
		resp := uno.DrawCard(req)
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

func (*unoService) SendCardAction(ctx context.Context, req *uno_pb.SendCardActionRequest) (*uno_pb.SendCardActionResponse, error) {
	type d struct {
		data *uno_pb.SendCardActionResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		resp := uno.SendCardAction(req)
		if resp == nil {
			ret.data = &uno_pb.SendCardActionResponse{}
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
		resp := uno.CallUNO(req)
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
		resp := uno.Challenge(req)
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
		resp := uno.IndicateUNO(req)
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
		resp := uno.TEST_SetPlayerCard(req)
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

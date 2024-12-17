package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	twoonone_pb "github.com/nanachi-sh/susubot-code/plugin/TwoOnOne/LLOneBot/protos/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/TwoOnOne/LLOneBot/twoonone"
	"google.golang.org/grpc"
)

type (
	twononeService struct {
		twoonone_pb.TwoOnOneServer
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
	twoonone_pb.RegisterTwoOnOneServer(gs, new(twononeService))
	return gs.Serve(l)
}

func (*twononeService) GetRooms(ctx context.Context, _ *twoonone_pb.Empty) (*twoonone_pb.GetRoomsResponse, error) {
	type d struct {
		data *twoonone_pb.GetRoomsResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		ris := twoonone.GetRooms()
		ret.data = &twoonone_pb.GetRoomsResponse{
			Rooms: ris,
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

func (*twononeService) CreateAccount(ctx context.Context, req *twoonone_pb.CreateAccountRequest) (*twoonone_pb.BasicResponse, error) {
	type d struct {
		data *twoonone_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		serr, err := twoonone.CreateAccount(req)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &twoonone_pb.BasicResponse{
			Err: serr,
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

func (*twononeService) GetAccount(ctx context.Context, req *twoonone_pb.GetAccountRequest) (*twoonone_pb.GetAccountResponse, error) {
	type d struct {
		data *twoonone_pb.GetAccountResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		ai, serr, err := twoonone.GetAccount(req)
		if err != nil {
			ret.err = err
			return
		}
		if serr != nil {
			ret.data = &twoonone_pb.GetAccountResponse{
				Err: serr,
			}
			return
		}
		ret.data = &twoonone_pb.GetAccountResponse{
			Info: &twoonone_pb.PlayerInfo{
				AccountInfo: ai,
				TableInfo:   nil,
			},
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

func (*twononeService) GetDailyCoin(ctx context.Context, req *twoonone_pb.GetDailyCoinRequest) (*twoonone_pb.BasicResponse, error) {
	type d struct {
		data *twoonone_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		ai, err := twoonone.GetDailyCoin(req)
		if err != nil {
			ret.err = err
			return
		}
		if ai == nil {
			ai = &twoonone_pb.BasicResponse{}
		}
		ret.data = ai
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

func (*twononeService) GetRoomInfo(ctx context.Context, req *twoonone_pb.GetRoomRequest) (*twoonone_pb.GetRoomResponse, error) {
	type d struct {
		data *twoonone_pb.GetRoomResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		resp := twoonone.GetRoom(req)
		ret.data = resp
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

func (*twononeService) CreateRoom(ctx context.Context, req *twoonone_pb.CreateRoomRequest) (*twoonone_pb.CreateRoomResponse, error) {
	type d struct {
		data *twoonone_pb.CreateRoomResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		resp := twoonone.CreateRoom(req)
		ret.data = resp
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

func (*twononeService) JoinRoom(ctx context.Context, req *twoonone_pb.JoinRoomRequest) (*twoonone_pb.JoinRoomResponse, error) {
	type d struct {
		data *twoonone_pb.JoinRoomResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		resp, err := twoonone.JoinRoom(req)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = resp
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

func (*twononeService) ExitRoom(ctx context.Context, req *twoonone_pb.ExitRoomRequest) (*twoonone_pb.ExitRoomResponse, error) {
	type d struct {
		data *twoonone_pb.ExitRoomResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		resp, err := twoonone.ExitRoom(req)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = resp
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

func (*twononeService) StartRoom(ctx context.Context, req *twoonone_pb.StartRoomRequest) (*twoonone_pb.StartRoomResponse, error) {
	type d struct {
		data *twoonone_pb.StartRoomResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		resp, err := twoonone.StartRoom(req)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = resp
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

func (*twononeService) RobLandownerAction(ctx context.Context, req *twoonone_pb.RobLandownerActionRequest) (*twoonone_pb.RobLandownerActionResponse, error) {
	type d struct {
		data *twoonone_pb.RobLandownerActionResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		resp, err := twoonone.RobLandownerAction(req)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = resp
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

func (*twononeService) SendCardAction(ctx context.Context, req *twoonone_pb.SendCardRequest) (*twoonone_pb.SendCardResponse, error) {
	type d struct {
		data *twoonone_pb.SendCardResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		resp, err := twoonone.SendCardAction(req)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = resp
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

package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	randomfortune_pb "github.com/nanachi-sh/susubot-code/plugin/randomfortune/protos/randomfortune"
	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/randomfortune"
	"google.golang.org/grpc"
)

type (
	randomfortuneService struct {
		randomfortune_pb.RandomFortuneServer
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
	randomfortune_pb.RegisterRandomFortuneServer(gs, new(randomfortuneService))
	return gs.Serve(l)
}

func (*randomfortuneService) GetFortune(ctx context.Context, req *randomfortune_pb.BasicRequest) (*randomfortune_pb.BasicResponse, error) {
	type d struct {
		data *randomfortune_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		resp, err := randomfortune.GetFortune(req)
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

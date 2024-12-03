package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	randomanimal_pb "github.com/nanachi-sh/susubot-code/plugin/randomAnimal/LLOneBot/protos/randomAnimal"
	randomanimal "github.com/nanachi-sh/susubot-code/plugin/randomAnimal/LLOneBot/randomAnimal"
	"google.golang.org/grpc"
)

type (
	randomanimalService struct {
		randomanimal_pb.RandomAnimalServer
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
	randomanimal_pb.RegisterRandomAnimalServer(gs, new(randomanimalService))
	return gs.Serve(l)
}

func (*randomanimalService) GetCat(ctx context.Context, _ *randomanimal_pb.Empty) (*randomanimal_pb.BasicResponse, error) {
	type d struct {
		data *randomanimal_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		cat, err := randomanimal.GetCat()
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &randomanimal_pb.BasicResponse{
			Type: cat.Type,
			Buf:  cat.Buf,
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

func (*randomanimalService) GetDog(ctx context.Context, _ *randomanimal_pb.Empty) (*randomanimal_pb.BasicResponse, error) {
	type d struct {
		data *randomanimal_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		cat, err := randomanimal.GetDog()
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &randomanimal_pb.BasicResponse{
			Type: cat.Type,
			Buf:  cat.Buf,
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

func (*randomanimalService) GetFox(ctx context.Context, _ *randomanimal_pb.Empty) (*randomanimal_pb.BasicResponse, error) {
	type d struct {
		data *randomanimal_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		cat, err := randomanimal.GetFox()
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &randomanimal_pb.BasicResponse{
			Type: cat.Type,
			Buf:  cat.Buf,
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

func (*randomanimalService) GetDuck(ctx context.Context, _ *randomanimal_pb.Empty) (*randomanimal_pb.BasicResponse, error) {
	type d struct {
		data *randomanimal_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		cat, err := randomanimal.GetDuck()
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &randomanimal_pb.BasicResponse{
			Type: cat.Type,
			Buf:  cat.Buf,
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

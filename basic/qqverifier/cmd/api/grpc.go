package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	qqverifier_pb "github.com/nanachi-sh/susubot-code/basic/qqverifier/protos/qqverifier"
	"github.com/nanachi-sh/susubot-code/basic/qqverifier/qqverifier"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type qqverifierService struct{ qqverifier_pb.QqverifierServer }

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
	qqverifier_pb.RegisterQqverifierServer(gs, new(qqverifierService))
	return gs.Serve(l)
}

func (*qqverifierService) NewVerify(ctx context.Context, req *qqverifier_pb.NewVerifyRequest) (*qqverifier_pb.NewVerifyResponse, error) {
	type d struct {
		data *qqverifier_pb.NewVerifyResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		resp, err := qqverifier.NewVerify(req)
		if err != nil {
			ret.err = err
			return
		}
		if resp == nil {
			ret.data = &qqverifier_pb.NewVerifyResponse{}
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

func (*qqverifierService) Verified(ctx context.Context, req *qqverifier_pb.VerifiedRequest) (*qqverifier_pb.VerifiedResponse, error) {
	type d struct {
		data *qqverifier_pb.VerifiedResponse
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
		for k, v := range md {
			fmt.Println(k, v)
		}
		resp, err := qqverifier.Verified(req)
		if err != nil {
			ret.err = err
			return
		}
		if resp == nil {
			ret.data = &qqverifier_pb.VerifiedResponse{}
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

func (*qqverifierService) Verify(ctx context.Context, req *qqverifier_pb.VerifyRequest) (*qqverifier_pb.VerifyResponse, error) {
	type d struct {
		data *qqverifier_pb.VerifyResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		resp, err := qqverifier.Verify(req)
		if err != nil {
			ret.err = err
			return
		}
		if resp == nil {
			ret.data = &qqverifier_pb.VerifyResponse{}
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

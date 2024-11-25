package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/nanachi-sh/susubot-code/handler/protos/handler"
	"github.com/nanachi-sh/susubot-code/handler/response"
	"google.golang.org/grpc"
)

type handlerService struct {
	handler.HandlerServer
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
	hs := new(handlerService)
	s := grpc.NewServer()
	s.RegisterService(&handler.Handler_ServiceDesc, hs)
	return s.Serve(l)
}

func (*handlerService) BotResponseUnmarshal(ctx context.Context, req *handler.BotResponseUnmarshalRequest) (*handler.BotResponseUnmarshalResponse, error) {
	type d struct {
		data *handler.BotResponseUnmarshalResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := &d{}
		defer func() { ch <- ret }()
		responseH, err := response.New(req.Buf, req.Type, req.CmdEventType)
		if err != nil {
			ret.err = err
			return
		}
		response, err := responseH.MarshalToResponse()
		if err != nil {
			ret.err = err
			return
		}
		ret.data = response
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

func (*handlerService) Check(ctx context.Context, req *handler.HealthCheckRequest) (*handler.HealthCheckResponse, error) {
	return &handler.HealthCheckResponse{
		Status: handler.HealthCheckResponse_SERVING,
	}, nil
}

func (*handlerService) Watch(req *handler.HealthCheckRequest, stream grpc.ServerStreamingServer[handler.HealthCheckResponse]) error {
	return errors.ErrUnsupported
}

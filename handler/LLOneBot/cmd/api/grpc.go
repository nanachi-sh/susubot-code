package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	request_pb "github.com/nanachi-sh/susubot-code/handler/LLOneBot/protos/handler/request"
	response_pb "github.com/nanachi-sh/susubot-code/handler/LLOneBot/protos/handler/response"
	"github.com/nanachi-sh/susubot-code/handler/LLOneBot/request"
	"github.com/nanachi-sh/susubot-code/handler/LLOneBot/response"
	"google.golang.org/grpc"
)

type (
	responseService struct {
		response_pb.ResponseHandlerServer
	}

	requestService struct {
		request_pb.RequestHandlerServer
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
	response_pb.RegisterResponseHandlerServer(gs, new(responseService))
	request_pb.RegisterRequestHandlerServer(gs, new(requestService))
	return gs.Serve(l)
}

func (*responseService) BotResponseUnmarshal(ctx context.Context, req *response_pb.BotResponseUnmarshalRequest) (*response_pb.BotResponseUnmarshalResponse, error) {
	type d struct {
		data *response_pb.BotResponseUnmarshalResponse
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

func (*requestService) BotRequestMarshal(ctx context.Context, req *request_pb.BotRequestMarshalRequest) (*request_pb.BotRequestMarshalResponse, error) {
	type d struct {
		data *request_pb.BotRequestMarshalResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		requestH, err := request.New(req)
		if err != nil {
			ret.err = err
			return
		}
		buf, err := requestH.Marshal()
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &request_pb.BotRequestMarshalResponse{
			Buf: buf,
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

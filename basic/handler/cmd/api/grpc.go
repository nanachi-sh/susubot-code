package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	request_pb "github.com/nanachi-sh/susubot-code/basic/handler/protos/handler/request"
	response_pb "github.com/nanachi-sh/susubot-code/basic/handler/protos/handler/response"
	"github.com/nanachi-sh/susubot-code/basic/handler/request"
	"github.com/nanachi-sh/susubot-code/basic/handler/response"
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

func (*responseService) Unmarshal(ctx context.Context, req *response_pb.UnmarshalRequest) (*response_pb.UnmarshalResponse, error) {
	type d struct {
		data *response_pb.UnmarshalResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := &d{}
		defer func() { ch <- ret }()
		responseH, err := response.New(req)
		if err != nil {
			ret.err = err
			return
		}
		if req.IgnoreCmdEvent && responseH.ResponseType() == response_pb.ResponseType_ResponseType_CmdEvent {
			ret.data = &response_pb.UnmarshalResponse{
				Type: response_pb.ResponseType_ResponseType_CmdEvent.Enum(),
			}
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

func (*requestService) SendGroupMessage(ctx context.Context, req *request_pb.SendGroupMessageRequest) (*request_pb.BasicResponse, error) {
	type d struct {
		data *request_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		buf, err := request.SendGroupMessage(req.GroupId, req.MessageChain, req.Echo)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &request_pb.BasicResponse{
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

func (*requestService) SendFriendMessage(ctx context.Context, req *request_pb.SendFriendMessageRequest) (*request_pb.BasicResponse, error) {
	type d struct {
		data *request_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		buf, err := request.SendFriendMessage(req.FriendId, req.MessageChain, req.Echo)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &request_pb.BasicResponse{
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

func (*requestService) MessageRecall(ctx context.Context, req *request_pb.MessageRecallRequest) (*request_pb.BasicResponse, error) {
	type d struct {
		data *request_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		buf, err := request.MessageRecall(req.MessageId, req.Echo)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &request_pb.BasicResponse{
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

func (*requestService) GetMessage(ctx context.Context, req *request_pb.GetMessageRequest) (*request_pb.BasicResponse, error) {
	type d struct {
		data *request_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		buf, err := request.GetMessage(req.MessageId, req.Echo)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &request_pb.BasicResponse{
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

func (*requestService) GetGroupInfo(ctx context.Context, req *request_pb.GetGroupInfoRequest) (*request_pb.BasicResponse, error) {
	type d struct {
		data *request_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		buf, err := request.GetGroupInfo(req.GroupId, req.Echo)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &request_pb.BasicResponse{
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

func (*requestService) GetGroupMemberInfo(ctx context.Context, req *request_pb.GetGroupMemberInfoRequest) (*request_pb.BasicResponse, error) {
	type d struct {
		data *request_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		buf, err := request.GetGroupMemberInfo(req.GroupId, req.UserId, req.Echo)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &request_pb.BasicResponse{
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

func (*requestService) GetFriendList(ctx context.Context, req *request_pb.BasicRequest) (*request_pb.BasicResponse, error) {
	type d struct {
		data *request_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		buf, err := request.GetFriendList(req.Echo)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &request_pb.BasicResponse{
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

func (*requestService) GetFriendInfo(ctx context.Context, req *request_pb.GetFriendInfoRequest) (*request_pb.BasicResponse, error) {
	type d struct {
		data *request_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		buf, err := request.GetFriendInfo(req.FriendId, req.Echo)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &request_pb.BasicResponse{
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

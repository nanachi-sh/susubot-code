// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5
// Source: request.proto

package server

import (
	"context"

	"github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/handler/request"
	"github.com/nanachi-sh/susubot-code/basic/handler/service/internal/logic/requesthandler"
	"github.com/nanachi-sh/susubot-code/basic/handler/service/internal/svc"
)

type RequestHandlerServer struct {
	svcCtx *svc.ServiceContext
	request.UnimplementedRequestHandlerServer
}

func NewRequestHandlerServer(svcCtx *svc.ServiceContext) *RequestHandlerServer {
	return &RequestHandlerServer{
		svcCtx: svcCtx,
	}
}

func (s *RequestHandlerServer) SendGroupMessage(ctx context.Context, in *request.SendGroupMessageRequest) (*request.BasicResponse, error) {
	l := requesthandlerlogic.NewSendGroupMessageLogic(ctx, s.svcCtx)
	return l.SendGroupMessage(in)
}

func (s *RequestHandlerServer) SendFriendMessage(ctx context.Context, in *request.SendFriendMessageRequest) (*request.BasicResponse, error) {
	l := requesthandlerlogic.NewSendFriendMessageLogic(ctx, s.svcCtx)
	return l.SendFriendMessage(in)
}

func (s *RequestHandlerServer) MessageRecall(ctx context.Context, in *request.MessageRecallRequest) (*request.BasicResponse, error) {
	l := requesthandlerlogic.NewMessageRecallLogic(ctx, s.svcCtx)
	return l.MessageRecall(in)
}

func (s *RequestHandlerServer) GetMessage(ctx context.Context, in *request.GetMessageRequest) (*request.BasicResponse, error) {
	l := requesthandlerlogic.NewGetMessageLogic(ctx, s.svcCtx)
	return l.GetMessage(in)
}

func (s *RequestHandlerServer) GetGroupInfo(ctx context.Context, in *request.GetGroupInfoRequest) (*request.BasicResponse, error) {
	l := requesthandlerlogic.NewGetGroupInfoLogic(ctx, s.svcCtx)
	return l.GetGroupInfo(in)
}

func (s *RequestHandlerServer) GetGroupMemberInfo(ctx context.Context, in *request.GetGroupMemberInfoRequest) (*request.BasicResponse, error) {
	l := requesthandlerlogic.NewGetGroupMemberInfoLogic(ctx, s.svcCtx)
	return l.GetGroupMemberInfo(in)
}

func (s *RequestHandlerServer) GetFriendList(ctx context.Context, in *request.BasicRequest) (*request.BasicResponse, error) {
	l := requesthandlerlogic.NewGetFriendListLogic(ctx, s.svcCtx)
	return l.GetFriendList(in)
}

func (s *RequestHandlerServer) GetFriendInfo(ctx context.Context, in *request.GetFriendInfoRequest) (*request.BasicResponse, error) {
	l := requesthandlerlogic.NewGetFriendInfoLogic(ctx, s.svcCtx)
	return l.GetFriendInfo(in)
}

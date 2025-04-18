// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5
// Source: response.proto

package server

import (
	"context"

	"github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/handler/response"
	"github.com/nanachi-sh/susubot-code/basic/handler/service/internal/logic/responsehandler"
	"github.com/nanachi-sh/susubot-code/basic/handler/service/internal/svc"
)

type ResponseHandlerServer struct {
	svcCtx *svc.ServiceContext
	response.UnimplementedResponseHandlerServer
}

func NewResponseHandlerServer(svcCtx *svc.ServiceContext) *ResponseHandlerServer {
	return &ResponseHandlerServer{
		svcCtx: svcCtx,
	}
}

func (s *ResponseHandlerServer) Unmarshal(ctx context.Context, in *response.UnmarshalRequest) (*response.UnmarshalResponse, error) {
	l := responsehandlerlogic.NewUnmarshalLogic(ctx, s.svcCtx)
	return l.Unmarshal(in)
}

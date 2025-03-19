package logic

import (
	"context"

	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReverseProxy_GETLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReverseProxy_GETLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReverseProxy_GETLogic {
	return &ReverseProxy_GETLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReverseProxy_GETLogic) ReverseProxy_GET() error {
	// todo: add your logic here and delete this line

	return nil
}

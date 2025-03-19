package logic

import (
	"context"

	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReverseProxy_POSTLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReverseProxy_POSTLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReverseProxy_POSTLogic {
	return &ReverseProxy_POSTLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReverseProxy_POSTLogic) ReverseProxy_POST() error {
	// todo: add your logic here and delete this line

	return nil
}

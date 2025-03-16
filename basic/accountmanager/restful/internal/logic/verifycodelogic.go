package logic

import (
	"context"

	inside "github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/accountmanager"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCodeLogic {
	return &VerifyCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyCodeLogic) VerifyCode() (resp any, err error) {
	// todo: add your logic here and delete this line

	return inside.NewRequest(l.Logger).VerifyCode()
}

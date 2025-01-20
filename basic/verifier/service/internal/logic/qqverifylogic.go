package logic

import (
	"context"

	v "github.com/nanachi-sh/susubot-code/basic/verifier/internal/verifier"
	"github.com/nanachi-sh/susubot-code/basic/verifier/pkg/protos/verifier"
	"github.com/nanachi-sh/susubot-code/basic/verifier/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QQVerifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQQVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QQVerifyLogic {
	return &QQVerifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QQVerifyLogic) QQ_Verify(in *verifier.QQ_VerifyRequest) (*verifier.QQ_VerifyResponse, error) {
	return v.NewRequest(l.Logger).QQ_Verify(in)
}

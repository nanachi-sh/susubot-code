package logic

import (
	"context"

	v "github.com/nanachi-sh/susubot-code/basic/verifier/internal/verifier"
	"github.com/nanachi-sh/susubot-code/basic/verifier/pkg/protos/verifier"
	"github.com/nanachi-sh/susubot-code/basic/verifier/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QQNewVerifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQQNewVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QQNewVerifyLogic {
	return &QQNewVerifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QQNewVerifyLogic) QQ_NewVerify(in *verifier.QQ_NewVerifyRequest) (*verifier.QQ_NewVerifyResponse, error) {
	return v.NewRequest(l.Logger).QQ_NewVerify(in)
}

package logic

import (
	"context"

	v "github.com/nanachi-sh/susubot-code/basic/verifier/internal/verifier"
	"github.com/nanachi-sh/susubot-code/basic/verifier/pkg/protos/verifier"
	"github.com/nanachi-sh/susubot-code/basic/verifier/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type QQVerifiedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQQVerifiedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QQVerifiedLogic {
	return &QQVerifiedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QQVerifiedLogic) QQ_Verified(in *verifier.QQ_VerifiedRequest) (*verifier.QQ_VerifiedResponse, error) {
	return v.NewRequest(l.Logger).QQ_Verified(in)
}

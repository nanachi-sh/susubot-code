package logic

import (
	"context"

	j "github.com/nanachi-sh/susubot-code/basic/jwt/internal/jwt"
	"github.com/nanachi-sh/susubot-code/basic/jwt/pkg/protos/jwt"
	"github.com/nanachi-sh/susubot-code/basic/jwt/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnoSignLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnoSignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnoSignLogic {
	return &UnoSignLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnoSignLogic) Uno_Sign(in *jwt.Uno_SignRequest) (*jwt.Uno_SignResponse, error) {
	return j.NewRequest(l.Logger).Uno_Sign(in)
}

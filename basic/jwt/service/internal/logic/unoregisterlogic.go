package logic

import (
	"context"

	j "github.com/nanachi-sh/susubot-code/basic/jwt/internal/jwt"
	"github.com/nanachi-sh/susubot-code/basic/jwt/pkg/protos/jwt"
	"github.com/nanachi-sh/susubot-code/basic/jwt/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnoRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnoRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnoRegisterLogic {
	return &UnoRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnoRegisterLogic) Uno_Register(in *jwt.Uno_RegisterRequest) (*jwt.Uno_RegisterResponse, error) {
	return j.NewRequest(l.Logger).Uno_Register(in)
}

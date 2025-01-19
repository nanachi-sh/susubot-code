package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &jwt.Uno_SignResponse{}, nil
}

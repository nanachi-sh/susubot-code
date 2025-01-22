package logic

import (
	"context"

	ra "github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/randomanimal"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/pkg/protos/randomanimal"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChikenCXKLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChikenCXKLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChikenCXKLogic {
	return &GetChikenCXKLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetChikenCXKLogic) GetChiken_CXK(in *randomanimal.BasicRequest) (*randomanimal.BasicResponse, error) {
	return ra.NewRequest(l.Logger).GetChiken_CXK(in)
}

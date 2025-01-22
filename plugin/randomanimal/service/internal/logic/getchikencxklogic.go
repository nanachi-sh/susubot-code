package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &randomanimal.BasicResponse{}, nil
}

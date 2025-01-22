package logic

import (
	"context"

	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/pkg/protos/randomanimal"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDuckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDuckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDuckLogic {
	return &GetDuckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDuckLogic) GetDuck(in *randomanimal.BasicRequest) (*randomanimal.BasicResponse, error) {
	// todo: add your logic here and delete this line

	return &randomanimal.BasicResponse{}, nil
}

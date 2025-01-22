package logic

import (
	"context"

	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/pkg/protos/randomanimal"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDogLogic {
	return &GetDogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDogLogic) GetDog(in *randomanimal.BasicRequest) (*randomanimal.BasicResponse, error) {
	// todo: add your logic here and delete this line

	return &randomanimal.BasicResponse{}, nil
}

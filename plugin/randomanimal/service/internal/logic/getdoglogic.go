package logic

import (
	"context"

	ra "github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/randomanimal"
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
	return ra.NewRequest(l.Logger).GetDog(in)

}

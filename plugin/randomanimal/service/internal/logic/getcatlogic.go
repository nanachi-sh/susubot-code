package logic

import (
	"context"

	ra "github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/randomanimal"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/pkg/protos/randomanimal"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCatLogic {
	return &GetCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCatLogic) GetCat(in *randomanimal.BasicRequest) (*randomanimal.BasicResponse, error) {
	return ra.NewRequest(l.Logger).GetCat(in)
}

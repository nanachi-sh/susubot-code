package logic

import (
	"context"

	ra "github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/randomanimal"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/pkg/protos/randomanimal"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFoxLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFoxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFoxLogic {
	return &GetFoxLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFoxLogic) GetFox(in *randomanimal.BasicRequest) (*randomanimal.BasicResponse, error) {
	return ra.NewRequest(l.Logger).GetFox(in)

}

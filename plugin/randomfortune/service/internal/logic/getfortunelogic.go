package logic

import (
	"context"

	rf "github.com/nanachi-sh/susubot-code/plugin/randomfortune/internal/randomfortune"
	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/pkg/protos/randomfortune"
	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFortuneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFortuneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFortuneLogic {
	return &GetFortuneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFortuneLogic) GetFortune(in *randomfortune.BasicRequest) (*randomfortune.BasicResponse, error) {
	return rf.NewRequest(l.Logger).GetFortune(in)
}

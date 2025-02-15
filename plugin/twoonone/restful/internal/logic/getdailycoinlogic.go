package logic

import (
	"context"

	inside "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDailyCoinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDailyCoinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDailyCoinLogic {
	return &GetDailyCoinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDailyCoinLogic) GetDailyCoin(req *types.GetDailyCoinRequest) (resp any, err error) {
	// todo: add your logic here and delete this line

	return inside.NewAPIRequest(l.Logger).GetDailyCoin(&twoonone.GetDailyCoinRequest{
		UserId: req.UserId,
	})
}

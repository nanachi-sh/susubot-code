package logic

import (
	"context"

	inside "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoRobLandownerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNoRobLandownerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoRobLandownerLogic {
	return &NoRobLandownerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoRobLandownerLogic) NoRobLandowner(req *types.NoRobLandownerRequest) (resp any, err error) {
	// todo: add your logic here and delete this line

	return inside.NewAPIRequest(l.Logger).NoRobLandowner(&twoonone.NoRobLandownerRequest{
		UserId:   req.UserId,
		RoomHash: req.RoomHash,
	})
}

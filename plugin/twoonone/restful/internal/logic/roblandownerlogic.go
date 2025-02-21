package logic

import (
	"context"

	inside "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone"
	pkg_types "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RobLandownerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRobLandownerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RobLandownerLogic {
	return &RobLandownerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RobLandownerLogic) RobLandowner(req *types.RobLandownerRequest) (resp any, err error) {
	// todo: add your logic here and delete this line

	return inside.NewAPIRequest(l.Logger).RobLandowner(&pkg_types.RobLandownerRequest{
		RoomHash: req.RoomHash,
		Extra:    pkg_types.Extra(req.Extra),
	})
}

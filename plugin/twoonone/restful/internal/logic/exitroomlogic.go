package logic

import (
	"context"

	inside "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone"
	pkg_types "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExitRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExitRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExitRoomLogic {
	return &ExitRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExitRoomLogic) ExitRoom(req *types.ExitRoomRequest) (resp any, err error) {
	// todo: add your logic here and delete this line
	return inside.NewAPIRequest(l.Logger).ExitRoom(&pkg_types.ExitRoomRequest{
		RoomHash: req.RoomHash,
		Extra:    pkg_types.Extra(req.Extra),
	})
}

package logic

import (
	"context"

	inside "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone"
	pkg_types "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartRoomLogic {
	return &StartRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartRoomLogic) StartRoom(req *types.StartRoomRequest) (resp any, err error) {
	// todo: add your logic here and delete this line

	return inside.NewAPIRequest(l.Logger).StartRoom(&pkg_types.StartRoomRequest{
		RoomHash: req.RoomHash,
		Extra:    pkg_types.Extra(req.Extra),
	})
}

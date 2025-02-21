package logic

import (
	"context"

	inside "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone"
	pkg_types "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCardLogic {
	return &SendCardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendCardLogic) SendCard(req *types.SendCardRequest) (resp any, err error) {
	// todo: add your logic here and delete this line

	return inside.NewAPIRequest(l.Logger).SendCard(&pkg_types.SendCardRequest{
		RoomHash:  req.RoomHash,
		SendCards: format(req.SendCards),
		Extra:     pkg_types.Extra(req.Extra),
	})
}

func format(cs []types.Card) []pkg_types.Card {
	ret := []pkg_types.Card{}
	for _, v := range cs {
		ret = append(ret, pkg_types.Card{
			Number: v.Number,
		})
	}
	return ret
}

package logic

import (
	"context"

	inside "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
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

	return inside.NewAPIRequest(l.Logger).SendCard(&twoonone.SendCardRequest{
		UserId:    "",
		RoomHash:  req.RoomHash,
		Sendcards: parseCard(req.SendCards),
	})
}

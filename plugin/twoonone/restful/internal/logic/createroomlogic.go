package logic

import (
	"context"

	inside "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoomLogic {
	return &CreateRoomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoomLogic) CreateRoom() (resp any, err error) {
	// todo: add your logic here and delete this line

	defer func() {
		l.Logger.Info(err == nil)
		l.Logger.Infof("%T", err)
	}()
	return inside.NewAPIRequest(l.Logger).CreateRoom(&twoonone.CreateRoomRequest{})
}

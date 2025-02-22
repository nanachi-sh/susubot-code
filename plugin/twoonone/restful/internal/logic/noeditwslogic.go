package logic

import (
	"context"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NOEDIT_wsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNOEDIT_wsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NOEDIT_wsLogic {
	return &NOEDIT_wsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NOEDIT_wsLogic) NOEDIT_ws(req *types.WebsocketHandShake) error {
	// todo: add your logic here and delete this line

	return nil
}

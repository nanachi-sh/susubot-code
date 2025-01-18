package logic

import (
	"context"

	c "github.com/nanachi-sh/susubot-code/basic/connector/internal/connector"
	"github.com/nanachi-sh/susubot-code/basic/connector/pkg/protos/connector"
	"github.com/nanachi-sh/susubot-code/basic/connector/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConnectLogic {
	return &ConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConnectLogic) Connect(in *connector.ConnectRequest) (*connector.ConnectResponse, error) {
	return c.NewRequest(l.Logger).Connect(in)
}

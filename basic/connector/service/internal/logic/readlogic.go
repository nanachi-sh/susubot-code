package logic

import (
	"context"

	c "github.com/nanachi-sh/susubot-code/basic/connector/internal/connector"
	"github.com/nanachi-sh/susubot-code/basic/connector/pkg/protos/connector"
	"github.com/nanachi-sh/susubot-code/basic/connector/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReadLogic) Read(in *connector.Empty, stream connector.Connector_ReadServer) error {
	return c.NewRequest(l.Logger).Read(in, stream)
}

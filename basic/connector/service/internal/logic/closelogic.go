package logic

import (
	"context"

	c "github.com/nanachi-sh/susubot-code/basic/connector/internal/connector"
	"github.com/nanachi-sh/susubot-code/basic/connector/pkg/protos/connector"
	"github.com/nanachi-sh/susubot-code/basic/connector/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCloseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloseLogic {
	return &CloseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CloseLogic) Close(in *connector.Empty) (*connector.BasicResponse, error) {
	return c.NewRequest(l.Logger).Close(in)
}

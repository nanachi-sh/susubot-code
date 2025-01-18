package logic

import (
	"context"

	c "github.com/nanachi-sh/susubot-code/basic/connector/internal/connector"
	"github.com/nanachi-sh/susubot-code/basic/connector/pkg/protos/connector"
	"github.com/nanachi-sh/susubot-code/basic/connector/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type WriteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWriteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WriteLogic {
	return &WriteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WriteLogic) Write(in *connector.WriteRequest) (*connector.BasicResponse, error) {
	return c.NewRequest(l.Logger).Write(in)
}

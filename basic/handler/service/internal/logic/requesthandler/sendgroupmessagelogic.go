package requesthandlerlogic

import (
	"context"

	r "github.com/nanachi-sh/susubot-code/basic/handler/internal/handler/request"
	"github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/handler/request"
	"github.com/nanachi-sh/susubot-code/basic/handler/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendGroupMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendGroupMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendGroupMessageLogic {
	return &SendGroupMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendGroupMessageLogic) SendGroupMessage(in *request.SendGroupMessageRequest) (*request.BasicResponse, error) {
	return r.NewRequest(l.Logger).SendGroupMessage(in)
}

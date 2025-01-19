package requesthandlerlogic

import (
	"context"

	r "github.com/nanachi-sh/susubot-code/basic/handler/internal/handler/request"
	"github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/handler/request"
	"github.com/nanachi-sh/susubot-code/basic/handler/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendFriendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendFriendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendFriendMessageLogic {
	return &SendFriendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendFriendMessageLogic) SendFriendMessage(in *request.SendFriendMessageRequest) (*request.BasicResponse, error) {
	return r.NewRequest(l.Logger).SendFriendMessage(in)
}

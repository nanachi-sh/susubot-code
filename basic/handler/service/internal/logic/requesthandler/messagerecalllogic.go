package requesthandlerlogic

import (
	"context"

	"github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/handler/request"
	"github.com/nanachi-sh/susubot-code/basic/handler/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageRecallLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessageRecallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageRecallLogic {
	return &MessageRecallLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MessageRecallLogic) MessageRecall(in *request.MessageRecallRequest) (*request.BasicResponse, error) {
	// todo: add your logic here and delete this line

	return &request.BasicResponse{}, nil
}

package responsehandlerlogic

import (
	"context"

	"github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/handler/response"
	"github.com/nanachi-sh/susubot-code/basic/handler/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnmarshalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnmarshalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnmarshalLogic {
	return &UnmarshalLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnmarshalLogic) Unmarshal(in *response.UnmarshalRequest) (*response.UnmarshalResponse, error) {
	// todo: add your logic here and delete this line

	return &response.UnmarshalResponse{}, nil
}

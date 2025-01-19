package requesthandlerlogic

import (
	"context"

	r "github.com/nanachi-sh/susubot-code/basic/handler/internal/handler/request"
	"github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/handler/request"
	"github.com/nanachi-sh/susubot-code/basic/handler/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendInfoLogic {
	return &GetFriendInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendInfoLogic) GetFriendInfo(in *request.GetFriendInfoRequest) (*request.BasicResponse, error) {
	return r.NewRequest(l.Logger).GetFriendInfo(in)
}

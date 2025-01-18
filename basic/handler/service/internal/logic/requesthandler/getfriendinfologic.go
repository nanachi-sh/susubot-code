package requesthandlerlogic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &request.BasicResponse{}, nil
}

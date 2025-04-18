package requesthandlerlogic

import (
	"context"

	r "github.com/nanachi-sh/susubot-code/basic/handler/internal/handler/request"
	"github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/handler/request"
	"github.com/nanachi-sh/susubot-code/basic/handler/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberInfoLogic {
	return &GetGroupMemberInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMemberInfoLogic) GetGroupMemberInfo(in *request.GetGroupMemberInfoRequest) (*request.BasicResponse, error) {
	return r.NewRequest(l.Logger).GetGroupMemberInfo(in)
}

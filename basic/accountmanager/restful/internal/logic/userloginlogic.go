package logic

import (
	"context"

	inside "github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/accountmanager"
	pkg_types "github.com/nanachi-sh/susubot-code/basic/accountmanager/pkg/types"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req *types.UserLoginRequest) (resp any, err error) {
	// todo: add your logic here and delete this line

	return inside.NewRequest(l.Logger).UserLogin(&pkg_types.UserLoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
}

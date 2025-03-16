package logic

import (
	"context"

	inside "github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/accountmanager"
	pkg_types "github.com/nanachi-sh/susubot-code/basic/accountmanager/pkg/types"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterRequest) (resp any, err error) {
	// todo: add your logic here and delete this line

	return inside.NewRequest(l.Logger).UserRegister(&pkg_types.UserRegisterRequest{
		VerifyCode: req.VerifyCode,
		Email:      req.Email,
		Password:   req.Password,
	})
}

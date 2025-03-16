package logic

import (
	"context"

	inside "github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/accountmanager"
	pkg_types "github.com/nanachi-sh/susubot-code/basic/accountmanager/pkg/types"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserVerifyCode_EmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserVerifyCode_EmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserVerifyCode_EmailLogic {
	return &UserVerifyCode_EmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserVerifyCode_EmailLogic) UserVerifyCode_Email(req *types.UserVerifyCodeEmailRequest) (resp any, err error) {
	// todo: add your logic here and delete this line

	return inside.NewRequest(l.Logger).UserVerifyCode_Email(&pkg_types.UserVerifyCodeEmailRequest{
		Email: req.Email,
	})
}

package database

import (
	twoonone_model "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/model/twoonone"
	pkg_types "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type User struct {
	twoonone_model.UserTwoonone
	UserPublic
}

type UserPublic struct {
	Id   string
	Name string
}

type DatabaseModel interface {
	CreateUser(logger logx.Logger, id string) *pkg_types.AppError
	GetUser(logger logx.Logger, id string) (User, *pkg_types.AppError)
	DeleteUser(logger logx.Logger, id string) *pkg_types.AppError
	UpdateUser(logger logx.Logger, id string, actions ...Action) *pkg_types.AppError
}

type Action interface {
	Merge(logx.Logger, *twoonone_model.UserTwoonone)
}

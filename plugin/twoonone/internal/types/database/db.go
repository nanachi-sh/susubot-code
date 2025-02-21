package database

import (
	twoonone_model "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/model/twoonone"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserPublic struct {
	Id   string
	Name string
}

type DatabaseModel interface {
	CreateUser(logger logx.Logger, id string) error
	GetUser(logger logx.Logger, id string) (twoonone_model.Twoonone, error)
	DeleteUser(logger logx.Logger, id string) error
	UpdateUser(logger logx.Logger, id string, actions ...Action) error
}

type Action interface {
	Merge(logx.Logger, *twoonone_model.Twoonone)
}

package svc

import "github.com/nanachi-sh/susubot-code/basic/fileweb/service/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}

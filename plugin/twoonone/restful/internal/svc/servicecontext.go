package svc

import (
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/config"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config             config.Config
	OIDCAuthMiddleware rest.Middleware
	HeaderMiddleware   rest.Middleware
	ExtraJWTMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		OIDCAuthMiddleware: middleware.NewOIDCAuthMiddleware().Handle,
		HeaderMiddleware:   middleware.NewHeaderMiddleware().Handle,
		ExtraJWTMiddleware: middleware.NewExtraJWTMiddleware().Handle,
	}
}

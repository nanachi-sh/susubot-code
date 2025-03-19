package svc

import (
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/config"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                   config.Config
	VerifyCodeAuthMiddleware rest.Middleware
	ReverseProxyMiddleware   rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                   c,
		VerifyCodeAuthMiddleware: middleware.NewVerifyCodeAuthMiddleware().Handle,
		ReverseProxyMiddleware:   middleware.NewReverseProxyMiddleware().Handle,
	}
}

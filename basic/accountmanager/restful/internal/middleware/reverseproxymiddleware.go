package middleware

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/middleware/reverseproxy"
)

type ReverseProxyMiddleware struct {
}

func NewReverseProxyMiddleware() *ReverseProxyMiddleware {
	return &ReverseProxyMiddleware{}
}

func (m *ReverseProxyMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reverseproxy.Handle(w, r, next)
	}
}

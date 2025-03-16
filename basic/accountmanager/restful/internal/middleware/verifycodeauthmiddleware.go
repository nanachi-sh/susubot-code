package middleware

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/middleware/verifycode"
)

type VerifyCodeAuthMiddleware struct {
}

func NewVerifyCodeAuthMiddleware() *VerifyCodeAuthMiddleware {
	return &VerifyCodeAuthMiddleware{}
}

func (m *VerifyCodeAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		verifycode.Handle(w, r, next)
	}
}

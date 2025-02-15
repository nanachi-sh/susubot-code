package middleware

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/middleware/auth"
)

type OIDCAuthMiddleware struct {
}

func NewOIDCAuthMiddleware() *OIDCAuthMiddleware {
	return &OIDCAuthMiddleware{}
}

func (m *OIDCAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth.Handle(w, r, next)
	}
}

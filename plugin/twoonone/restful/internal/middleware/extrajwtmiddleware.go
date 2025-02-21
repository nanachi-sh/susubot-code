package middleware

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/middleware/jwt"
)

type ExtraJWTMiddleware struct {
}

func NewExtraJWTMiddleware() *ExtraJWTMiddleware {
	return &ExtraJWTMiddleware{}
}

func (m *ExtraJWTMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need
		jwt.Handle(w, r, next)
	}
}

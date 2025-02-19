package middleware

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/middleware/header"
)

type ResponseHeaderMiddleware struct {
}

func NewResponseHeaderMiddleware() *ResponseHeaderMiddleware {
	return &ResponseHeaderMiddleware{}
}

func (m *ResponseHeaderMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need
		header.Handle(w, r, next)
	}
}

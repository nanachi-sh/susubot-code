package header

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/types"
)

func Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !request_verify(r) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Access-Control-Allow-Origin", "https://twoonone.unturned.fun:8080")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	next(w, r)
}

// 检查请求头
func request_verify(r *http.Request) bool {
	if r.Header.Get(types.HEADER_CUSTOM_KEY_extra_update) != "" {
		return false
	}
	return true
}

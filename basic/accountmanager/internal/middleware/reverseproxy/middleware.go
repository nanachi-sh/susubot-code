package reverseproxy

import (
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// API
	if len(r.RequestURI) >= 3 && r.RequestURI[:3] == "/v1" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	next(w, r)
}

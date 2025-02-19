package header

import "net/http"

func Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST")
	next(w, r)
}

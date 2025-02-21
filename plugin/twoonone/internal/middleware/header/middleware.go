package header

import "net/http"

func Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("Access-Control-Allow-Origin", "https://twoonone.unturned.fun:8080")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	next(w, r)
}

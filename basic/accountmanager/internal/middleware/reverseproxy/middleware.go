package reverseproxy

import (
	"fmt"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// API
	fmt.Println(r.RequestURI)
	next(w, r)
}

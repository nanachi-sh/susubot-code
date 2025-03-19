package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/zeromicro/go-zero/core/logx"
)

func Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	logger := logx.WithContext(r.Context())
	logger.Info("in")
	// API
	if len(r.RequestURI) >= 3 && r.RequestURI[:3] == "/v1" {
		w.Write([]byte("404 page not found\n"))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	logger.Info("s1")
	// u, err := url.Parse(configs.OIDC_ISSUER)
	// if err != nil {
	// 	logger.Error(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	u, err := url.Parse("https://test.unturned.fun:1080/v1/verify-code")
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Info("s2")
	reverse := httputil.NewSingleHostReverseProxy(u)
	reverse.ServeHTTP(w, r)
	logger.Info("s3")
}

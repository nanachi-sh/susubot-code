package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/configs"
	"github.com/zeromicro/go-zero/core/logx"
)

func Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	logger := logx.WithContext(r.Context())
	logger.Info("in")
	// API请求
	logger.Info(r.RequestURI, r.URL)
	if len(r.RequestURI) >= 3 && r.RequestURI[:3] == "/v1" {
		w.Write([]byte("404 page not found\n"))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	logger.Info("s1")
	u, err := url.Parse(configs.OIDC_ISSUER)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	r.URL.Scheme = u.Scheme
	r.URL.Host = u.Host
	r.Host = u.Host
	logger.Info("s2")

	reverse := httputil.NewSingleHostReverseProxy(u)
	reverse.ServeHTTP(w, r)
	logger.Info("s3")
}

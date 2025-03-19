package reverseproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/configs"
	"github.com/zeromicro/go-zero/core/logx"
)

func Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	logger := logx.WithContext(r.Context())
	// API请求
	if len(r.RequestURI) >= 3 && r.RequestURI[:3] == "/v1" {
		w.Write([]byte("404 page not found\n"))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Auth/Token
	if len(r.RequestURI) >= 4 && r.RequestURI[:4] == "/auth" {

	} else if len(r.RequestURI) >= 5 && r.RequestURI[:5] == "/token" {

	}

	u, err := url.Parse(configs.OIDC_ISSUER)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	r.URL.Scheme = u.Scheme
	r.URL.Host = u.Host
	r.Host = u.Host

	reverse := httputil.NewSingleHostReverseProxy(u)
	reverse.ModifyResponse = func(r *http.Response) error {
		r.Body.Close()
		buf := make([]byte, 512)
		fmt.Println(buf)
		l, _ := r.Body.Read(buf)
		fmt.Println(buf, l)
		buf = make([]byte, l)
		_, err := r.Body.Read(buf)
		if err != nil {
			logger.Error(err)
		}
		fmt.Println(string(buf))
		return nil
	}
	reverse.ServeHTTP(w, r)
}

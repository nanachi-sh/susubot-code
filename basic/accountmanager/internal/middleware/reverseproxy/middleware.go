package reverseproxy

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/configs"
	"github.com/zeromicro/go-zero/core/logx"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	logger := logx.WithContext(r.Context())
	// API请求
	if len(r.RequestURI) >= 3 && r.RequestURI[:3] == "/v1" {
		w.Write([]byte("404 page not found\n"))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var modifyResponse func(*bytes.Buffer) = func(*bytes.Buffer) {}

	// Auth/Token
	if len(r.RequestURI) >= 4 && r.RequestURI[:4] == "/auth" {

	} else if len(r.RequestURI) >= 5 && r.RequestURI[:5] == "/token" {
		// 将dex响应内容保存至redis，并给用户添加唯一标识符(感觉又梦回到session了)
		modifyResponse = func(responseBody *bytes.Buffer) {
			responseBody.Bytes()
		}
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
		return modifyResponseBase(r, modifyResponse)
	}
	reverse.ServeHTTP(w, r)
}

func modifyResponseBase(r *http.Response, modifyResponse func(buffer *bytes.Buffer)) error {
	// 获取响应
	buffer := new(bytes.Buffer)
	_, err := io.Copy(buffer, r.Body)
	if err != nil {
		return err
	}
	r.Body = io.NopCloser(buffer)
	modifyResponse(buffer)
	return nil
}

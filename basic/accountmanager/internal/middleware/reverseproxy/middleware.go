package reverseproxy

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/configs"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	logger := logx.WithContext(r.Context())
	// API请求
	if len(r.RequestURI) >= 3 && r.RequestURI[:3] == "/v1" {
		w.Write([]byte("404 page not found\n"))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sid := func(r *http.Request) string {
		sid, err := r.Cookie("SID")
		if err != nil {
			if err != http.ErrNoCookie {
				logger.Error(err)
			}
			return ""
		}
		return sid.Value
	}(r)
	if sid == "" {
		uuid, err := uuid.NewRandom()
		if err != nil {
			logger.Error(err)
			return
		}
		sid = uuid.String()
		http.SetCookie(w, &http.Cookie{
			Name:     "SID",
			Value:    sid,
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().AddDate(0, 0, 1),
		})
	}

	var modifyResponse func(*http.Response, *bytes.Buffer) = func(r *http.Response, b *bytes.Buffer) {}

	// Auth/Token
	if len(r.RequestURI) >= 4 && r.RequestURI[:4] == "/auth" {
		body, err := configs.Redis.Hget(sid, "cache_session")
		if err == nil {
			httpx.OkJson(w, body)
			return
		}
	} else if len(r.RequestURI) >= 5 && r.RequestURI[:5] == "/token" {
		// 将dex响应内容保存至redis，并给用户添加唯一标识符(感觉梦回到session)
		modifyResponse = func(resp *http.Response, responseBody *bytes.Buffer) {
			if resp.StatusCode != http.StatusOK {
				return
			}
			configs.Redis.Hset(sid, "cache_session", responseBody.String())
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

func modifyResponseBase(r *http.Response, modifyResponse func(r *http.Response, b *bytes.Buffer)) error {
	// 获取响应
	buffer := new(bytes.Buffer)
	_, err := io.Copy(buffer, r.Body)
	if err != nil {
		return err
	}
	r.Body = io.NopCloser(buffer)
	modifyResponse(r, buffer)
	return nil
}

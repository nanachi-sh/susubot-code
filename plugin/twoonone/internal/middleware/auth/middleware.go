package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/configs"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/types"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"golang.org/x/oauth2"
)

var (
	once sync.Once

	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	jwks     *jwk.Cache
	o2cfg    oauth2.Config

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	claims = struct {
		JWKSUri string `json:"jwks_uri"`
	}{}

	xlogger = logx.WithContext(context.Background())
)

func initialize() {
	provider, err := oidc.NewProvider(context.Background(), configs.OIDC_ISSUER)
	if err != nil {
		logger.Fatalln(err)
	}
	if err := provider.Claims(&claims); err != nil {
		logger.Fatalln(err)
	}
	o2cfg = oauth2.Config{
		ClientID:     configs.OIDC_CLIENT_ID,
		ClientSecret: configs.OIDC_CLIENT_SECRET,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  configs.OIDC_REDIRECT,
		Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "profile", "email"},
	}
	verifier = provider.Verifier(&oidc.Config{
		ClientID: configs.OIDC_CLIENT_ID,
	})
	{
		c, err := jwk.NewCache(context.Background(), httprc.NewClient())
		if err != nil {
			logger.Fatalln(err)
		}
		jwks = c
	}
	if err := jwks.Register(context.Background(), claims.JWKSUri); err != nil {
		logger.Fatalln(err)
	}
	for {
		if jwks.Ready(context.Background(), claims.JWKSUri) {
			break
		}
	}
}

func Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	xlogger.Info("in handle")
	defer xlogger.Info("out handle")
	if !configs.MIDDLEWARE_AuthHandlerStatus {
		next(w, r)
		return
	}
	once.Do(initialize)

	xlogger.Info("test")
	// 判断是否需要更新session_id
	session_id, ok := getSessionId(r)
	if !ok { //无session_id
		if !newSession(w, r) {
			return
		}
	} else if !verifySession(session_id) { //session_id过期
		if !newSession(w, r) {
			return
		}
	}
	// 与OIDC认证相关
	switch r.URL.Path {
	case "/login":
		if !loginHandle(w, r) {
			return
		}
	case "/callback":
		if !callbackHandle(w, r) {
			return
		}
	}
	// 普通请求
	if !handle(w, r) {
		return
	}
	next(w, r)
}

func getSessionId(r *http.Request) (string, bool) {
	c, err := r.Cookie(types.COOKIE_KEY_session_id)
	if err != nil {
		return "", false
	}
	return c.Value, true
}

func verifySession(session_id string) bool {
	_, err := getSessionFromRedis(session_id)
	return err == nil
}

func getSessionFromRedis(session_id string) (*types.REDIS_KEY_SESSION, error) {
	infos, err := configs.Redis.Hgetall(session_id)
	if err != nil {
		return nil, err
	}
	if len(infos) == 0 {
		return nil, redis.Nil
	}
	var (
		refresh_token *string
		expires_at    *int64
		state         *string
	)
	if str, ok := infos[types.REDIS_KEY_SESSION_refresh_token]; ok {
		refresh_token = &str
	}
	if str, ok := infos[types.REDIS_KEY_SESSION_expires_at]; ok {
		i64, err := strconv.ParseInt(str, 10, 0)
		if err != nil {
			return nil, err
		}
		expires_at = &i64
	}
	if str, ok := infos[types.REDIS_KEY_SESSION_State]; ok {
		state = &str
	}
	return &types.REDIS_KEY_SESSION{
		RefreshToken: refresh_token,
		ExpiresAt:    expires_at,
		State:        state,
	}, nil
}

func newSession(w http.ResponseWriter, r *http.Request) bool {
	session_id := generateSessionId()
	if err := configs.Redis.Hset(session_id, types.REDIS_KEY_SESSION_valid, ""); err != nil {
		xlogger.Error(err)
		return false
	}
	http.SetCookie(w, formatCookie(
		types.COOKIE_KEY_session_id,
		session_id,
		cutDomain(r.Host),
		0,
	))
	return true
}

func cutDomain(domain string) string {
	domain = regexp.
		MustCompile(`^(?:[a-z0-9-]+\.)+([a-z0-9-]+\.[a-z]{2,})$`).
		ReplaceAllString(domain, `$1`)
	if domain[0] != []byte(".")[0] {
		domain = "." + domain
	}
	return domain
}

func generateSessionId() string {
	return utils.RandomString(32, utils.Dict_Mixed)
}

func verifyToken(token_raw string) bool {
	token, err := jwt.Parse(token_raw, func(t *jwt.Token) (interface{}, error) {
		set, err := jwks.Lookup(context.Background(), claims.JWKSUri)
		if err != nil {
			return nil, err
		}
		kid, ok := t.Header["kid"].(string)
		if !ok {
			return nil, errors.New("token kid empty or !string")
		}
		k, ok := set.LookupKeyID(kid)
		if !ok {
			return nil, errors.New("no exist token pub key")
		}
		var pk any
		if err := jwk.Export(k, &pk); err != nil {
			return nil, err
		}
		return pk, nil
	})
	if err != nil {
		xlogger.Error(err)
		return false
	}
	return token.Valid
}

func handle(w http.ResponseWriter, r *http.Request) bool {
	access_token, err := r.Cookie(types.COOKIE_KEY_access_token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	if !verifyToken(access_token.Value) {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	return true
}

func loginHandle(w http.ResponseWriter, r *http.Request) bool {
	session_id, ok := getSessionId(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	state := utils.RandomString(16, utils.Dict_Mixed)
	if err := configs.Redis.Hset(session_id, types.REDIS_KEY_SESSION_State, state); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		xlogger.Error(err)
		return false
	}
	http.Redirect(w, r, o2cfg.AuthCodeURL(state), http.StatusFound)
	return true
}

func callbackHandle(w http.ResponseWriter, r *http.Request) bool {
	auth_code := r.URL.Query().Get("code")
	if auth_code == "" {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	state := r.URL.Query().Get("state")
	if state == "" {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	session_id, ok := getSessionId(r)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	session, err := getSessionFromRedis(session_id)
	if err != nil {
		xlogger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	// 会话state空或与用户state不一致
	if session.State == nil || *session.State != state {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	// 将authorization code发送至sso服务获取token
	token, err := o2cfg.Exchange(context.Background(), auth_code)
	if err != nil {
		xlogger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	// 获取id_token
	id_token_raw, ok := token.Extra("id_token").(string)
	if !ok {
		xlogger.Error("获取id_token失败")
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	// redis操作
	if _, err := configs.Redis.Hdel(session_id, types.REDIS_KEY_SESSION_State); err != nil {
		xlogger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	if err := configs.Redis.Hset(session_id, types.REDIS_KEY_SESSION_refresh_token, token.RefreshToken); err != nil {
		xlogger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	if err := configs.Redis.Hset(session_id, types.REDIS_KEY_SESSION_expires_at, strconv.FormatInt(token.Expiry.Unix(), 10)); err != nil {
		xlogger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	// 将cookie写入用户客户端
	http.SetCookie(w, formatCookie(
		types.COOKIE_KEY_access_token,
		token.AccessToken,
		cutDomain(r.Host),
		int(token.ExpiresIn),
	))
	http.SetCookie(w, formatCookie(
		types.COOKIE_KEY_id_token,
		id_token_raw,
		cutDomain(r.Host),
		int(token.ExpiresIn),
	))
	return true
}

func formatCookie(key, value, domain string, expires_in int) *http.Cookie {
	return &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		Domain:   domain,
		MaxAge:   expires_in,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
}

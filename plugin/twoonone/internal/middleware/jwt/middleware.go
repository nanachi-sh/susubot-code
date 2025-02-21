package jwt

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/configs"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/middleware/sql"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	db      = sql.NewHandler()
	xlogger = logx.WithContext(context.Background())
)

func update(w http.ResponseWriter, r *http.Request) bool {
	// 从数据库拉取数据，并重写请求头，后续应加上限速器，避免被穿透攻击(Cache MISS)
	c, _ := r.Cookie(types.COOKIE_KEY_access_token)
	m := jwt.MapClaims{}
	jwt.ParseWithClaims(c.Value, m, nil)
	user_id, _ := m["federated_claims"].(map[string]any)["user_id"].(string)
	i, err := db.GetUser(xlogger, user_id)
	if err != nil {
		xlogger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	// 有效时间跟随访问token
	exp, _ := m.GetExpirationTime()
	jwt_raw, err := sign(types.JWT_EXTRA{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "none",
			Subject:   "none",
			Audience:  jwt.ClaimStrings{},
			ExpiresAt: &jwt.NumericDate{Time: exp.Time},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ID:        "none",
		},
		WinCount:         int(i.Wincount),
		LoseCount:        int(i.Losecount),
		Coin:             i.Coin,
		LastGetdaliyTime: i.LastGetdaliyTime,
	})
	if err != nil {
		xlogger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	r.Header.Add("authorization", "Bearer "+jwt_raw)
	r.Header.Add(types.HEADER_CUSTOM_KEY_extra_update, "true")
	return true
}

func Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// 检查是否存在extra
	if extra_raw := r.Header.Get("authorization"); extra_raw != "" { //存在
		// 检查是否有效
		err := verify(extra_raw)
		if err != nil {
			switch err {
			case jwt.ErrTokenExpired: //extra过期,更新
				if !update(w, r) {
					return
				}
			default: //其他原因
				xlogger.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	} else { //不存在
		if !update(w, r) {
			return
		}
	}
	next(w, r)
}

func sign(extra types.JWT_EXTRA) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, extra)
	return token.SignedString([]byte(configs.JWT_SignKey))
}

func verify(raw string) error {
	token, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
		return []byte(configs.JWT_SignKey), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return jwt.ErrTokenNotValidYet
	}
	return nil
}

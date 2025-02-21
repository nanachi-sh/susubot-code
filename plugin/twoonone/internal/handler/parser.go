package handler

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/types"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
	pkg_types "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/zeromicro/go-zero/core/mapping"
)

const (
	customKey = "custom"
)

var (
	customUnmarshaler = mapping.NewUnmarshaler(
		customKey,
		mapping.WithOpaqueKeys(),
	)
)

func ParseCustom(r *http.Request, v any) error {
	var (
		user_id string
		email   string
		name    string

		wincount  int
		losecount int
		coin      float64

		extra_update bool
	)
	// 检查是否有access_token
	{
		c, err := r.Cookie(types.COOKIE_KEY_access_token)
		if err == nil && c.Value != "" {
			token_raw := c.Value
			m := jwt.MapClaims{}
			jwt.ParseWithClaims(token_raw, m, nil)
			ok := false
			email, ok = m["email"].(string)
			if !ok {
				return pkg_types.NewError(twoonone.Error_ERROR_UNDEFINED, "从访问Token获取email失败")
			}
			name, ok = m["name"].(string)
			if !ok {
				return pkg_types.NewError(twoonone.Error_ERROR_UNDEFINED, "从访问Token获取name失败")
			}
			if mi, ok := m["federated_claims"].(map[string]any); !ok {
				return pkg_types.NewError(twoonone.Error_ERROR_UNDEFINED, "从访问Token获取federated_claims失败")
			} else {
				user_id, ok = mi["user_id"].(string)
				if !ok {
					return pkg_types.NewError(twoonone.Error_ERROR_UNDEFINED, "从federated_claims获取user_id失败")
				}
			}
		}
	}
	// 检查是否有extra
	if extra_raw := r.Header.Get("authorization"); extra_raw != "" {
		if len(extra_raw) < 10 {
			return pkg_types.NewError(twoonone.Error_ERROR_UNDEFINED, "authorization format error")
		}
		if extra_raw[:7] != "Bearer " {
			return pkg_types.NewError(twoonone.Error_ERROR_UNDEFINED, "authorization not Bearer Token")
		}
		extra_raw = extra_raw[7:]
		v := &types.JWT_EXTRA{}
		jwt.ParseWithClaims(extra_raw, v, nil)
		wincount = v.WinCount
		losecount = v.LoseCount
		coin = v.Coin
	}
	// 检查extra是否需更新
	if need := r.Header.Get(types.HEADER_CUSTOM_KEY_extra_update); need != "" {
		extra_update = true
	}
	m := map[string]any{
		types.PARSE_CUSTOM_INTO: map[string]any{
			types.PARSE_CUSTOM_KEY_email:        email,
			types.PARSE_CUSTOM_KEY_name:         name,
			types.PARSE_CUSTOM_KEY_wincount:     wincount,
			types.PARSE_CUSTOM_KEY_losecount:    losecount,
			types.PARSE_CUSTOM_KEY_coin:         coin,
			types.PARSE_CUSTOM_KEY_user_id:      user_id,
			types.PARSE_CUSTOM_KEY_extra_update: extra_update,
		},
	}
	fmt.Println(m)
	if err := customUnmarshaler.Unmarshal(
		m,
		v,
	); err != nil {
		return err
	}
	return nil
}

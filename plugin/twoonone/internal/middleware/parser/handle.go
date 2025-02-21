package parser

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
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
		mapping.WithStringValues(),
		mapping.WithOpaqueKeys(),
	)
)

func ParseCustom(r *http.Request, v any) error {
	// 检查是否有access_token
	{
		c, err := r.Cookie(types.COOKIE_KEY_access_token)
		if err == nil && c.Value != "" {
			token_raw := c.Value
			m := jwt.MapClaims{}
			jwt.ParseWithClaims(token_raw, m, nil)
			email, ok := m["email"].(string)
			if !ok {
				return pkg_types.NewError(twoonone.Error_ERROR_UNDEFINED, "从访问Token获取email失败")
			}
			name, ok := m["name"].(string)
			if !ok {
				return pkg_types.NewError(twoonone.Error_ERROR_UNDEFINED, "从访问Token获取name失败")
			}
			if err := customUnmarshaler.Unmarshal(
				map[string]any{
					types.PARSE_CUSTOM_KEY_email: email,
					types.PARSE_CUSTOM_KEY_name:  name,
				},
				&v,
			); err != nil {
				return err
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
		if err := customUnmarshaler.Unmarshal(
			map[string]any{
				types.PARSE_CUSTOM_KEY_wincount:  v.WinCount,
				types.PARSE_CUSTOM_KEY_losecount: v.LoseCount,
				types.PARSE_CUSTOM_KEY_coin:      v.Coin,
			},
			&v,
		); err != nil {
			return err
		}
	}
	return nil
}

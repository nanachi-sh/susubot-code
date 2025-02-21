package types

import "github.com/golang-jwt/jwt/v5"

const (
	COOKIE_KEY_access_token = "access_token"
	COOKIE_KEY_id_token     = "id_token"
)

const (
	PARSE_CUSTOM_INTO = "extra"

	PARSE_CUSTOM_KEY_wincount           = "wincount"
	PARSE_CUSTOM_KEY_losecount          = "losecount"
	PARSE_CUSTOM_KEY_coin               = "coin"
	PARSE_CUSTOM_KEY_name               = "name"
	PARSE_CUSTOM_KEY_email              = "email"
	PARSE_CUSTOM_KEY_user_id            = "user_id"
	PARSE_CUSTOM_KEY_new_extra          = "new_extra"
	PARSE_CUSTOM_KEY_last_getdaliy_time = "last_getdaliy_time"
)

const (
	HEADER_CUSTOM_KEY_extra_update = "extra_update"
)

const (
	EXTRA_KEY_extra = "extra"
)

type JWT_EXTRA struct {
	jwt.RegisteredClaims
	WinCount         int     `json:"wincount"`
	LoseCount        int     `json:"losecount"`
	Coin             float64 `json:"coin"`
	LastGetdaliyTime int64   `json:"last_getdaliy_time"`
}

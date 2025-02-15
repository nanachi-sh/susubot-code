package types

const (
	COOKIE_KEY_session_id   = "session_id"
	COOKIE_KEY_access_token = "access_token"
	COOKIE_KEY_id_token     = "id_token"
)

const (
	REDIS_KEY_SESSION_refresh_token = "refresh_token" //
	REDIS_KEY_SESSION_expires_at    = "expires_at"    //timestamp
	REDIS_KEY_SESSION_State         = "state"         //oidc
	REDIS_KEY_SESSION_valid         = "valid"
)

const (
	LDAP_DN = "ou=user,dc=unturned,dc=fun"
)

type REDIS_KEY_SESSION struct {
	RefreshToken *string
	ExpiresAt    *int64
	State        *string
}

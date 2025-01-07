package uno

import (
	"database/sql"
	"net/http"
	"regexp"

	"github.com/nanachi-sh/susubot-code/plugin/uno/db"
	"github.com/nanachi-sh/susubot-code/plugin/uno/define"
)

const (
	key_accountHash = "account_hash"
	key_playerHash  = "player_hash"
)

// 检查是否为特权用户
func CheckPrivilegeUser(cs []*http.Cookie) bool {
	for _, v := range cs {
		if v.Name == key_accountHash {
			return v.Value == define.PrivilegeUserHash
		}
	}
	return false
}

// 检查是否为已定义来源的用户
func CheckNormalUserFromSource(cs []*http.Cookie) (bool, error) {
	for _, v := range cs {
		if v.Name == key_accountHash {
			if _, err := db.FindUser("", v.Value); err != nil {
				if err == sql.ErrNoRows {
					return false, nil
				}
				return false, err
			}
			return true, nil
		}
	}
	return false, nil
}

// 检查是否为临时玩家
func CheckTempUser(id string) bool {
	ok, err := regexp.MatchString(`^web[0-9a-zA-Z]{10,10}$`, id)
	if err != nil {
		panic("")
	}
	return ok
}

// 获取玩家哈希
func GetPlayerHash(cs []*http.Cookie) (string, bool) {
	for _, v := range cs {
		if v.Name == key_playerHash {
			return v.Value, true
		}
	}
	return "", false
}

package uno

import (
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
func CheckPrivilegeUser(uhash string) bool {
	return define.PrivilegeUserHash == uhash
}

// 检查是否为已定义来源的用户
func CheckNormalUserFromSource(uhash string) (bool, error) {
	if _, err := db.FindUser("", uhash); err != nil {
		return false, err
	}
	return true, nil
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

// 获取用户哈希
func GetUserHash(cs []*http.Cookie) (string, bool) {
	for _, v := range cs {
		if v.Name == key_accountHash {
			return v.Value, true
		}
	}
	return "", false
}

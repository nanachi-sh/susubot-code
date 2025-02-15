// 一般用于其他sql中间件调用，如MySQL，不应让sql中间件以外调用
package ldap

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/configs"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/types/database"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/zeromicro/go-zero/core/logx"
)

func FindUser(logger logx.Logger, id string) (database.UserPublic, *types.AppError) {
	result, err := configs.LDAP.Search(ldap.NewSearchRequest(
		fmt.Sprintf("uid=%s,ou=user,%s", id, configs.LDAP_DN),
		ldap.ScopeBaseObject,
		ldap.DerefFindingBaseObj,
		1, 5, false, "(objectClass=inetOrgPerson)",
		[]string{},
		nil,
	))
	if err != nil {
		if lerr, ok := err.(*ldap.Error); ok {
			switch lerr.ResultCode {
			case ldap.LDAPResultNoSuchObject:
				return database.UserPublic{}, types.NewError(twoonone.Error_ERROR_USER_NO_EXIST, "LDAP no exist the user")
			default:
				logger.Errorf("未处理错误: %s", lerr.Error())
			}
		} else {
			logger.Error(err)
		}
		return database.UserPublic{}, types.NewError(twoonone.Error_ERROR_UNDEFINED, "")
	}
	var (
		name string
	)
	for _, v := range result.Entries {
		if str := v.GetAttributeValue("displayName"); str != "" {
			name = str
		}
	}
	ok := true
	switch {
	case name == "":
		ok = false
	}
	if !ok {
		return database.UserPublic{}, types.NewError(twoonone.Error_ERROR_USER_INCOMPLETE, "")
	}
	return database.UserPublic{
		Id:   id,
		Name: name,
	}, nil
}

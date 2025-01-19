package db

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/nanachi-sh/susubot-code/basic/jwt/internal/configs"
	unomodel "github.com/nanachi-sh/susubot-code/basic/jwt/internal/model/uno"
	jwt_pb "github.com/nanachi-sh/susubot-code/basic/jwt/pkg/protos/jwt"
	"github.com/nanachi-sh/susubot-code/basic/jwt/pkg/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var requestTimeout time.Duration = time.Second * 30

func Uno_VerifyUser(logger logx.Logger, id string, password string) (bool, *jwt_pb.Errors) {
	u, serr := Uno_GetUser(logger, id)
	if serr != nil {
		return false, serr
	}
	return uno_VerifyUser(u, password), nil
}

func Uno_CreateUser(logger logx.Logger, userid, username, password string) *jwt_pb.Errors {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	salt := generateSalt()
	passwordEncrypt := encryptPassword(password, salt)
	if _, err := configs.Model_Uno.Insert(ctx, &unomodel.Users{
		Id:        userid,
		Name:      username,
		Password:  passwordEncrypt,
		Salt:      salt,
		WinCount:  0,
		LoseCount: 0,
	}); err != nil {
		switch err {
		default:
			logger.Errorf("未处理错误: %s", err.Error())
		case sqlx.ErrNotSettable:
			return jwt_pb.Errors_UserExist.Enum()
		}
		return jwt_pb.Errors_Undefined.Enum()
	}
	return nil
}

func Uno_GetUser(logger logx.Logger, userid string) (unomodel.Users, *jwt_pb.Errors) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	u, err := configs.Model_Uno.FindOne(ctx, userid)
	if err != nil {
		switch err {
		default:
			logger.Errorf("未处理错误: %s", err.Error())
		case sqlx.ErrNotFound:
			return unomodel.Users{}, jwt_pb.Errors_UserNoExist.Enum()
		}
		return unomodel.Users{}, jwt_pb.Errors_Undefined.Enum()
	}
	return *u, nil
}

func generateSalt() string {
	return utils.RandomString(6, utils.Dict)
}

func encryptPassword(password, salt string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(password+salt)))
}

func uno_VerifyUser(ui unomodel.Users, password string) bool {
	return encryptPassword(password, ui.Salt) == ui.Password
}

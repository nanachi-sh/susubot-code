package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/internal/configs"
	randomfortunemodel "github.com/nanachi-sh/susubot-code/plugin/randomfortune/internal/model/randomfortune"
	randomfortune_pb "github.com/nanachi-sh/susubot-code/plugin/randomfortune/pkg/protos/randomfortune"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func CheckPlayerTime(logger logx.Logger, id string) *randomfortune_pb.Errors {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	u, serr := getUser(logger, ctx, id)
	if serr != nil {
		if *serr == randomfortune_pb.Errors_UserNoExist {
			if serr := createUser(logger, ctx, randomfortunemodel.Users{
				Id:                 id,
				LastGetFortuneTime: time.Now().Unix(),
			}); serr != nil {
				return serr
			}
			return nil
		}
		return serr
	}
	now := time.Now()
	last := time.Unix(u.LastGetFortuneTime, 0)
	if last.Year() == now.Year() && last.Month() == now.Month() {
		fmt.Println(now.Day() > last.Day())
		if now.Day() > last.Day() {
			return nil
		} else {
			return randomfortune_pb.Errors_AlreadyGetFortune.Enum()
		}
	} else {
		return nil
	}
}

func UpdatePlayerTime(logger logx.Logger, id string) *randomfortune_pb.Errors {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if serr := updateUser(logger, ctx, randomfortunemodel.Users{
		Id:                 id,
		LastGetFortuneTime: time.Now().Unix(),
	}); serr != nil {
		return serr
	}
	return nil
}

func getUser(logger logx.Logger, ctx context.Context, id string) (randomfortunemodel.Users, *randomfortune_pb.Errors) {
	u, err := configs.Model_Randomfortune.FindOne(ctx, id)
	if err != nil {
		if myerr, ok := err.(*mysql.MySQLError); ok {
			switch myerr.Number {
			default:
				logger.Errorf("未处理错误: %v", myerr)
			}
		} else {
			switch err {
			default:
				logger.Errorf("未处理错误: %v", err)
			case sqlx.ErrNotFound:
				return randomfortunemodel.Users{}, randomfortune_pb.Errors_UserNoExist.Enum()
			}
		}
		return randomfortunemodel.Users{}, randomfortune_pb.Errors_Undefined.Enum()
	}
	return *u, nil
}

func updateUser(logger logx.Logger, ctx context.Context, u randomfortunemodel.Users) *randomfortune_pb.Errors {
	if err := configs.Model_Randomfortune.Update(ctx, &u); err != nil {
		if myerr, ok := err.(*mysql.MySQLError); ok {
			switch myerr.Number {
			default:
				logger.Errorf("未处理错误: %v", myerr)
			}
		} else {
			switch err {
			default:
				logger.Errorf("未处理错误: %v", err)
			case sqlx.ErrNotFound:
				return randomfortune_pb.Errors_UserNoExist.Enum()
			}
		}
		return randomfortune_pb.Errors_Undefined.Enum()
	}
	return nil
}

func createUser(logger logx.Logger, ctx context.Context, u randomfortunemodel.Users) *randomfortune_pb.Errors {
	if _, err := configs.Model_Randomfortune.Insert(ctx, &u); err != nil {
		if myerr, ok := err.(*mysql.MySQLError); ok {
			switch myerr.Number {
			default:
				logger.Errorf("未处理错误: %v", myerr)
			}
		} else {
			switch err {
			default:
				logger.Errorf("未处理错误: %v", err)
			}
		}
		return randomfortune_pb.Errors_Undefined.Enum()
	}
	return nil
}

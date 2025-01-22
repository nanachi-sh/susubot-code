package db

import (
	"context"
	"time"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/configs"
	randomanimalmodel "github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/model/randomanimal"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func FindCache(logger logx.Logger, id string) (randomanimalmodel.Caches, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	cache, err := configs.Model_Randomanimal.FindOne(ctx, id)
	if err != nil {
		if myerr, ok := err.(*mysql.MySQLError); ok {
			switch myerr.Number {
			default:
				logger.Errorf("未处理错误：%s", myerr.Error())
			}
		} else {
			switch err {
			default:
				logger.Errorf("未处理错误：%s", err.Error())
			case sqlx.ErrNotFound:
			}
		}
		return randomanimalmodel.Caches{}, false
	}
	return *cache, true
}

func AddCache(logger logx.Logger, id, hash string, atype randomanimalmodel.AssetType) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if _, err := configs.Model_Randomanimal.Insert(ctx, &randomanimalmodel.Caches{
		AssetType: string(atype),
		AssetId:   id,
		AssetHash: hash,
	}); err != nil {
		if myerr, ok := err.(*mysql.MySQLError); ok {
			switch myerr.Number {
			default:
				logger.Errorf("未处理错误：%s", myerr.Error())
			case mysqlerr.ER_DUP_ENTRY:
				logger.Errorf("重复缓存, Hash: %s", hash)
			}
		} else {
			switch err {
			default:
				logger.Errorf("未处理错误：%s", err.Error())
			case sqlx.ErrNotFound:
			}
		}
	}
}

func DeleteCache(logger logx.Logger, id string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if err := configs.Model_Randomanimal.Delete(ctx, id); err != nil {
		if myerr, ok := err.(*mysql.MySQLError); ok {
			switch myerr.Number {
			default:
				logger.Errorf("未处理错误：%s", myerr.Error())
			}
		} else {
			switch err {
			default:
				logger.Errorf("未处理错误：%s", err.Error())
			case sqlx.ErrNotFound:
				logger.Errorf("不存在缓存，Id: %s", id)
			}
		}
	}
}

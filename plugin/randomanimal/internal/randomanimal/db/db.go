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

func FindCache(logger logx.Logger, hash string) (randomanimalmodel.Caches, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	cache, err := configs.Model_Randomanimal.FindOne(ctx, hash)
	if err != nil {
		if myerr, ok := err.(*mysql.MySQLError); ok {
			if myerr.Is(sqlx.ErrNotFound) {
				return randomanimalmodel.Caches{}, false
			}
			switch myerr.Number {
			default:
				logger.Errorf("未处理错误：%s", myerr.Error())
			}
		} else {
			logger.Error(err)
		}
		return randomanimalmodel.Caches{}, false
	}
	return *cache, true
}

func AddCache(logger logx.Logger, hash string, atype randomanimalmodel.AssetType) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if _, err := configs.Model_Randomanimal.Insert(ctx, &randomanimalmodel.Caches{
		AssetType: string(atype),
		AssetHash: hash,
	}); err != nil {
		if myerr, ok := err.(*mysql.MySQLError); ok {
			if myerr.Is(sqlx.ErrNotFound) {
				return
			}
			switch myerr.Number {
			default:
				logger.Errorf("未处理错误：%s", myerr.Error())
			case mysqlerr.ER_DUP_ENTRY:
				logger.Errorf("重复缓存, Hash: %s", hash)
			}
		} else {
			logger.Error(err)
		}
	}
}

func DeleteCache(logger logx.Logger, hash string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if err := configs.Model_Randomanimal.Delete(ctx, hash); err != nil {
		if myerr, ok := err.(*mysql.MySQLError); ok {
			if myerr.Is(sqlx.ErrNotFound) {
				logger.Errorf("不存在缓存，Hash: %s", hash)
				return
			}
			switch myerr.Number {
			default:
				logger.Errorf("未处理错误：%s", myerr.Error())
			}
		} else {
			logger.Error(err)
		}
	}
}

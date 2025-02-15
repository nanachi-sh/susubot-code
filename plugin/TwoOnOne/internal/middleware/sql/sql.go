package sql

import (
	"time"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/middleware/sql/mysql"
	database_types "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/types/database"
)

func NewHandler() database_types.DatabaseModel {
	return mysql.NewHandler()
}

func IncCoin(n float64) database_types.Action {
	return mysql.IncCoin(n)
}

func DecCoin(n float64) database_types.Action {
	return mysql.DecCoin(n)
}

func IncWinCount(n ...uint) database_types.Action {
	return mysql.IncWinCount(n...)
}

func DecWinCount(n ...uint) database_types.Action {
	return mysql.DecWinCount(n...)
}

func IncLoseCount(n ...uint) database_types.Action {
	return mysql.IncLoseCount(n...)
}

func DecLoseCount(n ...uint) database_types.Action {
	return mysql.DecLoseCount(n...)
}

func UpdateGetDailyTime(n ...time.Time) database_types.Action {
	return mysql.UpdateGetDailyTime(n...)
}

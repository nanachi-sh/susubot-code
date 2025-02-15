package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nanachi-sh/susubot-code/plugin/TwoOnOne/log"
	twoonone_pb "github.com/nanachi-sh/susubot-code/plugin/TwoOnOne/protos/twoonone"
)

var (
	database *sql.DB
	logger   = log.Get()

	ignore any
)

func initDB() error {
	db, err := sql.Open("sqlite3", dbPosition)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	// 创建表
	if _, err := db.Exec(`CREATE TABLE Players (
		Id TEXT NOT NULL UNIQUE,
		Name TEXT NOT NULL,
		WinCount INT NOT NULL DEFAULT 0,
		LoseCount INT NOT NULL DEFAULT 0,
		Coin DOUBLE NOT NULL DEFAULT 0.0,
		LastGetDailyTimestamp BIGINT
	);`); err != nil {
		return err
	}
	return nil
}

const (
	dbPosition = "/databases/twoonone.db"
)

func init() {
	_, err := os.Lstat(dbPosition)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(dbPosition)
			if err != nil {
				logger.Fatalln(err)
			}
			f.Close()
			if err := initDB(); err != nil {
				logger.Fatalln(err)
			}
		}
	}
	db, err := sql.Open("sqlite3", dbPosition)
	if err != nil {
		logger.Fatalln(err)
	}
	if err := db.Ping(); err != nil {
		logger.Fatalln(err)
	}
	database = db
}

func CreateAccount(id, name string, initialCoin float64) error {
	values := fmt.Sprintf(`( "%v", "%v", %v )`, id, name, initialCoin)
	if _, err := database.Exec(fmt.Sprintf(`INSERT INTO Players (Id, Name, Coin) VALUES %v;`, values)); err != nil {
		return err
	}
	return nil
}

func IncCoin(id string, coin float64) error {
	if _, err := database.Exec(fmt.Sprintf(`UPDATE Players SET Coin=Coin+%v WHERE Id="%v";`, coin, id)); err != nil {
		return err
	}
	return nil
}

func DecCoin(id string, coin float64) error {
	if _, err := database.Exec(fmt.Sprintf(`UPDATE Players SET Coin=Coin-%v WHERE Id="%v";`, coin, id)); err != nil {
		return err
	}
	return nil
}

func IncWinCount(id string, count int) error {
	if _, err := database.Exec(fmt.Sprintf(`UPDATE Players SET WinCount=WinCount+%v WHERE Id="%v";`, count, id)); err != nil {
		return err
	}
	return nil
}

func IncLoseCount(id string, count int) error {
	if _, err := database.Exec(fmt.Sprintf(`UPDATE Players SET LoseCount=LoseCount+%v WHERE Id="%v";`, count, id)); err != nil {
		return err
	}
	return nil
}

func UpdateLastGetDailyTimestamp(id string, t time.Time) error {
	if t.IsZero() {
		t = time.Now()
	}
	if _, err := database.Exec(fmt.Sprintf(`UPDATE Players SET LastGetDailyTimestamp=%v WHERE Id="%v";`, t.Unix(), id)); err != nil {
		return err
	}
	return nil
}

func GetPlayer(id string) (*twoonone_pb.PlayerAccountInfo, error) {
	row := database.QueryRow(fmt.Sprintf(`SELECT * FROM Players WHERE Id="%v";`, id))
	var (
		name                  string
		wincount, losecount   int
		coin                  float64
		lastGetDailyTimestamp *int64
	)
	if err := row.Scan(&ignore, &name, &wincount, &losecount, &coin, &lastGetDailyTimestamp); err != nil {
		return nil, err
	}
	if lastGetDailyTimestamp == nil {
		lastGetDailyTimestamp = new(int64)
	}
	return &twoonone_pb.PlayerAccountInfo{
		Id:                    id,
		Name:                  name,
		WinCount:              int32(wincount),
		LoseCount:             int32(losecount),
		Coin:                  coin,
		LastGetDailyTimestamp: *lastGetDailyTimestamp,
	}, nil
}

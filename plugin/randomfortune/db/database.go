package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/log"
)

var (
	database *sql.DB

	logger = log.Get()

	ignore any
)

const (
	dbPosition = "/databases/randomfortune.db"
)

func initDB() error {
	db, err := sql.Open("sqlite3", dbPosition)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE GetFortuneMembers (
		Id TEXT NOT NULL UNIQUE,
		LastGetFortuneTime INTEGER NOT NULL DEFAULT 0
	);`); err != nil {
		return err
	}
	return nil
}

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

func AddMember(id string) error {
	if _, err := database.Exec(fmt.Sprintf(`INSERT INTO GetFortuneMembers (Id) VALUES ("%v");`, id)); err != nil {
		return err
	}
	return nil
}

func UpdateLastGetFortuneTime(id string, ts int64) error {
	if ts == 0 {
		ts = time.Now().Unix()
	}
	if _, err := database.Exec(fmt.Sprintf(`UPDATE GetFortuneMembers SET LastGetFortuneTime=%v WHERE Id="%v";`, ts, id)); err != nil {
		return err
	}
	return nil
}

func GetLastGetFortuneTime(id string) (int64, error) {
	row := database.QueryRow(fmt.Sprintf(`SELECT * FROM GetFortuneMembers WHERE Id="%v"`, id))
	if row.Err() != nil {
		return 0, row.Err()
	}
	var (
		ts int64
	)
	if err := row.Scan(&ignore, &ts); err != nil {
		return 0, err
	}
	return ts, nil
}

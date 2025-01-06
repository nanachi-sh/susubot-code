package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nanachi-sh/susubot-code/plugin/uno/log"
	uno_pb "github.com/nanachi-sh/susubot-code/plugin/uno/protos/uno"
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
	if _, err := db.Exec(`PRAGMA foreign_keys = 1;`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE SourceEnum (
		Source TEXT PRIMARY KEY NOT NULL,
		Seq INTEGER NOT NULL UNIQUE
	);`); err != nil {
		return err
	}
	if _, err := db.Exec(`INSERT INTO SourceEnum (Source, Seq) VALUES ("QQ", 0);`); err != nil {
		return err
	}
	// 创建表
	if _, err := db.Exec(`CREATE TABLE Players (
		Id TEXT NOT NULL UNIQUE,
		Name TEXT NOT NULL,
		WinCount INT NOT NULL DEFAULT 0,
		LoseCount INT NOT NULL DEFAULT 0,
		Source TEXT NOT NULL REFERENCES SourceEnum(Source)
	);`); err != nil {
		return err
	}
	return nil
}

const (
	dbPosition = "/databases/uno.db"
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

func GetUser(userid string) (*uno_pb.PlayerAccountInfo, error) {
	row := database.QueryRow(fmt.Sprintf(`SELECT * FROM Players WHERE Id="%v";`, userid))
	var (
		name                string
		wincount, losecount int
	)
	if err := row.Scan(&ignore, &name, &wincount, &losecount); err != nil {
		return nil, err
	}
	return &uno_pb.PlayerAccountInfo{
		Id:        userid,
		Name:      name,
		WinCount:  int32(wincount),
		LoseCount: int32(losecount),
	}, nil
}

type Source int

const (
	Source_QQ = iota
)

func CreateUser(userid, username string, source Source) error {
	sourceStr := ""
	switch source {
	default:
		return errors.New("Source有误")
	case Source_QQ:
		sourceStr = "QQ"
	}
	values := fmt.Sprintf(`( "%v", "%v", "%v" )`, userid, username, sourceStr)
	if _, err := database.Exec(fmt.Sprintf(`INSERT INTO Players (Id, Name, Source) VALUES %v;`, values)); err != nil {
		return err
	}
	return nil
}

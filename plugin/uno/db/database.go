package db

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nanachi-sh/susubot-code/plugin/uno/log"
	uno_pb "github.com/nanachi-sh/susubot-code/plugin/uno/protos/uno"
	"github.com/twmb/murmur3"
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
		Source TEXT NOT NULL REFERENCES SourceEnum(Source),
		Hash TEXT NOT NULL UNIQUE
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
	values := fmt.Sprintf(`( "%v", "%v", "%v", "%v" )`, userid, username, sourceStr, hash(userid))
	if _, err := database.Exec(fmt.Sprintf(`INSERT INTO Players (Id, Name, Source, Hash) VALUES %v;`, values)); err != nil {
		return err
	}
	return nil
}

func hash(id string) string {
	h1, h2 := murmur3.SeedStringSum128(rand.Uint64(), rand.Uint64(), id)
	return fmt.Sprintf("%v%v", strconv.FormatUint(h1, 16), strconv.FormatUint(h2, 16))
}

type UserInfo struct {
	AI   *uno_pb.PlayerAccountInfo
	Hash string
	S    Source
}

func FindUser(userid, userhash string) (*UserInfo, error) {
	key := ""
	if userid != "" {
		key = userid
	} else if userhash != "" {
		key = userhash
	} else {
		return nil, errors.New("UserId和UserHash都为空")
	}
	row := database.QueryRow(fmt.Sprintf(`SELECT * FROM Players WHERE Id="%v";`, key))
	var (
		name                string
		wincount, losecount int
		s                   Source
		hash                string
	)
	if err := row.Scan(&ignore, &name, &wincount, &losecount, &s, &hash); err != nil {
		return nil, err
	}
	return &UserInfo{
		AI: &uno_pb.PlayerAccountInfo{
			Id:        userid,
			Name:      name,
			WinCount:  int32(wincount),
			LoseCount: int32(losecount),
		},
		Hash: hash,
		S:    s,
	}, nil
}

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
		Password TEXT NOT NULL,
		SEED1 INTEGER NOT NULL,
		SEED2 INTEGER NOT NULL,
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

func CreateUser(userid, username string, source uno_pb.Source) error {
	sourceStr := ""
	switch source {
	default:
		return errors.New("Source有误")
	case uno_pb.Source_QQ:
		sourceStr = "QQ"
	}
	ss, pwd := generatePassword(userid, username)
	values := fmt.Sprintf(`( "%v", "%v", "%v", "%v", "%v", %v, %v )`, userid, username, sourceStr, hash(userid), pwd, ss[0], ss[1])
	if _, err := database.Exec(fmt.Sprintf(`INSERT INTO Players (Id, Name, Source, Hash, Password, SEED1, SEED2) VALUES %v;`, values)); err != nil {
		return err
	}
	return nil
}

func hash(id string) string {
	h1, h2 := murmur3.SeedStringSum128(rand.Uint64(), rand.Uint64(), id)
	return fmt.Sprintf("%v%v", strconv.FormatUint(h1, 16), strconv.FormatUint(h2, 16))
}

func generatePassword(id, name string) ([2]uint64, string) {
	s1, s2 := rand.Uint64(), rand.Uint64()
	h1, h2 := murmur3.SeedStringSum128(s1, s2, id+name)
	return [2]uint64{s1, s2}, fmt.Sprintf("%v%v", strconv.FormatUint(h1, 16), strconv.FormatUint(h2, 16))
}

func CalcPassword(s1, s2 uint64, id, name string) string {
	h1, h2 := murmur3.SeedStringSum128(s1, s2, id+name)
	return fmt.Sprintf("%v%v", strconv.FormatUint(h1, 16), strconv.FormatUint(h2, 16))
}

type UserInfo struct {
	AI    *uno_pb.PlayerAccountInfo
	Hash  string
	S     uno_pb.Source
	SEEDs [2]uint64
}

func FindUser(userid, userhash string) (*UserInfo, error) {
	key := ""
	value := ""
	if userid != "" {
		key = "Id"
		value = userid
	} else if userhash != "" {
		key = "Hash"
		value = userhash
	} else {
		return nil, errors.New("UserId和UserHash都为空")
	}
	row := database.QueryRow(fmt.Sprintf(`SELECT * FROM Players WHERE %v="%v";`, key, value))
	var (
		name         string
		s            uno_pb.Source
		hash         string
		seed1, seed2 uint64
		passwordHASH string
	)
	if err := row.Scan(&ignore, &name, &passwordHASH, &seed1, &seed2, &s, &hash); err != nil {
		return nil, err
	}
	return &UserInfo{
		AI: &uno_pb.PlayerAccountInfo{
			Id:   userid,
			Name: name,
		},
		Hash: hash,
		S:    s,
	}, nil
}

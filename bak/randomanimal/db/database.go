package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/log"
	randomanimal_pb "github.com/nanachi-sh/susubot-code/plugin/randomanimal/protos/randomanimal"
)

var (
	database *sql.DB

	logger = log.Get()

	ignore any
)

const (
	dbPosition = "/databases/randomanimal.db"
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
	if _, err := db.Exec(`CREATE TABLE TypeEnum (
		Type TEXT PRIMARY KEY NOT NULL,
		Seq INTEGER NOT NULL UNIQUE
	);`); err != nil {
		return err
	}
	if _, err := db.Exec(`INSERT INTO TypeEnum (Type, Seq) VALUES ("Image", 0);`); err != nil {
		return err
	}
	if _, err := db.Exec(`INSERT INTO TypeEnum (Type, Seq) VALUES ("Video", 1);`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE DogCache (
		IDHash TEXT NOT NULL UNIQUE,
		AssetHash TEXT NOT NULL,
		AssetType TEXT NOT NULL REFERENCES TypeEnum(AssetType)
	);`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE CatCache (
		IDHash TEXT NOT NULL UNIQUE,
		AssetHash TEXT NOT NULL,
		AssetType TEXT NOT NULL REFERENCES TypeEnum(AssetType)
	);`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE FoxCache (
		IDHash TEXT NOT NULL UNIQUE,
		AssetHash TEXT NOT NULL,
		AssetType TEXT NOT NULL REFERENCES TypeEnum(AssetType)
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

type Cache struct {
	AssetHash string
	Type      randomanimal_pb.Type
}

type AnimalType int

const (
	Cat AnimalType = iota
	Dog
	Fox
	Duck
	Chicken_CXK
)

func AnimalType2TableName(at AnimalType) string {
	switch at {
	case Dog:
		return "DogCache"
	case Fox:
		return "FoxCache"
	case Cat:
		return "CatCache"
	default:
		return ""
	}
}

func FindCache(idhash string, at AnimalType) (*Cache, error) {
	tn := AnimalType2TableName(at)
	if tn == "" {
		return nil, errors.New("不支持缓存")
	}
	row := database.QueryRow(fmt.Sprintf(`SELECT * FROM %v WHERE IDHash="%v";`, tn, idhash))
	if row.Err() != nil {
		return nil, row.Err()
	}
	c := new(Cache)
	var (
		assethash string
		typeStr   string
	)
	if err := row.Scan(&ignore, &assethash, &typeStr); err != nil {
		return nil, err
	}
	switch typeStr {
	case "Image":
		c.Type = randomanimal_pb.Type_Image
	case "Video":
		c.Type = randomanimal_pb.Type_Video
	default:
		return nil, errors.New("unexpected error")
	}
	c.AssetHash = assethash
	return c, nil
}

func DeleteCache(idhash string, at AnimalType) error {
	tn := AnimalType2TableName(at)
	if tn == "" {
		return errors.New("不支持缓存")
	}
	r, err := database.Exec(fmt.Sprintf(`DELETE FROM %v WHERE IDHash="%v";`, tn, idhash))
	if err != nil {
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("无匹配缓存")
	}
	return nil
}

func AddCache(idhash string, at AnimalType, c Cache) error {
	tn := AnimalType2TableName(at)
	if tn == "" {
		return errors.New("不支持缓存")
	}
	if _, err := database.Exec(fmt.Sprintf(`INSERT INTO %v (IDHash, AssetHash, AssetType) VALUES ("%v", "%v", "%v");`, tn, idhash, c.AssetHash, c.Type.String())); err != nil {
		return err
	}
	return nil
}

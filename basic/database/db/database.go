package db

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nanachi-sh/susubot-code/basic/database/log"
)

var (
	uno_database *sql.DB
	logger       = log.Get()

	ignore any
)

func initDB() error {
	db, err := sql.Open("sqlite3", uno_dbPosition)
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
		Password TEXT NOT NULL,
		Salt TEXT NOT NULL,
		WinCount INTEGER NOT NULL DEFAULT 0,
		LoseCount INTEGER NOT NULL DEFAULT 0
	);`); err != nil {
		return err
	}
	return nil
}

const (
	uno_dbPosition = "/databases/uno.db"
)

func init() {
	_, err := os.Lstat(uno_dbPosition)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(uno_dbPosition)
			if err != nil {
				logger.Fatalln(err)
			}
			f.Close()
			if err := initDB(); err != nil {
				logger.Fatalln(err)
			}
		}
	}
	db, err := sql.Open("sqlite3", uno_dbPosition)
	if err != nil {
		logger.Fatalln(err)
	}
	if err := db.Ping(); err != nil {
		logger.Fatalln(err)
	}
	uno_database = db
}

var dict = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateSalt() string {
	length := 6
	ret := new(strings.Builder)
	for i := 0; i < length; i++ {
		ret.WriteByte(dict[rand.Intn(len(dict))])
	}
	return ret.String()
}

func encryptPassword(password, salt string) string {
	return fmt.Sprintf("%x", password+salt)
}

// 数据库中顺序：
// 0: Id
// 1: Name
// 2: Password
// 3: Salt
// 4: WinCount
// 5: LoseCount
type Uno_UserInfo struct {
	Id, Name            string
	WinCount, LoseCount int

	password, salt string
}

func Uno_CreateUser(userid, username, password string) error {
	salt := generateSalt()
	passwordEncrypt := encryptPassword(password, salt)
	values := fmt.Sprintf(`( "%v", "%v", "%v" )`, userid, username, passwordEncrypt)
	if _, err := uno_database.Exec(fmt.Sprintf(`INSERT INTO Players (Id, Name, Password) VALUES %v;`, values)); err != nil {
		return err
	}
	return nil
}

func Uno_GetUser(userid string) (Uno_UserInfo, error) {
	row := uno_database.QueryRow(fmt.Sprintf(`SELECT * FROM Players WHERE Id="%v";`, userid))
	if err := row.Err(); err != nil {
		return Uno_UserInfo{}, err
	}
	var (
		id, name            string
		wincount, losecount int
		password, salt      string
	)
	if err := row.Scan(&id, &name, &password, &salt, &wincount, &losecount); err != nil {
		return Uno_UserInfo{}, err
	}
	return Uno_UserInfo{
		Id:        id,
		Name:      name,
		WinCount:  wincount,
		LoseCount: losecount,
	}, nil
}

func Uno_VerifyUser(userid string, password string) (bool, error) {
	ui, err := Uno_GetUser(userid)
	if err != nil {
		return false, err
	}
	return uno_VerifyUser(ui, password), nil
}

func uno_VerifyUser(ui Uno_UserInfo, password string) bool {
	return encryptPassword(password, ui.salt) == ui.password
}

func Uno_UpdateUser(ui Uno_UserInfo) error {
	if _, err := Uno_GetUser(ui.Id); err != nil {
		return err
	}
	if _, err := uno_database.Exec(fmt.Sprintf(`UPDATE Players SET Name="%v", WinCount=%v, LoseCount=%v WHERE Id="%v";`, ui.Name, ui.WinCount, ui.LoseCount, ui.Id)); err != nil {
		return err
	}
	return nil
}

func Uno_ChangePassword(userid, newPwd, oldPwd string) (bool, error) {
	ui, err := Uno_GetUser(userid)
	if err != nil {
		return false, err
	}
	if !uno_VerifyUser(ui, oldPwd) {
		return false, nil
	}
	salt := generateSalt()
	newPwdEncrypt := encryptPassword(newPwd, salt)
	if _, err := uno_database.Exec(fmt.Sprintf(`UPDATE Players SET Password="%v", Salt="%v" WHERE Id="%v";`, newPwdEncrypt, salt, ui.Id)); err != nil {
		return false, err
	}
	return true, nil
}

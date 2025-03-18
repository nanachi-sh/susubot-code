package configs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strconv"

	dex "github.com/dexidp/dex/api/v2"
	"github.com/go-ldap/ldap/v3"
	"github.com/mojocn/base64Captcha"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/model/applications"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/utils"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	DEBUG bool

	Captcha = base64Captcha.NewCaptcha(base64Captcha.NewDriverString(
		64,
		128,
		3,
		base64Captcha.OptionShowHollowLine,
		4,
		utils.Dict_Number,
		nil,
		base64Captcha.DefaultEmbeddedFonts,
		nil,
	), base64Captcha.DefaultMemStore)

	Redis *redis.Redis

	Model_Applications applications.TwoononeModel

	Call_Dex dex.DexClient

	LDAP *ldap.Conn

	HTTP_LISTEN_PORT int

	SMTP_USERNAME string
	SMTP_AUTHCODE string

	REDIS_HOST     netip.Addr
	REDIS_PORT     int
	REDIS_PASSWORD string

	MYSQL_HOST     netip.Addr
	MYSQL_PORT     int
	MYSQL_USERNAME string
	MYSQL_PASSWORD string

	LDAP_HOST     netip.Addr
	LDAP_PORT     int
	LDAP_USERNAME string
	LDAP_PASSWORD string
	LDAP_BASIC_DN string

	DEX_HOST netip.Addr
	DEX_PORT int

	APIServer_Config string = "api_server.yaml"
)

const (
	ConfigDir = "/config"
)

// 获取环境变量
func init() {
	port, err := utils.EnvPortToPort("HTTP_LISTEN_PORT")
	if err != nil {
		logger.Fatalf("HTTP监听端口获取出错，err：%v", err)
	}
	HTTP_LISTEN_PORT = int(port)

	str, err := utils.EnvToString("SMTP_USERNAME")
	if err != nil {
		logger.Fatalln(err)
	}
	SMTP_USERNAME = str
	str, err = utils.EnvToString("SMTP_AUTHCODE")
	if err != nil {
		logger.Fatalln(err)
	}
	SMTP_AUTHCODE = str

	str, err = utils.EnvToString("REDIS_HOST")
	if err != nil {
		logger.Fatalln(err)
	}
	addr, err := utils.ResolvIP(str)
	if err != nil {
		logger.Fatalln(err)
	}
	REDIS_HOST = addr

	port, err = utils.EnvPortToPort("REDIS_PORT")
	if err != nil {
		logger.Fatalln(err)
	}
	REDIS_PORT = int(port)

	str, err = utils.EnvToString("REDIS_PASSWORD")
	if err != nil {
		logger.Fatalln(err)
	}
	REDIS_PASSWORD = str

	str, err = utils.EnvToString("MYSQL_HOST")
	if err != nil {
		logger.Fatalln(err)
	}
	addr, err = utils.ResolvIP(str)
	if err != nil {
		logger.Fatalln(err)
	}
	MYSQL_HOST = addr

	port, err = utils.EnvPortToPort("MYSQL_PORT")
	if err != nil {
		logger.Fatalln(err)
	}
	MYSQL_PORT = int(port)

	str, err = utils.EnvToString("MYSQL_USERNAME")
	if err != nil {
		logger.Fatalln(err)
	}
	MYSQL_USERNAME = str

	str, err = utils.EnvToString("MYSQL_PASSWORD")
	if err != nil {
		logger.Fatalln(err)
	}
	MYSQL_PASSWORD = str

	str, err = utils.EnvToString("LDAP_HOST")
	if err != nil {
		logger.Fatalln(err)
	}
	addr, err = utils.ResolvIP(str)
	if err != nil {
		logger.Fatalln(err)
	}
	LDAP_HOST = addr

	port, err = utils.EnvPortToPort("LDAP_PORT")
	if err != nil {
		logger.Fatalln(err)
	}
	LDAP_PORT = int(port)

	str, err = utils.EnvToString("LDAP_USERNAME")
	if err != nil {
		logger.Fatalln(err)
	}
	LDAP_USERNAME = str

	str, err = utils.EnvToString("LDAP_PASSWORD")
	if err != nil {
		logger.Fatalln(err)
	}
	LDAP_PASSWORD = str
	str, err = utils.EnvToString("LDAP_BASIC_DN")
	if err != nil {
		logger.Fatalln(err)
	}
	LDAP_BASIC_DN = str

	str, err = utils.EnvToString("DEX_HOST")
	if err != nil {
		logger.Fatalln(err)
	}
	addr, err = utils.ResolvIP(str)
	if err != nil {
		logger.Fatalln(err)
	}
	DEX_HOST = addr

	port, err = utils.EnvPortToPort("DEX_PORT")
	if err != nil {
		logger.Fatalln(err)
	}
	DEX_PORT = int(port)

	if d := os.Getenv("DEBUG"); d != "" {
		if debug, err := strconv.ParseBool(d); err != nil {
			logger.Fatalln("Debug状态设置不正确")
		} else {
			DEBUG = debug
		}
	}
}

// 初始化gRPC配置
func init() {
	m := map[string]any{
		"Name": "accountmanager-api",
		"Host": "0.0.0.0",
		"Port": HTTP_LISTEN_PORT,
		"Log": map[string]any{
			"MaxContentLength": 16 * 1024,
		},
	}
	buf, err := json.Marshal(m)
	if err != nil {
		logger.Fatalln(err)
	}
	if err := os.WriteFile(APIServer_Config, buf, 0744); err != nil {
		logger.Fatalln(err)
	}
}

// 初始化数据库相关
func init() {
	Redis = redis.New(fmt.Sprintf("%s:%d", REDIS_HOST, REDIS_PORT), redis.WithPass(REDIS_PASSWORD))

	sqlconn := sqlx.NewMysql(fmt.Sprintf("%s:%s@tcp(%s:%d)/applications", MYSQL_USERNAME, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_PORT))
	Model_Applications = applications.NewTwoononeModel(sqlconn)

	ldapconn, err := ldap.DialURL(fmt.Sprintf("ldap://%s:%d", LDAP_HOST, LDAP_PORT))
	if err != nil {
		logger.Fatalln(err)
	}
	if err := ldapconn.Bind(LDAP_USERNAME, LDAP_PASSWORD); err != nil {
		logger.Fatalln(err)
	}
	LDAP = ldapconn
}

// call
func init() {
	client, err := grpc.NewClient(fmt.Sprintf("%s:%d", DEX_HOST, DEX_PORT))
	if err != nil {
		logger.Fatalln(err)
	}
	Call_Dex = dex.NewDexClient(client)
}

func GRPCOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	return opts
}

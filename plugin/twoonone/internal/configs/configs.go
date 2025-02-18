package configs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	twoononemodel "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/model/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/utils"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	DEBUG bool

	HTTPAPI_LISTEN_PORT int

	DATABASE_IP       netip.Addr
	DATABASE_PORT     int
	DATABASE_USER     string
	DATABASE_PASSWORD string

	REDIS_IP       netip.Addr
	REDIS_PORT     int
	REDIS_PASSWORD string

	LDAP_HOST     string
	LDAP_PORT     int
	LDAP_TLS      bool = true
	LDAP_USER     string
	LDAP_PASSWORD string
	LDAP_DN       string

	OIDC_ISSUER        string
	OIDC_CLIENT_ID     string
	OIDC_CLIENT_SECRET string
	OIDC_REDIRECT      string

	MIDDLEWARE_AuthHandlerStatus bool = true
	MIDDLEWARE_SQLHandlerStatus  bool = true

	Model_TwoOnOne twoononemodel.TwoononeModel
	Redis          *redis.Redis
	LDAP           ldap.Client

	DefaultCtx = context.Background()

	APIServer_Config string = "api_server.json"
)

const (
	ConfigDir                  = "/config"
	Randomfortune_HashListFile = ConfigDir + "/fortune_HashList.json"
)

// 获取环境变量
func init() {
	port, err := utils.EnvPortToPort("HTTPAPI_LISTEN_PORT")
	if err != nil {
		logger.Fatalf("HTTP API监听端口获取失败，err: %v", err)
	}
	HTTPAPI_LISTEN_PORT = int(port)

	if d := os.Getenv("DEBUG"); d != "" {
		if debug, err := strconv.ParseBool(d); err != nil {
			logger.Fatalln("Debug状态设置不正确")
		} else {
			DEBUG = debug
		}
	}

	if host := os.Getenv("DATABASE_HOST"); host == "" {
		logger.Fatalln("database host未设置")
	} else {
		addr, err := utils.ResolvIP(host)
		if err != nil {
			logger.Fatalf("database host获取失败，err: %v", err)
		}
		DATABASE_IP = addr
	}
	port, err = utils.EnvPortToPort("DATABASE_PORT")
	if err != nil {
		logger.Fatalf("database port获取失败, err: %v", err)
	}
	DATABASE_PORT = int(port)
	if DATABASE_USER = os.Getenv("DATABASE_USER"); DATABASE_USER == "" {
		logger.Fatalln("Database 用户名未设置")
	}
	if DATABASE_PASSWORD = os.Getenv("DATABASE_PASSWORD"); DATABASE_PASSWORD == "" {
		logger.Fatalln("Database 用户密码未设置")
	}

	if host := os.Getenv("REDIS_HOST"); host == "" {
		logger.Fatalln("redis host未设置")
	} else {
		addr, err := utils.ResolvIP(host)
		if err != nil {
			logger.Fatalf("redis host获取失败，err: %v", err)
		}
		REDIS_IP = addr
	}
	if REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD"); REDIS_PASSWORD == "" {
		logger.Fatalln("Redis 用户密码未设置")
	}
	port, err = utils.EnvPortToPort("REDIS_PORT")
	if err != nil {
		logger.Fatalln("redis port获取失败")
	}
	REDIS_PORT = int(port)

	if str, err := utils.EnvToString("LDAP_HOST"); err != nil {
		logger.Fatalln(err)
	} else {
		LDAP_HOST = str
	}
	port, err = utils.EnvPortToPort("LDAP_PORT")
	if err != nil {
		logger.Fatalln("ldap port获取失败")
	}
	LDAP_PORT = int(port)
	if b, err := strconv.ParseBool(os.Getenv("LDAP_TLS")); err == nil {
		LDAP_TLS = b
	}
	if LDAP_PASSWORD = os.Getenv("LDAP_PASSWORD"); LDAP_PASSWORD == "" {
		logger.Fatalln("LDAP password 未设置")
	}
	if LDAP_USER = os.Getenv("LDAP_USER"); LDAP_USER == "" {
		logger.Fatalln("LDAP user 未设置")
	}
	if str, err := utils.EnvToString("LDAP_DN"); err != nil {
		logger.Fatalln(err)
	} else {
		LDAP_DN = str
	}

	if OIDC_CLIENT_ID = os.Getenv("OIDC_CLIENT_ID"); OIDC_CLIENT_ID == "" {
		logger.Fatalln("OIDC CLIENT ID未设置")
	}
	if OIDC_CLIENT_SECRET = os.Getenv("OIDC_CLIENT_SECRET"); OIDC_CLIENT_SECRET == "" {
		logger.Fatalln("OIDC CLIENT SECRET未设置")
	}
	if OIDC_ISSUER = os.Getenv("OIDC_ISSUER"); OIDC_ISSUER == "" {
		logger.Fatalln("OIDC ISSUER未设置")
	}
	if OIDC_REDIRECT = os.Getenv("OIDC_REDIRECT"); OIDC_REDIRECT == "" {
		logger.Fatalln("OIDC REDIRECT未设置")
	}

	if b, err := strconv.ParseBool(os.Getenv("MIDDLEWARE_AUTH_HANDLER_STATUS")); err == nil {
		MIDDLEWARE_AuthHandlerStatus = b
	}
}

// 初始化gRPC配置
func init() {
	// m := map[string]any{
	// 	"Name":     "connector.rpc",
	// 	"ListenOn": fmt.Sprintf("0.0.0.0:%d", GRPC_LISTEN_PORT),
	// 	"Log": map[string]any{
	// 		"MaxContentLength": 16 * 1024,
	// 	},
	// }
	// buf, err := json.Marshal(m)
	// if err != nil {
	// 	logger.Fatalln(err)
	// }
	// if err := os.WriteFile(RPCServer_Config, buf, 0744); err != nil {
	// 	logger.Fatalln(err)
	// }
	m := map[string]any{
		"Name": "twoonone-api",
		"Host": "0.0.0.0",
		"Port": HTTPAPI_LISTEN_PORT,
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

// 初始化SQL Models
func init() {
	sqlconn := sqlx.NewMysql(fmt.Sprintf("%s:%s@tcp(%s:%d)/applications", DATABASE_USER, DATABASE_PASSWORD, DATABASE_IP, DATABASE_PORT))
	Model_TwoOnOne = twoononemodel.NewTwoononeModel(sqlconn)
	Redis = redis.New(fmt.Sprintf("%s:%d", REDIS_IP, REDIS_PORT), redis.WithPass("lsusu"))
	ldap_url := ""
	if LDAP_TLS {
		ldap_url = fmt.Sprintf("ldaps://%s:%d", LDAP_HOST, LDAP_PORT)
	} else {
		ldap_url = fmt.Sprintf("ldap://%s:%d", LDAP_HOST, LDAP_PORT)
	}
	ldapconn, err := ldap.DialURL(ldap_url)
	if err != nil {
		logger.Fatalln(err)
	}
	LDAP = ldapconn
	if err := LDAP.Bind(LDAP_USER, LDAP_PASSWORD); err != nil {
		logger.Fatalln(err)
	}
}

func GRPCOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	return opts
}

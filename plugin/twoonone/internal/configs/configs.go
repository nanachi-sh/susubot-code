package configs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strconv"

	twoononemodel "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/model/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/utils"
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

	OIDC_ISSUER        string
	OIDC_CLIENT_ID     string
	OIDC_CLIENT_SECRET string
	OIDC_REDIRECT      string

	JWT_SignKey string

	Model_TwoOnOne twoononemodel.TwoononeModel

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

	if str, err := utils.EnvToString("JWT_SIGN_KEY"); err != nil {
		logger.Fatalln(err)
	} else {
		JWT_SignKey = str
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
		"Middlewares": map[string]any{
			"Timeout": false,
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
}

func GRPCOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	return opts
}

package configs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strconv"

	randomfortunemodel "github.com/nanachi-sh/susubot-code/plugin/randomfortune/internal/model/randomfortune"
	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/internal/utils"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	DEBUG bool

	GRPC_LISTEN_PORT int

	GATEWAY_IP        netip.Addr
	GATEWAY_GRPC_PORT int

	DATABASE_IP       netip.Addr
	DATABASE_PORT     int
	DATABASE_USER     string
	DATABASE_PASSWORD string

	ASSETS_URL string

	Model_Randomfortune randomfortunemodel.UsersModel

	DefaultCtx = context.Background()

	RPCServer_Config string = "rpc_server.json"
)

const (
	ConfigDir                  = "/config"
	Randomfortune_HashListFile = ConfigDir + "/fortune_HashList.json"
)

// 获取环境变量
func init() {
	port, err := utils.EnvPortToPort("GRPC_LISTEN_PORT")
	if err != nil {
		logger.Fatalf("grpc监听端口获取失败，err: %v", err)
	}
	GRPC_LISTEN_PORT = int(port)

	if d := os.Getenv("DEBUG"); d != "" {
		if debug, err := strconv.ParseBool(d); err != nil {
			logger.Fatalln("Debug状态设置不正确")
		} else {
			DEBUG = debug
		}
	}

	if host := os.Getenv("GATEWAY_HOST"); host == "" {
		logger.Fatalln("gateway host未设置")
	} else {
		addr, err := utils.ResolvIP(host)
		if err != nil {
			logger.Fatalf("gateway host获取失败，err: %v", err)
		}
		GATEWAY_IP = addr
	}

	if ASSETS_URL = os.Getenv("ASSETS_URL"); ASSETS_URL == "" {
		logger.Fatalln("assets url未设置")
	}

	port, err = utils.EnvPortToPort("GATEWAY_GRPC_PORT")
	if err != nil {
		logger.Fatalf("gateway grpc port获取失败，err: %v", err)
	}
	GATEWAY_GRPC_PORT = int(port)

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
}

// 初始化gRPC配置
func init() {
	m := map[string]any{
		"Name":     "connector.rpc",
		"ListenOn": fmt.Sprintf("0.0.0.0:%d", GRPC_LISTEN_PORT),
		"Log": map[string]any{
			"MaxContentLength": 16 * 1024,
		},
	}
	buf, err := json.Marshal(m)
	if err != nil {
		logger.Fatalln(err)
	}
	if err := os.WriteFile(RPCServer_Config, buf, 0744); err != nil {
		logger.Fatalln(err)
	}
}

// 初始化SQL Models
func init() {
	conn := sqlx.NewMysql(fmt.Sprintf("%s:%s@tcp(%s:%d)/randomfortune", DATABASE_USER, DATABASE_PASSWORD, DATABASE_IP, DATABASE_PORT))
	Model_Randomfortune = randomfortunemodel.NewUsersModel(conn)
}

func GRPCOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	return opts
}

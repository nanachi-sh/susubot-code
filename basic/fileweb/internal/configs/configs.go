package configs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/nanachi-sh/susubot-code/basic/fileweb/internal/utils"
	"google.golang.org/grpc"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	DEBUG bool

	GRPC_LISTEN_PORT int
	HTTP_LISTEN_PORT int

	SEED1, SEED2 uint64

	RPCServer_Config string = "rpc_server.yaml"
)

const (
	WebDir    = "/var/www/html"
	ConfigDir = "/config"
)

// 获取环境变量
func init() {
	port, err := utils.EnvPortToPort("GRPC_LISTEN_PORT")
	if err != nil {
		logger.Fatalf("grpc监听端口获取出错，err：%v", err)
	}
	GRPC_LISTEN_PORT = int(port)

	port, err = utils.EnvPortToPort("HTTP_LISTEN_PORT")
	if err != nil {
		logger.Fatalf("HTTP监听端口获取出错，err：%v", err)
	}
	HTTP_LISTEN_PORT = int(port)

	if s, err := strconv.ParseUint(os.Getenv("SEED1"), 10, 0); err != nil {
		logger.Fatalln("SEED未设置或设置不正确")
	} else {
		SEED1 = s
	}
	if s, err := strconv.ParseUint(os.Getenv("SEED2"), 10, 0); err != nil {
		logger.Fatalln("SEED未设置或设置不正确")
	} else {
		SEED2 = s
	}

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

func GRPCOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	opts = append(opts, grpc.MaxRecvMsgSize(128*1024*1024))
	return opts
}

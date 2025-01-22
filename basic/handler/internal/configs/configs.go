package configs

import (
	"context"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strconv"

	connectorclient "github.com/nanachi-sh/susubot-code/basic/handler/internal/caller/connector"
	filewebclient "github.com/nanachi-sh/susubot-code/basic/handler/internal/caller/fileweb"
	"github.com/nanachi-sh/susubot-code/basic/handler/internal/types"
	"github.com/nanachi-sh/susubot-code/basic/handler/internal/utils"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	DEBUG bool

	GRPC_LISTEN_PORT  int
	GATEWAY_IP        netip.Addr
	GATEWAY_GRPC_PORT int
	ASSETS_URL        string

	RPCServer_Config string = "rpc_server.yaml"
	RPCClient_Config string = "rpc_client.yaml"

	Call_Connector connectorclient.Connector
	Call_FileWeb   filewebclient.FileWeb
	DefaultCtx     context.Context = context.Background()
)

const (
	ConfigDir = "/config"
)

// 获取环境变量
func init() {
	port, err := utils.EnvPortToPort(os.Getenv("GRPC_LISTEN_PORT"))
	if err != nil {
		logger.Fatalf("gRPC监听端口有误，Err: %s\n", err.Error())
	}
	GRPC_LISTEN_PORT = int(port)

	if d := os.Getenv("DEBUG"); d != "" {
		if debug, err := strconv.ParseBool(d); err != nil {
			logger.Fatalln("Debug状态设置不正确")
		} else {
			DEBUG = debug
		}
	}

	if gatewayHost := os.Getenv("GATEWAY_HOST"); gatewayHost == "" {
		logger.Fatalln("Gateway API Host为空")
	} else {
		ip, err := utils.ResolvIP(gatewayHost)
		if err != nil {
			logger.Fatalf("Gateway API Host解析出错，Err: %s\n", err.Error())
		}
		GATEWAY_IP = ip
	}
	port, err = utils.EnvPortToPort(os.Getenv("GATEWAY_GRPC_PORT"))
	if err != nil {
		logger.Fatalf("Gateway gRPC服务端口有误，Err: %s\n", err.Error())
	}
	GATEWAY_GRPC_PORT = int(port)

	if ASSETS_URL = os.Getenv("ASSETS_URL"); ASSETS_URL == "" {
		logger.Fatalln("Assets URL未设置")
	}
}

// 初始化gRPC配置
func init() {
	config := fmt.Sprintf(`Name: connector.rpc
ListenOn: 0.0.0.0:%d`, GRPC_LISTEN_PORT)
	if err := os.WriteFile(RPCServer_Config, []byte(config), 0744); err != nil {
		logger.Fatalln(err)
	}
	config = fmt.Sprintf(`Target: %s:%d`, GATEWAY_IP.String(), GATEWAY_GRPC_PORT)
	if err := os.WriteFile(RPCClient_Config, []byte(config), 0744); err != nil {
		logger.Fatalln(err)
	}
}

// 初始化gRPC Callers
func init() {
	var c types.Config
	if err := conf.LoadConfig(RPCClient_Config, &c); err != nil {
		logger.Fatalln(err)
	}
	client, err := zrpc.NewClient(c.RpcClientConf)
	if err != nil {
		logger.Fatalln(err)
	}
	Call_Connector = connectorclient.NewConnector(client)
	Call_FileWeb = filewebclient.NewFileWeb(client)
}

func GRPCOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	return opts
}

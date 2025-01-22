package configs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strconv"

	connectorclient "github.com/nanachi-sh/susubot-code/basic/verifier/internal/caller/connector"
	requesthandler "github.com/nanachi-sh/susubot-code/basic/verifier/internal/caller/handler/request"
	responsehandler "github.com/nanachi-sh/susubot-code/basic/verifier/internal/caller/handler/response"
	mock_connectorclient "github.com/nanachi-sh/susubot-code/basic/verifier/internal/mock/connector"
	"github.com/nanachi-sh/susubot-code/basic/verifier/internal/types"
	"github.com/nanachi-sh/susubot-code/basic/verifier/internal/utils"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	DEBUG bool

	GRPC_LISTEN_PORT  int
	HTTP_LISTEN_PORT  int
	GATEWAY_IP        netip.Addr
	GATEWAY_GRPC_PORT int

	Call_Connector        connectorclient.Connector
	Call_Handler_Request  requesthandler.RequestHandler
	Call_Handler_Response responsehandler.ResponseHandler

	RPCServer_Config  string = "rpc_server.yaml"
	RPCClient_Config  string = "rpc_client.yaml"
	RPCGateway_Config string = "rpc_gateway.yaml"

	DefaultCtx context.Context = context.Background()
)

const (
	ConfigDir = "/config"
)

// 获取环境变量
func init() {
	port, err := utils.EnvPortToPort("GRPC_LISTEN_PORT")
	if err != nil {
		logger.Fatalf("grpc监听端口获取失败，err： %v", err)
	}
	GRPC_LISTEN_PORT = int(port)

	port, err = utils.EnvPortToPort("HTTP_LISTEN_PORT")
	if err != nil {
		logger.Fatalf("http监听端口获取失败，err: %v", err)
	}
	HTTP_LISTEN_PORT = int(port)

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
	port, err = utils.EnvPortToPort("GATEWAY_GRPC_PORT")
	if err != nil {
		logger.Fatalf("gateway grpc port获取失败，err: %v", err)
	}
	GATEWAY_GRPC_PORT = int(port)
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
	if !DEBUG {
		var c types.Config
		if err := conf.LoadConfig(RPCClient_Config, &c); err != nil {
			logger.Fatalln(err)
		}
		client, err := zrpc.NewClient(c.RpcClientConf)
		if err != nil {
			logger.Fatalln(err)
		}
		Call_Connector = connectorclient.NewConnector(client)
		Call_Handler_Request = requesthandler.NewRequestHandler(client)
		Call_Handler_Response = responsehandler.NewResponseHandler(client)
	} else {
		var c types.Config
		if err := conf.LoadConfig(RPCClient_Config, &c); err != nil {
			logger.Fatalln(err)
		}
		cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
		if err != nil {
			logger.Fatalln(err)
		}
		caPool := x509.NewCertPool()
		caCert, err := os.ReadFile("ca.crt")
		if err != nil {
			logger.Fatalln(err)
		}
		caPool.AppendCertsFromPEM(caCert)
		cred := credentials.NewTLS(&tls.Config{
			RootCAs:      caPool,
			Certificates: []tls.Certificate{cert},
			ServerName:   "mtls.susu",
		})
		client, err := zrpc.NewClient(c.RpcClientConf, zrpc.WithDialOption(grpc.WithTransportCredentials(cred)))
		if err != nil {
			logger.Fatalln(err)
		}
		Call_Connector = mock_connectorclient.DefaultMock()
		Call_Handler_Request = requesthandler.NewRequestHandler(client)
		Call_Handler_Response = responsehandler.NewResponseHandler(client)
	}
}

func GRPCOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	return opts
}

package configs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"regexp"
	"strconv"
	"time"

	connectorclient "github.com/nanachi-sh/susubot-code/basic/verifier/internal/caller/connector"
	requesthandler "github.com/nanachi-sh/susubot-code/basic/verifier/internal/caller/handler/request"
	responsehandler "github.com/nanachi-sh/susubot-code/basic/verifier/internal/caller/handler/response"
	"github.com/nanachi-sh/susubot-code/basic/verifier/internal/types"
	"github.com/nanachi-sh/susubot-code/basic/verifier/pkg/utils"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	DEBUG             bool
	GRPC_LISTEN_PORT  int
	HTTP_LISTEN_PORT  int
	GRPC_mTLS         bool
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

	CertsDir           = ConfigDir + "/certs"
	GRPCServerCertFile = CertsDir + "/mtls_server.crt"
	GRPCServerKeyFile  = CertsDir + "/mtls_server.key"
	GRPCClientCertFile = CertsDir + "/mtls_client.crt"
	GRPCClientKeyFile  = CertsDir + "/mtls_client.key"
	GRPCCaCertFile     = CertsDir + "/mtls_ca.crt"
)

// 获取环境变量
func init() {
	portStr := os.Getenv("GRPC_LISTEN_PORT")
	if portStr == "" {
		logger.Fatalln("gRPC服务监听端口未设置")
	}
	port, err := strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		logger.Fatalln(err)
	}
	if !utils.PortRangeCheck(port) {
		logger.Fatalln("gRPC服务监听端口范围不正确")
	}
	GRPC_LISTEN_PORT = int(port)

	portStr = os.Getenv("HTTP_LISTEN_PORT")
	if portStr == "" {
		logger.Fatalln("HTTP服务监听端口未设置")
	}
	port, err = strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		logger.Fatalln(err)
	}
	if port <= 0 || port > 65535 {
		logger.Fatalln("HTTP服务监听端口范围不正确")
	}
	HTTP_LISTEN_PORT = int(port)

	if mtls, err := strconv.ParseBool(os.Getenv("GRPC_mTLS")); err != nil {
		logger.Fatalln("gRPC mTLS状态未设置或设置不正确")
	} else {
		GRPC_mTLS = mtls
	}

	d := os.Getenv("DEBUG")
	if d != "" {
		if debug, err := strconv.ParseBool(d); err != nil {
			logger.Fatalln("Debug状态设置不正确")
		} else {
			DEBUG = debug
		}
	}

	for {
		gatewayHost := os.Getenv("GATEWAY_HOST")
		if gatewayHost == "" {
			logger.Fatalln("Gateway API Host为空")
		}
		if ip := net.ParseIP(gatewayHost); ip != nil { //为IP
			a, err := netip.ParseAddr(ip.String())
			if err != nil {
				logger.Fatalln(err)
			}
			GATEWAY_IP = a
			break
		} else if ok, err := regexp.MatchString(`^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$`, gatewayHost); ok { //为域名
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			ips, err := net.DefaultResolver.LookupIP(ctx, "ip", gatewayHost)
			if err != nil {
				logger.Fatalln(err)
			}
			if len(ips) == 0 {
				logger.Fatalln("Gateway API Host有误，无解析结果")
			}
			a, err := netip.ParseAddr(ips[0].String())
			if err != nil {
				logger.Fatalln(err)
			}
			GATEWAY_IP = a
			break
		} else { //若无错误，为未知
			if err != nil {
				logger.Fatalln(err)
			} else {
				logger.Fatalln("Gateway API Host有误，非域名或IP")
			}
		}
	}
	portStr = os.Getenv("GATEWAY_GRPC_PORT")
	if portStr == "" {
		logger.Fatalln("Gateway gRPC服务端口未设置")
	}
	port, err = strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		logger.Fatalln(err)
	}
	if !utils.PortRangeCheck(port) {
		logger.Fatalln("Gateway gRPC服务端口范围不正确")
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
		cert, err := tls.LoadX509KeyPair(GRPCClientCertFile, GRPCClientKeyFile)
		if err != nil {
			logger.Fatalln(err)
		}
		caPool := x509.NewCertPool()
		buf, err := os.ReadFile(GRPCCaCertFile)
		if err != nil {
			logger.Fatalln(err)
		}
		if !caPool.AppendCertsFromPEM(buf) {
			logger.Fatalln("添加CA证书失败")
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs:      caPool,
			Certificates: []tls.Certificate{cert},
			ServerName:   "mtls.susu",
		})
		client, err := zrpc.NewClient(c.RpcClientConf, zrpc.WithDialOption(grpc.WithTransportCredentials(cred)))
		if err != nil {
			logger.Fatalln(err)
		}
		Call_Connector = connectorclient.NewConnector(client)
		Call_Handler_Request = requesthandler.NewRequestHandler(client)
		Call_Handler_Response = responsehandler.NewResponseHandler(client)
	} else {
		var c types.Config
		client, err := zrpc.NewClient(c.RpcClientConf)
		if err != nil {
			logger.Fatalln(err)
		}
		Call_Connector = connectorclient.NewConnector(client)
		Call_Handler_Request = requesthandler.NewRequestHandler(client)
		Call_Handler_Response = responsehandler.NewResponseHandler(client)
	}
}

func GRPCOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	if GRPC_mTLS {
		cert, err := tls.LoadX509KeyPair(GRPCServerCertFile, GRPCServerKeyFile)
		if err != nil {
			logger.Fatalln(err)
		}
		caPool := x509.NewCertPool()
		buf, err := os.ReadFile(GRPCCaCertFile)
		if err != nil {
			logger.Fatalln(err)
		}
		if !caPool.AppendCertsFromPEM(buf) {
			logger.Fatalln("添加CA证书失败")
		}
		cred := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    caPool,
		})
		opts = append(opts, grpc.Creds(cred))
	}
	return opts
}

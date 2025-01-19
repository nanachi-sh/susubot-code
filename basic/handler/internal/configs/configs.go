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

	connectorclient "github.com/nanachi-sh/susubot-code/basic/handler/internal/caller/connector"
	filewebclient "github.com/nanachi-sh/susubot-code/basic/handler/internal/caller/fileweb"
	"github.com/nanachi-sh/susubot-code/basic/handler/internal/types"
	"github.com/nanachi-sh/susubot-code/basic/handler/pkg/utils"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	GRPC_LISTEN_PORT  int
	GRPC_mTLS         bool
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

	CertsDir           = ConfigDir + "/certs"
	GRPCServerCertFile = CertsDir + "/mtls_server.crt"
	GRPCServerKeyFile  = CertsDir + "/mtls_server.key"
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

	if mtls, err := strconv.ParseBool(os.Getenv("GRPC_mTLS")); err != nil {
		logger.Fatalln("gRPC mTLS状态未设置或设置不正确")
	} else {
		GRPC_mTLS = mtls
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

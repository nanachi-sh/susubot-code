package configs

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	GRPC_LISTEN_PORT int
	HTTP_LISTEN_PORT int
	GRPC_mTLS        bool

	SEED1, SEED2 uint64

	RPCServer_Config string = "rpc_server.yaml"
)

const (
	WebDir    = "/var/www/html"
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
	if port <= 0 || port > 65535 {
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
}

// 初始化gRPC配置
func init() {
	config := fmt.Sprintf(`Name: connector.rpc
ListenOn: 0.0.0.0:%d`, GRPC_LISTEN_PORT)
	if err := os.WriteFile(RPCServer_Config, []byte(config), 0744); err != nil {
		logger.Fatalln(err)
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

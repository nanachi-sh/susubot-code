package configs

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"regexp"
	"strconv"
	"time"

	qqverifierclient "github.com/nanachi-sh/susubot-code/basic/jwt/internal/caller/qqverifier"
	unomodel "github.com/nanachi-sh/susubot-code/basic/jwt/internal/model/uno"
	"github.com/nanachi-sh/susubot-code/basic/jwt/internal/types"
	"github.com/nanachi-sh/susubot-code/basic/jwt/pkg/utils"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	GRPC_LISTEN_PORT  int
	HTTP_LISTEN_PORT  int
	GRPC_mTLS         bool
	GATEWAY_IP        netip.Addr
	GATEWAY_GRPC_PORT int
	ASSETS_URL        string
	DATABASE_IP       netip.Addr
	DATABASE_PORT     int
	DATABASE_USER     string
	DATABASE_PASSWORD string

	RPCServer_Config  string = "rpc_server.yaml"
	RPCClient_Config  string = "rpc_client.yaml"
	RPCGateway_Config string = "rpc_gateway.yaml"

	Call_QQVerifier qqverifierclient.Qqverifier
	Model_Uno       unomodel.UsersModel
	DefaultCtx      context.Context = context.Background()

	JWTKey *ecdsa.PrivateKey
)

const (
	ConfigDir = "/config"

	CertsDir           = ConfigDir + "/certs"
	GRPCServerCertFile = CertsDir + "/mtls_server.crt"
	GRPCServerKeyFile  = CertsDir + "/mtls_server.key"
	GRPCClientCertFile = CertsDir + "/mtls_client.crt"
	GRPCClientKeyFile  = CertsDir + "/mtls_client.key"
	GRPCCaCertFile     = CertsDir + "/mtls_ca.crt"
	JWTKeyFile         = CertsDir + "/jwt.key"
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

	for {
		host := os.Getenv("DATABASE_HOST")
		if host == "" {
			logger.Fatalln("Database Host为空")
		}
		if ip := net.ParseIP(host); ip != nil { //为IP
			a, err := netip.ParseAddr(ip.String())
			if err != nil {
				logger.Fatalln(err)
			}
			DATABASE_IP = a
			break
		} else if ok, err := regexp.MatchString(`^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$`, host); ok { //为域名
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
			if err != nil {
				logger.Fatalln(err)
			}
			if len(ips) == 0 {
				logger.Fatalln("Database Host有误，无解析结果")
			}
			a, err := netip.ParseAddr(ips[0].String())
			if err != nil {
				logger.Fatalln(err)
			}
			DATABASE_IP = a
			break
		} else { //若无错误，为未知
			if err != nil {
				logger.Fatalln(err)
			} else {
				logger.Fatalln("Database Host有误，非域名或IP")
			}
		}
	}
	portStr = os.Getenv("DATABASE_PORT")
	if portStr == "" {
		logger.Fatalln("Database服务端口未设置")
	}
	port, err = strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		logger.Fatalln(err)
	}
	if !utils.PortRangeCheck(port) {
		logger.Fatalln("Database服务端口范围不正确")
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
	// cert, err := tls.LoadX509KeyPair(GRPCClientCertFile, GRPCClientKeyFile)
	// if err != nil {
	// 	logger.Fatalln(err)
	// }
	// caPool := x509.NewCertPool()
	// buf, err := os.ReadFile(GRPCCaCertFile)
	// if err != nil {
	// 	logger.Fatalln(err)
	// }
	// if !caPool.AppendCertsFromPEM(buf) {
	// 	logger.Fatalln("添加CA证书失败")
	// }
	// cred := credentials.NewTLS(&tls.Config{
	// 	RootCAs:      caPool,
	// 	Certificates: []tls.Certificate{cert},
	// 	ServerName:   "mtls.susu",
	// })
	// client, err := zrpc.NewClient(c.RpcClientConf, zrpc.WithDialOption(grpc.WithTransportCredentials(cred)))
	// if err != nil {
	// 	logger.Fatalln(err)
	// }
	client, err := zrpc.NewClient(c.RpcClientConf)
	if err != nil {
		logger.Fatalln(err)
	}
	Call_QQVerifier = qqverifierclient.NewQqverifier(client)
}

// 初始化SQL Models
func init() {
	conn := sqlx.NewMysql(fmt.Sprintf("%s:%s@tcp(%s:%d)/uno", DATABASE_USER, DATABASE_PASSWORD, DATABASE_IP, DATABASE_PORT))
	Model_Uno = unomodel.NewUsersModel(conn)
}

func init() {
	buf, err := os.ReadFile(JWTKeyFile)
	if err != nil {
		logger.Fatalln(err)
	}
	pblock, _ := pem.Decode(buf)
	if pblock == nil {
		logger.Fatalln("JWT私钥不正确")
	}
	pk, err := x509.ParseECPrivateKey(pblock.Bytes)
	if err != nil {
		logger.Fatalln(err)
	}
	JWTKey = pk
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

package configs

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strconv"

	qqverifierclient "github.com/nanachi-sh/susubot-code/basic/jwt/internal/caller/qqverifier"
	mock_qqverifierclient "github.com/nanachi-sh/susubot-code/basic/jwt/internal/mock/qqverifier"
	unomodel "github.com/nanachi-sh/susubot-code/basic/jwt/internal/model/uno"
	"github.com/nanachi-sh/susubot-code/basic/jwt/internal/types"
	"github.com/nanachi-sh/susubot-code/basic/jwt/internal/utils"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	DEBUG bool

	GRPC_LISTEN_PORT  int
	HTTP_LISTEN_PORT  int
	GATEWAY_IP        netip.Addr
	GATEWAY_GRPC_PORT int
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

	CertsDir   = ConfigDir + "/certs"
	JWTKeyFile = CertsDir + "/jwt.key"
)

// 获取环境变量
func init() {
	port, err := utils.EnvPortToPort("GRPC_LISTEN_PORT")
	if err != nil {
		logger.Fatalf("grpc监听端口获取失败")
	}
	GRPC_LISTEN_PORT = int(port)

	port, err = utils.EnvPortToPort("HTTP_LISTEN_PORT")
	if err != nil {
		logger.Fatalf("HTTP监听端口获取失败")
	}
	HTTP_LISTEN_PORT = int(port)

	if d := os.Getenv("DEBUG"); d != "" {
		if debug, err := strconv.ParseBool(d); err != nil {
			logger.Fatalln("Debug状态设置不正确")
		} else {
			DEBUG = debug
		}
	}

	if host := os.Getenv("GATEWAY_HOST"); host == "" {
		logger.Fatalln("Gateway Host未设置")
	} else {
		addr, err := utils.ResolvIP(host)
		if err != nil {
			logger.Fatalf("gateway host获取失败，err: %v", err)
		}
		GATEWAY_IP = addr
	}
	port, err = utils.EnvPortToPort(os.Getenv("GATEWAY_GRPC_PORT"))
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
		Call_QQVerifier = qqverifierclient.NewQqverifier(client)
	} else {
		Call_QQVerifier = mock_qqverifierclient.DefaultMock()
	}
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
	return opts
}

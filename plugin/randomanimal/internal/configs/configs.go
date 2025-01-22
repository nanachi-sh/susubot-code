package configs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strconv"

	filewebclient "github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/caller/fileweb"
	randomanimalmodel "github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/model/randomanimal"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/types"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/utils"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	DEBUG bool

	GRPC_LISTEN_PORT int

	GATEWAY_IP        netip.Addr
	GATEWAY_GRPC_PORT int
	ASSETS_URL        string

	DATABASE_IP       netip.Addr
	DATABASE_PORT     int
	DATABASE_USER     string
	DATABASE_PASSWORD string

	SEED1, SEED2 uint64

	Call_Fileweb filewebclient.FileWeb

	Model_Randomanimal randomanimalmodel.CachesModel

	DefaultCtx = context.Background()

	RPCServer_Config string = "rpc_server.json"
	RPCClient_Config string = "rpc_client.yaml"
)

const (
	ConfigDir = "/config"
)

const (
	CatAPI  = "https://api.thecatapi.com/v1/images/search"
	DogAPI  = "https://random.dog/woof.json"
	FoxAPI  = "https://randomfox.ca/floof/"
	DuckAPI = "https://random-d.uk/api/randomimg"
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

	if ASSETS_URL = os.Getenv("ASSETS_URL"); ASSETS_URL == "" {
		logger.Fatalln("assets url未设置")
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

	port, err = utils.EnvPortToPort("GATEWAY_GRPC_PORT")
	if err != nil {
		logger.Fatalf("gateway grpc port获取失败，err: %v", err)
	}
	GATEWAY_GRPC_PORT = int(port)

	if s, err := strconv.ParseUint(os.Getenv("SEED1"), 10, 0); err != nil {
		logger.Fatalf("seed未设置或设置有误")
	} else {
		SEED1 = s
	}
	if s, err := strconv.ParseUint(os.Getenv("SEED2"), 10, 0); err != nil {
		logger.Fatalf("seed未设置或设置有误")
	} else {
		SEED2 = s
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
}

// 初始化gRPC配置
func init() {
	configm := make(map[string]any)
	configm["Timeout"] = 10000
	configm["ListenOn"] = fmt.Sprintf("0.0.0.0:%d", GRPC_LISTEN_PORT)
	configm["Name"] = "randomanimal.rpc"
	configm["Log"] = struct {
		Mode string
		Path string
	}{
		"file",
		"logs",
	}
	configbs, err := json.Marshal(configm)
	if err != nil {
		logger.Fatalln(err)
	}
	if err := os.WriteFile(RPCServer_Config, configbs, 0744); err != nil {
		logger.Fatalln(err)
	}
	config := fmt.Sprintf(`Target: %s:%d`, GATEWAY_IP.String(), GATEWAY_GRPC_PORT)
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
		Call_Fileweb = filewebclient.NewFileWeb(client)
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
		Call_Fileweb = filewebclient.NewFileWeb(client)
	}
}

// 初始化SQL Models
func init() {
	conn := sqlx.NewMysql(fmt.Sprintf("%s:%s@tcp(%s:%d)/randomanimal", DATABASE_USER, DATABASE_PASSWORD, DATABASE_IP, DATABASE_PORT))
	Model_Randomanimal = randomanimalmodel.NewCachesModel(conn)
}

func GRPCOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	return opts
}

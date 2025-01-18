package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"

	"github.com/nanachi-sh/susubot-code/basic/connector/internal/configs"
	"github.com/nanachi-sh/susubot-code/basic/connector/pkg/protos/connector"
	"github.com/nanachi-sh/susubot-code/basic/connector/service/internal/config"
	"github.com/nanachi-sh/susubot-code/basic/connector/service/internal/server"
	"github.com/nanachi-sh/susubot-code/basic/connector/service/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", configs.RPC_Config, "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		connector.RegisterConnectorServer(grpcServer, server.NewConnectorServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	if configs.GRPC_mTLS {
		cert, err := tls.LoadX509KeyPair(configs.GRPCCertFile, configs.GRPCKeyFile)
		if err != nil {
			panic(err)
		}
		caPool := x509.NewCertPool()
		buf, err := os.ReadFile(configs.GRPCCaFile)
		if err != nil {
			panic(err)
		}
		ca, _ := pem.Decode(buf)
		if ca == nil {
			panic(err)
		}
		if !caPool.AppendCertsFromPEM(ca.Bytes) {
			panic("")
		}
		cred := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    caPool,
		})
		s.AddOptions(grpc.Creds(cred))
	}
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

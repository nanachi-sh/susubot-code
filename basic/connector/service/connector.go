package main

import (
	"flag"
	"fmt"

	"github.com/nanachi-sh/susubot-code/basic/connector/internal/configs"
	"github.com/nanachi-sh/susubot-code/basic/connector/pkg/protos/connector"
	"github.com/nanachi-sh/susubot-code/basic/connector/service/internal/config"
	"github.com/nanachi-sh/susubot-code/basic/connector/service/internal/server"
	"github.com/nanachi-sh/susubot-code/basic/connector/service/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "rpc.yaml", "the config file")

func main() {
	_ = configs.LOAD
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

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

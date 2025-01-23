package main

import (
	"flag"
	"fmt"

	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/internal/configs"
	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/pkg/protos/randomfortune"
	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/service/internal/config"
	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/service/internal/server"
	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/service/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", configs.RPCServer_Config, "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		randomfortune.RegisterRandomFortuneServer(grpcServer, server.NewRandomFortuneServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddOptions(configs.GRPCOptions()...)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

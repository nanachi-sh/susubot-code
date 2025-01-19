package main

import (
	"flag"
	"fmt"

	"github.com/nanachi-sh/susubot-code/basic/jwt/gateway"
	"github.com/nanachi-sh/susubot-code/basic/jwt/internal/configs"
	"github.com/nanachi-sh/susubot-code/basic/jwt/pkg/protos/jwt"
	"github.com/nanachi-sh/susubot-code/basic/jwt/service/internal/config"
	"github.com/nanachi-sh/susubot-code/basic/jwt/service/internal/server"
	"github.com/nanachi-sh/susubot-code/basic/jwt/service/internal/svc"

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
		jwt.RegisterJwtServer(grpcServer, server.NewJwtServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddOptions(configs.GRPCOptions()...)
	go func() {
		gateway.Serve()
	}()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

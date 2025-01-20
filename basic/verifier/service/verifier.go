package main

import (
	"flag"
	"fmt"

	"github.com/nanachi-sh/susubot-code/basic/verifier/gateway"
	"github.com/nanachi-sh/susubot-code/basic/verifier/internal/configs"
	"github.com/nanachi-sh/susubot-code/basic/verifier/pkg/protos/verifier"
	"github.com/nanachi-sh/susubot-code/basic/verifier/service/internal/config"
	"github.com/nanachi-sh/susubot-code/basic/verifier/service/internal/server"
	"github.com/nanachi-sh/susubot-code/basic/verifier/service/internal/svc"

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
		verifier.RegisterVerifierServer(grpcServer, server.NewVerifierServer(ctx))

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

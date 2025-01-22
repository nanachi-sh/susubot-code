package main

import (
	"fmt"

	"github.com/nanachi-sh/susubot-code/basic/fileweb/internal/configs"
	"github.com/nanachi-sh/susubot-code/basic/fileweb/internal/web"
	"github.com/nanachi-sh/susubot-code/basic/fileweb/pkg/protos/fileweb"
	"github.com/nanachi-sh/susubot-code/basic/fileweb/service/internal/config"
	"github.com/nanachi-sh/susubot-code/basic/fileweb/service/internal/server"
	"github.com/nanachi-sh/susubot-code/basic/fileweb/service/internal/svc"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var c config.Config
	c.RpcServerConf = configs.RPCServer_Config
	ctx := svc.NewServiceContext(c)

	go func() { web.Serve() }()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		fileweb.RegisterFileWebServer(grpcServer, server.NewFileWebServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddOptions(configs.GRPCOptions()...)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

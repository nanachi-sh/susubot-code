package main

import (
	"flag"
	"fmt"

	"github.com/nanachi-sh/susubot-code/basic/handler/internal/configs"
	"github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/handler/request"
	"github.com/nanachi-sh/susubot-code/basic/handler/pkg/protos/handler/response"
	"github.com/nanachi-sh/susubot-code/basic/handler/service/internal/config"
	requesthandlerServer "github.com/nanachi-sh/susubot-code/basic/handler/service/internal/server/requesthandler"
	responsehandlerServer "github.com/nanachi-sh/susubot-code/basic/handler/service/internal/server/responsehandler"
	"github.com/nanachi-sh/susubot-code/basic/handler/service/internal/svc"

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
		request.RegisterRequestHandlerServer(grpcServer, requesthandlerServer.NewRequestHandlerServer(ctx))
		response.RegisterResponseHandlerServer(grpcServer, responsehandlerServer.NewResponseHandlerServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddOptions(configs.GRPCOptions()...)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

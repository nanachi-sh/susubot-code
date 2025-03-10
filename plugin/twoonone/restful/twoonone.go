package main

import (
	"flag"
	"fmt"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/configs"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/config"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/handler"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", configs.APIServer_Config, "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors("https://twoonone.unturned.fun:8080", "https://twoonone.susu.love", "https://twoonone.unturned.fun"))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

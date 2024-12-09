package define

import (
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/LLOneBot/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	GatewayIP  net.IP
	GRPCClient *grpc.ClientConn
	FilewebCtx context.Context

	logger = log.Get()
)

func init() {
	logger.Println("init run")
	gatewayHost := os.Getenv("GATEWAY_HOST")
	if gatewayHost == "" {
		logger.Fatalln("Gateway API Host为空")
	}
	if ip := net.ParseIP(gatewayHost); ip != nil { //为IP
		GatewayIP = ip
		return
	} else if ok, err := regexp.MatchString(`^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$`, gatewayHost); ok { //为域名
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		ips, err := net.DefaultResolver.LookupIP(ctx, "ip", gatewayHost)
		if err != nil {
			logger.Fatalln(err)
			return
		}
		if len(ips) == 0 {
			logger.Fatalln("Gateway API Host有误，无解析结果")
			return
		}
		GatewayIP = ips[0]
		return
	} else { //若无错误，为未知
		if err != nil {
			logger.Fatalln(err)
		} else {
			logger.Fatalln("Gateway API Host有误，非域名或IP")
		}
	}
	fmt.Println(GatewayIP.String())
	c, err := grpc.NewClient(fmt.Sprintf("%v:2080", GatewayIP.String()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalln(err)
	}
	fmt.Println(c == nil)
	GRPCClient = c
	FilewebCtx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"service-target": "fileweb",
	}))
}

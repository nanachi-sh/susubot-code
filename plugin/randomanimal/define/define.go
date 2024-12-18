package define

import (
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	GatewayIP  net.IP
	GRPCClient *grpc.ClientConn
	FilewebCtx context.Context

	ExternalHost     string
	ExternalHTTPPort int

	logger = log.Get()
)

func init() {
	gatewayHost := os.Getenv("GATEWAY_HOST")
	if gatewayHost == "" {
		logger.Fatalln("Gateway API Host为空")
	}
	for {
		if ip := net.ParseIP(gatewayHost); ip != nil { //为IP
			GatewayIP = ip
			break
		} else if ok, err := regexp.MatchString(`^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$`, gatewayHost); ok { //为域名
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			ips, err := net.DefaultResolver.LookupIP(ctx, "ip", gatewayHost)
			if err != nil {
				logger.Fatalln(err)
			}
			if len(ips) == 0 {
				logger.Fatalln("Gateway API Host有误，无解析结果")
			}
			GatewayIP = ips[0]
			break
		} else { //若无错误，为未知
			if err != nil {
				logger.Fatalln(err)
			} else {
				logger.Fatalln("Gateway API Host有误，非域名或IP")
			}
		}
	}
	c, err := grpc.NewClient(fmt.Sprintf("%v:2080", GatewayIP.String()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalln(err)
	}
	GRPCClient = c
	FilewebCtx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"service-target": "fileweb",
	}))
	ExternalHost = os.Getenv("EXTERNAL_HOST")
	if ExternalHost == "" {
		logger.Fatalln("External Host未设置")
	}
	ExternalHTTPPort_str := os.Getenv("EXTERNAL_HTTP_PORT")
	if ExternalHTTPPort_str == "" {
		logger.Fatalln("External HTTPPort未设置")
	}
	httpport, err := strconv.ParseInt(ExternalHTTPPort_str, 10, 0)
	if err != nil {
		logger.Fatalln("External HTTPPort不为纯整数")
	}
	if httpport <= 0 || httpport > 65535 {
		logger.Fatalln("External HTTPPort范围不正确")
	}
	ExternalHTTPPort = int(httpport)
}

package define

import (
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/nanachi-sh/susubot-code/basic/handler/log"
	connector_pb "github.com/nanachi-sh/susubot-code/basic/handler/protos/connector"
	fileweb_pb "github.com/nanachi-sh/susubot-code/basic/handler/protos/fileweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	GatewayIP  net.IP
	GRPCClient *grpc.ClientConn
	ConnectorC connector_pb.ConnectorClient
	FilewebC   fileweb_pb.FileWebClient
	DefaultCtx context.Context

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
	ConnectorC = connector_pb.NewConnectorClient(GRPCClient)
	FilewebC = fileweb_pb.NewFileWebClient(GRPCClient)
}

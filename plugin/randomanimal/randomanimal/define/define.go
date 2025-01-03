package define

import (
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/log"
	fileweb_pb "github.com/nanachi-sh/susubot-code/plugin/randomanimal/protos/fileweb"
	randomanimal_pb "github.com/nanachi-sh/susubot-code/plugin/randomanimal/protos/randomanimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type BasicReturn struct {
	Buf  []byte
	Hash *string
	Type randomanimal_pb.Type
}

var (
	GatewayIP  net.IP
	FilewebC   fileweb_pb.FileWebClient
	FilewebCtx context.Context
	GRPCClient *grpc.ClientConn

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
	FilewebC = fileweb_pb.NewFileWebClient(GRPCClient)
	FilewebCtx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"service-target": "fileweb",
		"version":        "stable",
	}))
}

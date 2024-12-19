package define

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/nanachi-sh/susubot-code/basic/qqinteraction/log"
	connector_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/connector"
	request_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/handler/request"
	response_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/handler/response"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	Conf = new(Config)

	GatewayIP         net.IP
	GRPCClient        *grpc.ClientConn
	ConnectorCtx      context.Context
	ConnectorC        connector_pb.ConnectorClient
	HandlerCtx        context.Context
	Handler_RequestC  request_pb.RequestHandlerClient
	Handler_ResponseC response_pb.ResponseHandlerClient

	logger = log.Get()
)

func init() {
	d, err := os.ReadFile("/config/qqinteraction_config.json")
	if err != nil {
		logger.Fatalln(err)
	}
	if err := json.Unmarshal(d, Conf); err != nil {
		logger.Fatalln(err)
	}

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
	ConnectorCtx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"service-target": "connector",
	}))
	ConnectorC = connector_pb.NewConnectorClient(GRPCClient)
	HandlerCtx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"service-target": "handler",
	}))
	Handler_RequestC = request_pb.NewRequestHandlerClient(GRPCClient)
	Handler_ResponseC = response_pb.NewResponseHandlerClient(GRPCClient)
}

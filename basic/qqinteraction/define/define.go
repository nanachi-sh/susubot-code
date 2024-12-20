package define

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/nanachi-sh/susubot-code/basic/qqinteraction/log"
	connector_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/connector"
	request_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/handler/request"
	response_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/handler/response"
	randomanimal_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/randomanimal"
	randomfortune_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/randomfortune"
	twoonone_pb "github.com/nanachi-sh/susubot-code/basic/qqinteraction/protos/twoonone"
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
	RandomAnimalC     randomanimal_pb.RandomAnimalClient
	RandomAnimalCtx   context.Context
	RandomFortuneC    randomfortune_pb.RandomFortuneClient
	RandomFortuneCtx  context.Context
	TwoOnOneC         twoonone_pb.TwoOnOneClient
	TwoOnOneCtx       context.Context

	ExternalHost     string
	ExternalHTTPPort int
	ExternalURL      string

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
		"version":        "stable",
	}))
	ConnectorC = connector_pb.NewConnectorClient(GRPCClient)
	HandlerCtx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"service-target": "handler",
		"version":        "stable",
	}))
	Handler_RequestC = request_pb.NewRequestHandlerClient(GRPCClient)
	Handler_ResponseC = response_pb.NewResponseHandlerClient(GRPCClient)
	RandomAnimalC = randomanimal_pb.NewRandomAnimalClient(GRPCClient)
	RandomAnimalCtx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"service-target": "randomanimal",
		"version":        "stable",
	}))
	RandomFortuneC = randomfortune_pb.NewRandomFortuneClient(GRPCClient)
	RandomFortuneCtx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"service-target": "randomfortune",
		"version":        "stable",
	}))
	TwoOnOneC = twoonone_pb.NewTwoOnOneClient(GRPCClient)
	TwoOnOneCtx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"service-target": "twoonone",
		"version":        "stable",
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
	ExternalURL = fmt.Sprintf("http://%v:%v", ExternalHost, ExternalHTTPPort)
}

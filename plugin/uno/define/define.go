package define

import (
	"context"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/nanachi-sh/susubot-code/plugin/uno/log"
	qqverifier_pb "github.com/nanachi-sh/susubot-code/plugin/uno/protos/qqverifier"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	AssetsDir = "/config"
	CertsDir  = AssetsDir + "/certs"
)

var (
	PrivilegeUserHash string
	EnableTLS         bool
	GatewayIP         net.IP
	GRPCClient        *grpc.ClientConn
	QQVerifierC       qqverifier_pb.QqverifierClient
	QQVerifierCtx     context.Context

	logger = log.Get()
)

func init() {
	PrivilegeUserHash = os.Getenv("PRIVILEGE_USERHASH")
	if PrivilegeUserHash == "" {
		logger.Fatalln("Privilege UserHash为空")
	}
	gatewayHost := os.Getenv("GATEWAY_HOST")
	if gatewayHost == "" {
		logger.Fatalln("Gateway API Host为空")
	}
	if b, err := strconv.ParseBool(os.Getenv("ENABLE_TLS")); err != nil {
		logger.Fatalln("Enable TLS未设置或设置有误")
	} else {
		EnableTLS = b
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
	QQVerifierC = qqverifier_pb.NewQqverifierClient(GRPCClient)
	QQVerifierCtx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		"service-target": "qqverifier",
		"version":        "stable",
	}))
}

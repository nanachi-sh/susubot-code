package define

import (
	"context"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/nanachi-sh/susubot-code/basic/handler/LLOneBot/log"
)

var (
	GatewayIP net.IP
	logger    = log.Get()
)

func init() {
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
}

const (
	Role_Member = "member"
	Role_Admin  = "admin"
	Role_Owner  = "owner"
)

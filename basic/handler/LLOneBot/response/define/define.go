package define

import (
	"net"
	"os"
)

var (
	GatewayIP = net.ParseIP(os.Getenv("GATEWAY_IP"))
)

func init() {
	if GatewayIP == nil {
		panic("Gateway API IP未设置")
	}
}

const (
	Role_Member = "member"
	Role_Admin  = "admin"
	Role_Owner  = "owner"
)

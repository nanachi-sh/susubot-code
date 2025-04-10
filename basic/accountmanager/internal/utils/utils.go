package utils

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"net/netip"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func PortRangeCheck(p int64) bool { return !(p <= 0 || p > 65535) }

const (
	Dict_Mixed   string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Dict_Number  string = "0123456789"
	Dict_English string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY"
)

func RandomString(length int, dict string) string {
	ret := new(strings.Builder)
	for i := 0; i < length; i++ {
		ret.WriteByte(dict[rand.Intn(len(dict))])
	}
	return ret.String()
}

func ResolvIP(addr string) (netip.Addr, error) {
	if ip := net.ParseIP(addr); ip != nil { //为IP
		return netip.ParseAddr(ip.String())
	} else if ok, err := regexp.MatchString(`^(([a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62}))+|localhost)$`, addr); ok { //为域名或localhost
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		ips, err := net.DefaultResolver.LookupIP(ctx, "ip", addr)
		if err != nil {
			return netip.Addr{}, err
		}
		if len(ips) == 0 {
			return netip.Addr{}, errors.New("无解析结果")
		}
		return netip.ParseAddr(ips[0].String())
	} else { //若无错误，为未知
		if err != nil {
			return netip.Addr{}, err
		} else {
			return netip.Addr{}, errors.New("非域名或IP")
		}
	}
}

func EnvPortToPort(envKey string) (uint16, error) {
	portStr, err := EnvToString(envKey)
	if err != nil {
		return 0, err
	}
	port, err := strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		return 0, err
	}
	if !PortRangeCheck(port) {
		return 0, errors.New("端口范围不正确")
	}
	return uint16(port), nil
}

func EnvToString(envKey string) (string, error) {
	if str := os.Getenv(envKey); str == "" {
		key_desc := ""
		{
			strs := strings.Split(envKey, "_")
			for _, v := range strs {
				key_desc += v
				key_desc += " "
			}
		}
		return "", errors.New(key_desc + " 未设置")
	} else {
		return str, nil
	}
}

package utils

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/netip"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/pkg/protos/randomanimal"
	"github.com/twmb/murmur3"
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
	portStr := os.Getenv(envKey)
	if portStr == "" {
		return 0, errors.New("未设置")
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

func Murmurhash128ToString(seed1, seed2 uint64, text string) string {
	h1, h2 := murmur3.SeedStringSum128(seed1, seed2, text)
	builder := new(strings.Builder)
	builder.WriteString(strconv.FormatUint(h1, 16))
	builder.WriteString(strconv.FormatUint(h2, 16))
	return builder.String()
}

func MatchMediaType(resp *http.Response) (randomanimal.Type, error) {
	if resp == nil || resp.Close {
		return 0, errors.New("response错误")
	}
	switch ct := strings.ToLower(resp.Header.Get("Content-Type")); ct {
	case "image/webp", "image/jpeg", "image/png", "image/gif":
		return randomanimal.Type_Image, nil
	case "video/mpeg4", "video/mp4", "video/webm":
		return randomanimal.Type_Video, nil
	default:
		return 0, fmt.Errorf("判断响应类型失败，Content-Type: %v", ct)
	}
}

package api

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/nanachi-sh/susubot-code/basic/fileweb/fileweb/define"
)

func HTTPServe() error {
	portStr := os.Getenv("HTTP_LISTEN_PORT")
	if portStr == "" {
		return errors.New("HTTP服务监听端口未设置")
	}
	port, err := strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		return err
	}
	if port <= 0 || port > 65535 {
		return errors.New("HTTP服务监听端口范围不正确")
	}
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		return errors.New("HTTP服务监听端口已被占用")
	}
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(define.WorkDir))))
	return http.Serve(l, nil)
}

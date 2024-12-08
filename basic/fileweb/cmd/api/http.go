package api

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
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
	if err := mkdir(); err != nil {
		return err
	}
	http.Handle("/", http.FileServer(http.Dir("/tmp")))
	return http.Serve(l, nil)
}

func mkdir() error {
	if _, err := os.Lstat("/var/www/html"); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll("/var/www/html", 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

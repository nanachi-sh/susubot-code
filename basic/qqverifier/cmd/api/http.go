package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	qqverifier_pb "github.com/nanachi-sh/susubot-code/basic/qqverifier/protos/qqverifier"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func HTTPServe() error {
	portStr := os.Getenv("HTTPAPI_LISTEN_PORT")
	if portStr == "" {
		return errors.New("HTTP API服务监听端口未设置")
	}
	port, err := strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		return err
	}
	if port <= 0 || port > 65535 {
		return errors.New("HTTP API服务监听端口范围不正确")
	}
	gRPCportStr := os.Getenv("GRPC_LISTEN_PORT")
	if portStr == "" {
		return errors.New("gRPC服务监听端口未设置")
	}
	gRPCport, err := strconv.ParseInt(gRPCportStr, 10, 0)
	if err != nil {
		return err
	}
	if gRPCport <= 0 || gRPCport > 65535 {
		return errors.New("gRPC服务监听端口范围不正确")
	}
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%v", gRPCport), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	sMux := runtime.NewServeMux()
	if err := qqverifier_pb.RegisterQqverifierHandler(context.Background(), sMux, conn); err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		return err
	}
	return http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") != "" {
			w.Header().Add("Access-Control-Allow-Origin", "https://192.168.1.254:8080")
			w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, OPTIONS, DELETE")
			w.Header().Add("Access-Control-Allow-Credentials", "true")
		}
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		sMux.ServeHTTP(w, r)
	}))
}

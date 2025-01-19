package gateway

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nanachi-sh/susubot-code/basic/jwt/internal/configs"
	"github.com/nanachi-sh/susubot-code/basic/jwt/pkg/protos/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

func Serve() {
	var conn *grpc.ClientConn
	if configs.GRPC_mTLS {
		cert, err := tls.LoadX509KeyPair(configs.GRPCClientCertFile, configs.GRPCClientKeyFile)
		if err != nil {
			logger.Fatalln(err)
		}
		caPool := x509.NewCertPool()
		buf, err := os.ReadFile(configs.GRPCCaCertFile)
		if err != nil {
			logger.Fatalln(err)
		}
		if !caPool.AppendCertsFromPEM(buf) {
			logger.Fatalln("添加CA证书失败")
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs:      caPool,
			Certificates: []tls.Certificate{cert},
			ServerName:   "mtls.susu",
		})
		c, err := grpc.NewClient(fmt.Sprintf("localhost:%v", configs.GRPC_LISTEN_PORT), grpc.WithTransportCredentials(cred))
		if err != nil {
			logger.Fatalln(err)
		}
		conn = c
	} else {
		c, err := grpc.NewClient(fmt.Sprintf("localhost:%v", configs.GRPC_LISTEN_PORT), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatalln(err)
		}
		conn = c
	}
	sMux := runtime.NewServeMux()
	if err := jwt.RegisterJwtHandler(context.Background(), sMux, conn); err != nil {
		logger.Fatalln(err)
	}
	logger.Fatalln(
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", configs.HTTP_LISTEN_PORT), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Origin") != "" {
				w.Header().Add("Access-Control-Allow-Origin", "*")
				w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, OPTIONS, DELETE")
			}
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			sMux.ServeHTTP(w, r)
		})),
	)
}

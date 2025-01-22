package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nanachi-sh/susubot-code/basic/jwt/internal/configs"
	"github.com/nanachi-sh/susubot-code/basic/jwt/pkg/protos/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

func Serve() {
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%v", configs.GRPC_LISTEN_PORT), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalln(err)
	}
	sMux := runtime.NewServeMux()
	if err := jwt.RegisterJwtHandler(context.Background(), sMux, conn); err != nil {
		logger.Fatalln(err)
	}
	addr := fmt.Sprintf("0.0.0.0:%d", configs.HTTP_LISTEN_PORT)
	fmt.Printf("Starting grpc gateway at %s...\n", addr)
	logger.Fatalln(
		http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

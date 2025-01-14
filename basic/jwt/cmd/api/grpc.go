package api

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/nanachi-sh/susubot-code/basic/jwt/define"
	"github.com/nanachi-sh/susubot-code/basic/jwt/jwt"
	jwt_pb "github.com/nanachi-sh/susubot-code/basic/jwt/protos/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type jwtService struct{ jwt_pb.JwtServer }

func GRPCServe() error {
	portStr := os.Getenv("GRPC_LISTEN_PORT")
	if portStr == "" {
		return errors.New("gRPC服务监听端口未设置")
	}
	port, err := strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		return err
	}
	if port <= 0 || port > 65535 {
		return errors.New("gRPC服务监听端口范围不正确")
	}
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		return err
	}
	opts := []grpc.ServerOption{}
	if define.EnableTLS {
		cert, err := tls.LoadX509KeyPair(fmt.Sprintf("%v/tls.pem", define.CertDir), fmt.Sprintf("%v/tls.key", define.CertDir))
		if err != nil {
			return err
		}
		cred := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
		})
		opts = append(opts, grpc.Creds(cred))
	}
	gs := grpc.NewServer(opts...)
	jwt_pb.RegisterJwtServer(gs, new(jwtService))
	return gs.Serve(l)
}

func (*jwtService) Uno_Sign(ctx context.Context, req *jwt_pb.Uno_SignRequest) (*jwt_pb.Uno_SignResponse, error) {
	type ret struct {
		data *jwt_pb.Uno_SignResponse
		err  error
	}
	ch := make(chan *ret, 1)
	go func() {
		ret := &ret{}
		defer func() { ch <- ret }()
		resp := jwt.Uno_Sign(req)
		ret.data = resp
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

func (*jwtService) Uno_Register(ctx context.Context, req *jwt_pb.Uno_RegisterRequest) (*jwt_pb.Uno_RegisterResponse, error) {
	type ret struct {
		data *jwt_pb.Uno_RegisterResponse
		err  error
	}
	ch := make(chan *ret, 1)
	go func() {
		ret := &ret{}
		defer func() { ch <- ret }()
		resp := jwt.Uno_Register(req)
		ret.data = resp
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return x.data, nil
	}
}

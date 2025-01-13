package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/nanachi-sh/susubot-code/basic/database/database"
	database_pb "github.com/nanachi-sh/susubot-code/basic/database/protos/database"
	"google.golang.org/grpc"
)

type databaseService struct{ database_pb.DatabaseServer }

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
	gs := grpc.NewServer()
	database_pb.RegisterDatabaseServer(gs, new(databaseService))
	return gs.Serve(l)
}

func (*databaseService) Uno_CreateUser(ctx context.Context, req *database_pb.Uno_CreateUserRequest) (*database_pb.Uno_CreateUserResponse, error) {
	type ret struct {
		data *database_pb.Uno_CreateUserResponse
		err  error
	}
	ch := make(chan *ret, 1)
	go func() {
		ret := &ret{}
		defer func() { ch <- ret }()
		resp := database.Uno_CreateUser(req)
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

func (*databaseService) Uno_GetUser(ctx context.Context, req *database_pb.Uno_GetUserRequest) (*database_pb.Uno_GetUserResponse, error) {
	type ret struct {
		data *database_pb.Uno_GetUserResponse
		err  error
	}
	ch := make(chan *ret, 1)
	go func() {
		ret := &ret{}
		defer func() { ch <- ret }()
		resp := database.Uno_GetUser(req)
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

func (*databaseService) Uno_VerifyUser(ctx context.Context, req *database_pb.Uno_VerifyUserRequest) (*database_pb.Uno_VerifyUserResponse, error) {
	type ret struct {
		data *database_pb.Uno_VerifyUserResponse
		err  error
	}
	ch := make(chan *ret, 1)
	go func() {
		ret := &ret{}
		defer func() { ch <- ret }()
		resp := database.Uno_VerifyUser(req)
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

func (*databaseService) Uno_UpdateUser(ctx context.Context, req *database_pb.Uno_UpdateUserRequest) (*database_pb.Uno_UpdateUserResponse, error) {
	type ret struct {
		data *database_pb.Uno_UpdateUserResponse
		err  error
	}
	ch := make(chan *ret, 1)
	go func() {
		ret := &ret{}
		defer func() { ch <- ret }()
		resp := database.Uno_UpdateUser(req)
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

func (*databaseService) Uno_ChangePassword(ctx context.Context, req *database_pb.Uno_ChangePasswordRequest) (*database_pb.Uno_ChangePasswordResponse, error) {
	type ret struct {
		data *database_pb.Uno_ChangePasswordResponse
		err  error
	}
	ch := make(chan *ret, 1)
	go func() {
		ret := &ret{}
		defer func() { ch <- ret }()
		resp := database.Uno_ChangePassword(req)
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

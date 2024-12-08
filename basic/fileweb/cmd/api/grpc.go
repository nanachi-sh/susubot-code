package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/nanachi-sh/susubot-code/basic/fileweb/fileweb"
	fileweb_pb "github.com/nanachi-sh/susubot-code/basic/fileweb/protos/fileweb"
	"google.golang.org/grpc"
)

type filewebService struct {
	fileweb_pb.FileWebServer
}

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
	fileweb_pb.RegisterFileWebServer(gs, new(filewebService))
	return gs.Serve(l)
}

func (*filewebService) Upload(ctx context.Context, req *fileweb_pb.UploadRequest) (*fileweb_pb.UploadResponse, error) {
	type ret struct {
		data *fileweb_pb.UploadResponse
		err  error
	}
	ch := make(chan *ret, 1)
	go func() {
		ret := &ret{}
		defer func() { ch <- ret }()
		resp, err := fileweb.Upload(req)
		if err != nil {
			ret.err = err
			return
		}
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

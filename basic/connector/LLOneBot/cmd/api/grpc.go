package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/nanachi-sh/susubot-code/basic/connector/LLOneBot/connector"
	connector_pb "github.com/nanachi-sh/susubot-code/basic/connector/LLOneBot/protos/connector"
	"google.golang.org/grpc"
)

type (
	connectorService struct {
		connectorH *connector.Connector
		connector_pb.ConnectorServer
	}
)

type myerror struct{}

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
	connector_pb.RegisterConnectorServer(gs, &connectorService{
		connectorH: connector.New(),
	})
	return gs.Serve(l)
}

func (cs *connectorService) Connect(ctx context.Context, req *connector_pb.ConnectRequest) (*connector_pb.ConnectResponse, error) {
	type ret struct {
		data *connector_pb.ConnectResponse
		err  error
	}
	ch := make(chan *ret, 1)
	go func() {
		ret := &ret{}
		defer func() { ch <- ret }()
		buf, err := cs.connectorH.Connect(req)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &connector_pb.ConnectResponse{
			Buf: buf,
		}
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

func (cs *connectorService) Read(_ *connector_pb.Empty, stream grpc.ServerStreamingServer[connector_pb.ReadResponse]) error {
	now := time.Now().UnixMicro()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan []byte)
	send := make(chan *connector_pb.ReadResponse, 1)
	go func() {
		for {
			fmt.Println(now, "send waiting")
			resp := <-send
			fmt.Println(now, "sending")
			if err := stream.Send(resp); err != nil {
				ctx = context.WithValue(ctx, myerror{}, err)
				cancel()
				return
			}
			fmt.Println(now, "sended")
		}
	}()
	for {
		go func() {
			buf, err := cs.connectorH.Read(now, 0)
			if err != nil {
				ctx = context.WithValue(ctx, myerror{}, err)
				cancel()
				return
			}
			ch <- buf
		}()
		select {
		case <-ctx.Done(): //结束
			switch x := ctx.Value(myerror{}).(type) {
			case error:
				return x
			default:
				return nil
			}
		case buf := <-ch: //新响应
			send <- &connector_pb.ReadResponse{
				IsClose: false,
				Buf:     buf,
			}
		}
	}
}

func (cs *connectorService) Write(ctx context.Context, req *connector_pb.WriteRequest) (*connector_pb.Empty, error) {
	type ret struct {
		err error
	}
	ch := make(chan *ret, 1)
	go func() {
		ret := &ret{}
		defer func() { ch <- ret }()
		if err := cs.connectorH.Write(req.Buf); err != nil {
			ret.err = err
			return
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return &connector_pb.Empty{}, nil
	}
}

func (cs *connectorService) Close(ctx context.Context, _ *connector_pb.Empty) (*connector_pb.Empty, error) {
	type ret struct {
		err error
	}
	ch := make(chan *ret, 1)
	go func() {
		ret := &ret{}
		defer func() { ch <- ret }()
		if err := cs.connectorH.Close(); err != nil {
			ret.err = err
			return
		}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case x := <-ch:
		if x.err != nil {
			return nil, x.err
		}
		return &connector_pb.Empty{}, nil
	}
}

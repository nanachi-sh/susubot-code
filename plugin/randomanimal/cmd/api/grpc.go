package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/define"
	randomanimal_pb "github.com/nanachi-sh/susubot-code/plugin/randomanimal/protos/randomanimal"
	randomanimal "github.com/nanachi-sh/susubot-code/plugin/randomanimal/randomanimal"
	"google.golang.org/grpc"
)

type (
	randomanimalService struct {
		randomanimal_pb.RandomAnimalServer
	}
)

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
	randomanimal_pb.RegisterRandomAnimalServer(gs, new(randomanimalService))
	return gs.Serve(l)
}

func (*randomanimalService) GetCat(ctx context.Context, req *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	type d struct {
		data *randomanimal_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		cat, err := randomanimal.GetCat(req.AutoUpload)
		if err != nil {
			ret.err = err
			return
		}
		if req.AutoUpload {
			ret.data = &randomanimal_pb.BasicResponse{
				Type: cat.Type,
				Response: &randomanimal_pb.BasicResponse_UploadResponse{
					Hash:    *cat.Hash,
					URLPath: fmt.Sprintf("/assets/%v", *cat.Hash),
				},
			}
		} else {
			ret.data = &randomanimal_pb.BasicResponse{
				Type: cat.Type,
				Buf:  cat.Buf,
			}
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

func (*randomanimalService) GetDog(ctx context.Context, req *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	type d struct {
		data *randomanimal_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		dog, err := randomanimal.GetDog(req.AutoUpload)
		if err != nil {
			ret.err = err
			return
		}
		if req.AutoUpload {
			ret.data = &randomanimal_pb.BasicResponse{
				Type: dog.Type,
				Response: &randomanimal_pb.BasicResponse_UploadResponse{
					Hash:    *dog.Hash,
					URLPath: fmt.Sprintf("/assets/%v", *dog.Hash),
				},
			}
		} else {
			ret.data = &randomanimal_pb.BasicResponse{
				Type: dog.Type,
				Buf:  dog.Buf,
			}
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

func (*randomanimalService) GetFox(ctx context.Context, req *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	type d struct {
		data *randomanimal_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		fox, err := randomanimal.GetFox(req.AutoUpload)
		if err != nil {
			ret.err = err
			return
		}
		if req.AutoUpload {
			ret.data = &randomanimal_pb.BasicResponse{
				Type: fox.Type,
				Response: &randomanimal_pb.BasicResponse_UploadResponse{
					Hash:    *fox.Hash,
					URLPath: fmt.Sprintf("/assets/%v", *fox.Hash),
				},
			}
		} else {
			ret.data = &randomanimal_pb.BasicResponse{
				Type: fox.Type,
				Buf:  fox.Buf,
			}
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

func (*randomanimalService) GetDuck(ctx context.Context, req *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	type d struct {
		data *randomanimal_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		duck, err := randomanimal.GetDuck(req.AutoUpload)
		if err != nil {
			ret.err = err
			return
		}
		if req.AutoUpload {
			ret.data = &randomanimal_pb.BasicResponse{
				Type: duck.Type,
				Response: &randomanimal_pb.BasicResponse_UploadResponse{
					Hash:    *duck.Hash,
					URLPath: fmt.Sprintf("/assets/%v", *duck.Hash),
				},
			}
		} else {
			ret.data = &randomanimal_pb.BasicResponse{
				Type: duck.Type,
				Buf:  duck.Buf,
			}
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

func (*randomanimalService) GetChiken_CXK(ctx context.Context, req *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	type d struct {
		data *randomanimal_pb.BasicResponse
		err  error
	}
	ch := make(chan *d, 1)
	go func() {
		ret := new(d)
		defer func() { ch <- ret }()
		chickenHash, err := randomanimal.GetChicken_CXK()
		if err != nil {
			ret.err = err
			return
		}
		if req.AutoUpload {
			ret.data = &randomanimal_pb.BasicResponse{
				Type: randomanimal_pb.Type_Image,
				Response: &randomanimal_pb.BasicResponse_UploadResponse{
					Hash:    chickenHash,
					URLPath: fmt.Sprintf("/assets/%v", chickenHash),
				},
			}
			return
		}
		resp, err := http.Get(fmt.Sprintf("http://%v:1080/assets/%v", define.GatewayIP.String(), chickenHash))
		if err != nil {
			ret.err = err
			return
		}
		defer resp.Body.Close()
		resp_body, err := io.ReadAll(resp.Body)
		if err != nil {
			ret.err = err
			return
		}
		ret.data = &randomanimal_pb.BasicResponse{
			Type: randomanimal_pb.Type_Image,
			Buf:  resp_body,
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

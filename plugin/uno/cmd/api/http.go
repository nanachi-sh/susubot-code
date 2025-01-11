package api

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nanachi-sh/susubot-code/plugin/uno/define"
	"github.com/nanachi-sh/susubot-code/plugin/uno/protos/uno"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type marshaler struct {
	*runtime.JSONPb
}

func GetMarshaler() *marshaler {
	return &marshaler{
		JSONPb: &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	}
}

func (m *marshaler) Marshal(v any) ([]byte, error) {
	var (
		code     int    = 200
		message  string = "successful"
		jwt      string
		response any
	)
	switch v := v.(type) {
	default:
		response = v
	case *uno.GetRoomsResponse:
		if len(v.Rooms) == 0 {
			response = []struct{}{}
		} else {
			response = v.Rooms
		}
	case *uno.CreateRoomResponse:
		if v.Err != nil {
			message = v.Err.String()
		} else {
			response = v
		}
	case *uno.GetRoomResponse:
		if v.Err != nil {
			message = v.Err.String()
		} else if v.Extra != nil {
			response = v.Extra
		} else {
			response = v.Simple
		}
	case *uno.GetPlayerResponse:
		if v.Err != nil {
			message = v.Err.String()
		} else if v.Extra != nil {
			response = v.Extra
		} else {
			response = v.Simple
		}
	case *uno.JoinRoomResponse:
		if v.Err != nil {
			message = v.Err.String()
		} else {
			response = v
		}
	case *uno.ExitRoomResponse:
		if v.Err != nil {
			message = v.Err.String()
		} else {
			response = v
		}
	case *uno.DrawCardResponse:
		if v.Err != nil {
			message = v.Err.String()
		} else {
			response = v
		}
	case *uno.SendCardResponse:
		if v.Err != nil {
			message = v.Err.String()
		} else {
			response = v
		}
	case *uno.NoSendCardResponse:
		if v.Err != nil {
			message = v.Err.String()
		} else {
			response = v
		}
	case *uno.CallUNOResponse:
		if v.Err != nil {
			message = v.Err.String()
		} else {
			response = v
		}
	case *uno.ChallengeResponse:
		if v.Err != nil {
			message = v.Err.String()
		} else {
			response = v
		}
	case *uno.IndicateUNOResponse:
		if v.Err != nil {
			message = v.Err.String()
		} else {
			response = v
		}
	case *uno.RoomEventResponse:
		if v.Err != nil {
			message = v.Err.String()
		} else {
			response = v
		}
	case *uno.GetUserResponse:
	case *uno.BasicResponse:
	}
	response = struct{}{}
	return m.JSONPb.Marshal(&struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		JWT     string `json:"jwt,omitempty"`
		Body    any    `json:"body"`
	}{
		code,
		message,
		jwt,
		response,
	})
}

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
	var conn *grpc.ClientConn
	if !define.EnableTLS {
		conn, err = grpc.NewClient(fmt.Sprintf("localhost:%v", gRPCport), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
	} else {
		cred, err := credentials.NewClientTLSFromFile(fmt.Sprintf("%v/tls.pem", define.CertsDir), "uno.unturned.fun")
		if err != nil {
			return err
		}
		conn, err = grpc.NewClient(fmt.Sprintf("localhost:%v", gRPCport), grpc.WithTransportCredentials(cred))
		if err != nil {
			return err
		}
	}
	sMux := runtime.NewServeMux(
		runtime.WithMarshalerOption("application/basicjson", &runtime.HTTPBodyMarshaler{Marshaler: GetMarshaler()}),
	)
	if err := uno.RegisterUnoHandler(context.Background(), sMux, conn); err != nil {
		return err
	}
	var l net.Listener
	if define.EnableTLS {
		cert, err := tls.LoadX509KeyPair(fmt.Sprintf("%v/tls.pem", define.CertsDir), fmt.Sprintf("%v/tls.key", define.CertsDir))
		if err != nil {
			return err
		}
		l, err = tls.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port), &tls.Config{
			Certificates: []tls.Certificate{cert},
		})
		if err != nil {
			return err
		}
	} else {
		l, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
		if err != nil {
			return err
		}
	}
	return http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") != "" {
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8080")
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

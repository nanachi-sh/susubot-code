package stream

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/configs"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  16 * 1024,
		WriteBufferSize: 16 * 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	xlogger = logx.WithContext(context.Background())
)

func Serve() error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", configs.WEBSOCKET_LISTEN_PORT))
	if err != nil {
		return err
	}
	mux := &http.ServeMux{}
	mux.HandleFunc("/ws", websocketHandle)
	return http.Serve(l, mux)
}

func websocketHandle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		xlogger.Error(err)
		return
	}
	defer conn.Close()
	// websocket stream
	for {
	}
}

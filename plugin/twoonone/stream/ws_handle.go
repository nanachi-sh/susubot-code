package stream

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/handler"
	inside "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  4 * 1024,
		WriteBufferSize: 8 * 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	//
	logger := logx.WithContext(r.Context())
	// 升级为websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	defer conn.Close()
	// 获取用户信息
	var req types.WebsocketHandShake
	if err := handler.ParseCustom(r, &req); err != nil {
		resp, _ := handler.Generate(nil, err)
		conn.WriteJSON(resp)
		return
	}
	// 获取对应桌事件流
	event, err := inside.NewAPIRequest(logger).RoomEvent(&req.Extra)
	if err != nil {
		resp, _ := handler.Generate(nil, err)
		conn.WriteJSON(resp)
		return
	}
	closed, close := context.WithCancel(context.Background())
	// websocket stream
	go func() {
		defer close()
		for {
			e, ok := event.Read()
			if !ok {
				// 桌事件结束
				return
			}
			select {
			case <-closed.Done():
				return
			default:
			}
			if ep := e.GetRoomExitPlayer(); ep != nil {
				if ep.LeaverInfo.User.Id == req.Extra.UserId {
					return
				}
			}
			resp, _ := handler.Generate(e, nil)
			if err := conn.WriteJSON(resp); err != nil {
				logger.Error(err)
				return
			}
		}
	}()
	go func() {
		defer close()
		for {
			time.Sleep(time.Second * 30)
			select {
			case <-closed.Done():
				return
			default:
			}
			if err := conn.WriteMessage(websocket.PingMessage, []byte("PING MESSAGE")); err != nil {
				logger.Error(err)
				return
			}
		}
	}()
	<-closed.Done()
}

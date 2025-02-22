package stream

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/handler"
	inside "github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/twoonone"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  16 * 1024,
		WriteBufferSize: 16 * 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	//
	logger := logx.WithContext(r.Context())
	fmt.Println(r.Header)
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
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	// 获取对应桌事件流
	event, err := inside.NewAPIRequest(logger).RoomEvent(&req.Extra)
	if err != nil {
		handler.Response(w, r, nil, err)
		return
	}
	// websocket stream
	for {
		e, ok := event.Read()
		if !ok {
			// 桌事件结束
			return
		}
		resp, _ := handler.Generate(e, nil)
		if err := conn.WriteJSON(resp); err != nil {
			logger.Error(err)
			return
		}
	}
}

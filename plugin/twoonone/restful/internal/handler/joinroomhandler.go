package handler

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/handler"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/logic"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func joinRoomHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JoinRoomRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewJoinRoomLogic(r.Context(), svcCtx)
		resp, err := l.JoinRoom(&req)

		handler.Response(w, r, resp, err)
	}
}

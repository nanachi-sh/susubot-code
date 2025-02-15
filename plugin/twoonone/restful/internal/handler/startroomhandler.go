package handler

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/handler"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/logic"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func startRoomHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StartRoomRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewStartRoomLogic(r.Context(), svcCtx)
		resp, err := l.StartRoom(&req)

		handler.Response(w, r, resp, err)
	}
}

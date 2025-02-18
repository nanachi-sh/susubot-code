package handler

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/handler"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/logic"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"
)

func createRoomHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewCreateRoomLogic(r.Context(), svcCtx)
		resp, err := l.CreateRoom()

		handler.Response(w, r, resp, err)
	}
}

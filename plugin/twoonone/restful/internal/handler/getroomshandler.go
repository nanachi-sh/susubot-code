package handler

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/handler"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/logic"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"
)

func getRoomsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetRoomsLogic(r.Context(), svcCtx)
		resp, err := l.GetRooms()

		handler.Response(w, r, resp, err)
	}
}

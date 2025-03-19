package handler

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/handler"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/logic"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/svc"
)

func ReverseProxy_GETHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewReverseProxy_GETLogic(r.Context(), svcCtx)
		err := l.ReverseProxy_GET()
		var resp any = nil
		handler.Response(w, r, resp, err)
	}
}

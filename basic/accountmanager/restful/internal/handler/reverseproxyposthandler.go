package handler

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/handler"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/logic"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/svc"
)

func ReverseProxy_POSTHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewReverseProxy_POSTLogic(r.Context(), svcCtx)
		err := l.ReverseProxy_POST()
		var resp any = nil
		handler.Response(w, r, resp, err)
	}
}

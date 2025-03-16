package handler

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/handler"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/logic"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/svc"
)

func VerifyCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewVerifyCodeLogic(r.Context(), svcCtx)
		resp, err := l.VerifyCode()

		handler.Response(w, r, resp, err)
	}
}

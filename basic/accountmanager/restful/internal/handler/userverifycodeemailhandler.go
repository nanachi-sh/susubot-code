package handler

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/handler"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/logic"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/restful/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserVerifyCode_EmailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserVerifyCodeEmailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserVerifyCode_EmailLogic(r.Context(), svcCtx)
		resp, err := l.UserVerifyCode_Email(&req)

		handler.Response(w, r, resp, err)
	}
}

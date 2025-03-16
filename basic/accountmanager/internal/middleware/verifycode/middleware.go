package verifycode

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/configs"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/internal/handler"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/pkg/protos/accountmanager"
	"github.com/nanachi-sh/susubot-code/basic/accountmanager/pkg/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	verify_id := r.Form.Get("verify_id")
	if verify_id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	answer := r.Form.Get("answer")
	if answer == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !configs.Captcha.Verify(verify_id, answer, false) {
		ret, httpCode := handler.Generate(nil, types.NewError(accountmanager.Error_ERROR_VERIFYCODE_ANSWER_FAIL, ""))
		httpx.WriteJsonCtx(r.Context(), w, httpCode, ret)
		return
	} else {
		configs.Captcha.Store.Get(verify_id, true)
	}
	next(w, r)
}

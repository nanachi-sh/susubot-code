package handler

import (
	"net/http"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type JSON_Response struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    any    `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, r *http.Request, resp any, err error) {
	ret := new(JSON_Response)
	statusCode := http.StatusOK
	if err != nil {
		if e, ok := err.(*types.AppError); ok {
			statusCode = e.StatusCode()
			ret.Code = int(e.Code)
			ret.Message = e.Message()
		} else {
			ret.Code = -1
			ret.Message = "failed"
		}
	} else {
		ret.Code = 0
		ret.Message = "OK"
		ret.Data = resp
	}
	ret.Data = resp
	httpx.WriteJsonCtx(r.Context(), w, statusCode, ret)
}

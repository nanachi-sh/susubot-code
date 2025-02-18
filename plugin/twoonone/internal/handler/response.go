package handler

import (
	"fmt"
	"net/http"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
)

type JSON_Response struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    any    `json:"data"`
}

func Response(w http.ResponseWriter, r *http.Request, resp any, err error) {
	ret := new(JSON_Response)
	if err != nil {
		if e, ok := err.(*types.AppError); ok {
			fmt.Println(e)
			w.WriteHeader(e.StatusCode())
			ret.Code = int(e.Code)
			ret.Message = e.Message()
		} else {
			ret.Code = -1
			ret.Message = "failed"
		}
	} else {
		ret.Code = 0
		ret.Message = "OK"
	}
}

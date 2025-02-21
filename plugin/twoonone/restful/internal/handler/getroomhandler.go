package handler

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/nanachi-sh/susubot-code/plugin/twoonone/internal/handler"
	pkg_types "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/types"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/logic"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/svc"
	"github.com/nanachi-sh/susubot-code/plugin/twoonone/restful/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func getRoomHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRoomRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		if err := handler.ParseCustom(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetRoomLogic(r.Context(), svcCtx)
		resp, err := l.GetRoom(&req)
		if req.Extra.Extra_update {
			m, ok := getExtraMap(resp)
			if !ok {
				httpx.ErrorCtx(r.Context(), w, errors.New("api handler error"))
			}
			m[pkg_types.EXTRA_KEY_extra] = w.Header().Get("authorization")[7:]
		}

		handler.Response(w, r, resp, err)
	}
}

func getExtraMap(resp any) (map[string]string, bool) {
	value := reflect.ValueOf(resp).Elem()
	fmt.Println(value.Type().Name())
	if extra := value.FieldByName("Extra"); !extra.IsNil() {
		m, ok := extra.Interface().(map[string]string)
		if !ok {
			return nil, false
		}
		return m, true
	} else {
		return nil, false
	}
}

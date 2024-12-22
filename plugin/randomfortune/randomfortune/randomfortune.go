package randomfortune

import (
	crand "crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/db"
	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/define"
	randomfortune_pb "github.com/nanachi-sh/susubot-code/plugin/randomfortune/protos/randomfortune"
)

func GetFortune(req *randomfortune_pb.BasicRequest) (*randomfortune_pb.BasicResponse, error) {
	if req.MemberId != nil {
		memberid := *req.MemberId
		ok := false
		for {
			ts, err := db.GetLastGetFortuneTime(memberid)
			if err != nil {
				if sql.ErrNoRows == err { //用户不存在
					ok = true
					if err := db.AddMember(memberid); err != nil {
						return nil, err
					}
					break
				}
				return nil, err
			}
			nowt := time.Now()
			mt := time.Unix(ts, 0)
			// 判断是否为同一年同一月
			if nowt.Month() == mt.Month() && nowt.Year() == mt.Year() {
				if nowt.Day() > mt.Day() {
					ok = true
					break
				}
			} else {
				//非同一年同一月必然不是同一天
				ok = true
				break
			}
		}
		if !ok {
			return &randomfortune_pb.BasicResponse{
				AlreadyGetFortune: true,
			}, nil
		}
	}
	d, err := os.ReadFile("/config/fortune_HashList.json")
	if err != nil {
		return nil, err
	}
	j := []string{}
	if err := json.Unmarshal(d, &j); err != nil {
		return nil, err
	}
	if len(j) == 0 {
		return nil, errors.New("fortune hashlist为空")
	}
	i, err := crand.Int(crand.Reader, big.NewInt(int64(len(j))))
	if err != nil {
		return nil, err
	}
	hash := j[i.Int64()]
	switch req.ReturnMethod {
	case randomfortune_pb.BasicRequest_Hash:
		if req.MemberId != nil {
			if err := db.UpdateLastGetFortuneTime(*req.MemberId, 0); err != nil {
				return nil, err
			}
		}
		return &randomfortune_pb.BasicResponse{
			Response: &randomfortune_pb.BasicResponse_UploadResponse{
				Hash:    hash,
				URLPath: fmt.Sprintf("/assets/%v", hash),
			},
		}, nil
	case randomfortune_pb.BasicRequest_Raw:
		resp, err := http.Get(fmt.Sprintf("%v:1080/assets/%v", define.GatewayIP.String(), hash))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		resp_body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if req.MemberId != nil {
			if err := db.UpdateLastGetFortuneTime(*req.MemberId, 0); err != nil {
				return nil, err
			}
		}
		return &randomfortune_pb.BasicResponse{
			Buf: resp_body,
		}, nil
	default:
		return nil, errors.New("未知ReturnMethod")
	}
}

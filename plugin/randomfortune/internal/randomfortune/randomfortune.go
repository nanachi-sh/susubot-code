package randomfortune

import (
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"

	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/internal/configs"
	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/internal/randomfortune/db"
	randomfortune_pb "github.com/nanachi-sh/susubot-code/plugin/randomfortune/pkg/protos/randomfortune"
	"github.com/zeromicro/go-zero/core/logx"
)

type Request struct {
	logger logx.Logger
}

func NewRequest(l logx.Logger) *Request {
	return &Request{logger: l}
}

func (r *Request) GetFortune(in *randomfortune_pb.BasicRequest) (*randomfortune_pb.BasicResponse, error) {
	resp, serr := getFortune(r.logger, in)
	if serr != nil {
		return &randomfortune_pb.BasicResponse{Body: &randomfortune_pb.BasicResponse_Err{Err: *serr}}, nil
	}
	return resp, nil
}

func getFortune(logger logx.Logger, in *randomfortune_pb.BasicRequest) (*randomfortune_pb.BasicResponse, *randomfortune_pb.Errors) {
	if in.MemberId != nil {
		memberid := *in.MemberId
		if serr := db.CheckPlayerTime(logger, memberid); serr != nil {
			return &randomfortune_pb.BasicResponse{Body: &randomfortune_pb.BasicResponse_Err{Err: *serr}}, nil
		}
	}
	d, err := os.ReadFile(configs.Randomfortune_HashListFile)
	if err != nil {
		logger.Error(err)
		return nil, randomfortune_pb.Errors_Undefined.Enum()
	}
	j := []string{}
	if err := json.Unmarshal(d, &j); err != nil {
		logger.Error(err)
		return nil, randomfortune_pb.Errors_Undefined.Enum()
	}
	if len(j) == 0 {
		logger.Error("hash表为空")
		return nil, randomfortune_pb.Errors_Undefined.Enum()
	}
	i, err := crand.Int(crand.Reader, big.NewInt(int64(len(j))))
	if err != nil {
		logger.Error(err)
		return nil, randomfortune_pb.Errors_Undefined.Enum()
	}
	hash := j[i.Int64()]
	switch in.ReturnMethod {
	case randomfortune_pb.BasicRequest_Hash:
		if in.MemberId != nil {
			memberid := *in.MemberId
			if serr := db.UpdatePlayerTime(logger, memberid); serr != nil {
				return nil, serr
			}
		}
		return &randomfortune_pb.BasicResponse{
			Body: &randomfortune_pb.BasicResponse_Hash{
				Hash: &randomfortune_pb.BasicResponse_UploadResponseFromHash{
					Hash: hash,
				},
			},
		}, nil
	case randomfortune_pb.BasicRequest_Raw:
		resp, err := http.Get(fmt.Sprintf("%s/%s", configs.ASSETS_URL, hash))
		if err != nil {
			logger.Error(err)
			return nil, randomfortune_pb.Errors_Undefined.Enum()
		}
		defer resp.Body.Close()
		resp_body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error(err)
			return nil, randomfortune_pb.Errors_Undefined.Enum()
		}
		if in.MemberId != nil {
			memberid := *in.MemberId
			if serr := db.UpdatePlayerTime(logger, memberid); serr != nil {
				return nil, serr
			}
		}
		return &randomfortune_pb.BasicResponse{
			Body: &randomfortune_pb.BasicResponse_Buf{
				Buf: &randomfortune_pb.BasicResponse_UploadResponseFromRaw{
					Buf: resp_body,
				},
			},
		}, nil
	default:
		logger.Error("未知ReturnMethod")
		return nil, randomfortune_pb.Errors_Undefined.Enum()
	}
}

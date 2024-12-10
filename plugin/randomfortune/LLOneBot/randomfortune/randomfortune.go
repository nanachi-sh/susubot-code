package randomfortune

import (
	crand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"

	"github.com/nanachi-sh/susubot-code/plugin/randomfortune/LLOneBot/define"
	randomfortune_pb "github.com/nanachi-sh/susubot-code/plugin/randomfortune/LLOneBot/protos/randomfortune"
)

func GetFortune(req *randomfortune_pb.BasicRequest) (*randomfortune_pb.BasicResponse, error) {
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
		return &randomfortune_pb.BasicResponse{
			Buf: resp_body,
		}, nil
	default:
		return nil, errors.New("未知ReturnMethod")
	}
}

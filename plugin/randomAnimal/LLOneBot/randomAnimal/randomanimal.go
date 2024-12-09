package randomanimal

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	randomanimal_pb "github.com/nanachi-sh/susubot-code/plugin/randomanimal/LLOneBot/protos/randomanimal"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/LLOneBot/randomanimal/define"
)

func GetCat() (*define.BasicReturn, error) {
	resp, err := http.Get("https://api.thecatapi.com/v1/images/search")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var j []*define.JSON_Cat
	if err := json.Unmarshal(resp_body, &j); err != nil {
		return nil, err
	}
	if len(j) == 0 {
		return nil, errors.New("获取猫图片失败")
	}
	url := j[0].URL
	assetResp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer assetResp.Body.Close()
	assetResp_body, err := io.ReadAll(assetResp.Body)
	if err != nil {
		return nil, err
	}
	var Type *randomanimal_pb.Type
	// 通过响应头判断类型
	switch strings.ToLower(assetResp.Header.Get("Content-Type")) {
	case "image/webp", "image/jpeg", "image/png", "image/gif":
		Type = randomanimal_pb.Type_Image.Enum()
	default:
		// 暂时懒得写
		return nil, errors.New("判断响应类型失败")
	}
	if Type == nil {
		return nil, errors.New("资源类型为nil")
	}
	return &define.BasicReturn{
		Buf:  assetResp_body,
		Type: *Type,
	}, nil
}

func GetDog() (*define.BasicReturn, error) {
	resp, err := http.Get("https://random.dog/woof.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	j := new(define.JSON_Dog)
	if err := json.Unmarshal(resp_body, j); err != nil {
		return nil, err
	}
	assetResp, err := http.Get(j.URL)
	if err != nil {
		return nil, err
	}
	defer assetResp.Body.Close()
	assetResp_body, err := io.ReadAll(assetResp.Body)
	if err != nil {
		return nil, err
	}
	var Type *randomanimal_pb.Type
	switch strings.ToLower(assetResp.Header.Get("Content-Type")) {
	case "image/webp", "image/jpeg", "image/png", "image/gif":
		Type = randomanimal_pb.Type_Image.Enum()
	case "video/mpeg4":
		Type = randomanimal_pb.Type_Video.Enum()
	default:
		return nil, errors.New("判断响应类型失败")
	}
	if Type == nil {
		return nil, errors.New("资源类型为nil")
	}
	return &define.BasicReturn{
		Buf:  assetResp_body,
		Type: *Type,
	}, nil
}

func GetFox() (*define.BasicReturn, error) {
	resp, err := http.Get("https://randomfox.ca/floof/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	j := new(define.JSON_Fox)
	if err := json.Unmarshal(resp_body, j); err != nil {
		return nil, err
	}
	assetResp, err := http.Get(j.URL)
	if err != nil {
		return nil, err
	}
	defer assetResp.Body.Close()
	var Type *randomanimal_pb.Type
	switch strings.ToLower(assetResp.Header.Get("Content-Type")) {
	case "image/webp", "image/jpeg", "image/png", "image/gif":
		Type = randomanimal_pb.Type_Image.Enum()
	default:
		return nil, errors.New("判断响应类型失败")
	}
	if Type == nil {
		return nil, errors.New("资源类型为nil")
	}
	assetResp_body, err := io.ReadAll(assetResp.Body)
	if err != nil {
		return nil, err
	}
	return &define.BasicReturn{
		Buf:  assetResp_body,
		Type: *Type,
	}, nil
}

func GetDuck() (*define.BasicReturn, error) {
	resp, err := http.Get("https://random-d.uk/api/randomimg")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var Type *randomanimal_pb.Type
	switch strings.ToUpper(resp.Header.Get("Content-Type")) {
	case "image/webp", "image/jpeg", "image/png", "image/gif":
		Type = randomanimal_pb.Type_Image.Enum()
	default:
		// 暂时懒得写
		return nil, errors.New("判断响应类型失败")
	}
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &define.BasicReturn{
		Buf:  resp_body,
		Type: *Type,
	}, nil
}

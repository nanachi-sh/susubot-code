package randomanimal

import (
	crand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"

	randomanimal_pb "github.com/nanachi-sh/susubot-code/plugin/randomanimal/protos/randomanimal"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/randomanimal/define"
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
	switch ct := strings.ToLower(assetResp.Header.Get("Content-Type")); ct {
	case "image/webp", "image/jpeg", "image/png", "image/gif":
		Type = randomanimal_pb.Type_Image.Enum()
	default:
		// 暂时懒得写
		return nil, fmt.Errorf("判断响应类型失败，Content-Type: %v", ct)
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
	switch ct := strings.ToLower(assetResp.Header.Get("Content-Type")); ct {
	case "image/webp", "image/jpeg", "image/png", "image/gif":
		Type = randomanimal_pb.Type_Image.Enum()
	case "video/mpeg4", "video/mp4":
		Type = randomanimal_pb.Type_Video.Enum()
	default:
		return nil, fmt.Errorf("判断响应类型失败，Content-Type: %v", ct)
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
	switch ct := strings.ToLower(assetResp.Header.Get("Content-Type")); ct {
	case "image/webp", "image/jpeg", "image/png", "image/gif":
		Type = randomanimal_pb.Type_Image.Enum()
	default:
		return nil, fmt.Errorf("判断响应类型失败，Content-Type: %v", ct)
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
	switch ct := strings.ToLower(resp.Header.Get("Content-Type")); ct {
	case "image/webp", "image/jpeg", "image/png", "image/gif":
		Type = randomanimal_pb.Type_Image.Enum()
	default:
		// 暂时懒得写
		return nil, fmt.Errorf("判断响应类型失败，Content-Type: %v", ct)
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

// return Hash
func GetChicken_CXK() (string, error) {
	d, err := os.ReadFile("/config/chickenCXK_HashList.json")
	if err != nil {
		return "", err
	}
	j := []string{}
	if err := json.Unmarshal(d, &j); err != nil {
		return "", err
	}
	if len(j) == 0 {
		return "", errors.New("Chicken_CXK Hash表为空")
	}
	i, err := crand.Int(crand.Reader, big.NewInt(int64(len(j))))
	if err != nil {
		return "", err
	}
	return j[i.Int64()], nil
}

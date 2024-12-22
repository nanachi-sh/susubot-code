package randomanimal

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
	"strconv"
	"strings"

	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/db"
	fileweb_pb "github.com/nanachi-sh/susubot-code/plugin/randomanimal/protos/fileweb"
	randomanimal_pb "github.com/nanachi-sh/susubot-code/plugin/randomanimal/protos/randomanimal"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/randomanimal/define"
	"github.com/twmb/murmur3"
)

const (
	s1 = 1752698478008528169
	s2 = 18166171646046525357
)

func murmurhash128ToString(seed1, seed2 uint64, id string) string {
	h1, h2 := murmur3.SeedStringSum128(seed1, seed2, id)
	builder := new(strings.Builder)
	builder.WriteString(strconv.FormatUint(h1, 16))
	builder.WriteString(strconv.FormatUint(h2, 16))
	return builder.String()
}

func cacheValidCheck(hash string) bool {
	resp, err := http.Head(fmt.Sprintf("%v:1080/assets/%v", define.GatewayIP.String(), hash))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func GetCat(autoupload bool) (*define.BasicReturn, error) {
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
	id := j[0].Id
	idhash := murmurhash128ToString(s1, s2, id)
	cache, err := db.FindCache(idhash, db.Cat)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	} else { //命中缓存
		if cacheValidCheck(cache.AssetHash) { //检查缓存有效性
			if autoupload {
				return &define.BasicReturn{
					Buf:  nil,
					Hash: &cache.AssetHash,
					Type: cache.Type,
				}, nil
			} else {
				resp, err := http.Get(fmt.Sprintf("http://%v:1080/assets/%v", define.GatewayIP.String(), cache.AssetHash))
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()
				resp_body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, err
				}
				return &define.BasicReturn{
					Buf:  resp_body,
					Hash: nil,
					Type: cache.Type,
				}, nil
			}
		} else { //缓存无效，删除记录
			if err := db.DeleteCache(idhash, db.Cat); err != nil {
				return nil, err
			}
		}
	}
	//未命中缓存
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
	if autoupload {
		resp, err := define.FilewebC.Upload(define.FilewebCtx, &fileweb_pb.UploadRequest{
			Buf: assetResp_body,
		})
		if err != nil {
			return nil, err
		}
		if err := db.AddCache(idhash, db.Cat, db.Cache{
			AssetHash: resp.Hash,
			Type:      *Type,
		}); err != nil {
			return nil, err
		}
		return &define.BasicReturn{
			Buf:  nil,
			Hash: &resp.Hash,
			Type: *Type,
		}, nil
	} else {
		return &define.BasicReturn{
			Buf:  assetResp_body,
			Type: *Type,
		}, nil
	}
}

func GetDog(autoupload bool) (*define.BasicReturn, error) {
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
	url := j.URL
	id := ""
	if len(url) > 19 {
		id = url[19:]
	} else {
		id = url
	}
	idhash := murmurhash128ToString(s1, s2, id)
	cache, err := db.FindCache(idhash, db.Dog)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	} else { //命中缓存
		if cacheValidCheck(cache.AssetHash) { //检查缓存有效性
			if autoupload {
				return &define.BasicReturn{
					Buf:  nil,
					Hash: &cache.AssetHash,
					Type: cache.Type,
				}, nil
			} else {
				resp, err := http.Get(fmt.Sprintf("http://%v:1080/assets/%v", define.GatewayIP.String(), cache.AssetHash))
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()
				resp_body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, err
				}
				return &define.BasicReturn{
					Buf:  resp_body,
					Hash: nil,
					Type: cache.Type,
				}, nil
			}
		} else { //缓存无效，删除记录
			if err := db.DeleteCache(idhash, db.Dog); err != nil {
				return nil, err
			}
		}
	}
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
	if autoupload {
		resp, err := define.FilewebC.Upload(define.FilewebCtx, &fileweb_pb.UploadRequest{
			Buf: assetResp_body,
		})
		if err != nil {
			return nil, err
		}
		if err := db.AddCache(idhash, db.Dog, db.Cache{
			AssetHash: resp.Hash,
			Type:      *Type,
		}); err != nil {
			return nil, err
		}
		return &define.BasicReturn{
			Buf:  nil,
			Hash: &resp.Hash,
			Type: *Type,
		}, nil
	} else {
		return &define.BasicReturn{
			Buf:  assetResp_body,
			Type: *Type,
		}, nil
	}
}

func GetFox(autoupload bool) (*define.BasicReturn, error) {
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
	id := ""
	if len(j.URL) > 28 {
		id = j.URL[28:]
	} else {
		id = j.URL
	}
	idhash := murmurhash128ToString(s1, s2, id)
	cache, err := db.FindCache(idhash, db.Fox)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	} else { //命中缓存
		if cacheValidCheck(cache.AssetHash) { //检查缓存有效性
			if autoupload {
				return &define.BasicReturn{
					Buf:  nil,
					Hash: &cache.AssetHash,
					Type: cache.Type,
				}, nil
			} else {
				resp, err := http.Get(fmt.Sprintf("http://%v:1080/assets/%v", define.GatewayIP.String(), cache.AssetHash))
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()
				resp_body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, err
				}
				return &define.BasicReturn{
					Buf:  resp_body,
					Hash: nil,
					Type: cache.Type,
				}, nil
			}
		} else { //缓存无效，删除记录
			if err := db.DeleteCache(idhash, db.Fox); err != nil {
				return nil, err
			}
		}
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
	if autoupload {
		resp, err := define.FilewebC.Upload(define.FilewebCtx, &fileweb_pb.UploadRequest{
			Buf: assetResp_body,
		})
		if err != nil {
			return nil, err
		}
		if err := db.AddCache(idhash, db.Fox, db.Cache{
			AssetHash: resp.Hash,
			Type:      *Type,
		}); err != nil {
			return nil, err
		}
		return &define.BasicReturn{
			Buf:  nil,
			Hash: &resp.Hash,
			Type: *Type,
		}, nil
	} else {
		return &define.BasicReturn{
			Buf:  assetResp_body,
			Type: *Type,
		}, nil
	}
}

func GetDuck(autoupload bool) (*define.BasicReturn, error) {
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
	if autoupload {
		resp, err := define.FilewebC.Upload(define.FilewebCtx, &fileweb_pb.UploadRequest{
			Buf: resp_body,
		})
		if err != nil {
			return nil, err
		}
		return &define.BasicReturn{
			Buf:  nil,
			Hash: &resp.Hash,
			Type: *Type,
		}, nil
	} else {
		return &define.BasicReturn{
			Buf:  resp_body,
			Type: *Type,
		}, nil
	}
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

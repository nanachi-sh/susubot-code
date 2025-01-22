package randomanimal

import (
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"

	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/configs"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/randomanimal/db"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/types"
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/utils"
	fileweb_pb "github.com/nanachi-sh/susubot-code/plugin/randomanimal/pkg/protos/fileweb"
	randomanimal_pb "github.com/nanachi-sh/susubot-code/plugin/randomanimal/pkg/protos/randomanimal"
	"github.com/zeromicro/go-zero/core/logx"
)

type Request struct {
	logger logx.Logger
}

func NewRequest(l logx.Logger) *Request {
	return &Request{logger: l}
}

func (r *Request) GetCat(in *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	ret, ok := getCat(r.logger)
	if !ok {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Err{Err: randomanimal_pb.Errors_NoGet},
		}, nil
	}
	if ret.Hash != "" {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Hash{
				Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
					Hash: ret.Hash,
					Type: ret.Type,
				},
			},
		}, nil
	} else {
		if in.AutoUpload {
			hash, ok := upload(r.logger, ret.Buf)
			if !ok {
				return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Err{Err: randomanimal_pb.Errors_NoUpload}}, nil
			}
			return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Hash{Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
				Hash: hash,
				Type: ret.Type,
			}}}, nil
		} else {
			return &randomanimal_pb.BasicResponse{
				Body: &randomanimal_pb.BasicResponse_Buf{Buf: &randomanimal_pb.BasicResponse_UploadResponseByBuf{
					Buf:  ret.Buf,
					Type: ret.Type,
				}},
			}, nil
		}
	}
}

func (r *Request) GetChiken_CXK(in *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	hash, ok := getChicken_CXK(r.logger)
	ret := types.BasicReturn{
		Hash: hash,
		Type: randomanimal_pb.Type_Image,
	}
	if !ok {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Err{Err: randomanimal_pb.Errors_NoGet},
		}, nil
	}
	if ret.Hash != "" {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Hash{
				Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
					Hash: ret.Hash,
					Type: ret.Type,
				},
			},
		}, nil
	} else {
		if in.AutoUpload {
			hash, ok := upload(r.logger, ret.Buf)
			if !ok {
				return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Err{Err: randomanimal_pb.Errors_NoUpload}}, nil
			}
			return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Hash{Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
				Hash: hash,
				Type: ret.Type,
			}}}, nil
		} else {
			return &randomanimal_pb.BasicResponse{
				Body: &randomanimal_pb.BasicResponse_Buf{Buf: &randomanimal_pb.BasicResponse_UploadResponseByBuf{
					Buf:  ret.Buf,
					Type: ret.Type,
				}},
			}, nil
		}
	}
}

func (r *Request) GetDog(in *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	ret, ok := getDog(r.logger)
	if !ok {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Err{Err: randomanimal_pb.Errors_NoGet},
		}, nil
	}
	if ret.Hash != "" {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Hash{
				Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
					Hash: ret.Hash,
					Type: ret.Type,
				},
			},
		}, nil
	} else {
		if in.AutoUpload {
			hash, ok := upload(r.logger, ret.Buf)
			if !ok {
				return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Err{Err: randomanimal_pb.Errors_NoUpload}}, nil
			}
			return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Hash{Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
				Hash: hash,
				Type: ret.Type,
			}}}, nil
		} else {
			return &randomanimal_pb.BasicResponse{
				Body: &randomanimal_pb.BasicResponse_Buf{Buf: &randomanimal_pb.BasicResponse_UploadResponseByBuf{
					Buf:  ret.Buf,
					Type: ret.Type,
				}},
			}, nil
		}
	}
}

func (r *Request) GetDuck(in *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	ret, ok := getDuck(r.logger)
	if !ok {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Err{Err: randomanimal_pb.Errors_NoGet},
		}, nil
	}
	if ret.Hash != "" {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Hash{
				Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
					Hash: ret.Hash,
					Type: ret.Type,
				},
			},
		}, nil
	} else {
		if in.AutoUpload {
			hash, ok := upload(r.logger, ret.Buf)
			if !ok {
				return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Err{Err: randomanimal_pb.Errors_NoUpload}}, nil
			}
			return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Hash{Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
				Hash: hash,
				Type: ret.Type,
			}}}, nil
		} else {
			return &randomanimal_pb.BasicResponse{
				Body: &randomanimal_pb.BasicResponse_Buf{Buf: &randomanimal_pb.BasicResponse_UploadResponseByBuf{
					Buf:  ret.Buf,
					Type: ret.Type,
				}},
			}, nil
		}
	}
}

func (r *Request) GetFox(in *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	ret, ok := getFox(r.logger)
	if !ok {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Err{Err: randomanimal_pb.Errors_NoGet},
		}, nil
	}
	if ret.Hash != "" {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Hash{
				Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
					Hash: ret.Hash,
					Type: ret.Type,
				},
			},
		}, nil
	} else {
		if in.AutoUpload {
			hash, ok := upload(r.logger, ret.Buf)
			if !ok {
				return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Err{Err: randomanimal_pb.Errors_NoUpload}}, nil
			}
			return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Hash{Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
				Hash: hash,
				Type: ret.Type,
			}}}, nil
		} else {
			return &randomanimal_pb.BasicResponse{
				Body: &randomanimal_pb.BasicResponse_Buf{Buf: &randomanimal_pb.BasicResponse_UploadResponseByBuf{
					Buf:  ret.Buf,
					Type: ret.Type,
				}},
			}, nil
		}
	}
}

func cacheValidCheck(hash string) bool {
	resp, err := http.Head(fmt.Sprintf("%s/%s", configs.ASSETS_URL, hash))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func upload(logger logx.Logger, buf []byte) (string, bool) {
	resp, err := configs.Call_Fileweb.Upload(configs.DefaultCtx, &fileweb_pb.UploadRequest{
		Buf: buf,
	})
	if err != nil {
		logger.Error(err)
		return "", false
	}
	if serr := resp.GetErr(); serr != fileweb_pb.Errors_EMPTY {
		switch serr {
		default:
			logger.Errorf("未处理错误：%s", serr.String())
		}
		return "", false
	}
	hash := resp.GetHash()
	return hash, true
}

func getCat(logger logx.Logger) (types.BasicReturn, bool) {
	resp, err := http.Get(configs.CatAPI)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	defer resp.Body.Close()
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	var j []*types.JSON_Cat
	if err := json.Unmarshal(resp_body, &j); err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	if len(j) == 0 {
		logger.Error("获取猫图片无结果")
		return types.BasicReturn{}, false
	}
	url := j[0].URL
	id := j[0].Id
	idhash := utils.Murmurhash128ToString(configs.SEED1, configs.SEED2, id)
	cache, ok := db.FindCache(logger, idhash)
	if ok {
		if cacheValidCheck(cache.AssetHash) { //检查缓存有效性
			return types.BasicReturn{
				Hash: cache.AssetHash,
				Type: randomanimal_pb.Type(randomanimal_pb.Type_value[cache.AssetType]),
			}, true
		} else { //缓存无效，删除记录
			db.DeleteCache(logger, idhash)
		}
	}
	//未命中缓存
	assetResp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	defer assetResp.Body.Close()
	assetResp_body, err := io.ReadAll(assetResp.Body)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	Type, err := utils.MatchMediaType(assetResp)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	return types.BasicReturn{
		Buf:  assetResp_body,
		Type: Type,
	}, true
}

func getDog(logger logx.Logger) (types.BasicReturn, bool) {
	resp, err := http.Get(configs.DogAPI)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	defer resp.Body.Close()
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	j := new(types.JSON_Dog)
	if err := json.Unmarshal(resp_body, j); err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	url := j.URL
	id := ""
	if len(url) > 19 {
		id = url[19:]
	} else {
		id = url
	}
	idhash := utils.Murmurhash128ToString(configs.SEED1, configs.SEED2, id)
	cache, ok := db.FindCache(logger, idhash)
	if ok {
		if cacheValidCheck(cache.AssetHash) { //检查缓存有效性
			return types.BasicReturn{
				Hash: cache.AssetHash,
				Type: randomanimal_pb.Type(randomanimal_pb.Type_value[cache.AssetType]),
			}, true
		} else { //缓存无效，删除记录
			db.DeleteCache(logger, idhash)
		}
	}
	assetResp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	defer assetResp.Body.Close()
	assetResp_body, err := io.ReadAll(assetResp.Body)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	Type, err := utils.MatchMediaType(assetResp)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	return types.BasicReturn{
		Buf:  assetResp_body,
		Type: Type,
	}, true
}

func getFox(logger logx.Logger) (types.BasicReturn, bool) {
	resp, err := http.Get(configs.FoxAPI)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	defer resp.Body.Close()
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	j := new(types.JSON_Fox)
	if err := json.Unmarshal(resp_body, j); err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	id := ""
	if len(j.URL) > 28 {
		id = j.URL[28:]
	} else {
		id = j.URL
	}
	idhash := utils.Murmurhash128ToString(configs.SEED1, configs.SEED2, id)
	cache, ok := db.FindCache(logger, idhash)
	if ok { //命中缓存
		if cacheValidCheck(cache.AssetHash) { //检查缓存有效性
			return types.BasicReturn{
				Hash: cache.AssetHash,
				Type: randomanimal_pb.Type(randomanimal_pb.Type_value[cache.AssetType]),
			}, true
		} else { //缓存无效，删除记录
			db.DeleteCache(logger, idhash)
		}
	}
	assetResp, err := http.Get(j.URL)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	defer assetResp.Body.Close()
	Type, err := utils.MatchMediaType(assetResp)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	assetResp_body, err := io.ReadAll(assetResp.Body)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	return types.BasicReturn{
		Buf:  assetResp_body,
		Type: Type,
	}, true
}

func getDuck(logger logx.Logger) (types.BasicReturn, bool) {
	resp, err := http.Get(configs.DuckAPI)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	defer resp.Body.Close()
	Type, err := utils.MatchMediaType(resp)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return types.BasicReturn{}, false
	}
	return types.BasicReturn{
		Buf:  resp_body,
		Type: Type,
	}, true
}

// return Hash
func getChicken_CXK(logger logx.Logger) (string, bool) {
	d, err := os.ReadFile("/config/chickenCXK_HashList.json")
	if err != nil {
		logger.Error(err)
		return "", false
	}
	j := []string{}
	if err := json.Unmarshal(d, &j); err != nil {
		logger.Error(err)
		return "", false
	}
	if len(j) == 0 {
		logger.Error("Chicken_CXK Hash表为空")
		return "", false
	}
	i, err := crand.Int(crand.Reader, big.NewInt(int64(len(j))))
	if err != nil {
		logger.Error(err)
		return "", false
	}
	return j[i.Int64()], true
}

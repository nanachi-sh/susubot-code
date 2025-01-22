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
	randomanimalmodel "github.com/nanachi-sh/susubot-code/plugin/randomanimal/internal/model/randomanimal"
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
	resp, serr := getCat(r.logger, in)
	if serr != nil {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Err{Err: *serr},
		}, nil
	}
	return resp, nil
}

func (r *Request) GetChiken_CXK(in *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	resp, serr := getChicken_CXK(r.logger)
	if serr != nil {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Err{Err: *serr},
		}, nil
	}
	return resp, nil
}

func (r *Request) GetDog(in *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	resp, serr := getDog(r.logger, in)
	if serr != nil {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Err{Err: *serr},
		}, nil
	}
	return resp, nil
}

func (r *Request) GetDuck(in *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	resp, serr := getDuck(r.logger, in)
	if serr != nil {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Err{Err: *serr},
		}, nil
	}
	return resp, nil
}

func (r *Request) GetFox(in *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, error) {
	resp, serr := getFox(r.logger, in)
	if serr != nil {
		return &randomanimal_pb.BasicResponse{
			Body: &randomanimal_pb.BasicResponse_Err{Err: *serr},
		}, nil
	}
	return resp, nil
}

func cacheValidCheck(hash string) bool {
	resp, err := http.Head(fmt.Sprintf("%s/%s", configs.ASSETS_URL, hash))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func upload(logger logx.Logger, buf []byte) (string, *randomanimal_pb.Errors) {
	resp, err := configs.Call_Fileweb.Upload(configs.DefaultCtx, &fileweb_pb.UploadRequest{
		Buf: buf,
	})
	if err != nil {
		logger.Error(err)
		return "", randomanimal_pb.Errors_Undefined.Enum()
	}
	if serr := resp.GetErr(); serr != fileweb_pb.Errors_EMPTY {
		switch serr {
		default:
			logger.Errorf("未处理错误：%s", serr.String())
		}
		return "", randomanimal_pb.Errors_Undefined.Enum()
	}
	hash := resp.GetHash()
	return hash, nil
}

func getCat(logger logx.Logger, in *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, *randomanimal_pb.Errors) {
	resp, err := http.Get(configs.CatAPI)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	defer resp.Body.Close()
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	var j []*types.JSON_Cat
	if err := json.Unmarshal(resp_body, &j); err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	if len(j) == 0 {
		logger.Error("获取猫图片无结果")
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	url := j[0].URL
	id := j[0].Id
	idhash := utils.Murmurhash128ToString(configs.SEED1, configs.SEED2, id)
	cache, ok := db.FindCache(logger, idhash)
	if ok {
		if cacheValidCheck(cache.AssetHash) { //检查缓存有效性
			return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Hash{Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
				Hash: cache.AssetHash,
				Type: randomanimal_pb.Type(randomanimal_pb.Type_value[cache.AssetType]),
			}}}, nil
		} else { //缓存无效，删除记录
			db.DeleteCache(logger, idhash)
		}
	}
	//未命中缓存
	assetResp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	defer assetResp.Body.Close()
	assetResp_body, err := io.ReadAll(assetResp.Body)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	Type, err := utils.MatchMediaType(assetResp)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	if in.AutoUpload {
		hash, serr := upload(logger, assetResp_body)
		if serr != nil {
			return nil, serr
		}
		db.AddCache(logger, idhash, hash, randomanimalmodel.AssetType(Type.String()))
		return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Hash{Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
			Hash: hash,
			Type: randomanimal_pb.Type(randomanimal_pb.Type_value[Type.String()]),
		}}}, nil
	}
	return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Buf{Buf: &randomanimal_pb.BasicResponse_UploadResponseByBuf{
		Buf:  assetResp_body,
		Type: Type,
	}}}, nil
}

func getDog(logger logx.Logger, in *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, *randomanimal_pb.Errors) {
	resp, err := http.Get(configs.DogAPI)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	defer resp.Body.Close()
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	j := new(types.JSON_Dog)
	if err := json.Unmarshal(resp_body, j); err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
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
			return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Hash{Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
				Hash: cache.AssetHash,
				Type: randomanimal_pb.Type(randomanimal_pb.Type_value[cache.AssetType]),
			}}}, nil
		} else { //缓存无效，删除记录
			db.DeleteCache(logger, idhash)
		}
	}
	assetResp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	defer assetResp.Body.Close()
	assetResp_body, err := io.ReadAll(assetResp.Body)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	Type, err := utils.MatchMediaType(assetResp)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	if in.AutoUpload {
		hash, serr := upload(logger, assetResp_body)
		if serr != nil {
			return nil, serr
		}
		db.AddCache(logger, idhash, hash, randomanimalmodel.AssetType(Type.String()))
		return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Hash{Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
			Hash: hash,
			Type: randomanimal_pb.Type(randomanimal_pb.Type_value[Type.String()]),
		}}}, nil
	}
	return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Buf{Buf: &randomanimal_pb.BasicResponse_UploadResponseByBuf{
		Buf:  assetResp_body,
		Type: Type,
	}}}, nil
}

func getFox(logger logx.Logger, in *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, *randomanimal_pb.Errors) {
	resp, err := http.Get(configs.FoxAPI)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	defer resp.Body.Close()
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	j := new(types.JSON_Fox)
	if err := json.Unmarshal(resp_body, j); err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
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
			return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Hash{Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
				Hash: cache.AssetHash,
				Type: randomanimal_pb.Type(randomanimal_pb.Type_value[cache.AssetType]),
			}}}, nil
		} else { //缓存无效，删除记录
			db.DeleteCache(logger, idhash)
		}
	}
	assetResp, err := http.Get(j.URL)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	defer assetResp.Body.Close()
	Type, err := utils.MatchMediaType(assetResp)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	assetResp_body, err := io.ReadAll(assetResp.Body)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	if in.AutoUpload {
		hash, serr := upload(logger, assetResp_body)
		if serr != nil {
			return nil, serr
		}
		db.AddCache(logger, idhash, hash, randomanimalmodel.AssetType(Type.String()))
		return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Hash{Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
			Hash: hash,
			Type: randomanimal_pb.Type(randomanimal_pb.Type_value[Type.String()]),
		}}}, nil
	}
	return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Buf{Buf: &randomanimal_pb.BasicResponse_UploadResponseByBuf{
		Buf:  assetResp_body,
		Type: Type,
	}}}, nil
}

func getDuck(logger logx.Logger, in *randomanimal_pb.BasicRequest) (*randomanimal_pb.BasicResponse, *randomanimal_pb.Errors) {
	resp, err := http.Get(configs.DuckAPI)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	defer resp.Body.Close()
	Type, err := utils.MatchMediaType(resp)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	if in.AutoUpload {
		hash, serr := upload(logger, resp_body)
		if serr != nil {
			return nil, serr
		}
		return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Hash{Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
			Hash: hash,
			Type: randomanimal_pb.Type(randomanimal_pb.Type_value[Type.String()]),
		}}}, nil
	}
	return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Buf{Buf: &randomanimal_pb.BasicResponse_UploadResponseByBuf{
		Buf:  resp_body,
		Type: Type,
	}}}, nil
}

// return Hash
func getChicken_CXK(logger logx.Logger) (*randomanimal_pb.BasicResponse, *randomanimal_pb.Errors) {
	d, err := os.ReadFile("/config/chickenCXK_HashList.json")
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	j := []string{}
	if err := json.Unmarshal(d, &j); err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	if len(j) == 0 {
		logger.Error("Chicken_CXK Hash表为空")
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	i, err := crand.Int(crand.Reader, big.NewInt(int64(len(j))))
	if err != nil {
		logger.Error(err)
		return nil, randomanimal_pb.Errors_Undefined.Enum()
	}
	hash := j[i.Int64()]
	return &randomanimal_pb.BasicResponse{Body: &randomanimal_pb.BasicResponse_Hash{Hash: &randomanimal_pb.BasicResponse_UploadResponseByHash{
		Hash: hash,
		Type: randomanimal_pb.Type_Image,
	}}}, nil
}

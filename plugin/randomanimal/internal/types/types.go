package types

import (
	"github.com/nanachi-sh/susubot-code/plugin/randomanimal/pkg/protos/randomanimal"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcClientConf
}

type (
	JSON_Cat struct {
		URL string `json:"url"`
		Id  string `json:"id"`
	}
	JSON_Dog struct {
		URL string `json:"url"`
	}
	JSON_Fox struct {
		URL string `json:"image"`
	}
)

type BasicReturn struct {
	Buf  []byte
	Hash string
	Type randomanimal.Type
}

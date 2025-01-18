package fileweb

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nanachi-sh/susubot-code/basic/fileweb/internal/configs"
	fileweb_pb "github.com/nanachi-sh/susubot-code/basic/fileweb/pkg/protos/fileweb"
	"github.com/twmb/murmur3"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Request struct {
	logger logx.Logger
}

func NewRequest(l logx.Logger) *Request {
	return &Request{
		logger: l,
	}
}

func murmurhash128ToString(seed1, seed2 uint64, buf []byte) string {
	h1, h2 := murmur3.SeedSum128(seed1, seed2, buf)
	builder := new(strings.Builder)
	builder.WriteString(strconv.FormatUint(h1, 16))
	builder.WriteString(strconv.FormatUint(h2, 16))
	return builder.String()
}

func (r *Request) Upload(in *fileweb_pb.UploadRequest) (*fileweb_pb.UploadResponse, error) {
	if len(in.Buf) == 0 {
		return &fileweb_pb.UploadResponse{}, status.Error(codes.InvalidArgument, "")
	}
	return upload(r.logger, in), nil
}

func upload(logger logx.Logger, in *fileweb_pb.UploadRequest) *fileweb_pb.UploadResponse {
	hash := strings.ToUpper(murmurhash128ToString(configs.SEED1, configs.SEED2, in.Buf))
	fi, err := os.Lstat(fmt.Sprintf("%v/%v", configs.WebDir, hash))
	if err != nil && os.IsExist(err) {
		logger.Error(err)
		return nil
	}
	if fi != nil {
		return &fileweb_pb.UploadResponse{
			Body: &fileweb_pb.UploadResponse_Hash{Hash: hash},
		}
	}
	if err := os.WriteFile(fmt.Sprintf("%v/%v", configs.WebDir, hash), in.Buf, 0744); err != nil {
		logger.Error(err)
		return &fileweb_pb.UploadResponse{
			Body: &fileweb_pb.UploadResponse_Err{Err: fileweb_pb.Errors_Undefined},
		}
	}
	return &fileweb_pb.UploadResponse{
		Body: &fileweb_pb.UploadResponse_Hash{Hash: hash},
	}
}

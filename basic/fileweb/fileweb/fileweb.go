package fileweb

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nanachi-sh/susubot-code/basic/fileweb/fileweb/define"
	"github.com/nanachi-sh/susubot-code/basic/fileweb/log"
	fileweb_pb "github.com/nanachi-sh/susubot-code/basic/fileweb/protos/fileweb"
	"github.com/twmb/murmur3"
)

var logger = log.Get()

const (
	s1 = 17322916394824533408
	s2 = 7960282306262499246
)

func murmurhash128ToString(seed1, seed2 uint64, buf []byte) string {
	h1, h2 := murmur3.SeedSum128(seed1, seed2, buf)
	builder := new(strings.Builder)
	builder.WriteString(strconv.FormatUint(h1, 16))
	builder.WriteString(strconv.FormatUint(h2, 16))
	return builder.String()
}

func Upload(req *fileweb_pb.UploadRequest) (*fileweb_pb.UploadResponse, error) {
	hash := strings.ToUpper(murmurhash128ToString(s1, s2, req.Buf))
	fi, err := os.Lstat(fmt.Sprintf("%v/%v", define.WorkDir, hash))
	if err != nil && os.IsExist(err) {
		logger.Println(err)
		return nil, err
	}
	if fi != nil {
		return &fileweb_pb.UploadResponse{
			Hash:    hash,
			URLPath: fmt.Sprintf("/assets/%v", hash),
		}, nil
	}
	if err := os.WriteFile(fmt.Sprintf("%v/%v", define.WorkDir, hash), req.Buf, 0744); err != nil {
		logger.Println(err)
		return nil, err
	}
	return &fileweb_pb.UploadResponse{
		Hash:    hash,
		URLPath: fmt.Sprintf("/assets/%v", hash),
	}, nil
}

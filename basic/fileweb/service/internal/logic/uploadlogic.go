package logic

import (
	"context"

	fw "github.com/nanachi-sh/susubot-code/basic/fileweb/internal/fileweb"
	"github.com/nanachi-sh/susubot-code/basic/fileweb/pkg/protos/fileweb"
	"github.com/nanachi-sh/susubot-code/basic/fileweb/service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadLogic) Upload(in *fileweb.UploadRequest) (*fileweb.UploadResponse, error) {
	return fw.NewRequest(l.Logger).Upload(in)
}

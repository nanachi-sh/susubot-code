package logic

import (
	"context"

	"service/C:/Users/User/works/susubot-code/basic/fileweb/pkg/protos/fileweb"
	"service/internal/svc"

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
	// todo: add your logic here and delete this line

	return &fileweb.UploadResponse{}, nil
}

// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5
// Source: fileweb.proto

package server

import (
	"context"

	"service/C:/Users/User/works/susubot-code/basic/fileweb/pkg/protos/fileweb"
	"service/internal/logic"
	"service/internal/svc"
)

type FileWebServer struct {
	svcCtx *svc.ServiceContext
	fileweb.UnimplementedFileWebServer
}

func NewFileWebServer(svcCtx *svc.ServiceContext) *FileWebServer {
	return &FileWebServer{
		svcCtx: svcCtx,
	}
}

func (s *FileWebServer) Upload(ctx context.Context, in *fileweb.UploadRequest) (*fileweb.UploadResponse, error) {
	l := logic.NewUploadLogic(ctx, s.svcCtx)
	return l.Upload(in)
}

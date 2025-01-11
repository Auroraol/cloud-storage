package logic

import (
	"context"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadByChunkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件分片上传
func NewFileUploadByChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadByChunkLogic {
	return &FileUploadByChunkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileUploadByChunkLogic) FileUploadByChunk(req *types.FileUploadByChunkRequest) (resp *types.FileUploadByChunkResponse, err error) {
	// todo: add your logic here and delete this line

	return
}

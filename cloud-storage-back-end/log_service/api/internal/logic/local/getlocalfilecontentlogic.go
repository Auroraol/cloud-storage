package local

import (
	"context"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLocalFileContentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLocalFileContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLocalFileContentLogic {
	return &GetLocalFileContentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLocalFileContentLogic) GetLocalFileContent(req *types.GetLocalFileContentReq) (resp *types.GetLocalFileContentRes, err error) {
	// 默认值
	if req.Limit <= 0 {
		req.Limit = 1024 * 100 // 默认读取100KB
	}

	// 读取文件内容
	content, err := l.svcCtx.LocalFileService.ReadFileContent(req.FilePath, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}

	return &types.GetLocalFileContentRes{
		Content: content,
		Length:  len(content),
	}, nil
}

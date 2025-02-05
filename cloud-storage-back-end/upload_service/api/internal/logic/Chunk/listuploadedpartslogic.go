package Chunk

import (
	"context"
	"errors"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListUploadedPartsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询分片上传状态
func NewListUploadedPartsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUploadedPartsLogic {
	return &ListUploadedPartsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUploadedPartsLogic) ListUploadedParts(req *types.ListPartsRequest) (resp *types.ListPartsResponse, err error) {
	// 1. 参数验证
	if req.UploadId == "" || req.Key == "" {
		return nil, errors.New("uploadId and key are required")
	}

	// 2. 获取OSS bucket
	bucket, err := l.svcCtx.OSSClient.Bucket(l.svcCtx.Config.OSSConfig.BucketName)
	if err != nil {
		l.Logger.Errorf("获取OSS bucket失败: %v", err)
		return nil, err
	}

	// 3. 列举已上传的分片
	lsRes, err := bucket.ListUploadedParts(req.Key, req.UploadId)
	if err != nil {
		l.Logger.Errorf("列举已上传分片失败: %v", err)
		return nil, err
	}

	// 4. 转换分片信息
	parts := make([]types.PartInfo, 0, len(lsRes.UploadedParts))
	var totalSize int64 = 0
	for _, part := range lsRes.UploadedParts {
		parts = append(parts, types.PartInfo{
			PartNumber: part.PartNumber,
			Size:       part.Size,
			ETag:       part.ETag,
		})
		totalSize += part.Size
	}

	// 5. 构建响应
	resp = &types.ListPartsResponse{
		Parts:      parts,
		TotalParts: len(parts),
		FileSize:   totalSize,
	}

	l.Logger.Infof("成功获取分片上传状态, uploadId: %s, key: %s, totalParts: %d",
		req.UploadId, req.Key, len(parts))

	return resp, nil
}

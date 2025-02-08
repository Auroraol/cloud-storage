package Chunk

import (
	"context"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/store/oss"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"
	ossSDK "github.com/aliyun/aliyun-oss-go-sdk/oss"
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
	// 获取OSS bucket
	bucket := oss.Bucket()
	if bucket == nil {
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "获取OSS Bucket失败")
	}

	// 构造分片上传初始化结果
	imur := ossSDK.InitiateMultipartUploadResult{
		Key:      req.Key,
		UploadID: req.UploadId,
		Bucket:   bucket.BucketName,
	}

	// 列举已上传的分片
	lsRes, err := bucket.ListUploadedParts(imur)
	if err != nil {
		l.Logger.Errorf("列举已上传分片失败: %v", err)
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "列举已上传分片失败")
	}

	// 转换分片信息
	parts := make([]types.PartInfo, 0, len(lsRes.UploadedParts))
	var totalSize int64 = 0
	for _, part := range lsRes.UploadedParts {
		parts = append(parts, types.PartInfo{
			PartNumber: part.PartNumber,
			Size:       int64(part.Size),
			ETag:       part.ETag,
		})
		totalSize += int64(part.Size)
	}

	// 构建响应
	resp = &types.ListPartsResponse{
		Parts:      parts,
		TotalParts: len(parts),
		FileSize:   totalSize,
	}

	l.Logger.Infof("成功获取分片上传状态, uploadId: %s, key: %s, totalParts: %d",
		req.UploadId, req.Key, len(parts))

	return resp, nil
}

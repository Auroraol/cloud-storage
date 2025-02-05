package Chunk

import (
	"context"
	"strconv"

	ossSDK "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/store/oss"
)

type CompleteMultipartUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 完成分片上传
func NewCompleteMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteMultipartUploadLogic {
	return &CompleteMultipartUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompleteMultipartUploadLogic) CompleteMultipartUpload(req *types.ChunkUploadCompleteRequest) (resp *types.ChunkUploadCompleteResponse, err error) {
	// 参数验证
	if req.UploadId == "" {
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "uploadId不能为空")
	}
	if req.Key == "" {
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "key不能为空")
	}
	if len(req.ETags) == 0 {
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "分片ETag列表不能为空")
	}

	// 获取OSS bucket
	bucket := oss.Bucket()
	if bucket == nil {
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "获取OSS Bucket失败")
	}

	// 构建完成分片上传的请求
	var parts []ossSDK.UploadPart
	for i, etag := range req.ETags {
		parts = append(parts, ossSDK.UploadPart{
			PartNumber: i + 1,
			ETag:       etag,
		})
	}

	// 完成分片上传
	imur := ossSDK.InitiateMultipartUploadResult{
		Bucket:   bucket.BucketName,
		Key:      req.Key,
		UploadID: req.UploadId,
	}

	result, err := bucket.CompleteMultipartUpload(imur, parts)
	if err != nil {
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "完成分片上传失败")
	}

	// 获取文件大小
	props, err := bucket.GetObjectDetailedMeta(req.Key)
	if err != nil {
		logx.Errorf("获取文件元数据失败: %v", err)
		// 继续执行，不影响返回结果
	}

	var size int64
	if sizeStr := props.Get("Content-Length"); sizeStr != "" {
		size, _ = strconv.ParseInt(sizeStr, 10, 64)
	}

	return &types.ChunkUploadCompleteResponse{
		URL:  result.Location,
		Size: size,
	}, nil
}

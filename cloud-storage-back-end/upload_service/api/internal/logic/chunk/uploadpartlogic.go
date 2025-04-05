package chunk

import (
	"context"
	"go.uber.org/zap"
	"mime/multipart"
	"os"

	ossSDK "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/store/oss"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/types"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/utils"
)

type UploadPartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 上传分片
func NewUploadPartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadPartLogic {
	return &UploadPartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 保存分片到临时文件
func (l *UploadPartLogic) UploadPart(req *types.ChunkUploadRequest, header *multipart.FileHeader) (resp *types.ChunkUploadResponse, err error) {

	// 保存分片到临时文件
	tempFile, err := utils.SaveUploadedFile(header)
	if err != nil {
		zap.S().Error("保存临时文件失败 err:%v", err)
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "保存临时文件失败")
	}
	defer utils.CleanupTempFile(tempFile) // 确保清理临时文件

	// 获取OSS bucket
	bucket := oss.Bucket()
	if bucket == nil {
		zap.S().Error("获取OSS Bucket失败")
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "获取OSS Bucket失败")
	}

	// 打开临时文件
	partFile, err := os.Open(tempFile)
	if err != nil {
		zap.S().Error("打开临时文件失败 err:%v", err)
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "打开临时文件失败")
	}
	defer partFile.Close()

	// 获取文件大小
	fileInfo, err := partFile.Stat()
	if err != nil {
		zap.S().Error("获取文件信息失败 err:%v", err)
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "获取文件信息失败")
	}

	// 上传分片
	imur := ossSDK.InitiateMultipartUploadResult{
		Bucket:   bucket.BucketName,
		Key:      req.Key,
		UploadID: req.UploadId,
	}

	part, err := bucket.UploadPart(imur, partFile, fileInfo.Size(), req.ChunkIndex)
	if err != nil {
		zap.S().Error("上传分片失败 err:%v", err)
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "上传分片失败")
	}

	return &types.ChunkUploadResponse{
		ETag: part.ETag,
	}, nil
}

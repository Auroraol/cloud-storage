package Chunk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Auroraol/cloud-storage/common/token"
	"go.uber.org/zap"
	"path/filepath"

	ossSDK "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/store/oss"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"
)

type InitiateMultipartUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 初始化分片上传
func NewInitiateMultipartUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitiateMultipartUploadLogic {
	return &InitiateMultipartUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitiateMultipartUploadLogic) InitiateMultipartUpload(req *types.ChunkUploadInitRequest) (resp *types.ChunkUploadInitResponse, err error) {
	// 检查文件大小
	if req.FileSize <= 0 {
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "文件大小必须大于0")
	}

	// 解析元数据
	var metadata map[string]string
	if req.Metadata != "" {
		if err := json.Unmarshal([]byte(req.Metadata), &metadata); err != nil {
			return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "解析元数据失败")
		}
	}

	// 构建OSS对象键（使用用户ID作为前缀）
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}
	objectKey := fmt.Sprintf("users/%d/files/%s", userId, filepath.Base(req.FileName))

	// 获取OSS bucket
	bucket := oss.Bucket()
	if bucket == nil {
		zap.S().Error("获取OSS Bucket失败")
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "获取OSS Bucket失败")
	}

	// 准备上传选项
	var options []ossSDK.Option
	for k, v := range metadata {
		options = append(options, ossSDK.Meta(k, v))
	}

	// 初始化分片上传
	imur, err := bucket.InitiateMultipartUpload(objectKey, options...)
	if err != nil {
		zap.S().Error("初始化分片上传失败 err:%v", err)
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "初始化分片上传失败")
	}

	return &types.ChunkUploadInitResponse{
		UploadId: imur.UploadID,
		Key:      objectKey,
	}, nil
}

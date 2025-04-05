package chunk

import (
	"context"
	"crypto/md5"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"go.uber.org/zap"

	ossSDK "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/store/oss"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/types"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/utils"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/model"
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
	// 获取OSS bucket
	bucket := oss.Bucket()
	if bucket == nil {
		zap.S().Error("获取OSS Bucket失败")
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
		zap.S().Error("完成分片上传失败 err:%s", err)
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "完成分片上传失败")
	}

	// 获取文件大小和其他元数据
	props, err := bucket.GetObjectDetailedMeta(req.Key)
	if err != nil {
		zap.S().Error("获取文件元数据失败 err:%v", err)
		// 继续执行，不影响返回结果
	}

	var size int64
	if sizeStr := props.Get("Content-Length"); sizeStr != "" {
		size, _ = strconv.ParseInt(sizeStr, 10, 64)
	}

	// 获取用户ID
	//userId := token.GetUidFromCtx(l.ctx)

	// 保存文件信息到数据库
	fileName := filepath.Base(req.Key)
	fileExt := filepath.Ext(fileName)
	identity := utils.IntUuid()

	// 生成文件哈希值
	// 由于分片上传无法直接计算整个文件的MD5，我们使用文件名、大小和上传时间的组合来生成唯一哈希值
	hashSource := fmt.Sprintf("%s_%d_%d", fileName, size, time.Now().UnixNano())
	hashMd5 := md5.Sum([]byte(hashSource))
	fileHash := fmt.Sprintf("%x", hashMd5)

	_, err = l.svcCtx.RepositoryPoolModel.InsertWithId(l.ctx, &model.RepositoryPool{
		Identity: uint64(identity),
		OssKey:   req.Key,
		Hash:     fileHash,
		Name:     fileName,
		Ext:      fileExt,
		Size:     size,
		Path:     result.Location,
	})
	if err != nil {
		zap.S().Error("保存文件信息失败, err: %s", err)
		// 继续执行，不影响返回结果
	}

	// 发送文件上传完成消息到 Pulsar (异步, 暂不实现)
	//if l.svcCtx.FilePublisher != nil {
	//	// 创建文件上传消息
	//	fileUploadedMsg := pulsar.NewFileUploadedMessage(
	//		strconv.FormatInt(int64(identity), 10),
	//		fileName,
	//		size,
	//		props.Get("Content-Type"),
	//		strconv.FormatInt(userId, 10),
	//		result.Location,
	//	)
	//
	//	// 设置文件哈希
	//	fileUploadedMsg.FileHash = fileHash
	//
	//	// 发送消息
	//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//	defer cancel()
	//
	//	_, err = l.svcCtx.FilePublisher.SendObject(ctx, fileUploadedMsg, map[string]string{
	//		"service":     "upload-service",
	//		"upload_type": "multipart",
	//	})
	//	if err != nil {
	//		// 只记录日志，不影响上传流程
	//		zap.S().Warnf("发送分片上传完成消息失败: %s", err)
	//	} else {
	//		zap.S().Infof("分片上传完成消息已发送: %s", fileUploadedMsg.FileID)
	//	}
	//}

	return &types.ChunkUploadCompleteResponse{
		URL:          result.Location,
		Size:         size,
		RepositoryId: int64(identity),
	}, nil
}

package logic

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/store/oss"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/token"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/rpc/client/auditservicerpc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/types"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/utils"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/model"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件上传
func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 计算文件的MD5
func calculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (l *FileUploadLogic) FileUpload(req *types.FileUploadRequest, r *http.Request) (resp *types.FileUploadResponse, err error) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		zap.S().Error("获取上传文件失败: %s", err)
		return nil, response.NewErrMsg("获取上传文件失败")
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			zap.S().Error("关闭文件失败: %s", err)
		}
	}(file)

	// 检查文件大小
	if fileHeader.Size > 20*1024*1024 { // 20MB
		return nil, response.NewErrMsg("文件大小超过限制，请使用分片上传")
	}

	// 判断是否已达用户容量上限
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {

		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 添加重试机制
	var volumeInfo *pb.FindVolumeResp
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		volumeInfo, err = l.svcCtx.UserCenterRpc.FindVolumeById(l.ctx, &pb.FindVolumeReq{Id: userId})
		if err == nil {
			break
		}
		if i < maxRetries-1 {
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return nil, response.NewErrMsg(fmt.Sprintf("调用用户中心服务RPC失败, err: %v", err))
	}

	if volumeInfo.NowVolume+fileHeader.Size > volumeInfo.TotalVolume {
		return nil, response.NewErrCode(response.FILE_TOO_LARGE_ERROR)
	}

	// 保存文件到临时目录
	tempFilePath, err := utils.SaveUploadedFile(fileHeader)
	if err != nil {
		zap.S().Error("保存临时文件失败, err: %s", err)
		return nil, response.NewErrMsg(fmt.Sprintf("保存临时文件失败, err: %v", err))
	}
	defer utils.CleanupTempFile(tempFilePath) // 确保清理临时文件

	// 计算文件MD5
	md5Str, err := calculateFileMD5(tempFilePath)
	if err != nil {
		zap.S().Error("计算文件MD5失败, err: %s", err)
		return nil, response.NewErrMsg(fmt.Sprintf("计算文件MD5失败: %v", err))
	}

	// 检查文件是否已存在（秒传）
	count, err := l.svcCtx.RepositoryPoolModel.CountByHash(l.ctx, md5Str)
	if err != nil {
		zap.S().Error("查询文件MD5失败, err: %s", err)
		return nil, response.NewErrMsg(fmt.Sprintf("查询文件MD5失败, err: %v", err))
	}
	if count > 0 {
		repositoryInfo, err := l.svcCtx.RepositoryPoolModel.FindRepositoryPoolByHash(l.ctx, md5Str)
		if err != nil {
			zap.S().Error("查询文件MD5失败, err: %s", err)
			return nil, response.NewErrMsg(fmt.Sprintf("查询文件MD5失败, err: %v", err))
		}
		return &types.FileUploadResponse{
			URL:          repositoryInfo.Path,
			Key:          repositoryInfo.OssKey,
			Size:         repositoryInfo.Size,
			RepositoryId: int64(repositoryInfo.Identity),
		}, nil
	}

	// 添加操作日志
	l.svcCtx.AuditLogServiceRpc.CreateOperationLog(l.ctx, &auditservicerpc.OperationLogReq{
		FileName: fileHeader.Filename,
		UserId:   userId,
		Content:  "上传文件",
		FileSize: int32(fileHeader.Size),
		Flag:     0,
	})

	// 准备OSS上传选项
	objectKey := fmt.Sprintf("files/%s/%s", strconv.FormatInt(userId, 10), utils.StringUuid())
	uploadOptions := oss.FileUploadOptions{
		FilePath:    tempFilePath,
		ObjectKey:   objectKey,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Metadata: map[string]string{
			"original-name": fileHeader.Filename,
			"upload-time":   time.Now().Format(time.RFC3339),
			"user-id":       fmt.Sprintf("%d", userId),
		},
	}

	// 上传文件
	fileUrl, err := oss.UploadFile(uploadOptions)
	if err != nil {
		zap.S().Error("文件上传失败, err: %s", err)
		return nil, response.NewErrMsg(fmt.Sprintf("文件上传失败: %v", err))
	}

	identity := utils.IntUuid()
	// 保存文件信息到数据库
	_, err = l.svcCtx.RepositoryPoolModel.InsertWithId(l.ctx, &model.RepositoryPool{
		Identity: uint64(identity),
		OssKey:   objectKey,
		Hash:     md5Str,
		Name:     fileHeader.Filename,
		Ext:      path.Ext(fileHeader.Filename),
		Size:     fileHeader.Size,
		Path:     fileUrl,
	})
	if err != nil {
		zap.S().Error("保存文件信息失败, err: %s", err)
		return nil, response.NewErrMsg(fmt.Sprintf("保存文件信息失败: %v", err))
	}

	// 更新用户存储容量
	_, err = l.svcCtx.UserCenterRpc.AddVolume(l.ctx, &pb.AddVolumeReq{
		Id:   userId,
		Size: fileHeader.Size,
	})
	if err != nil {
		zap.S().Error("更新用户存储容量失败, err: %s", err)
		return nil, response.NewErrMsg(fmt.Sprintf("更新用户存储容量失败: %v", err))
	}

	// 发送文件上传完成消息到 Pulsar (异步, 暂不实现)
	//if l.svcCtx.FilePublisher != nil {
	//	// 创建文件上传消息
	//	fileUploadedMsg := pulsar.NewFileUploadedMessage(
	//		strconv.FormatInt(int64(identity), 10),
	//		fileHeader.Filename,
	//		fileHeader.Size,
	//		fileHeader.Header.Get("Content-Type"),
	//		strconv.FormatInt(userId, 10),
	//		fileUrl,
	//	)
	//	// 设置文件哈希
	//	fileUploadedMsg.FileHash = md5Str
	//
	//	// 发送消息
	//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//	defer cancel()
	//
	//	_, err = l.svcCtx.FilePublisher.SendObject(ctx, fileUploadedMsg, map[string]string{
	//		"service": "upload-service",
	//	})
	//	if err != nil {
	//		// 只记录日志，不影响上传流程
	//		zap.S().Warnf("发送文件上传消息失败: %s", err)
	//	} else {
	//		zap.S().Infof("文件上传消息已发送: %s", fileUploadedMsg.FileID)
	//	}
	//}

	return &types.FileUploadResponse{
		URL:          fileUrl,
		Key:          objectKey,
		Size:         fileHeader.Size,
		RepositoryId: int64(identity),
	}, nil
}

package logic

//import (
//	"context"
//	"crypto/md5"
//	"fmt"
//	"io"
//	"mime/multipart"
//	"os"
//	"path"
//	"path/filepath"
//	"time"
//
//	"github.com/Auroraol/cloud-storage/common/response"
//	"github.com/Auroraol/cloud-storage/common/store/oss"
//	"github.com/Auroraol/cloud-storage/common/token"
//	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
//	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"
//	"github.com/Auroraol/cloud-storage/upload_service/model"
//	"github.com/Auroraol/cloud-storage/user_center/rpc/pb"
//
//	"github.com/google/uuid"
//	"github.com/zeromicro/go-zero/core/logx"
//)
//
//const (
//	// 上传方式的阈值
//	SimpleUploadLimit = 5 * 1024 * 1024   // 5MB以下使用简单上传
//	MultipartLimit    = 100 * 1024 * 1024 // 100MB以下使用断点续传
//	ChunkSize         = 5 * 1024 * 1024   // 分片大小为5MB
//)
//
//type FileUploadLogic struct {
//	logx.Logger
//	ctx    context.Context
//	svcCtx *svc.ServiceContext
//}
//
//// 文件上传
//func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
//	return &FileUploadLogic{
//		Logger: logx.WithContext(ctx),
//		ctx:    ctx,
//		svcCtx: svcCtx,
//	}
//}
//
//// 计算文件的MD5
//func calculateFileMD5(filePath string) (string, error) {
//	file, err := os.Open(filePath)
//	if err != nil {
//		return "", err
//	}
//	defer file.Close()
//
//	hash := md5.New()
//	if _, err := io.Copy(hash, file); err != nil {
//		return "", err
//	}
//
//	return fmt.Sprintf("%x", hash.Sum(nil)), nil
//}
//
//func (l *FileUploadLogic) FileUpload(req *types.FileUploadRequest, file multipart.File, fileHeader *multipart.FileHeader) (resp *types.FileUploadResponse, err error) {
//	// 判断是否已达用户容量上限
//	userId := token.GetUidFromCtx(l.ctx)
//	if userId == 0 {
//		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
//	}
//	volumeInfo, err := l.svcCtx.UserCenterRpc.FindVolumeById(l.ctx, &pb.FindVolumeReq{Id: userId})
//	if err != nil {
//		return nil, err
//	}
//	if volumeInfo.NowVolume+fileHeader.Size > volumeInfo.TotalVolume {
//		return nil, response.NewErrCode(response.FILE_TOO_LARGE_ERROR)
//	}
//
//	// 生成临时文件路径
//	tempDir := filepath.Join(os.TempDir(), "cloud-storage-uploads")
//	if err := os.MkdirAll(tempDir, 0755); err != nil {
//		return nil, fmt.Errorf("创建临时目录失败: %v", err)
//	}
//	tempFilePath := filepath.Join(tempDir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename))
//
//	// 保存上传的文件到临时文件
//	tempFile, err := os.Create(tempFilePath)
//	if err != nil {
//		return nil, fmt.Errorf("创建临时文件失败: %v", err)
//	}
//	defer func() {
//		tempFile.Close()
//		os.Remove(tempFilePath)
//	}()
//
//	if _, err := io.Copy(tempFile, file); err != nil {
//		return nil, fmt.Errorf("保存临时文件失败: %v", err)
//	}
//
//	// 计算文件MD5
//	md5Str, err := calculateFileMD5(tempFilePath)
//	if err != nil {
//		return nil, fmt.Errorf("计算文件MD5失败: %v", err)
//	}
//
//	// 检查文件是否已存在（秒传）
//	count, err := l.svcCtx.RepositoryPoolModel.CountByHash(l.ctx, md5Str)
//	if err != nil {
//		return nil, err
//	}
//	if count > 0 {
//		repositoryInfo, err := l.svcCtx.RepositoryPoolModel.FindRepositoryPoolByHash(l.ctx, md5Str)
//		if err != nil {
//			return nil, err
//		}
//		return &types.FileUploadResponse{
//			URL:  repositoryInfo.Path,
//			Key:  repositoryInfo.Identity,
//			Size: repositoryInfo.Size,
//		}, nil
//	}
//
//	// 准备OSS上传选项
//	objectKey := fmt.Sprintf("files/%s/%s", userId, uuid.New().String())
//	uploadOptions := oss.FileUploadOptions{
//		FilePath:    tempFilePath,
//		ObjectKey:   objectKey,
//		ContentType: fileHeader.Header.Get("Content-Type"),
//		Metadata: map[string]string{
//			"original-name": fileHeader.Filename,
//			"upload-time":   time.Now().Format(time.RFC3339),
//			"user-id":       fmt.Sprintf("%d", userId),
//			"upload-type":   l.determineUploadType(fileHeader.Size),
//		},
//	}
//
//	// 根据文件大小选择上传方式
//	var fileUrl string
//	switch l.determineUploadType(fileHeader.Size) {
//	case "simple":
//		// 简单上传
//		fileUrl, err = oss.NormalUpload(uploadOptions)
//	case "multipart":
//		// 断点续传
//		fileUrl, err = oss.ResumeUpload(uploadOptions)
//	case "chunk":
//		// 分片上传
//		fileUrl, err = oss.MultipartUpload(uploadOptions)
//	}
//
//	if err != nil {
//		return nil, fmt.Errorf("文件上传失败: %v", err)
//	}
//
//	// 保存文件信息到数据库
//	_, err = l.svcCtx.RepositoryPoolModel.Insert(l.ctx, &model.RepositoryPool{
//		Identity: objectKey,
//		Hash:     md5Str,
//		Name:     fileHeader.Filename,
//		Ext:      path.Ext(fileHeader.Filename),
//		Size:     fileHeader.Size,
//		Path:     fileUrl,
//	})
//	if err != nil {
//		return nil, fmt.Errorf("保存文件信息失败: %v", err)
//	}
//
//	// 更新用户存储容量
//	_, err = l.svcCtx.UserCenterRpc.AddVolume(l.ctx, &pb.AddVolumeReq{
//		Id:   userId,
//		Size: fileHeader.Size,
//	})
//	if err != nil {
//		return nil, fmt.Errorf("更新用户存储容量失败: %v", err)
//	}
//
//	return &types.FileUploadResponse{
//		URL:  fileUrl,
//		Key:  objectKey,
//		Size: fileHeader.Size,
//	}, nil
//}
//
//// 根据文件大小确定上传方式
//func (l *FileUploadLogic) determineUploadType(fileSize int64) string {
//	switch {
//	case fileSize <= SimpleUploadLimit:
//		return "simple"
//	case fileSize <= MultipartLimit:
//		return "multipart"
//	default:
//		return "chunk"
//	}
//}

//
//// 文件上传
//func CosUpload(fileHeader *multipart.FileHeader, newId int64, b []byte) (string, string, error) {
//	u, _ := url.Parse(CosUrl)
//	bs := &cos.BaseURL{BucketURL: u}
//	c := cos.NewClient(bs, &http.Client{
//		Transport: &cos.AuthorizationTransport{
//			SecretID:  SecretID,
//			SecretKey: SecretKey,
//		},
//	})
//	baseName := path.Base(fileHeader.Filename)
//	name := "butane-netdisk/" + strconv.FormatInt(newId, 10) + baseName
//	_, err := c.Object.Put(context.Background(), name, bytes.NewReader(b), nil)
//	if err != nil {
//		return "", "", err
//	}
//	filePath := CosUrl + "/" + name
//	return filePath, baseName, nil
//}

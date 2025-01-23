package logic

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	"github.com/Auroraol/cloud-storage/common/tool"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"
	"github.com/Auroraol/cloud-storage/upload_service/model"
	"github.com/Auroraol/cloud-storage/user_center/rpc/client/user"
	uuid "github.com/satori/go.uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"mime/multipart"
	"path"
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

// 根据md5码查看文件是否存在，若存在则秒传成功。若不存在则插入数据库数据并拼接name和ext，上传文件到cos
func (l *FileUploadLogic) FileUpload(req *types.FileUploadRequest, file multipart.File, fileHeader *multipart.FileHeader) (resp *types.FileUploadResponse, err error) {
	// 判断是否已达用户容量上限
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}
	volumeInfo, err := l.svcCtx.UserCenterRpc.FindVolumeById(l.ctx, &user.FindVolumeReq{Id: userId})
	if err != nil {
		return nil, err
	}
	if volumeInfo.NowVolume+fileHeader.Size > volumeInfo.TotalVolume {
		return nil, response.NewErrCode(response.FILE_TOO_LARGE_ERROR)
	}
	// 增加用户当前已存储容量
	_, err = l.svcCtx.UserCenterRpc.AddVolume(l.ctx, &user.AddVolumeReq{
		Id:   userId,
		Size: fileHeader.Size,
	})
	if err != nil {
		return nil, err
	}

	//判断文件是否已存在，若已存在则为秒传成功
	b := make([]byte, fileHeader.Size)
	_, err = file.Read(b)
	if err != nil {
		return nil, err
	}
	md5Str := tool.Md5ByBytes(b)
	//
	count, err := l.svcCtx.RepositoryPoolModel.CountByHash(l.ctx, md5Str)
	if count > 0 {
		repositoryInfo, err := l.svcCtx.RepositoryPoolModel.FindRepositoryPoolByHash(l.ctx, md5Str)
		if err != nil {
			return nil, err
		}
		return &types.FileUploadResponse{Id: repositoryInfo.Identity}, err
	}

	// 判断使用的存储引擎，默认使用COS
	// 上传文件到cos，并得到filepath
	//filePath, baseName, err := utils.CosUpload(fileHeader, newId, b)
	//if err != nil {
	//	return nil, err
	//}
	//写入异步转移任务队列-->mq

	var baseName string
	var filePath string

	// 插入数据
	newId := uuid.NewV4().String()
	_, err = l.svcCtx.RepositoryPoolModel.InsertWithId(l.ctx, &model.RepositoryPool{
		Identity: newId,
		Hash:     md5Str,
		Ext:      path.Ext(baseName),
		Size:     fileHeader.Size,
		Path:     filePath,
		Name:     fileHeader.Filename,
	})
	if err != nil {
		return nil, response.NewErrCode(response.FILE_UPLOAD_ERROR)
	}
	return &types.FileUploadResponse{Id: newId}, err
}

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

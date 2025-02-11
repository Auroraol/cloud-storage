package repository

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	uploadServicePb "github.com/Auroraol/cloud-storage/upload_service/rpc/pb"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户文件列表
func NewUserFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileListLogic {
	return &UserFileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileListLogic) UserFileList(req *types.UserFileListRequest) (resp *types.UserFileListResponse, err error) {
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	allUserRepository, err := l.svcCtx.UserRepositoryModel.FindAllFileByParentId(l.ctx, req.Id, userId)
	if err != nil {
		return nil, response.NewErrMsg("该文件夹下搜索文件失败！")
	}
	// 获得所有文件信息
	newList := make([]*types.UserFile, 0)
	for _, userRepository := range allUserRepository {
		repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadServicePb.RepositoryReq{RepositoryId: int64(userRepository.RepositoryId)})
		if err != nil {
			logx.Errorf("获取文件信息失败 err:%v", err)
			continue
			//return nil, err
		}
		newList = append(newList, &types.UserFile{
			Id:           int64(userRepository.Id),
			RepositoryId: int64(userRepository.RepositoryId),
			Name:         userRepository.Name,
			Ext:          repositoryInfo.Ext,
			Path:         repositoryInfo.Path,
			Size:         repositoryInfo.Size,
			UpdateTime:   repositoryInfo.UpdateTime,
		})
	}
	return &types.UserFileListResponse{
		List:  newList,
		Count: int64(len(allUserRepository)),
	}, err
}

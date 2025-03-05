package repository

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	uploadServicePb "github.com/Auroraol/cloud-storage/upload_service/rpc/pb"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFolderSizeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户文件移动
func NewUserFolderSizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFolderSizeLogic {
	return &UserFolderSizeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFolderSizeLogic) UserFolderSize(req *types.UserFolderSizeRequest) (resp *types.UserFolderSizeResponse, err error) {

	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 查询同一目录下的文件/文件夹
	allUserRepository, err := l.svcCtx.UserRepositoryModel.FindAllFolderAndByParentId(l.ctx, req.Id, userId)
	if err != nil {
		zap.S().Error("该文件夹下搜索文件夹失败！")
		return nil, response.NewErrMsg("该文件夹下搜索文件夹失败！")
	}
	size := 0
	for _, userRepository := range allUserRepository {
		repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadServicePb.RepositoryReq{RepositoryId: int64(userRepository.RepositoryId)})
		if err != nil {
			zap.S().Error("获取文件信息失败 err:%v", err)
			continue
			//return nil, err
		}
		size += int(repositoryInfo.Size)
	}
	return &types.UserFolderSizeResponse{
		Size: size,
	}, nil
}

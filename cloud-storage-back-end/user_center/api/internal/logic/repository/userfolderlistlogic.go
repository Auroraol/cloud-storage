package repository

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFolderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户文件夹列表
func NewUserFolderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFolderListLogic {
	return &UserFolderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFolderListLogic) UserFolderList(req *types.UserFolderListRequest) (resp *types.UserFolderListResponse, err error) {
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	allUserRepository, err := l.svcCtx.UserRepositoryModel.FindAllFolderByParentId(l.ctx, req.Id, userId)
	if err != nil {
		return nil, response.NewErrMsg("该文件夹下搜索文件夹失败！")
	}
	newList := make([]*types.UserFolder, 0)
	for _, userRepository := range allUserRepository {
		newList = append(newList, &types.UserFolder{
			Id:   int64(userRepository.Id),
			Name: userRepository.Name,
		})
	}
	return &types.UserFolderListResponse{List: newList}, nil
}

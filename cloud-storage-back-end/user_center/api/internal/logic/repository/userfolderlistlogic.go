package repository

import (
	"context"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/time"
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
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 只查询正常状态的文件夹
	allUserRepository, err := l.svcCtx.UserRepositoryModel.FindAllNormalFolderByParentId(l.ctx, req.Id, userId)
	if err != nil {
		zap.S().Error("该文件夹下搜索文件夹失败！")
		return nil, response.NewErrMsg("该文件夹下搜索文件夹失败！")
	}
	// 获得所有文件夹信息
	newList := make([]*types.UserFolder, 0)
	for _, userRepository := range allUserRepository {
		// 忽略不是文件夹
		if userRepository.RepositoryId != 0 {
			continue
		}
		s := userRepository.UpdateTime.String()
		timePart := s[:19]
		timestamp, err := time.StringTimeToTimestamp(timePart)
		if err != nil {
			zap.S().Error("转化时间戳失败 err:%s", err)
			return nil, err
		}
		newList = append(newList, &types.UserFolder{
			Id:         int64(userRepository.Id),
			Name:       userRepository.Name,
			UpdateTime: timestamp,
		})
	}
	return &types.UserFolderListResponse{List: newList}, nil
}

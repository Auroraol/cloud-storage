package recycle

import (
	"context"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/token"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/types"

	uploadServicePb "github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/rpc/pb"
	userCenterPb "github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserRecycleRestoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户回收站文件恢复
func NewUserRecycleRestoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRecycleRestoreLogic {
	return &UserRecycleRestoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRecycleRestoreLogic) UserRecycleRestore(req *types.UserRecycleRestoreRequest) (resp *types.UserRecycleRestoreResponse, err error) {
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 获取文件信息
	fileInfo, err := l.svcCtx.UserRepositoryModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		zap.S().Error("文件不存在 err:%v", err)
		return nil, response.NewErrMsg("文件不存在")
	}

	if fileInfo.UserId != uint64(userId) {
		zap.S().Error("无权操作此文件")
		return nil, response.NewErrMsg("无权操作此文件")
	}

	// 更新状态为正常
	fileInfo.Status = 0
	err = l.svcCtx.UserRepositoryModel.Update(l.ctx, fileInfo)
	if err != nil {
		zap.S().Error("恢复文件失败 err:%v", err)
		return nil, response.NewErrMsg("恢复文件失败")
	}

	// 如果是文件，更新用户存储容量
	if fileInfo.RepositoryId != 0 {
		repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadServicePb.RepositoryReq{
			RepositoryId: int64(fileInfo.RepositoryId),
		})
		if err == nil {
			_, err = l.svcCtx.UserCenterRpc.AddVolume(l.ctx, &userCenterPb.AddVolumeReq{
				Id:   userId,
				Size: repositoryInfo.Size,
			})
			if err != nil {
				zap.S().Error("更新用户存储容量失败 err:%v", err)
			}
		}
	}

	// 如果是文件夹，递归恢复其下的所有文件和文件夹
	if fileInfo.RepositoryId == 0 {
		err = l.restoreFolderContents(l.ctx, int64(fileInfo.Id), userId)
		if err != nil {
			return nil, err
		}
	}

	return &types.UserRecycleRestoreResponse{
		Success: true,
	}, nil
}

func (l *UserRecycleRestoreLogic) restoreFolderContents(ctx context.Context, parentId int64, userId int64) error {
	children, err := l.svcCtx.UserRepositoryModel.FindAllDeletedByParentId(ctx, parentId, userId)
	if err != nil {
		return err
	}

	for _, child := range children {
		child.Status = 0
		err = l.svcCtx.UserRepositoryModel.Update(ctx, child)
		if err != nil {
			zap.S().Error("恢复文件夹内容失败 err:%v", err)
			return err
		}

		// 如果是文件，更新用户存储容量
		if child.RepositoryId != 0 {
			repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(ctx, &uploadServicePb.RepositoryReq{
				RepositoryId: int64(child.RepositoryId),
			})
			if err == nil {
				_, err = l.svcCtx.UserCenterRpc.AddVolume(ctx, &userCenterPb.AddVolumeReq{
					Id:   userId,
					Size: repositoryInfo.Size,
				})
				if err != nil {
					zap.S().Error("更新用户存储容量失败 err:%v", err)
				}
			}
		} else if child.RepositoryId == 0 {
			err = l.restoreFolderContents(ctx, int64(child.Id), userId)
			if err != nil {
				zap.S().Error("恢复文件夹内容失败 err:%v", err)
				return err
			}
		}
	}
	return nil
}

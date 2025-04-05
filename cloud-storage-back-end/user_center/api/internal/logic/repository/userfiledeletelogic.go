package repository

import (
	"context"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/token"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/rpc/client/auditservicerpc"
	uploadServicePb "github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/rpc/pb"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/types"
	userCenterPb "github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/pb"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UserFileDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户文件删除
func NewUserFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileDeleteLogic {
	return &UserFileDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileDeleteLogic) UserFileDelete(req *types.UserFileDeleteRequest) (resp *types.UserFileDeleteResponse, err error) {
	// 先获取文件信息
	userFileInfo, err := l.svcCtx.UserRepositoryModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		zap.S().Error("文件不存在 err:%s", err)
		return nil, err
	}

	// 更新status为已删除状态,而不是真正删除
	userFileInfo.Status = 1 // 1表示已删除
	err = l.svcCtx.UserRepositoryModel.Update(l.ctx, userFileInfo)
	if err != nil {
		zap.S().Error("更新文件状态失败 err:%v", err)
		return nil, response.NewErrMsg("更新文件状态失败")
	}

	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 添加操作日志文件id
	fileId := userFileInfo.RepositoryId
	if userFileInfo.ParentId != 0 {
		fileId = userFileInfo.ParentId
	}

	// 如果是文件夹,递归更新子文件和文件夹的状态
	if userFileInfo.RepositoryId == 0 {
		err = l.updateFolderContentsStatus(l.ctx, int64(userFileInfo.Id), userId)
		if err != nil {
			zap.S().Error("更新文件夹状态失败 err:%v", err)
			return nil, err
		}

		// 添加文件夹删除操作日志
		l.svcCtx.AuditLogServiceRpc.CreateOperationLog(l.ctx, &auditservicerpc.OperationLogReq{
			UserId:   userId,
			Content:  "删除文件夹",
			Flag:     2,
			FileId:   int64(fileId),
			FileName: userFileInfo.Name,
		})

		return &types.UserFileDeleteResponse{}, nil
	}

	// 从中心存储池中取 size (文件信息)
	repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadServicePb.RepositoryReq{RepositoryId: int64(userFileInfo.RepositoryId)})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		zap.S().Error("中心存储池找不到该数据 err:%v", err)
		return nil, response.NewErrMsg("中心存储池找不到该数据")
	}

	// 更新用户存储容量
	_, err = l.svcCtx.UserCenterRpc.DecreaseVolume(l.ctx, &userCenterPb.DecreaseVolumeReq{
		Id:   userId,
		Size: repositoryInfo.Size,
	})
	if err != nil {
		zap.S().Error("更新容量失败 err:%v", err)
		return nil, response.NewErrMsg("更新容量失败！")
	}

	// 添加文件删除操作日志
	l.svcCtx.AuditLogServiceRpc.CreateOperationLog(l.ctx, &auditservicerpc.OperationLogReq{
		UserId:   userId,
		Content:  "删除文件",
		FileSize: int32(repositoryInfo.Size),
		Flag:     2,
		FileId:   int64(fileId),
		FileName: userFileInfo.Name,
	})

	return &types.UserFileDeleteResponse{}, nil
}

// 新增递归更新文件夹内容状态的方法
func (l *UserFileDeleteLogic) updateFolderContentsStatus(ctx context.Context, parentId int64, userId int64) error {
	children, err := l.svcCtx.UserRepositoryModel.FindAllFolderAndByParentId(ctx, parentId, userId)
	if err != nil {
		zap.S().Error("更新文件夹状态失败 err:%v", err)
		return err
	}

	for _, child := range children {
		// 更新子项的状态为已删除
		child.Status = 1
		err = l.svcCtx.UserRepositoryModel.Update(ctx, child)
		if err != nil {
			zap.S().Error("更新子项状态失败: %v", err)
			return err
		}
		if child.RepositoryId == 0 {
			// 如果是文件夹,递归更新其内容
			err = l.updateFolderContentsStatus(ctx, int64(child.Id), userId)
			if err != nil {
				zap.S().Error("更新子项状态失败: %v", err)
				return err
			}
		} else {
			// 如果是文件，获取文件大小并更新用户存储容量
			repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(ctx, &uploadServicePb.RepositoryReq{RepositoryId: int64(child.RepositoryId)})
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 如果文件在中心存储池中找不到，可以选择跳过或记录日志
					zap.S().Error("File with RepositoryId %d not found in center storage pool", child.RepositoryId)
					continue
				}
				return err
			}

			// 更新用户存储容量
			_, err = l.svcCtx.UserCenterRpc.DecreaseVolume(ctx, &userCenterPb.DecreaseVolumeReq{
				Id:   userId,
				Size: repositoryInfo.Size,
			})
			if err != nil {
				zap.S().Error("更新容量失败 err:%v", err)
				return err
			}
		}
	}
	return nil
}

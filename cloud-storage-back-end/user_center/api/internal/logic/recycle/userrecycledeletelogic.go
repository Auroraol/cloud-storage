package recycle

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	"github.com/Auroraol/cloud-storage/log_service/rpc/client/auditservicerpc"
	uploadServicePb "github.com/Auroraol/cloud-storage/upload_service/rpc/pb"
	userCenterPb "github.com/Auroraol/cloud-storage/user_center/rpc/pb"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRecycleDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户回收站文件删除
func NewUserRecycleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRecycleDeleteLogic {
	return &UserRecycleDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRecycleDeleteLogic) UserRecycleDelete(req *types.UserRecycleDeleteRequest) (resp *types.UserRecycleDeleteResponse, err error) {
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 获取文件信息
	fileInfo, err := l.svcCtx.UserRepositoryModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		return nil, response.NewErrMsg("文件不存在")
	}

	if fileInfo.UserId != uint64(userId) {
		return nil, response.NewErrMsg("无权操作此文件")
	}

	// 如果是文件夹，递归删除其内容
	if fileInfo.RepositoryId == 0 {
		err = l.deleteFolderContents(l.ctx, int64(fileInfo.Id), userId)
		if err != nil {
			return nil, err
		}
	} else {
		// 如果是文件，更新用户存储容量
		repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadServicePb.RepositoryReq{RepositoryId: int64(fileInfo.RepositoryId)})
		if err == nil {
			_, err = l.svcCtx.UserCenterRpc.DecreaseVolume(l.ctx, &userCenterPb.DecreaseVolumeReq{
				Id:   userId,
				Size: repositoryInfo.Size,
			})
			if err != nil {
				logx.Errorf("更新用户存储容量失败 err:%v", err)
			}
		}
	}

	// 彻底删除文件记录
	err = l.svcCtx.UserRepositoryModel.Delete(l.ctx, fileInfo.Id)
	if err != nil {
		return nil, response.NewErrMsg("删除文件失败")
	}

	l.svcCtx.AuditLogServiceRpc.CreateOperationLog(l.ctx, &auditservicerpc.OperationLogReq{
		UserId:   userId,
		Content:  "删除回收站文件",
		Flag:     2,
		FileId:   int64(fileInfo.Id),
		FileName: fileInfo.Name,
	})
	return &types.UserRecycleDeleteResponse{}, nil
}

func (l *UserRecycleDeleteLogic) deleteFolderContents(ctx context.Context, parentId int64, userId int64) error {
	children, err := l.svcCtx.UserRepositoryModel.FindAllDeletedByParentId(ctx, parentId, userId)
	if err != nil {
		return err
	}

	for _, child := range children {
		if child.RepositoryId == 0 {
			err = l.deleteFolderContents(ctx, int64(child.Id), userId)
			if err != nil {
				return err
			}
		} else {
			// 更新用户存储容量
			repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(ctx, &uploadServicePb.RepositoryReq{RepositoryId: int64(child.RepositoryId)})
			if err == nil {
				_, err = l.svcCtx.UserCenterRpc.DecreaseVolume(ctx, &userCenterPb.DecreaseVolumeReq{
					Id:   userId,
					Size: repositoryInfo.Size,
				})
				if err != nil {
					logx.Errorf("更新用户存储容量失败 err:%v", err)
				}
			}
		}

		err = l.svcCtx.UserRepositoryModel.Delete(ctx, child.Id)
		if err != nil {
			return err
		}
		l.svcCtx.AuditLogServiceRpc.CreateOperationLog(l.ctx, &auditservicerpc.OperationLogReq{
			UserId:  userId,
			Content: "更新用户存储容量",
			Flag:    3,
		})
	}
	return nil
}

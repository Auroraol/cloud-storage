package repository

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	uploadServicePb "github.com/Auroraol/cloud-storage/upload_service/rpc/pb"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"
	userCenterPb "github.com/Auroraol/cloud-storage/user_center/rpc/pb"
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

//
//func (l *UserFileDeleteLogic) UserFileDelete(req *types.UserFileDeleteRequest) (resp *types.UserFileDeleteResponse, err error) {
//	// 先删 user_repository (文件夹/文件)
//	userFileInfo, err := l.svcCtx.UserRepositoryModel.FindOne(l.ctx, uint64(req.Id))
//	if err != nil {
//		return nil, err
//	}
//	err = l.svcCtx.UserRepositoryModel.Delete(l.ctx, userFileInfo.Id) // 删除 user_repository记录
//	if err != nil {
//		return nil, response.NewErrMsg("更新个人存储池失败")
//	}
//
//	userId := token.GetUidFromCtx(l.ctx)
//	if userId == 0 {
//		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
//	}
//	// 文件夹
//	if userFileInfo.RepositoryId == 0 {
//		// 递归删除文件夹下的所有文件和子文件夹
//		err = l.deleteFolderContents(l.ctx, int64(userFileInfo.Id), userId)
//		if err != nil {
//			return nil, err
//		}
//		return &types.UserFileDeleteResponse{}, nil
//	}
//
//	// 从中心存储池中取 size (文件信息)
//	repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadServicePb.RepositoryReq{RepositoryId: int64(userFileInfo.RepositoryId)})
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return nil, nil
//		}
//		return nil, response.NewErrMsg("中心存储池找不到该数据")
//	}
//
//	_, err = l.svcCtx.UserCenterRpc.DecreaseVolume(l.ctx, &userCenterPb.DecreaseVolumeReq{
//		Id:   userId,
//		Size: repositoryInfo.Size,
//	})
//	if err != nil {
//		return nil, response.NewErrMsg("更新容量失败！")
//	}
//	return &types.UserFileDeleteResponse{}, nil
//}
//
//func (l *UserFileDeleteLogic) deleteFolderContents(ctx context.Context, parentId int64, userId int64) error {
//	// 获取文件夹下的所有子文件和子文件夹
//	children, err := l.svcCtx.UserRepositoryModel.FindAllFolderByParentId(ctx, parentId, userId)
//	if err != nil {
//		return err
//	}
//
//	for _, child := range children {
//		// 如果是文件夹，递归删除其内容
//		if child.RepositoryId == 0 {
//			err = l.deleteFolderContents(ctx, int64(child.Id), userId)
//			if err != nil {
//				return err
//			}
//		}
//		// 删除子文件或子文件夹
//		err = l.svcCtx.UserRepositoryModel.Delete(ctx, child.Id)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (l *UserFileDeleteLogic) UserFileDelete(req *types.UserFileDeleteRequest) (resp *types.UserFileDeleteResponse, err error) {
	// 先删 user_repository (文件夹/文件)
	userFileInfo, err := l.svcCtx.UserRepositoryModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.UserRepositoryModel.Delete(l.ctx, userFileInfo.Id) // 删除 user_repository记录
	if err != nil {
		return nil, response.NewErrMsg("更新个人存储池失败")
	}

	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 文件夹
	if userFileInfo.RepositoryId == 0 {
		// 递归删除文件夹下的所有文件和子文件夹
		err = l.deleteFolderContents(l.ctx, int64(userFileInfo.Id), userId)
		if err != nil {
			return nil, err
		}
		return &types.UserFileDeleteResponse{}, nil
	}

	// 从中心存储池中取 size (文件信息)
	repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadServicePb.RepositoryReq{RepositoryId: int64(userFileInfo.RepositoryId)})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, response.NewErrMsg("中心存储池找不到该数据")
	}

	// 更新用户存储容量
	_, err = l.svcCtx.UserCenterRpc.DecreaseVolume(l.ctx, &userCenterPb.DecreaseVolumeReq{
		Id:   userId,
		Size: repositoryInfo.Size,
	})
	if err != nil {
		return nil, response.NewErrMsg("更新容量失败！")
	}

	return &types.UserFileDeleteResponse{}, nil
}

func (l *UserFileDeleteLogic) deleteFolderContents(ctx context.Context, parentId int64, userId int64) error {
	// 获取文件夹下的所有子文件和子文件夹
	children, err := l.svcCtx.UserRepositoryModel.FindAllFolderByParentId(ctx, parentId, userId)
	if err != nil {
		return err
	}

	for _, child := range children {
		if child.RepositoryId == 0 {
			// 如果是文件夹，递归删除其内容
			err = l.deleteFolderContents(ctx, int64(child.Id), userId)
			if err != nil {
				return err
			}
		} else {
			// 如果是文件，获取文件大小并更新用户存储容量
			repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(ctx, &uploadServicePb.RepositoryReq{RepositoryId: int64(child.RepositoryId)})
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 如果文件在中心存储池中找不到，可以选择跳过或记录日志
					logx.Infof("File with RepositoryId %d not found in center storage pool", child.RepositoryId)
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
				return err
			}
		}

		// 删除子文件或子文件夹
		err = l.svcCtx.UserRepositoryModel.Delete(ctx, child.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

package repository

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	uploadServicePb "github.com/Auroraol/cloud-storage/upload_service/rpc/pb"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"
	userCenterPb "github.com/Auroraol/cloud-storage/user_center/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
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
	//先删user_repository
	userFileInfo, err := l.svcCtx.UserRepositoryModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.UserRepositoryModel.Delete(l.ctx, userFileInfo.Id)
	if err != nil {
		return nil, response.NewErrMsg("更新个人存储池失败！")
	}
	//从中心存储池中取size
	repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadServicePb.RepositoryReq{RepositoryId: int64(userFileInfo.RepositoryId)})
	if err != nil {
		return nil, response.NewErrMsg("中心存储池找不到该数据！")
	}

	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	_, err = l.svcCtx.UserCenterRpc.DecreaseVolume(l.ctx, &userCenterPb.DecreaseVolumeReq{
		Id:   userId,
		Size: repositoryInfo.Size,
	})
	if err != nil {
		return nil, response.NewErrMsg("更新容量失败！")
	}
	return &types.UserFileDeleteResponse{}, nil
}

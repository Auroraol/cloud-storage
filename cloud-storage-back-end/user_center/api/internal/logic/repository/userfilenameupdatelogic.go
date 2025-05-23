package repository

import (
	"context"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileNameUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户文件名称修改
func NewUserFileNameUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileNameUpdateLogic {
	return &UserFileNameUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileNameUpdateLogic) UserFileNameUpdate(req *types.UserFileNameUpdateRequest) (resp *types.UserFileNameUpdateResponse, err error) {
	userFileInfo, err := l.svcCtx.UserRepositoryModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		zap.S().Error("UserRepositoryModel.FindOne err:%s", err)
		return nil, err
	}
	userFileInfo.Name = req.Name
	err = l.svcCtx.UserRepositoryModel.Update(l.ctx, userFileInfo)
	if err != nil {
		zap.S().Error("UserRepositoryModel.Update err:%s", err)
		return nil, err
	}
	return &types.UserFileNameUpdateResponse{}, nil
}

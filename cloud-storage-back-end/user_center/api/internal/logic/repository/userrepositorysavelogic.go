package repository

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	uploadServicePb "github.com/Auroraol/cloud-storage/upload_service/rpc/pb"
	"github.com/Auroraol/cloud-storage/user_center/model"
	userCenterPb "github.com/Auroraol/cloud-storage/user_center/rpc/pb"
	"github.com/bwmarrin/snowflake"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRepositorySaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户文件的关联存储
func NewUserRepositorySaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRepositorySaveLogic {
	return &UserRepositorySaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRepositorySaveLogic) UserRepositorySave(req *types.UserRepositorySaveRequest) (resp *types.UserRepositorySaveResponse, err error) {
	// 判断文件是否超容量
	repositoryPoolInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadServicePb.RepositoryReq{RepositoryId: req.RepositoryId})
	if err != nil {
		zap.S().Error("调用上传服务RPC失败, err: %s", err)
		return nil, response.NewErrMsg(fmt.Sprintf("调用上传服务RPC失败, err: %v", err))
	}

	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	userInfo, err := l.svcCtx.UserCenterRpc.FindVolumeById(l.ctx, &userCenterPb.FindVolumeReq{Id: userId})
	if err != nil {
		zap.S().Error("调用用户中心RPC失败, err: %v", err)
		return nil, err
	}
	if repositoryPoolInfo.Size+userInfo.NowVolume > userInfo.TotalVolume {
		zap.S().Error("文件超出容量限制！")
		return nil, response.NewErrMsg("文件超出容量限制！")
	}
	// 更新当前容量
	_, err = l.svcCtx.UserCenterRpc.AddVolume(l.ctx, &userCenterPb.AddVolumeReq{
		Id:   userId,
		Size: repositoryPoolInfo.Size,
	})
	if err != nil {
		zap.S().Error("更新当前容量失败, err: %v", err)
		return nil, err
	}
	// 新增关联记录
	node, err := snowflake.NewNode(1)
	if err != nil {
		zap.S().Error("failed to create snowflake node: %w", err)
		return nil, fmt.Errorf("failed to create snowflake node: %w", err)
	}
	newId := node.Generate().Int64() // 生成一个新的唯一 ID
	_, err = l.svcCtx.UserRepositoryModel.Insert(l.ctx, &model.UserRepository{
		Id:           uint64(newId),
		UserId:       uint64(userId),
		ParentId:     uint64(req.ParentId),
		RepositoryId: uint64(req.RepositoryId),
		Name:         req.Name,
	})

	if err != nil {
		zap.S().Error("UserRepositoryModel.Insert err:%v", err)
		return nil, response.NewErrMsg("存储失败！")
	}

	return &types.UserRepositorySaveResponse{}, nil
}

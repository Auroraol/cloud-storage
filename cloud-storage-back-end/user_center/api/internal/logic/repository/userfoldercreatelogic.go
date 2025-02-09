package repository

import (
	"context"
	"fmt"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	"github.com/Auroraol/cloud-storage/user_center/model"
	"github.com/bwmarrin/snowflake"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFolderCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户文件夹创建
func NewUserFolderCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFolderCreateLogic {
	return &UserFolderCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFolderCreateLogic) UserFolderCreate(req *types.UserFolderCreateRequest) (resp *types.UserFolderCreateResponse, err error) {
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 验证文件夹名字是否存在
	existCount, err := l.svcCtx.UserRepositoryModel.CountByParentIdAndName(l.ctx, req.ParentId, userId, req.Name)
	if err != nil {
		return nil, response.NewErrMsg("验证文件夹名字不存在失败！")
	}
	if existCount > 0 {
		return nil, response.NewErrMsg("已存在相同名称的文件夹！")
	}
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, fmt.Errorf("failed to create snowflake node: %w", err)
	}
	newId := node.Generate().Int64() // 生成一个新的唯一 ID
	_, err = l.svcCtx.UserRepositoryModel.Insert(l.ctx, &model.UserRepository{
		Id:       uint64(newId),
		UserId:   uint64(userId),
		ParentId: uint64(req.ParentId),
		Name:     req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &types.UserFolderCreateResponse{Id: newId}, nil
}

package logic

import (
	"context"
	"fmt"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	"github.com/Auroraol/cloud-storage/share_service/model"
	"github.com/bwmarrin/snowflake"

	"github.com/Auroraol/cloud-storage/share_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/share_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareBasicCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建分享记录
func NewShareBasicCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicCreateLogic {
	return &ShareBasicCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicCreateLogic) ShareBasicCreate(req *types.ShareBasicCreateRequest) (resp *types.ShareBasicCreateResponse, err error) {
	//UserRepository, err := l.svcCtx.UserCenterRepositoryRpc.FindRepositoryIdById(l.ctx, &user.FindRepositoryIdReq{Id: req.UserRepositoryId})
	//if err != nil {
	//	logx.Errorf("failed to find repository id by id: %w", err)
	//	return nil, err
	//}
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, fmt.Errorf("failed to create snowflake node: %w", err)
	}
	newId := node.Generate().Int64()
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}
	_, err = l.svcCtx.ShareBasicModel.InsertWithId(l.ctx, &model.ShareBasic{
		Id:               uint64(newId),
		UserId:           uint64(userId),
		RepositoryId:     uint64(req.RepositoryId),
		UserRepositoryId: uint64(req.UserRepositoryId),
		ExpiredTime:      req.ExpiredTime,
	})
	if err != nil {
		logx.Errorf("failed to insert share basic: %w", err)
		return nil, err
	}
	return &types.ShareBasicCreateResponse{Id: newId}, nil
}

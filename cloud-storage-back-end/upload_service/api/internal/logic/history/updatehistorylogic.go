package history

import (
	"context"
	"strconv"

	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/token"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/types"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/model"
	"github.com/bwmarrin/snowflake"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新历史上传记录
func NewUpdateHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateHistoryLogic {
	return &UpdateHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateHistoryLogic) UpdateHistory(req *types.UpdateHistoryRequest) (resp *types.UpdateHistoryResponse, err error) {
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}
	node, err := snowflake.NewNode(1)
	if err != nil {
		zap.S().Error("创建雪花节点失败 err:%v", err)
		return nil, response.NewErrMsg("failed to create snowflake node")
	}
	newId := node.Generate().Int64()

	// 创建上传历史记录
	repositoryId := uint64(0)
	if req.RepositoryId > 0 {
		repositoryId = uint64(req.RepositoryId)
	}

	res, err := l.svcCtx.UploadHistoryModel.UpdateHistory(l.ctx, &model.UploadHistory{
		Id:           uint64(newId),
		UserId:       uint64(userId),
		Status:       req.Status,
		Size:         req.Size,
		FileName:     req.FileName,
		RepositoryId: repositoryId,
	})
	if err != nil {
		zap.S().Error("更新失败 err:%v", err)
		return nil, response.NewErrMsg("更新失败")
	}
	id, err := res.LastInsertId()
	return &types.UpdateHistoryResponse{
		Id: strconv.FormatInt(id, 10),
	}, nil
}

package history

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HistoryDeleteByIdListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除所有历史记录
func NewHistoryDeleteByIdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HistoryDeleteByIdListLogic {
	return &HistoryDeleteByIdListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HistoryDeleteByIdListLogic) HistoryDeleteByIdList(req *types.HistoryDeleteAllRequest) (resp *types.HistoryDeleteAllResponse, err error) {
	if err = l.svcCtx.UploadHistoryModel.DeleteAllByIdList(l.ctx, req.Ids); err != nil {
		zap.S().Error("删除失败 err:%v", err)
		return nil, response.NewErrCodeMsg(response.DB_UPDATE_AFFECTED_ZERO_ERROR, "删除失败")
	}

	return &types.HistoryDeleteAllResponse{}, nil
}

package history

import (
	"context"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/token"
	"go.uber.org/zap"
	"strconv"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HistoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 历史记录列表
func NewHistoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HistoryListLogic {
	return &HistoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HistoryListLogic) HistoryList(req *types.HistoryListRequest) (resp *types.HistoryListResponse, err error) {
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}
	page, err := l.svcCtx.UploadHistoryModel.FindAllInPage(l.ctx, userId, req.Page, req.Size)
	if err != nil {
		zap.S().Error("获取历史记录失败 err:%v", err)
	}
	historyList := make([]*types.History, 0)
	for _, history := range page {
		historyList = append(historyList, &types.History{
			FileName:     history.FileName,
			Size:         history.Size,
			Status:       history.Status,
			UpdateTime:   history.UpdateTime.String(),
			Id:           strconv.FormatUint(history.Id, 10),
			RepositoryId: int64(history.RepositoryId),
		})
	}

	total, err := l.svcCtx.UploadHistoryModel.CountTotalByUserIdId(l.ctx, userId)
	resp = &types.HistoryListResponse{
		HistoryList: historyList,
		Total:       total,
	}

	return resp, nil
}

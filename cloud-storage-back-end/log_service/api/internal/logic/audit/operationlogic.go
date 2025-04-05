package audit

import (
	"context"
	"fmt"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"go.uber.org/zap"
	"strconv"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获得操作日志
func NewOperationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationLogic {
	return &OperationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperationLogic) Operation(req *types.GetOperationLogReq) (resp *types.GetOperationLogRes, err error) {
	// 参数校验
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取总数
	auditModel := l.svcCtx.AuditModel
	var total int64
	total, err = auditModel.CountByTimeRange(l.ctx, int32(req.Flag), req.StartTime, req.EndTime)
	if err != nil {
		zap.S().Errorf("获取总数失败 err: %v", err)
		return nil, response.NewErrMsg("获取总数失败")
	}

	if total == 0 {
		return &types.GetOperationLogRes{
			Total:         int(total),
			OperationLogs: make([]types.OperationLog, 0),
		}, nil
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	var audits []*model.Audit
	// 按时间范围分页查询
	audits, err = auditModel.FindByTimeRange(l.ctx, offset, req.PageSize, int32(req.Flag), req.StartTime, req.EndTime)
	if err != nil {
		zap.S().Errorf("分页查询失败 err: %v", err)
		return nil, response.NewErrMsg("分页查询失败")
	}

	// 组装返回数据
	operationLogs := make([]types.OperationLog, 0)
	for _, audit := range audits {
		operationLogs = append(operationLogs, types.OperationLog{
			Content:   audit.Content,
			CreatedAt: fmt.Sprintf("%d", audit.CreateTime),
			Flag:      int(audit.Flag),
			FileId:    strconv.FormatUint(audit.RepositoryId, 10),
			FileName:  audit.FileName,
		})
	}

	return &types.GetOperationLogRes{
		Total:         int(total),
		OperationLogs: operationLogs,
	}, nil
}

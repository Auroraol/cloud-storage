package audit

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"
	"github.com/Auroraol/cloud-storage/log_service/model"

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

	//if req.StartTime > 0 && req.EndTime > 0 {
	// 按时间范围查询总数
	total, err = auditModel.CountByTimeRange(l.ctx, int32(req.Flag), req.StartTime, req.EndTime)
	//} else {
	// 普通查询总数
	//total, err = auditModel.Count(l.ctx, int(req.Flag))
	//}

	if err != nil {
		return nil, fmt.Errorf("获取总数失败: %v", err)
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

	//if req.StartTime > 0 && req.EndTime > 0 {
	// 按时间范围分页查询
	audits, err = auditModel.FindByTimeRange(l.ctx, offset, req.PageSize, int32(req.Flag), req.StartTime, req.EndTime)
	//} else {
	//	// 普通分页查询
	//	audits, err = auditModel.FindByPage(l.ctx, offset, req.PageSize, int(req.Flag))
	//}

	if err != nil {
		return nil, fmt.Errorf("分页查询失败: %v", err)
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

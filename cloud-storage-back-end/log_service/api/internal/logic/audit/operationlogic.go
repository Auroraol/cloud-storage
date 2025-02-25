package audit

import (
	"context"

	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"

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
	// 分页查询
	auditModel := l.svcCtx.AuditModel
	println(auditModel)

	return
}

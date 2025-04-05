package auditservicerpclogic

import (
	"context"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/time"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/model"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/rpc/pb"
	"go.uber.org/zap"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOperationLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOperationLogLogic {
	return &CreateOperationLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建操作记录
func (l *CreateOperationLogLogic) CreateOperationLog(in *pb.OperationLogReq) (*pb.OperationLogResp, error) {
	insert, err := l.svcCtx.AuditModel.Insert(l.ctx, &model.Audit{
		UserId:       uint64(in.UserId),
		Content:      in.Content,
		Flag:         int64(in.Flag),
		CreateTime:   time.LocalTimeNow().Unix(),
		FileName:     in.FileName,
		RepositoryId: uint64(in.FileId),
	})
	if err != nil {
		zap.S().Errorf("创建操作记录失败: %v", err)
		return nil, err
	}

	id, err := insert.LastInsertId()
	if err != nil {
		zap.S().Errorf("获取操作记录ID失败: %v", err)
		return nil, err
	}

	return &pb.OperationLogResp{
		Id: id,
	}, nil
}

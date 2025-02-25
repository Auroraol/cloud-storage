package auditservicerpclogic

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/time"
	"github.com/Auroraol/cloud-storage/log_service/model"
	"github.com/Auroraol/cloud-storage/log_service/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/rpc/pb"

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
		UserId:     uint64(in.UserId),
		Content:    in.Content,
		FileSize:   int64(in.FileSize),
		Flag:       int64(in.Flag),
		CreateTime: time.LocalTimeNow().Unix(),
	})
	if err != nil {
		return nil, err
	}

	id, err := insert.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &pb.OperationLogResp{
		Id: id,
	}, nil
}

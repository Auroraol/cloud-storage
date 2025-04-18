// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3
// Source: log_service.proto

package server

import (
	"context"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/rpc/internal/logic/auditservicerpc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/rpc/pb"
)

type AuditServiceRpcServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedAuditServiceRpcServer
}

func NewAuditServiceRpcServer(svcCtx *svc.ServiceContext) *AuditServiceRpcServer {
	return &AuditServiceRpcServer{
		svcCtx: svcCtx,
	}
}

// 创建操作记录
func (s *AuditServiceRpcServer) CreateOperationLog(ctx context.Context, in *pb.OperationLogReq) (*pb.OperationLogResp, error) {
	l := auditservicerpclogic.NewCreateOperationLogLogic(ctx, s.svcCtx)
	return l.CreateOperationLog(in)
}

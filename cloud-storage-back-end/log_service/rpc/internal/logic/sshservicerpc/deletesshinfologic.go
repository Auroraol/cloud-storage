package sshservicerpclogic

import (
	"context"

	"github.com/Auroraol/cloud-storage/log_service/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSshInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSshInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSshInfoLogic {
	return &DeleteSshInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除ssh信息
func (l *DeleteSshInfoLogic) DeleteSshInfo(in *pb.DeleteSshInfoReq) (*pb.SshInfoResp, error) {
	// 根据 ID 删除 SSH 信息
	err := l.svcCtx.SshInfoModel.Delete(l.ctx, in.SshId)
	if err != nil {
		logx.Errorf("删除 SSH 信息失败: %v", err)
		return nil, err
	}

	return &pb.SshInfoResp{Success: true, Message: "删除 SSH 信息成功", Id: in.SshId}, nil
}

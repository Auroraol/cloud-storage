package ssh

import (
	"context"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/rpc/client/sshservicerpc"
	"go.uber.org/zap"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSSHConnectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除SSH连接信息
func NewDeleteSSHConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSSHConnectLogic {
	return &DeleteSSHConnectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSSHConnectLogic) DeleteSSHConnect(req *types.DeleteSSHConnectReq) (resp *types.DeleteSSHConnectRes, err error) {
	res, err := l.svcCtx.SshServiceRpc.DeleteSshInfo(l.ctx, &sshservicerpc.DeleteSshInfoReq{
		SshId: req.SshId,
	})
	if err != nil || res.Success == false {
		zap.S().Error("保存ssh信息失败 SshServiceRpc err:%v", err)
		return nil, err
	}
	return &types.DeleteSSHConnectRes{
		Success: true,
		Message: "删除ssh信息成功",
		SshId:   req.SshId,
	}, nil
}

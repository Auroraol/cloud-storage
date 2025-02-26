package ssh

import (
	"context"
	"fmt"

	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConnectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// SSH连接
func NewConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConnectLogic {
	return &ConnectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConnectLogic) Connect(req *types.SSHConnectReq) (resp *types.SSHConnectRes, err error) {
	// 参数校验
	if req.Host == "" {
		return &types.SSHConnectRes{
			Success: false,
			Message: "主机地址不能为空",
		}, nil
	}

	if req.User == "" {
		return &types.SSHConnectRes{
			Success: false,
			Message: "用户名不能为空",
		}, nil
	}

	// 连接主机
	err = l.svcCtx.SSHService.Connect(req.Host, req.User, req.Password, req.PrivateKeyPath)
	if err != nil {
		return &types.SSHConnectRes{
			Success: false,
			Message: fmt.Sprintf("连接失败: %v", err),
		}, nil
	}

	return &types.SSHConnectRes{
		Success: true,
		Message: "连接成功",
	}, nil
}

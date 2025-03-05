package ssh

import (
	"context"
	"fmt"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	"github.com/Auroraol/cloud-storage/log_service/rpc/client/sshservicerpc"
	"go.uber.org/zap"
	"strconv"

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
	err = l.svcCtx.SSHService.Connect(req.Host, strconv.Itoa(req.Port), req.User, req.Password, req.PrivateKeyPath)
	if err != nil {
		zap.S().Errorf("连接主机失败 err: %v", err)
		return &types.SSHConnectRes{
			Success: false,
			Message: fmt.Sprintf("连接失败: %v", err),
		}, nil
	}

	// 保存ssh信息
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	res, err := l.svcCtx.SshServiceRpc.SaveSshInfo(l.ctx, &sshservicerpc.SshInfoReq{
		UserId:   userId,
		Host:     req.Host,
		Port:     int32(int64(req.Port)),
		User:     req.User,
		Password: req.Password,
	})
	if err != nil || res.Success == false {
		zap.S().Error("保存ssh信息失败 SshServiceRpc err:%v", err)
		return nil, err
	}

	return &types.SSHConnectRes{
		Success: true,
		Message: "连接成功",
	}, nil
}

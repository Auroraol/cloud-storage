package ssh

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	"github.com/Auroraol/cloud-storage/log_service/rpc/client/sshservicerpc"
	"go.uber.org/zap"

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
		zap.S().Errorf("连接主机失败 err: %s", err)
		return &types.SSHConnectRes{
			Success: false,
			Message: fmt.Sprintf("连接失败: %s", err),
		}, nil
	}

	// 获取用户ID
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 开启协程保存SSH信息
	go func() {
		// 创建新的上下文，避免使用已经完成的请求上下文
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 保存SSH连接信息
		sshReq := &sshservicerpc.SshInfoReq{
			UserId:   userId,
			Host:     req.Host,
			Port:     int32(req.Port),
			User:     req.User,
			Password: req.Password,
		}

		// 如果记录已存在会自动失败，无需显式检查
		_, saveErr := l.svcCtx.SshServiceRpc.SaveSshInfo(ctx, sshReq)
		if saveErr != nil {
			// 记录错误但不影响主流程，保存失败可能是因为记录已存在
			zap.S().Infof("SSH信息保存失败，可能是记录已存在: %v", saveErr)
		}
	}()

	return &types.SSHConnectRes{
		Success: true,
		Message: "连接成功",
	}, nil
}

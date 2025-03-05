package svc

import (
	"github.com/Auroraol/cloud-storage/common/logx"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/config"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/service"
	"github.com/Auroraol/cloud-storage/log_service/model"
	"github.com/Auroraol/cloud-storage/log_service/rpc/client/sshservicerpc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config        config.Config
	AuditModel    model.AuditModel
	LogfileModel  model.LogfileModel
	SSHService    service.SSHService
	SshServiceRpc sshservicerpc.SshServiceRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Options.Dsn)
	c.LogConfig.CustomLevels = map[string]string{
		"requests": "info", // 自定义业务日志级别，对应info级别
	}
	if err := logx.InitLogger(c.LogConfig); err != nil {
		zap.S().Errorf("日志初始化失败: %v", err)
		panic(err)
	}
	return &ServiceContext{
		Config:       c,
		AuditModel:   model.NewAuditModel(conn, c.CacheRedis),
		LogfileModel: model.NewLogfileModel(conn, c.CacheRedis),
		SSHService:   service.NewSSHService(),
		SshServiceRpc: sshservicerpc.NewSshServiceRpc(
			zrpc.MustNewClient(c.SshServiceRpcConf)),
	}
}

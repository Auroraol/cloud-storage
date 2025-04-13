package svc

import (
	"path/filepath"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/logx"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/config"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/service"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/localfile"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/model"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/rpc/client/sshservicerpc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config           config.Config
	AuditModel       model.AuditModel
	LogfileModel     model.LogfileModel
	SSHService       service.SSHService
	SshServiceRpc    sshservicerpc.SshServiceRpc
	LocalFileService *localfile.LocalFileService
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

	// 为本地日志服务设置基础路径
	basePath := c.LogConfig.LogPath
	if !filepath.IsAbs(basePath) {
		// 如果不是绝对路径，将其转换为绝对路径
		absPath, err := filepath.Abs(basePath)
		if err == nil {
			basePath = absPath
		}
	}

	// 创建本地文件服务
	localFileService, err := localfile.NewLocalFileService(basePath)
	if err != nil {
		zap.S().Errorf("创建本地文件服务失败: %v", err)
		panic(err)
	}

	return &ServiceContext{
		Config:           c,
		AuditModel:       model.NewAuditModel(conn, c.CacheRedis),
		LogfileModel:     model.NewLogfileModel(conn, c.CacheRedis),
		SSHService:       service.NewSSHService(),
		SshServiceRpc:    sshservicerpc.NewSshServiceRpc(zrpc.MustNewClient(c.SshServiceRpcConf)),
		LocalFileService: localFileService,
	}
}

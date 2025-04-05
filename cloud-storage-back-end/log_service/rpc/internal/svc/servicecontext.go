package svc

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/logx"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/model"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config       config.Config
	AuditModel   model.AuditModel
	SshInfoModel model.SshInfoModel
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
		SshInfoModel: model.NewSshInfoModel(conn, c.CacheRedis),
	}
}

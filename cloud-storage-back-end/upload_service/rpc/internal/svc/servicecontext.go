package svc

import (
	"github.com/Auroraol/cloud-storage/common/logx"
	"github.com/Auroraol/cloud-storage/upload_service/model"
	"github.com/Auroraol/cloud-storage/upload_service/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config              config.Config
	RepositoryPoolModel model.RepositoryPoolModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Options.Dsn)

	// 初始化日志
	c.Log.CustomLevels = map[string]string{
		"requests": "info", // 自定义业务日志级别，对应info级别
		"audit":    "warn", // 自定义审计日志级别，对应warn级别
	}
	if err := logx.InitLogger(c.Log); err != nil {
		zap.S().Errorf("日志初始化失败: %v", err)
		panic(err)
	}

	return &ServiceContext{
		Config:              c,
		RepositoryPoolModel: model.NewRepositoryPoolModel(conn, c.CacheRedis),
	}
}

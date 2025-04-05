package svc

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/logx"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/model"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config              config.Config
	UserModel           model.UserModel
	UserRepositoryModel model.UserRepositoryModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Options.Dsn) // 原生连接

	// 初始化日志
	c.LogConfig.CustomLevels = map[string]string{
		"requests": "info", // 自定义业务日志级别，对应info级别
		"audit":    "warn", // 自定义审计日志级别，对应warn级别
	}
	if err := logx.InitLogger(c.LogConfig); err != nil {
		zap.S().Errorf("日志初始化失败: %s", err)
		panic(err)
	}

	return &ServiceContext{
		Config:              c,
		UserModel:           model.NewUserModel(sqlConn),
		UserRepositoryModel: model.NewUserRepositoryModel(sqlConn, c.CacheRedis),
	}
}

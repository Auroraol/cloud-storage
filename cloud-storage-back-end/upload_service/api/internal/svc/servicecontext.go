package svc

import (
	"github.com/Auroraol/cloud-storage/common/logx"
	"github.com/Auroraol/cloud-storage/common/mq"
	"github.com/Auroraol/cloud-storage/log_service/rpc/client/auditservicerpc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/config"
	"github.com/Auroraol/cloud-storage/upload_service/model"
	"github.com/Auroraol/cloud-storage/user_center/rpc/client/userservicerpc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config              config.Config
	UserCenterRpc       userservicerpc.UserServiceRpc
	RepositoryPoolModel model.RepositoryPoolModel
	//RedisClient         *redis.Redis
	UploadHistoryModel model.UploadHistoryModel
	AuditLogServiceRpc auditservicerpc.AuditServiceRpc
	PulsarManager      *mq.PulsarManager // Pulsar消息队列管理器
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Options.Dsn)

	// 初始化日志
	c.LogConfig.CustomLevels = map[string]string{
		"requests": "info", // 自定义业务日志级别，对应info级别
	}
	if err := logx.InitLogger(c.LogConfig); err != nil {
		zap.S().Errorf("日志初始化失败: %s", err)
		panic(err)
	}

	// 初始化Pulsar管理器
	var pulsarManager *mq.PulsarManager
	if c.Pulsar.Enabled {
		var err error
		pulsarManager, err = mq.NewPulsarManager(mq.PulsarConfig{
			URL: c.Pulsar.URL,
		})
		if err != nil {
			zap.S().Errorf("Pulsar管理器初始化失败: %s", err)
		} else {
			zap.S().Info("Pulsar管理器初始化成功")
		}
	}

	return &ServiceContext{
		Config:              c,
		UserCenterRpc:       userservicerpc.NewUserServiceRpc(zrpc.MustNewClient(c.UserCenterRpcConf)),
		RepositoryPoolModel: model.NewRepositoryPoolModel(conn, c.CacheRedis),
		UploadHistoryModel:  model.NewUploadHistoryModel(conn, c.CacheRedis),
		AuditLogServiceRpc:  auditservicerpc.NewAuditServiceRpc(zrpc.MustNewClient(c.AuditServiceRpcConf)),
		PulsarManager:       pulsarManager,
		//RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
		//	r.Type = c.Redis.Type
		//	r.Pass = c.Redis.Pass
		//}),
	}
}

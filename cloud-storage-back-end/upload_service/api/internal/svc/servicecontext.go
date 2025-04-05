package svc

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/logx"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/rpc/client/auditservicerpc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/config"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/model"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/client/userservicerpc"
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
	//PulsarManager      *pulsar.PulsarManager
	//FilePublisher      *pulsar.Publisher
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

	//// 初始化 Pulsar 管理器和发布者
	//var pulsarManager *pulsar.PulsarManager
	//var filePublisher *pulsar.Publisher
	//
	//if c.PubConfig.Enabled {
	//	var err error
	//	pulsarManager, err = pulsar.NewPulsarManager(pulsar.Config{
	//		URL: c.PubConfig.URL,
	//	})
	//	if err != nil {
	//		zap.S().Errorf("Pulsar 管理器初始化失败: %s", err)
	//	} else {
	//		// 创建文件上传消息的发布者
	//		filePublisher, err = pulsar.NewPublisher(pulsarManager, pulsar.PubConfig{
	//			Topic:           "file-uploaded",
	//			BatchingEnabled: true,
	//		})
	//		if err != nil {
	//			zap.S().Errorf("Pulsar 发布者初始化失败: %s", err)
	//		}
	//	}
	//}

	return &ServiceContext{
		Config:              c,
		UserCenterRpc:       userservicerpc.NewUserServiceRpc(zrpc.MustNewClient(c.UserCenterRpcConf)),
		RepositoryPoolModel: model.NewRepositoryPoolModel(conn, c.CacheRedis),
		UploadHistoryModel:  model.NewUploadHistoryModel(conn, c.CacheRedis),
		AuditLogServiceRpc:  auditservicerpc.NewAuditServiceRpc(zrpc.MustNewClient(c.AuditServiceRpcConf)),
		//PulsarManager:       pulsarManager,
		//FilePublisher:       filePublisher,
	}
}

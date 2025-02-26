package svc

import (
	"github.com/Auroraol/cloud-storage/log_service/rpc/client/auditservicerpc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/config"
	"github.com/Auroraol/cloud-storage/upload_service/model"
	"github.com/Auroraol/cloud-storage/user_center/rpc/client/userservicerpc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	UserCenterRpc       userservicerpc.UserServiceRpc
	RepositoryPoolModel model.RepositoryPoolModel
	//RedisClient         *redis.Redis
	UploadHistoryModel model.UploadHistoryModel
	AuditLogServiceRpc auditservicerpc.AuditServiceRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Options.Dsn)
	return &ServiceContext{
		Config:              c,
		UserCenterRpc:       userservicerpc.NewUserServiceRpc(zrpc.MustNewClient(c.UserCenterRpcConf)),
		RepositoryPoolModel: model.NewRepositoryPoolModel(conn, c.CacheRedis),
		UploadHistoryModel:  model.NewUploadHistoryModel(conn, c.CacheRedis),
		AuditLogServiceRpc:  auditservicerpc.NewAuditServiceRpc(zrpc.MustNewClient(c.AuditServiceRpcConf)),
		//RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
		//	r.Type = c.Redis.Type
		//	r.Pass = c.Redis.Pass
		//}),
	}
}

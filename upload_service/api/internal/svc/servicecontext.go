package svc

import (
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/config"
	"github.com/Auroraol/cloud-storage/upload_service/model"
	"github.com/Auroraol/cloud-storage/user_center/rpc/usercenter"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	UserCenterRpc       usercenter.UserCenter
	RepositoryPoolModel model.RepositoryPoolModel
	//RedisClient         *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Options.Dsn)
	return &ServiceContext{
		Config: c,
		UserCenterRpc: usercenter.NewUserCenter(
			zrpc.MustNewClient(c.UserCenterRpcConf)),
		RepositoryPoolModel: model.NewRepositoryPoolModel(conn, c.CacheRedis),
		//RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
		//	r.Type = c.Redis.Type
		//	r.Pass = c.Redis.Pass
		//}),
	}
}

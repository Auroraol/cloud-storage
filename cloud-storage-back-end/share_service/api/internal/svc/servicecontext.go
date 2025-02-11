package svc

import (
	"github.com/Auroraol/cloud-storage/common/orm"
	"github.com/Auroraol/cloud-storage/share_service/api/internal/config"
	"github.com/Auroraol/cloud-storage/share_service/model"
	"github.com/Auroraol/cloud-storage/upload_service/rpc/uploadservice"
	"github.com/Auroraol/cloud-storage/user_center/rpc/client/userrepository"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config                  config.Config
	ShareBasicModel         model.ShareBasicModel
	UserCenterRepositoryRpc userrepository.UserRepository // 用户中心储存
	UploadServiceRpc        uploadservice.UploadService   // 上传服务
	RedisClient             *redis.Redis                  // redis
	Engine                  *gorm.DB                      // orm 联表查询
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Options.Dsn)
	client, _ := orm.NewClient(c.Options)
	return &ServiceContext{
		Config:          c,
		ShareBasicModel: model.NewShareBasicModel(conn, c.CacheRedis),
		Engine:          client.GetGormDB(),
		UserCenterRepositoryRpc: userrepository.NewUserRepository(
			zrpc.MustNewClient(c.UserCenterRpcConf)),
		UploadServiceRpc: uploadservice.NewUploadService(
			zrpc.MustNewClient(c.UploadServiceRpcConf)),
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
	}
}

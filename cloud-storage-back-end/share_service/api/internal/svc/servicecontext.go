package svc

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/logx"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/orm"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/share_service/api/internal/config"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/share_service/model"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/rpc/client/uploadservicerpc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/client/userrepositoryrpc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config                  config.Config
	ShareBasicModel         model.ShareBasicModel
	UserCenterRepositoryRpc userrepositoryrpc.UserRepositoryRpc // 用户中心储存
	UploadServiceRpc        uploadservicerpc.UploadServiceRpc   // 上传服务
	RedisClient             *redis.Redis                        // redis
	Engine                  *gorm.DB                            // orm 联表查询
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Options.Dsn)
	client, _ := orm.NewClient(c.Options)
	c.LogConfig.CustomLevels = map[string]string{
		"requests": "info", // 自定义业务日志级别，对应info级别
	}
	if err := logx.InitLogger(c.LogConfig); err != nil {
		zap.S().Errorf("日志初始化失败: %s", err)
		panic(err)
	}
	return &ServiceContext{
		Config:          c,
		ShareBasicModel: model.NewShareBasicModel(conn, c.CacheRedis),
		Engine:          client.GetGormDB(),
		UserCenterRepositoryRpc: userrepositoryrpc.NewUserRepositoryRpc(
			zrpc.MustNewClient(c.UserCenterRpcConf)),
		UploadServiceRpc: uploadservicerpc.NewUploadServiceRpc(
			zrpc.MustNewClient(c.UploadServiceRpcConf)),
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
	}
}

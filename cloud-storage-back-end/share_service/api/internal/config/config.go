package config

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/logx"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/orm"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	Options              orm.Options
	CacheRedis           cache.CacheConf // 缓存
	UserCenterRpcConf    zrpc.RpcClientConf
	UploadServiceRpcConf zrpc.RpcClientConf
	Redis                redis.RedisConf
	LogConfig            logx.LogConfig
}

package config

import (
	"github.com/Auroraol/cloud-storage/common/orm"
	"github.com/zeromicro/go-zero/core/stores/cache"
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
	CacheRedis           cache.CacheConf
	UploadServiceRpcConf zrpc.RpcClientConf
	UserCenterRpcConf    zrpc.RpcClientConf
	AuditServiceRpcConf  zrpc.RpcClientConf
}

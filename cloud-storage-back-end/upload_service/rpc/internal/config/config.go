package config

import (
	"github.com/Auroraol/cloud-storage/common/logx"
	"github.com/Auroraol/cloud-storage/common/orm"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Options    orm.Options
	CacheRedis cache.CacheConf // 缓存
	Log        logx.LogConfig  // 日志配置
}

package config

import (
	"github.com/Auroraol/cloud-storage/common/logx"
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
	UserCenterRpcConf   zrpc.RpcClientConf
	AuditServiceRpcConf zrpc.RpcClientConf
	Options             orm.Options
	CacheRedis          cache.CacheConf // 缓存
	LogConfig           logx.LogConfig  // 日志配置
	Pulsar              struct {
		Enabled     bool   // 是否启用Pulsar
		URL         string // Pulsar服务地址
		ServiceName string // 服务名称
	}
}

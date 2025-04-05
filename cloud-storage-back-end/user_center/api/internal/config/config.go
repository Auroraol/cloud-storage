package config

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/logx"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/orm"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/sms"
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
	LogConfig            logx.LogConfig
	Sms                  sms.SmsConfig
}

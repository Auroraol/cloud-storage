package config

import (
	"github.com/Auroraol/cloud-storage/common/orm"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Options orm.Options
}

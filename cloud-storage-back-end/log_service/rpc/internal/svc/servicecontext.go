package svc

import (
	"github.com/Auroraol/cloud-storage/log_service/model"
	"github.com/Auroraol/cloud-storage/log_service/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	AuditModel model.AuditModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Options.Dsn)
	return &ServiceContext{
		Config: c, AuditModel: model.NewAuditModel(conn, c.CacheRedis),
	}
}

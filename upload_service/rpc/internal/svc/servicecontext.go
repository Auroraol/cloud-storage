package svc

import (
	"github.com/Auroraol/cloud-storage/upload_service/model"
	"github.com/Auroraol/cloud-storage/upload_service/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config              config.Config
	RepositoryPoolModel model.RepositoryPoolModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Options.Dsn)
	return &ServiceContext{
		Config:              c,
		RepositoryPoolModel: model.NewRepositoryPoolModel(conn, c.CacheRedis),
	}
}

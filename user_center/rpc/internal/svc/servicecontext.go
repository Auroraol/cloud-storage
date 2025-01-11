package svc

import (
	"github.com/Auroraol/cloud-storage/user_center/model"
	"github.com/Auroraol/cloud-storage/user_center/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Options.Dsn) // 原生连接
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(sqlConn),
	}
}

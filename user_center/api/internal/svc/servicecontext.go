package svc

import (
	"github.com/Auroraol/cloud-storage/common/orm"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/config"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	DB       *gorm.DB
	DbClient orm.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, err := orm.NewClient(c.Options)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:   c,
		DbClient: client,
	}
}

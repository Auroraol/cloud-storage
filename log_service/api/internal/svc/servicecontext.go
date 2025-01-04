package svc

import (
	"cloud-storage/log_service/api/internal/config"
	"cloud-storage/log_service/api/internal/db"
	"cloud-storage/log_service/model"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {

	db := db.InitGorm(c.Mysql.DataSource) // init自定义的
	db.AutoMigrate(&model.LogfileModel{})
	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}

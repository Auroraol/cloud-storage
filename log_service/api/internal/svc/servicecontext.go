package svc

import (
	"cloud-storage/log_service/api/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {

	//db := init_gorm.InitGorm(c.Mysql.DataSource) // init自定义的
	//db.AutoMigrate(&models.UserModel{})
	return &ServiceContext{
		Config: c,
		//DB:     db,
	}
}

package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitGorm gorm初始化
func InitGorm(MysqlDataSource string) *gorm.DB {
	//ctx := context.Background()
	db, err := gorm.Open(mysql.Open(MysqlDataSource), &gorm.Config{})
	if err != nil {
		//logger.Errorf(ctx, "gorm.Open err : %+v", err)
	}
	return db
}

package orm

import (
	"context"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"time"
	//"time"
)

// 接口
type Client interface {
	GetGormDB(ctx context.Context) *gorm.DB
}

// 实现
type clientImpl struct {
	db *gorm.DB
}

func NewClient(options Options) (Client, error) {
	config := &gorm.Config{
		SkipDefaultTransaction:                   options.SkipDefaultTransaction,
		FullSaveAssociations:                     options.FullSaveAssociations,
		DryRun:                                   options.DryRun,
		PrepareStmt:                              options.PrepareStmt,
		DisableAutomaticPing:                     options.DisableAutomaticPing,
		DisableForeignKeyConstraintWhenMigrating: options.DisableForeignKeyConstraintWhenMigrating,
		DisableNestedTransaction:                 options.DisableNestedTransaction,
		AllowGlobalUpdate:                        options.AllowGlobalUpdate,
		QueryFields:                              options.QueryFields,
		CreateBatchSize:                          options.CreateBatchSize,
	}

	DB, err := gorm.Open(GetDialector(options), config)
	if err != nil {
		return nil, err
	}
	//设置连接池参数
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, err
	}

	if options.MaxOpenConns == nil {
		defaultMaxOpenConns := DefaultMaxOpenConns
		options.MaxOpenConns = &defaultMaxOpenConns
	}
	sqlDB.SetMaxOpenConns(*options.MaxOpenConns)

	if options.MaxIdleConns != nil {
		sqlDB.SetMaxIdleConns(*options.MaxIdleConns)
	}

	if options.ConnMaxIdleTime != nil {
		sqlDB.SetConnMaxIdleTime(time.Duration(*options.ConnMaxIdleTime) * time.Second)
	}

	if options.ConnMaxLifetime != nil {
		sqlDB.SetConnMaxLifetime(time.Duration(*options.ConnMaxLifetime) * time.Second)
	}
	impl := &clientImpl{
		db: DB,
	}
	return impl, nil
}

func GetDialector(option Options) gorm.Dialector {
	var dialector gorm.Dialector
	switch option.DBType {
	case MySQL:
		dialector = mysql.Open(option.Dsn)
	case Postgres:
		dialector = postgres.Open(option.Dsn)
	case SQLite:
		dialector = sqlite.Open(option.Dsn)
	case ClickHouse:
		dialector = clickhouse.Open(option.Dsn)
	case SQLServer:
		dialector = sqlserver.Open(option.Dsn)
	default:
		dialector = postgres.Open(option.Dsn)
	}
	return dialector
}

func (c *clientImpl) GetGormDB(ctx context.Context) *gorm.DB {
	return c.db
}

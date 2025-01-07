package orm

type Config struct {
	Options Options `mapstructure:",squash"`
}

type Options struct {
	Dsn                                      string
	DBType                                   string
	DisableAutoInit                          bool `json:",optional"`
	SkipDefaultTransaction                   bool `json:",optional"`
	FullSaveAssociations                     bool `json:",optional"`
	PrepareStmt                              bool `json:",optional"`
	DisableAutomaticPing                     bool `json:",optional"`
	DisableForeignKeyConstraintWhenMigrating bool `json:",optional"`
	DryRun                                   bool `json:",optional"`
	DisableNestedTransaction                 bool `json:",optional"`
	AllowGlobalUpdate                        bool `json:",optional"`
	QueryFields                              bool `json:",optional"`
	CreateBatchSize                          int  `json:",optional"`
	// 连接池参数
	MaxIdleConns    *int `json:",optional"`
	MaxOpenConns    *int `json:",optional"`
	ConnMaxIdleTime *int `json:",optional"`
	ConnMaxLifetime *int `json:",optional"`
}

package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LogfileModel = (*customLogfileModel)(nil)

type (
	// LogfileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLogfileModel.
	LogfileModel interface {
		logfileModel
	}

	customLogfileModel struct {
		*defaultLogfileModel
	}
)

// NewLogfileModel returns a model for the database table.
func NewLogfileModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LogfileModel {
	return &customLogfileModel{
		defaultLogfileModel: newLogfileModel(conn, c, opts...),
	}
}

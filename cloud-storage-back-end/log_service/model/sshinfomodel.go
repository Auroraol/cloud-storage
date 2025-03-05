package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SshInfoModel = (*customSshInfoModel)(nil)

type (
	// SshInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSshInfoModel.
	SshInfoModel interface {
		sshInfoModel
		FindAll(ctx context.Context) ([]*SshInfo, error) // 新增 FindAll 方法
	}

	customSshInfoModel struct {
		*defaultSshInfoModel
	}
)

// NewSshInfoModel returns a model for the database table.
func NewSshInfoModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SshInfoModel {
	return &customSshInfoModel{
		defaultSshInfoModel: newSshInfoModel(conn, c, opts...),
	}
}

func (m *defaultSshInfoModel) FindAll(ctx context.Context) ([]*SshInfo, error) {
	query := "select id, user_id, host, port, username, password, created_at, updated_at from " + m.table
	var resp []*SshInfo
	err := m.QueryRowNoCacheCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

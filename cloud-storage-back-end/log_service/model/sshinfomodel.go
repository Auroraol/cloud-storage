package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SshInfoModel = (*customSshInfoModel)(nil)

type (
	// SshInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSshInfoModel.
	SshInfoModel interface {
		sshInfoModel
		FindAll(ctx context.Context, userId string) ([]*SshInfo, error) // 新增 FindAll 方法
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

// 查询所有
func (m *defaultSshInfoModel) FindAll(ctx context.Context, userId string) ([]*SshInfo, error) {
	query := "SELECT * FROM " + m.table + " WHERE user_id = ?" // 使用参数化占位符
	var resp []*SshInfo

	// 使用 QueryRows 处理多行结果
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId) // 注意方法名改为 QueryRows...
	if err != nil {
		return nil, fmt.Errorf("FindAll failed: %w", err) // 增强错误信息
	}
	return resp, nil
}

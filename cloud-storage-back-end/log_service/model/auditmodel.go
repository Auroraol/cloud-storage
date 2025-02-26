package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AuditModel = (*customAuditModel)(nil)

type (
	// AuditModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuditModel.
	AuditModel interface {
		auditModel
		Count(ctx context.Context, flag int) (int64, error)
		FindByPage(ctx context.Context, offset, limit int, flag int) ([]*Audit, error)
		CountByTimeRange(ctx context.Context, flag int32, startTime, endTime int64) (int64, error)
		FindByTimeRange(ctx context.Context, offset, limit int, flag int32, startTime, endTime int64) ([]*Audit, error)
	}

	customAuditModel struct {
		*defaultAuditModel
	}
)

// NewAuditModel returns a model for the database table.
func NewAuditModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AuditModel {
	return &customAuditModel{
		defaultAuditModel: newAuditModel(conn, c, opts...),
	}
}

// Count 获取总数
func (m *customAuditModel) Count(ctx context.Context, flag int) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", m.table)
	if flag >= 0 {
		query += " WHERE `flag` = ?"
		err := m.QueryRowNoCacheCtx(ctx, &count, query, flag)
		if err != nil {
			return 0, err
		}
	} else {
		err := m.QueryRowNoCacheCtx(ctx, &count, query)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

// FindByPage 分页查询
func (m *customAuditModel) FindByPage(ctx context.Context, offset, limit int, flag int) ([]*Audit, error) {
	var resp []*Audit
	query := fmt.Sprintf("SELECT %s FROM %s", auditRows, m.table)
	if flag >= 0 {
		query += " WHERE `flag` = ?"
	}
	query += " ORDER BY id DESC LIMIT ?,?"

	if flag >= 0 {
		err := m.QueryRowsNoCacheCtx(ctx, &resp, query, flag, offset, limit)
		if err != nil {
			return nil, err
		}
	} else {
		err := m.QueryRowsNoCacheCtx(ctx, &resp, query, offset, limit)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

// CountByTimeRange 按时间范围统计数量
func (m *customAuditModel) CountByTimeRange(ctx context.Context, flag int32, startTime, endTime int64) (int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", m.table)
	args := []interface{}{}

	if flag >= 0 {
		query += " WHERE flag = ?"
		args = append(args, flag)
	}

	if startTime > 0 && endTime > 0 {
		if len(args) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " create_time >= ? AND create_time < ?"
		args = append(args, startTime, endTime)
	}

	err := m.QueryRowNoCacheCtx(ctx, &count, query, args...)
	return count, err
}

// FindByTimeRange 按时间范围分页查询
func (m *customAuditModel) FindByTimeRange(ctx context.Context, offset, limit int, flag int32, startTime, endTime int64) ([]*Audit, error) {
	var resp []*Audit
	query := fmt.Sprintf("SELECT %s FROM %s", auditRows, m.table)
	args := []interface{}{}

	if flag >= 0 {
		query += " WHERE flag = ?"
		args = append(args, flag)
	}

	if startTime > 0 && endTime > 0 {
		if len(args) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " create_time >= ? AND create_time < ?"
		args = append(args, startTime, endTime)
	}

	query += " ORDER BY create_time DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

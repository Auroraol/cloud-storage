package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ShareBasicModel = (*customShareBasicModel)(nil)

type (
	// ShareBasicModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShareBasicModel.
	ShareBasicModel interface {
		shareBasicModel
		InsertWithId(ctx context.Context, data *ShareBasic) (sql.Result, error)
		AddOneClick(ctx context.Context, id int64) error
		FindOneByIdentity(ctx context.Context, identity uint64) (*ShareBasic, error)
	}

	customShareBasicModel struct {
		*defaultShareBasicModel
	}
)

// NewShareBasicModel returns a model for the database table.
func NewShareBasicModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ShareBasicModel {
	return &customShareBasicModel{
		defaultShareBasicModel: newShareBasicModel(conn, c, opts...),
	}
}

func (m *defaultShareBasicModel) InsertWithId(ctx context.Context, data *ShareBasic) (sql.Result, error) {
	shareBasicIdKey := fmt.Sprintf("%s%s", cacheShareBasicIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (id, user_id, repository_id, user_repository_id, expired_time, click_num, code) values (?, ?, ?, ?, ?, ?, ?)", m.table)
		return conn.ExecCtx(ctx, query, data.Id, data.UserId, data.RepositoryId, data.UserRepositoryId, data.ExpiredTime, data.ClickNum, data.Code)
	}, shareBasicIdKey)

	return ret, err
}

func (m *defaultShareBasicModel) AddOneClick(ctx context.Context, id int64) error {
	shareBasicIdKey := fmt.Sprintf("%s%s", cacheShareBasicIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set click_num = click_num + 1 where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, shareBasicIdKey)
	return err
}

func (m *defaultShareBasicModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(shareBasicRows).From(m.table)
}

func (m *defaultShareBasicModel) FindOneByIdentity(ctx context.Context, repositoryId uint64) (*ShareBasic, error) {
	var resp ShareBasic
	rowBuilder := m.RowBuilder().Where("repository_id = ?", repositoryId)
	query, values, err := rowBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	if err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...); err != nil {
		if errors.Is(err, sqlc.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &resp, nil
}

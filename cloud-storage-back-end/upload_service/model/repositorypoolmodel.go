package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RepositoryPoolModel = (*customRepositoryPoolModel)(nil)

type (
	// RepositoryPoolModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRepositoryPoolModel.
	RepositoryPoolModel interface {
		repositoryPoolModel
		InsertWithId(ctx context.Context, data *RepositoryPool) (sql.Result, error)
		CountByHash(ctx context.Context, hash string) (int64, error)
		FindRepositoryPoolByHash(ctx context.Context, hash string) (*RepositoryPool, error)
	}

	customRepositoryPoolModel struct {
		*defaultRepositoryPoolModel
	}
)

// NewRepositoryPoolModel returns a model for the database table.
func NewRepositoryPoolModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RepositoryPoolModel {
	return &customRepositoryPoolModel{
		defaultRepositoryPoolModel: newRepositoryPoolModel(conn, c, opts...),
	}
}

func (m *defaultRepositoryPoolModel) InsertWithId(ctx context.Context, data *RepositoryPool) (sql.Result, error) {
	repositoryPoolHashKey := fmt.Sprintf("%s%v", cacheRepositoryPoolHashPrefix, data.Hash)
	repositoryPoolIdKey := fmt.Sprintf("%s%v", cacheRepositoryPoolIdPrefix, data.Id)

	// 使用 squirrel 构建插入查询，添加 id 字段
	query, args, err := squirrel.Insert(m.table).
		Columns("id", "identity", "hash", "ext", "size", "path", "name").
		Values(data.Id, data.Identity, data.Hash, data.Ext, data.Size, data.Path, data.Name).
		ToSql()
	if err != nil {
		return nil, err
	}

	// 执行插入操作并更新缓存
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, repositoryPoolHashKey, repositoryPoolIdKey)

	return ret, err
}

func (m *defaultRepositoryPoolModel) CountByHash(ctx context.Context, hash string) (int64, error) {
	countBuilder := m.CountBuilder("id")
	query, values, err := countBuilder.Where("hash = ?", hash).ToSql()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

func (m *defaultRepositoryPoolModel) FindRepositoryPoolByHash(ctx context.Context, hash string) (*RepositoryPool, error) {
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.Where("hash = ?", hash).ToSql()
	if err != nil {
		return nil, err
	}
	var resp RepositoryPool
	if err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...); err != nil {
		if errors.Is(err, sqlc.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &resp, nil
}

// 构建查询 repository_pool 表的 SELECT 语句
func (m *defaultRepositoryPoolModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(repositoryPoolRows).From(m.table)
}

// 生成类似 "SELECT count(field) FROM table" 的 SQL 语句
func (m *defaultRepositoryPoolModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("count(" + field + ")").From(m.table)
}

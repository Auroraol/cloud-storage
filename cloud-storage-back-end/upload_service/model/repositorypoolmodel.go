package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/jsonx"

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
		FindOneByIdentity(ctx context.Context, identity uint64) (*RepositoryPool, error)
		DeleteByIdentity(ctx context.Context, identity uint64) error
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

// InsertWithId 插入数据并返回结果(缓存hash和文件信息) 默认不过期
func (m *defaultRepositoryPoolModel) InsertWithId(ctx context.Context, data *RepositoryPool) (sql.Result, error) {
	// ，ExecCtx 方法用于执行插入操作，并且设置了两个缓存键：repositoryPoolIdKey 和 repositoryPoolHashKey。
	repositoryPoolHashKey := fmt.Sprintf("%s%v", cacheRepositoryPoolHashPrefix, data.Hash)
	repositoryPoolIdKey := fmt.Sprintf("%s%v", cacheRepositoryPoolIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, "`identity`,`hash`,`ext`,`size`,`path`,`name`, `oss_key`")
		return conn.ExecCtx(ctx, query, data.Identity, data.Hash, data.Ext, data.Size, data.Path, data.Name, data.OssKey)
	}, repositoryPoolIdKey, repositoryPoolHashKey)

	switch {
	case err == nil:
		return ret, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultRepositoryPoolModel) CountByHash(ctx context.Context, hash string) (int64, error) {
	countBuilder := m.CountBuilder("id")
	query, values, err := countBuilder.Where("hash = ?", hash).ToSql()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
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

func (m *defaultRepositoryPoolModel) FindOneByIdentity(ctx context.Context, identity uint64) (*RepositoryPool, error) {
	repositoryPoolIdKey := fmt.Sprintf("%s%v", cacheRepositoryPoolIdPrefix, identity)
	var resp RepositoryPool
	// 先调用 m.QueryRowCtx, 从缓存中获取不到数据, 再会执行conn.QueryRowCtx(ctx, v, query, identity)
	err := m.QueryRowCtx(ctx, &resp, repositoryPoolIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `identity` = ? limit 1", repositoryPoolRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, identity)
	})
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
	// 不需要
	//// 如果缓存中没有数据，则查询数据库
	//if errors.Is(err, sqlc.ErrNotFound) {
	//	query := fmt.Sprintf("select %s from %s where `identity` = ? limit 1", repositoryPoolRows, m.table)
	//	if err := m.QueryRowNoCacheCtx(ctx, &resp, query, identity); err != nil {
	//		if errors.Is(err, sqlc.ErrNotFound) {
	//			return nil, ErrNotFound
	//		}
	//		return nil, err
	//	}
	//	// 将查询结果存入缓存
	//	if err := m.SetWithExpireCtx(ctx, repositoryPoolIdKey, resp, 0); err != nil {
	//		return nil, err
	//	}
	//	return &resp, nil
	//}
}

func (m *defaultRepositoryPoolModel) DeleteByIdentity(ctx context.Context, identity uint64) error {
	data, err := m.FindOneByIdentity(ctx, identity)
	if err != nil {
		return err
	}

	repositoryPoolHashKey := fmt.Sprintf("%s%v", cacheRepositoryPoolHashPrefix, data.Hash)
	repositoryPoolIdKey := fmt.Sprintf("%s%v", cacheRepositoryPoolIdPrefix, identity)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `identity` = ?", m.table)
		return conn.ExecCtx(ctx, query, identity)
	}, repositoryPoolHashKey, repositoryPoolIdKey)
	return err
}

// expiree 参数表示缓存的过期时间，0:短期缓存：2分钟; 1:中期缓存：6小时，适用于不常更新但对实时性要求不高的数据 其他:长期缓存：12小时到24小时，适用于几乎不变的数据。
func (m *defaultRepositoryPoolModel) SetWithExpireCtx(ctx context.Context, key string, val any,
	expire int) error {
	data, err := jsonx.Marshal(val)
	if err != nil {
		return err
	}
	// 根据业务需求动态设置缓存时间
	var expiration time.Duration
	switch expire {
	case 0:
		expiration = 2 * time.Minute
	case 1:
		expiration = 16 * time.Hour
	default:
		expiration = 6 * time.Hour
	}
	if err := m.CachedConn.SetCacheWithExpireCtx(ctx, key, string(data), expiration); err != nil {
		return err
	}
	return nil
}

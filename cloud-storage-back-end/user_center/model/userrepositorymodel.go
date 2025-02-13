package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Auroraol/cloud-storage/upload_service/model"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserRepositoryModel = (*customUserRepositoryModel)(nil)

type (
	// UserRepositoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserRepositoryModel.
	UserRepositoryModel interface {
		userRepositoryModel
		InsertWithId(ctx context.Context, data *UserRepository) (sql.Result, error)
		FindByRepositoryId(ctx context.Context, repositoryId int64) (*UserRepository, error)
		FindAllInPage(ctx context.Context, parentId int64, userId int64, startIndex int64, pageSize int64) ([]*UserRepository, error)
		FindAllFileByParentId(ctx context.Context, parentId int64, userId int64) ([]*UserRepository, error)
		FindAllFolderByParentId(ctx context.Context, parentId int64, userId int64) ([]*UserRepository, error)
		FindAllFolderAndByParentId(ctx context.Context, parentId int64, userId int64) ([]*UserRepository, error)
		FindAllFolderById(ctx context.Context, id int64, userId int64) ([]*UserRepository, error)
		CountByParentIdAndName(ctx context.Context, parentId int64, userId int64, Name string) (int64, error)
		CountByIdAndParentId(ctx context.Context, id int64, userId int64, parentId int64) (int64, error)
		CountTotalByIdAndParentId(ctx context.Context, userId int64, parentId int64) (int64, error)
		FindAllNormalInPage(ctx context.Context, parentId int64, userId int64, startIndex int64, pageSize int64) ([]*UserRepository, error)
		FindAllNormalFolderByParentId(ctx context.Context, parentId int64, userId int64) ([]*UserRepository, error)
		CountTotalNormalByIdAndParentId(ctx context.Context, userId int64, parentId int64) (int64, error)
		FindAllDeletedInPage(ctx context.Context, parentId int64, userId int64, startIndex int64, pageSize int64) ([]*UserRepository, error)
		FindAllDeletedByParentId(ctx context.Context, parentId int64, userId int64) ([]*UserRepository, error)
		CountTotalDeletedByUserId(ctx context.Context, userId int64) (int64, error)
	}

	customUserRepositoryModel struct {
		*defaultUserRepositoryModel
	}
)

// NewUserRepositoryModel returns a model for the database table.
func NewUserRepositoryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserRepositoryModel {
	return &customUserRepositoryModel{
		defaultUserRepositoryModel: newUserRepositoryModel(conn, c, opts...),
	}
}

func (m *defaultUserRepositoryModel) InsertWithId(ctx context.Context, data *UserRepository) (sql.Result, error) {
	// 构建缓存键
	userRepositoryIdKey := fmt.Sprintf("%s%v", cacheUserRepositoryIdPrefix, data.Id)

	// 执行插入操作
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("INSERT INTO %s (id, user_id, parent_id, repository_id, name) VALUES (?, ?, ?, ?, ?)", m.table)
		result, err = conn.ExecCtx(ctx, query, data.Id, data.UserId, data.ParentId, data.RepositoryId, data.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to insert user repository: %w", err)
		}
		return result, nil
	}, userRepositoryIdKey)

	if err != nil {
		return nil, fmt.Errorf("failed to insert user repository: %w", err)
	}

	return ret, nil
}

func (m *defaultUserRepositoryModel) FindByRepositoryId(ctx context.Context, repositoryId int64) (*UserRepository, error) {
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.Where("repository_id = ?", repositoryId).ToSql()
	var resp UserRepository
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, model.ErrNotFound
	default:
		return nil, err
	}
}

// 查询定文件夹下的所有子文件和子文件夹(分页)
func (m *defaultUserRepositoryModel) FindAllInPage(ctx context.Context, parentId int64, userId int64, startIndex int64, pageSize int64) ([]*UserRepository, error) {
	var resp []*UserRepository
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.Where("parent_id = ?", parentId).Where("user_id = ?", userId).OrderBy("repository_id").Offset(uint64(startIndex)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, model.ErrNotFound
	default:
		return nil, err
	}
}

// 查询定文件夹下的所有子文件和子文件夹
func (m *defaultUserRepositoryModel) FindAllFileByParentId(ctx context.Context, parentId int64, userId int64) ([]*UserRepository, error) {
	var resp []*UserRepository
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.Where("parent_id = ?", parentId).Where("user_id = ?", userId).Where("repository_id != ?", 0).ToSql()
	if err != nil {
		return nil, err
	}
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, model.ErrNotFound
	default:
		return nil, err
	}
}

// 查询定文件夹下的所有子文件和子文件夹
func (m *defaultUserRepositoryModel) FindAllFolderByParentId(ctx context.Context, parentId int64, userId int64) ([]*UserRepository, error) {
	var resp []*UserRepository
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.Where("parent_id = ?", parentId).Where("user_id = ?", userId).Where("repository_id = ?", 0).ToSql()
	if err != nil {
		return nil, err
	}
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, model.ErrNotFound
	default:
		return nil, err
	}
}

//

func (m *defaultUserRepositoryModel) FindAllFolderAndByParentId(ctx context.Context, parentId int64, userId int64) ([]*UserRepository, error) {
	var resp []*UserRepository
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.Where("parent_id = ?", parentId).Where("user_id = ?", userId).ToSql()
	if err != nil {
		return nil, err
	}
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, model.ErrNotFound
	default:
		return nil, err
	}
}

// 临时
func (m *defaultUserRepositoryModel) FindAllFolderById(ctx context.Context, id int64, userId int64) ([]*UserRepository, error) {
	var resp []*UserRepository
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.Where("id = ?", id).Where("user_id = ?", userId).ToSql()
	if err != nil {
		return nil, err
	}
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, model.ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserRepositoryModel) CountByParentIdAndName(ctx context.Context, parentId int64, userId int64, Name string) (int64, error) {
	countBuilder := m.CountBuilder("id")
	query, values, err := countBuilder.Where("parent_id = ?", parentId).Where("user_id = ?", userId).Where("name = ?", Name).ToSql()
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, model.ErrNotFound
	default:
		return 0, err
	}
}

func (m *defaultUserRepositoryModel) CountByIdAndParentId(ctx context.Context, id int64, userId int64, parentId int64) (int64, error) {
	countBuilder := m.CountBuilder("id")
	query, values, err := countBuilder.Where("parent_id = ?", parentId).Where("user_id = ?", userId).Where("id = ?", id).ToSql()
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return 0, model.ErrNotFound
	default:
		return 0, err
	}
}

// 计算目录下文件数量
func (m *defaultUserRepositoryModel) CountTotalByIdAndParentId(ctx context.Context, userId int64, parentId int64) (int64, error) {
	countBuilder := m.CountBuilder("id")
	query, values, err := countBuilder.Where("parent_id = ?", parentId).Where("user_id = ?", userId).ToSql()
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return 0, model.ErrNotFound
	default:
		return 0, err
	}
}

func (m *defaultUserRepositoryModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(userRepositoryRows).From(m.table)
}

func (m *defaultUserRepositoryModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("count(" + field + ")").From(m.table)
}

// 分页查询正常状态的文件和文件夹
func (m *defaultUserRepositoryModel) FindAllNormalInPage(ctx context.Context, parentId int64, userId int64, startIndex int64, pageSize int64) ([]*UserRepository, error) {
	var resp []*UserRepository
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.
		Where("parent_id = ?", parentId).
		Where("user_id = ?", userId).
		Where("status = ?", 0). // 只查询正常状态
		OrderBy("repository_id").
		Offset(uint64(startIndex)).
		Limit(uint64(pageSize)).
		ToSql()
	if err != nil {
		return nil, err
	}

	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, model.ErrNotFound
	default:
		return nil, err
	}
}

// 查询正常状态的文件夹
func (m *defaultUserRepositoryModel) FindAllNormalFolderByParentId(ctx context.Context, parentId int64, userId int64) ([]*UserRepository, error) {
	var resp []*UserRepository
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.
		Where("parent_id = ?", parentId).
		Where("user_id = ?", userId).
		Where("repository_id = ?", 0).
		Where("status = ?", 0). // 只查询正常状态
		ToSql()
	if err != nil {
		return nil, err
	}

	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, model.ErrNotFound
	default:
		return nil, err
	}
}

// 计算正常状态的文件和文件夹数量
func (m *defaultUserRepositoryModel) CountTotalNormalByIdAndParentId(ctx context.Context, userId int64, parentId int64) (int64, error) {
	countBuilder := m.CountBuilder("id")
	query, values, err := countBuilder.
		Where("parent_id = ?", parentId).
		Where("user_id = ?", userId).
		Where("status = ?", 0). // 只统计正常状态
		ToSql()

	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return 0, model.ErrNotFound
	default:
		return 0, err
	}
}

// 分页查询已删除状态的文件和文件夹
func (m *defaultUserRepositoryModel) FindAllDeletedInPage(ctx context.Context, parentId int64, userId int64, startIndex int64, pageSize int64) ([]*UserRepository, error) {
	var resp []*UserRepository
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.
		Where("parent_id = ?", parentId).
		Where("user_id = ?", userId).
		Where("status = ?", 1). // 查询已删除状态
		OrderBy("update_time DESC").
		Offset(uint64(startIndex)).
		Limit(uint64(pageSize)).
		ToSql()
	if err != nil {
		return nil, err
	}

	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, model.ErrNotFound
	default:
		return nil, err
	}
}

// 查询已删除状态的子文件和文件夹
func (m *defaultUserRepositoryModel) FindAllDeletedByParentId(ctx context.Context, parentId int64, userId int64) ([]*UserRepository, error) {
	var resp []*UserRepository
	rowBuilder := m.RowBuilder()
	query, values, err := rowBuilder.
		Where("parent_id = ?", parentId).
		Where("user_id = ?", userId).
		Where("status = ?", 1). // 查询已删除状态
		ToSql()
	if err != nil {
		return nil, err
	}

	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, model.ErrNotFound
	default:
		return nil, err
	}
}

// 统计已删除文件总数
func (m *defaultUserRepositoryModel) CountTotalDeletedByUserId(ctx context.Context, userId int64) (int64, error) {
	countBuilder := m.CountBuilder("id")
	query, values, err := countBuilder.
		Where("user_id = ?", userId).
		Where("status = ?", 1). // 统计已删除状态
		ToSql()

	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return 0, model.ErrNotFound
	default:
		return 0, err
	}
}

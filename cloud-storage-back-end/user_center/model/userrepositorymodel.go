package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/model"
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
		SearchFilesByKeywordInPage(ctx context.Context, parentId int64, userId int64, keyword string, startIndex int64, pageSize int64) ([]*UserRepository, error)
		CountSearchResultsByKeyword(ctx context.Context, parentId int64, userId int64, keyword string) (int64, error)
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
	userRepositoryIdKey := fmt.Sprintf("%s%s", cacheUserRepositoryIdPrefix, data.Id)

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

// 根据关键字搜索文件和文件夹(分页)
func (m *defaultUserRepositoryModel) SearchFilesByKeywordInPage(ctx context.Context, parentId int64, userId int64, keyword string, startIndex int64, pageSize int64) ([]*UserRepository, error) {
	var resp []*UserRepository
	rowBuilder := m.RowBuilder()

	var query string
	var values []interface{}
	var err error

	if parentId == 0 {
		// 全局搜索：在所有目录下搜索
		query, values, err = rowBuilder.
			Where("user_id = ?", userId).
			Where("status = ?", 0).                // 只搜索正常状态的文件和文件夹
			Where("name LIKE ?", "%"+keyword+"%"). // 使用LIKE进行模糊匹配
			OrderBy("repository_id").              // 按repository_id排序，文件夹在前，文件在后
			Offset(uint64(startIndex)).
			Limit(uint64(pageSize)).
			ToSql()
	} else {
		// 在指定目录及其所有子目录下搜索
		// 首先获取指定目录的所有子目录ID
		var folderIds []int64
		folderIds = append(folderIds, parentId)

		// 递归查找所有子文件夹
		subFolders, err := m.findAllSubFolders(ctx, parentId, userId)
		if err != nil && !errors.Is(err, model.ErrNotFound) {
			return nil, err
		}

		for _, folder := range subFolders {
			folderIds = append(folderIds, int64(folder.Id))
		}

		// 在这些文件夹中搜索
		query, values, err = rowBuilder.
			Where("user_id = ?", userId).
			Where("status = ?", 0).                     // 只搜索正常状态的文件和文件夹
			Where("name LIKE ?", "%"+keyword+"%").      // 使用LIKE进行模糊匹配
			Where(squirrel.Eq{"parent_id": folderIds}). // 在指定的文件夹列表中搜索
			OrderBy("repository_id").                   // 按repository_id排序，文件夹在前，文件在后
			Offset(uint64(startIndex)).
			Limit(uint64(pageSize)).
			ToSql()
	}

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

// 递归查找所有子文件夹
func (m *defaultUserRepositoryModel) findAllSubFolders(ctx context.Context, parentId int64, userId int64) ([]*UserRepository, error) {
	var allFolders []*UserRepository

	// 查找直接子文件夹
	folders, err := m.FindAllFolderByParentId(ctx, parentId, userId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return allFolders, nil // 没有子文件夹，返回空列表
		}
		return nil, err
	}

	allFolders = append(allFolders, folders...)

	// 递归查找每个子文件夹的子文件夹
	for _, folder := range folders {
		subFolders, err := m.findAllSubFolders(ctx, int64(folder.Id), userId)
		if err != nil && !errors.Is(err, model.ErrNotFound) {
			return nil, err
		}
		allFolders = append(allFolders, subFolders...)
	}

	return allFolders, nil
}

// 统计搜索结果总数
func (m *defaultUserRepositoryModel) CountSearchResultsByKeyword(ctx context.Context, parentId int64, userId int64, keyword string) (int64, error) {
	countBuilder := m.CountBuilder("id")

	var query string
	var values []interface{}
	var err error

	if parentId == 0 {
		// 全局搜索：统计所有目录下的结果
		query, values, err = countBuilder.
			Where("user_id = ?", userId).
			Where("status = ?", 0).                // 只统计正常状态的文件和文件夹
			Where("name LIKE ?", "%"+keyword+"%"). // 使用LIKE进行模糊匹配
			ToSql()
	} else {
		// 在指定目录及其所有子目录下统计
		// 首先获取指定目录的所有子目录ID
		var folderIds []int64
		folderIds = append(folderIds, parentId)

		// 递归查找所有子文件夹
		subFolders, err := m.findAllSubFolders(ctx, parentId, userId)
		if err != nil && !errors.Is(err, model.ErrNotFound) {
			return 0, err
		}

		for _, folder := range subFolders {
			folderIds = append(folderIds, int64(folder.Id))
		}

		// 在这些文件夹中统计
		query, values, err = countBuilder.
			Where("user_id = ?", userId).
			Where("status = ?", 0).                     // 只统计正常状态的文件和文件夹
			Where("name LIKE ?", "%"+keyword+"%").      // 使用LIKE进行模糊匹配
			Where(squirrel.Eq{"parent_id": folderIds}). // 在指定的文件夹列表中统计
			ToSql()
	}

	if err != nil {
		return 0, err
	}

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

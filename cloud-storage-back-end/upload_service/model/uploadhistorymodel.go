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
	"strings"
)

var _ UploadHistoryModel = (*customUploadHistoryModel)(nil)

type (
	// UploadHistoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUploadHistoryModel.
	UploadHistoryModel interface {
		uploadHistoryModel
		FindAllInPage(ctx context.Context, userId int64, startIndex int64, pageSize int64) ([]*UploadHistory, error)
		//UpdateHistory(ctx context.Context, id int64, data *UploadHistory) error
		UpdateHistory(ctx context.Context, data *UploadHistory) (sql.Result, error)
		CountTotalByUserIdId(ctx context.Context, userId int64) (int64, error)
		DeleteAllByIdList(ctx context.Context, ids []string) error
	}

	customUploadHistoryModel struct {
		*defaultUploadHistoryModel
	}
)

// NewUploadHistoryModel returns a model for the database table.
func NewUploadHistoryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UploadHistoryModel {
	return &customUploadHistoryModel{
		defaultUploadHistoryModel: newUploadHistoryModel(conn, c, opts...),
	}
}

func (m *defaultUploadHistoryModel) FindAllInPage(ctx context.Context, userId int64, startIndex int64, pageSize int64) ([]*UploadHistory, error) {
	query := "select id, repository_id, user_id, file_name, size, status, create_time, update_time from upload_history where user_id = ? order by update_time desc limit ?, ?"
	var resp []*UploadHistory
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId, startIndex, pageSize)
	return resp, err
}

//func (m *defaultUploadHistoryModel) UpdateHistory(ctx context.Context, id int64, data *UploadHistory) error {
//	// 检查记录是否存在
//	var exists bool
//	err := m.QueryRowNoCacheCtx(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM upload_history WHERE id = ?)", id)
//	if err != nil {
//		return err
//	}
//
//	if !exists {
//		// 如果记录不存在，则插入数据
//		_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
//			query := "INSERT INTO upload_history (user_id, file_name, size, status, repository_id) VALUES (?, ?, ?, ?, ?)"
//			result, err = conn.ExecCtx(ctx, query, data.UserId, data.FileName, data.Size, data.Status, data.RepositoryId)
//			return
//		})
//		return err
//	} else {
//		// 如果记录存在，则更新数据
//		_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
//			query := "UPDATE upload_history SET file_name = ?, size = ?, status = ?"
//			args := []interface{}{data.FileName, data.Size, data.Status}
//
//			if data.RepositoryId != 0 {
//				query += ", repository_id = ?"
//				args = append(args, data.RepositoryId)
//			}
//
//			query += " WHERE id = ?"
//			args = append(args, id)
//
//			result, err = conn.ExecCtx(ctx, query, args...)
//			return
//		})
//		return err
//	}
//}

func (m *defaultUploadHistoryModel) UpdateHistory(ctx context.Context, data *UploadHistory) (sql.Result, error) {
	uploadHistoryIdKey := fmt.Sprintf("%s%v", cacheUploadHistoryIdPrefix, data.Id)
	uploadHistoryRepositoryIdKey := fmt.Sprintf("%s%v", cacheUploadHistoryRepositoryIdPrefix, data.RepositoryId)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, "id, user_id, file_name, size, repository_id, status")
		return conn.ExecCtx(ctx, query, data.Id, data.UserId, data.FileName, data.Size, data.RepositoryId, data.Status)
	}, uploadHistoryIdKey, uploadHistoryRepositoryIdKey)
	return ret, err
}

// 计算目录下文件数量
func (m *defaultUploadHistoryModel) CountTotalByUserIdId(ctx context.Context, userId int64) (int64, error) {
	countBuilder := m.CountBuilder("id")
	query, values, err := countBuilder.Where("user_id = ?", userId).ToSql()
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, values...)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return 0, err
	default:
		return 0, err
	}
}
func (m *defaultUploadHistoryModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(repositoryPoolRows).From(m.table)
}

func (m *defaultUploadHistoryModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("count(" + field + ")").From(m.table)
}

func (m *defaultUploadHistoryModel) DeleteAllByIdList(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil // 如果用户ID列表为空，直接返回成功
	}

	// 动态生成 IN 子句中的占位符
	placeholders := make([]string, len(ids))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE id IN (%s)", m.table, strings.Join(placeholders, ","))

	// 将 []int64 转换为 []interface{}
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// 执行删除操作
	result, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.Exec(query, args...)
	})
	if err != nil {
		return fmt.Errorf("删除失败: %v", err)
	}

	// 检查受影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取受影响行数失败: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("没有找到匹配的记录")
	}

	return nil
}

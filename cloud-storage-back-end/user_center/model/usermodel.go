package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil) //是一个类型断言的技巧，用于确保 customUserModel 结构体实现了 UserModel

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		withSession(session sqlx.Session) UserModel
		FindUserIdByUsernameAndPassword(ctx context.Context, username string, password string) (int64, error)
		UpdateVolume(ctx context.Context, id int64, size int64) (result sql.Result, err error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	// 创建的默认用户模型实例
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

func (m *customUserModel) withSession(session sqlx.Session) UserModel {
	return NewUserModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultUserModel) FindUserIdByUsernameAndPassword(ctx context.Context, username string, password string) (int64, error) {
	var resp User
	query := fmt.Sprintf("select %s from %s where `username` = ? and password = ? limit 1", userRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, username, password)
	switch {
	case err == nil:
		return resp.Id, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return -1, ErrNotFound
	default:
		return -1, err
	}
}

func (m *defaultUserModel) UpdateVolume(ctx context.Context, id int64, size int64) (result sql.Result, err error) {
	query := fmt.Sprintf("update %s set now_volume = now_volume + ? where `id` = ? and `now_volume` + ? <= `total_volume`", m.table)
	res, err := m.conn.ExecCtx(ctx, query, size, id, size)
	if err != nil {
		return nil, err
	}
	return res, nil
}

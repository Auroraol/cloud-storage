package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Auroraol/cloud-storage/common/time"

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
		UpdateAvatar(ctx context.Context, id int64, avatar string) (result sql.Result, err error)
		UpdatePassword(ctx context.Context, id int64, password string) (err error)
		UpdateInfo(ctx context.Context, id int64, nickname string, brief string, birthday string, gender int64, email string, mobile string) (err error)
		FindOneByPassword(ctx context.Context, id int64, password string) (*User, error)
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

func (m *defaultUserModel) UpdateAvatar(ctx context.Context, id int64, avatar string) (result sql.Result, err error) {
	query := fmt.Sprintf("update %s set `avatar` = ? where `id` = ?", m.table)
	res, err := m.conn.ExecCtx(ctx, query, avatar, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *defaultUserModel) FindOneByPassword(ctx context.Context, id int64, password string) (*User, error) {
	var resp User
	query := fmt.Sprintf("select %s from %s where `id` = ? and `password` = ? limit 1", userRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, id, password)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) UpdateInfo(ctx context.Context, id int64, nickname string, brief string, birthday string, gender int64, email string, mobile string) (err error) {
	var setStmt []string
	var args []interface{}

	// 动态构建更新字段
	if nickname != "" {
		setStmt = append(setStmt, "`nickname` = ?")
		args = append(args, nickname)
	}
	if brief != "" {
		setStmt = append(setStmt, "`brief` = ?")
		args = append(args, brief)
	}
	if birthday != "" {
		timestamp, err := strconv.ParseInt(birthday, 10, 64)
		if err != nil {
			return err
		}
		birthday = time.TimestampToStringTime(timestamp)
		setStmt = append(setStmt, "`birthday` = ?")
		args = append(args, birthday)
	}
	if gender != 0 {
		setStmt = append(setStmt, "`gender` = ?")
		args = append(args, gender)
	}
	if email != "" {
		setStmt = append(setStmt, "`email` = ?")
		args = append(args, email)
	}
	if mobile != "" {
		setStmt = append(setStmt, "`mobile` = ?")
		args = append(args, mobile)
	}

	// 如果没有需要更新的字段，直接返回
	if len(setStmt) == 0 {
		return errors.New("没有需要更新的信息")
	}

	// 构建完整的 SQL 语句
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setStmt, ", "))
	args = append(args, id)

	_, err = m.conn.ExecCtx(ctx, query, args...)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, sqlx.ErrNotFound):
		return ErrNotFound
	default:
		return err
	}
}

func (m *defaultUserModel) UpdatePassword(ctx context.Context, id int64, password string) error {
	query := fmt.Sprintf("UPDATE %s SET `password` = ? WHERE `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, password, id)
	if err != nil {
		return err
	}
	return nil
}

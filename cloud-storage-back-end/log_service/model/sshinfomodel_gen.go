// Code generated by goctl. DO NOT EDIT.
// versions:
//  goctl version: 1.7.3

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	sshInfoFieldNames          = builder.RawFieldNames(&SshInfo{})
	sshInfoRows                = strings.Join(sshInfoFieldNames, ",")
	sshInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(sshInfoFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	sshInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(sshInfoFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheSshInfoIdPrefix         = "cache:sshInfo:id:"
	cacheSshInfoUserIdHostPrefix = "cache:sshInfo:userId:host:"
)

type (
	sshInfoModel interface {
		Insert(ctx context.Context, data *SshInfo) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*SshInfo, error)
		FindOneByUserIdHost(ctx context.Context, userId int64, host string) (*SshInfo, error)
		Update(ctx context.Context, data *SshInfo) error
		Delete(ctx context.Context, id int64) error
	}

	defaultSshInfoModel struct {
		sqlc.CachedConn
		table string
	}

	SshInfo struct {
		Id        int64     `db:"id"`         // 主键ID
		UserId    int64     `db:"user_id"`    // 用户ID
		Host      string    `db:"host"`       // 主机地址
		Port      int64     `db:"port"`       // 端口号
		Username  string    `db:"username"`   // 用户名
		Password  string    `db:"password"`   // 密码
		CreatedAt time.Time `db:"created_at"` // 创建时间
		UpdatedAt time.Time `db:"updated_at"` // 更新时间
	}
)

func newSshInfoModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultSshInfoModel {
	return &defaultSshInfoModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`ssh_info`",
	}
}

func (m *defaultSshInfoModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	sshInfoIdKey := fmt.Sprintf("%s%v", cacheSshInfoIdPrefix, id)
	sshInfoUserIdHostKey := fmt.Sprintf("%s%v:%v", cacheSshInfoUserIdHostPrefix, data.UserId, data.Host)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, sshInfoIdKey, sshInfoUserIdHostKey)
	return err
}

func (m *defaultSshInfoModel) FindOne(ctx context.Context, id int64) (*SshInfo, error) {
	sshInfoIdKey := fmt.Sprintf("%s%v", cacheSshInfoIdPrefix, id)
	var resp SshInfo
	err := m.QueryRowCtx(ctx, &resp, sshInfoIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", sshInfoRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultSshInfoModel) FindOneByUserIdHost(ctx context.Context, userId int64, host string) (*SshInfo, error) {
	sshInfoUserIdHostKey := fmt.Sprintf("%s%v:%v", cacheSshInfoUserIdHostPrefix, userId, host)
	var resp SshInfo
	err := m.QueryRowIndexCtx(ctx, &resp, sshInfoUserIdHostKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? and `host` = ? limit 1", sshInfoRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, userId, host); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultSshInfoModel) Insert(ctx context.Context, data *SshInfo) (sql.Result, error) {
	sshInfoIdKey := fmt.Sprintf("%s%v", cacheSshInfoIdPrefix, data.Id)
	sshInfoUserIdHostKey := fmt.Sprintf("%s%v:%v", cacheSshInfoUserIdHostPrefix, data.UserId, data.Host)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, sshInfoRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.UserId, data.Host, data.Port, data.Username, data.Password)
	}, sshInfoIdKey, sshInfoUserIdHostKey)
	return ret, err
}

func (m *defaultSshInfoModel) Update(ctx context.Context, newData *SshInfo) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	sshInfoIdKey := fmt.Sprintf("%s%v", cacheSshInfoIdPrefix, data.Id)
	sshInfoUserIdHostKey := fmt.Sprintf("%s%v:%v", cacheSshInfoUserIdHostPrefix, data.UserId, data.Host)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, sshInfoRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.UserId, newData.Host, newData.Port, newData.Username, newData.Password, newData.Id)
	}, sshInfoIdKey, sshInfoUserIdHostKey)
	return err
}

func (m *defaultSshInfoModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheSshInfoIdPrefix, primary)
}

func (m *defaultSshInfoModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", sshInfoRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultSshInfoModel) tableName() string {
	return m.table
}

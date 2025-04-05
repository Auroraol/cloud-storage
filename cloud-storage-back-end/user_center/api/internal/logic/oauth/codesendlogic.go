package oauth

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"strconv"
	"time"

	"github.com/Auroraol/cloud-storage/common/tool"
	"github.com/Auroraol/cloud-storage/user_center/model"
	"github.com/pkg/errors"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type CodeSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 验证码发送
func NewCodeSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CodeSendLogic {
	return &CodeSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CodeSend 发送验证码
func (l *CodeSendLogic) CodeSend(req *types.CodeSendRequest) (resp *types.CodeSendResponse, err error) {
	_, err = l.svcCtx.Sms.SendVerificationCode(req.Mobile)
	if err != nil {
		zap.S().Errorf("验证码发送失败: %s", err.Error())
		return nil, errors.New("验证码发送失败")
	}
	return &types.CodeSendResponse{}, nil
}

// 注册新用户
func (l *LoginByMobileCodeLogic) register(mobile int64) (int64, error) {
	// 1. 创建用户数据
	now := time.Now()
	mobileStr := strconv.FormatInt(mobile, 10)

	user := &model.User{
		Version: 1,
		Username: sql.NullString{
			String: mobileStr,
			Valid:  true,
		},
		Password: sql.NullString{
			String: "", // 不设置密码
			Valid:  false,
		},
		Mobile: sql.NullInt64{
			Int64: mobile,
			Valid: true,
		},
		Nickname:    tool.Krand(8, tool.KC_RAND_KIND_ALL),
		Gender:      1, // 默认性别男
		Avatar:      sql.NullString{String: "", Valid: false},
		Birthday:    sql.NullTime{Time: time.Time{}, Valid: false},
		Email:       sql.NullString{String: "", Valid: false},
		Brief:       sql.NullString{String: "", Valid: false},
		Info:        sql.NullString{String: "", Valid: false},
		DelState:    sql.NullInt64{Int64: 0, Valid: true},
		DeleteTime:  sql.NullTime{Time: time.Time{}, Valid: false},
		Status:      0, // 正常状态
		Admin:       0, // 非管理员
		NowVolume:   0,
		TotalVolume: 1073741824, // 默认10GB空间
		CreateTime:  now,
		UpdateTime:  now,
	}

	// 2. 插入用户数据
	result, err := l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		zap.S().Errorf("新用户注册失败: %s", err.Error())
		return 0, errors.New("注册失败")
	}

	// 3. 获取新用户ID
	userId, err := result.LastInsertId()
	if err != nil {
		zap.S().Errorf("获取新用户ID失败: %s", err.Error())
		return 0, errors.New("注册失败")
	}

	zap.S().Infof("用户[%d]注册成功，ID: %d", mobile, userId)
	return userId, nil
}

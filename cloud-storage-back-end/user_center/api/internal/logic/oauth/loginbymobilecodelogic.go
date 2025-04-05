package oauth

import (
	"context"
	"database/sql"

	"strconv"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/common/token"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginByMobileCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 手机号登陆/注册
func NewLoginByMobileCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginByMobileCodeLogic {
	return &LoginByMobileCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginByMobileCodeLogic) LoginByMobileCode(req *types.MobileLoginReq) (resp *types.MobileLoginResp, err error) {
	// 1. 参数校验
	if req.Mobile == "" {
		return nil, response.NewErrCodeMsg(response.DATA_PARAM_ERROR, "手机号不能为空")
	}
	if req.Code == "" {
		return nil, response.NewErrCodeMsg(response.DATA_PARAM_ERROR, "验证码不能为空")
	}

	// 2. 获取短信客户端
	smsClient := l.svcCtx.Sms

	// 3. 验证短信验证码
	if !smsClient.VerifyCode(req.Mobile, req.Code) {
		return nil, response.NewErrMsg("验证码错误或已过期")
	}

	// 验证成功后移除验证码
	smsClient.RemoveVerificationCode(req.Mobile)

	// 4. 查询用户是否存在
	mobile, _ := strconv.ParseInt(req.Mobile, 10, 64)
	mobileSQL := sql.NullInt64{
		Int64: mobile,
		Valid: true,
	}
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, mobileSQL)
	if err != nil && !errors.Is(err, sqlc.ErrNotFound) {
		zap.S().Error("查询用户信息失败: %s", err.Error())
		return nil, response.NewErrCode(response.SYSTEM_ERROR)
	}

	// 5. 用户不存在，进行注册
	if errors.Is(err, sqlc.ErrNotFound) {
		zap.S().Errorf("用户[%s]不存在，进行注册", req.Mobile)

		// 创建新用户
		userId, err := l.register(mobile)
		if err != nil {
			zap.S().Error("注册用户失败: %s", err)
			return nil, response.NewErrCode(response.SYSTEM_ERROR)
		}

		// 重新查询用户信息
		user, err = l.svcCtx.UserModel.FindOne(l.ctx, userId)
		if err != nil {
			zap.S().Error("查询新注册用户信息失败: %s", err)
			return nil, response.NewErrCode(response.SYSTEM_ERROR)
		}
	}

	// 6. 生成token
	tokenResp, err := token.GenerateToken(&token.GenerateTokenReq{
		AccessExpire: l.svcCtx.Config.JwtAuth.AccessExpire,
		AccessSecret: l.svcCtx.Config.JwtAuth.AccessSecret,
		UserId:       user.Id,
	})
	if err != nil {
		zap.S().Error("生成token失败: %s", err.Error())
		return nil, response.NewErrCode(response.SYSTEM_ERROR)
	}

	// 7. 返回登录结果
	return &types.MobileLoginResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}

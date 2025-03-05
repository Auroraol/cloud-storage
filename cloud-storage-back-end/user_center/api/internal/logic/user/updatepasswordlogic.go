package user

import (
	"context"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	"github.com/Auroraol/cloud-storage/common/utils"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改密码
func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePasswordLogic) UpdatePassword(req *types.UserPasswordReq) (resp *types.UserPasswordResp, err error) {
	// 获取用户ID
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 获取用户信息
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		zap.S().Error("获取用户信息失败: %v", err)
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "获取用户信息失败")
	}

	// 验证原密码
	if !user.Password.Valid || !utils.ComparePassword(user.Password.String, req.OldPassword) {
		zap.S().Error("原密码不正确")
		return nil, response.NewErrCodeMsg(response.INVALID_REQUEST, "原密码不正确")
	}

	// 加密新密码
	hashedPassword := utils.Md5ByString(req.NewPassword)
	// 更新密码
	err = l.svcCtx.UserModel.UpdatePassword(l.ctx, userId, hashedPassword)
	if err != nil {
		zap.S().Error("更新密码失败: %v", err)
		return nil, response.NewErrCode(response.SYSTEM_ERROR)
	}

	//// 清除用户缓存
	//if err := l.svcCtx.Cache.Del(fmt.Sprintf("user:%d", userId)); err != nil {
	//	l.Logger.Errorf("清除用户缓存失败: %v", err)
	//}

	return &types.UserPasswordResp{}, nil
}

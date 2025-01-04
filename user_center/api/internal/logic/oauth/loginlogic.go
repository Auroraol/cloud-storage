package oauth

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/token"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 账号密码登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.AccountLoginReq) (resp *types.AccountLoginResp, err error) {
	var userId int64
	// 1、从数据库中查询当前用户

	// 2、生成token
	tokenResp, _ := token.GenerateToken(&token.GenerateTokenReq{
		AccessExpire: l.svcCtx.Config.JwtAuth.AccessExpire,
		AccessSecret: l.svcCtx.Config.JwtAuth.AccessSecret,
		UserId:       userId,
	})

	return &types.AccountLoginResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}

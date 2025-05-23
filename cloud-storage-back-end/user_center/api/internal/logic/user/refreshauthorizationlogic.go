package user

import (
	"context"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/token"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/types"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshAuthorizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 刷新Authorization
func NewRefreshAuthorizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshAuthorizationLogic {
	return &RefreshAuthorizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshAuthorizationLogic) RefreshAuthorization(req *types.RefreshAuthRequest, Authorization string) (resp *types.RefreshAuthResponse, err error) {
	//获得原token的剩余信息
	restClaims := make(jwt.MapClaims)
	judgeValid, err := jwt.ParseWithClaims(Authorization, restClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(l.svcCtx.Config.JwtAuth.AccessSecret), nil
	})
	if err != nil {
		zap.S().Error("jwt.ParseWithClaims err:%s", err)
		return nil, err
	}
	//判断是否token有效
	if !judgeValid.Valid {
		zap.S().Error("token无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 利用过期token的其他值，生成新token等信息
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.REFRESH_CREDENTIALS_INVALID)
	}

	tokenResp, _ := token.GenerateToken(&token.GenerateTokenReq{
		AccessExpire: l.svcCtx.Config.JwtAuth.AccessExpire,
		AccessSecret: l.svcCtx.Config.JwtAuth.AccessSecret,
		UserId:       userId,
	})

	return &types.RefreshAuthResponse{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}

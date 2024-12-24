package token

import (
	"context"

	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

/**
使用: GenerateJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, userInfo.Id)
获得自定义的负载数据: ctx.Value(CtxKeyJwtUserId).(json.Number)
*/

// 从任务上下文获取JWT(token)中的uid, 用户信息优先使用上下文获取
var CtxKeyJwtUserId = "jwtUserId"

// 生成JWTToken
func GenerateJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[CtxKeyJwtUserId] = userId // 自定义的负载（payload），可以设置为任何信息，例如用户名、用户ID等
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

// 从任务上下文获取JWT(token)中的uid
func GetUidFromCtx(ctx context.Context) int64 {
	var uid int64
	if jsonUid, ok := ctx.Value(CtxKeyJwtUserId).(json.Number); ok {
		if int64Uid, err := jsonUid.Int64(); err == nil {
			uid = int64Uid
		} else {
			logx.WithContext(ctx).Errorf("GetUidFromCtx err : %+v", err)
		}
	}
	return uid
}

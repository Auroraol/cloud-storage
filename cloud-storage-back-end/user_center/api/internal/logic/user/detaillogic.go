package user

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户详情
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		zap.S().Error("UserModel.FindOne err:%s", err)
		return nil, err
	}
	if user == nil {
		zap.S().Error("UserModel.FindOne err:%s", err)
		return nil, errors.Wrapf(response.NewErrCode(response.ACCOUNT_NOT_FOUND), "id:%d", userId)
	}
	resp = &types.UserInfoResp{}
	_ = copier.Copy(&resp, user)
	// 返回字段roles数组
	// resp.admin 存在 则加入roles数组
	if user.Admin == 1 {
		resp.Roles = append(resp.Roles, "admin")
	}
	resp.Roles = []string{"admin"}
	return resp, nil
}

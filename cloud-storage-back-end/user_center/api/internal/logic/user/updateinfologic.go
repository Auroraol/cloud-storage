package user

import (
	"context"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/token"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/types"

	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改用户信息
func NewUpdateInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInfoLogic {
	return &UpdateInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateInfoLogic) UpdateInfo(req *types.UpdateInfoReq) (resp *types.UpdateInfoResp, err error) {
	// 从context中获取用户ID
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 调用 UserModel 的 UpdateInfo 方法
	err = l.svcCtx.UserModel.UpdateInfo(l.ctx, userId, req.Nickname, req.Brief, req.Birthday, req.Gender, req.Email, req.Mobile)
	if err != nil {
		zap.S().Error("更新用户信息失败, err: %s", err)
		return nil, errors.New("更新用户信息失败")
	}

	return &types.UpdateInfoResp{}, nil
}

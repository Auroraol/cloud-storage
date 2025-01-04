package oauth

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/tool"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 账号密码注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.AccountRegisterReq) (resp *types.AccountRegisterResp, err error) {
	// 1、从数据库中检查是否已注册

	// 2、生成新用户，插入数据库中
	tool.Md5ByString("!111")
	return &types.AccountRegisterResp{}, nil
}

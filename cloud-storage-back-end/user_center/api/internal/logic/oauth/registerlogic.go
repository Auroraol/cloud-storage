package oauth

import (
	"context"
	"database/sql"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/tool"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/utils"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/types"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"

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
	userName := sql.NullString{
		String: req.Name,
		Valid:  true,
	}
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, userName)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		zap.S().Error("用户存在 err:%s", err)
		return nil, errors.Wrapf(response.NewErrCode(response.DB_ERROR), "Register user exists name:%s,err:%v", req.Name, err)
	}
	if user != nil {
		return nil, response.NewErrCode(response.ACCOUNT_REGISTERED_ERROR)
	}
	// 2、生成新用户，插入数据库中
	user = &model.User{}
	user.Username = sql.NullString{
		String: req.Name,
		Valid:  true,
	}
	if len(user.Nickname) == 0 {
		user.Nickname = tool.Krand(8, tool.KC_RAND_KIND_ALL)
	}
	if len(req.Password) > 0 {
		password := sql.NullString{
			String: utils.Md5ByString(req.Password),
			Valid:  true,
		}
		user.Password = password
	}
	user.TotalVolume = 1073741824 // 用户初始大小1G
	_, err = l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		zap.S().Error("用户注册失败 err:%s", err)
		return nil, errors.Wrapf(response.NewErrCode(response.DB_ERROR), "Register db user Insert err:%v,user:%+v", err, user)
	}
	return &types.AccountRegisterResp{}, nil
}

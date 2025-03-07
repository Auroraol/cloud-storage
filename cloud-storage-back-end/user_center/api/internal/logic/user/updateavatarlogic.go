package user

import (
	"context"
	"fmt"
	"github.com/Auroraol/cloud-storage/common/store/oss"
	"go.uber.org/zap"
	"mime/multipart"
	"path/filepath"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/time"
	"github.com/Auroraol/cloud-storage/common/token"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更换头像
func NewUpdateAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAvatarLogic {
	return &UpdateAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAvatarLogic) UpdateAvatar(req *types.UserAvatarReq, file multipart.File, fileHeader *multipart.FileHeader) (resp *types.UserAvatarResp, err error) {
	// 获取用户ID
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 获取文件信息
	filename := fileHeader.Filename
	contentType := fileHeader.Header.Get("Content-Type")
	extension := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%d%s", time.LocalTimeNow().UnixNano(), extension)

	// 上传文件到存储服务
	avatarUrl, err := oss.Upload(file, newFilename, contentType)
	if err != nil {
		zap.S().Error("上传文件失败: %s", err)
		return nil, response.NewErrCode(response.SYSTEM_ERROR)
	}

	// 删除旧头像文件
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		zap.S().Error("获取用户信息失败: %v", err)
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "获取用户信息失败")
	}
	if _, err := oss.Delete(user.Avatar.String); err != nil {
		zap.S().Error("删除旧头像失败: %v", err)
	}

	// 更新用户头像
	_, err = l.svcCtx.UserModel.UpdateAvatar(l.ctx, userId, avatarUrl)
	if err != nil {
		zap.S().Error("更新头像失败: %v", err)
		// 删除已上传的新文件
		if _, err := oss.Delete(avatarUrl); err != nil {
			zap.S().Error("删除新上传的头像失败: %v", err)
		}
		return nil, response.NewErrCodeMsg(response.SYSTEM_ERROR, "更新头像失败")
	}

	//// 清除用户缓存
	//if err := l.svcCtx.Cache.Del(fmt.Sprintf("user:%d", userId)); err != nil {
	//	l.Logger.Errorf("清除用户缓存失败: %v", err)
	//}

	return &types.UserAvatarResp{}, nil
}

package user

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	"mime/multipart"

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
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	resp = &types.UserAvatarResp{}

	// 获取文件信息
	//filename := fileHeader.Filename
	//contentType := fileHeader.Header.Get("Content-Type")
	//extension := filepath.Ext(filename)
	//newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
	//
	//// 上传文件到存储服务
	//avatarUrl, err := l.svcCtx.Storage.Upload(file, newFilename, contentType)
	//if err != nil {
	//	l.Logger.Errorf("上传文件失败: %v", err)
	//	return nil, response.NewErrCode(response.SYSTEM_ERROR)
	//}
	//
	//// 更新用户头像
	//_, err = l.svcCtx.UserModel.UpdateAvatar(l.ctx, userId, avatarUrl)
	//if err != nil {
	//	l.Logger.Errorf("更新头像失败: %v", err)
	//	// 删除已上传的新文件
	//	if err := l.svcCtx.Storage.Delete(avatarUrl); err != nil {
	//		l.Logger.Errorf("删除新上传的头像失败: %v", err)
	//		return nil, response.NewErrCode(response.SYSTEM_ERROR)
	//	}
	//}
	//
	//// 删除旧头像文件
	//if err := l.svcCtx.Storage.Delete(oldAvatar); err != nil {
	//	l.Logger.Errorf("删除旧头像失败: %v", err)
	//}
	//
	//// 清除用户缓存
	//if err := l.svcCtx.Cache.Del(fmt.Sprintf("user:%d", userId)); err != nil {
	//	l.Logger.Errorf("清除用户缓存失败: %v", err)
	//}

	return &types.UserAvatarResp{}, nil
}

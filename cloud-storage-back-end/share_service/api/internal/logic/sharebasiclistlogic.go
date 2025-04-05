package logic

import (
	"context"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/token"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/share_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/share_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareBasicListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户分享列表
func NewShareBasicListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicListLogic {
	return &ShareBasicListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicListLogic) ShareBasicList(req *types.ShareBasicListRequest) (resp *types.ShareBasicListResponse, err error) {
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 计算Offset和Limit
	page := req.Page
	pageSize := req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 5 // 默认每页5条记录
	}
	offset := (page - 1) * pageSize

	// 获取分享列表
	shareFile := make([]*types.ShareBasicDetailReply, 0)
	var total int64

	// 查询总记录数
	err = l.svcCtx.Engine.
		Table("share_basic").
		Where("share_basic.user_id = ? AND share_basic.deleted_at IS NULL", userId).
		Count(&total).Error
	if err != nil {
		zap.S().Error("查询分享列表总数失败")
		return nil, err
	}

	// 分页查询
	err = l.svcCtx.Engine.
		Table("share_basic").
		Select("share_basic.id, share_basic.repository_id, user_repository.name, repository_pool.ext, repository_pool.path, repository_pool.size, share_basic.click_num, user.username AS owner, user.avatar, share_basic.expired_time, share_basic.update_time, share_basic.code").
		Joins("LEFT JOIN repository_pool ON repository_pool.identity = share_basic.repository_id").
		Joins("LEFT JOIN user_repository ON user_repository.id = share_basic.user_repository_id").
		Joins("LEFT JOIN user ON share_basic.user_id = user.id").
		Where("share_basic.user_id = ? AND share_basic.deleted_at IS NULL", userId).
		Order("share_basic.update_time DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&shareFile).Error

	if err != nil {
		zap.S().Error("查询分享列表失败 err:%v", err)
		return nil, err
	}

	return &types.ShareBasicListResponse{
		List:  shareFile,
		Total: total,
	}, nil
}

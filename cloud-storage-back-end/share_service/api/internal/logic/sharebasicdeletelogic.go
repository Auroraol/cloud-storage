package logic

import (
	"context"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/token"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/share_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/share_service/api/internal/types"
	"go.uber.org/zap"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareBasicDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 资源删除
func NewShareBasicDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicDeleteLogic {
	return &ShareBasicDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicDeleteLogic) ShareBasicDelete(req *types.ShareBasicDeleteRequest) (resp *types.ShareBasicDeleteResponse, err error) {
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}
	id, err := strconv.ParseInt(req.Id, 10, 64)
	if err != nil {
		zap.S().Error("转换数字失败！")
		return nil, err
	}
	err = l.svcCtx.ShareBasicModel.Delete(l.ctx, uint64(id))
	if err != nil {
		zap.S().Error("删除失败！")
		return nil, response.NewErrCodeMsg(response.DB_UPDATE_AFFECTED_ZERO_ERROR, "删除失败！")
	}

	return &types.ShareBasicDeleteResponse{}, nil
}

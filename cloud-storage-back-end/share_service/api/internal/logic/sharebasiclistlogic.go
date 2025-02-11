package logic

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"

	"github.com/Auroraol/cloud-storage/share_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/share_service/api/internal/types"

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
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 获取分享列表

	return &types.ShareBasicListResponse{
		List: []*types.ShareBasicListResponseData{
			{
				Id:          1,
				UpdataTime:  1,
				ExpireTime:  1,
				Filename:    "test.txt",
				BrowseCount: 1,
			},
		},
	}, nil
}

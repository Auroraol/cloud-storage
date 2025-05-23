package logic

import (
	"context"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/token"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/share_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/share_service/api/internal/types"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/rpc/client/uploadservicerpc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/client/userrepositoryrpc"
	"go.uber.org/zap"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareBasicSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 资源保存
func NewShareBasicSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicSaveLogic {
	return &ShareBasicSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicSaveLogic) ShareBasicSave(req *types.ShareBasicSaveRequest) (resp *types.ShareBasicSaveResponse, err error) {
	nameInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadservicerpc.RepositoryReq{RepositoryId: req.RepositoryId})
	if err != nil {
		zap.S().Error("获取资源池信息失败 err:%s", err)
		return nil, err
	}
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	idInfo, err := l.svcCtx.UserCenterRepositoryRpc.CreateByShare(l.ctx, &userrepositoryrpc.CreateByShareReq{
		UserId:       userId,
		ParentId:     req.ParentId,
		RepositoryId: req.RepositoryId,
		Name:         nameInfo.Name,
	})
	if err != nil {
		zap.S().Error("创建资源失败 err:%s", err)
		return nil, err
	}
	return &types.ShareBasicSaveResponse{Id: strconv.FormatInt(idInfo.Id, 10)}, nil
}

package uploadservicerpclogic

import (
	"context"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteByIdLogic {
	return &DeleteByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteByIdLogic) DeleteById(in *pb.DeleteByIdReq) (*pb.DeleteByIdResp, error) {
	repositoryInfo, err := l.svcCtx.RepositoryPoolModel.FindOneByIdentity(l.ctx, uint64(in.RepositoryId))
	if err != nil {
		zap.S().Error("删除失败 err:%v", err)
		return nil, err
	}
	err = l.svcCtx.RepositoryPoolModel.DeleteByIdentity(l.ctx, uint64(in.RepositoryId))
	if err != nil {
		zap.S().Error("删除失败 err:%v", err)
		return nil, err
	}
	return &pb.DeleteByIdResp{Size: repositoryInfo.Size}, nil
}

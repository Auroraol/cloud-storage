package logic

import (
	"context"

	"github.com/Auroraol/cloud-storage/upload_service/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/rpc/pb"

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
	repositoryInfo, err := l.svcCtx.RepositoryPoolModel.FindOne(l.ctx, uint64(in.RepositoryId))
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.RepositoryPoolModel.Delete(l.ctx, uint64(in.RepositoryId))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteByIdResp{Size: repositoryInfo.Size}, nil
}

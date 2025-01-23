package logic

import (
	"context"

	"github.com/Auroraol/cloud-storage/upload_service/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRepositoryPoolByRepositoryIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRepositoryPoolByRepositoryIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRepositoryPoolByRepositoryIdLogic {
	return &GetRepositoryPoolByRepositoryIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRepositoryPoolByRepositoryIdLogic) GetRepositoryPoolByRepositoryId(in *pb.RepositoryReq) (*pb.RepositoryResp, error) {
	repositoryPoolInfo, err := l.svcCtx.RepositoryPoolModel.FindOne(l.ctx, uint64(in.RepositoryId))
	if err != nil {
		return nil, err
	}

	return &pb.RepositoryResp{Ext: repositoryPoolInfo.Ext,
		Size: repositoryPoolInfo.Size,
		Path: repositoryPoolInfo.Path,
		Name: repositoryPoolInfo.Name}, nil
}

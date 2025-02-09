package uploadservicelogic

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
	// todo: add your logic here and delete this line

	return &pb.RepositoryResp{}, nil
}

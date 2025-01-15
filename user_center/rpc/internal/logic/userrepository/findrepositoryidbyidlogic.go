package userrepositorylogic

import (
	"context"

	"github.com/Auroraol/cloud-storage/user_center/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindRepositoryIdByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindRepositoryIdByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindRepositoryIdByIdLogic {
	return &FindRepositoryIdByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindRepositoryIdByIdLogic) FindRepositoryIdById(in *pb.FindRepositoryIdReq) (*pb.FindRepositoryIdReply, error) {
	userRepositoryInfo, err := l.svcCtx.UserRepositoryModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, err
	}
	return &pb.FindRepositoryIdReply{RepositoryId: int64(userRepositoryInfo.RepositoryId)}, nil
}

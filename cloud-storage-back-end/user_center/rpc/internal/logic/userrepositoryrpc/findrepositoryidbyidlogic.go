package userrepositoryrpclogic

import (
	"context"
	"go.uber.org/zap"

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
	userRepositoryInfo, err := l.svcCtx.UserRepositoryModel.FindByRepositoryId(l.ctx, in.Id)
	if err != nil {
		zap.S().Error("UserRepositoryModel.FindByRepositoryId err:%v", err)
		return nil, err
	}
	return &pb.FindRepositoryIdReply{RepositoryId: int64(userRepositoryInfo.RepositoryId)}, nil
}

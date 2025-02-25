package userrepositoryrpclogic

import (
	"context"

	"github.com/Auroraol/cloud-storage/user_center/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserRepositoryNameByRepositoryIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserRepositoryNameByRepositoryIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserRepositoryNameByRepositoryIdLogic {
	return &GetUserRepositoryNameByRepositoryIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserRepositoryNameByRepositoryIdLogic) GetUserRepositoryNameByRepositoryId(in *pb.RepositoryIdReq) (*pb.UserRepositoryNameReply, error) {
	userInfo, err := l.svcCtx.UserRepositoryModel.FindByRepositoryId(l.ctx, in.RepositoryId)
	if err != nil {
		return nil, err
	}
	return &pb.UserRepositoryNameReply{RepositoryName: userInfo.Name}, nil
}

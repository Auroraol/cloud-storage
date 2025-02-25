package userrepositoryrpclogic

import (
	"context"
	"fmt"
	"github.com/Auroraol/cloud-storage/user_center/model"
	"github.com/bwmarrin/snowflake"

	"github.com/Auroraol/cloud-storage/user_center/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateByShareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateByShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateByShareLogic {
	return &CreateByShareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateByShareLogic) CreateByShare(in *pb.CreateByShareReq) (*pb.CreateByShareReply, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, fmt.Errorf("failed to create snowflake node: %w", err)
	}
	newId := node.Generate().Int64()
	_, err = l.svcCtx.UserRepositoryModel.InsertWithId(l.ctx, &model.UserRepository{
		Id:           uint64(newId),
		UserId:       uint64(in.UserId),
		ParentId:     uint64(in.ParentId),
		RepositoryId: uint64(in.RepositoryId),
		Name:         in.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateByShareReply{Id: newId}, nil
}

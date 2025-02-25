package userservicerpclogic

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"

	"github.com/Auroraol/cloud-storage/user_center/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddVolumeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddVolumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddVolumeLogic {
	return &AddVolumeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddVolumeLogic) AddVolume(in *pb.AddVolumeReq) (*pb.AddVolumeResp, error) {
	res, err := l.svcCtx.UserModel.UpdateVolume(l.ctx, in.Id, in.Size)
	if err != nil {
		return nil, err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, response.NewErrMsg("更新失败")
	}
	return &pb.AddVolumeResp{}, nil
}

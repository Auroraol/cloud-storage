package userlogic

import (
	"context"

	"github.com/Auroraol/cloud-storage/user_center/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindVolumeByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindVolumeByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindVolumeByIdLogic {
	return &FindVolumeByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindVolumeByIdLogic) FindVolumeById(in *pb.FindVolumeReq) (*pb.FindVolumeResp, error) {
	userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.FindVolumeResp{
		NowVolume:   userInfo.NowVolume,
		TotalVolume: userInfo.TotalVolume,
	}, nil
}

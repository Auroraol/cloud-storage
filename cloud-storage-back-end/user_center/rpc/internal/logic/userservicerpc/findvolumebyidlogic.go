package userservicerpclogic

import (
	"context"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/pb"

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
	userVolume, err := l.svcCtx.UserModel.FindVolume(l.ctx, in.Id)
	if err != nil {
		zap.S().Error("UserModel.FindVolume err:%v", err)
		return nil, err
	}

	return &pb.FindVolumeResp{
		NowVolume:   userVolume.NowVolume,
		TotalVolume: userVolume.TotalVolume,
	}, nil
}

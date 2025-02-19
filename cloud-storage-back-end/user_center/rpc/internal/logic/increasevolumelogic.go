package logic

import (
	"context"

	"github.com/Auroraol/cloud-storage/user_center/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type IncreaseVolumeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIncreaseVolumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IncreaseVolumeLogic {
	return &IncreaseVolumeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IncreaseVolumeLogic) IncreaseVolume(in *pb.IncreaseVolumeReq) (*pb.IncreaseVolumeResp, error) {
	// 获取用户信息
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, err
	}

	// 增加用户已使用容量
	user.UsedVolume += uint64(in.Size)

	// 更新用户信息
	err = l.svcCtx.UserModel.Update(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.IncreaseVolumeResp{
		Success: true,
	}, nil
}

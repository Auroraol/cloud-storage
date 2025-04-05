package userservicerpclogic

import (
	"context"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DecreaseVolumeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecreaseVolumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecreaseVolumeLogic {
	return &DecreaseVolumeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DecreaseVolumeLogic) DecreaseVolume(in *pb.DecreaseVolumeReq) (*pb.DecreaseVolumeResp, error) {
	// 减少用户容量
	res, err := l.svcCtx.UserModel.UpdateVolume(l.ctx, in.Id, -in.Size)
	if err != nil {
		zap.S().Error("UserModel.UpdateVolume err:%v", err)
		return nil, err
	}
	num, err := res.RowsAffected()
	if err != nil {
		zap.S().Error("UserModel.UpdateVolume err:%v", err)
		return nil, err
	}
	if num == 0 {
		zap.S().Error("删除失败！")
		return nil, response.NewErrMsg("删除失败！")
	}
	return &pb.DecreaseVolumeResp{}, nil
}

package logic

import (
	"context"

	"github.com/Auroraol/cloud-storage/user_center/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/rpc/pb"

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
	// todo: add your logic here and delete this line
	//res, err := l.svcCtx.UserBasicModel.UpdateVolume(l.ctx, in.Id, -in.Size)
	//if err != nil {
	//	return nil, err
	//}
	//num, err := res.RowsAffected()
	//if err != nil {
	//	return nil, err
	//}
	//if num == 0 {
	//	return nil, errorx.NewDefaultError("删除失败！")
	//}
	return &pb.DecreaseVolumeResp{}, nil
}

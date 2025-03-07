package sshservicerpclogic

import (
	"context"
	"github.com/Auroraol/cloud-storage/log_service/model"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/log_service/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveSshInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveSshInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveSshInfoLogic {
	return &SaveSshInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存ssh信息
func (l *SaveSshInfoLogic) SaveSshInfo(in *pb.SshInfoReq) (*pb.SshInfoResp, error) {
	// 保存 SSH 信息到数据库
	res, err := l.svcCtx.SshInfoModel.Insert(l.ctx, &model.SshInfo{
		UserId:   in.UserId,
		Host:     in.Host,
		Port:     int64(in.Port),
		Username: in.User,
		Password: in.Password,
	})
	if err != nil {
		zap.S().Errorf("保存ssh信息失败: %s", err)
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		zap.S().Errorf("获取ssh信息ID失败: %s", err)
		return nil, err
	}

	return &pb.SshInfoResp{Success: true, Message: "保存ssh信息成功", Id: id}, nil
}

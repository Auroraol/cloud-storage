package sshservicerpclogic

import (
	"context"
	"go.uber.org/zap"
	"strconv"

	"github.com/Auroraol/cloud-storage/log_service/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSshInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSshInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSshInfoLogic {
	return &GetSshInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询ssh信息
func (l *GetSshInfoLogic) GetSshInfo(in *pb.GetSshInfosReq) (*pb.SshInfoListResp, error) {
	// 查询 SSH 信息列表
	sshInfos, err := l.svcCtx.SshInfoModel.FindAll(l.ctx, strconv.FormatInt(in.UserId, 10))
	if err != nil {
		zap.S().Errorf("查询 SSH 信息列表失败: %s", err)
		return nil, err
	}

	var resp pb.SshInfoListResp
	for _, sshInfo := range sshInfos {
		resp.Items = append(resp.Items, &pb.SshInfoDetailResp{
			SshId:    sshInfo.Id,
			UserId:   sshInfo.UserId,
			Host:     sshInfo.Host,
			Port:     int32(sshInfo.Port),
			User:     sshInfo.Username,
			Password: sshInfo.Password,
		})
	}

	return &resp, nil
}

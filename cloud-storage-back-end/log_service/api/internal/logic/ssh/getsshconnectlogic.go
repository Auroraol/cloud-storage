package ssh

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	"github.com/Auroraol/cloud-storage/log_service/rpc/client/sshservicerpc"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSSHConnectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取SSH连接信息
func NewGetSSHConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSSHConnectLogic {
	return &GetSSHConnectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSSHConnectLogic) GetSSHConnect(req *types.GetSSHConnectReq) (resp *types.SshInfoListResp, err error) {
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	res, err := l.svcCtx.SshServiceRpc.GetSshInfo(l.ctx, &sshservicerpc.GetSshInfosReq{
		UserId: userId,
	})
	if err != nil {
		zap.S().Error("获取ssh信息失败 SshServiceRpc err:%v", err)
		return nil, err
	}
	var items []*types.SshInfoDetailResp
	for _, item := range res.Items {
		items = append(items, &types.SshInfoDetailResp{
			UserId:   item.UserId,
			SshId:    item.SshId,
			Host:     item.Host,
			Port:     item.Port,
			User:     item.User,
			Password: item.Password,
		})
	}

	return &types.SshInfoListResp{
		Items: items,
	}, nil
}

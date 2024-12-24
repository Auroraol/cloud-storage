package logic

import (
	"context"

	"cloud-storage/log_service/api/internal/svc"
	"cloud-storage/log_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 路径文件
func NewPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PathLogic {
	return &PathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PathLogic) Path(req *types.GetPathsFileReq) (resp *types.GetPathsFileRes, err error) {
	// todo: add your logic here and delete this line

	// 检测logfile表中, 有没有该文件

	return
}

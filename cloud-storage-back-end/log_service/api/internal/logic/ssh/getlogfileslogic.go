package ssh

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogFilesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取日志文件列表
func NewGetLogFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogFilesLogic {
	return &GetLogFilesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLogFilesLogic) GetLogFiles(req *types.GetLogFilesReq) (resp *types.GetLogFilesRes, err error) {
	// 参数校验
	if req.Host == "" {
		zap.S().Errorf("主机地址不能为空")
		return &types.GetLogFilesRes{
			Files:   []string{},
			Success: false,
		}, fmt.Errorf("主机地址不能为空")
	}

	if req.Path == "" {
		zap.S().Errorf("日志路径不能为空")
		return &types.GetLogFilesRes{
			Files:   []string{},
			Success: false,
		}, fmt.Errorf("日志路径不能为空")
	}

	// 获取日志文件列表
	files, err := l.svcCtx.SSHService.GetLogFiles(req.Path)
	if err != nil {
		zap.S().Errorf("获取日志文件列表失败: %v", err)
		return &types.GetLogFilesRes{
			Files:   []string{},
			Success: false,
		}, fmt.Errorf("获取日志文件列表失败: %v", err)
	}

	return &types.GetLogFilesRes{
		Files:   files,
		Success: true,
	}, nil
}

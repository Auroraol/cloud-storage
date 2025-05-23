package ssh

import (
	"context"
	"fmt"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 读取日志文件
func NewReadLogFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogFileLogic {
	return &ReadLogFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogFileLogic) ReadLogFile(req *types.ReadLogFileReq) (resp *types.ReadLogFileRes, err error) {
	// 参数校验
	if req.Host == "" {
		zap.S().Errorf("主机地址不能为空")
		return &types.ReadLogFileRes{
			Contents:   []string{},
			TotalLines: 0,
			Page:       req.Page,
			PageSize:   req.PageSize,
			Success:    false,
		}, fmt.Errorf("主机地址不能为空")
	}

	if req.Path == "" {
		zap.S().Errorf("日志路径不能为空")
		return &types.ReadLogFileRes{
			Contents:   []string{},
			TotalLines: 0,
			Page:       req.Page,
			PageSize:   req.PageSize,
			Success:    false,
		}, fmt.Errorf("日志路径不能为空")
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 读取日志文件
	contents, totalLines, err := l.svcCtx.SSHService.ReadLogFile(req.Path, req.Match, req.Page, req.PageSize)
	if err != nil {
		zap.S().Errorf("读取日志文件失败: %v", err)
		return &types.ReadLogFileRes{
			Contents:   []string{},
			TotalLines: 0,
			Page:       req.Page,
			PageSize:   req.PageSize,
			Success:    false,
		}, response.NewErrMsg("读取日志文件失败")
	}

	return &types.ReadLogFileRes{
		Contents:   contents,
		TotalLines: totalLines,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Success:    true,
	}, nil
}

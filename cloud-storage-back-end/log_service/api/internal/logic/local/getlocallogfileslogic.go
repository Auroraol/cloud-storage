package local

import (
	"context"
	"time"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLocalLogFilesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLocalLogFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLocalLogFilesLogic {
	return &GetLocalLogFilesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLocalLogFilesLogic) GetLocalLogFiles(req *types.GetLocalLogFilesReq) (resp *types.GetLocalLogFilesRes, err error) {
	// 获取文件列表
	files, err := l.svcCtx.LocalFileService.ListFiles(req.Path)
	if err != nil {
		return nil, err
	}

	// 获取文件统计
	stat, err := l.svcCtx.LocalFileService.GetFileStat(req.Path)
	if err != nil {
		return nil, err
	}

	// 转换文件信息为API类型
	fileInfos := make([]types.LocalFileInfo, 0, len(files))
	for _, file := range files {
		fileInfos = append(fileInfos, types.LocalFileInfo{
			Path:      file.Path,
			Name:      file.Name,
			Size:      file.Size,
			IsDir:     file.IsDir,
			ModTime:   file.ModTime.Format(time.RFC3339),
			Extension: file.Extension,
		})
	}

	// 转换统计信息为API类型
	fileStat := types.FileStat{
		TotalFiles:     stat.TotalFiles,
		TotalDirs:      stat.TotalDirs,
		TotalSize:      stat.TotalSize,
		LogFileCount:   stat.LogFileCount,
		RecentModified: stat.RecentModified,
	}

	return &types.GetLocalLogFilesRes{
		Files: fileInfos,
		Stat:  fileStat,
	}, nil
}

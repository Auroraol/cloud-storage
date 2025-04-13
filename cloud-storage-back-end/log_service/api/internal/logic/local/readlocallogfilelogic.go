package local

import (
	"context"
	"path/filepath"
	"time"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/localfile"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLocalLogFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLocalLogFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLocalLogFileLogic {
	return &ReadLocalLogFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLocalLogFileLogic) ReadLocalLogFile(req *types.ReadLocalLogFileReq) (resp *types.ReadLocalLogFileRes, err error) {
	// 从完整路径中提取文件名
	fileName := filepath.Base(req.FilePath)

	// 构建过滤条件
	filter := localfile.LogFilter{
		Level:      req.Level,
		Keyword:    req.Keyword,
		MaxResults: req.MaxResults,
	}

	// 解析开始时间
	if req.StartTime != "" {
		t, err := time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			return nil, err
		}
		filter.StartTime = &t
	}

	// 解析结束时间
	if req.EndTime != "" {
		t, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			return nil, err
		}
		filter.EndTime = &t
	}

	// 读取日志文件
	entries, err := l.svcCtx.LocalFileService.ReadLogFile(fileName, filter)
	if err != nil {
		return nil, err
	}

	// 转换日志条目为API类型
	logEntries := make([]types.LocalLogEntry, 0, len(entries))
	for _, entry := range entries {
		logEntries = append(logEntries, types.LocalLogEntry{
			Timestamp: entry.Timestamp.Format(time.RFC3339),
			Level:     entry.Level,
			Content:   entry.Content,
			Source:    entry.Source,
			LineNum:   entry.LineNum,
		})
	}

	return &types.ReadLocalLogFileRes{
		Entries: logEntries,
		Total:   len(logEntries),
	}, nil
}

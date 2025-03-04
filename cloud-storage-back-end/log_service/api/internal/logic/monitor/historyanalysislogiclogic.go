package monitor

import (
	"context"
	"fmt"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

type HistoryAnalysisLogicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 历史分析

func NewHistoryAnalysisLogicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HistoryAnalysisLogicLogic {
	return &HistoryAnalysisLogicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HistoryAnalysisLogicLogic) HistoryAnalysisLogic(req *types.HistoryAnalysisReq) (resp *types.HistoryAnalysisRes, err error) {
	// 参数校验
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	if req.Host == "" {
		return nil, fmt.Errorf("主机地址不能为空")
	}

	if req.LogFile == "" {
		return nil, fmt.Errorf("日志文件不能为空")
	}

	// 如果没有指定时间范围，默认为最近24小时
	if req.StartTime <= 0 || req.EndTime <= 0 {
		req.EndTime = time.Now().Unix()
		req.StartTime = req.EndTime - 86400 // 24小时前
	}

	// 使用SSH服务读取日志文件
	contents, _, err := l.svcCtx.SSHService.ReadLogFile(req.LogFile, req.Keywords, req.Page, req.PageSize)
	if err != nil {
		// 如果读取失败，使用模拟数据
		l.Logger.Errorf("读取日志文件失败: %v", err)
		return nil, response.NewErrMsg("读取日志文件失败")
	}

	// 创建一个map来存储按时间戳分组的日志条目
	timestampToLogEntries := make(map[int64]*types.LogEntry)

	// 解析日志内容
	for _, content := range contents {
		// 简单解析日志行，实际应根据日志格式进行解析
		parts := strings.SplitN(content, " ", 3)
		timestamp := time.Now().Unix()
		level := "INFO"
		message := content

		if len(parts) >= 3 {
			// 尝试解析时间戳
			if t, err := time.Parse("2006-01-02T15:04:05", parts[0]); err == nil {
				timestamp = t.Unix()
			}

			// 尝试解析日志级别
			level = parseLogLevel(parts[1])

			message = parts[2]
		}

		// 如果相同时间戳的日志条目已经存在，合并它们并累加Value
		if existingEntry, exists := timestampToLogEntries[timestamp]; exists {
			// 合并日志内容
			existingEntry.Content += " | " + message

			// 选择最严重的日志级别
			if compareLogLevels(level, existingEntry.Level) > 0 {
				existingEntry.Level = level
			}

			// 累加Value
			existingEntry.Value += 1
		} else {
			// 如果该时间戳日志条目还不存在，直接添加
			timestampToLogEntries[timestamp] = &types.LogEntry{
				Timestamp: timestamp,
				Content:   message,
				Level:     level,
				Value:     1,
			}
		}
	}

	// 将合并后的日志条目添加到结果数据中
	data := make([]types.LogEntry, 0, len(timestampToLogEntries))
	for _, entry := range timestampToLogEntries {
		data = append(data, *entry)
	}

	return &types.HistoryAnalysisRes{
		Data:     data,
		Total:    len(data),
		Page:     req.Page,
		PageSize: req.PageSize,
		Success:  true,
	}, nil
}

// 比较日志级别的严重程度，返回1表示level1比level2更严重，-1表示level1比level2轻，0表示一样
func compareLogLevels(level1, level2 string) int {
	// 日志级别的严重程度顺序：CRITICAL > ERROR > WARN > INFO > DEBUG
	levels := map[string]int{
		"CRITICAL": 5,
		"ERROR":    4,
		"WARN":     3,
		"INFO":     2,
		"DEBUG":    1,
	}

	if levels[level1] > levels[level2] {
		return 1
	} else if levels[level1] < levels[level2] {
		return -1
	} else {
		return 0
	}
}

// 创建日志级别映射表（可扩展）
var levelMap = map[string]string{
	"error":    "ERROR",
	"err":      "ERROR",
	"warning":  "WARN",
	"warn":     "WARN",
	"debug":    "DEBUG",
	"dbg":      "DEBUG",
	"info":     "INFO",
	"notice":   "INFO",
	"critical": "CRITICAL",
}

func parseLogLevel(s string) string {
	// 清理字符串并转换为小写
	cleanStr := strings.ToLower(strings.Trim(s, "[]"))

	// 优先完全匹配
	if level, exists := levelMap[cleanStr]; exists {
		return level
	}

	// 部分匹配（针对带格式的日志）
	for k, v := range levelMap {
		if strings.Contains(cleanStr, k) {
			return v
		}
	}

	// 默认为INFO
	return "INFO"
}

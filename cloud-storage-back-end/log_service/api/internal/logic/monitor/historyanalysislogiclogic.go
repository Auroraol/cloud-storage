package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"sort"
	"strings"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/time"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
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
		return nil, response.NewErrCodeMsg(response.DATA_PARAM_ERROR, "主机地址不能为空")
	}

	if req.LogFile == "" {
		return nil, response.NewErrCodeMsg(response.DATA_PARAM_ERROR, "日志文件不能为空")
	}

	// 如果没有指定时间范围，默认为最近24小时
	if req.StartTime <= 0 || req.EndTime <= 0 {
		//req.EndTime = time.Now().Unix()
		req.StartTime = req.EndTime - 86400 // 24小时前
	}

	// 使用SSH服务读取日志文件
	contents, _, err := l.svcCtx.SSHService.ReadLogFile(req.LogFile, req.Keywords, req.Page, req.PageSize)
	if err != nil {
		// 如果读取失败，使用模拟数据
		zap.S().Errorf("读取日志文件失败: %v", err)
		return nil, response.NewErrMsg("读取日志文件失败")
	}

	// 创建一个map来存储按时间戳分组的日志条目
	timestampToLogEntries := make(map[int64]*types.LogEntry)

	// 解析日志内容
	for _, content := range contents {
		// 尝试解析JSON格式日志
		timestamp, level, message, fields, err := parseJSONLog(content)
		if err != nil {
			zap.S().Errorf("解析日志内容失败: %v", err)
			return nil, response.NewErrMsg("解析日志内容失败")
		}

		// JSON解析成功，如果有额外字段，将它们添加到消息中
		if len(fields) > 0 {
			extraFields := ""
			for k, v := range fields {
				extraFields += fmt.Sprintf(" %s=%s", k, v)
			}
			if extraFields != "" {
				message += extraFields
			}
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

	// 按时间戳排序
	sort.Slice(data, func(i, j int) bool {
		return data[i].Timestamp < data[j].Timestamp
	})

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

// 解析JSON格式的日志
func parseJSONLog(content string) (timestamp int64, level string, message string, fields map[string]string, err error) {
	// 初始化返回值
	fields = make(map[string]string)

	// 解析JSON
	var logEntry map[string]interface{}
	err = json.Unmarshal([]byte(content), &logEntry)
	if err != nil {
		return 0, "INFO", content, fields, err
	}

	// 提取日志级别
	if levelValue, exists := logEntry["level"]; exists {
		if levelStr, ok := levelValue.(string); ok {
			level = parseLogLevel(levelStr)
		}
	}

	// 提取时间
	if timeValue, exists := logEntry["time"]; exists {
		if timeStr, ok := timeValue.(string); ok {
			// 尝试解析时间
			//time -> 2025-03-04 23:22:53
			toTimestamp, err := time.StringTimeToTimestamp(timeStr)
			if err == nil {
				timestamp = toTimestamp
			} else {
				// 如果解析失败，使用当前时间
				timestamp = time.LocalTimeNow().Unix()
			}
		}
	}

	// 提取消息
	if msgValue, exists := logEntry["msg"]; exists {
		if msgStr, ok := msgValue.(string); ok {
			message = msgStr
		}
	}

	// 提取其他字段
	for key, value := range logEntry {
		if key != "level" && key != "time" && key != "msg" {
			switch v := value.(type) {
			case string:
				fields[key] = v
			default:
				// 对于非字符串类型，转换为字符串
				fields[key] = fmt.Sprintf("%v", v)
			}
		}
	}

	return timestamp, level, message, fields, nil
}

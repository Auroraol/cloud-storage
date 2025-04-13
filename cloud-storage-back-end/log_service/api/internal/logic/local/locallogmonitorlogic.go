package local

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	stdtime "time"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/time"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LocalLogMonitorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 本地日志监控
func NewLocalLogMonitorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LocalLogMonitorLogic {
	return &LocalLogMonitorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LocalLogMonitorLogic) LocalLogMonitor(req *types.LocalRealTimeMonitorReq) (resp *types.LocalRealTimeMonitorRes, err error) {
	// 检查文件是否存在
	if _, err := os.Stat(req.LogFile); os.IsNotExist(err) {
		return nil, err
	}

	// 读取文件内容
	content, err := os.ReadFile(req.LogFile)
	if err != nil {
		return nil, err
	}

	// 计算时间范围
	now := time.LocalTimeNow()
	startTime := now.Add(-stdtime.Duration(req.TimeRange) * stdtime.Hour)

	// 使用双层 map 进行计数 [监控项][时间戳]计数
	counter := make(map[string]map[int64]int)
	for _, item := range req.MonitorItems {
		counter[item] = make(map[int64]int)
	}

	// 创建一个map来存储调用者统计信息
	callerStats := make(map[string]int)

	// 日志处理
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		// 尝试解析为JSON格式
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err == nil {
			// 提取时间
			var logTime stdtime.Time
			var toTimestamp int64
			if timeStr, ok := logEntry["time"].(string); ok {
				// 尝试解析时间，支持多种格式
				toTimestamp, err = time.StringTimeToTimestamp(timeStr)
				if err != nil {
					// 如果解析失败，使用当前时间
					logTime = time.LocalTimeNow()
				} else {
					logTime = stdtime.Unix(toTimestamp, 0)
				}
			} else {
				// 如果没有时间字段，使用当前时间
				logTime = time.LocalTimeNow()
			}

			// 检查日志是否在指定的时间范围内
			if logTime.Before(startTime) {
				continue // 跳过不在时间范围内的日志
			}

			// 提取日志级别
			var level string
			if levelValue, ok := logEntry["level"].(string); ok {
				level = strings.ToUpper(levelValue)
			}

			// 提取消息
			var message string
			if msgValue, ok := logEntry["msg"].(string); ok {
				message = msgValue
			}

			// 提取调用者信息并统计
			if callerValue, ok := logEntry["caller"].(string); ok {
				// 统计调用者
				callerStats[callerValue]++
			}

			// 按分钟取整
			timestamp := toTimestamp

			// 按照监控项进行统计
			for _, item := range req.MonitorItems {
				switch item {
				case "requests":
					// 统计所有请求
					if level == "REQUESTS" {
						counter["requests"][timestamp]++
					}
				case "errors":
					// 统计错误日志
					if level == "ERROR" {
						counter["errors"][timestamp]++
					}
				case "response_time":
					// 尝试从消息中提取响应时间
					if responseTime, err := extractResponseTime(message); err == nil {
						counter["response_time"][timestamp] += responseTime
					}
				case "debug_logs":
					// 统计调试日志
					if level == "DEBUG" {
						counter["debug_logs"][timestamp]++
					}
				case "warn_logs":
					// 统计警告日志
					if level == "WARN" {
						counter["warn_logs"][timestamp]++
					}
				case "info_logs":
					// 统计信息日志
					if level == "INFO" {
						counter["info_logs"][timestamp]++
					}
				}
			}
		}
	}

	// 构建监控数据
	monitorData := make([]types.LocalMonitorData, 0)
	for item, timestamps := range counter {
		for ts, count := range timestamps {
			monitorData = append(monitorData, types.LocalMonitorData{
				Timestamp: ts,
				Value:     count,
				Type:      item,
				Caller:    req.LogFile,
			})
		}
	}

	// 按时间戳排序
	sort.Slice(monitorData, func(i, j int) bool {
		return monitorData[i].Timestamp < monitorData[j].Timestamp
	})

	return &types.LocalRealTimeMonitorRes{
		Data:    monitorData,
		Total:   len(monitorData),
		Success: true,
	}, nil
}

// 辅助函数：从日志行提取响应时间
func extractResponseTime(line string) (int, error) {
	// 假设日志格式包含 response_time=123ms
	re := regexp.MustCompile(`response_time=(\d+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 2 {
		return 0, fmt.Errorf("未找到响应时间")
	}

	val, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return val, nil
}

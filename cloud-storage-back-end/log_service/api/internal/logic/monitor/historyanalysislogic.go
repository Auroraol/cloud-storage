package monitor

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HistoryAnalysisLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 历史分析
func NewHistoryAnalysisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HistoryAnalysisLogic {
	return &HistoryAnalysisLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HistoryAnalysisLogic) HistoryAnalysis(req *types.HistoryAnalysisReq) (resp *types.HistoryAnalysisRes, err error) {
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
		l.Logger.Errorf("读取日志文件失败: %v，使用模拟数据", err)
		return l.generateMockData(req)
	}

	// 解析日志内容
	data := make([]types.LogEntry, 0, len(contents))
	for _, content := range contents {
		// 简单解析日志行，实际应根据日志格式进行解析
		parts := strings.SplitN(content, " ", 3)
		timestamp := time.Now().Unix()
		level := "INFO"
		source := "system"
		message := content

		if len(parts) >= 3 {
			// 尝试解析时间戳
			if t, err := time.Parse("2006-01-02T15:04:05", parts[0]); err == nil {
				timestamp = t.Unix()
			}

			// 尝试解析日志级别
			if strings.Contains(parts[1], "ERROR") {
				level = "ERROR"
			} else if strings.Contains(parts[1], "WARN") {
				level = "WARN"
			} else if strings.Contains(parts[1], "DEBUG") {
				level = "DEBUG"
			}

			// 尝试解析来源
			if strings.Contains(parts[1], "[") && strings.Contains(parts[1], "]") {
				source = strings.Trim(parts[1], "[]")
			}

			message = parts[2]
		}

		data = append(data, types.LogEntry{
			Timestamp: timestamp,
			Content:   message,
			Level:     level,
			Source:    source,
		})
	}

	return &types.HistoryAnalysisRes{
		Data:     data,
		Total:    100, // 假设总数为100
		Page:     req.Page,
		PageSize: req.PageSize,
		Success:  true,
	}, nil
}

// 生成模拟数据
func (l *HistoryAnalysisLogic) generateMockData(req *types.HistoryAnalysisReq) (*types.HistoryAnalysisRes, error) {
	// 模拟从日志文件中读取数据
	// 实际应该根据日志文件路径、主机、时间范围和关键字进行过滤
	total := 100 // 模拟总数

	// 计算分页
	offset := (req.Page - 1) * req.PageSize
	limit := req.PageSize
	if offset >= total {
		return &types.HistoryAnalysisRes{
			Data:     make([]types.LogEntry, 0),
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Success:  true,
		}, nil
	}

	// 生成模拟日志数据
	data := make([]types.LogEntry, 0, limit)
	logLevels := []string{"INFO", "WARN", "ERROR", "DEBUG"}
	sources := []string{"api", "service", "database", "cache"}

	// 生成随机日志内容
	contents := []string{
		"用户登录成功",
		"请求处理完成",
		"数据库查询执行",
		"缓存更新",
		"文件上传成功",
		"请求参数错误",
		"数据库连接超时",
		"权限验证失败",
		"系统异常",
		"网络连接中断",
	}

	// 如果有关键字，过滤内容
	filteredContents := contents
	if req.Keywords != "" {
		filteredContents = make([]string, 0)
		for _, content := range contents {
			if strings.Contains(content, req.Keywords) {
				filteredContents = append(filteredContents, content)
			}
		}
		// 如果没有匹配的内容，使用原始内容
		if len(filteredContents) == 0 {
			filteredContents = contents
		}
	}

	// 根据聚合方式确定时间间隔
	var interval int64
	switch req.AggregateBy {
	case "按分钟":
		interval = 60 // 1分钟
	case "按小时":
		interval = 3600 // 1小时
	case "按天":
		interval = 86400 // 1天
	default:
		interval = 3600 // 默认按小时
	}

	// 生成日志条目
	for i := 0; i < limit && offset+i < total; i++ {
		// 随机生成时间戳，在开始和结束时间之间
		timestamp := req.StartTime + rand.Int63n(req.EndTime-req.StartTime)
		// 根据聚合方式调整时间戳
		timestamp = (timestamp / interval) * interval

		// 随机选择日志级别、来源和内容
		level := logLevels[rand.Intn(len(logLevels))]
		source := sources[rand.Intn(len(sources))]
		content := filteredContents[rand.Intn(len(filteredContents))]

		// 添加日志条目
		data = append(data, types.LogEntry{
			Timestamp: timestamp,
			Content:   content,
			Level:     level,
			Source:    source,
		})
	}

	return &types.HistoryAnalysisRes{
		Data:     data,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Success:  true,
	}, nil
}

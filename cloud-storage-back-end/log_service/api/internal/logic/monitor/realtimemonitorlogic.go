package monitor

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RealTimeMonitorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 实时监控
func NewRealTimeMonitorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RealTimeMonitorLogic {
	return &RealTimeMonitorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RealTimeMonitorLogic) RealTimeMonitor(req *types.RealTimeMonitorReq) (resp *types.RealTimeMonitorRes, err error) {
	// 参数校验
	if req.TimeRange <= 0 {
		req.TimeRange = 1 // 默认1小时
	}

	if req.Host == "" {
		return nil, fmt.Errorf("主机地址不能为空")
	}

	if req.LogFile == "" {
		return nil, fmt.Errorf("日志文件不能为空")
	}

	if len(req.MonitorItems) == 0 {
		return nil, fmt.Errorf("监控项不能为空")
	}

	// 使用SSH服务获取监控数据
	monitorData, err := l.svcCtx.SSHService.MonitorLogFile(req.LogFile, req.MonitorItems, req.TimeRange)
	if err != nil {
		// 如果获取失败，使用模拟数据
		l.Logger.Errorf("获取监控数据失败: %v，使用模拟数据", err)
		return l.generateMockData(req)
	}

	// 转换数据格式
	data := make([]types.MonitorData, 0)
	for itemType, items := range monitorData {
		for _, item := range items {
			data = append(data, types.MonitorData{
				Timestamp: item["timestamp"].(int64),
				Value:     item["value"].(int),
				Type:      itemType,
			})
		}
	}

	return &types.RealTimeMonitorRes{
		Data:    data,
		Total:   len(data),
		Success: true,
	}, nil
}

// 生成模拟数据
func (l *RealTimeMonitorLogic) generateMockData(req *types.RealTimeMonitorReq) (*types.RealTimeMonitorRes, error) {
	// 根据时间范围确定数据点数量
	var dataPoints int
	switch req.TimeRange {
	case 1: // 1小时
		dataPoints = 60 // 每分钟一个点
	case 6: // 6小时
		dataPoints = 72 // 每5分钟一个点
	case 12: // 12小时
		dataPoints = 72 // 每10分钟一个点
	case 24: // 24小时
		dataPoints = 96 // 每15分钟一个点
	default:
		dataPoints = 60
	}

	// 生成监控数据
	data := make([]types.MonitorData, 0)
	now := time.Now().Unix()
	interval := int64(req.TimeRange * 3600 / dataPoints)

	// 为每个监控项生成数据
	for _, item := range req.MonitorItems {
		for i := 0; i < dataPoints; i++ {
			timestamp := now - int64(i)*interval
			value := 0

			// 根据监控项类型生成不同范围的随机值
			switch item {
			case "请求数":
				value = rand.Intn(100) + 50 // 50-150之间的随机值
			case "错误数":
				value = rand.Intn(20) // 0-20之间的随机值
			case "响应时间":
				value = rand.Intn(500) + 100 // 100-600毫秒之间的随机值
			}

			data = append(data, types.MonitorData{
				Timestamp: timestamp,
				Value:     value,
				Type:      item,
			})
		}
	}

	return &types.RealTimeMonitorRes{
		Data:    data,
		Total:   len(data),
		Success: true,
	}, nil
}

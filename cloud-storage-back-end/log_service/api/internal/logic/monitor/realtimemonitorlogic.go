package monitor

import (
	"context"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/time"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"

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
		return nil, response.NewErrCodeMsg(response.DATA_PARAM_ERROR, "主机地址不能为空")
	}

	if req.LogFile == "" {
		return nil, response.NewErrCodeMsg(response.DATA_PARAM_ERROR, "日志文件不能为空")
	}

	if len(req.MonitorItems) == 0 {
		return nil, response.NewErrCodeMsg(response.DATA_PARAM_ERROR, "监控项不能为空")
	}

	// 使用SSH服务获取监控数据
	monitorData, err := l.svcCtx.SSHService.MonitorLogFile(req.LogFile, req.MonitorItems, req.TimeRange)
	if err != nil {
		// 如果获取失败，记录错误
		zap.S().Errorf("获取监控数据失败: %v", err)
		return nil, response.NewErrMsg("获取监控数据失败")
	}

	// 转换数据格式
	data := make([]types.MonitorData, 0)
	for itemType, items := range monitorData {
		// 如果是调用者统计信息，需要特殊处理
		if itemType == "caller_stats" {
			// 处理调用者统计信息
			callerData := processCallerStats(items)
			data = append(data, callerData...)
			continue
		}

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

// 处理调用者统计信息
func processCallerStats(items []map[string]interface{}) []types.MonitorData {
	result := make([]types.MonitorData, 0, len(items))

	// 当前时间戳，用于所有调用者统计
	now := time.LocalTimeNow().Unix()

	for _, item := range items {
		caller, ok := item["caller"].(string)
		if !ok {
			continue
		}

		value, ok := item["value"].(int)
		if !ok {
			// 尝试转换为其他数字类型
			if floatVal, ok := item["value"].(float64); ok {
				value = int(floatVal)
			} else {
				continue
			}
		}

		result = append(result, types.MonitorData{
			Timestamp: now,
			Value:     value,
			Type:      "caller_stats",
			Caller:    caller,
		})
	}

	return result
}

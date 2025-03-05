package monitor

import (
	"github.com/Auroraol/cloud-storage/common/logx"
	"go.uber.org/zap"
	"net/http"

	"github.com/Auroraol/cloud-storage/common/response"

	"github.com/Auroraol/cloud-storage/log_service/api/internal/logic/monitor"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 实时监控
func RealTimeMonitorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.LogWithCustomLevel("requests", r.Host+" ["+r.RequestURI+"]")
		var req types.RealTimeMonitorReq
		if err := httpx.Parse(r, &req); err != nil {
			zap.S().Errorf("parse param error: %v", err)
			response.ParamErrorResult(r, w, err)
			return
		}

		l := monitor.NewRealTimeMonitorLogic(r.Context(), svcCtx)
		resp, err := l.RealTimeMonitor(&req)
		response.HttpResult(r, w, resp, err)
	}
}

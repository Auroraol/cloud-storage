package local

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/logx"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"go.uber.org/zap"
	"net/http"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/logic/local"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 本地日志监控
func LocalLogMonitorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LocalRealTimeMonitorReq
		logx.LogWithCustomLevel("requests", r.Host+" ["+r.RequestURI+"]")
		if err := httpx.Parse(r, &req); err != nil {
			zap.S().Errorf("parse param error: %v", err)
			response.ParamErrorResult(r, w, err)
			return
		}

		l := local.NewLocalLogMonitorLogic(r.Context(), svcCtx)
		resp, err := l.LocalLogMonitor(&req)
		response.HttpResult(r, w, resp, err)

	}
}

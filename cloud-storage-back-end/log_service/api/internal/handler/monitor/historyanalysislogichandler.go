package monitor

import (
	"github.com/Auroraol/cloud-storage/common/response"
	"net/http"

	"github.com/Auroraol/cloud-storage/log_service/api/internal/logic/monitor"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 历史分析
func HistoryAnalysisLogicHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HistoryAnalysisReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := monitor.NewHistoryAnalysisLogicLogic(r.Context(), svcCtx)
		resp, err := l.HistoryAnalysisLogic(&req)
		response.HttpResult(r, w, resp, err)
	}
}

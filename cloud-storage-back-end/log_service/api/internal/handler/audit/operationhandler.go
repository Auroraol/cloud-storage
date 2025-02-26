package audit

import (
	"github.com/Auroraol/cloud-storage/common/response"
	"net/http"

	"github.com/Auroraol/cloud-storage/log_service/api/internal/logic/audit"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获得操作日志
func OperationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetOperationLogReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := audit.NewOperationLogic(r.Context(), svcCtx)
		resp, err := l.Operation(&req)
		response.HttpResult(r, w, resp, err)
	}
}

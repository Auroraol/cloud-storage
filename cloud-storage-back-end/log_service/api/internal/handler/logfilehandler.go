package handler

import (
	"net/http"

	"cloud-storage/log_service/api/internal/logic"
	"cloud-storage/log_service/api/internal/svc"
	"cloud-storage/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 日志文件
func logfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetLogfileReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := logic.NewLogfileLogic(r.Context(), svcCtx)
		resp, err := l.Logfile(&req)
		response.HttpResult(r, w, resp, err)
	}
}

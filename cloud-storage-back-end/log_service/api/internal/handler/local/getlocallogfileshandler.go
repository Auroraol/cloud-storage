package local

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"net/http"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/logic/local"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetLocalLogFilesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetLocalLogFilesReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := local.NewGetLocalLogFilesLogic(r.Context(), svcCtx)
		resp, err := l.GetLocalLogFiles(&req)
		response.HttpResult(r, w, resp, err)
	}
}

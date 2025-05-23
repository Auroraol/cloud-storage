package local

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"net/http"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/logic/local"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ReadLocalLogFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReadLocalLogFileReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := local.NewReadLocalLogFileLogic(r.Context(), svcCtx)
		resp, err := l.ReadLocalLogFile(&req)
		response.HttpResult(r, w, resp, err)
	}
}

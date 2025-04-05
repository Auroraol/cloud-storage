package history

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"net/http"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/logic/history"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新历史上传记录
func UpdateHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateHistoryRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := history.NewUpdateHistoryLogic(r.Context(), svcCtx)
		resp, err := l.UpdateHistory(&req)
		response.HttpResult(r, w, resp, err)
	}
}

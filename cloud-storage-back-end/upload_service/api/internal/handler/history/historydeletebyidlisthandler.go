package history

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"net/http"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/logic/history"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除所有历史记录
func HistoryDeleteByIdListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HistoryDeleteAllRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := history.NewHistoryDeleteByIdListLogic(r.Context(), svcCtx)
		resp, err := l.HistoryDeleteByIdList(&req)
		response.HttpResult(r, w, resp, err)
	}
}

package handler

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/logx"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"net/http"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/share_service/api/internal/logic"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/share_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/share_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 用户分享列表
func ShareBasicListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.LogWithCustomLevel("requests", r.Host+" ["+r.RequestURI+"]")
		var req types.ShareBasicListRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := logic.NewShareBasicListLogic(r.Context(), svcCtx)
		resp, err := l.ShareBasicList(&req)
		response.HttpResult(r, w, resp, err)
	}
}

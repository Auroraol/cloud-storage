package handler

import (
	"github.com/Auroraol/cloud-storage/common/logx"
	"github.com/Auroraol/cloud-storage/common/response"
	"net/http"

	"github.com/Auroraol/cloud-storage/share_service/api/internal/logic"
	"github.com/Auroraol/cloud-storage/share_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/share_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取资源详情
func ShareBasicDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.LogWithCustomLevel("requests", r.Host+" ["+r.RequestURI+"]")
		var req types.DetailRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := logic.NewShareBasicDetailLogic(r.Context(), svcCtx)
		resp, err := l.ShareBasicDetail(&req)
		response.HttpResult(r, w, resp, err)
	}
}

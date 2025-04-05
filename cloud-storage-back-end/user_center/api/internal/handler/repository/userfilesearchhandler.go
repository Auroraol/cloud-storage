package repository

import (
	"net/http"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/logic/repository"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 搜索用户文件和文件夹
func UserFileSearchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFileSearchRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := repository.NewUserFileSearchLogic(r.Context(), svcCtx)
		resp, err := l.UserFileSearch(&req)
		response.HttpResult(r, w, resp, err)
	}
}

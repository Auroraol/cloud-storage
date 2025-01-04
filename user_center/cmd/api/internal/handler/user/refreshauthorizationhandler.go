package user

import (
	"net/http"

	"cloud-storage/user_center/cmd/api/internal/logic/user"
	"cloud-storage/user_center/cmd/api/internal/svc"
	"cloud-storage/user_center/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 刷新Authorization
func RefreshAuthorizationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshAuthRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewRefreshAuthorizationLogic(r.Context(), svcCtx)
		resp, err := l.RefreshAuthorization(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

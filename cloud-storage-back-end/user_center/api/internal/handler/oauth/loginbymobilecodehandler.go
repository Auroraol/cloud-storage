package oauth

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"net/http"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/logic/oauth"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 手机号登陆/注册
func LoginByMobileCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MobileLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := oauth.NewLoginByMobileCodeLogic(r.Context(), svcCtx)
		resp, err := l.LoginByMobileCode(&req)
		response.HttpResult(r, w, resp, err)
	}
}

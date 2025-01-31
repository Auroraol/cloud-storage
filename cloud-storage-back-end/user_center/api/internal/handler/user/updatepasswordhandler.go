package user

import (
	"github.com/Auroraol/cloud-storage/common/response"
	"net/http"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/logic/user"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 修改密码
func UpdatePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserPasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewUpdatePasswordLogic(r.Context(), svcCtx)
		resp, err := l.UpdatePassword(&req)
		response.HttpResult(r, w, resp, err)
	}
}

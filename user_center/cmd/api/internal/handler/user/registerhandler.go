package user

import (
	"net/http"

	"cloud-storage/common/result"
	"cloud-storage/user_center/cmd/api/internal/logic/user"
	"cloud-storage/user_center/cmd/api/internal/svc"
	"cloud-storage/user_center/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 用户注册
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		result.HttpResult(r, w, resp, err)
	}
}

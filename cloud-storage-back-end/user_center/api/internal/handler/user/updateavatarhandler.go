package user

import (
	"github.com/Auroraol/cloud-storage/common/response"
	"net/http"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/logic/user"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更换头像
func UpdateAvatarHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserAvatarReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}
		l := user.NewUpdateAvatarLogic(r.Context(), svcCtx)
		resp, err := l.UpdateAvatar(&req, file, fileHeader)
		response.HttpResult(r, w, resp, err)
	}
}

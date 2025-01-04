package user

import (
	"net/http"

	"cloud-storage/user_center/cmd/api/internal/logic/user"
	"cloud-storage/user_center/cmd/api/internal/svc"
	"cloud-storage/user_center/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 验证码发送
func CodeSendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CodeSendRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewCodeSendLogic(r.Context(), svcCtx)
		resp, err := l.CodeSend(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

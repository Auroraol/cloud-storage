package ssh

import (
	"net/http"

	"github.com/Auroraol/cloud-storage/log_service/api/internal/logic/ssh"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除SSH连接信息
func DeleteSSHConnectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteSSHConnectReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := ssh.NewDeleteSSHConnectLogic(r.Context(), svcCtx)
		resp, err := l.DeleteSSHConnect(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

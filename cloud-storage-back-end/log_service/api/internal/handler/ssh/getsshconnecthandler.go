package ssh

import (
	"net/http"

	"github.com/Auroraol/cloud-storage/log_service/api/internal/logic/ssh"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取SSH连接信息
func GetSSHConnectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSSHConnectReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := ssh.NewGetSSHConnectLogic(r.Context(), svcCtx)
		resp, err := l.GetSSHConnect(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

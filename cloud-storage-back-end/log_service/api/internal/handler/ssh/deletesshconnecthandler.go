package ssh

import (
	"github.com/Auroraol/cloud-storage/common/logx"
	"github.com/Auroraol/cloud-storage/common/response"
	"go.uber.org/zap"
	"net/http"

	"github.com/Auroraol/cloud-storage/log_service/api/internal/logic/ssh"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除SSH连接信息
func DeleteSSHConnectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.LogWithCustomLevel("requests", r.Host+" ["+r.RequestURI+"]")
		var req types.DeleteSSHConnectReq
		if err := httpx.Parse(r, &req); err != nil {
			zap.S().Errorf("parse param error: %v", err)
			response.ParamErrorResult(r, w, err)
			return
		}

		l := ssh.NewDeleteSSHConnectLogic(r.Context(), svcCtx)
		resp, err := l.DeleteSSHConnect(&req)
		response.HttpResult(r, w, resp, err)
	}
}

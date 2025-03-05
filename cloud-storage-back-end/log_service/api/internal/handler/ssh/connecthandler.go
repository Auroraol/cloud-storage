package ssh

import (
	"github.com/Auroraol/cloud-storage/common/logx"
	"go.uber.org/zap"
	"net/http"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/logic/ssh"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// SSH连接
func ConnectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.LogWithCustomLevel("requests", r.Host+" ["+r.RequestURI+"]")
		var req types.SSHConnectReq
		if err := httpx.Parse(r, &req); err != nil {
			zap.S().Errorf("parse param error: %v", err)
			response.ParamErrorResult(r, w, err)
			return
		}
		zap.S().Infof("req: %+v", req)

		l := ssh.NewConnectLogic(r.Context(), svcCtx)
		resp, err := l.Connect(&req)
		response.HttpResult(r, w, resp, err)
	}
}

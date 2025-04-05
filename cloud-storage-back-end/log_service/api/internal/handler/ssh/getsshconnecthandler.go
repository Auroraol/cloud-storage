package ssh

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/logx"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"go.uber.org/zap"
	"net/http"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/logic/ssh"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取SSH连接信息
func GetSSHConnectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.LogWithCustomLevel("requests", r.Host+" ["+r.RequestURI+"]")
		var req types.GetSSHConnectReq
		if err := httpx.Parse(r, &req); err != nil {
			zap.S().Errorf("parse param error: %v", err)
			response.ParamErrorResult(r, w, err)
			return
		}

		l := ssh.NewGetSSHConnectLogic(r.Context(), svcCtx)
		resp, err := l.GetSSHConnect(&req)
		response.HttpResult(r, w, resp, err)
	}
}

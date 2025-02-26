package ssh

import (
	"net/http"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/logic/ssh"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取日志文件列表
func GetLogFilesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetLogFilesReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := ssh.NewGetLogFilesLogic(r.Context(), svcCtx)
		resp, err := l.GetLogFiles(&req)
		response.HttpResult(r, w, resp, err)
	}
}

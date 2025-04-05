package handler

import (
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/logx"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/logic"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/api/internal/types"
)

// 文件上传
func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.LogWithCustomLevel("requests", r.Host+" ["+r.RequestURI+"]")
		var req types.FileUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}
		l := logic.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req, r)
		response.HttpResult(r, w, resp, err)
	}
}

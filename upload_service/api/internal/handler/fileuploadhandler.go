package handler

import (
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/logic"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"
)

// 文件上传
func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}
		l := logic.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(&req, file, fileHeader)
		response.HttpResult(r, w, resp, err)
	}
}

package Chunk

import (
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/logic/Chunk"
	"net/http"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 完成分片上传
func CompleteMultipartUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChunkUploadCompleteRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := Chunk.NewCompleteMultipartUploadLogic(r.Context(), svcCtx)
		resp, err := l.CompleteMultipartUpload(&req)
		response.HttpResult(r, w, resp, err)
	}
}

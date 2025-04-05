package chunk

import (
	"github.com/Auroraol/cloud-storage/common/logx"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/logic/chunk"
	"net/http"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 初始化分片上传
func InitiateMultipartUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.LogWithCustomLevel("requests", r.Host+" ["+r.RequestURI+"]")
		var req types.ChunkUploadInitRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := chunk.NewInitiateMultipartUploadLogic(r.Context(), svcCtx)
		resp, err := l.InitiateMultipartUpload(&req)
		response.HttpResult(r, w, resp, err)
	}
}

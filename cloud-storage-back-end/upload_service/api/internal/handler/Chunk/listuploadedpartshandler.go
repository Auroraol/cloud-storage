package Chunk

import (
	"github.com/Auroraol/cloud-storage/common/response"
	"net/http"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/logic/Chunk"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 查询分片上传状态
func ListUploadedPartsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListPartsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := Chunk.NewListUploadedPartsLogic(r.Context(), svcCtx)
		resp, err := l.ListUploadedParts(&req)
		response.HttpResult(r, w, resp, err)
	}
}

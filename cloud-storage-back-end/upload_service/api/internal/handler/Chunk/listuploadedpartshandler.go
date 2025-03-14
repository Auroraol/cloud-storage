package Chunk

import (
	"github.com/Auroraol/cloud-storage/common/logx"
	"net/http"

	"github.com/Auroraol/cloud-storage/common/response"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/logic/Chunk"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 查询分片上传状态
func ListUploadedPartsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.LogWithCustomLevel("requests", r.Host+" ["+r.RequestURI+"]")
		var req types.ListPartsRequest
		if err := httpx.Parse(r, &req); err != nil {
			// 1. 参数验证
			if req.UploadId == "" {
				response.ParamErrorResult(r, w, response.NewErrCodeMsg(response.SYSTEM_ERROR, "uploadId不能为空"))
				return
			}
			if req.Key == "" {
				response.ParamErrorResult(r, w, response.NewErrCodeMsg(response.SYSTEM_ERROR, "key不能为空"))
				return
			}
			response.ParamErrorResult(r, w, err)
			return
		}

		l := Chunk.NewListUploadedPartsLogic(r.Context(), svcCtx)
		resp, err := l.ListUploadedParts(&req)
		response.HttpResult(r, w, resp, err)
	}
}

package chunk

import (
	"github.com/Auroraol/cloud-storage/common/logx"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/logic/chunk"

	"github.com/Auroraol/cloud-storage/common/response"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 上传分片
func UploadPartHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.LogWithCustomLevel("requests", r.Host+" ["+r.RequestURI+"]")
		var req types.ChunkUploadRequest
		// 参数验证
		if err := httpx.Parse(r, &req); err != nil {
			if req.UploadId == "" {
				response.ParamErrorResult(r, w, response.NewErrCodeMsg(response.SYSTEM_ERROR, "uploadId不能为空"))
				return
			}
			if req.ChunkIndex <= 0 {
				response.ParamErrorResult(r, w, response.NewErrCodeMsg(response.SYSTEM_ERROR, "分片索引必须大于0"))
				return
			}
			if req.Key == "" {
				response.ParamErrorResult(r, w, response.NewErrCodeMsg(response.SYSTEM_ERROR, "key不能为空"))
				return
			}
			response.ParamErrorResult(r, w, err)
			return
		}

		// 获取上传的文件分片
		file, header, err := r.FormFile("file")
		if err != nil {
			response.ParamErrorResult(r, w, response.NewErrMsg("获取上传文件分片失败"))
			return
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				zap.S().Errorf("关闭文件失败: %v", err)
			}
		}(file)

		l := chunk.NewUploadPartLogic(r.Context(), svcCtx)
		resp, err := l.UploadPart(&req, header)
		response.HttpResult(r, w, resp, err)
	}
}

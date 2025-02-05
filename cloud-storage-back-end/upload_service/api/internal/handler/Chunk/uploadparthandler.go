package Chunk

import (
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/logic/Chunk"
	"github.com/zeromicro/go-zero/core/logx"
	"mime/multipart"
	"net/http"

	"github.com/Auroraol/cloud-storage/common/response"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 上传分片
func UploadPartHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChunkUploadRequest
		// 参数验证
		if err := httpx.Parse(r, &req); err != nil {
			if req.UploadId == "" {
				logx.Errorf("uploadId不能为空: %v", err)
			}
			if req.ChunkIndex <= 0 {
				logx.Errorf("分片索引必须大于0: %v", err)
			}
			if req.Key == "" {
				logx.Errorf("key不能为空: %v", err)
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
				logx.Errorf("关闭文件失败: %v", err)
			}
		}(file)

		l := Chunk.NewUploadPartLogic(r.Context(), svcCtx)
		resp, err := l.UploadPart(&req, header)
		response.HttpResult(r, w, resp, err)
	}
}

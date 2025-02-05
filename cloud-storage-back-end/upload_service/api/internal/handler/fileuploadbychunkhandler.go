package handler

//// 文件分片上传
//func FileUploadByChunkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var req types.FileUploadByChunkRequest
//		if err := httpx.Parse(r, &req); err != nil {
//			response.ParamErrorResult(r, w, err)
//			return
//		}
//		file, fileHeader, err := r.FormFile("file")
//		if err != nil {
//			response.ParamErrorResult(r, w, err)
//			return
//		}
//		l := logic.NewFileUploadByChunkLogic(r.Context(), svcCtx)
//		resp, err := l.FileUploadByChunk(&req, file, fileHeader)
//		response.HttpResult(r, w, resp, err)
//
//	}
//}

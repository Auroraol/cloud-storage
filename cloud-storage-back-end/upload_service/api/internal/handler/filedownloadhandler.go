package handler

import (
	"github.com/Auroraol/cloud-storage/common/logx"
	"net/http"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
)

// 文件下载
func FileDownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.LogWithCustomLevel("requests", r.Host+" ["+r.RequestURI+"]")
		//var req types.FileDownloadRequest
		//if err := httpx.Parse(r, &req); err != nil {
		//	httpx.Error(w, err)
		//	return
		//}

		// 获取用户ID
		//userId := token.GetUidFromCtx(r.Context())
		//if userId == 0 {
		//	httpx.Error(w, response.NewErrCode(response.CREDENTIALS_INVALID))
		//	return
		//}
		//
		//// 查询文件信息
		//repositoryInfo, err := svcCtx.RepositoryPoolModel.FindOne(r.Context(), uint64(req.RepositoryId))
		//if err != nil {
		//	httpx.Error(w, response.NewErrMsg("文件不存在"))
		//	return
		//}
		//
		//// 生成下载URL
		//downloadURL := oss.DownloadURL(repositoryInfo.OssKey)
		//if downloadURL == "" {
		//	httpx.Error(w, response.NewErrMsg("生成下载链接失败"))
		//	return
		//}
		//
		//// 添加操作日志
		//svcCtx.AuditLogServiceRpc.CreateOperationLog(r.Context(), &auditservicerpc.OperationLogReq{
		//	UserId:   userId,
		//	Content:  "下载文件",
		//	FileSize: int32(repositoryInfo.Size),
		//	Flag:     1, // 1表示下载操作
		//})
		//
		//httpx.OkJson(w, &types.FileDownloadResponse{
		//	URL: downloadURL,
		//})
	}
}

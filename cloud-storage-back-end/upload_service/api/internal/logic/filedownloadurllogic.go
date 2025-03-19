package logic

import (
	"context"
	"net/http"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/store/oss"
	"github.com/Auroraol/cloud-storage/common/token"
	"github.com/Auroraol/cloud-storage/log_service/rpc/client/auditservicerpc"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileDownloadUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取文件下载链接
func NewFileDownloadUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDownloadUrlLogic {
	return &FileDownloadUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDownloadUrlLogic) FileDownloadUrl(req *types.FileDownloadUrlRequest, r *http.Request) (resp *types.FileDownloadUrlResponse, err error) {
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("用户(%d)获取文件(%s)的下载链接失败，用户未登录", userId, req.RepositoryId)
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 查询文件信息
	repositoryInfo, err := l.svcCtx.RepositoryPoolModel.FindOneByIdentity(l.ctx, uint64(req.RepositoryId))
	if err != nil {
		zap.S().Errorf("查询文件信息失败 err: %s", err)
		return nil, response.NewErrMsg("查询文件信息失败")
	}

	// 生成下载URL
	downloadURL := oss.DownloadURL(repositoryInfo.OssKey)
	if downloadURL == "" {
		zap.S().Error("生成下载链接失败")
		return nil, response.NewErrMsg("生成下载链接失败")
	}

	// 添加操作日志
	l.svcCtx.AuditLogServiceRpc.CreateOperationLog(l.ctx, &auditservicerpc.OperationLogReq{
		UserId:   userId,
		Content:  "获取文件下载链接",
		FileSize: int32(repositoryInfo.Size),
		Flag:     1, // 1表示下载操作
		FileName: repositoryInfo.Name,
	})

	zap.S().Infof("用户(%d)获取文件(%s)的下载链接", userId, repositoryInfo.Name)

	return &types.FileDownloadUrlResponse{
		URL: downloadURL,
	}, nil
}

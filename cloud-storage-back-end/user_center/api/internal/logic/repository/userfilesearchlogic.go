package repository

import (
	"context"

	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/response"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/token"
	uploadServicePb "github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/rpc/pb"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 搜索用户文件和文件夹
func NewUserFileSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileSearchLogic {
	return &UserFileSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileSearchLogic) UserFileSearch(req *types.UserFileSearchRequest) (resp *types.UserFileSearchResponse, err error) {
	// 处理分页参数
	pageSize := req.Size
	if pageSize == 0 {
		pageSize = 10
	}
	startPage := req.Page
	if startPage == 0 {
		startPage = 1
	}
	startIndex := pageSize * (startPage - 1)

	// 获取用户ID
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 执行搜索
	searchResults, err := l.svcCtx.UserRepositoryModel.SearchFilesByKeywordInPage(
		l.ctx,
		req.ParentId,
		userId,
		req.Keyword,
		startIndex,
		pageSize,
	)
	if err != nil {
		zap.S().Errorf("搜索文件失败: %v", err)
		return nil, response.NewErrMsg("搜索文件失败")
	}

	// 获取总数
	total, err := l.svcCtx.UserRepositoryModel.CountSearchResultsByKeyword(
		l.ctx,
		req.ParentId,
		userId,
		req.Keyword,
	)
	if err != nil {
		zap.S().Errorf("统计搜索结果失败: %v", err)
	}

	// 组装返回结果
	fileList := make([]*types.UserFile, 0)
	for _, item := range searchResults {
		// 文件夹
		if item.RepositoryId == 0 {
			fileList = append(fileList, &types.UserFile{
				Id:           int64(item.Id),
				RepositoryId: 0,
				Name:         item.Name,
				Ext:          "",
				Path:         "",
				Size:         0,
				UpdateTime:   item.UpdateTime.Unix(),
			})
		} else {
			// 文件，需要通过RPC获取文件详情
			repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(
				l.ctx,
				&uploadServicePb.RepositoryReq{RepositoryId: int64(item.RepositoryId)},
			)
			if err != nil {
				zap.S().Errorf("获取文件详情失败: %v", err)
				continue
			}
			fileList = append(fileList, &types.UserFile{
				Id:           int64(item.Id),
				RepositoryId: int64(item.RepositoryId),
				Name:         item.Name,
				Ext:          repositoryInfo.Ext,
				Path:         repositoryInfo.Path,
				Size:         repositoryInfo.Size,
				UpdateTime:   repositoryInfo.UpdateTime,
			})
		}
	}

	return &types.UserFileSearchResponse{
		List:  fileList,
		Count: total,
	}, nil
}

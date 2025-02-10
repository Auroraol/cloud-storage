package repository

import (
	"context"
	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/time"
	"github.com/Auroraol/cloud-storage/common/token"
	uploadServicePb "github.com/Auroraol/cloud-storage/upload_service/rpc/pb"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileAndFolderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户文件和文件夹列表
func NewUserFileAndFolderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileAndFolderListLogic {
	return &UserFileAndFolderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileAndFolderListLogic) UserFileAndFolderList(req *types.UserFileAndFolderListRequest) (resp *types.UserFileAndFolderListResponse, err error) {
	pageSize := req.Size
	if req.Size == 0 {
		pageSize = 10
	}
	startPage := req.Page
	if startPage == 0 {
		startPage = 1
	}
	startIndex := pageSize * (startPage - 1)

	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 先获取文件夹列表
	folders, err := l.svcCtx.UserRepositoryModel.FindAllFolderByParentId(l.ctx, req.Id, userId)
	if err != nil {
		return nil, response.NewErrMsg("获取文件夹列表失败！")
	}

	// 获取文件列表
	files, err := l.svcCtx.UserRepositoryModel.FindAllFileByParentId(l.ctx, req.Id, userId)
	if err != nil {
		return nil, response.NewErrMsg("获取文件列表失败！")
	}

	// 合并文件夹和文件列表
	totalList := make([]*types.UserFile, 0)

	// 添加文件夹
	for _, folder := range folders {
		s := folder.UpdateTime.String()
		timePart := s[:19]
		timestamp, err := time.StringTimeToTimestamp(timePart)

		if err != nil {
			return nil, err
		}

		totalList = append(totalList, &types.UserFile{
			Id:           int64(folder.Id),
			RepositoryId: 0, // 文件夹的 RepositoryId 为 0
			Name:         folder.Name,
			Ext:          "",
			Path:         "",
			Size:         0,
			UpdateTime:   timestamp,
		})
	}

	// 添加文件
	for _, file := range files {
		repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadServicePb.RepositoryReq{RepositoryId: int64(file.RepositoryId)})
		if err != nil {
			continue
		}
		totalList = append(totalList, &types.UserFile{
			Id:           int64(file.Id),
			RepositoryId: int64(file.RepositoryId),
			Name:         file.Name,
			Ext:          repositoryInfo.Ext,
			Path:         repositoryInfo.Path,
			Size:         repositoryInfo.Size,
			UpdateTime:   repositoryInfo.UpdateTime,
		})
	}

	// 计算分页
	totalCount := len(totalList)
	start := startIndex
	end := start + pageSize
	if start >= int64(totalCount) {
		return &types.UserFileAndFolderListResponse{
			List:  []*types.UserFile{},
			Count: int64(totalCount),
		}, nil
	}
	if end > int64(totalCount) {
		end = int64(totalCount)
	}

	return &types.UserFileAndFolderListResponse{
		List:  totalList[start:end],
		Count: int64(totalCount),
	}, nil
}

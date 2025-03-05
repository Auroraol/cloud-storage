package recycle

import (
	"context"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/time"
	"github.com/Auroraol/cloud-storage/common/token"
	uploadServicePb "github.com/Auroraol/cloud-storage/upload_service/rpc/pb"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRecycleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户回收站列表
func NewUserRecycleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRecycleListLogic {
	return &UserRecycleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRecycleListLogic) UserRecycleList(req *types.UserRecycleListRequest) (resp *types.UserRecycleListResponse, err error) {
	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	pageSize := req.Size
	if pageSize == 0 {
		pageSize = 10
	}
	startPage := req.Page
	if startPage == 0 {
		startPage = 1
	}
	startIndex := pageSize * (startPage - 1)

	// 获取已删除状态的文件/文件夹
	deletedFiles, err := l.svcCtx.UserRepositoryModel.FindAllDeletedInPage(l.ctx, req.Id, userId, startIndex, pageSize)
	if err != nil {
		zap.S().Error("获取回收站列表失败 err:%v", err)
		return nil, response.NewErrMsg("获取回收站列表失败")
	}

	total, err := l.svcCtx.UserRepositoryModel.CountTotalDeletedByUserId(l.ctx, userId)
	if err != nil {
		zap.S().Error("获取回收站文件总数失败 err:%v", err)
		return nil, response.NewErrMsg("获取回收站文件总数失败")
	}

	fileList := make([]*types.UserRecycleFile, 0)
	for _, file := range deletedFiles {
		if file.RepositoryId == 0 {
			// 文件夹
			s := file.UpdateTime.String()
			timePart := s[:19]
			timestamp, err := time.StringTimeToTimestamp(timePart)
			if err != nil {
				zap.S().Error("转化时间戳失败 err:%v", err)
				continue
			}
			fileList = append(fileList, &types.UserRecycleFile{
				Id:           int64(file.Id),
				RepositoryId: 0,
				Name:         file.Name,
				Ext:          "",
				Path:         "",
				Size:         0,
				UpdateTime:   timestamp,
			})
		} else {
			// 文件
			repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadServicePb.RepositoryReq{RepositoryId: int64(file.RepositoryId)})
			if err != nil {
				continue
			}
			fileList = append(fileList, &types.UserRecycleFile{
				Id:           int64(file.Id),
				RepositoryId: int64(file.RepositoryId),
				Name:         file.Name,
				Ext:          repositoryInfo.Ext,
				Path:         repositoryInfo.Path,
				Size:         repositoryInfo.Size,
				UpdateTime:   repositoryInfo.UpdateTime,
			})
		}
	}

	return &types.UserRecycleListResponse{
		List:  fileList,
		Total: total,
	}, nil
}

package repository

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
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}

	// 获取所有正常状态的文件/文件夹信息
	allUserRepository, err := l.svcCtx.UserRepositoryModel.FindAllNormalInPage(l.ctx, req.Id, userId, startIndex, pageSize)
	total, err := l.svcCtx.UserRepositoryModel.CountTotalNormalByIdAndParentId(l.ctx, userId, req.Id)
	if err != nil {
		zap.S().Error("获取用户文件列表失败 err:%v", err)
	}

	// 合并文件夹和文件列表
	totalList := make([]*types.UserFile, 0)
	for _, userRepository := range allUserRepository {
		// 文件夹
		if userRepository.RepositoryId == 0 {
			s := userRepository.UpdateTime.String()
			timePart := s[:19]
			timestamp, err := time.StringTimeToTimestamp(timePart)
			if err != nil {
				zap.S().Error("转化时间戳失败 err:%v", err)
				continue
			}
			totalList = append(totalList, &types.UserFile{
				Id:           int64(userRepository.Id),
				RepositoryId: 0, // 文件夹的 RepositoryId 为 0
				Name:         userRepository.Name,
				Ext:          "",
				Path:         "",
				Size:         0,
				UpdateTime:   timestamp,
			})
		} else {
			// 文件
			repositoryInfo, err := l.svcCtx.UploadServiceRpc.GetRepositoryPoolByRepositoryId(l.ctx, &uploadServicePb.RepositoryReq{RepositoryId: int64(userRepository.RepositoryId)})
			if err != nil {
				continue
			}
			totalList = append(totalList, &types.UserFile{
				Id:           int64(userRepository.Id),
				RepositoryId: int64(userRepository.RepositoryId),
				Name:         userRepository.Name,
				Ext:          repositoryInfo.Ext,
				Path:         repositoryInfo.Path,
				Size:         repositoryInfo.Size,
				UpdateTime:   repositoryInfo.UpdateTime,
			})
		}
	}

	return &types.UserFileAndFolderListResponse{
		List:  totalList,
		Count: total,
	}, nil
}

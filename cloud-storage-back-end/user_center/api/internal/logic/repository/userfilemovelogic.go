package repository

import (
	"context"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/common/response"
	"github.com/Auroraol/cloud-storage/common/token"
	"github.com/Auroraol/cloud-storage/log_service/rpc/client/auditservicerpc"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFileMoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户文件移动
func NewUserFileMoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFileMoveLogic {
	return &UserFileMoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFileMoveLogic) UserFileMove(req *types.UserFileMoveRequest) (resp *types.UserFileMoveResponse, err error) {
	//检测该文件是否存在
	userFileInfo, err := l.svcCtx.UserRepositoryModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil {
		zap.S().Error("原文件不存在！")
		return nil, response.NewErrMsg("原文件不存在！")
	}

	userId := token.GetUidFromCtx(l.ctx)
	if userId == 0 {
		zap.S().Error("凭证无效")
		return nil, response.NewErrCode(response.CREDENTIALS_INVALID)
	}
	count, err := l.svcCtx.UserRepositoryModel.CountByIdAndParentId(l.ctx, req.Id, userId, req.ParentId)
	if err != nil {
		zap.S().Error("UserRepositoryModel.CountByIdAndParentId err:%s", err)
		return nil, err
	}
	if count > 0 {
		zap.S().Error("已存在相同名称的文件！")
		return nil, response.NewErrMsg("已存在相同名称的文件！")
	}

	// 添加操作日志文件id
	fileId := userFileInfo.RepositoryId
	if userFileInfo.ParentId != 0 {
		fileId = userFileInfo.ParentId
	}

	//修改
	userFileInfo.ParentId = uint64(req.ParentId)
	err = l.svcCtx.UserRepositoryModel.Update(l.ctx, userFileInfo)
	if err != nil {
		zap.S().Error("UserRepositoryModel.Update err:%v", err)
		return nil, err
	}

	// 添加操作日志
	l.svcCtx.AuditLogServiceRpc.CreateOperationLog(l.ctx, &auditservicerpc.OperationLogReq{
		UserId:   userId,
		Content:  "移动文件",
		Flag:     5,
		FileId:   int64(fileId),
		FileName: userFileInfo.Name,
	})

	return &types.UserFileMoveResponse{}, nil
}

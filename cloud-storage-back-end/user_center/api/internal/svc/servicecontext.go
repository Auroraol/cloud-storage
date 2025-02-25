package svc

import (
	"github.com/Auroraol/cloud-storage/log_service/rpc/client/auditservicerpc"
	"github.com/Auroraol/cloud-storage/upload_service/rpc/uploadservicerpc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/config"
	"github.com/Auroraol/cloud-storage/user_center/model"
	"github.com/Auroraol/cloud-storage/user_center/rpc/client/userservicerpc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	UserModel           model.UserModel
	UserRepositoryModel model.UserRepositoryModel
	UploadServiceRpc    uploadservicerpc.UploadServiceRpc
	UserCenterRpc       userservicerpc.UserServiceRpc
	AuditLogServiceRpc  auditservicerpc.AuditServiceRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlConn := sqlx.NewMysql(c.Options.Dsn)
	return &ServiceContext{
		Config:              c,
		UserModel:           model.NewUserModel(mysqlConn),
		UserRepositoryModel: model.NewUserRepositoryModel(mysqlConn, c.CacheRedis),
		UploadServiceRpc: uploadservicerpc.NewUploadServiceRpc(
			zrpc.MustNewClient(c.UploadServiceRpcConf)),
		UserCenterRpc: userservicerpc.NewUserServiceRpc(
			zrpc.MustNewClient(c.UserCenterRpcConf)),
		AuditLogServiceRpc: auditservicerpc.NewAuditServiceRpc(
			zrpc.MustNewClient(c.AuditServiceRpcConf)),
	}
}

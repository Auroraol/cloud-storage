package svc

import (
	"github.com/Auroraol/cloud-storage/upload_service/rpc/uploadservice"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/config"
	"github.com/Auroraol/cloud-storage/user_center/model"
	"github.com/Auroraol/cloud-storage/user_center/rpc/usercenter"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config              config.Config
	UserModel           model.UserModel
	UserRepositoryModel model.UserRepositoryModel
	UploadServiceRpc    uploadservice.UploadService
	UserCenterRpc       usercenter.UserCenter
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlConn := sqlx.NewMysql(c.Options.Dsn)
	return &ServiceContext{
		Config:              c,
		UserModel:           model.NewUserModel(mysqlConn),
		UserRepositoryModel: model.NewUserRepositoryModel(mysqlConn, c.CacheRedis),
		UploadServiceRpc: uploadservice.NewUploadService(
			zrpc.MustNewClient(c.UploadServiceRpcConf)),
		UserCenterRpc: usercenter.NewUserCenter(
			zrpc.MustNewClient(c.UserCenterRpcConf)),
	}
}

package svc

import (
	"github.com/Auroraol/cloud-storage/common/logx"
	"github.com/Auroraol/cloud-storage/common/sms"
	"github.com/Auroraol/cloud-storage/log_service/rpc/client/auditservicerpc"
	"github.com/Auroraol/cloud-storage/upload_service/rpc/client/uploadservicerpc"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/config"
	"github.com/Auroraol/cloud-storage/user_center/model"
	"github.com/Auroraol/cloud-storage/user_center/rpc/client/userservicerpc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"go.uber.org/zap"
)

type ServiceContext struct {
	Config              config.Config
	UserModel           model.UserModel
	UserRepositoryModel model.UserRepositoryModel
	UploadServiceRpc    uploadservicerpc.UploadServiceRpc
	UserCenterRpc       userservicerpc.UserServiceRpc
	AuditLogServiceRpc  auditservicerpc.AuditServiceRpc
	Sms                 sms.SmsClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlConn := sqlx.NewMysql(c.Options.Dsn)

	// 初始化日志
	c.LogConfig.CustomLevels = map[string]string{
		"requests": "info", // 自定义业务日志级别，对应info级别
		"audit":    "warn", // 自定义审计日志级别，对应warn级别
	}
	if err := logx.InitLogger(c.LogConfig); err != nil {
		zap.S().Errorf("日志初始化失败: %s", err)
		panic(err)
	}

	SmsClient, err := sms.NewClient(c.Sms)
	if err != nil {
		zap.S().Errorf("创建SMS客户端失败: %s", err.Error())
		panic(err)
	}

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
		Sms: SmsClient,
	}
}

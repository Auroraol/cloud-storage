package main

import (
	"flag"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/log_service/rpc/internal/config"
	auditservicerpcServer "github.com/Auroraol/cloud-storage/log_service/rpc/internal/server/auditservicerpc"
	"github.com/Auroraol/cloud-storage/log_service/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/log_service/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "log_service/rpc/etc/logservice.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterAuditServiceRpcServer(grpcServer, auditservicerpcServer.NewAuditServiceRpcServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	zap.S().Infof("Starting server at %s...", c.ListenOn)
	s.Start()
}

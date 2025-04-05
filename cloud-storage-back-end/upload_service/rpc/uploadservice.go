package main

import (
	"flag"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/rpc/internal/config"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/rpc/internal/server/uploadservicerpc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/upload_service/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "upload_service/rpc/etc/uploadservice.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterUploadServiceRpcServer(grpcServer, server.NewUploadServiceRpcServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	zap.S().Infof("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

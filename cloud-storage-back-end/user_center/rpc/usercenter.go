package main

import (
	"flag"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/common/logx"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/internal/config"
	userrepositoryrpcServer "github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/internal/server/userrepositoryrpc"
	userservicerpcServer "github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/internal/server/userservicerpc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/tree/main/cloud-storage-back-end/user_center/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "user_center/rpc/etc/usercenter.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterUserServiceRpcServer(grpcServer, userservicerpcServer.NewUserServiceRpcServer(ctx))
		pb.RegisterUserRepositoryRpcServer(grpcServer, userrepositoryrpcServer.NewUserRepositoryRpcServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	logx.LogWithCustomLevel("Starting rpc server at %s...\n", c.ListenOn)

	s.Start()
}

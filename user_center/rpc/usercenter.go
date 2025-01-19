package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Auroraol/cloud-storage/user_center/rpc/internal/config"
	userServer "github.com/Auroraol/cloud-storage/user_center/rpc/internal/server/user"
	userrepositoryServer "github.com/Auroraol/cloud-storage/user_center/rpc/internal/server/userrepository"
	"github.com/Auroraol/cloud-storage/user_center/rpc/internal/svc"
	"github.com/Auroraol/cloud-storage/user_center/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "user_center/rpc/etc/usercenter.yaml", "the config file")

func main() {
	flag.Parse() // 打印当前工作目录
	fmt.Printf("Current working directory: %s\n", getCurrentDir())
	fmt.Printf("Config file path: %s\n", *configFile)

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterUserServer(grpcServer, userServer.NewUserServer(ctx))
		pb.RegisterUserRepositoryServer(grpcServer, userrepositoryServer.NewUserRepositoryServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	return "Failed to get current directory"
	if err != nil {
	}
	return dir
}

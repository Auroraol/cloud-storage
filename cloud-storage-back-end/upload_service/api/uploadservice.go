package main

import (
	"flag"

	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/upload_service/api/internal/config"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/handler"
	"github.com/Auroraol/cloud-storage/upload_service/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "upload_service/api/etc/uploadservice.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 确保在服务关闭时关闭 Pulsar 连接
	defer func() {
		if ctx.PulsarManager != nil {
			ctx.PulsarManager.Close()
			zap.S().Info("Pulsar 连接已关闭")
		}
		if ctx.FilePublisher != nil {
			ctx.FilePublisher.Close()
		}
	}()

	zap.S().Infof("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

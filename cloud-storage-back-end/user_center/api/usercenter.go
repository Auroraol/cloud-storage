package main

import (
	"flag"
	"go.uber.org/zap"

	"github.com/Auroraol/cloud-storage/user_center/api/internal/config"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/handler"
	"github.com/Auroraol/cloud-storage/user_center/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "user_center/api/etc/usercenter.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	zap.S().Info("Starting server at %s:%d...", c.Host, c.Port)
	server.Start()
}

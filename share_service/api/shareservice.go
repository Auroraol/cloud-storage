package main

import (
	"flag"
	"fmt"

	"github.com/Auroraol/cloud-storage/share_service/api/internal/config"
	"github.com/Auroraol/cloud-storage/share_service/api/internal/handler"
	"github.com/Auroraol/cloud-storage/share_service/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "share_service/api/etc/shareservice.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

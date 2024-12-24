package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/qjpcpu/filelog"
	"path/filepath"

	"cloud-storage/log_service/api/internal/config"
	"cloud-storage/log_service/api/internal/handler"
	"cloud-storage/log_service/api/internal/svc"
	"github.com/qjpcpu/log"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
)

type app struct {
	SERVER_NAME string
}

const (
	SERVER_NAME = "sdk-tb-api"
	PLATFORM    = "tb"
)

type Config struct {
	Log struct {
		Dir   string `toml:"dir"`
		Level string `toml:"level"`
	} `toml:"log"`
}

var configFile = flag.String("f", "log_service/api/etc/logservice.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	gConfig := &Config{}
	gConfig.Log.Level = "debug"
	gConfig.Log.Dir = "log_service/"

	if gConfig.Log.Dir != "" {
		appFile := filepath.Join(gConfig.Log.Dir, SERVER_NAME+".log")
		errFile := filepath.Join(gConfig.Log.Dir, SERVER_NAME+".err.log")

		lcfg := logger.Standard{}
		lcfg.Level = gConfig.Log.Level
		lcfg.File = appFile
		lcfg.ErrFile = errFile
		lcfg.AppName = SERVER_NAME // 该项配置生产环境会注入环境变量
		logger.ResetStandardWithOptions(logger.Options(lcfg))

		log.GetBuilder().
			SetFile(appFile).
			SetErrorLog(errFile).
			SetFormat(log.DebugFormat).
			SetRotate(log.RotateNone).
			SetLevel(gConfig.Log.Level).
			Submit()

		//apiWriter, _ = filelog.NewWriter(filepath.Join(gConfig.Log.Dir, SERVER_NAME+".api.log"))
		_, _ = filelog.NewWriter(filepath.Join(gConfig.Log.Dir, SERVER_NAME+".api.log"))
		//println(apiWriter)
	}
	logger.Info(context.Background(), "config: ", gConfig)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

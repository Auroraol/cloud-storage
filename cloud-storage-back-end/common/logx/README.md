使用zap日志库和lumberjack进行日志切割

可用于gin，kratos框架中

需要安装：

go get -u go.uber.org/zap
go get -u github.com/natefinch/lumberjack

	github.com/stretchr/testify v1.8.3 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1

1. 能够将事件记录到文件中，而不是应用程序控制台;
2. 日志切割-能够根据文件大小来切割日志文件;
3. 支持不同的日志级别。例如INFO，DEBUG，ERROR等;
4. 能够打印基本信息，如调用文件/函数名和行号，日志时间等;



	// 日志初始化
	conf := logx.LogConfig{
		LogLevel:          "debug",    // 输出日志级别 "debug" "info" "warn" "error"
		LogFormat:         "json",     // 输出日志格式  json
		LogPath:           "./",       // 输出日志文件位置
		LogFileName:       "test.log", // 输出日志文件名称
		LogFileMaxSize:    1,          // 输出单个日志文件大小，单位MB
		LogFileMaxBackups: 10,         // 输出最大日志备份个数
		LogMaxAge:         1000,       // 日志保留时间，单位: 天 (day)
		LogCompress:       false,      // 是否压缩日志
		LogStdout:         false,      // 是否输出到控制台
	}
	// 2. 初始化log
	if err := logx.InitLogger(conf); err != nil {
		panic(err)
	}

	zap.S().Debugf("测试 Debugf 用法：%s", "111") // logger Debugf 用法
	zap.S().Errorf("测试 Errorf 用法：%s", "111") // logger Errorf 用法
	zap.S().Warnf("测试 Warnf 用法：%s", "111")   // logger Warnf 用法
	zap.S().Infof("测试 Infof 用法：%s, %d, %v, %f", "111", 1111, errors.New("collector returned no data"), 3333.33)
	// logger With 用法
	logger := zap.S().With("collector", "cpu", "name", "主机")
	logger.Infof("测试 (With + Infof) 用法：%s", "测试")
	zap.S().Errorf("测试 Errorf 用法：%s", "111")


```
// 配置日志
conf := logx.LogConfig{
    LogLevel:      "info",
    LogFormat:     "json",
    LogPath:       "./logs",
    LogFileName:   "app.log",
    LogFileMaxSize: 100,
    CustomLevels: map[string]string{
        "business": "info", // 自定义业务日志级别，对应info级别
        "audit":    "warn", // 自定义审计日志级别，对应warn级别
    },
    SeparateLevel: false, //是否将不同级别的日志分开存储到不同文件
}
logx.InitLogger(conf)

// 方法1：使用LogWithCustomLevel函数
logx.LogWithCustomLevel("business", "这是一条业务日志")
// 输出: {"level":"BUSINESS","time":"2025-03-05 14:54:43","caller":"logx/logger.go:232","msg":"这是一条业务日志"}

// 方法2：创建自定义级别的Logger
businessLogger, _ := logx.NewCustomLevelLogger("business")
businessLogger.Info("这是一条业务日志")
// 输出: {"level":"BUSINESS","time":"2025-03-05 14:54:43","caller":"logx/logger.go:232","msg":"这是一条业务日志"}
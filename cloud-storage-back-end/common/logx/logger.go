package logx

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const DefaultLogPath = "./logs" // 默认输出日志文件路径

type LogConfig struct {
	LogLevel          string            // 日志打印级别 debug  info  warning  error
	LogFormat         string            // 输出日志格式 json
	LogPath           string            // 输出日志文件路径
	LogFileName       string            // 输出日志文件名称
	LogFileMaxSize    int               // 【日志分割】单个日志文件最多存储量 单位(mb)
	LogFileMaxBackups int               // 【日志分割】日志备份文件最多数量
	LogMaxAge         int               // 日志保留时间，单位: 天 (day)
	LogCompress       bool              // 是否压缩日志
	LogStdout         bool              // 是否输出到控制台
	SeparateLevel     bool              // 是否将不同级别的日志分开存储到不同文件
	CustomLevels      map[string]string // 自定义日志级别，key为级别名称，value为对应的zapcore.Level字符串
}

// 自定义日志级别
type CustomLevel struct {
	Name  string
	Level zapcore.Level
}

// 全局变量，存储自定义日志级别
var customLevels = make(map[string]zapcore.Level)

// 全局变量，存储级别到自定义名称的映射
var levelToCustomName = make(map[zapcore.Level]string)

// RegisterCustomLevel 注册自定义日志级别
func RegisterCustomLevel(name string, level zapcore.Level) {
	customLevels[name] = level
	levelToCustomName[level] = name
}

// GetCustomLevel 获取自定义日志级别
func GetCustomLevel(name string) (zapcore.Level, bool) {
	level, ok := customLevels[name]
	return level, ok
}

// GetCustomLevelName 获取自定义日志级别名称
func GetCustomLevelName(level zapcore.Level) string {
	if name, ok := levelToCustomName[level]; ok {
		return name
	}
	return level.String()
}

// 自定义Level编码器
func customLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	levelName := GetCustomLevelName(l)
	enc.AppendString(strings.ToUpper(levelName))
}

// PulsarCore 实现zapcore.Core接口，将日志发送到Pulsar
type PulsarCore struct {
	zapcore.LevelEnabler
	encoder     zapcore.Encoder
	client      pulsar.Client
	producer    pulsar.Producer
	serviceName string
}

// With 实现zapcore.Core接口
func (c *PulsarCore) With(fields []zapcore.Field) zapcore.Core {
	clone := *c
	clone.encoder = c.encoder.Clone()
	for _, field := range fields {
		field.AddTo(clone.encoder)
	}
	return &clone
}

// Check 实现zapcore.Core接口
func (c *PulsarCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

// Write 实现zapcore.Core接口
func (c *PulsarCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	buf, err := c.encoder.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}

	// 添加服务名称标识
	logData := map[string]interface{}{
		"service": c.serviceName,
		"log":     buf.String(),
		"level":   ent.Level.String(),
		"time":    ent.Time.Format(time.RFC3339),
	}

	// 将日志数据序列化为JSON
	jsonData, err := json.Marshal(logData)
	if err != nil {
		return err
	}

	// 异步发送到Pulsar
	c.producer.SendAsync(context.Background(), &pulsar.ProducerMessage{
		Payload: jsonData,
	}, func(id pulsar.MessageID, message *pulsar.ProducerMessage, err error) {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to send log to Pulsar: %v\n", err)
		}
	})

	return nil
}

// Sync 实现zapcore.Core接口
func (c *PulsarCore) Sync() error {
	c.producer.Flush()
	return nil
}

// InitLogger 初始化Logger
func InitLogger(conf LogConfig) (err error) {
	// 解析全局日志级别
	var globalLevel zapcore.Level
	if err := globalLevel.UnmarshalText([]byte(conf.LogLevel)); err != nil {
		return err
	}

	// 注册自定义日志级别
	if conf.CustomLevels != nil {
		for name, levelStr := range conf.CustomLevels {
			var level zapcore.Level
			if err := level.UnmarshalText([]byte(levelStr)); err != nil {
				return fmt.Errorf("invalid custom level %s: %v", name, err)
			}
			RegisterCustomLevel(name, level)
		}
	}

	// 创建所有日志核心
	var cores []zapcore.Core

	// 标准日志级别
	levels := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.DPanicLevel,
		zapcore.PanicLevel,
		zapcore.FatalLevel,
	}

	// 如果不需要分离日志级别，则创建一个统一的WriteSyncer
	if !conf.SeparateLevel {
		// 创建统一的WriteSyncer
		syncer, err := getLogWriter(conf, "")
		if err != nil {
			return err
		}

		// 创建统一的Enabler
		enabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= globalLevel
		})

		// 创建核心
		encoder := getEncoder(conf)
		core := zapcore.NewCore(encoder, syncer, enabler)
		cores = append(cores, core)
	} else {
		// 分离日志级别，为每个级别创建单独的核心
		for _, lvl := range levels {
			// 如果当前级别低于全局日志级别，则跳过
			if lvl < globalLevel {
				continue
			}

			// 创建当前级别的WriteSyncer
			syncer, err := getLogWriter(conf, lvl.String())
			if err != nil {
				return err
			}

			// 创建当前级别的Enabler（仅允许当前级别）
			enabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
				return l == lvl
			})

			// 创建核心
			encoder := getEncoder(conf)
			core := zapcore.NewCore(encoder, syncer, enabler)
			cores = append(cores, core)
		}

		// 为自定义级别创建核心
		for name, level := range customLevels {
			// 如果当前级别低于全局日志级别，则跳过
			if level < globalLevel {
				continue
			}

			// 创建当前级别的WriteSyncer
			syncer, err := getLogWriter(conf, name)
			if err != nil {
				return err
			}

			// 创建当前级别的Enabler
			enabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
				return l == level
			})

			// 创建核心
			encoder := getEncoder(conf)
			core := zapcore.NewCore(encoder, syncer, enabler)
			cores = append(cores, core)
		}
	}

	// 添加控制台核心（如果需要）
	if conf.LogStdout {
		consoleEncoder := getEncoder(conf)
		consoleSyncer := zapcore.AddSync(os.Stdout)
		consoleEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= globalLevel
		})
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleSyncer, consoleEnabler))
	}

	// 合并所有核心
	core := zapcore.NewTee(cores...)

	// 创建Logger实例
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)

	return nil
}

// getEncoder 编码器配置
func getEncoder(conf LogConfig) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	encoderConfig.TimeKey = "time"

	// 使用自定义Level编码器
	encoderConfig.EncodeLevel = customLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	if conf.LogFormat == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// getLogWriter 创建日志Writer
func getLogWriter(conf LogConfig, levelName string) (zapcore.WriteSyncer, error) {
	// 确保日志目录存在
	if conf.LogPath == "" {
		conf.LogPath = DefaultLogPath
	}
	if err := os.MkdirAll(conf.LogPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create log directory failed: %v", err)
	}

	// 生成文件名
	ext := filepath.Ext(conf.LogFileName)
	base := strings.TrimSuffix(conf.LogFileName, ext)
	filename := conf.LogFileName

	// 如果需要分离级别，则添加级别后缀
	if levelName != "" && conf.SeparateLevel {
		filename = fmt.Sprintf("%s-%s%s", base, strings.ToLower(levelName), ext)
	}

	fullPath := filepath.Join(conf.LogPath, filename)

	// 创建Lumberjack日志切割器
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fullPath,
		MaxSize:    conf.LogFileMaxSize,
		MaxBackups: conf.LogFileMaxBackups,
		MaxAge:     conf.LogMaxAge,
		Compress:   conf.LogCompress,
	}

	return zapcore.AddSync(lumberJackLogger), nil
}

// IsExist 判断文件/目录是否存在（保持不变）
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// 以下是自定义日志级别的使用方法

// LogWithCustomLevel 使用自定义级别记录日志
func LogWithCustomLevel(name string, msg string, fields ...zap.Field) {
	if level, ok := GetCustomLevel(name); ok {
		zap.L().Check(level, msg).Write(fields...)
	}
}

// 创建自定义级别的Logger
func NewCustomLevelLogger(name string) (*zap.Logger, error) {
	level, ok := GetCustomLevel(name)
	if !ok {
		return nil, fmt.Errorf("custom level %s not registered", name)
	}

	return zap.L().WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &customLevelCore{
			Core:  core,
			level: level,
		}
	})), nil
}

// 自定义级别的Core
type customLevelCore struct {
	zapcore.Core
	level zapcore.Level
}

// 重写Check方法，使用自定义级别
func (c *customLevelCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(c.level) {
		ent.Level = c.level
		return ce.AddCore(ent, c)
	}
	return ce
}

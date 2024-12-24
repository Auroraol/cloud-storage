package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(opts ...Option) XdLogger {
	options := newOptions(opts...)
	logger := newZapLogger(os.Stdout, options)
	// if len(opts) > 0 {
	// 	if err := initLoggerWithOptions(logger, options); err != nil {
	// 		logger.Error(context.Background(), "NewZapLogger initLoggerWithOptions fail: ", options)
	// 	}
	// }
	return logger
}

func newZapLogger(writer io.Writer, opt Options) *zapLogger {
	opts := []zap.Option{zap.Fields(zap.String(hostKey, os.Getenv(xdRealHost)))}
	if opt.WithStack {
		opts = append(opts, zap.AddStacktrace(zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})))
	}
	if opt.AppName != "" {
		AppName = opt.AppName
	}
	if AppName != "" {
		opts = append(opts, zap.Fields(zap.String(appKey, AppName)))
	}
	log := &zapLogger{
		opt:    opt,
		zapOpt: opts,
		writer: writer,
		level:  parseZapLevel(opt.Level),
		zapLoggerEntity: zapLoggerEntity{
			withAppName: AppName != "",
		},
	}
	log.logger = zap.New(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(writer), zap.LevelEnablerFunc(log.levelEnableFunc())),
		opts...,
	).Named("zap").With()
	return log
}

var (
	jsonEncoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "file",
		MessageKey:    "msg",
		StacktraceKey: "stack",
		FunctionKey:   "func",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeName:    zapcore.FullNameEncoder,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(time.RFC3339))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
)

func parseZapLevel(in string) (out zapcore.Level) {
	switch in {
	case "panic":
		out = zapcore.PanicLevel
	case "fatal":
		out = zapcore.FatalLevel
	case "error":
		out = zapcore.ErrorLevel
	case "warn":
		out = zapcore.WarnLevel
	case "info":
		out = zapcore.InfoLevel
	case "debug":
		out = zapcore.DebugLevel
	case "trace":
		out = zapcore.DebugLevel
	default:
		out = zapcore.InfoLevel
	}
	return
}

func zapLevelToLogrusLevel(in zapcore.Level) (out Level) {
	out, err := ParseLevel(in.String())
	if err != nil {
		out = InfoLevel
	}
	return out
}

func levelEnableFunc(level zapcore.Level) zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return lvl >= level
	}
}
func (l *zapLogger) levelEnableFunc() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return lvl >= l.level
	}
}
func (l *zapLogger) errLoglevelEnableFunc() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	}
}

var _ XdLogger = &zapLogger{}

// zapLogger 是 ZapLogger 的实现
type zapLogger struct {
	zapLoggerEntity

	writer    io.Writer
	errWriter io.Writer

	level  zapcore.Level
	opt    Options
	zapOpt []zap.Option
}

func (l *zapLogger) GetZapLogger() *zap.Logger {
	return l.zapLoggerEntity.logger
}
func (l *zapLogger) SetOutput(out, shadowOut io.Writer) {
	l.writer = out
	l.logger = l.logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewCore(jsonEncoder, zapcore.AddSync(out), zap.LevelEnablerFunc(l.levelEnableFunc()))
	}))
}
func (l *zapLogger) GetOutput() (out, shadowOut io.Writer) {
	return l.writer, nil
}
func (l *zapLogger) SetLevel(level, shadowLevel Level) {
	l.level = parseZapLevel(level.String())
}
func (l *zapLogger) GetLevel() Level {
	lv, _ := logrus.ParseLevel(l.level.String())
	return Level(lv)
}
func (l *zapLogger) SetFormatter(formatter, shadowFormatter Formatter) {} // todo
func (l *zapLogger) SetReportCaller(include, shadowInclude bool)       {} // todo
func (l *zapLogger) AddHook(hook, shadowHook Hook) {
	if !l.zapLoggerEntity.withAppName {
		l.zapLoggerEntity.withAppName = true
		l.zapOpt = append(l.zapOpt, zap.Fields(zap.String(appKey, AppName)))
		l.logger = l.logger.With(zap.String(appKey, AppName))
	}
	switch v := hook.(type) {
	case *TraceHook:
		l.withTrace = true
	case *FileLineHook:
		l.logger = l.logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(2))
	case *MergeHook:
		l.mergeFields = true
	case *XdLfsHook:
		l.errWriter = v.writer
		l.errLogger = zap.New(
			zapcore.NewCore(jsonEncoder, zapcore.AddSync(l.errWriter), zap.LevelEnablerFunc(l.errLoglevelEnableFunc())),
			l.zapOpt...,
		).Named("zap").With()
	}
}
func (l *zapLogger) ResetHooks() {} // todo

var _ XdLoggerEntry = &zapLoggerEntity{}

type zapLoggerEntity struct {
	logger      *zap.Logger
	errLogger   *zap.Logger
	withTrace   bool
	mergeFields bool
	withAppName bool
	custom      map[string]interface{}
}

func (l *zapLoggerEntity) WithField(key string, value interface{}) XdLoggerEntry {
	cp := *l
	if l.mergeFields {
		if cp.custom == nil {
			cp.custom = make(map[string]interface{})
		}
		cp.custom[key] = value
	} else {
		cp.logger = cp.logger.With(zap.Any(key, value))
		if cp.errLogger != nil {
			cp.errLogger = cp.errLogger.With(zap.Any(key, value))
		}
	}
	return &cp
}
func (l *zapLoggerEntity) WithFields(fields Fields) XdLoggerEntry {
	cp := *l
	if cp.mergeFields {
		if cp.custom == nil {
			cp.custom = make(map[string]interface{})
		}
		for k, v := range fields {
			cp.custom[k] = v
		}
	} else {
		fs := []zap.Field{}
		for k, v := range fields {
			fs = append(fs, zap.Any(k, v))
		}
		cp.logger = cp.logger.With(fs...)
		if cp.errLogger != nil {
			cp.errLogger = cp.errLogger.With(fs...)
		}
	}
	return &cp
}
func (l *zapLoggerEntity) WithError(err error) XdLoggerEntry {
	cp := *l
	cp.logger = l.logger.With(zap.Error(err))
	if cp.errLogger != nil {
		cp.errLogger = cp.errLogger.With(zap.Error(err))
	}
	return &cp
}
func (l *zapLoggerEntity) WithTime(t time.Time) XdLoggerEntry {
	cp := *l
	cp.logger = l.logger.With(zap.Time("time", t))
	if cp.errLogger != nil {
		cp.errLogger = cp.errLogger.With(zap.Time("time", t))
	}
	return &cp
}
func (l *zapLoggerEntity) WithObject(obj interface{}) XdLoggerEntry {
	fields := parseFieldsFromObj(obj) // todo: 优化
	return l.WithFields(fields)
}

func (l *zapLoggerEntity) logf(ctx context.Context, level zapcore.Level, format string, args ...interface{}) {
	if format != "" {
		if strings.Contains(format, "%") && len(args) > 0 {
			format = fmt.Sprintf(format, args...)
		}
	} else {
		for i, arg := range args {
			v := zap.Any("arg_"+strconv.Itoa(i+1), arg)
			if v.String != "" {
				format += v.String
			} else if v.Integer != 0 {
				format += strconv.Itoa(int(v.Integer))
			} else {
				format += fmt.Sprint(arg)
			}
		}
	}
	var fields = []zap.Field{zap.Int64(tsKey, int64(time.Now().UnixNano()/1e6))}
	if len(l.custom) > 0 {
		fields = append(fields, zap.Any(customKey, l.custom)) // todo: 优化性能
	}
	if l.withTrace {
		for k, v := range GetTraceKvs(ctx) {
			fields = append(fields, zap.String(k, v))
		}
	}
	switch level {
	case zapcore.DebugLevel:
		l.logger.Debug(format, fields...)
	case zapcore.InfoLevel:
		l.logger.Info(format, fields...)
	case zapcore.WarnLevel:
		l.logger.Warn(format, fields...)
	case zapcore.ErrorLevel:
		l.logger.Error(format, fields...)
		if l.errLogger != nil {
			l.errLogger.Error(format, fields...)
		}
	case zapcore.DPanicLevel:
		l.logger.Panic(format, fields...)
		if l.errLogger != nil {
			l.errLogger.Panic(format, fields...)
		}
	case zapcore.FatalLevel:
		l.logger.Fatal(format, fields...)
		if l.errLogger != nil {
			l.errLogger.Fatal(format, fields...)
		}
	default:
		l.logger.Info(format, fields...)
	}
}
func (l *zapLoggerEntity) Tracef(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, zapcore.InfoLevel, format, args...)
}
func (l *zapLoggerEntity) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, zapcore.DebugLevel, format, args...)
}
func (l *zapLoggerEntity) Infof(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, zapcore.InfoLevel, format, args...)
}
func (l *zapLoggerEntity) Printf(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, zapcore.InfoLevel, format, args...)
}
func (l *zapLoggerEntity) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, zapcore.WarnLevel, format, args...)
}
func (l *zapLoggerEntity) Warningf(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, zapcore.WarnLevel, format, args...)
}
func (l *zapLoggerEntity) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, zapcore.ErrorLevel, format, args...)
}
func (l *zapLoggerEntity) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, zapcore.FatalLevel, format, args...)
}
func (l *zapLoggerEntity) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, zapcore.DPanicLevel, format, args...)
}
func (l *zapLoggerEntity) Logf(ctx context.Context, level Level, format string, args ...interface{}) {
	l.logf(ctx, zapcore.InfoLevel, format, args...)
}
func (l *zapLoggerEntity) Log(ctx context.Context, level Level, args ...interface{}) {
	l.logf(ctx, zapcore.InfoLevel, "", args...)
}
func (l *zapLoggerEntity) Trace(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.InfoLevel, "", args...)
}
func (l *zapLoggerEntity) Debug(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.DebugLevel, "", args...)
}
func (l *zapLoggerEntity) Info(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.InfoLevel, "", args...)
}
func (l *zapLoggerEntity) Print(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.InfoLevel, "", args...)
}
func (l *zapLoggerEntity) Warn(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.WarnLevel, "", args...)
}
func (l *zapLoggerEntity) Warning(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.WarnLevel, "", args...)
}
func (l *zapLoggerEntity) Error(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.ErrorLevel, "", args...)
}
func (l *zapLoggerEntity) Fatal(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.FatalLevel, "", args...)
}
func (l *zapLoggerEntity) Panic(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.DPanicLevel, "", args...)
}
func (l *zapLoggerEntity) Logln(ctx context.Context, level Level, args ...interface{}) {
	l.logf(ctx, zapcore.InfoLevel, "", args...)
}
func (l *zapLoggerEntity) Traceln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.InfoLevel, "", args...)
}
func (l *zapLoggerEntity) Debugln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.DebugLevel, "", args...)
}
func (l *zapLoggerEntity) Infoln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.InfoLevel, "", args...)
}
func (l *zapLoggerEntity) Println(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.InfoLevel, "", args...)
}
func (l *zapLoggerEntity) Warnln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.WarnLevel, "", args...)
}
func (l *zapLoggerEntity) Warningln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.WarnLevel, "", args...)
}
func (l *zapLoggerEntity) Errorln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.ErrorLevel, "", args...)
}
func (l *zapLoggerEntity) Fatalln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.FatalLevel, "", args...)
}
func (l *zapLoggerEntity) Panicln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, zapcore.DPanicLevel, "", args...)
}

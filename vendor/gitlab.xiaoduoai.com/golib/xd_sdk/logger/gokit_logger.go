package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-kit/log"
	loglevel "github.com/go-kit/log/level"
)

func NewGoKitLogger(opts ...Option) XdLogger {
	options := newOptions(opts...)
	logger := newGoKitLogger(os.Stdout, options)
	// if len(opts) > 0 {
	// 	if err := initLoggerWithOptions(logger, options); err != nil {
	// 		logger.Error(context.Background(), "NewGoKitLogger initLoggerWithOptions fail: ", options)
	// 	}
	// }
	return logger
}

func newGoKitLogger(writer io.Writer, opt Options) *goKitLogger {
	level := ParseLevelOrInfo(opt.Level)
	logger := log.With(loglevel.NewFilter(log.NewJSONLogger(writer), getGoKitLevelOpt(level)), "logger", "go-kit")
	return &goKitLogger{
		goKitEntity: goKitEntity{
			host:   os.Getenv(hostKey),
			Logger: logger,
			opt:    &opt,
		},
		writer: writer,
		level:  level,
	}
}

var _ XdLogger = &goKitLogger{}

type goKitLogger struct {
	goKitEntity
	writer io.Writer
	level  Level
}

func (l *goKitLogger) GetOutput() (out, shadowOut io.Writer) {
	return l.writer, nil
}
func (l *goKitLogger) SetOutput(out, shadowOut io.Writer) {
	*l = *newGoKitLogger(out, *l.opt)
}

func getGoKitLevelOpt(level Level) (opt loglevel.Option) {
	switch level {
	case DebugLevel, TraceLevel:
		opt = loglevel.AllowDebug()
	case InfoLevel:
		opt = loglevel.AllowInfo()
	case WarnLevel:
		opt = loglevel.AllowWarn()
	case ErrorLevel, PanicLevel, FatalLevel:
		opt = loglevel.AllowError()
	default:
		opt = loglevel.AllowInfo()
	}
	return
}

func (l *goKitLogger) SetLevel(level, shadowLevel Level) {
	l.Logger = loglevel.NewFilter(l.Logger, getGoKitLevelOpt(level))
}
func (l *goKitLogger) GetLevel() Level {
	return Level(l.level)
}
func (l *goKitLogger) SetFormatter(formatter, shadowFormatter Formatter) {} // todo
func (l *goKitLogger) SetReportCaller(include, shadowInclude bool)       {} // todo
func (l *goKitLogger) AddHook(hook, shadowHook Hook) {
	switch hook.(type) {
	case *TraceHook:
		l.withTrace = true
	case *FileLineHook:
		// l.logger = l.logger.With().CallerWithSkipFrameCount(4).Logger()
		// 暂未实现
	case *MergeHook:
		l.mergeFields = true
		l.Logger = log.With(l.Logger, hostKey, os.Getenv(xdRealHost), appKey, AppName)
		// l.logger = l.logger.With().Str(hostKey, os.Getenv(xdRealHost)).Logger()
	}
}                                  // todo
func (l *goKitLogger) ResetHooks() {} // todo

type goKitEntity struct {
	log.Logger
	host        string
	withTrace   bool
	mergeFields bool
	custom      map[string]interface{}

	opt *Options
}

func (l *goKitEntity) getLevelFunc(lv Level) log.Logger {
	switch lv {
	case DebugLevel, TraceLevel:
		return loglevel.Debug(l.Logger)
	case InfoLevel:
		return loglevel.Info(l.Logger)
	case WarnLevel:
		return loglevel.Warn(l.Logger)
	case ErrorLevel, PanicLevel, FatalLevel:
		return loglevel.Error(l.Logger)
	default:
		return loglevel.Info(l.Logger)
	}
}

func (l *goKitEntity) log(ctx context.Context, level Level, args ...interface{}) {
	msg := ""
	for _, v := range args {
		switch vv := v.(type) {
		case string:
			msg += vv
		default:
			msg += fmt.Sprint(vv)
		}
	}
	kvs := append(make([]interface{}, 0, 6), hostKey, l.host, timeKey, time.Now(), tsKey, time.Now().UnixNano()/1e6, msgKey, msg)
	if l.withTrace {
		for k, v := range GetTraceKvs(ctx) {
			kvs = append(kvs, k, v)
		}
	}
	if len(l.custom) > 0 {
		kvs = append(kvs, customKey, l.custom)
	}
	if err := l.getLevelFunc(level).Log(kvs...); err != nil {
		fmt.Println("goKitEntity.log fail: ", err)
	}
}

func (l *goKitEntity) WithField(key string, value interface{}) XdLoggerEntry {
	cp := *l
	if cp.mergeFields {
		if cp.custom == nil {
			cp.custom = make(map[string]interface{})
		}
		cp.custom[key] = value
	} else {
		cp.Logger = log.With(l.Logger, key, value)
	}
	return &cp
}
func (l *goKitEntity) WithFields(fields Fields) XdLoggerEntry {
	cp := *l
	if cp.mergeFields {
		if cp.custom == nil {
			cp.custom = make(map[string]interface{})
		}
		for k, v := range fields {
			cp.custom[k] = v
		}
	} else {
		fs := make([]interface{}, 0, len(fields))
		for k, v := range fields {
			fs = append(fs, k, v)
		}
		cp.Logger = log.With(l.Logger, fs...)
	}
	return &cp
}
func (l *goKitEntity) WithError(err error) XdLoggerEntry {
	return l.WithField("err", err)
}
func (l *goKitEntity) WithTime(t time.Time) XdLoggerEntry {
	return l.WithField("time", t)
}
func (l *goKitEntity) WithObject(obj interface{}) XdLoggerEntry {
	return l.WithField("obj", obj)
}

func (l *goKitEntity) Tracef(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, DebugLevel, format, fmt.Sprintf(format, args...))
}
func (l *goKitEntity) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, DebugLevel, format, fmt.Sprintf(format, args...))
}
func (l *goKitEntity) Infof(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, InfoLevel, format, fmt.Sprintf(format, args...))
}
func (l *goKitEntity) Printf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, DebugLevel, format, fmt.Sprintf(format, args...))
}
func (l *goKitEntity) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, WarnLevel, format, fmt.Sprintf(format, args...))
}
func (l *goKitEntity) Warningf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, WarnLevel, format, fmt.Sprintf(format, args...))
}
func (l *goKitEntity) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, ErrorLevel, format, fmt.Sprintf(format, args...))
}
func (l *goKitEntity) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, ErrorLevel, format, fmt.Sprintf(format, args...))
	os.Exit(1)
}
func (l *goKitEntity) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, ErrorLevel, format, fmt.Sprintf(format, args...))
	panic(fmt.Sprintf(format, args...))
}

func (l *goKitEntity) Logf(ctx context.Context, lv Level, format string, args ...interface{}) {
	l.log(ctx, DebugLevel, format, fmt.Sprintf(format, args...))
}
func (l *goKitEntity) Log(ctx context.Context, lv Level, args ...interface{}) {
	l.log(ctx, DebugLevel, args...)
}
func (l *goKitEntity) Trace(ctx context.Context, args ...interface{}) {
	l.log(ctx, DebugLevel, args...)
}
func (l *goKitEntity) Debug(ctx context.Context, args ...interface{}) {
	l.log(ctx, DebugLevel, args...)
}
func (l *goKitEntity) Info(ctx context.Context, args ...interface{}) {
	l.log(ctx, InfoLevel, args...)
}
func (l *goKitEntity) Print(ctx context.Context, args ...interface{}) {
	l.log(ctx, DebugLevel, args...)
}
func (l *goKitEntity) Warn(ctx context.Context, args ...interface{}) {
	l.log(ctx, WarnLevel, args...)
}
func (l *goKitEntity) Warning(ctx context.Context, args ...interface{}) {
	l.log(ctx, WarnLevel, args...)
}
func (l *goKitEntity) Error(ctx context.Context, args ...interface{}) {
	l.log(ctx, ErrorLevel, args...)
}
func (l *goKitEntity) Fatal(ctx context.Context, args ...interface{}) {
	l.log(ctx, ErrorLevel, args...)
	os.Exit(1)
}
func (l *goKitEntity) Panic(ctx context.Context, args ...interface{}) {
	l.log(ctx, DebugLevel, args...)
	panic(fmt.Sprint(args...))
}
func (l *goKitEntity) Logln(ctx context.Context, level Level, args ...interface{}) {
	l.log(ctx, level, args...)
}
func (l *goKitEntity) Traceln(ctx context.Context, args ...interface{}) {
	l.log(ctx, DebugLevel, args...)
}
func (l *goKitEntity) Debugln(ctx context.Context, args ...interface{}) {
	l.log(ctx, DebugLevel, args...)
}
func (l *goKitEntity) Infoln(ctx context.Context, args ...interface{}) {
	l.log(ctx, InfoLevel, args...)
}
func (l *goKitEntity) Println(ctx context.Context, args ...interface{}) {
	l.log(ctx, DebugLevel, args...)
}
func (l *goKitEntity) Warnln(ctx context.Context, args ...interface{}) {
	l.log(ctx, WarnLevel, args...)
}
func (l *goKitEntity) Warningln(ctx context.Context, args ...interface{}) {
	l.log(ctx, WarnLevel, args...)
}
func (l *goKitEntity) Errorln(ctx context.Context, args ...interface{}) {
	l.log(ctx, ErrorLevel, args...)
}
func (l *goKitEntity) Fatalln(ctx context.Context, args ...interface{}) {
	l.log(ctx, ErrorLevel, args...)
	os.Exit(1)
}
func (l *goKitEntity) Panicln(ctx context.Context, args ...interface{}) {
	l.log(ctx, ErrorLevel, args...)
	panic(fmt.Sprint(args...))
}

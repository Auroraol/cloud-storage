package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
)

func NewZerologger(opts ...Option) XdLogger {
	options := newOptions(opts...)
	logger := newZerolog(os.Stdout, options)
	// if len(opts) > 0 {
	// 	if err := initLoggerWithOptions(logger, options); err != nil {
	// 		logger.Error(context.Background(), "NewZerolog initLoggerWithOptions fail: ", options)
	// 	}
	// }
	return logger
}

var _ XdLogger = &zeroLogger{}

type zeroLogger struct {
	zerologEntity
	writer    io.Writer
	errWriter io.Writer

	opt   Options
	level Level
}

// type HookFunc func(e *Event, level Level, message string)
func zerologFieldHook() zerolog.HookFunc {
	return func(e *zerolog.Event, level zerolog.Level, message string) {
		e.Str("logger", "zerolog")
		e.Time("time", time.Now())
		e.Int64(tsKey, time.Now().UnixNano()/1e6)
		e.Str(appKey, AppName)
		e.Str(hostKey, os.Getenv(xdRealHost))
	}
}
func newZerolog(writer io.Writer, opt Options) *zeroLogger {
	zerolog.MessageFieldName = "msg"
	zerolog.CallerFieldName = "file"
	level := ParseLevelOrInfo(opt.Level)
	zerolog.SetGlobalLevel(logrusLevelToZeroLevel(level))

	logCtx := zerolog.New(writer).Hook(zerologFieldHook())
	if opt.WithStack {
		// todo: 此处无效，因为 zerolog 的 stack 是从 .Err 中获取的
		logCtx = logCtx.With().Stack().Logger()
	}
	return &zeroLogger{
		zerologEntity: zerologEntity{
			logger: logCtx,
		},
		opt:    opt,
		writer: writer,
		level:  level,
	}
}

func (l *zeroLogger) GetOutput() (out, shadowOut io.Writer) {
	return l.writer, nil
}
func (l *zeroLogger) SetOutput(out, shadowOut io.Writer) {
	l.logger = l.logger.Output(out)
}
func (l *zeroLogger) SetLevel(level, shadowLevel Level) {
	l.logger.Level(logrusLevelToZeroLevel(level))
}
func (l *zeroLogger) GetLevel() Level {
	return Level(l.level)
}
func (l *zeroLogger) SetFormatter(formatter, shadowFormatter Formatter) {} // todo
func (l *zeroLogger) SetReportCaller(include, shadowInclude bool)       {} // todo
func (l *zeroLogger) AddHook(hook, shadowHook Hook) {
	switch v := hook.(type) {
	case *TraceHook:
		l.withTrace = true
	case *FileLineHook:
		l.logger = l.logger.With().CallerWithSkipFrameCount(4).Logger()
		// 暂未实现
	case *MergeHook:
		l.mergeFields = true
	case *XdLfsHook:
		l.errWriter = v.writer
		l.withErrLogger = true
		l.errLogger = zerolog.New(l.errWriter).Level(zerolog.ErrorLevel).Hook(zerologFieldHook())
	}
}

func (l *zeroLogger) ResetHooks() {} // todo

func logrusLevelToZeroLevel(in Level) (out zerolog.Level) {
	switch in {
	case logrus.PanicLevel:
		out = zerolog.PanicLevel
	case logrus.FatalLevel:
		out = zerolog.FatalLevel
	case logrus.ErrorLevel:
		out = zerolog.ErrorLevel
	case logrus.WarnLevel:
		out = zerolog.WarnLevel
	case logrus.InfoLevel:
		out = zerolog.InfoLevel
	case logrus.DebugLevel:
		out = zerolog.DebugLevel
	case logrus.TraceLevel:
		out = zerolog.TraceLevel
	default:
		out = zerolog.InfoLevel
	}
	return
}

type zerologEntity struct {
	logger        zerolog.Logger
	errLogger     zerolog.Logger
	withErrLogger bool
	withTrace     bool
	mergeFields   bool
	fields        map[string]interface{}
	custom        map[string]interface{}
}

// Only map[string]interface{} and []interface{} are accepted. []interface{} must
// alternate string keys and arbitrary values, and extraneous ones are ignored.
func (l *zerologEntity) WithField(key string, value interface{}) XdLoggerEntry {
	cp := *l
	if cp.mergeFields {
		if cp.custom == nil {
			cp.custom = make(map[string]interface{})
		}
		cp.custom[key] = value
	} else {
		cp.logger = withField(cp.logger, key, value)
		if cp.withErrLogger {
			cp.errLogger = withField(cp.errLogger, key, value)
		}
	}
	return &cp
}
func withField(in zerolog.Logger, key string, value interface{}) (out zerolog.Logger) {
	switch v := value.(type) {
	case int64:
		out = in.With().Int64(key, v).Logger()
	case int:
		out = in.With().Int(key, v).Logger()
	case float32, float64:
		out = in.With().Float64(key, v.(float64)).Logger()
	case string:
		out = in.With().Str(key, v).Logger()
	case []string:
		out = in.With().Strs(key, v).Logger()
	case fmt.Stringer:
		out = in.With().Stringer(key, v).Logger()
	case zerolog.LogObjectMarshaler:
		out = in.With().Object(key, v).Logger()
	case error:
		out = in.With().AnErr(key, v).Logger()
	case time.Time:
		out = in.With().Time(key, v).Logger()
	default:
		out = in.With().Interface(key, v).Logger()
	}
	return out
}
func (l *zerologEntity) WithFields(fields Fields) XdLoggerEntry {
	cp := *l
	if cp.mergeFields {
		if cp.custom == nil {
			cp.custom = make(map[string]interface{})
		}
		for k, v := range fields {
			cp.custom[k] = v
		}
	} else {
		cp.logger = l.logger.With().Fields(map[string]interface{}(fields)).Logger()
		if cp.withErrLogger {
			cp.errLogger = l.errLogger.With().Fields(map[string]interface{}(fields)).Logger()
		}
	}
	return &cp
}
func (l *zerologEntity) WithError(err error) XdLoggerEntry {
	cp := *l
	cp.logger = cp.logger.With().Err(err).Logger()
	if cp.withErrLogger {
		cp.errLogger = cp.errLogger.With().Err(err).Logger()
	}
	return &cp
}
func (l *zerologEntity) WithTime(t time.Time) XdLoggerEntry {
	cp := *l
	cp.logger = cp.logger.With().Time("time", t).Logger()
	if cp.withErrLogger {
		cp.errLogger = cp.errLogger.With().Time("time", t).Logger()
	}
	return &cp
}
func (l *zerologEntity) WithObject(obj interface{}) XdLoggerEntry {
	cp := *l
	cp.logger = cp.logger.With().Fields(obj).Logger()
	if cp.withErrLogger {
		cp.errLogger = cp.errLogger.With().Fields(obj).Logger()
	}
	return &cp
}

func (l *zerologEntity) getLevelFunc(lv Level) (app, err *zerolog.Event) {
	switch lv {
	case TraceLevel:
		app = l.logger.Trace()
	case DebugLevel:
		app = l.logger.Debug()
	case InfoLevel:
		app = l.logger.Info()
	case WarnLevel:
		app = l.logger.Warn()
	case ErrorLevel:
		app = l.logger.Error()
		if l.withErrLogger {
			err = l.errLogger.Error()
		}
	case PanicLevel:
		app = l.logger.Panic()
		if l.withErrLogger {
			err = l.errLogger.Panic()
		}
	case FatalLevel:
		app = l.logger.Fatal()
		if l.withErrLogger {
			err = l.errLogger.Fatal()
		}
	default:
		app = l.logger.Log()
	}
	return app, err
}

func (l *zerologEntity) logf(ctx context.Context, lv Level, format string, args ...interface{}) {
	fields := make(map[string]interface{})
	app, err := l.getLevelFunc(lv)
	if l.withTrace {
		for k, v := range GetTraceKvs(ctx) {
			fields[k] = v
		}
	}
	if len(l.custom) > 0 {
		fields[customKey] = l.custom
	}

	app = app.Fields(fields)
	if format != "" {
		app.Msgf(format, args...)
		if err != nil {
			err.Msgf(format, args...)
		}
	} else {
		msg := ""
		for _, v := range args {
			msg += typeToStr(v)
		}
		app.Msg(msg)
		if err != nil {
			err.Msg(msg)
		}
	}
}

func (l *zerologEntity) Tracef(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, TraceLevel, format, args...)
}
func (l *zerologEntity) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, DebugLevel, format, args...)
}
func (l *zerologEntity) Infof(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, InfoLevel, format, args...)
}
func (l *zerologEntity) Printf(ctx context.Context, format string, args ...interface{}) {
}
func (l *zerologEntity) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, WarnLevel, format, args...)
}
func (l *zerologEntity) Warningf(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, WarnLevel, format, args...)
}
func (l *zerologEntity) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, ErrorLevel, format, args...)
}
func (l *zerologEntity) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, FatalLevel, format, args...)
}
func (l *zerologEntity) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.logf(ctx, PanicLevel, format, args...)
}
func (l *zerologEntity) Logf(ctx context.Context, level Level, format string, args ...interface{}) {
	l.logf(ctx, level, format, args...)
}

var emptyFormat string

func (l *zerologEntity) Log(ctx context.Context, level Level, args ...interface{}) {
	l.logf(ctx, level, emptyFormat, args...)
}
func (l *zerologEntity) Trace(ctx context.Context, args ...interface{}) {
	l.logf(ctx, TraceLevel, emptyFormat, args...)
}
func (l *zerologEntity) Debug(ctx context.Context, args ...interface{}) {
	l.logf(ctx, DebugLevel, emptyFormat, args...)
}
func (l *zerologEntity) Info(ctx context.Context, args ...interface{}) {
	l.logf(ctx, InfoLevel, emptyFormat, args...)
}
func (l *zerologEntity) Print(ctx context.Context, args ...interface{}) {
	l.logf(ctx, DebugLevel, emptyFormat, args...)
}
func (l *zerologEntity) Warn(ctx context.Context, args ...interface{}) {
	l.logf(ctx, WarnLevel, emptyFormat, args...)
}
func (l *zerologEntity) Warning(ctx context.Context, args ...interface{}) {
	l.logf(ctx, WarnLevel, emptyFormat, args...)
}
func (l *zerologEntity) Error(ctx context.Context, args ...interface{}) {
	l.logf(ctx, ErrorLevel, emptyFormat, args...)
}
func (l *zerologEntity) Fatal(ctx context.Context, args ...interface{}) {
	l.logf(ctx, FatalLevel, emptyFormat, args...)
}
func (l *zerologEntity) Panic(ctx context.Context, args ...interface{}) {
	l.logf(ctx, PanicLevel, emptyFormat, args...)
}
func (l *zerologEntity) Logln(ctx context.Context, level Level, args ...interface{}) {
	l.logf(ctx, level, emptyFormat, args...)
}
func (l *zerologEntity) Traceln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, TraceLevel, emptyFormat, args...)
}
func (l *zerologEntity) Debugln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, DebugLevel, emptyFormat, args...)
}
func (l *zerologEntity) Infoln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, InfoLevel, emptyFormat, args...)
}
func (l *zerologEntity) Println(ctx context.Context, args ...interface{}) {
	l.logf(ctx, DebugLevel, emptyFormat, args...)
}
func (l *zerologEntity) Warnln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, WarnLevel, emptyFormat, args...)
}
func (l *zerologEntity) Warningln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, WarnLevel, emptyFormat, args...)
}
func (l *zerologEntity) Errorln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, ErrorLevel, emptyFormat, args...)
}
func (l *zerologEntity) Fatalln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, FatalLevel, emptyFormat, args...)
}
func (l *zerologEntity) Panicln(ctx context.Context, args ...interface{}) {
	l.logf(ctx, PanicLevel, emptyFormat, args...)
}

func typeToStr(in interface{}) (out string) {
	// todo
	switch v := in.(type) {
	case string:
		out = v
	case int, int16, int32, int64, int8:
		out = strconv.FormatInt(v.(int64), 10)
	case uint, uint8, uint16, uint32, uint64:
		out = strconv.FormatUint(v.(uint64), 10)
	case error:
		out = v.Error()
	default:
		out = fmt.Sprint(in)
	}
	return out
}

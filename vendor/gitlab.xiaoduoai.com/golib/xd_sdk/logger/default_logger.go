package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Entry = logrus.Entry
type Ext1FieldLogger = logrus.Ext1FieldLogger
type FieldLogger = logrus.FieldLogger
type Fields = logrus.Fields
type Formatter = logrus.Formatter
type Hook = logrus.Hook
type Logger = XdLogger
type Level = logrus.Level
type LevelHooks = logrus.LevelHooks
type MutexWrap = logrus.MutexWrap

const PanicLevel = logrus.PanicLevel
const FatalLevel = logrus.FatalLevel
const ErrorLevel = logrus.ErrorLevel
const WarnLevel = logrus.WarnLevel
const InfoLevel = logrus.InfoLevel
const DebugLevel = logrus.DebugLevel
const TraceLevel = logrus.TraceLevel

var AllLevels = logrus.AllLevels
var xdStdLogger = xdStdLoggerNew()

func newJSONFormatter() logrus.Formatter {
	formatter := new(JSONFormatter)
	formatter.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	return formatter
}

func newTextFormatter() logrus.Formatter {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	return formatter
}

// New 生成带有指定格式的标准logger
func xdStdLoggerNew() Logger {
	formatter := newJSONFormatter()

	nl := logrus.Logger{
		Out:          os.Stderr,
		Formatter:    formatter,
		Hooks:        make(LevelHooks),
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	sl := logrus.Logger{
		Out:          os.Stderr,
		Formatter:    formatter,
		Hooks:        make(LevelHooks),
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	return &CtxLogger{&nl, &sl}
}

func StandardLogger() Logger {
	return xdStdLogger
}

func SetOutput(out, shadowOut io.Writer) {
	xdStdLogger.SetOutput(out, shadowOut)
}

func SetAppName(app string) {
	AppName = app
}

func GetOutput() (out, shadowOut io.Writer) {
	return xdStdLogger.GetOutput()
}

func SetFormatter(formatter, shadowFormatter Formatter) {
	xdStdLogger.SetFormatter(formatter, shadowFormatter)
}

func SetReportCaller(include, shadowInclude bool) {
	xdStdLogger.SetReportCaller(include, shadowInclude)
}

func SetLevel(level Level) {
	xdStdLogger.SetLevel(level, level)
}

func GetLevel() Level {
	return xdStdLogger.GetLevel()
}

func SetLevelWithShadow(level, shadowLevel Level) {
	xdStdLogger.SetLevel(level, shadowLevel)
}

func AddHook(hook, shadowHook Hook) {
	xdStdLogger.AddHook(hook, shadowHook)
}

func ParseLevel(level string) (Level, error) {
	return logrus.ParseLevel(level)
}

func ParseLevelOrInfo(level string) Level {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return InfoLevel
	}
	return l
}

func NewLogrusEntry(l *logrus.Logger) *Entry {
	return logrus.NewEntry(l)
}

func WithError(err error) XdLoggerEntry {
	return xdStdLogger.WithError(err)
}

func WithField(key string, value interface{}) XdLoggerEntry {
	return xdStdLogger.WithField(key, value)
}

func WithFields(fields Fields) XdLoggerEntry {
	return xdStdLogger.WithFields(fields)
}

func WithTime(t time.Time) XdLoggerEntry {
	return xdStdLogger.WithTime(t)
}

func WithObject(obj interface{}) XdLoggerEntry {
	return xdStdLogger.WithObject(obj)
}

func Trace(ctx context.Context, args ...interface{}) {
	xdStdLogger.Trace(ctx, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	xdStdLogger.Debug(ctx, args...)
}

func Print(ctx context.Context, args ...interface{}) {
	xdStdLogger.Print(ctx, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	xdStdLogger.Info(ctx, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	xdStdLogger.Warn(ctx, args...)
}

func Warning(ctx context.Context, args ...interface{}) {
	xdStdLogger.Warning(ctx, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	xdStdLogger.Error(ctx, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	xdStdLogger.Panic(ctx, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	xdStdLogger.Fatal(ctx, args...)
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	xdStdLogger.Tracef(ctx, format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	xdStdLogger.Debugf(ctx, format, args...)
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	xdStdLogger.Printf(ctx, format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	xdStdLogger.Infof(ctx, format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	xdStdLogger.Warnf(ctx, format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	xdStdLogger.Warningf(ctx, format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	xdStdLogger.Errorf(ctx, format, args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	xdStdLogger.Panicf(ctx, format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	xdStdLogger.Fatalf(ctx, format, args...)
}

func Traceln(ctx context.Context, args ...interface{}) {
	xdStdLogger.Traceln(ctx, args...)
}

func Debugln(ctx context.Context, args ...interface{}) {
	xdStdLogger.Debugln(ctx, args...)
}

func Println(ctx context.Context, args ...interface{}) {
	xdStdLogger.Println(ctx, args...)
}

func Infoln(ctx context.Context, args ...interface{}) {
	xdStdLogger.Infoln(ctx, args...)
}

func Warnln(ctx context.Context, args ...interface{}) {
	xdStdLogger.Warnln(ctx, args...)
}

func Warningln(ctx context.Context, args ...interface{}) {
	xdStdLogger.Warningln(ctx, args...)
}

func Errorln(ctx context.Context, args ...interface{}) {
	xdStdLogger.Errorln(ctx, args...)
}

func Panicln(ctx context.Context, args ...interface{}) {
	xdStdLogger.Panicln(ctx, args...)
}

func Fatalln(ctx context.Context, args ...interface{}) {
	xdStdLogger.Fatalln(ctx, args...)
}

package logger

import (
	"context"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	"gitlab.xiaoduoai.com/golib/xd_sdk/metadata"

	"github.com/sirupsen/logrus"
)

type MergeFlagKey struct{}
type MergeInfoKey struct{}

var (
	mergeFlagKey MergeFlagKey
	mergeInfoKey MergeInfoKey
)

const (
	mergeLineLimit = 100
	mergeLenLimit  = 64 * 1024 //最大长度64kb
)

func WithMergeFlag(ctx context.Context) context.Context {
	return context.WithValue(ctx, mergeFlagKey, 1)
}

func UnsetMergeFlag(ctx context.Context) context.Context {
	return context.WithValue(ctx, mergeFlagKey, nil)
}

func IsMerge(ctx context.Context) bool {
	return ctx.Value(mergeFlagKey) != nil
}

func InitMergeInfo(ctx context.Context) context.Context {
	m := &MergeInfo{
		Mutex: sync.Mutex{},
		Logs:  make([]string, 0, mergeLineLimit),
	}
	return context.WithValue(ctx, mergeInfoKey, m)
}

func MergeInfoFromCtx(ctx context.Context) *MergeInfo {
	iface := ctx.Value(mergeInfoKey)
	info, ok := iface.(*MergeInfo)
	if !ok {
		return nil
	}
	return info
}

type MergeInfo struct {
	Mutex   sync.Mutex
	LogsLen int
	Logs    []string
}

func (m *MergeInfo) FormatLog() string {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	if len(m.Logs) > 0 {
		res := "[" + strings.Join(m.Logs, ",") + "]"
		m.Logs, m.LogsLen = m.Logs[:0], 0
		return res
	}
	return ""
}

func (m *MergeInfo) AppendLog(l string) string {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.Logs = append(m.Logs, l)
	m.LogsLen += len(l)
	if len(m.Logs) >= mergeLineLimit || m.LogsLen > mergeLenLimit {
		res := "[" + strings.Join(m.Logs, ",") + "]"
		m.Logs, m.LogsLen = m.Logs[:0], 0
		return res
	}
	return ""
}

var wipeKeys = map[string]int8{
	appKey:   1, //app
	TraceKey: 1, //trace
	hostKey:  1, //host
}

// 实现接口 XdLoggerEntry
type XdLogShadowEntry struct {
	*logrus.Entry
	nl *logrus.Logger
	sl *logrus.Logger
}

// logger 以及 entry的公用接口，方法一致，业务使用体验一致
type XdLoggerEntry interface {
	WithField(key string, value interface{}) XdLoggerEntry
	WithFields(fields Fields) XdLoggerEntry
	WithError(err error) XdLoggerEntry
	WithTime(t time.Time) XdLoggerEntry
	WithObject(obj interface{}) XdLoggerEntry
	Tracef(ctx context.Context, format string, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Printf(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Warningf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
	Panicf(ctx context.Context, format string, args ...interface{})
	Logf(ctx context.Context, level Level, format string, args ...interface{})
	Log(ctx context.Context, level Level, args ...interface{})
	Trace(ctx context.Context, args ...interface{})
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Print(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warning(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})
	Panic(ctx context.Context, args ...interface{})
	Logln(ctx context.Context, level Level, args ...interface{})
	Traceln(ctx context.Context, args ...interface{})
	Debugln(ctx context.Context, args ...interface{})
	Infoln(ctx context.Context, args ...interface{})
	Println(ctx context.Context, args ...interface{})
	Warnln(ctx context.Context, args ...interface{})
	Warningln(ctx context.Context, args ...interface{})
	Errorln(ctx context.Context, args ...interface{})
	Fatalln(ctx context.Context, args ...interface{})
	Panicln(ctx context.Context, args ...interface{})
}

func (en XdLogShadowEntry) WithField(key string, value interface{}) XdLoggerEntry {
	return &XdLogShadowEntry{en.Entry.WithField(key, value), en.nl, en.sl}
}

func (en XdLogShadowEntry) WithFields(fields Fields) XdLoggerEntry {
	return &XdLogShadowEntry{en.Entry.WithFields(fields), en.nl, en.sl}
}

func (en XdLogShadowEntry) WithError(err error) XdLoggerEntry {
	return &XdLogShadowEntry{en.Entry.WithError(err), en.nl, en.sl}
}

func (en XdLogShadowEntry) WithTime(t time.Time) XdLoggerEntry {
	return &XdLogShadowEntry{en.Entry.WithTime(t), en.nl, en.sl}
}

func (en XdLogShadowEntry) WithObject(obj interface{}) XdLoggerEntry {
	fields := parseFieldsFromObj(obj)
	return &XdLogShadowEntry{en.Entry.WithFields(fields), en.nl, en.sl}
}

// setlogger，根据ctx压测标记，切换entry绑定的logger,将日志输出到不同的地方
func (en XdLogShadowEntry) setLogger(ctx context.Context) {
	if metadata.IsTestFlow(ctx) {
		en.Entry.Logger = en.sl
		return
	} else {
		en.Entry.Logger = en.nl
	}
}

func logLevelTrans(ctx context.Context, originLevel Level) Level {
	if originLevel > logrus.InfoLevel && metadata.IsPrintAtLeastInfo(ctx) {
		return logrus.InfoLevel
	}

	return originLevel
}

func shadowLogf(entry *logrus.Entry, level logrus.Level, format string, args ...interface{}) {
	// Fatal和Panic不进行干扰
	if IsMerge(entry.Context) && entry.Logger.IsLevelEnabled(level) && level > FatalLevel {
		entry.Level = level
		entry.Message = fmt.Sprintf(format, args...)
		if appendLogToCtx(entry) {
			return
		}
	}
	entry.Logf(level, format, args...)
}

func shadowLog(entry *logrus.Entry, level logrus.Level, args ...interface{}) {
	// Fatal和Panic不进行干扰
	if IsMerge(entry.Context) && entry.Logger.IsLevelEnabled(level) && level > FatalLevel {
		entry.Level = level
		entry.Message = fmt.Sprint(args...)
		if appendLogToCtx(entry) {
			return
		}
	}
	entry.Log(level, args...)
}

func shadowLogln(entry *logrus.Entry, level logrus.Level, args ...interface{}) {
	// Fatal和Panic不进行干扰
	if IsMerge(entry.Context) && entry.Logger.IsLevelEnabled(level) && level > FatalLevel {
		entry.Level = level
		entry.Message = fmt.Sprintln(args...)
		if appendLogToCtx(entry) {
			return
		}
	}
	entry.Logln(level, args...)
}

func appendLogToCtx(entry *logrus.Entry) bool {
	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}

	entry.Logger.Hooks.Fire(entry.Level, entry)
	// 原hook调用较深, 所以是skip=6, 此处提高层级为2
	file, line, fn := searchFileLineWithSkip(2)
	entry.Data[FileKey] = fmt.Sprintf("%v:%v", path.Base(file), line)
	idx := strings.LastIndex(fn, "/")
	entry.Data[FuncKey] = fn[idx+1:]
	//host, trace, app信息冗余, 不需要打印
	for k := range wipeKeys {
		delete(entry.Data, k)
	}

	mergeInfo := MergeInfoFromCtx(entry.Context)
	if mergeInfo == nil {
		return false
	}
	log, err := entry.Logger.Formatter.Format(entry)
	if err != nil || len(log) < 2 {
		return false
	}
	// 去除末尾的\n
	log = log[:len(log)-1]

	if res := mergeInfo.AppendLog(string(log)); res != "" {
		entry.Log(entry.Level, res)
	}
	return true
}

func (en XdLogShadowEntry) Tracef(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	shadowLogf(en.Entry.WithContext(ctx), logLevelTrans(ctx, TraceLevel), format, args...)
}

func (en XdLogShadowEntry) Debugf(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	shadowLogf(en.Entry.WithContext(ctx), logLevelTrans(ctx, DebugLevel), format, args...)
}

func (en XdLogShadowEntry) Infof(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	shadowLogf(en.Entry.WithContext(ctx), InfoLevel, format, args...)
}

func (en XdLogShadowEntry) Printf(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Printf(format, args...)
}

func (en XdLogShadowEntry) Warnf(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	shadowLogf(en.Entry.WithContext(ctx), WarnLevel, format, args...)
}

func (en XdLogShadowEntry) Warningf(ctx context.Context, format string, args ...interface{}) {
	en.Warnf(ctx, format, args...)
}

func (en XdLogShadowEntry) Errorf(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	shadowLogf(en.Entry.WithContext(ctx), ErrorLevel, format, args...)
}

func (en XdLogShadowEntry) Fatalf(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Fatalf(format, args...)
}

func (en XdLogShadowEntry) Panicf(ctx context.Context, format string, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logf(PanicLevel, format, args...)
}

func (en XdLogShadowEntry) Logf(ctx context.Context, level Level, format string, args ...interface{}) {
	en.setLogger(ctx)
	shadowLogf(en.Entry.WithContext(ctx), logLevelTrans(ctx, level), format, args...)
}

func (en XdLogShadowEntry) Log(ctx context.Context, level Level, args ...interface{}) {
	en.setLogger(ctx)
	shadowLog(en.Entry.WithContext(ctx), logLevelTrans(ctx, level), args...)
}

func (en XdLogShadowEntry) Trace(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	shadowLog(en.Entry.WithContext(ctx), logLevelTrans(ctx, TraceLevel), args...)
}

func (en XdLogShadowEntry) Debug(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	shadowLog(en.Entry.WithContext(ctx), logLevelTrans(ctx, DebugLevel), args...)
}

func (en XdLogShadowEntry) Info(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	shadowLog(en.Entry.WithContext(ctx), InfoLevel, args...)
}

func (en XdLogShadowEntry) Print(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Print(args...)
}

func (en XdLogShadowEntry) Warn(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	shadowLog(en.Entry.WithContext(ctx), WarnLevel, args...)
}

func (en XdLogShadowEntry) Warning(ctx context.Context, args ...interface{}) {
	en.Warn(ctx, args...)
}

func (en XdLogShadowEntry) Error(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	shadowLog(en.Entry.WithContext(ctx), ErrorLevel, args...)
}

func (en XdLogShadowEntry) Fatal(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Fatal(args...)
}

func (en XdLogShadowEntry) Panic(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Panic(args...)
}

func (en XdLogShadowEntry) Logln(ctx context.Context, level Level, args ...interface{}) {
	en.setLogger(ctx)
	shadowLogln(en.Entry.WithContext(ctx), logLevelTrans(ctx, level), args...)
}

func (en XdLogShadowEntry) Traceln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	shadowLogln(en.Entry.WithContext(ctx), logLevelTrans(ctx, TraceLevel), args...)
}

func (en XdLogShadowEntry) Debugln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	shadowLogln(en.Entry.WithContext(ctx), logLevelTrans(ctx, DebugLevel), args...)
}

func (en XdLogShadowEntry) Infoln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	shadowLogln(en.Entry.WithContext(ctx), InfoLevel, args...)
}

func (en XdLogShadowEntry) Println(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Println(args...)
}

func (en XdLogShadowEntry) Warnln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	shadowLogln(en.Entry.WithContext(ctx), WarnLevel, args...)
}

func (en XdLogShadowEntry) Warningln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	shadowLogln(en.Entry.WithContext(ctx), WarnLevel, args...)
}

func (en XdLogShadowEntry) Errorln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	shadowLogln(en.Entry.WithContext(ctx), ErrorLevel, args...)
}

func (en XdLogShadowEntry) Fatalln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Fatalln(args...)
}

func (en XdLogShadowEntry) Panicln(ctx context.Context, args ...interface{}) {
	en.setLogger(ctx)
	en.Entry.WithContext(ctx).Logln(PanicLevel, args...)
}

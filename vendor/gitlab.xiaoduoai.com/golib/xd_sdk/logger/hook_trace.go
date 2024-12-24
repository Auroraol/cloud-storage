package logger

import (
	"context"

	"gitlab.xiaoduoai.com/golib/xd_sdk/octrace"
	"gitlab.xiaoduoai.com/golib/xd_sdk/public"
)

// 从context中提取出与trace相关的一些字段写入日志中。
const (
	TraceKey      = "trace"
	SpanKey       = "span"
	SampledKey    = "sampled"
	ParentSpanKey = "pspan"
	TimeStamp     = "timestamp"
	customKey     = "custom"
	xdRealHost    = "XD_REAL_HOST"

	CallerHostName = "xd_caller_host"
	CallerSvcName  = "xd_caller_svc"
)

type TraceHook struct{}

func NewTraceHook() *TraceHook {
	return &TraceHook{}
}

func (h *TraceHook) Fire(entry *Entry) error {
	if ctx := entry.Context; ctx != nil {
		span := octrace.ExtractSpan(ctx)
		if span != nil {
			entry.Data[TraceKey] = span.SpanContext().TraceID.String()
		}

		info := public.GetCallInfoFromCtx(entry.Context)
		if info.SvcName != "" {
			entry.Data[CallerSvcName] = info.SvcName
		}
		if info.HostName != "" {
			entry.Data[CallerHostName] = info.HostName

		}
	}
	return nil
}

func (h *TraceHook) Levels() []Level {
	return AllLevels
}

func GetTraceKvs(ctx context.Context) (kvs map[string]string) {
	if ctx == nil {
		return
	}
	kvs = make(map[string]string)
	if span := octrace.ExtractSpan(ctx); span != nil {
		kvs[TraceKey] = span.SpanContext().TraceID.String()
	}

	info := public.GetCallInfoFromCtx(ctx)
	if info.SvcName != "" {
		kvs[CallerSvcName] = info.SvcName
	}
	if info.HostName != "" {
		kvs[CallerHostName] = info.HostName
	}
	return kvs
}

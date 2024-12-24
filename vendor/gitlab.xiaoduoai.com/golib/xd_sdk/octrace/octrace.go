package octrace

import (
	"context"
	"fmt"
	"os"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/sirupsen/logrus"
	"gitlab.xiaoduoai.com/golib/xd_sdk/metadata"
	"gitlab.xiaoduoai.com/golib/xd_sdk/octrace/config"
	"gitlab.xiaoduoai.com/golib/xd_sdk/octrace/filter"
	"gitlab.xiaoduoai.com/golib/xd_sdk/xd_error"
	"go.opencensus.io/trace"
)

var enable = true

func Enabled() bool {
	return enable
}

func SetFilter(f filter.TraceFilter) {
	filter.SetFilter(f)
}

func InitWithConfig(cfg config.Config, onErr func(error)) (func(), error) {
	var opts []Option
	addOpt := func(opt Option) {
		opts = append(opts, opt)
	}
	addOpt(WithSample(SampleType(cfg.SampleType), cfg.Fraction))
	addOpt(WithServiceName(cfg.ServiceName))
	addOpt(WithEndPoint(cfg.AgentEndPoint, cfg.CollectorEndPoint))
	addOpt(WithAuthor(cfg.Username, cfg.Password))
	addOpt(WithTraceLogFile(cfg.TraceLogFile))
	if onErr == nil {
		onErr = func(err error) {
			logrus.StandardLogger().Errorf("jaeger exporter err: %v\n", err)
		}
	}
	addOpt(WithOnErr(onErr))

	return Init(opts...)
}

func Init(opts ...Option) (func(), error) {
	hostname, _ := os.Hostname()
	options := Options{
		agentEndPoint: "localhost:6831",
		sampler:       trace.NeverSample(),
		serviceName:   hostname,
		onError: func(err error) {
			logrus.StandardLogger().Errorf("jaeger exporter err: %v\n", err)
		},
	}
	for _, opt := range opts {
		opt(&options)
	}
	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     options.agentEndPoint,
		CollectorEndpoint: options.collectorEndPoint,
		Username:          options.username,
		Password:          options.password,
		Process: jaeger.Process{
			ServiceName: options.serviceName,
			Tags:        options.tags,
		},
		OnError: options.onError,
	})
	if err != nil {
		return nil, err
	}

	//var le *LogExporter
	//le, _ = NewLogExporter(LoggerOptions{TracesLogFile: options.traceLogFile})
	//exporter := newExporterWrapper(je, le, options.sampleType, options.sampleParam)

	xdGen := initIDGenerator()
	trace.RegisterExporter(je)
	trace.ApplyConfig(trace.Config{DefaultSampler: options.sampler, IDGenerator: xdGen})

	enable = true
	return func() {
		je.Flush()
	}, nil
}

func ExtractSpan(ctx context.Context) *trace.Span {
	span := trace.FromContext(ctx)
	return span
}

type TraceID trace.TraceID

func (t TraceID) String() string {
	if t == EmptyTraceID {
		return ""
	}
	return fmt.Sprintf("%02x", t[:])
}

func (t TraceID) Empty() bool {
	return t == EmptyTraceID
}

var EmptyTraceID = TraceID([16]byte{})

func ExtractTraceID(ctx context.Context) TraceID {
	span := trace.FromContext(ctx)
	if span == nil {
		return EmptyTraceID
	}
	return TraceID(span.SpanContext().TraceID)
}

type Attribute = trace.Attribute
type Status = trace.Status

var (
	BoolAttribute    = trace.BoolAttribute
	Int64Attribute   = trace.Int64Attribute
	Float64Attribute = trace.Float64Attribute
	StringAttribute  = trace.StringAttribute
)

func MustNewSpan(ctx context.Context, name string) (context.Context, func()) {
	ctx, span := trace.StartSpan(ctx, name)
	return ctx, span.End
}

// NewSpan crete a child span if ctx has trace info
func NewSpan(ctx context.Context, name string) (context.Context, func()) {
	span := trace.FromContext(ctx)
	if span == nil {
		return ctx, func() {}
	}
	ctx, span = trace.StartSpan(ctx, name)
	return ctx, span.End
}

func AddAttributes(ctx context.Context, attributes ...Attribute) {
	span := trace.FromContext(ctx)
	if span == nil {
		return
	}
	span.AddAttributes(attributes...)
}

func SetStatus(ctx context.Context, status Status) {
	span := trace.FromContext(ctx)
	if span == nil {
		return
	}
	span.SetStatus(status)
}

func SetError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	span := trace.FromContext(ctx)
	if span == nil {
		return
	}
	xderr, ok := err.(*xd_error.XDError)
	if ok {
		span.SetStatus(trace.Status{Code: int32(xderr.Code), Message: xderr.Message})
	} else {
		span.SetStatus(trace.Status{Code: trace.StatusCodeUnknown, Message: err.Error()})
	}
}

func Log(ctx context.Context, fields []Attribute, msg string) {
	span := trace.FromContext(ctx)
	if span == nil {
		return
	}
	span.Annotate(fields, msg)
}

func Logf(ctx context.Context, fields []Attribute, format string, info ...interface{}) {
	span := trace.FromContext(ctx)
	if span == nil {
		return
	}
	span.Annotatef(fields, format, info...)
}

func CopySpanToNewCtx(ctx context.Context) context.Context {
	// 灰度流量标签在metadata中，copy trace时一并复制到新的ctx中
	md := metadata.FromContext(ctx)
	newMd := copyMetadata(md)
	newCtx := context.Background()
	newCtx = metadata.WithMetadata(newCtx, newMd)

	span := trace.FromContext(ctx)
	if span == nil {
		return newCtx
	}
	return trace.NewContext(newCtx, span)
}

func copyMetadata(m metadata.Metadata) metadata.Metadata {
	return metadata.NewMetadata(m.GetD())
}

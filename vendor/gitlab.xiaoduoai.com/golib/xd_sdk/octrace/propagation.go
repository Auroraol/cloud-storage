package octrace

import (
	"encoding/hex"
	"errors"

	"go.opencensus.io/trace"
)

type SpanContextMap map[string]string

const (
	TraceIDKey = "X_OCTRACE_TRACEID"
	SpanIDKey  = "X_OCTRACE_SPANID"
	SampleKey  = "X_OCTRACE_SAMPLE"
)

var (
	MapNilError       = errors.New("map is nil")
	ParseTraceIDError = errors.New("parse traceid failed.")
	ParseSpanIDError  = errors.New("parse spanid failed.")
	ParseSampledError = errors.New("parse sampled failed.")
)

func (s *SpanContextMap) FromSpanContext(sc trace.SpanContext) {
	if *s == nil {
		*s = make(SpanContextMap)
	} else {
		delete(*s, TraceIDKey)
		delete(*s, SpanIDKey)
		delete(*s, SampleKey)
	}
	t := trace.SpanContext{}
	if sc == t {
		return
	}
	(*s)[TraceIDKey] = hex.EncodeToString(sc.TraceID[:])
	(*s)[SpanIDKey] = hex.EncodeToString(sc.SpanID[:])
	(*s)[SampleKey] = EncodeSampled(sc.TraceOptions)
}

func (s SpanContextMap) ToSpanContext() (trace.SpanContext, error) {
	var ctx trace.SpanContext
	if s == nil {
		return ctx, MapNilError
	}
	traceID, ok := ParseTraceID(s[TraceIDKey])
	if !ok {
		return ctx, ParseTraceIDError
	}
	spanID, ok := ParseSpanID(s[SpanIDKey])
	if !ok {
		return ctx, ParseSpanIDError
	}
	traceOpt, ok := ParseSampled(s[SampleKey])
	if !ok {
		return ctx, ParseSampledError
	}

	return trace.SpanContext{
		TraceID:      traceID,
		SpanID:       spanID,
		TraceOptions: traceOpt,
	}, nil
}

func (s SpanContextMap) TraceID() string {
	t, _ := s[TraceIDKey]
	return t
}

func (s SpanContextMap) SpanID() string {
	t, _ := s[SpanIDKey]
	return t
}

func (s SpanContextMap) Sampled() string {
	t, _ := s[SampleKey]
	return t
}

// ParseTraceID parses the value of the X-B3-TraceId header.
func ParseTraceID(tid string) (trace.TraceID, bool) {
	if tid == "" {
		return trace.TraceID{}, false
	}
	b, err := hex.DecodeString(tid)
	if err != nil {
		return trace.TraceID{}, false
	}
	var traceID trace.TraceID
	if len(b) <= 8 {
		// The lower 64-bits.
		start := 8 + (8 - len(b))
		copy(traceID[start:], b)
	} else {
		start := 16 - len(b)
		copy(traceID[start:], b)
	}

	return traceID, true
}

// ParseSpanID parses the value of the X-B3-SpanId or X-B3-ParentSpanId headers.
func ParseSpanID(sid string) (spanID trace.SpanID, ok bool) {
	if sid == "" {
		return trace.SpanID{}, false
	}
	b, err := hex.DecodeString(sid)
	if err != nil {
		return trace.SpanID{}, false
	}
	start := 8 - len(b)
	copy(spanID[start:], b)
	return spanID, true
}

// ParseSampled parses the value of the X-B3-Sampled header.
func ParseSampled(sampled string) (trace.TraceOptions, bool) {
	switch sampled {
	case "true", "1":
		return trace.TraceOptions(1), true
	default:
		return trace.TraceOptions(0), true
	}
}

func EncodeSampled(opt trace.TraceOptions) string {
	if opt.IsSampled() {
		return "true"
	}
	return "false"
}

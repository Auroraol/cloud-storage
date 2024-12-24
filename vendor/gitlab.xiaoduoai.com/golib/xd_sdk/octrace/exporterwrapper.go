package octrace

import (
	"encoding/binary"

	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

type checker func(traceID trace.TraceID) bool
type ExporterWrapper struct {
	jaeger   *jaeger.Exporter
	logger   *LogExporter
	needPush checker
}

func newExporterWrapper(jaeger *jaeger.Exporter, logger *LogExporter, sampleType SampleType, fraction float64) *ExporterWrapper {
	ret := &ExporterWrapper{jaeger: jaeger, logger: logger}
	switch sampleType {
	case AlwaysSample:
		ret.needPush = func(traceID trace.TraceID) bool { return true }
	case NeverSample:
		ret.needPush = func(traceID trace.TraceID) bool { return false }
	default:
		if fraction >= 1 {
			ret.needPush = func(traceID trace.TraceID) bool { return true }
		} else if fraction <= 0 {
			ret.needPush = func(traceID trace.TraceID) bool { return false }
		} else {
			ret.needPush = func() checker {
				traceIDUpperBound := uint64(fraction * (1 << 63))
				return func(traceID trace.TraceID) bool {
					x := binary.BigEndian.Uint64(traceID[0:8]) >> 1
					return x < traceIDUpperBound
				}
			}()
		}
	}
	return ret
}

func (e *ExporterWrapper) ExportSpan(s *trace.SpanData) {
	if e.logger != nil {
		e.logger.ExportSpan(s)
	}
	if e.needPush(s.TraceID) {
		e.jaeger.ExportSpan(s)
	}
}

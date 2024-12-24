package filter

import (
	"strings"
)

type TraceFilter interface {
	NotTrace(path string) bool
}

var traceFilter TraceFilter

func NotTrace(path string) bool {
	if traceFilter == nil {
		return false
	}
	return traceFilter.NotTrace(path)
}

func SetFilter(filter TraceFilter) {
	traceFilter = filter
}

type DefaultFilter struct {
	filter map[string]struct{}
}

func (f DefaultFilter) NotTrace(path string) bool {
	p := strings.ToLower(path)
	_, ok := f.filter[p]
	return ok
}

func init() {
	traceFilter = DefaultFilter{
		filter: map[string]struct{}{
			strings.ToLower("/healthcheck1"): struct{}{},
			strings.ToLower("Heartbeat"):     struct{}{},
			//strings.ToLower("sayhello"):    struct{}{},
		},
	}
}

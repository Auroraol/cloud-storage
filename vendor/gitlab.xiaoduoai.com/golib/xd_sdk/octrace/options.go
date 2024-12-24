package octrace

import "contrib.go.opencensus.io/exporter/jaeger"
import "go.opencensus.io/trace"

type SampleType int

const (
	AlwaysSample SampleType = iota
	NeverSample
	ProbabilitySample
	RateLimitSample
)

type Options struct {
	sampler           trace.Sampler
	sampleType        SampleType
	sampleParam       float64
	agentEndPoint     string
	collectorEndPoint string
	serviceName       string
	tags              []jaeger.Tag
	username          string
	password          string
	onError           func(error)
	traceLogFile      string
}

type Option func(opts *Options)

func WithTraceLogFile(file string) Option {
	return func(opts *Options) {
		opts.traceLogFile = file
	}
}

func WithSample(sample SampleType, fraction float64) Option {
	return func(opts *Options) {
		switch sample {
		case AlwaysSample:
			opts.sampler = trace.AlwaysSample()
		case NeverSample:
			opts.sampler = trace.NeverSample()
		case ProbabilitySample:
			opts.sampler = trace.ProbabilitySampler(fraction)
		case RateLimitSample:
			opts.sampler = RateLimitingSample(int64(fraction))
		}
	}
}

func WithEndPoint(agent string, collector string) Option {
	return func(opts *Options) {
		opts.agentEndPoint = agent
		opts.collectorEndPoint = collector
	}
}

func WithServiceName(serviceName string) Option {
	return func(opts *Options) {
		opts.serviceName = serviceName
	}
}

func WithTags(tags ...jaeger.Tag) Option {
	return func(opts *Options) {
		opts.tags = append(opts.tags, tags...)
	}
}

func WithAuthor(username, password string) Option {
	return func(opts *Options) {
		opts.username = username
		opts.password = password
	}
}

func WithOnErr(errFunc func(error)) Option {
	return func(opts *Options) {
		opts.onError = errFunc
	}
}

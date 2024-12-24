package octrace

import (
	"math"
	"sync"
	"time"

	"go.opencensus.io/trace"
)

// RateLimitingSample only samples with given rate(qps).
// It also samples spans whose parents are sampled.
func RateLimitingSample(qps int64) trace.Sampler {
	lastSampleTime := time.Now()
	lock := sync.Mutex{}
	return func(p trace.SamplingParameters) trace.SamplingDecision {
		if p.ParentContext.IsSampled() {
			return trace.SamplingDecision{Sample: true}
		}

		lock.Lock()
		defer lock.Unlock()
		elapseSec := time.Since(lastSampleTime).Seconds()
		fraction := math.Min(elapseSec*float64(qps), 1.0)
		if fraction >= 1 {
			lastSampleTime = time.Now()
		}

		return trace.SamplingDecision{Sample: fraction >= 1}
	}
}

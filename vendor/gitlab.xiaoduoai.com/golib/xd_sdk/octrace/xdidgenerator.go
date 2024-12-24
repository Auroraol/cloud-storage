package octrace

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type xdIDGenerator struct {
	sync.Mutex

	// Please keep these as the first fields
	// so that these 8 byte fields will be aligned on addresses
	// divisible by 8, on both 32-bit and 64-bit machines when
	// performing atomic increments and accesses.
	// See:
	// * https://github.com/census-instrumentation/opencensus-go/issues/587
	// * https://github.com/census-instrumentation/opencensus-go/issues/865
	// * https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	nextSpanID uint64
	spanIDInc  uint64

	traceIDAdd  [2]uint64
	traceIDRand *rand.Rand
}


func initIDGenerator() *xdIDGenerator {
	gen := &xdIDGenerator{}
	// initialize traceID and spanID generators.
	var rngSeed int64
	for _, p := range []interface{}{
		&rngSeed, &gen.traceIDAdd, &gen.nextSpanID, &gen.spanIDInc,
	} {
		binary.Read(crand.Reader, binary.LittleEndian, p)
	}
	gen.traceIDRand = rand.New(rand.NewSource(rngSeed))
	gen.spanIDInc |= 1
	return gen
}

// NewSpanID returns a non-zero span ID from a randomly-chosen sequence.
func (gen *xdIDGenerator) NewSpanID() [8]byte {
	var id uint64
	for id == 0 {
		id = atomic.AddUint64(&gen.nextSpanID, gen.spanIDInc)
	}
	var sid [8]byte
	binary.LittleEndian.PutUint64(sid[:], id)
	return sid
}

// NewTraceID returns a non-zero trace ID from a randomly-chosen sequence.
// mu should be held while this function is called.
func (gen *xdIDGenerator) NewTraceID() [16]byte {
	var tid [16]byte
	// Construct the trace ID from two outputs of traceIDRand, with a constant
	// added to each half for additional entropy.
	gen.Lock()
	n, _ := strconv.ParseUint(fmt.Sprintf("%d", time.Now().Unix()), 16, 64)
	binary.LittleEndian.PutUint64(tid[0:8], gen.traceIDRand.Uint64()+gen.traceIDAdd[0])
	binary.BigEndian.PutUint64(tid[8:16], n)
	gen.Unlock()
	return tid
}
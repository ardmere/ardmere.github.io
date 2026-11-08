package verifier

import (
	"context"
	"log"
	"sync/atomic"
	"time"
)

const (
	defaultProgressEveryN      = 25
	defaultProgressEveryNHuge  = 50000 // when total is unknown (e.g. CSV stream)
	defaultProgressInterval    = 30 * time.Second
)

// Progress logs heartbeat lines for long-running on-chain balance work.
// Safe for concurrent Tick from worker goroutines.
type Progress struct {
	verifier string
	phase    string
	total    int64
	done     atomic.Int64
	start    time.Time
	everyN   int64
	interval time.Duration
	lastLog  atomic.Int64 // unix nanos of last progress log
	coin     atomic.Value // string
	network  atomic.Value // string
}

// NewProgress starts a progress tracker. Call Start immediately after construction.
// Pass total=0 when the row count is not known yet (CSV streaming); logging then
// relies on a larger every-N stride plus the ~30s heartbeat.
func NewProgress(verifierID, phase string, total int) *Progress {
	everyN := int64(defaultProgressEveryN)
	if total <= 0 {
		everyN = defaultProgressEveryNHuge
	}
	p := &Progress{
		verifier: verifierID,
		phase:    phase,
		total:    int64(total),
		start:    time.Now(),
		everyN:   everyN,
		interval: defaultProgressInterval,
	}
	p.coin.Store("")
	p.network.Store("")
	return p
}

// Start logs the beginning of a phase.
func (p *Progress) Start(detail string) {
	if p == nil {
		return
	}
	p.start = time.Now()
	p.lastLog.Store(p.start.UnixNano())
	if detail == "" {
		log.Printf("  %s: start %s total=%d", p.verifier, p.phase, p.total)
		return
	}
	log.Printf("  %s: start %s total=%d %s", p.verifier, p.phase, p.total, detail)
}

// Tick records one completed unit and may emit a progress line.
func (p *Progress) Tick(coin, network string) {
	if p == nil {
		return
	}
	n := p.done.Add(1)
	if coin != "" {
		p.coin.Store(coin)
	}
	if network != "" {
		p.network.Store(network)
	}
	if !p.shouldLog(n) {
		return
	}
	p.emit("progress", n)
}

// Heartbeat emits a progress line every interval until stop is called.
// Use around opaque work (CSV sample) or RPC fan-out so a hang before the
// first Tick still surfaces which verifier/phase is running.
func (p *Progress) Heartbeat(ctx context.Context) (stop func()) {
	if p == nil {
		return func() {}
	}
	done := make(chan struct{})
	go func() {
		t := time.NewTicker(p.interval)
		defer t.Stop()
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			case <-t.C:
				p.emit("heartbeat", p.done.Load())
			}
		}
	}()
	return func() { close(done) }
}

// Finish logs completion of the phase.
func (p *Progress) Finish() {
	if p == nil {
		return
	}
	n := p.done.Load()
	if n == 0 && p.total > 0 {
		// Still emit finish so hang-before-first-row is distinguishable from
		// a finish that never arrived.
	}
	p.emit("finish", n)
}

func (p *Progress) shouldLog(n int64) bool {
	if n == p.total {
		return true
	}
	if p.everyN > 0 && n%p.everyN == 0 {
		return true
	}
	now := time.Now().UnixNano()
	last := p.lastLog.Load()
	if time.Duration(now-last) >= p.interval {
		if p.lastLog.CompareAndSwap(last, now) {
			return true
		}
	}
	return false
}

func (p *Progress) emit(kind string, n int64) {
	p.lastLog.Store(time.Now().UnixNano())
	coin, _ := p.coin.Load().(string)
	network, _ := p.network.Load().(string)
	elapsed := time.Since(p.start).Truncate(time.Millisecond)
	pair := ""
	if coin != "" || network != "" {
		pair = " coin=" + coin + " network=" + network
	}
	total := p.total
	if total <= 0 {
		log.Printf("  %s: %s %s completed=%d elapsed=%s%s",
			p.verifier, kind, p.phase, n, elapsed, pair)
		return
	}
	log.Printf("  %s: %s %s completed=%d/%d elapsed=%s%s",
		p.verifier, kind, p.phase, n, total, elapsed, pair)
}

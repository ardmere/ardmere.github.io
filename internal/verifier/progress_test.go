package verifier

import (
	"sync"
	"testing"
	"time"
)

func TestProgressTickFinish(t *testing.T) {
	p := NewProgress("onchain-balance-ledger", "rpc", 3)
	p.everyN = 1
	p.interval = time.Hour
	p.Start("test")
	p.Tick("BTC", "BTC")
	p.Tick("ETH", "ETH")
	p.Tick("SOL", "SOL")
	p.Finish()
	if got := p.done.Load(); got != 3 {
		t.Fatalf("done=%d want 3", got)
	}
}

func TestProgressConcurrentTicks(t *testing.T) {
	const n = 100
	p := NewProgress("onchain-balance-deposit", "rpc", n)
	p.everyN = 1000
	p.interval = time.Hour
	p.Start("concurrent")
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.Tick("USDT", "ETH")
		}()
	}
	wg.Wait()
	p.Finish()
	if got := p.done.Load(); got != n {
		t.Fatalf("done=%d want %d", got, n)
	}
}

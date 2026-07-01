package rpc

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ResultCache stores immutable RPC results keyed by (network, target, method, data, height).
type ResultCache struct {
	dir string
	mu  sync.Mutex
	mem map[string][]byte
}

type cacheEntry struct {
	Result []byte `json:"result"`
	Used   string `json:"used,omitempty"`
}

// NewResultCache returns an in-memory cache, optionally backed by disk under dir.
func NewResultCache(dir string) *ResultCache {
	if dir == "" {
		dir = os.Getenv("RPC_CACHE_DIR")
	}
	if dir == "" {
		dir = ".cache/rpc"
	}
	return &ResultCache{dir: dir, mem: map[string][]byte{}}
}

func (c *ResultCache) key(net Network, method, target, dataHex string, height int64) string {
	raw := fmt.Sprintf("%s|%s|%s|%s|%d", net, method, strings.ToLower(target), dataHex, height)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (c *ResultCache) Get(net Network, method, target, dataHex string, height int64) ([]byte, bool) {
	key := c.key(net, method, target, dataHex, height)
	c.mu.Lock()
	if raw, ok := c.mem[key]; ok {
		c.mu.Unlock()
		return raw, true
	}
	c.mu.Unlock()

	if c.dir == "" {
		return nil, false
	}
	path := filepath.Join(c.dir, key+".json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	var entry cacheEntry
	if err := json.Unmarshal(raw, &entry); err != nil {
		return nil, false
	}
	c.mu.Lock()
	c.mem[key] = entry.Result
	c.mu.Unlock()
	return entry.Result, true
}

func (c *ResultCache) Put(net Network, method, target, dataHex string, height int64, result []byte) {
	if len(result) == 0 {
		return
	}
	key := c.key(net, method, target, dataHex, height)
	c.mu.Lock()
	c.mem[key] = append([]byte(nil), result...)
	c.mu.Unlock()

	if c.dir == "" {
		return
	}
	_ = os.MkdirAll(c.dir, 0o755)
	path := filepath.Join(c.dir, key+".json")
	entry, _ := json.Marshal(cacheEntry{Result: result})
	_ = os.WriteFile(path, entry, 0o644)
}

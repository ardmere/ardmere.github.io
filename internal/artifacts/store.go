package artifacts

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/exchange"
)

// StoreOpts controls optional FetchAndStore behavior.
// Reserved for future flags; callers currently pass a zero value.
type StoreOpts struct{}

// FetchMeta is written to fetch.json after a successful archive.
type FetchMeta struct {
	Exchange   string            `json:"exchange"`
	SnapshotID string            `json:"snapshotId"`
	FetchedAt  string            `json:"fetchedAt"`
	Source     string            `json:"source"`
	Artifacts  []bundle.Artifact `json:"artifacts,omitempty"`
}

// FetchAndStore runs adapter.FetchCurrentSnapshot into a temp dir, then
// archives files under artifacts/<id>/<auditId>/raw/ and writes fetch.json.
func FetchAndStore(ctx context.Context, adapter exchange.Adapter, artifactsBase string, opts exchange.FetchOpts, _ StoreOpts) (exchange.StoreResult, error) {
	if artifactsBase == "" {
		artifactsBase = "./artifacts"
	}
	tmp, err := os.MkdirTemp("", "por-fetch-*")
	if err != nil {
		return exchange.StoreResult{}, fmt.Errorf("temp dir: %w", err)
	}
	defer os.RemoveAll(tmp)

	fetch, err := adapter.FetchCurrentSnapshot(ctx, tmp, opts)
	if err != nil {
		return exchange.StoreResult{}, err
	}
	if fetch.Snapshot.ID == "" {
		return exchange.StoreResult{}, fmt.Errorf("fetch returned empty snapshot id")
	}

	snapshotDir := SnapshotRoot(artifactsBase, adapter.ID(), fetch.Snapshot.ID)
	rawDir := filepath.Join(snapshotDir, "raw")
	bundlesDir := filepath.Join(snapshotDir, "bundles")
	if err := os.MkdirAll(rawDir, 0o755); err != nil {
		return exchange.StoreResult{}, fmt.Errorf("mkdir raw: %w", err)
	}
	if err := os.MkdirAll(bundlesDir, 0o755); err != nil {
		return exchange.StoreResult{}, fmt.Errorf("mkdir bundles: %w", err)
	}

	// Copy every file present in the temp dir into raw/ (content-addressed names).
	files, err := listFiles(tmp)
	if err != nil {
		return exchange.StoreResult{}, err
	}
	copied := map[string]string{} // abs tmp path → abs raw path
	for _, src := range files {
		base := filepath.Base(src)
		dst := filepath.Join(rawDir, base)
		if err := copyFile(src, dst); err != nil {
			return exchange.StoreResult{}, fmt.Errorf("copy %s: %w", base, err)
		}
		copied[src] = dst
	}

	rewrite := func(p string) string {
		if p == "" {
			return ""
		}
		if abs, ok := copied[p]; ok {
			return abs
		}
		// Match by basename when paths differ slightly.
		base := filepath.Base(p)
		cand := filepath.Join(rawDir, base)
		if _, err := os.Stat(cand); err == nil {
			return cand
		}
		return p
	}

	for i := range fetch.Artifacts {
		abs := rewrite(fetch.Artifacts[i].LocalPath)
		fetch.Artifacts[i].LocalPath = filepath.Join("raw", filepath.Base(abs))
		// Keep absolute path available via Summary/Wallet fields below.
		_ = abs
	}
	fetch.SummaryPath = rewrite(fetch.SummaryPath)
	fetch.WalletZipPath = rewrite(fetch.WalletZipPath)

	meta := FetchMeta{
		Exchange:   adapter.ID(),
		SnapshotID: fetch.Snapshot.ID,
		FetchedAt:  time.Now().UTC().Format(time.RFC3339),
		Source:     defaultSource(opts),
		Artifacts:  fetch.Artifacts,
	}
	metaBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return exchange.StoreResult{}, fmt.Errorf("marshal fetch.json: %w", err)
	}
	if err := os.WriteFile(filepath.Join(snapshotDir, "fetch.json"), metaBytes, 0o644); err != nil {
		return exchange.StoreResult{}, fmt.Errorf("write fetch.json: %w", err)
	}

	return exchange.StoreResult{
		Fetch:       fetch,
		SnapshotDir: snapshotDir,
		RawDir:      rawDir,
		BundlesDir:  bundlesDir,
	}, nil
}

func defaultSource(opts exchange.FetchOpts) string {
	if opts.ImportSource != "" {
		return opts.ImportSource
	}
	if opts.SummaryPath != "" || opts.WalletZipPath != "" || opts.ZkBundlePath != "" || opts.UserProofPath != "" {
		return "local-import"
	}
	return "public-api"
}

func listFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		out = append(out, filepath.Join(dir, e.Name()))
	}
	return out, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

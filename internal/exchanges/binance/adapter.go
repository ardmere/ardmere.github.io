package binance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ardmere/ardmere/internal/exchanges/binance/bapi"
	"github.com/ardmere/ardmere/internal/artifacts"
	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/exchange"
	"github.com/ardmere/ardmere/internal/por"
	"github.com/ardmere/ardmere/internal/walletzip"
)

const ID = "binance"

type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) ID() string { return ID }

func (a *Adapter) Capabilities() exchange.Capabilities {
	return exchange.Capabilities{
		Tier:             exchange.TierOnchain,
		WalletZip:        true,
		OnchainHot:       true,
		OnchainToken:     true,
		OnchainLedger:    true,
	}
}

func (a *Adapter) FetchAndStore(ctx context.Context, artifactsBase string, opts exchange.FetchOpts) (exchange.StoreResult, error) {
	return artifacts.FetchAndStore(ctx, a, artifactsBase, opts, artifacts.StoreOpts{})
}

func (a *Adapter) AggregateWalletZip(path string) (*walletzip.Aggregate, error) {
	return walletzip.AggregateFiles(path, walletzip.HotCold, walletzip.Deposit)
}

func (a *Adapter) VerifierProfile() exchange.VerifierProfile {
	return exchange.VerifierProfile{
		Shared: []string{
			"artifact-integrity@1",
			"internal-consistency@1.1",
			"btc-anchor@1",
			"solvency-claim@1",
			"onchain-balance-hot@2.1",
			"onchain-balance-token@2.0",
			"onchain-balance-ledger@1.4",
			"onchain-balance-deposit@1.2",
		},
		Stubs: []string{
			"address-ownership@0",
			"global-zk-proof@0",
			"third-party-attestation@0",
			"cross-chain-wrapped@0",
		},
	}
}

func (a *Adapter) FetchCurrentSnapshot(ctx context.Context, outDir string, opts exchange.FetchOpts) (exchange.FetchResult, error) {
	if opts.AuditID != "" {
		return a.fetchByAuditID(ctx, outDir, opts)
	}

	bapiBytes, raw, err := bapi.FetchSnapshot(ctx)
	if err != nil {
		return exchange.FetchResult{}, err
	}
	return a.storeBAPISnapshot(ctx, outDir, opts, bapiBytes, raw)
}

func (a *Adapter) fetchByAuditID(ctx context.Context, outDir string, opts exchange.FetchOpts) (exchange.FetchResult, error) {
	bapiBytes, raw, err := bapi.FetchSnapshot(ctx)
	if err == nil && raw.AuditID == opts.AuditID {
		return a.storeBAPISnapshot(ctx, outDir, opts, bapiBytes, raw)
	}

	entries, err := bapi.FetchRecentListEntries(ctx, 36)
	if err != nil {
		return exchange.FetchResult{}, err
	}
	var entry *bapi.ListEntry
	for i := range entries {
		if entries[i].AuditID == opts.AuditID {
			entry = &entries[i]
			break
		}
	}
	if entry == nil {
		return exchange.FetchResult{}, fmt.Errorf("binance audit %s not found in public snapshot list", opts.AuditID)
	}

	raw = bapi.Snapshot{
		AuditID:      opts.AuditID,
		SnapshotTime: entry.SnapshotTime,
		AuditDate:    entry.SnapshotTimeUTC.Format("01/02/06"),
	}
	bapiBytes, err = json.Marshal(map[string]any{
		"code": "000000",
		"data": raw,
		"note": "wallet-only historical fetch; BAPI summary unavailable for non-current snapshot",
	})
	if err != nil {
		return exchange.FetchResult{}, err
	}
	return a.storeBAPISnapshot(ctx, outDir, opts, bapiBytes, raw)
}

func (a *Adapter) storeBAPISnapshot(ctx context.Context, outDir string, opts exchange.FetchOpts, bapiBytes []byte, raw bapi.Snapshot) (exchange.FetchResult, error) {
	bapiSum := sha256.Sum256(bapiBytes)
	bapiSumHex := hex.EncodeToString(bapiSum[:])
	bapiPath := filepath.Join(outDir, bapiSumHex+".json")
	if err := os.WriteFile(bapiPath, bapiBytes, 0o644); err != nil {
		return exchange.FetchResult{}, err
	}

	meta, err := bapi.FetchSnapshotMeta(ctx, raw.AuditID, raw)
	if err != nil {
		meta = bapi.SnapshotMeta{}
	}
	if meta.SnapshotTime.IsZero() {
		if t, e := bapi.ParseSnapshotTime(raw.SnapshotTime); e == nil {
			meta.SnapshotTime = t
		}
	}

	snap := Normalize(raw, Meta{
		Exchange:       ID,
		PeriodSeq:      meta.PeriodSeq,
		BTCBlockHeight: meta.BTCBlockHeight,
		SnapshotTime:   meta.SnapshotTime,
	})

	arts := []bundle.Artifact{{
		Kind:      por.KindBapiSnapshot,
		SHA256:    bapiSumHex,
		URL:       bapi.SnapshotURL(),
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
		SizeBytes: int64(len(bapiBytes)),
		LocalPath: bapiPath,
	}}

	out := exchange.FetchResult{
		Snapshot:    snap,
		Artifacts:   arts,
		SummarySha:  bapiSumHex,
		SummaryPath: bapiPath,
	}

	if opts.SkipWalletZip {
		return out, nil
	}

	zipURL, err := bapi.FetchWalletZipURL(ctx, raw.AuditID)
	if err != nil {
		return exchange.FetchResult{}, fmt.Errorf("wallet zip url: %w", err)
	}
	path, sumHex, size, err := walletzip.Download(ctx, zipURL, outDir)
	if err != nil {
		return exchange.FetchResult{}, fmt.Errorf("wallet zip: %w", err)
	}

	out.WalletZipPath = path
	out.WalletZipSha = sumHex
	out.Artifacts = append(out.Artifacts, bundle.Artifact{
		Kind:      por.KindWalletZip,
		SHA256:    sumHex,
		URL:       zipURL,
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
		SizeBytes: size,
		LocalPath: path,
	})
	return out, nil
}

func (a *Adapter) ParseSnapshotFromArtifacts(artBundle bundle.ArtifactBundle, artifactsDir string) (por.Snapshot, exchange.ParsedArtifacts, error) {
	var parsed exchange.ParsedArtifacts
	var summaryPath string

	for _, art := range artBundle.Artifacts {
		switch art.Kind {
		case por.KindBapiSnapshot, por.KindSummarySnapshot:
			parsed.SummarySha = art.SHA256
			summaryPath = resolvePath(artifactsDir, art.LocalPath)
		case por.KindWalletZip, por.KindWalletAddressBundle:
			parsed.WalletZipSha = art.SHA256
			parsed.WalletZipPath = resolvePath(artifactsDir, art.LocalPath)
		}
	}
	if summaryPath == "" {
		return por.Snapshot{}, parsed, fmt.Errorf("artifact bundle missing summary snapshot (bapiSnapshot or summarySnapshot)")
	}
	parsed.SummaryPath = summaryPath

	raw, err := bapi.LoadSnapshot(summaryPath)
	if err != nil {
		return por.Snapshot{}, parsed, err
	}

	snapTime, _ := time.Parse(time.RFC3339, artBundle.SnapshotTime)
	snap := Normalize(raw, Meta{
		Exchange:       artBundle.Exchange,
		PeriodSeq:      artBundle.PeriodSeq,
		BTCBlockHeight: artBundle.BTCBlockHeight,
		SnapshotTime:   snapTime,
	})
	if artBundle.Exchange == "" {
		snap.Exchange = ID
	}
	return snap, parsed, nil
}

func resolvePath(artifactsDir, localPath string) string {
	if localPath == "" {
		return ""
	}
	if filepath.IsAbs(localPath) {
		return localPath
	}
	if _, err := os.Stat(localPath); err == nil {
		return localPath
	}
	return filepath.Join(artifactsDir, filepath.Base(localPath))
}

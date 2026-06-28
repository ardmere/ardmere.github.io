package okx

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ardmere/ardmere/internal/artifacts"
	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/exchange"
	"github.com/ardmere/ardmere/internal/exchanges/okx/okxapi"
	"github.com/ardmere/ardmere/internal/por"
	"github.com/ardmere/ardmere/internal/walletzip"
)

type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) ID() string { return ID }

func (a *Adapter) Capabilities() exchange.Capabilities {
	return exchange.Capabilities{
		Tier:             exchange.TierWallet,
		WalletZip:        true,
		AddressOwnership: true,
		GlobalZK:         true,
		OnchainHot:       true,
		OnchainToken:     true,
	}
}

func (a *Adapter) FetchAndStore(ctx context.Context, artifactsBase string, opts exchange.FetchOpts) (exchange.StoreResult, error) {
	return artifacts.FetchAndStore(ctx, a, artifactsBase, opts, artifacts.StoreOpts{})
}

func (a *Adapter) AggregateWalletZip(path string) (*walletzip.Aggregate, error) {
	return walletzip.AggregateOKXZip(path)
}

func (a *Adapter) VerifierProfile() exchange.VerifierProfile {
	return exchange.VerifierProfile{
		Shared: []string{
			"artifact-integrity@1",
			"solvency-claim@1",
			"internal-consistency@1",
			"address-ownership@okx-1",
			"onchain-balance-hot@2.1",
			"onchain-balance-token@2.0",
			"onchain-balance-ledger@1.2",
			"global-zk-proof@okx-1",
		},
		Stubs: []string{
			"btc-anchor@0",
			"third-party-attestation@0",
			"cross-chain-wrapped@0",
		},
	}
}

func (a *Adapter) FetchCurrentSnapshot(ctx context.Context, outDir string, opts exchange.FetchOpts) (exchange.FetchResult, error) {
	var summaryBytes []byte
	var parsed okxapi.SummaryBundle
	var err error

	if opts.SummaryPath != "" {
		summaryBytes, err = os.ReadFile(opts.SummaryPath)
		if err != nil {
			return exchange.FetchResult{}, fmt.Errorf("read summary: %w", err)
		}
		parsed, err = okxapi.ParseSummaryBytes(summaryBytes)
		if err != nil {
			return exchange.FetchResult{}, err
		}
	} else {
		parsed, summaryBytes, err = okxapi.FetchCurrentSummary(ctx)
		if err != nil {
			return exchange.FetchResult{}, fmt.Errorf("fetch okx summary: %w", err)
		}
	}

	sum := sha256.Sum256(summaryBytes)
	sumHex := hex.EncodeToString(sum[:])
	sumPath := filepath.Join(outDir, sumHex+".json")
	if err := os.WriteFile(sumPath, summaryBytes, 0o644); err != nil {
		return exchange.FetchResult{}, err
	}

	snap := Normalize(parsed.Audit, Meta{Exchange: ID})
	arts := []bundle.Artifact{{
		Kind:      por.KindSummarySnapshot,
		SHA256:    sumHex,
		URL:       okxapi.URLDetail,
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
		SizeBytes: int64(len(summaryBytes)),
		LocalPath: sumPath,
	}}

	out := exchange.FetchResult{
		Snapshot:    snap,
		Artifacts:   arts,
		SummarySha:  sumHex,
		SummaryPath: sumPath,
	}

	if opts.SkipWalletZip && opts.WalletZipPath == "" {
		return a.maybeAttachLiability(ctx, outDir, opts, snap, out)
	}

	walletURL := opts.WalletZipPath
	if walletURL == "" {
		walletURL, _, err = okxapi.ResolveDownloads(ctx, snap.ID)
		if err != nil {
			return exchange.FetchResult{}, fmt.Errorf("resolve wallet url: %w", err)
		}
	}

	var walletPath, walletSha string
	var walletSize int64
	if opts.WalletZipPath != "" && fileExists(opts.WalletZipPath) {
		walletPath, walletSha, walletSize, err = importLocalFile(opts.WalletZipPath, outDir, ".zip")
	} else {
		walletPath, walletSha, walletSize, err = okxapi.DownloadFile(ctx, walletURL, outDir, ".zip")
	}
	if err != nil {
		return exchange.FetchResult{}, fmt.Errorf("wallet zip: %w", err)
	}
	out.WalletZipPath = walletPath
	out.WalletZipSha = walletSha
	out.Artifacts = append(out.Artifacts, bundle.Artifact{
		Kind:      por.KindWalletZip,
		SHA256:    walletSha,
		URL:       walletURL,
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
		SizeBytes: walletSize,
		LocalPath: walletPath,
	})

	return a.maybeAttachLiability(ctx, outDir, opts, snap, out)
}

func (a *Adapter) maybeAttachLiability(ctx context.Context, outDir string, opts exchange.FetchOpts, snap por.Snapshot, out exchange.FetchResult) (exchange.FetchResult, error) {
	if opts.SkipLiabilityZip {
		return out, nil
	}
	liabRef := opts.LiabilityZipPath
	if liabRef == "" {
		var err error
		_, liabRef, err = okxapi.ResolveDownloads(ctx, snap.ID)
		if err != nil {
			return exchange.FetchResult{}, fmt.Errorf("resolve liability url: %w", err)
		}
	}

	var path, sha string
	var size int64
	var err error
	switch {
	case fileExists(liabRef):
		path, sha, size, err = importLocalFile(liabRef, outDir, ".zip")
	case isHTTP(liabRef):
		path, sha, size, err = okxapi.DownloadFile(ctx, liabRef, outDir, ".zip")
	default:
		return exchange.FetchResult{}, fmt.Errorf("liability zip not found: %s", liabRef)
	}
	if err != nil {
		return exchange.FetchResult{}, fmt.Errorf("liability zip: %w", err)
	}
	out.Artifacts = append(out.Artifacts, bundle.Artifact{
		Kind:      por.KindGlobalProofBundle,
		SHA256:    sha,
		URL:       liabRef,
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
		case por.KindSummarySnapshot, por.KindBapiSnapshot:
			parsed.SummarySha = art.SHA256
			summaryPath = resolvePath(artifactsDir, art.LocalPath)
		case por.KindWalletZip, por.KindWalletAddressBundle:
			parsed.WalletZipSha = art.SHA256
			parsed.WalletZipPath = resolvePath(artifactsDir, art.LocalPath)
		}
	}
	if summaryPath == "" {
		return por.Snapshot{}, parsed, fmt.Errorf("artifact bundle missing summarySnapshot")
	}
	parsed.SummaryPath = summaryPath

	raw, err := os.ReadFile(summaryPath)
	if err != nil {
		return por.Snapshot{}, parsed, err
	}
	okxParsed, err := okxapi.ParseSummaryBytes(raw)
	if err != nil {
		return por.Snapshot{}, parsed, err
	}
	snap := Normalize(okxParsed.Audit, Meta{
		Exchange:  artBundle.Exchange,
		PeriodSeq: artBundle.PeriodSeq,
	})
	if artBundle.Exchange == "" {
		snap.Exchange = ID
	}
	return snap, parsed, nil
}

func importLocalFile(srcPath, outDir, ext string) (path, sumHex string, size int64, err error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return "", "", 0, err
	}
	defer f.Close()
	tmp, err := os.CreateTemp(outDir, "okx-import-*")
	if err != nil {
		return "", "", 0, err
	}
	tmpPath := tmp.Name()
	h := sha256.New()
	size, err = io.Copy(io.MultiWriter(tmp, h), f)
	tmp.Close()
	if err != nil {
		os.Remove(tmpPath)
		return "", "", 0, err
	}
	sumHex = hex.EncodeToString(h.Sum(nil))
	final := filepath.Join(outDir, sumHex+ext)
	if err := os.Rename(tmpPath, final); err != nil {
		os.Remove(tmpPath)
		return "", "", 0, err
	}
	return final, sumHex, size, nil
}

func resolvePath(artifactsDir, localPath string) string {
	if localPath == "" {
		return ""
	}
	if filepath.IsAbs(localPath) {
		return localPath
	}
	if p := filepath.Join(artifactsDir, localPath); fileExists(p) {
		return p
	}
	if p := filepath.Join(artifactsDir, "raw", filepath.Base(localPath)); fileExists(p) {
		return p
	}
	if fileExists(localPath) {
		return localPath
	}
	return filepath.Join(artifactsDir, filepath.Base(localPath))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isHTTP(s string) bool {
	return len(s) > 8 && (s[:7] == "http://" || s[:8] == "https://")
}

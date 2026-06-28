package gateio

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
	"github.com/ardmere/ardmere/internal/exchanges/gateio/gateapi"
	"github.com/ardmere/ardmere/internal/por"
	"github.com/ardmere/ardmere/internal/walletzip"
)

type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) ID() string { return ID }

func (a *Adapter) Capabilities() exchange.Capabilities {
	return exchange.Capabilities{
		Tier: exchange.TierSummary,
	}
}

func (a *Adapter) FetchAndStore(ctx context.Context, artifactsBase string, opts exchange.FetchOpts) (exchange.StoreResult, error) {
	return artifacts.FetchAndStore(ctx, a, artifactsBase, opts, artifacts.StoreOpts{})
}

func (a *Adapter) AggregateWalletZip(string) (*walletzip.Aggregate, error) {
	return nil, fmt.Errorf("gate.io public PoR has no wallet zip")
}

func (a *Adapter) VerifierProfile() exchange.VerifierProfile {
	return exchange.VerifierProfile{
		Shared: []string{
			"artifact-integrity@1",
			"solvency-claim@1",
		},
		Stubs: []string{
			"internal-consistency@0",
			"btc-anchor@0",
			"onchain-balance-hot@0",
			"onchain-balance-token@0",
			"address-ownership@0",
			"global-zk-proof@gateio-0",
			"third-party-attestation@0",
			"cross-chain-wrapped@0",
		},
	}
}

func (a *Adapter) FetchCurrentSnapshot(ctx context.Context, outDir string, opts exchange.FetchOpts) (exchange.FetchResult, error) {
	var summaryBytes []byte
	var parsed gateapi.ParsedSummary
	var err error

	if opts.SummaryPath != "" {
		summaryBytes, err = os.ReadFile(opts.SummaryPath)
		if err != nil {
			return exchange.FetchResult{}, fmt.Errorf("read summary: %w", err)
		}
		parsed, err = gateapi.ParseSummaryBytes(summaryBytes)
		if err != nil {
			return exchange.FetchResult{}, err
		}
	} else {
		summaryBytes, parsed, err = gateapi.FetchPublicSummary(ctx)
		if err != nil {
			return exchange.FetchResult{}, fmt.Errorf("fetch gate summary (Akamai may block datacenter IPs; use -summary-path): %w", err)
		}
	}

	sum := sha256.Sum256(summaryBytes)
	sumHex := hex.EncodeToString(sum[:])
	sumPath := filepath.Join(outDir, sumHex+".json")
	if err := os.WriteFile(sumPath, summaryBytes, 0o644); err != nil {
		return exchange.FetchResult{}, err
	}

	snap := Normalize(parsed, Meta{Exchange: ID})
	if snap.ID == "" {
		snap.ID = snap.SnapshotTime.Format("20060102")
	}

	arts := []bundle.Artifact{{
		Kind:      por.KindSummarySnapshot,
		SHA256:    sumHex,
		URL:       gateapi.SummaryURL(),
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

	if opts.ZkBundlePath != "" {
		art, err := importFileArtifact(opts.ZkBundlePath, outDir, por.KindGlobalProofBundle)
		if err != nil {
			return exchange.FetchResult{}, err
		}
		out.Artifacts = append(out.Artifacts, art)
	}
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
		case por.KindGlobalProofBundle:
			// reserved for future global-zk-proof verifier
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
	gateParsed, err := gateapi.ParseSummaryBytes(raw)
	if err != nil {
		return por.Snapshot{}, parsed, err
	}

	snap := Normalize(gateParsed, Meta{
		Exchange:  artBundle.Exchange,
		PeriodSeq: artBundle.PeriodSeq,
	})
	if artBundle.Exchange == "" {
		snap.Exchange = ID
	}
	return snap, parsed, nil
}

func importFileArtifact(srcPath, outDir, kind string) (bundle.Artifact, error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return bundle.Artifact{}, err
	}
	defer f.Close()

	tmp, err := os.CreateTemp(outDir, "import-*")
	if err != nil {
		return bundle.Artifact{}, err
	}
	tmpPath := tmp.Name()
	h := sha256.New()
	size, err := io.Copy(io.MultiWriter(tmp, h), f)
	tmp.Close()
	if err != nil {
		os.Remove(tmpPath)
		return bundle.Artifact{}, err
	}
	sumHex := hex.EncodeToString(h.Sum(nil))
	finalPath := filepath.Join(outDir, sumHex+filepath.Ext(srcPath))
	if err := os.Rename(tmpPath, finalPath); err != nil {
		os.Remove(tmpPath)
		return bundle.Artifact{}, err
	}
	return bundle.Artifact{
		Kind:      kind,
		SHA256:    sumHex,
		URL:       "file://" + srcPath,
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
		SizeBytes: size,
		LocalPath: finalPath,
	}, nil
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
	if _, err := os.Stat(localPath); err == nil {
		return localPath
	}
	return filepath.Join(artifactsDir, filepath.Base(localPath))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

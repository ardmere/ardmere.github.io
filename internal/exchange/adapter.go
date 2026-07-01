package exchange

import (
	"context"

	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/por"
	"github.com/ardmere/ardmere/internal/walletzip"
)

// Adapter fetches exchange artifacts and normalizes them into por.Snapshot.
type Adapter interface {
	ID() string

	Capabilities() Capabilities

	// FetchCurrentSnapshot downloads snapshot artifacts into outDir (flat temp dir).
	FetchCurrentSnapshot(ctx context.Context, outDir string, opts FetchOpts) (FetchResult, error)

	// FetchAndStore fetches into artifacts/<id>/<auditId>/raw/ and writes fetch.json.
	FetchAndStore(ctx context.Context, artifactsBase string, opts FetchOpts) (StoreResult, error)

	// AggregateWalletZip parses the exchange wallet bundle for consistency verifiers.
	AggregateWalletZip(path string) (*walletzip.Aggregate, error)

	ParseSnapshotFromArtifacts(artBundle bundle.ArtifactBundle, artifactsDir string) (por.Snapshot, ParsedArtifacts, error)

	VerifierProfile() VerifierProfile
}

// StoreResult is the outcome of FetchAndStore.
type StoreResult struct {
	Fetch       FetchResult
	SnapshotDir string
	RawDir      string
	BundlesDir  string
}

// FetchOpts controls optional fetch steps.
type FetchOpts struct {
	SkipWalletZip bool
	// SummaryPath imports a saved summary JSON instead of calling the public API.
	SummaryPath string
	// ZkBundlePath imports a local zkmerkle_cex_*.tar.gz (login-gated on Gate; manual drop-in).
	ZkBundlePath string
	// WalletZipPath imports a local wallet zip instead of downloading.
	WalletZipPath string
	// LiabilityZipPath imports or overrides the OKX liability zk-STARK zip URL/path.
	LiabilityZipPath string
	SkipLiabilityZip bool
	// ImportSource labels fetch.json when SummaryPath is supplied (e.g. browser-api).
	ImportSource string
	// UserProofPath imports a login-gated myProof.json (Bybit user Merkle path).
	UserProofPath string
	// AuditID fetches a specific snapshot when the exchange supports it (OKX, Binance).
	AuditID string
}

// FetchResult is the outcome of FetchCurrentSnapshot.
type FetchResult struct {
	Snapshot      por.Snapshot
	Artifacts     []bundle.Artifact
	WalletZipPath string
	WalletZipSha  string
	SummarySha    string
	SummaryPath   string
}

// ParsedArtifacts holds resolved local paths and content hashes from a bundle.
type ParsedArtifacts struct {
	SummaryPath   string
	SummarySha    string
	WalletZipPath string
	WalletZipSha  string
	UserProofPath string
	UserProofSha  string
}

// VerifierProfile lists verifiers to run for this exchange (P0: documentation + stubs hook).
type VerifierProfile struct {
	Shared []string
	Stubs  []string
}

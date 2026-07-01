package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ardmere/ardmere/internal/artifacts"
	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/exchange"
	"github.com/ardmere/ardmere/internal/exchangereg"
	"github.com/ardmere/ardmere/internal/por"
	"github.com/ardmere/ardmere/internal/verifier"
)

// LoadResult is parsed state from a saved artifact bundle.
type LoadResult struct {
	Adapter      exchange.Adapter
	Snapshot     por.Snapshot
	Artifacts    []bundle.Artifact
	ArtifactsDir string
	BundlesDir   string
	Parsed       exchange.ParsedArtifacts
	ArtBundle    bundle.ArtifactBundle
}

// LoadFromBundle reads an artifact bundle and resolves snapshot paths.
func LoadFromBundle(artifactsDir, snapshotID, exchangeID string) (LoadResult, error) {
	artBundlePath := artifacts.ResolveBundlePath(artifactsDir, snapshotID, ".artifact-bundle.json")
	snapshotRoot := artifacts.SnapshotRootFromDir(artBundleDir(artifactsDir))

	artBundleBytes, err := os.ReadFile(artBundlePath)
	if err != nil {
		return LoadResult{}, fmt.Errorf("read artifact bundle: %w", err)
	}
	var artBundle bundle.ArtifactBundle
	if err := json.Unmarshal(artBundleBytes, &artBundle); err != nil {
		return LoadResult{}, fmt.Errorf("decode artifact bundle: %w", err)
	}

	exID := exchangeID
	if exID == "" {
		exID = artBundle.Exchange
	}
	if exID == "" {
		exID = "binance"
	}
	adapter, err := exchangereg.Get(exID)
	if err != nil {
		return LoadResult{}, err
	}

	snap, parsed, err := adapter.ParseSnapshotFromArtifacts(artBundle, snapshotRoot)
	if err != nil {
		return LoadResult{}, err
	}

	bundlesDir := artifactsDir
	if filepath.Base(bundlesDir) != "bundles" {
		if _, err := os.Stat(filepath.Join(bundlesDir, "bundles")); err == nil {
			bundlesDir = filepath.Join(bundlesDir, "bundles")
		}
	}

	return LoadResult{
		Adapter:      adapter,
		Snapshot:     snap,
		Artifacts:    artBundle.Artifacts,
		ArtifactsDir: snapshotRoot,
		BundlesDir:   bundlesDir,
		Parsed:       parsed,
		ArtBundle:    artBundle,
	}, nil
}

func artBundleDir(artifactsDir string) string {
	return artifacts.SnapshotRootFromDir(artifactsDir)
}

// BTCBlockHeight returns the snapshot BTC anchor height when present.
func BTCBlockHeight(snap por.Snapshot) uint32 {
	if snap.TimeAnchor != nil && snap.TimeAnchor.Kind == "btc_block" {
		return snap.TimeAnchor.Height
	}
	return 0
}

// SummarizeVerifications compresses verifier outcomes for on-chain anchor.
func SummarizeVerifications(vs []verifier.Verification) (summary uint8, coverageBps uint16) {
	var maxCov float64
	for _, v := range vs {
		if v.Coverage > maxCov {
			maxCov = v.Coverage
		}
		switch v.VerifierID {
		case "internal-consistency":
			if v.Verdict == verifier.VerdictPass {
				summary |= 1 << 0
			}
		case "onchain-balance-hot":
			if v.Verdict == verifier.VerdictPass {
				summary |= 1 << 1
			}
		case "solvency-claim":
			if v.Verdict == verifier.VerdictPass {
				summary |= 1 << 2
			}
		}
	}
	if summary&(1<<0) != 0 {
		summary |= 1 << 3
	}
	return summary, uint16(maxCov * 10000)
}

// BuildBundles constructs artifact and verification bundle manifests.
func BuildBundles(snap por.Snapshot, arts []bundle.Artifact, verifications []verifier.Verification, verdictSummary uint8, coverageBps uint16) (bundle.ArtifactBundle, bundle.VerificationBundle) {
	snapshotTimeRFC := snap.SnapshotTime.UTC().Format(time.RFC3339)
	builtAt := time.Now().UTC().Format(time.RFC3339)

	artBundle := bundle.ArtifactBundle{
		Exchange:           snap.Exchange,
		SnapshotID:         snap.ID,
		PeriodSeq:          snap.PeriodSeq,
		SnapshotTime:       snapshotTimeRFC,
		BTCBlockHeight:     BTCBlockHeight(snap),
		ExchangeMerkleRoot: snap.MerkleRoot,
		BuiltAt:            builtAt,
		Kind:               "artifact-bundle",
		Artifacts:          arts,
		MerkleRoot:         bundle.HexRoot(bundle.ArtifactRoot(arts)),
	}
	verBundle := bundle.VerificationBundle{
		Exchange:           snap.Exchange,
		SnapshotID:         snap.ID,
		PeriodSeq:          snap.PeriodSeq,
		SnapshotTime:       snapshotTimeRFC,
		BTCBlockHeight:     BTCBlockHeight(snap),
		ExchangeMerkleRoot: snap.MerkleRoot,
		VerdictSummary:     verdictSummary,
		CoverageBps:        coverageBps,
		BuiltAt:            builtAt,
		Kind:               "verification-bundle",
		Verifications:      verifications,
		MerkleRoot:         bundle.HexRoot(bundle.VerificationRoot(verifications)),
	}
	return artBundle, verBundle
}

// WriteJSON writes indented JSON to path.
func WriteJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

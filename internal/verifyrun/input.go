package verifyrun

import (
	"os"
	"path/filepath"

	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/exchange"
	"github.com/ardmere/ardmere/internal/por"
	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/verifier"
	"github.com/ardmere/ardmere/internal/walletzip"
)

// Input holds everything needed to execute an exchange VerifierProfile.
type Input struct {
	Exchange      string
	Profile       exchange.VerifierProfile
	Snapshot      por.Snapshot
	Artifacts     []bundle.Artifact
	ArtifactsDir  string
	SummarySha    string
	WalletZipPath string
	WalletZipSha  string
	WalletAgg     *walletzip.Aggregate
	SkipRPC       bool
	RPC           *rpc.Client
	ETHDeposits        verifier.ETHDepositIndexer
	DepositMaxSamples  int
	DepositTopKPerCoin int
	DepositOnly        bool
	UserProofPath      string
	UserProofSha       string
}

func userProofFromArtifacts(arts []bundle.Artifact, artifactsDir string) (path, sha string) {
	for _, art := range arts {
		if art.Kind != "userMerkleProof" {
			continue
		}
		p := art.LocalPath
		if !filepath.IsAbs(p) {
			p = filepath.Join(artifactsDir, p)
		}
		if _, err := os.Stat(p); err == nil {
			return p, art.SHA256
		}
		if p2 := filepath.Join(artifactsDir, "raw", filepath.Base(art.LocalPath)); true {
			if _, err := os.Stat(p2); err == nil {
				return p2, art.SHA256
			}
		}
	}
	return "", ""
}

func proofBundleFromArtifacts(arts []bundle.Artifact, artifactsDir string) (path, sha string) {
	for _, art := range arts {
		if art.Kind != "globalProofBundle" {
			continue
		}
		p := art.LocalPath
		if !filepath.IsAbs(p) {
			p = filepath.Join(artifactsDir, p)
		}
		if _, err := os.Stat(p); err == nil {
			return p, art.SHA256
		}
		if p2 := filepath.Join(artifactsDir, "raw", filepath.Base(art.LocalPath)); true {
			if _, err := os.Stat(p2); err == nil {
				return p2, art.SHA256
			}
		}
	}
	return "", ""
}

func artifactRefs(arts []bundle.Artifact) []verifier.ArtifactRef {
	out := make([]verifier.ArtifactRef, len(arts))
	for i, a := range arts {
		out[i] = verifier.ArtifactRef{
			Kind:      a.Kind,
			SHA256:    a.SHA256,
			LocalPath: a.LocalPath,
		}
	}
	return out
}

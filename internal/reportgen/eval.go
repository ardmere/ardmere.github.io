package reportgen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/artifacts"
	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/verifier"
)

type stageEval struct {
	Stage         string
	StageName     string
	Gen           string
	EvidenceLevel string
	Confidence    string
	EffectivePoR  bool
	Headline      string
	Rationale     string
	Blocked       []blockedReason
	RiskFlags     []riskFlag
}

type blockedReason struct {
	RuleID    string `json:"ruleId"`
	Evidence  string `json:"evidence,omitempty"`
	MaxStage  string `json:"maxStage"`
	Reason    string `json:"reason"`
	RiskFlag  string `json:"riskFlag"`
}

type riskFlag struct {
	ID          string `json:"id"`
	Severity    string `json:"severity"`
	StageEffect string `json:"stageEffect"`
	Message     string `json:"message"`
}

func evaluate(exchangeID string, art bundle.ArtifactBundle, vers []verifier.Verification) stageEval {
	vmap := map[string]verifier.Verification{}
	for _, v := range vers {
		vmap[v.VerifierID] = v
	}
	has := func(kind string) bool {
		for _, a := range art.Artifacts {
			if a.Kind == kind || strings.Contains(a.Kind, kind) {
				return true
			}
		}
		return false
	}
	verdict := func(id string) string {
		if v, ok := vmap[id]; ok {
			return string(v.Verdict)
		}
		return "UNVERIFIABLE"
	}

	e := stageEval{
		Stage:         "Stage 0",
		StageName:     "Trust the Exchange",
		Gen:           "Gen 1",
		EvidenceLevel: "E1",
		Confidence:    "low",
		EffectivePoR:  false,
	}

	switch exchangeID {
	case "okx":
		e.Gen = "Gen 2"
		e.EvidenceLevel = "E2"
		e.Confidence = "medium"
		if has("walletZip") && has("globalProofBundle") &&
			verdict("address-ownership") == "PASS" &&
			(verdict("global-zk-proof") == "PASS" || verdict("global-zk-proof") == "PARTIAL") {
			e.Stage = "Stage 1"
			e.StageName = "Verifiable Disclosure"
			e.EffectivePoR = true
			e.Headline = "OKX publishes public summary, wallet_address_list, wallet_ownership_proof, and global zk proof bundle for audit " + art.SnapshotID + ". ardmere verifies address ownership and global zk binding, so this snapshot reaches Stage 1."
			e.Rationale = "Stage 1 is supported by public wallet_address_list, wallet_ownership_proof, global_proof, and parameter/proof artifacts. Stage 2 is not reached because canonical official anchoring, stable DA, low-friction user inclusion proof, stronger publication frequency, and full business-consistent constraints are not established in this report."
		} else {
			e.Headline = "OKX audit " + art.SnapshotID + " does not meet all Stage 1 public artifact thresholds in this pipeline."
			e.Rationale = "Missing or failing verifiers block Stage 1 effective PoR for this snapshot."
		}
		e.Blocked = []blockedReason{
			{RuleID: "floor_server_only_distribution", Evidence: "canonical on-chain/DA anchor", MaxStage: "Stage 1", Reason: "No official exchange canonical on-chain/DA anchor established.", RiskFlag: "NO_CANONICAL_ANCHOR"},
			{RuleID: "s2_frequency", Evidence: "weekly full PoR / daily anchor", MaxStage: "Stage 1", Reason: "Monthly cadence does not satisfy Stage 2 frequency expectations.", RiskFlag: "HIGH_FREQUENCY_GAP"},
		}
	case "binance":
		e.Gen = "Gen 2"
		e.EvidenceLevel = "E2"
		e.Confidence = "medium"
		e.Headline = "Binance snapshot " + art.SnapshotID + " provides public reserve summary and wallet_address_list where available, but Stage 1 is blocked by missing wallet_ownership_proof, opaque trusted setup, and unavailable public global proof/vk artifacts."
		e.Rationale = "The available artifacts support Gen 2 / E2 classification and Stage 0 PoR disclosure. Users still need to trust Binance for wallet control, trusted setup honesty, and public availability of the full global proof stack."
		if verdict("global-zk-proof") == "UNVERIFIABLE" || !has("globalProofBundle") {
			e.Blocked = append(e.Blocked,
				blockedReason{RuleID: "floor_no_wallet_ownership_proof", Evidence: "wallet_ownership_proof", MaxStage: "Stage 0", Reason: "No public batch-verifiable wallet_ownership_proof.", RiskFlag: "NO_WALLET_OWNERSHIP_PROOF"},
				blockedReason{RuleID: "floor_login_wall_global_proof", Evidence: "global_proof.csv, verifying_key", MaxStage: "Stage 0", Reason: "Public global proof/vk not available.", RiskFlag: "UNVERIFIABLE"},
				blockedReason{RuleID: "floor_opaque_trusted_setup", Evidence: "trusted_setup_transcript", MaxStage: "Stage 0", Reason: "Trusted setup transcript is not public.", RiskFlag: "OPAQUE_TRUSTED_SETUP"},
			)
		}
	default:
		e.Headline = exchangeID + " snapshot " + art.SnapshotID + " remains Stage 0 in the ardmere public evaluation set."
		e.Rationale = "Public artifacts do not support independently reproducible global solvency verification for this snapshot."
		if !has("walletZip") {
			e.Blocked = append(e.Blocked, blockedReason{RuleID: "floor_no_public_wallet_address_list", Evidence: "wallet_address_list", MaxStage: "Stage 0", Reason: "No public wallet_address_list.", RiskFlag: "UNVERIFIABLE"})
		}
	}

	for _, b := range e.Blocked {
		e.RiskFlags = append(e.RiskFlags, riskFlag{
			ID:          b.RiskFlag,
			Severity:    "high",
			StageEffect: "max_stage_0",
			Message:     b.Reason,
		})
	}
	if e.Stage == "Stage 1" {
		for i := range e.RiskFlags {
			if e.RiskFlags[i].ID == "NO_CANONICAL_ANCHOR" || e.RiskFlags[i].ID == "HIGH_FREQUENCY_GAP" {
				e.RiskFlags[i].StageEffect = "max_stage_1"
				e.RiskFlags[i].Severity = "medium"
			}
		}
	}
	return e
}

func loadVerificationBundle(dir, snapshotID string) (bundle.VerificationBundle, error) {
	path := artifacts.ResolveBundlePath(dir, snapshotID, ".verification-bundle.v2.json")
	if _, err := os.Stat(path); err != nil {
		path = artifacts.ResolveBundlePath(dir, snapshotID, ".verification-bundle.json")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return bundle.VerificationBundle{}, err
	}
	var vb bundle.VerificationBundle
	if err := json.Unmarshal(raw, &vb); err != nil {
		return bundle.VerificationBundle{}, err
	}
	return vb, nil
}

func loadArtifactBundle(dir, snapshotID string) (bundle.ArtifactBundle, error) {
	path := artifacts.ResolveBundlePath(dir, snapshotID, ".artifact-bundle.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return bundle.ArtifactBundle{}, err
	}
	var ab bundle.ArtifactBundle
	if err := json.Unmarshal(raw, &ab); err != nil {
		return bundle.ArtifactBundle{}, err
	}
	return ab, nil
}

func snapshotDir(artifactsBase, exchangeID, snapshotID string) string {
	return filepath.Join(artifacts.ExchangeBase(artifactsBase, exchangeID), snapshotID)
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}

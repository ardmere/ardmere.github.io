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
	RuleID   string `json:"ruleId"`
	Evidence string `json:"evidence,omitempty"`
	MaxStage string `json:"maxStage"`
	Reason   string `json:"reason"`
	RiskFlag string `json:"riskFlag"`
}

type riskFlag struct {
	ID          string `json:"id"`
	Severity    string `json:"severity"`
	StageEffect string `json:"stageEffect"`
	Message     string `json:"message"`
}

// EvaluateStage runs the same stage determination used by transparency reports.
func EvaluateStage(exchangeID string, art bundle.ArtifactBundle, vers []verifier.Verification) (stage string, effectivePoR bool) {
	e := evaluate(exchangeID, art, vers)
	return e.Stage, e.EffectivePoR
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
	case "mexc":
		e.Headline = "MEXC snapshot " + art.SnapshotID + " is a weak Stage 0 (press-release / import summary): self-reported ratios only; no public period merkle root in the archived package; user Merkle remains login-gated."
		e.Rationale = "solvency-claim PASS means self-reported reserve ratios ≥ 100% only — not independently reproduced reserves. Datacenter fetches of www.mexc.com are WAF-blocked; Hacken PDF attestation is not machine-parsed here."
		e.Blocked = append(e.Blocked,
			blockedReason{RuleID: "floor_no_public_wallet_address_list", Evidence: "wallet_address_list", MaxStage: "Stage 0", Reason: "No public wallet_address_list.", RiskFlag: "UNVERIFIABLE"},
			blockedReason{RuleID: "floor_no_public_merkle_root", Evidence: "merkle_root", MaxStage: "Stage 0", Reason: "Press-release summary has no public period merkle root to bind user proofs.", RiskFlag: "NO_PUBLIC_MERKLE_ROOT"},
			blockedReason{RuleID: "floor_waf_gated_dashboard", Evidence: "por_dashboard_api", MaxStage: "Stage 0", Reason: "www.mexc.com PoR JSON is WAF-gated from datacenters; summary imported from press materials.", RiskFlag: "WAF_GATED"},
		)
		if verdict("user-merkle-proof") == "UNVERIFIABLE" {
			e.Blocked = append(e.Blocked, blockedReason{RuleID: "floor_login_gated_user_proof", Evidence: "userMerkleProof", MaxStage: "Stage 0", Reason: "User myProof.json not present in the public evaluation package (login-gated).", RiskFlag: "LOGIN_GATED_USER_PROOF"})
		}
	case "ourbit":
		e.Headline = "Ourbit snapshot " + art.SnapshotID + " is Stage 0: public stock_info ratios + live multi-chain wallet table; no snapshot-bound walletZip/ownership; user Merkle login-gated; hash truncated to 128 bits."
		e.Rationale = "solvency-claim PASS means self-reported stock_info reserve ratios ≥ 100% only. wallet_detail is archived for transparency but updates live and lacks ownership signatures, so it is not Stage 1 wallet_address_list."
		e.Blocked = append(e.Blocked,
			blockedReason{RuleID: "floor_no_snapshot_wallet_address_list", Evidence: "wallet_address_list", MaxStage: "Stage 0", Reason: "Public wallet_detail is live/unsigned — not a snapshot-bound wallet_address_list.", RiskFlag: "UNVERIFIABLE"},
			blockedReason{RuleID: "floor_no_public_merkle_root", Evidence: "merkle_root", MaxStage: "Stage 0", Reason: "stock_info has no public period merkle root to bind user proofs.", RiskFlag: "NO_PUBLIC_MERKLE_ROOT"},
			blockedReason{RuleID: "floor_truncated_hash", Evidence: "sha256_truncation_128", MaxStage: "Stage 0", Reason: "Ourbit Merkle hashes truncate SHA-256 to 32 hex chars (128-bit).", RiskFlag: "TRUNCATED_HASH"},
		)
		if verdict("user-merkle-proof") == "UNVERIFIABLE" {
			e.Blocked = append(e.Blocked, blockedReason{RuleID: "floor_login_gated_user_proof", Evidence: "userMerkleProof", MaxStage: "Stage 0", Reason: "User merkel_tree_ourbit.json not present in the public evaluation package (login-gated).", RiskFlag: "LOGIN_GATED_USER_PROOF"})
		}
	case "coinex":
		e.Headline = "CoinEx snapshot " + art.SnapshotID + " is Stage 0: public reserve ratios + period merkle root + XLSX address list; user Merkle login-gated; no ownership proofs or global ZK."
		e.Rationale = "solvency-claim PASS means self-reported reserve_rate ≥ 100% only. Public wallet_proof_file is an XLSX address dump without ownership signatures, so it is not Stage 1 wallet_address_list/walletZip. Official Go OSS (coinexcom/proof-of-reserves) is ported as user-merkle-proof@coinex-1 for imported proofs."
		e.Blocked = append(e.Blocked,
			blockedReason{RuleID: "floor_no_snapshot_wallet_address_list", Evidence: "wallet_address_list", MaxStage: "Stage 0", Reason: "Public wallet_proof_file is XLSX without ownership signatures — not Stage 1 walletZip.", RiskFlag: "UNVERIFIABLE"},
		)
		if verdict("user-merkle-proof") == "UNVERIFIABLE" {
			e.Blocked = append(e.Blocked, blockedReason{RuleID: "floor_login_gated_user_proof", Evidence: "userMerkleProof", MaxStage: "Stage 0", Reason: "User merkle proof JSON not present in the public evaluation package (login-gated).", RiskFlag: "LOGIN_GATED_USER_PROOF"})
		}

	case "coinw":
		e.Headline = "CoinW snapshot " + art.SnapshotID + " is a weak Stage 0 (announcement / import summary): self-reported total USDT reserve ratio only; Cloudflare WAF blocks stable datacenter fetches of the PoR dashboard; no public wallet list."
		e.Rationale = "solvency-claim PASS means the September 2026 announcement total reserve ratio ≥ 100% only — not independently reproduced reserves or a bound merkle root. Official por-verify binaries exist without published source; user merkel_tree_cw.json remains login-gated. A later Oct-24 SSR panel was observed under geo-restriction with animated/conflicting merkle and ratio UI — left observational, not assessed."
		e.Blocked = append(e.Blocked,
			blockedReason{RuleID: "floor_no_public_wallet_address_list", Evidence: "wallet_address_list", MaxStage: "Stage 0", Reason: "CoinW publishes marketing PoR ratios without a permissionless wallet_address_list.", RiskFlag: "UNVERIFIABLE"},
			blockedReason{RuleID: "floor_no_public_merkle_root", Evidence: "merkle_root", MaxStage: "Stage 0", Reason: "Assessed announcement summary has no public period merkle root to bind user proofs; dashboard merkle UI was not machine-stable under WAF/geo probe.", RiskFlag: "NO_PUBLIC_MERKLE_ROOT"},
			blockedReason{RuleID: "floor_waf_gated_dashboard", Evidence: "por_dashboard", MaxStage: "Stage 0", Reason: "www.coinw.com PoR pages return Cloudflare WAF 403 from datacenters; summary imported from announcement / browser probe.", RiskFlag: "WAF_GATED"},
		)
		if verdict("user-merkle-proof") == "UNVERIFIABLE" {
			e.Blocked = append(e.Blocked, blockedReason{RuleID: "floor_login_gated_user_proof", Evidence: "userMerkleProof", MaxStage: "Stage 0", Reason: "User merkel_tree_cw.json + automatic verify require login; OSS release ships binaries only (no source to port).", RiskFlag: "LOGIN_GATED_USER_PROOF"})
		}

	case "kraken":
		e.Headline = "Kraken snapshot " + art.SnapshotID + " is Stage 0 attestation-model: quarterly accountant Merkle (LedgerLens / The Network Firm) plus CMS reserve ratios — no permissionless wallet ZIP or machine-verified attestation PDF."
		e.Rationale = "solvency-claim PASS means self-reported CMS reserve ratios ≥ 100% (capped stables treated as published 105%+ floor). LedgerLens exposes a full combined Merkle root and names The Network Firm LLP, but user leaf/path remains login-gated, wallet ownership evidence stays inside the accountant workflow, and third-party-attestation is not machine-parsed — Stage 1 stays out of reach by disclosure design."
		e.Blocked = append(e.Blocked,
			blockedReason{RuleID: "floor_no_public_wallet_address_list", Evidence: "wallet_address_list", MaxStage: "Stage 0", Reason: "Kraken publishes ratios + accountant Merkle root only — no permissionless wallet addresses or reserve/liability amounts.", RiskFlag: "UNVERIFIABLE"},
			blockedReason{RuleID: "floor_attestation_not_machine_verified", Evidence: "third_party_attestation", MaxStage: "Stage 0", Reason: "Accountant is named (The Network Firm LLP / LedgerLens) but the CPA attestation report is not ingested as a machine-verifiable public artifact.", RiskFlag: "UNVERIFIABLE"},
			blockedReason{RuleID: "floor_quarterly_cadence", Evidence: "publication_frequency", MaxStage: "Stage 0", Reason: "Quarterly accountant reviews are too infrequent for Stage 2 frequency expectations and are a weak monitor-primary path.", RiskFlag: "HIGH_FREQUENCY_GAP"},
		)
		if verdict("user-merkle-proof") == "UNVERIFIABLE" {
			e.Blocked = append(e.Blocked, blockedReason{RuleID: "floor_login_gated_user_proof", Evidence: "userMerkleProof", MaxStage: "Stage 0", Reason: "User Merkle Leaf ID / path self-check requires Kraken Pro login (or personalized accountant tool input); no public user proof in the evaluation package.", RiskFlag: "LOGIN_GATED_USER_PROOF"})
		}

	case "btcc":
		e.Headline = "BTCC snapshot " + art.SnapshotID + " is Stage 0 hard-capped: permissionless reserveInfo exposes per-coin ratios and a truncated (~16 hex) merkle root only — no wallet amounts/addresses, no OSS verifier."
		e.Rationale = "solvency-claim PASS means self-reported reserve ratios ≥ 100% from https://www.btcc.com/social-api/contract/reserve/reserveInfo only. Cadence is monthly ~day-15. Truncated root cannot bind a full Merkle path offline; leaf validate and history APIs need login; Stage 1 walletZip/ownership/ZK stay unwired by design of the public disclosure."
		e.Blocked = append(e.Blocked,
			blockedReason{RuleID: "floor_no_public_wallet_address_list", Evidence: "wallet_address_list", MaxStage: "Stage 0", Reason: "BTCC publishes ratios only — no wallet addresses or reserve/liability amounts.", RiskFlag: "UNVERIFIABLE"},
			blockedReason{RuleID: "floor_truncated_merkle_root", Evidence: "merkle_root", MaxStage: "Stage 0", Reason: "Published merkle root is truncated to ~16 hex characters; insufficient to bind a full Merkle proof offline.", RiskFlag: "TRUNCATED_MERKLE_ROOT"},
			blockedReason{RuleID: "floor_no_oss_user_merkle_verifier", Evidence: "user_merkle_verifier", MaxStage: "Stage 0", Reason: "No official BTCC OSS Merkle verifier repo found; user-merkle-proof stays stubbed.", RiskFlag: "UNVERIFIABLE"},
		)
		if verdict("user-merkle-proof") == "UNVERIFIABLE" {
			e.Blocked = append(e.Blocked, blockedReason{RuleID: "floor_login_gated_user_proof", Evidence: "userMerkleProof", MaxStage: "Stage 0", Reason: "Leaf hash self-check and /reserve/validate require login; no public user proof artifact in the evaluation package.", RiskFlag: "LOGIN_GATED_USER_PROOF"})
		}

	case "phemex":
		e.Headline = "Phemex snapshot " + art.SnapshotID + " is Stage 0: permissionless porVersions balances → computed reserve ratios; monthly Merkle inclusion and cold-wallet addresses appear on the public PoR UI, but merkle root / wallet-addresses APIs are CloudFront-gated and user proofs stay login-gated."
		e.Rationale = "solvency-claim PASS means chainBalance/cexBalance ratios from https://por.phemex.com/porVersions only. UI shows Part of Our Cold Wallets and a Merkle self-check, but there is no permissionless wallet zip, no public period merkle root in the archived package, and no official OSS verifier — Stage 1 walletZip/ownership/ZK stay unwired."
		e.Blocked = append(e.Blocked,
			blockedReason{RuleID: "floor_no_machine_readable_wallet_list", Evidence: "wallet_address_list", MaxStage: "Stage 0", Reason: "Cold-wallet addresses are UI-visible only; api10 wallet-addresses is CloudFront/auth gated and there is no downloadable wallet zip.", RiskFlag: "UNVERIFIABLE"},
			blockedReason{RuleID: "floor_no_public_merkle_root", Evidence: "merkle_root", MaxStage: "Stage 0", Reason: "Period merkle root is not in public porVersions; api10 merkle-tree/root is CloudFront/auth gated from datacenters.", RiskFlag: "NO_PUBLIC_MERKLE_ROOT"},
			blockedReason{RuleID: "floor_no_oss_user_merkle_verifier", Evidence: "user_merkle_verifier", MaxStage: "Stage 0", Reason: "No official Phemex OSS Merkle verifier repo found; user-merkle-proof stays stubbed.", RiskFlag: "UNVERIFIABLE"},
		)
		if verdict("user-merkle-proof") == "UNVERIFIABLE" {
			e.Blocked = append(e.Blocked, blockedReason{RuleID: "floor_login_gated_user_proof", Evidence: "userMerkleProof", MaxStage: "Stage 0", Reason: "Hashed-client-ID Merkle inclusion tool requires login; no public user proof artifact in the evaluation package.", RiskFlag: "LOGIN_GATED_USER_PROOF"})
		}

	case "toobit":
		e.Headline = "Toobit snapshot " + art.SnapshotID + " is Stage 0: public BTC/ETH/USDT/USDC reserve ratios + period merkle root via latest-audit; user SMT proof JSON login-gated; no ownership proofs or global ZK."
		e.Rationale = "solvency-claim PASS means self-reported reserve ratios ≥ 100% only. Official Go OSS (toobit-docs/por-verifier) is ported as user-merkle-proof@toobit-1 for imported proofs. No public wallet address list — Stage 1 walletZip stays unwired."
		e.Blocked = append(e.Blocked,
			blockedReason{RuleID: "floor_no_public_wallet_address_list", Evidence: "wallet_address_list", MaxStage: "Stage 0", Reason: "No public wallet address list observed for Toobit PoR.", RiskFlag: "UNVERIFIABLE"},
		)
		if verdict("user-merkle-proof") == "UNVERIFIABLE" {
			e.Blocked = append(e.Blocked, blockedReason{RuleID: "floor_login_gated_user_proof", Evidence: "userMerkleProof", MaxStage: "Stage 0", Reason: "User proof_of_reserve_{auditId}.json not present in the public evaluation package (login-gated).", RiskFlag: "LOGIN_GATED_USER_PROOF"})
		}

	case "woox":
		e.Headline = "WOO X snapshot " + art.SnapshotID + " is Stage 0: public live dashboard + monthly merkle audit overview with permissionless full-tree download; user Merkle hash login-gated; assets largely on custody/CEX venues — Stage 1 ownership/on-chain path hard."
		e.Rationale = "solvency-claim reflects self-reported USD-notional assets/liabilities from snapshot_overview (UI formula). Full protobuf Merkle trees are public on oss.woo.network (~200MB) and woo-x/merkle-tree documents validation, but ardmere does not yet port the Java checker and user hashes remain login-gated. No public wallet_address_list / ownership / global ZK."
		e.Blocked = append(e.Blocked,
			blockedReason{RuleID: "floor_no_public_wallet_address_list", Evidence: "wallet_address_list", MaxStage: "Stage 0", Reason: "No public wallet_address_list; live assets cite custody/CEX venues via read-only API.", RiskFlag: "UNVERIFIABLE"},
			blockedReason{RuleID: "floor_no_public_merkle_root_hex", Evidence: "merkle_root", MaxStage: "Stage 0", Reason: "Period merkle root hex is not published separately; only the full protobuf tree URL is public.", RiskFlag: "NO_PUBLIC_MERKLE_ROOT"},
			blockedReason{RuleID: "floor_assets_on_third_party_venues", Evidence: "asset_venues", MaxStage: "Stage 0", Reason: "Dashboard assets include custodian / third-party venue balances — Stage 1 on-chain ownership is hard.", RiskFlag: "THIRD_PARTY_VENUE_ASSETS"},
		)
		if verdict("user-merkle-proof") == "UNVERIFIABLE" {
			e.Blocked = append(e.Blocked, blockedReason{RuleID: "floor_login_gated_user_proof", Evidence: "userMerkleProof", MaxStage: "Stage 0", Reason: "User Merkle hash is login-gated; woo-x/merkle-tree Java not yet ported against the ~200MB tree.", RiskFlag: "LOGIN_GATED_USER_PROOF"})
		}

	case "bitunix":
		e.Headline = "Bitunix snapshot " + art.SnapshotID + " is Stage 0: public BTC/ETH/USDT reserve ratios + period merkle root; user merkle_tree_bitunix.json login-gated; no ownership proofs or global ZK."
		e.Rationale = "solvency-claim PASS means self-reported reserve ratios ≥ 100% only. Official Java OSS (BitunixOfficial/proof-of-reserves) is ported as user-merkle-proof@bitunix-1 for imported proofs. Recent periods expose empty walletAddress — no Stage 1 walletZip."
		e.Blocked = append(e.Blocked,
			blockedReason{RuleID: "floor_no_public_wallet_address_list", Evidence: "wallet_address_list", MaxStage: "Stage 0", Reason: "Public walletAddress field is empty for recent audits; no Stage 1 walletZip.", RiskFlag: "UNVERIFIABLE"},
		)
		if verdict("user-merkle-proof") == "UNVERIFIABLE" {
			e.Blocked = append(e.Blocked, blockedReason{RuleID: "floor_login_gated_user_proof", Evidence: "userMerkleProof", MaxStage: "Stage 0", Reason: "User merkle_tree_bitunix.json not present in the public evaluation package (login-gated).", RiskFlag: "LOGIN_GATED_USER_PROOF"})
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

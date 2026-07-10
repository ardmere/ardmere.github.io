package reportgen

import (
	"strings"

	"github.com/ardmere/ardmere/internal/bundle"
	"github.com/ardmere/ardmere/internal/por"
)

type reportContext struct {
	Exchange string
	Art      bundle.ArtifactBundle
	Ver      bundle.VerificationBundle
	Eval     stageEval
	Freq     frequencyInfo
}

type checklistItem struct {
	ID       string
	Stage    string
	Question string
	Status   string
	Notes    string
}

type recommendationItem struct {
	Priority         string
	Text             string
	RelatedRiskFlags []string
}

func newReportContext(exchangeID string, art bundle.ArtifactBundle, ver bundle.VerificationBundle, eval stageEval, freq frequencyInfo) reportContext {
	return reportContext{
		Exchange: exchangeID,
		Art:      art,
		Ver:      ver,
		Eval:     eval,
		Freq:     freq,
	}
}

func (c reportContext) hasArtifact(kind string) bool {
	for _, a := range c.Art.Artifacts {
		if a.Kind == kind || strings.Contains(a.Kind, kind) {
			return true
		}
	}
	return false
}

func (c reportContext) hasSummary() bool {
	for _, a := range c.Art.Artifacts {
		if por.IsSummaryKind(a.Kind) {
			return true
		}
	}
	return false
}

func (c reportContext) hasWalletList() bool {
	for _, a := range c.Art.Artifacts {
		if por.IsWalletKind(a.Kind) {
			return true
		}
	}
	return false
}

func (c reportContext) verdict(id string) string {
	for _, v := range c.Ver.Verifications {
		if v.VerifierID == id {
			return strings.ToUpper(string(v.Verdict))
		}
	}
	return "UNVERIFIABLE"
}

func (c reportContext) hasRiskFlag(id string) bool {
	for _, rf := range c.Eval.RiskFlags {
		if rf.ID == id {
			return true
		}
	}
	return false
}

func collectChecklist(c reportContext) []checklistItem {
	items := []checklistItem{
		{
			ID:       "S0-1",
			Stage:    "Stage 0",
			Question: "Has public PoR been published during the last 12 months?",
			Status:   "pass",
			Notes:    "Snapshot " + c.Art.SnapshotID + " is in the public evaluation set.",
		},
		{
			ID:       "S0-2",
			Stage:    "Stage 0",
			Question: "Is snapshot time explicit?",
			Status:   checklistStatus(c.Art.SnapshotTime != "", "unverifiable"),
			Notes:    c.Art.SnapshotTime,
		},
		{
			ID:       "S0-3",
			Stage:    "Stage 0",
			Question: "Is there reserve summary or wallet_address_list?",
			Status:   checklistStatus(c.hasSummary() || c.hasWalletList(), "fail"),
			Notes:    artifactNote(c.hasSummary(), c.hasWalletList()),
		},
		{
			ID:       "S1-1",
			Stage:    "Stage 1",
			Question: "Is wallet_address_list public and on-chain verifiable?",
			Status:   checklistStatus(c.hasWalletList(), "fail"),
			Notes:    walletListNote(c.hasWalletList(), c.verdict("onchain-balance-hot")),
		},
		{
			ID:       "S1-2",
			Stage:    "Stage 1",
			Question: "Is wallet_ownership_proof public and batch-verifiable?",
			Status:   verdictChecklistStatus(c.verdict("address-ownership")),
			Notes:    verifierNote(c.verdict("address-ownership"), "address-ownership"),
		},
		{
			ID:       "S1-3",
			Stage:    "Stage 1",
			Question: "Is global_proof public and independently reproducible?",
			Status:   globalProofChecklistStatus(c.verdict("global-zk-proof")),
			Notes:    verifierNote(c.verdict("global-zk-proof"), "global-zk-proof"),
		},
		{
			ID:       "S1-4",
			Stage:    "Stage 1",
			Question: "Is independent third-party review available with stated scope?",
			Status:   verdictChecklistStatus(c.verdict("third-party-attestation")),
			Notes:    verifierNote(c.verdict("third-party-attestation"), "third-party-attestation"),
		},
		{
			ID:       "S1-5",
			Stage:    "Stage 1",
			Question: "If trusted setup is required, is transcript public; otherwise is the proof system transparent setup?",
			Status:   trustedSetupStatus(c),
			Notes:    trustedSetupNote(c),
		},
	}

	if c.Eval.EffectivePoR {
		items = append(items,
			checklistItem{
				ID:       "S2-1",
				Stage:    "Stage 2",
				Question: "Is root/proof/vk canonically anchored outside exchange servers?",
				Status:   checklistStatus(!c.hasRiskFlag("NO_CANONICAL_ANCHOR"), "fail"),
				Notes:    "No official exchange canonical on-chain/DA anchor established.",
			},
			checklistItem{
				ID:       "S2-2",
				Stage:    "Stage 2",
				Question: "Is publication frequency sufficient for Stage 2 (weekly full PoR / daily anchor)?",
				Status:   frequencyStage2Status(c.Freq.ObservedCadence),
				Notes:    c.Freq.StageImpact,
			},
		)
	} else {
		items = append(items,
			checklistItem{
				ID:       "S2-1",
				Stage:    "Stage 2",
				Question: "Is root/proof/vk canonically anchored outside exchange servers?",
				Status:   "not_applicable",
				Notes:    "Stage 2 is not evaluated until Stage 1 blockers are resolved.",
			},
			checklistItem{
				ID:       "S2-2",
				Stage:    "Stage 2",
				Question: "Is publication frequency sufficient for Stage 2 (weekly full PoR / daily anchor)?",
				Status:   "not_applicable",
				Notes:    "Stage 2 is not evaluated until Stage 1 blockers are resolved.",
			},
		)
	}
	return items
}

func collectRecommendations(c reportContext) []recommendationItem {
	seen := map[string]bool{}
	var out []recommendationItem
	add := func(priority, text string, flags ...string) {
		if text == "" || seen[text] {
			return
		}
		seen[text] = true
		out = append(out, recommendationItem{
			Priority:         priority,
			Text:             text,
			RelatedRiskFlags: flags,
		})
	}

	for _, rf := range c.Eval.RiskFlags {
		if rec, ok := recommendationForRiskFlag(rf.ID, rf.Message); ok {
			add(rec.Priority, rec.Text, rf.ID)
		}
	}

	if exchangeProofSetup(c.Exchange) == proofSetupTrusted && trustedSetupStatus(c) == "unverifiable" {
		add("P0", "Publish trusted setup transcript or migrate to a transparent-setup proof system.", "OPAQUE_TRUSTED_SETUP")
	}

	if c.hasOnchainFail() {
		add("P2", "Review onchain-balance verifier FAIL/WARN outliers in the verification bundle and reconcile CSV vs chain deltas.", "IMPLEMENTATION_SEMANTICS_RISK")
	}
	if len(out) == 0 {
		add("P2", "Maintain public artifact archives and reproducible verifier output for each snapshot.", "UNVERIFIABLE")
	}
	return out
}

func (c reportContext) hasOnchainFail() bool {
	for _, id := range []string{"onchain-balance-hot", "onchain-balance-token", "onchain-balance-ledger", "onchain-balance-deposit"} {
		v := c.verdict(id)
		if v == "FAIL" || v == "WARN" {
			return true
		}
	}
	return false
}

type riskRecommendation struct {
	Priority string
	Text     string
}

func recommendationForRiskFlag(id, message string) (riskRecommendation, bool) {
	switch id {
	case "NO_WALLET_OWNERSHIP_PROOF":
		return riskRecommendation{"P0", "Publish batch-verifiable wallet_ownership_proof for wallet_address_list."}, true
	case "OPAQUE_TRUSTED_SETUP":
		return riskRecommendation{"P0", "Publish trusted setup transcript or migrate to a transparent-setup proof system."}, true
	case "LOGIN_WALL", "UNVERIFIABLE":
		if strings.Contains(strings.ToLower(message), "global proof") || strings.Contains(strings.ToLower(message), "proof/vk") {
			return riskRecommendation{"P0", "Make global proof.csv, verifying key, and parameter sources publicly downloadable."}, true
		}
		return riskRecommendation{"P0", "Publish missing public artifacts required for independent reproduction."}, true
	case "NO_CANONICAL_ANCHOR":
		return riskRecommendation{"P1", "Anchor root/proof/commitment on-chain or publish immutable DA/archive records."}, true
	case "HIGH_FREQUENCY_GAP":
		return riskRecommendation{"P1", "Add daily anchoring, weekly full PoR, or event-triggered updates to reduce timing risk."}, true
	case "NO_INDEPENDENT_REVIEW":
		return riskRecommendation{"P1", "Publish scoped third-party review with stated assurance type and scope."}, true
	case "HIGH_FRICTION_INCLUSION_PROOF":
		return riskRecommendation{"P1", "Provide Web/WASM/GUI one-click inclusion verification with exportable proof."}, true
	case "SAMPLE_ONLY":
		return riskRecommendation{"P0", "Publish the current production proof bundle and bind it to the snapshot summary."}, true
	case "BLACK_BOX_INCLUSION_PROOF":
		return riskRecommendation{"P0", "Provide exportable user inclusion proof and a local/open verifier."}, true
	default:
		return riskRecommendation{}, false
	}
}

func checklistStatus(ok bool, failStatus string) string {
	if ok {
		return "pass"
	}
	if failStatus == "" {
		failStatus = "fail"
	}
	return failStatus
}

func verdictChecklistStatus(verdict string) string {
	switch verdict {
	case "PASS":
		return "pass"
	case "PARTIAL", "WARN":
		return "warn"
	case "FAIL":
		return "fail"
	default:
		return "unverifiable"
	}
}

func globalProofChecklistStatus(verdict string) string {
	switch verdict {
	case "PASS", "PARTIAL":
		return "pass"
	case "WARN":
		return "warn"
	case "FAIL":
		return "fail"
	default:
		return "unverifiable"
	}
}

func trustedSetupStatus(c reportContext) string {
	switch exchangeProofSetup(c.Exchange) {
	case proofSetupTransparent:
		return "pass"
	case proofSetupTrusted:
		if c.verdict("global-zk-proof") == "PASS" {
			return "warn"
		}
		return "unverifiable"
	default:
		return "not_applicable"
	}
}

func trustedSetupNote(c reportContext) string {
	switch exchangeProofSetup(c.Exchange) {
	case proofSetupTransparent:
		return "Plonky2 / zk-STARK uses transparent setup; no trusted setup ceremony required."
	case proofSetupTrusted:
		if c.Exchange == "binance" {
			return "No public trusted setup transcript observed for zk-SNARK PoS."
		}
		return "No public trusted setup transcript observed for zk-SNARK global proof."
	default:
		return "No global ZK proof in the assessed public artifact set."
	}
}

func frequencyStage2Status(cadence string) string {
	switch cadence {
	case "~weekly", "sub-daily":
		return "warn"
	case "~monthly", "UNVERIFIABLE":
		return "fail"
	default:
		if strings.Contains(cadence, "day") {
			return "fail"
		}
		return "fail"
	}
}

func artifactNote(hasSummary, hasWallet bool) string {
	switch {
	case hasSummary && hasWallet:
		return "Summary and wallet_address_list present."
	case hasSummary:
		return "Summary present; wallet_address_list absent."
	case hasWallet:
		return "wallet_address_list present; summary absent or stub."
	default:
		return "Neither summary nor wallet list observed."
	}
}

func walletListNote(hasWallet bool, hotVerdict string) string {
	if !hasWallet {
		return "No public wallet_address_list artifact."
	}
	if hotVerdict == "UNVERIFIABLE" {
		return "Wallet list present; on-chain replay skipped or incomplete."
	}
	return "Wallet list present; hot-wallet on-chain verifier ran."
}

func verifierNote(verdict, id string) string {
	if verdict == "PASS" {
		return id + " verifier PASS."
	}
	if verdict == "PARTIAL" {
		return id + " verifier PARTIAL."
	}
	return id + " verifier " + verdict + "."
}

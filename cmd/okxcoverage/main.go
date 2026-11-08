// Command okxcoverage dry-runs schema-v4 wallet/reserve coverage arithmetic for an OKX audit.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ardmere/ardmere/internal/coingecko"
	okx "github.com/ardmere/ardmere/internal/exchanges/okx"
	"github.com/ardmere/ardmere/internal/exchanges/okx/okxapi"
	"github.com/ardmere/ardmere/internal/pipeline"
	"github.com/ardmere/ardmere/internal/verifier"
)

func main() {
	auditID := flag.String("audit", "508756245", "OKX audit id")
	assessmentPath := flag.String("assessment", "", "optional assessment JSON (default docs/reports/okx/<audit>-assessment.json)")
	flag.Parse()

	if *assessmentPath == "" {
		*assessmentPath = filepath.Join("docs", "reports", "okx", *auditID+"-assessment.json")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	assessRaw, err := os.ReadFile(*assessmentPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read assessment: %v\n", err)
		os.Exit(1)
	}
	var assess map[string]any
	if err := json.Unmarshal(assessRaw, &assess); err != nil {
		fmt.Fprintf(os.Stderr, "decode assessment: %v\n", err)
		os.Exit(1)
	}

	bundle, _, err := okxapi.FetchSummaryByAuditID(ctx, *auditID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch okx summary %s: %v\n", *auditID, err)
		os.Exit(1)
	}
	snap := okx.Normalize(bundle.Audit, okx.Meta{Exchange: "okx"})

	published := publishedFromAssessment(assess)
	vs := verificationsFromAssessment(assess)

	fmt.Printf("OKX coverage dry-run audit=%s snapshotTime=%s\n", snap.ID, snap.SnapshotTime.UTC().Format(time.RFC3339))
	fmt.Printf("assessment porStage=%v effectivePoR=%v\n", assess["porStage"], assess["effectivePoR"])
	fmt.Printf("publishedAddressCount (from assessment ownership reason / fallback)=%d\n", published)
	fmt.Printf("ownership verdict from assessment=%s\n", ownershipVerdict(vs))
	fmt.Printf("on-chain balance verifiers from assessment: ")
	for _, v := range vs {
		switch v.VerifierID {
		case "onchain-balance-hot", "onchain-balance-token", "onchain-balance-ledger", "onchain-balance-deposit":
			fmt.Printf("%s=%s ", v.VerifierID, v.Verdict)
		}
	}
	fmt.Println()

	ratios := pipeline.ComputeCoverageRatios(ctx, pipeline.CoverageInput{
		Snapshot:              snap,
		Verifications:         vs,
		PublishedAddressCount: published,
		PublishedUnknown:      published == 0,
		AddressBalances:       nil, // no zip in this dry-run; ownership PASS alone cannot invent balances
		Prices:                coingecko.New(),
	})

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	fmt.Println("--- ratios ---")
	_ = enc.Encode(map[string]any{
		"walletCoverageBps":  ratios.WalletCoverageBps,
		"reserveCoverageBps": ratios.ReserveCoverageBps,
		"coverageDetail":     ratios.Detail,
		"note":               "This dry-run uses assessment verifier verdicts (ownership PASS, on-chain often UNVERIFIABLE when --skip-rpc). Verifiable set requires BOTH ownership PASS and balance PASS — so ratios stay 0 when RPC coverage is absent. ClaimedReserveUSD still prices exchange claimed asset amounts with CoinGecko at snapshotTime.",
	})
}

func ownershipVerdict(vs []verifier.Verification) string {
	for _, v := range vs {
		if v.VerifierID == "address-ownership" {
			return string(v.Verdict)
		}
	}
	return "missing"
}

func publishedFromAssessment(assess map[string]any) int {
	// Prefer ownership reason like "verified 682007/682007 address signatures".
	results, _ := assess["verifierResults"].([]any)
	for _, raw := range results {
		m, _ := raw.(map[string]any)
		if m["verifierId"] != "address-ownership" {
			continue
		}
		reason, _ := m["reason"].(string)
		var pass, total int
		if _, err := fmt.Sscanf(reason, "verified %d/%d", &pass, &total); err == nil && total > 0 {
			return total
		}
	}
	return 0
}

func verificationsFromAssessment(assess map[string]any) []verifier.Verification {
	results, _ := assess["verifierResults"].([]any)
	out := make([]verifier.Verification, 0, len(results))
	for _, raw := range results {
		m, _ := raw.(map[string]any)
		id, _ := m["verifierId"].(string)
		verdict, _ := m["verdict"].(string)
		reason, _ := m["reason"].(string)
		version, _ := m["version"].(string)
		out = append(out, verifier.Verification{
			VerifierID: id,
			Version:    version,
			Verdict:    verifier.Verdict(verdict),
			Reason:     reason,
		})
	}
	return out
}

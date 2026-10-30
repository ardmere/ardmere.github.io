package sitegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ardmere/ardmere/internal/exchangeregistry"
)

func formatStageMD(stage string) string {
	if stage == "Stage 1" {
		return "**Stage 1**"
	}
	return stage
}

func formatEffectivePoR(yes bool) string {
	if yes {
		return "Yes"
	}
	return "No"
}

func formatGenEvidence(a Assessment) string {
	gen := a.Gen
	if gen == "" {
		gen = "Gen 1"
	}
	ev := a.EvidenceLevel
	if ev == "" {
		ev = "E1"
	}
	return fmt.Sprintf("`%s / %s`", gen, ev)
}

func formatSnapshotDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func formatSnapshotMonthYear(t time.Time) string {
	return t.Format("Jan 2006")
}

func formatBinanceListTime(t time.Time) string {
	return t.Format("02/01/06")
}

func reportLinkMD(rel string) string {
	return fmt.Sprintf("[report](%s)", rel)
}

func fullReportLinkMD(rel string) string {
	return fmt.Sprintf("[Full report](%s)", rel)
}

func binanceSnapshotNote(a Assessment) string {
	hasBAPI := hasArtifactKind(a, "bapiSnapshot")
	hasWallet := hasArtifactKind(a, "walletZip")
	switch {
	case hasBAPI && hasWallet:
		if a.PeriodSeq > 0 {
			return fmt.Sprintf("Full BAPI summary + wallet zip (period %d)", a.PeriodSeq)
		}
		return "Full BAPI summary + wallet zip"
	case hasWallet:
		return "Wallet zip + list metadata; BAPI summary not public for historical period"
	default:
		return "Public PoR artifacts archived"
	}
}

func secondarySnapshotNote(a Assessment) string {
	if note := snapshotNoteOverride(a.Exchange, a.SnapshotID); note != "" {
		return note
	}
	if fixtureImportExists(a.Exchange, a.SnapshotID) {
		return "Fixture import"
	}
	if hasArtifactKind(a, "summarySnapshot") && !hasArtifactKind(a, "walletZip") {
		return "Public summary snapshot (scraped)"
	}
	if hasArtifactKind(a, "walletZip") && hasArtifactKind(a, "globalProofBundle") {
		return "Full public artifact stack"
	}
	return "Public PoR artifacts archived"
}

func snapshotNoteOverride(exchange, snapshotID string) string {
	notes, err := loadSnapshotNotes()
	if err != nil || notes == nil {
		return ""
	}
	return notes.Overrides[exchange+"/"+snapshotID]
}

func fixtureImportExists(exchange, snapshotID string) bool {
	candidates := []string{
		filepath.Join("fixtures", exchange, snapshotID),
		filepath.Join("fixtures", exchange, snapshotID+"-summary.json"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

func upstreamVerifierMarkdown(ex exchangeregistry.Exchange) string {
	if ex.ID == "okx" {
		var v1, v2 string
		for _, repo := range ex.UpstreamRepos {
			short := shortenRepoURL(repo.URL)
			link := fmt.Sprintf("[%s](%s)", short, repo.URL)
			switch repo.Role {
			case "global_zk_verifier_spec":
				v2 = link
			case "address_ownership_verifier":
				v1 = link
			}
		}
		if v2 != "" && v1 != "" {
			var v1URL string
			for _, repo := range ex.UpstreamRepos {
				if repo.Role == "address_ownership_verifier" {
					v1URL = repo.URL
					break
				}
			}
			return v2 + " (+ [v1 sigs](" + v1URL + "))"
		}
	}
	for _, repo := range ex.UpstreamRepos {
		short := shortenRepoURL(repo.URL)
		return fmt.Sprintf("[%s](%s)", short, repo.URL)
	}
	return "—"
}

func shortenRepoURL(u string) string {
	const prefix = "https://github.com/"
	if strings.HasPrefix(u, prefix) {
		return u[len(prefix):]
	}
	return u
}

func exchangeHistoryAnchor(exchange string) string {
	switch exchange {
	case "binance":
		return "./snapshot-history.md#binance-monthly-bapi-summary-only-for-current-snapshot"
	case "okx":
		return "./snapshot-history.md#okx-monthly-full-public-artifact-stack"
	default:
		return ""
	}
}

func snapshotLabel(a Assessment) string {
	label := formatSnapshotMonthYear(a.SnapshotTime)
	if a.Exchange == "binance" && a.PeriodSeq > 0 {
		return fmt.Sprintf("%s, period %d", label, a.PeriodSeq)
	}
	return label
}

func exchangeBullet(a Assessment, multiSnapshot bool) string {
	history := ""
	if multiSnapshot {
		if anchor := exchangeHistoryAnchor(a.Exchange); anchor != "" {
			history = fmt.Sprintf(" · [history](%s)", anchor)
		}
	}
	return fmt.Sprintf("- **%s (`%s`)** — Latest `%s` (%s): %s %s%s",
		a.ExchangeName,
		a.Stage,
		a.SnapshotID,
		snapshotLabel(a),
		a.Headline,
		fullReportLinkMD(a.ReportRel),
		history,
	)
}

func countMultiSnapshotExchanges(rows []Assessment) map[string]bool {
	counts := map[string]int{}
	for _, a := range rows {
		counts[a.Exchange]++
	}
	out := map[string]bool{}
	for ex, n := range counts {
		out[ex] = n > 1
	}
	return out
}

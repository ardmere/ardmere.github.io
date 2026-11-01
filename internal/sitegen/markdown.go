package sitegen

import (
	"fmt"
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

func reportLinkMD(rel string) string {
	return fmt.Sprintf("[report](%s)", rel)
}

func fullReportLinkMD(rel string) string {
	return fmt.Sprintf("[Full report](%s)", rel)
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

func exchangeHistoryAnchor(_ string) string {
	return "./snapshot-history.md#all-public-reports"
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
		history = fmt.Sprintf(" · [history](%s)", exchangeHistoryAnchor(a.Exchange))
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

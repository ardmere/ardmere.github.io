package sitegen

import (
	"fmt"
	"strings"

	"github.com/ardmere/ardmere/internal/exchangeregistry"
)

// WriteExchangeComparison generates docs/reports/exchange-comparison.md from assessments.
func WriteExchangeComparison(path string, rows []Assessment, reg exchangeregistry.Registry) error {
	latest := latestByExchange(rows, reg)
	multi := countMultiSnapshotExchanges(rows)

	var b strings.Builder
	b.WriteString("# Exchange PoR Transparency Comparison\n\n")
	b.WriteString("> Methodology: [PoR Stage Framework](../por-transparency-framework.md)  \n")
	b.WriteString("> Policy standard: [Effective PoR Standard](../effective-por-standard.md)  \n")
	b.WriteString(fmt.Sprintf("> **All %d public reports:** [Snapshot history](./snapshot-history.md)\n\n", len(rows)))

	b.WriteString("## Summary (latest snapshot per exchange)\n\n")
	b.WriteString("| Exchange | Snapshot | PoR Stage | Gen / Evidence | Effective PoR | Confidence | Last Verified Snapshot | Upstream verifier | Report |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- | --- | --- | --- |\n")
	for _, a := range latest {
		ex := reg.Find(a.Exchange)
		upstream := "—"
		if ex != nil {
			upstream = upstreamVerifierMarkdown(*ex)
		}
		effective := "false"
		if a.EffectivePoR {
			effective = "true"
		}
		conf := a.Confidence
		if conf == "" {
			conf = "low"
		}
		b.WriteString(fmt.Sprintf("| %s | `%s` | `%s` | %s | `%s` | `%s` | `%s` | %s | %s |\n",
			a.ExchangeName,
			a.SnapshotID,
			a.Stage,
			formatGenEvidence(a),
			effective,
			conf,
			a.SnapshotTime.Format("2006-01-02T15:04:05Z"),
			upstream,
			reportLinkMD(a.ReportRel),
		))
	}

	b.WriteString("\n## By exchange\n\n")
	for _, a := range latest {
		b.WriteString(exchangeBullet(a, multi[a.Exchange]))
		b.WriteString("\n")
	}

	b.WriteString("\nGap detail, risk flags, and verifier output are in each report. Missing data is marked `UNVERIFIABLE`, not `PASS` ([methodology §1.3](../por-transparency-framework.md#gap-reporting-discipline)).\n\n")
	b.WriteString("Canonical upstream registry (PoR pages + GitHub repos): [`config/exchanges/registry.yaml`](../../config/exchanges/registry.yaml) · CLI: `por exchanges`.\n\n")
	b.WriteString("## Evidence archive\n\n")
	b.WriteString("Assessment JSON, bundle roots, and hashes for this report set: [Artifact archive index](./artifact-archive-index.md).\n")

	return osWriteFile(path, b.String())
}

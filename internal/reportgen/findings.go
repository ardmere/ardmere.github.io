package reportgen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ardmere/ardmere/internal/verifier"
)

const (
	maxWarnDetailRows           = 15
	maxUnverifiableDetailRows   = 20
	maxUnsupportedPairRows      = 25
)

type statusCounts map[verifier.Verdict]int

func countFindingStatuses(findings []verifier.Finding) statusCounts {
	c := statusCounts{}
	for _, f := range findings {
		c[f.Status]++
	}
	return c
}

func writeSection6FindingDetails(b *strings.Builder, vers []verifier.Verification, snapshotID string) {
	b.WriteString("## 6. Verifier Finding Details\n\n")
	b.WriteString("This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. ")
	b.WriteString("Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. ")
	fmt.Fprintf(b, "Full machine-readable output: `%s.verification-bundle.v2.json` in the local artifact bundle.\n\n", snapshotID)

	writeStubUnverifiableVerifiers(b, vers)

	var detail []verifier.Verification
	for _, v := range vers {
		if v.Verdict == verifier.VerdictPass && len(findingsOfStatus(v.Findings, verifier.VerdictFail)) == 0 &&
			len(findingsOfStatus(v.Findings, verifier.VerdictWarn)) == 0 {
			continue
		}
		if v.Verdict == verifier.VerdictUnverifiable && len(v.Findings) == 0 {
			continue
		}
		detail = append(detail, v)
	}

	if len(detail) == 0 {
		b.WriteString("No additional row-level `FAIL` or `WARN` findings beyond the summary table in §5.\n\n")
		return
	}

	b.WriteString("### Row-level findings\n\n")
	for _, v := range detail {
		writeVerifierFindingBlock(b, v)
	}
}

func writeStubUnverifiableVerifiers(b *strings.Builder, vers []verifier.Verification) {
	var stubs []verifier.Verification
	for _, v := range vers {
		if v.Verdict == verifier.VerdictUnverifiable && len(v.Findings) == 0 && strings.TrimSpace(v.Reason) != "" {
			stubs = append(stubs, v)
		}
	}
	if len(stubs) == 0 {
		return
	}
	b.WriteString("### Capability and artifact gaps (`UNVERIFIABLE`)\n\n")
	b.WriteString("| Verifier | Explanation |\n| --- | --- |\n")
	for _, v := range stubs {
		fmt.Fprintf(b, "| `%s` | %s |\n", v.VerifierID, escapeCell(v.Reason))
	}
	b.WriteString("\n")
}

func writeVerifierFindingBlock(b *strings.Builder, v verifier.Verification) {
	fmt.Fprintf(b, "#### `%s` (`%s`)\n\n", v.VerifierID, strings.ToUpper(string(v.Verdict)))
	if v.Reason != "" {
		fmt.Fprintf(b, "Summary: %s\n\n", escapeCell(v.Reason))
	}
	counts := countFindingStatuses(v.Findings)
	if len(counts) > 0 {
		fmt.Fprintf(b, "Finding counts: %s\n\n", formatStatusCounts(counts))
	}

	fails := findingsOfStatus(v.Findings, verifier.VerdictFail)
	warns := findingsOfStatus(v.Findings, verifier.VerdictWarn)
	unver := findingsOfStatus(v.Findings, verifier.VerdictUnverifiable)

	if len(fails) > 0 {
		b.WriteString("**FAIL**\n\n")
		writeFindingTable(b, fails)
	}
	if len(warns) > 0 {
		b.WriteString("**WARN**\n\n")
		rows := warns
		extra := ""
		if len(warns) > maxWarnDetailRows {
			rows = warns[:maxWarnDetailRows]
			extra = fmt.Sprintf("\n_%d additional `WARN` rows omitted; see verification bundle._\n\n", len(warns)-maxWarnDetailRows)
		}
		writeFindingTable(b, rows)
		b.WriteString(extra)
	}
	if len(unver) > 0 {
		b.WriteString("**UNVERIFIABLE (row-level)**\n\n")
		writeUnverifiableFindings(b, unver)
	}
}

func writeUnverifiableFindings(b *strings.Builder, findings []verifier.Finding) {
	var rowCount []verifier.Finding
	var substantive []verifier.Finding
	for _, f := range findings {
		if f.Field == "row_count" {
			rowCount = append(rowCount, f)
			continue
		}
		substantive = append(substantive, f)
	}

	if len(rowCount) > 0 {
		b.WriteString("Unsupported (coin, network) pairs:\n\n")
		b.WriteString("| Pair | Rows | Note |\n| --- | --- | --- |\n")
		sort.Slice(rowCount, func(i, j int) bool {
			return rowCount[i].Subject < rowCount[j].Subject
		})
		shown := rowCount
		extra := ""
		if len(rowCount) > maxUnsupportedPairRows {
			shown = rowCount[:maxUnsupportedPairRows]
			extra = fmt.Sprintf("\n_%d additional unsupported pairs omitted; see verification bundle._\n\n", len(rowCount)-maxUnsupportedPairRows)
		}
		for _, f := range shown {
			note := f.Note
			if note == "" {
				note = "no native verifier for this (coin,network) pair yet"
			}
			fmt.Fprintf(b, "| `%s` | %s | %s |\n", escapeCell(f.Subject), escapeCell(f.Actual), escapeCell(note))
		}
		b.WriteString(extra)
	}

	if len(substantive) == 0 {
		return
	}
	rows := substantive
	extra := ""
	if len(substantive) > maxUnverifiableDetailRows {
		rows = substantive[:maxUnverifiableDetailRows]
		extra = fmt.Sprintf("\n_%d additional row-level `UNVERIFIABLE` rows omitted; see verification bundle._\n\n", len(substantive)-maxUnverifiableDetailRows)
	}
	b.WriteString("Other unverifiable rows:\n\n")
	writeFindingTable(b, rows)
	b.WriteString(extra)
}

func writeFindingTable(b *strings.Builder, findings []verifier.Finding) {
	b.WriteString("| Subject | Field | Claim | Actual | Note |\n| --- | --- | --- | --- | --- |\n")
	for _, f := range findings {
		fmt.Fprintf(b, "| `%s` | `%s` | %s | %s | %s |\n",
			escapeCell(f.Subject),
			escapeCell(f.Field),
			escapeCell(f.Claim),
			escapeCell(f.Actual),
			escapeCell(f.Note),
		)
	}
	b.WriteString("\n")
}

func findingsOfStatus(findings []verifier.Finding, status verifier.Verdict) []verifier.Finding {
	var out []verifier.Finding
	for _, f := range findings {
		if f.Status == status {
			out = append(out, f)
		}
	}
	return out
}

func formatStatusCounts(c statusCounts) string {
	order := []verifier.Verdict{verifier.VerdictFail, verifier.VerdictWarn, verifier.VerdictUnverifiable, verifier.VerdictPass, verifier.VerdictPartial}
	var parts []string
	for _, s := range order {
		if n := c[s]; n > 0 {
			parts = append(parts, fmt.Sprintf("%s %d", strings.ToUpper(string(s)), n))
		}
	}
	return strings.Join(parts, ", ")
}

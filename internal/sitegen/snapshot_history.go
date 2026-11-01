package sitegen

import (
	"fmt"
	"strings"

	"github.com/ardmere/ardmere/internal/exchangeregistry"
)

// WriteSnapshotHistory generates docs/reports/snapshot-history.md from assessments.
func WriteSnapshotHistory(path string, rows []Assessment, _ exchangeregistry.Registry) error {
	assessed := latestGeneratedAt(rows)
	var b strings.Builder
	b.WriteString("# PoR Snapshot History (Public Evaluation Set)\n\n")
	b.WriteString("> Latest snapshot per exchange: [Exchange comparison](./exchange-comparison.md)  \n")
	b.WriteString("> Methodology: [PoR Stage Framework](../por-transparency-framework.md)\n\n")
	b.WriteString("This index lists every assessed snapshot in the ardmere public evaluation set. ")
	b.WriteString("The [homepage](../../index.html) shows the six most recent reports (auto-generated via `por site`). ")
	b.WriteString("This file is also auto-generated; edit `config/site/snapshot-notes.yaml` for page footers.\n\n")
	if !assessed.IsZero() {
		b.WriteString(fmt.Sprintf("**Assessed:** %s · Schema: `ardmere/exchange-assessment@1`\n\n", assessed.Format("2006-01-02")))
	} else {
		b.WriteString("**Schema:** `ardmere/exchange-assessment@1`\n\n")
	}

	b.WriteString("## All public reports\n\n")
	b.WriteString("| Exchange | Snapshot | Date | Stage | Report |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")
	for _, a := range rows {
		b.WriteString(fmt.Sprintf("| %s | `%s` | %s | %s | %s |\n",
			a.ExchangeName, a.SnapshotID, formatSnapshotDate(a.SnapshotTime), formatStageMD(a.Stage), reportLinkMD(a.ReportRel)))
	}
	b.WriteString("\nDates are snapshot times from each assessment JSON (`snapshotTime`, UTC date).\n\n")

	for _, footer := range snapshotFooters() {
		b.WriteString(footer)
		b.WriteString("\n\n")
	}

	b.WriteString("## Regenerating this set\n\n")
	b.WriteString("Ensure `~/.zshenv` exports RPC keys (`INFURA_KEY`, `ALCHEMY_KEY`, `BLOCKCHAIR_API_KEY`, …) — see [ledger-rpc-runbook.md](../ledger-rpc-runbook.md). `./scripts/por.sh` loads `~/.zshenv` and repo `.env` automatically.\n\n")
	b.WriteString("```bash\n")
	b.WriteString("./scripts/por-batch-recent.sh 3 okx,binance\n")
	b.WriteString("go run ./cmd/por site\n")
	b.WriteString("./scripts/build-public-html.sh\n")
	b.WriteString("```\n\n")
	b.WriteString("Use `go run ./cmd/por batch -help` for flags (`-skip-rpc`, `-full-rpc-latest-only`, etc.).\n")

	return osWriteFile(path, b.String())
}

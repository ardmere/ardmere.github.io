# Documentation map

## Public site product

These pages are part of the ardmere public transparency registry. They are linked from [index.html](../index.html) and the site README.

**Site URLs use `.html`** (styled pages for GitHub Pages). Markdown remains the editable source; regenerate with [`scripts/build-public-html.sh`](../scripts/build-public-html.sh).

| Document | Purpose |
| --- | --- |
| [por-transparency-framework.html](./por-transparency-framework.html) | Methodology and effective PoR standard |
| [reports/exchange-comparison.html](./reports/exchange-comparison.html) | Exchange comparison table |
| [reports/artifact-archive-index.html](./reports/artifact-archive-index.html) | Artifact archive index |
| [reports/*/transparency-report.html](./reports/) | Exchange transparency reports |
| [reports/*/*-assessment.json](./reports/) | Machine-readable assessments bound to each report |
| [ardmere-service-audience.html](./ardmere-service-audience.html) | Product audience and positioning |
| [deposit-spot-check.md](./deposit-spot-check.md) | Deposit spot-check spec ([`/verify/deposit/`](../verify/deposit/)) |

## Repository-only documentation

These files stay in the public git repository for contributors and operators, but they are **not** linked from the public site navigation.

| Document | Purpose |
| --- | --- |
| [STRUCTURE.md](./STRUCTURE.md) | Repository layout |
| [por-cli.md](./por-cli.md) | CLI reference |
| [verifier-architecture.md](./verifier-architecture.md) | Verifier architecture |
| [deployments.md](./deployments.md) | On-chain deployment records |
| [decisions.md](./decisions.md) | Architecture decision records |
| [exchange-tiers.md](./exchange-tiers.md) | Verifier coverage tiers (VC) |
| [binance-por-data-guide.md](./binance-por-data-guide.md) | Binance data integration |
| [okx-por-data-guide.md](./okx-por-data-guide.md) | OKX data integration |
| [bybit-por-data-guide.md](./bybit-por-data-guide.md) | Bybit data integration |
| [bitget-por-data-guide.md](./bitget-por-data-guide.md) | Bitget data integration |
| [htx-por-data-guide.md](./htx-por-data-guide.md) | HTX data integration |
| [gate-por-data-guide.md](./gate-por-data-guide.md) | Gate.io data integration |

## Internal / non-product

Not part of the public site. Kept for report generation only.

| Path | Purpose |
| --- | --- |
| [templates/](./templates/) | Report and assessment generation templates |

When adding new docs, put public product pages in the first table, contributor docs in the second, and generation or obsolete material in the third.

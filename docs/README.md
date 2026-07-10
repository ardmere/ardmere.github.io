# Documentation map

## Public site product

These pages are part of the ardmere public transparency registry. They are linked from [index.html](../index.html) and the site README.

**Site URLs use `.html`** (styled pages for GitHub Pages). Markdown remains the editable source; regenerate with [`scripts/build-public-html.sh`](../scripts/build-public-html.sh).

| Document | Purpose |
| --- | --- |
| [por-transparency-framework.html](./por-transparency-framework.html) | Methodology and effective PoR standard |
| [effective-por-standard.html](./effective-por-standard.html) | Minimum effective PoR policy (public standard page) |
| [insights/index.html](./insights/index.html) | Editorial insights hub (not exchange ratings) |
| [reports/exchange-comparison.html](./reports/exchange-comparison.html) | Exchange comparison table and per-exchange reports |
| [reports/artifact-archive-index.html](./reports/artifact-archive-index.html) | Evidence archive (Reports sub-page): bundles, hashes, assessment JSON |
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
| [ledger-rpc-runbook.md](./ledger-rpc-runbook.md) | On-chain / ledger RPC ops: env vars, failure modes, fix order |
| [deployments.md](./deployments.md) | On-chain deployment records |
| [anchor-query-api.md](./anchor-query-api.md) | ArdmerePoRAnchor on-chain query guide (`cast` / `ethers.js`) |
| [decisions.md](./decisions.md) | Architecture decision records |
| [binance-por-data-guide.md](./binance-por-data-guide.md) | Binance data integration |
| [okx-por-data-guide.md](./okx-por-data-guide.md) | OKX data integration |
| [bybit-por-data-guide.md](./bybit-por-data-guide.md) | Bybit data integration |
| [bitget-por-data-guide.md](./bitget-por-data-guide.md) | Bitget data integration |
| [htx-por-data-guide.md](./htx-por-data-guide.md) | HTX data integration |
| [gate-por-data-guide.md](./gate-por-data-guide.md) | Gate.io data integration |
| [exchange-upstream-registry.md](./exchange-upstream-registry.md) | Official PoR pages + GitHub repos (human index) |
| [../config/exchanges/registry.yaml](../config/exchanges/registry.yaml) | Machine-readable upstream registry (`por exchanges`) |
| [exchange-reserve-transparency-whitepaper.md](./exchange-reserve-transparency-whitepaper.md) | **Whitepaper body (English draft v1.0)** — exchange reserve transparency framework and methodology |

Project-specific Cursor Agent Skills live in `.cursor/skills/` (gitignored; not published to the site).

## Internal / non-product

Not part of the public site. Kept for report generation only.

| Path | Purpose |
| --- | --- |
| [templates/](./templates/) | Report and assessment generation templates |

When adding new docs, put public product pages in the first table, contributor docs in the second, and generation or obsolete material in the third.

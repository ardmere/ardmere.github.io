# Documentation map

## Public site product

These pages are part of the ardmere public transparency registry. They are linked from [index.html](../index.html) and the site README.

**Site URLs use `.html`** (styled pages for GitHub Pages). Markdown remains the editable source; regenerate with [`scripts/build-public-html.sh`](../scripts/build-public-html.sh) using **pandoc 3.1.3** (see [`por-cli.md`](./por-cli.md#regenerating-reports-and-public-html) — other versions change table markup and add unrelated diff noise).

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
| [contract-query-guide.html](./contract-query-guide.html) | How to query `ArdmerePoRAnchor` for exchange audit / anchor records |
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
| [anchor-query-api.md](./anchor-query-api.md) | Extended ABI reference (`cast` / `ethers.js`); public summary is [contract-query-guide.html](./contract-query-guide.html) |
| [decisions.md](./decisions.md) | Architecture decision records |
| [binance-por-data-guide.md](./binance-por-data-guide.md) | Binance data integration |
| [okx-por-data-guide.md](./okx-por-data-guide.md) | OKX data integration |
| [bybit-por-data-guide.md](./bybit-por-data-guide.md) | Bybit data integration |
| [bitget-por-data-guide.md](./bitget-por-data-guide.md) | Bitget data integration |
| [htx-por-data-guide.md](./htx-por-data-guide.md) | HTX data integration |
| [gate-por-data-guide.md](./gate-por-data-guide.md) | Gate.io data integration |
| [poloniex-por-data-guide.md](./poloniex-por-data-guide.md) | Poloniex data integration |
| [ourbit-por-data-guide.md](./ourbit-por-data-guide.md) | Ourbit data integration |
| [coinex-por-data-guide.md](./coinex-por-data-guide.md) | CoinEx data integration |
| [bitunix-por-data-guide.md](./bitunix-por-data-guide.md) | Bitunix data integration |
| [woox-por-data-guide.md](./woox-por-data-guide.md) | WOO X data integration |
| [toobit-por-data-guide.md](./toobit-por-data-guide.md) | Toobit data integration |
| [phemex-por-data-guide.md](./phemex-por-data-guide.md) | Phemex data integration |
| [btcc-por-data-guide.md](./btcc-por-data-guide.md) | BTCC data integration |
| [kraken-por-data-guide.md](./kraken-por-data-guide.md) | Kraken data integration |
| [coinw-por-data-guide.md](./coinw-por-data-guide.md) | CoinW data integration |
| [exchange-upstream-registry.md](./exchange-upstream-registry.md) | Official PoR pages + GitHub repos (human index) |
| [../config/exchanges/registry.yaml](../config/exchanges/registry.yaml) | Machine-readable upstream registry (`por exchanges`) |
| [por-monitor-schedule.md](./por-monitor-schedule.md) | Per-exchange PoR scan windows (replaces daily all-P0+P1) |
| [exchange-onboarding-shortlist.md](./exchange-onboarding-shortlist.md) | Next-10 exchanges to integrate (after current registry) |
| [../config/por-monitor/schedule.yaml](../config/por-monitor/schedule.yaml) | Machine-readable monitor cadence |
| [exchange-reserve-transparency-whitepaper.md](./exchange-reserve-transparency-whitepaper.md) | **Whitepaper body (English draft v1.0)** — exchange reserve transparency framework and methodology |

Project-specific Cursor Agent Skills live in `skills/` and are symlinked from `.cursor/skills/`:

| Skill | Purpose |
| --- | --- |
| [../skills/por-custody-risk/SKILL.md](../skills/por-custody-risk/SKILL.md) | Score an assessment against its `ArdmerePoRAnchor` record; no `decision` field |
| [../skills/por-inclusion-proof/SKILL.md](../skills/por-inclusion-proof/SKILL.md) | Verify a user's own inclusion proof locally and bind it to the published root |

Run `./scripts/check-skills.sh` after changing either one; `--live` also reads Sepolia.

## Internal / non-product

Not part of the public site. Kept for report generation only.

| Path | Purpose |
| --- | --- |
| [templates/](./templates/) | Report and assessment generation templates |

When adding new docs, put public product pages in the first table, contributor docs in the second, and generation or obsolete material in the third.

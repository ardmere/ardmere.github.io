# CoinW Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: coinw`).

CoinW publishes **monthly** Merkle-tree PoR with a marketing-page total USDT reserve ratio and login-gated user inclusion (`merkel_tree_cw.json` + closed-source `por-verify-*` binaries). ardmere tier: **1** (summary only). **Stage 0 hard cap** — Cloudflare WAF / geo-restriction blocks stable datacenter dashboard capture; no permissionless wallet ZIP.

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Total USDT reserve ratio | Help-center PoR announcement (e.g. Sep 2026) | Yes (HTML; often WAF from datacenters) | `summarySnapshot` |
| Snapshot time + merkle root UI | https://www.coinw.com/proof-of-reserves (`__NUXT_DATA__` SSR) | Partial (browser; geo modal; UI animation) | observational only |
| User Merkle JSON | Assets → Overview → Proof of Reserves → Download Data (`merkel_tree_cw.json`) | **Login** | `userMerkleProof` if imported; **not verified** |
| Offline verifier binaries | https://github.com/connectCoinw/proof-of-reserves | Yes (zip assets; **no source**) | — |
| Wallet addresses / amounts | Not permissionlessly published in assessed package | — | — |

## Cadence

| | |
|--|--|
| Typical publish | **Monthly** (announcement + PoR page panel) |
| Monitor window (suggested) | days around month-end / early next month **+ ~1 week** buffer |
| Latest assessed | `202609` (September 2026 announcement total **109.65%**) |

Monitor entry: [`schedule.yaml`](../config/por-monitor/schedule.yaml) `id: coinw` — `month_days` 20–31 with
−1/+3 buffer (spills into the first days of the next month), `capture: browser`,
`probe: announcement_browser_capture`, `confidence: low` (only one committed assessment plus one
unassessed browser observation). The entry carries a `waf:` note; see
[por-monitor-schedule.md](./por-monitor-schedule.md#coinw-为什么必须-capture-browser).

## Official references

| Resource | URL |
|----------|-----|
| PoR page | https://www.coinw.com/proof-of-reserves |
| Sep 2026 announcement | https://www.coinw.com/help-center/announcement-center/maintenance-updates/project-upgrades/2026-09-26/proof-of-reserves-update-september-2026/8128 |
| Offline verifier release | https://github.com/connectCoinw/proof-of-reserves/releases/tag/v1.0.0 |

## Fetch (import-first)

Datacenter IPs commonly receive **Cloudflare WAF 403 / “Just a moment…”** from `www.coinw.com`. Prefer the announcement fixture (or a browser-captured info JSON):

```bash
# Announcement / curated fixture (recommended)
por fetch coinw -summary-path ./fixtures/coinw/202609-summary.json

# Optional: wrap a raw info JSON capture
por fetch coinw -info-file ./coinw-info.json

# Live probe (usually fails with WAF tip)
por fetch coinw
```

Snapshot id for the assessed package is the compact month label (`202609`). Exact September snapshot clock was **not** published in the announcement (fixture uses `2026-09-01 00:00:00` as a month anchor and records that caveat).

### Oct 2026 browser probe (not assessed)

On 2026-11-06 an interactive browser session bypassed Cloudflare but hit a **Service Restricted** geo modal. The PoR panel showed Snapshot Time `2026-10-24 04:00:01 (UTC+8)` while merkle root and total ratio **changed during the session** (conflicting 64-hex strings; ratio range ~100.35–115.75%). Those observations are archived at `fixtures/coinw/browser-por-probe-20261024.json` and are **not** used as the machine-assessed root/rate.

## Verify

```bash
por verify -exchange coinw -snapshot 202609 -artifacts ./artifacts/coinw/202609
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported total reserve rate ≥ 100% (announcement TOTAL / USDT ratio) |

### Stubbed (UNVERIFIABLE) — honest blockers

| Gap | Why |
|-----|-----|
| No wallet list / amounts | Marketing ratio only → Stage 1 `walletZip` / on-chain unwired. |
| `user-merkle-proof@0` | `merkel_tree_cw.json` login-gated; GitHub ships binaries only (no source to port). |
| Dashboard merkle root | Announcement omits root; live SSR merkle UI was animated/unstable under geo-restriction. |
| Cloudflare WAF | Datacenter curl cannot stably refresh the PoR page / APIs. |
| Global ZK / third-party attestation | Not part of this disclosure model. |

## Boundaries

- **PASS on solvency-claim** = self-reported announcement total ratio only.
- Do **not** treat closed-source `por-verify-*` binaries as an ardmere-ported `user-merkle-proof@coinw-1` verifier.
- Do **not** promote the Oct-24 browser panel to an assessed snapshot until merkle + ratio are machine-stable from a permissionless channel.

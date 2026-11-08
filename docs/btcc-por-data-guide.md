# BTCC Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: btcc`).

BTCC publishes **monthly** Merkle-tree PoR snapshots dated the **15th** (UTC `createDate`, continuous since 2025-04-15). The public PoR page shows per-asset reserve ratios and a truncated merkle root. ardmere tier: **1** (summary only). **Stage 0 hard cap** — truncated root, no wallets/amounts.

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Per-period ratios + truncated root | `GET https://www.btcc.com/social-api/contract/reserve/reserveInfo?reportDate=YYYY-MM-DD` | Yes (JSON) | `summarySnapshot` |
| Public PoR overview | https://www.btcc.com/en-US/proof-of-reserves | Yes (SPA) | observational |
| Report history page | `GET …/getReserveHistoryPage` | **Login** (token) | — |
| Leaf validate | `POST …/validate` `{leafHash, reportDate}` | **Login** | `userMerkleProof` if imported; **not verified** (no OSS) |
| Wallet amounts / addresses | Not published | — | — |
| Global ZK / third-party attestation | Not observed | — | — |

## Cadence

| | |
|--|--|
| Typical publish | **Monthly ~day 15** (`createDate` = `YYYY-MM-DD` with day 15) |
| Monitor window (suggested) | days **14–17** UTC (+buffer) |
| Latest assessed | `20261015` (2026-10-15); `2026-11-15` not published as of 2026-11-06 |

UI default date rule (from PoR page JS): on/after day 15 of the month → that month’s 15th; otherwise prior month’s 15th.

Confirmed by per-month `reserveInfo` probes on 2026-11-07: `createDate` is the 15th for twelve consecutive
months (2025-12-15 … 2026-10-15), and 2026-11-15 returns `data: null` (not yet published).

Monitor entry: [`schedule.yaml`](../config/por-monitor/schedule.yaml) `id: btcc` — `month_days` 14–17,
`confidence: high`. See [por-monitor-schedule.md](./por-monitor-schedule.md).

## Official references

| Resource | URL |
|----------|-----|
| PoR page | https://www.btcc.com/en-US/proof-of-reserves |
| Support: PoR feature | https://www.btcc.com/en-US/support-center/detail/46861585161113 |
| Example press (Sep 2026) | https://www.btcc.com/en-US/newsroom/company-news/btcc-september-2026-por-report |

## Fetch

```bash
# Live public API (latest published day-15)
por fetch btcc
por fetch btcc -audit-id 2026-10-15

# Or import a saved bundle fixture
por fetch btcc -summary-path ./fixtures/btcc/20261015-summary.json
```

`reserveInfo` returns `rate` strings like `"147%"` for BTC/ETH/USDT/USDC/XRP/ADA and `rootHash` (~16 hex). Snapshot id is the compact calendar form of `createDate` (e.g. `2026-10-15` → `20261015`). **No** reserve or liability amounts are published.

## Verify

```bash
por verify -exchange btcc -snapshot 20261015 -artifacts ./artifacts/btcc/20261015
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported per-coin reserve rates ≥ 100% (min rate also stored as total) |

### Stubbed (UNVERIFIABLE) — honest blockers

| Gap | Why |
|-----|-----|
| Truncated merkle root | Public `rootHash` is ~16 hex chars — not a full root for offline Merkle binding. |
| No wallet list / amounts | Ratios only → Stage 1 `walletZip` / on-chain unwired by design. |
| `user-merkle-proof@0` | Leaf hash from “View My Report” + `/validate` need login. **No official OSS Merkle verifier**. |
| History list API | `getReserveHistoryPage` returns login-expired without token (permissionless probing uses per-date `reserveInfo`). |
| Global ZK / attestation | Not published as reproducible public artifacts. |

## Boundaries

- **PASS on solvency-claim** = self-reported ratios from permissionless `reserveInfo` only.
- Do **not** treat the truncated `rootHash` as a full-period Merkle commitment suitable for Stage 1.
- Do **not** expect Stage 1 without permissionless wallet amounts/addresses and a full root (shortlist: Stage 0 hard cap).

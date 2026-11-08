# Phemex Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: phemex`).

Phemex publishes **monthly** Merkle-tree PoR. The public PoR UI shows per-asset reserve ratios, a logged-in Merkle inclusion tool, and **Part of Our Cold Wallets**. ardmere tier: **1** (summary only for Stage 0).

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Version history + chain/cex balances | `GET https://por.phemex.com/porVersions` | Yes (JSON) | `summarySnapshot` |
| Public PoR overview + cold wallets UI | https://phemex.com/proof-of-reserves | Yes (SPA; wallets listed, **no download**) | observational only |
| Merkle root / reserve-funds / wallet-addresses APIs | `api10.phemex.com/phemex-user/public/v2/por/*` | CloudFront / auth gated from datacenters | — |
| User Merkle inclusion (hashed client ID) | Logged-in PoR tool | Login | `userMerkleProof` if imported; **not verified** (no OSS) |
| Global ZK / third-party attestation | Not observed as public artifacts | — | — |

## Cadence

| | |
|--|--|
| Typical publish | 月更，1 日 |
| Monitor window (UTC) | 每月 1–3 日 |

`por.phemex.com/porVersions` 2026-11-07 实探 42 期：2024-06 起连续 30 个月全部 1 日；更早的 2022–2024-04 曾用 24–26 日。

Monitor entry: [`schedule.yaml`](../config/por-monitor/schedule.yaml) `id: phemex`. See [por-monitor-schedule.md](./por-monitor-schedule.md) for the confidence policy.

## Official references

| Resource | URL |
|----------|-----|
| PoR page | https://phemex.com/proof-of-reserves |
| porVersions API | https://por.phemex.com/porVersions |
| Help: what is PoR | https://phemex.com/help-center/what-is-the-phemex-proof-of-reserves-por-tool |
| Help: how to use PoR tool | https://phemex.com/help-center/how-do-i-use-the-phemex-proof-of-reserves-por-tool |
| Nov 2026 press (ratios) | https://www.globenewswire.com/news-release/2026/11/03/3394447/0/en/phemex-november-proof-of-reserves-highlights-asset-transparency-as-platform-expands-financial-services.html |

## Fetch

```bash
# Live public API (latest version with totalBalances)
por fetch phemex
por fetch phemex -audit-id 2026-11-01

# Or import a saved bundle fixture
por fetch phemex -summary-path ./fixtures/phemex/20261101-summary.json
```

`porVersions` entries expose `chainBalance` (wallet assets) and `cexBalance` (customer liabilities). The adapter computes per-coin `reserveRate = chain/cex * 100` for `solvency-claim@1`. Snapshot id is the compact calendar form of `versionID` (e.g. `2026-11-01` → `20261101`).

## Verify

```bash
por verify -exchange phemex -snapshot 20261101 -artifacts ./artifacts/phemex/20261101
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported per-coin reserve rates ≥ 100% (min rate also stored as total) |

### Stubbed (UNVERIFIABLE) — honest blockers

| Gap | Why |
|-----|-----|
| `user-merkle-proof@0` | Inclusion tool needs login (hashed client ID). **No official OSS Merkle verifier** under Phemex GitHub. |
| Period merkle root | Not present in `porVersions`; `api10 …/merkle-tree/root` CloudFront/auth gated from datacenters. |
| Wallet ZIP / on-chain | Cold wallets are **UI-visible** (“Part of Our Cold Wallets”) but there is no downloadable list; `api10 …/wallet-addresses` gated → Stage 1 `walletZip` unwired. |
| Global ZK / attestation | Not published as reproducible public artifacts. |

## UI notes (shortlist)

- **Monthly Merkle:** academy/help copy and the logged-in PoR tool describe Merkle inclusion; public overview does not print a period root hash for archival without the gated API.
- **Wallet lists on UI:** addresses are shown without login; many EVM L2s share one address. Treat as disclosure UI, not a permissionless machine-readable custody export.

## Boundaries

- **PASS on solvency-claim** = computed ratios from permissionless `porVersions` balances only.
- Press/UI may round or truncate ratios to 2 decimals slightly differently than full-precision `chain/cex`.
- Do not treat UI wallet visibility as Stage 1 until a current-period wallet list is permissionlessly fetchable and ownership-bound.

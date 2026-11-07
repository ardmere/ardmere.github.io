# Ourbit Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: ourbit`).

Ourbit publishes Merkle-tree PoR with public reserve ratios (`stock_info`) and a multi-chain wallet table. ardmere tier: **1** (summary Stage 0). User inclusion uses OSS [`OurbitOfficial/merkle_validate`](https://github.com/OurbitOfficial/merkle_validate) (Java; SHA-256 truncated to 32 hex / 128 bits).

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Reserve ratio + amounts | `api/assetbussiness/merkle/stock_info` | Yes | `summarySnapshot` |
| Multi-chain wallets | `api/assetbussiness/merkle/wallet_detail` (+ currencies) | Yes (live) | archived inside summary bundle; **not** Stage 1 `walletZip` |
| User Merkle leaf + path | Logged-in PoR → `merkel_tree_ourbit.json` | Login | stored + verified if imported |
| Period merkle root | Not in `stock_info` | No stable public field | — |
| Ownership / global ZK | — | No | — |

## Official references

| Resource | URL |
|----------|-----|
| PoR page | https://www.ourbit.com/assets/reserve-proof |
| Reserve ratio page | https://www.ourbit.com/assets/reserve-ratio |
| User merkle page | https://www.ourbit.com/assets/merkle |
| `stock_info` API | https://www.ourbit.com/api/assetbussiness/merkle/stock_info |
| `wallet_detail` API | https://www.ourbit.com/api/assetbussiness/merkle/wallet_detail |
| OSS verifier | https://github.com/OurbitOfficial/merkle_validate |
| Verify guide | https://www.ourbit.com/support/articles/17827791510977 |

## Fetch

```bash
# Live public APIs (no login)
por fetch ourbit

# Or import a saved bundle fixture
por fetch ourbit -summary-path ./fixtures/ourbit/20261101-summary.json

# Optional: import login-gated user proof
por fetch ourbit -summary-path ./fixtures/ourbit/20261101-summary.json \
  -user-proof ./fixtures/ourbit/merkel_tree_ourbit.json
```

API `stockRate` arrives as a fraction (`1.10` = 110%). The adapter normalizes to percentage units for `solvency-claim@1`.  
`totalReserveRate` is the **min** per-coin rate (conservative floor; API has no aggregate total).

Snapshot id is the UTC+8 calendar day of `snapshotDate` (e.g. epoch `…T16:00:00Z` → `20261101`).

## Verify

```bash
por verify -exchange ourbit -snapshot 20261101 -artifacts ./artifacts/ourbit/20261101
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported min/per-coin reserve rates ≥ 100% |
| `user-merkle-proof@ourbit-1` | Imported `merkel_tree_ourbit.json` recomputes under Ourbit's 128-bit truncated scheme |

### Stubbed (UNVERIFIABLE) — honest blockers

| Gap | Why |
|-----|-----|
| Public period merkle root | `stock_info` exposes ratios/amounts only — no root to bind user proofs without the login file |
| `user-merkle-proof` without import | Leaf/path download (`merkel_tree_ourbit.json`) and merkle list/download APIs are **login-gated** |
| Hash truncation | Leaf/parent hashes are SHA-256 truncated to **32 hex (128-bit)** — collision resistance ~64 bits; report as risk |
| Stage 1 wallets / ownership | `wallet_detail` is a **live** multi-chain table (amounts/heights update continuously; not snapshot-bound) and has **no ownership signatures** |
| Global ZK / attestation | Not published as reproducible public artifacts |

## Boundaries

- **PASS on solvency-claim** = self-reported `stock_info` ratios only.
- Live wallet rows are archived for transparency; ardmere does **not** treat them as permissionless Stage 1 `wallet_address_list`.
- Do not treat OSS `merkle_validate` availability as Stage 1 — it only proves login-gated user inclusion under a truncated hash scheme.

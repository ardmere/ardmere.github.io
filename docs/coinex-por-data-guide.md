# CoinEx Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: coinex`).

CoinEx publishes Merkle-tree PoR with public reserve ratios, a period merkle root, and a downloadable XLSX wallet address list. ardmere tier: **1** (summary Stage 0). User inclusion uses OSS [`coinexcom/proof-of-reserves`](https://github.com/coinexcom/proof-of-reserves) (Go; full SHA-256).

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Reserve ratio + amounts + root | `/res/account/balance/merkle/reserve` | Yes | `summarySnapshot` |
| Wallet address XLSX | `wallet_proof_file` URL on each period | Yes | URL archived in summary Extra; **not** Stage 1 `walletZip` |
| User Merkle leaf + path | Logged-in PoR → Copy my audition JSON | Login | stored + verified if imported |
| Ownership / global ZK | — | No | — |

## Cadence

| | |
|--|--|
| Typical publish | ~34 天（中位），日分散 |
| Monitor window (UTC) | `last + 34d ±5`；≥45 天静默强制补检 |

公开 `/res/account/balance/merkle/reserve` 提供完整期次列表：2026-11-07 实探 20 期（2024-11 → 2026-09），间隔中位 34 天、均值 35.3 天、范围 20–69 天，日分布 4–31 分散，**没有**可用的日历窗。注意 `docs/reports/` 只有一份，但 API 历史足够，故 `confidence: high`。

Monitor entry: [`schedule.yaml`](../config/por-monitor/schedule.yaml) `id: coinex`. See [por-monitor-schedule.md](./por-monitor-schedule.md) for the confidence policy.

## Official references

| Resource | URL |
|----------|-----|
| PoR page | https://www.coinex.com/en/reserve-proof |
| Check Merkle help | https://www.coinex.com/en/reserve-proof/help/check-merkle |
| Verify address help | https://www.coinex.com/en/reserve-proof/help/verify-address |
| Reserve API | https://www.coinex.com/res/account/balance/merkle/reserve |
| User proof API | https://www.coinex.com/res/account/balance/merkle/proof |
| OSS verifier | https://github.com/coinexcom/proof-of-reserves |
| Latest announcement (Aug 25, 2026) | https://www.coinex.com/en/announcements/detail/52911154866708 |

## Fetch

```bash
# Live public API (no login) — archives latest period
por fetch coinex

# Or import a saved bundle fixture
por fetch coinex -summary-path ./fixtures/coinex/20260914-summary.json

# Optional: import login-gated user proof
por fetch coinex -summary-path ./fixtures/coinex/20260914-summary.json \
  -user-proof ./fixtures/coinex/mock_user_proof.json
```

`reserve_rate` arrives already as a percentage (`107.42` = 107.42%).  
`totalReserveRate` is the **min** per-coin rate (conservative floor; API has no aggregate total).

Snapshot id is the UTC calendar day of `time` (unix seconds), e.g. `2026-09-14T07:00:00Z` → `20260914`.

## Verify

```bash
por verify -exchange coinex -snapshot 20260914 -artifacts ./artifacts/coinex/20260914
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported min/per-coin reserve rates ≥ 100% |
| `user-merkle-proof@coinex-1` | Imported audition JSON recomputes under CoinEx's full SHA-256 scheme and binds to summary `root_hash` when present |

### Stubbed (UNVERIFIABLE) — honest blockers

| Gap | Why |
|-----|-----|
| `user-merkle-proof` without import | `/res/account/balance/merkle/proof` returns **401** without login |
| Stage 1 wallets / ownership | Public `wallet_proof_file` is an **XLSX** address dump with **no ownership signatures** — not machine-verified as Stage 1 `walletZip` |
| Global ZK / attestation | Not published as reproducible public artifacts |
| Upstream example JSON | `coinexcom/proof-of-reserves` example currently fails its own validator; ardmere unit-tests a synthetic self-consistent fixture |

## Boundaries

- **PASS on solvency-claim** = self-reported API ratios only.
- Public XLSX is archived by URL for transparency; ardmere does **not** treat it as permissionless Stage 1 `wallet_address_list` with ownership proofs.
- Do not treat OSS verifier availability as Stage 1 — it only proves login-gated user inclusion under full SHA-256.

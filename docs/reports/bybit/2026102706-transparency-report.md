# BYBIT 2026102706 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [2026102706-assessment.json](./2026102706-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `bybit` |
| Snapshot | `2026102706` |
| Snapshot time | `2026-10-27T06:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

bybit snapshot 2026102706 remains Stage 0 in the ardmere public evaluation set.

Public artifacts do not support independently reproducible global solvency verification for this snapshot.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `wallet_address_list` | `UNVERIFIABLE` | Stage 0 | No public wallet_address_list. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-09-23T03:00:00Z` |
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `~231 days between snapshots` |
| History available | `3 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `8fe7cb998465dfdb0f2be8ad662c4880c6098b099e9c92c9b34af0e61371f3a0` | https://www.bybit.com/x-api/por/public/v1/reserve-ratio/latest |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x7b3625510a187a5f476c09cd37e916b23f72ec1bbc2371a2382b325d8896c5f5` |
| Verification bundle root | `0x8b93005b5808b08d3493cb0687178e1b18279ad8a38990b14500dfcc135c9754` |
| Artifact bundle SHA-256 | `61275b1063e64c845081cc06b66ed9e160a4deacd29b26d48146981797b6416c` |
| Verification bundle SHA-256 | `f7121c49d6f91195875c10aba66c0e1fe735ae6508593fa37e1facc8581c1126` |

Local bundle paths: [2026102706.artifact-bundle.json](../../../artifacts/bybit/2026102706/bundles/2026102706.artifact-bundle.json), [2026102706.verification-bundle.v2.json](../../../artifacts/bybit/2026102706/bundles/2026102706.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `FAIL` | `1.0000` |  |
| `user-merkle-proof` | `bybit-1` | `UNVERIFIABLE` | `0.0000` | myProof.json not present; download from Bybit PoR page while logged in and import via -user-proof |
| `internal-consistency` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | `0` | `UNVERIFIABLE` | `0.0000` | No public HotCold wallet address list |
| `onchain-balance-token` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address list for ERC20/BEP20 audit |
| `onchain-balance-ledger` | `0` | `UNVERIFIABLE` | `0.0000` | Verifier not available for this exchange snapshot |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Bybit uses pure Merkle Tree (no ZK); user myProof.json proves inclusion only |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | Hacken monthly PoR reports are public but not machine-readable in this pipeline |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `2026102706.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `user-merkle-proof` | myProof.json not present; download from Bybit PoR page while logged in and import via -user-proof |
| `internal-consistency` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | No public HotCold wallet address list |
| `onchain-balance-token` | No public wallet address list for ERC20/BEP20 audit |
| `onchain-balance-ledger` | Verifier not available for this exchange snapshot |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Bybit uses pure Merkle Tree (no ZK); user myProof.json proves inclusion only |
| `third-party-attestation` | Hacken monthly PoR reports are public but not machine-readable in this pipeline |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `solvency-claim` (`FAIL`)

Finding counts: FAIL 5

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `total` | `totalReserveRate` | 0 | 0.00 | self-reported exchange summary only; does not prove reserve authenticity |
| `BTC` | `reserveRate` | 0 | 0.00 | self-reported exchange summary only; does not prove reserve authenticity |
| `ETH` | `reserveRate` | 0 | 0.00 | self-reported exchange summary only; does not prove reserve authenticity |
| `USDT` | `reserveRate` | 0 | 0.00 | self-reported exchange summary only; does not prove reserve authenticity |
| `USDC` | `reserveRate` | 0 | 0.00 | self-reported exchange summary only; does not prove reserve authenticity |

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 2026102706 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-10-27T06:00:00Z |
| Is there reserve summary or wallet_address_list? (S0-3) | `pass` | Summary present; wallet_address_list absent. |
| Is wallet_address_list public and on-chain verifiable? (S1-1) | `fail` | No public wallet_address_list artifact. |
| Is wallet_ownership_proof public and batch-verifiable? (S1-2) | `unverifiable` | address-ownership verifier UNVERIFIABLE. |
| Is global_proof public and independently reproducible? (S1-3) | `unverifiable` | global-zk-proof verifier UNVERIFIABLE. |
| Is independent third-party review available with stated scope? (S1-4) | `unverifiable` | third-party-attestation verifier UNVERIFIABLE. |
| If trusted setup is required, is transcript public; otherwise is the proof system transparent setup? (S1-5) | `not_applicable` | No global ZK proof in the assessed public artifact set. |
| Is root/proof/vk canonically anchored outside exchange servers? (S2-1) | `not_applicable` | Stage 2 is not evaluated until Stage 1 blockers are resolved. |
| Is publication frequency sufficient for Stage 2 (weekly full PoR / daily anchor)? (S2-2) | `not_applicable` | Stage 2 is not evaluated until Stage 1 blockers are resolved. |

## 8. Recommendations

| Priority | Recommendation | Related Risk |
| --- | --- | --- |
| `P0` | Publish missing public artifacts required for independent reproduction. | `UNVERIFIABLE` |

## 9. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

## Capture notes / gaps (2026-10-31)

- `x-api/por/public/v1/reserve-ratio/latest` returns **404**; public page uses `queryHisBalancesInfo` / `queryHisChainInfo` (liability/audit-scope amounts only).
- Exact October reserve ratios were **not published** in press/Hacken PDF located at capture time — **not invented**. `solvency-claim` FAIL reflects missing ratios (rates stored as 0).
- Merkle root not present on public APIs. User `myProof.json` login-gated — `user-merkle-proof` UNVERIFIABLE.

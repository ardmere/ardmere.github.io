# BITGET 202605 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [202605-assessment.json](./202605-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `bitget` |
| Snapshot | `202605` |
| Snapshot time | `2026-05-27T10:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

bitget snapshot 202605 remains Stage 0 in the ardmere public evaluation set.

Public artifacts do not support independently reproducible global solvency verification for this snapshot.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| — | `UNVERIFIABLE` | Stage 0 | No public wallet_address_list. |

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `user-merkle-proof` | `bitget-1` | `PASS` | `1.0000` | Bitget user Merkle path valid under weak 64-bit truncated-hash scheme; proves inclusion only for supplied account |
| `internal-consistency` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | `0` | `UNVERIFIABLE` | `0.0000` | No public HotCold wallet address list |
| `onchain-balance-token` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address list for ERC20/BEP20 audit |
| `onchain-balance-ledger` | `0` | `UNVERIFIABLE` | `0.0000` | Verifier not available for this exchange snapshot |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Bitget uses pure Merkle Tree with 64-bit truncated SHA-256 (no ZK) |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No clearly identified independent third-party PoR attestation in machine-readable form |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | Wrapped token reconciliation rules not finalized |

## 8. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

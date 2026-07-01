# BYBIT 2025061709 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [2025061709-assessment.json](./2025061709-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `bybit` |
| Snapshot | `2025061709` |
| Snapshot time | `2025-06-17T09:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

bybit snapshot 2025061709 remains Stage 0 in the ardmere public evaluation set.

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
| `user-merkle-proof` | `bybit-1` | `PASS` | `1.0000` | user Merkle path valid (Gen 1.5 SHA-256); proves inclusion only for supplied account |
| `internal-consistency` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | `0` | `UNVERIFIABLE` | `0.0000` | No public HotCold wallet address list |
| `onchain-balance-token` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address list for ERC20/BEP20 audit |
| `onchain-balance-ledger` | `0` | `UNVERIFIABLE` | `0.0000` | Verifier not available for this exchange snapshot |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Bybit uses pure Merkle Tree (no ZK); user myProof.json proves inclusion only |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | Hacken monthly PoR reports are public but not machine-readable in this pipeline |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | Wrapped token reconciliation rules not finalized |

## 8. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

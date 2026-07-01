# OKX 508399035 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [508399035-assessment.json](./508399035-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `okx` |
| Snapshot | `508399035` |
| Snapshot time | `2026-06-18T16:00:00Z` |
| PoR Stage | `Stage 1 — Verifiable Disclosure` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `true` |

OKX publishes public summary, wallet_address_list, wallet_ownership_proof, and global zk proof bundle for audit 508399035. ardmere verifies address ownership and global zk binding, so this snapshot reaches Stage 1.

Stage 1 is supported by public wallet_address_list, wallet_ownership_proof, global_proof, and parameter/proof artifacts. Stage 2 is not reached because canonical official anchoring, stable DA, low-friction user inclusion proof, stronger publication frequency, and full business-consistent constraints are not established in this report.

## 2. Stage Decision

### Stage 1, Effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| — | `NO_CANONICAL_ANCHOR` | Stage 1 | No official exchange canonical on-chain/DA anchor established. |
| — | `HIGH_FREQUENCY_GAP` | Stage 1 | Monthly cadence does not satisfy Stage 2 frequency expectations. |

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.0` | `PASS` | `1.0000` |  |
| `address-ownership` | `okx-1` | `PASS` | `1.0000` | verified 196411/196411 address signatures |
| `onchain-balance-hot` | `2.1` | `FAIL` | `0.0022` |  |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.0024` |  |
| `onchain-balance-ledger` | `1.4` | `FAIL` | `0.0015` |  |
| `global-zk-proof` | `okx-1` | `PARTIAL` | `0.7500` | structure and summary binding OK; set OKX_ZK_STARK_VALIDATOR for cryptographic verification |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | BTC block time anchor verifier not implemented |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | Wrapped token reconciliation rules not finalized |

## 8. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

# OKX 506872725 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [506872725-assessment.json](./506872725-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `okx` |
| Snapshot | `506872725` |
| Snapshot time | `2026-05-06T16:00:00Z` |
| PoR Stage | `Stage 1 — Verifiable Disclosure` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `true` |

OKX publishes public summary, wallet_address_list, wallet_ownership_proof, and global zk proof bundle for audit 506872725. ardmere verifies address ownership and global zk binding, so this snapshot reaches Stage 1.

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
| `solvency-claim` | `1.0` | `FAIL` | `1.0000` | okx summary has no reserve currencies |
| `internal-consistency` | `1.0` | `PASS` | `1.0000` |  |
| `address-ownership` | `okx-1` | `PASS` | `1.0000` | verified 320917/320917 address signatures |
| `onchain-balance-hot` | `2.1` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-token` | `2.0` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-ledger` | `1.2` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `global-zk-proof` | `okx-1` | `PARTIAL` | `0.7500` | structure and summary binding OK; set OKX_ZK_STARK_VALIDATOR for cryptographic verification |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | BTC block time anchor verifier not implemented |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | Wrapped token reconciliation rules not finalized |

## 8. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

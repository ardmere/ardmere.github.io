# BINANCE PR01MAY26 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [PR01MAY26-assessment.json](./PR01MAY26-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `binance` |
| Snapshot | `PR01MAY26` |
| Snapshot time | `2026-01-05T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

Binance snapshot PR01MAY26 provides public reserve summary and wallet_address_list where available, but Stage 1 is blocked by missing wallet_ownership_proof, opaque trusted setup, and unavailable public global proof/vk artifacts.

The available artifacts support Gen 2 / E2 classification and Stage 0 PoR disclosure. Users still need to trust Binance for wallet control, trusted setup honesty, and public availability of the full global proof stack.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| — | `NO_WALLET_OWNERSHIP_PROOF` | Stage 0 | No public batch-verifiable wallet_ownership_proof. |
| — | `UNVERIFIABLE` | Stage 0 | Public global proof/vk not available. |

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.1` | `PASS` | `1.0000` |  |
| `btc-anchor` | `1` | `UNVERIFIABLE` | `0.0000` | BTC block time anchor verifier not implemented |
| `solvency-claim` | `1.0` | `FAIL` | `1.0000` | summary has no coin rows |
| `onchain-balance-hot` | `2.1` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-token` | `2.0` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-ledger` | `1.4` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-deposit` | `1.2` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | Wrapped token reconciliation rules not finalized |

## 8. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

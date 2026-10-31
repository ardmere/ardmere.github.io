# KUCOIN 20260630 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [20260630-assessment.json](./20260630-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `kucoin` |
| Snapshot | `20260630` |
| Snapshot time | `2026-06-30T15:59:59Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

kucoin snapshot 20260630 remains Stage 0 in the ardmere public evaluation set.

Public artifacts do not support independently reproducible global solvency verification for this snapshot.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `wallet_address_list` | `UNVERIFIABLE` | Stage 0 | No public wallet_address_list. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-09-30T15:59:59Z` |
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `UNVERIFIABLE` |
| History available | `1 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Only one kucoin snapshot is in the ardmere public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `39b1c81d41cc09feee1fc41f4b513643950c868715d85e8f1f35ea8a10811b14` | https://www.kucoin.com/proof-of-reserves |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0xc66b0b1522108f5c8eb3d47377228b6f6a0ebfe0fd19ab50cd785e1461646985` |
| Verification bundle root | `0xb5ddffb52336c397e1afdb2631b1aa1fe11c41e80392199e43d2b64da0a86367` |
| Artifact bundle SHA-256 | `76ec6ed4531c19ae5033eabcf6e903a34df04bafae14c98bc9c87f7f62be3e06` |
| Verification bundle SHA-256 | `9d972d9c1bbb5146a8eecb118471be688ec812ffce9710f336b08caf66a6c5f8` |

Local bundle paths: [20260630.artifact-bundle.json](../../../artifacts/kucoin/20260630/bundles/20260630.artifact-bundle.json), [20260630.verification-bundle.v2.json](../../../artifacts/kucoin/20260630/bundles/20260630.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `user-merkle-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Verifier not available for this exchange snapshot |
| `internal-consistency` | `0` | `UNVERIFIABLE` | `0.0000` | Wallet address bundle not available |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | BTC block time anchor verifier not implemented |
| `onchain-balance-hot` | `0` | `UNVERIFIABLE` | `0.0000` | Wallet address bundle not available for on-chain audit |
| `onchain-balance-token` | `0` | `UNVERIFIABLE` | `0.0000` | Wallet address bundle not available for token audit |
| `onchain-balance-ledger` | `0` | `UNVERIFIABLE` | `0.0000` | Verifier not available for this exchange snapshot |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `20260630.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `user-merkle-proof` | Verifier not available for this exchange snapshot |
| `internal-consistency` | Wallet address bundle not available |
| `btc-anchor` | BTC block time anchor verifier not implemented |
| `onchain-balance-hot` | Wallet address bundle not available for on-chain audit |
| `onchain-balance-token` | Wallet address bundle not available for token audit |
| `onchain-balance-ledger` | Verifier not available for this exchange snapshot |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

No additional row-level `FAIL` or `WARN` findings beyond the summary table in §5.

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 20260630 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-06-30T15:59:59Z |
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

# OKX 502299735 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [502299735-assessment.json](./502299735-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `okx` |
| Snapshot | `502299735` |
| Snapshot time | `2026-09-07T16:00:00Z` |
| PoR Stage | `Stage 1 — Verifiable Disclosure` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `true` |

OKX publishes public summary, wallet_address_list, wallet_ownership_proof, and global zk proof bundle for audit 502299735. ardmere verifies address ownership and global zk binding, so this snapshot reaches Stage 1.

Stage 1 is supported by public wallet_address_list, wallet_ownership_proof, global_proof, and parameter/proof artifacts. Stage 2 is not reached because canonical official anchoring, stable DA, low-friction user inclusion proof, stronger publication frequency, and full business-consistent constraints are not established in this report.

## 2. Stage Decision

### Stage 1, Effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `canonical on-chain/DA anchor` | `NO_CANONICAL_ANCHOR` | Stage 1 | No official exchange canonical on-chain/DA anchor established. |
| `weekly full PoR / daily anchor` | `HIGH_FREQUENCY_GAP` | Stage 1 | Monthly cadence does not satisfy Stage 2 frequency expectations. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-10-02T16:00:00Z` |
| Previous snapshot | `2026-08-10T16:00:00Z` |
| Observed cadence | `~monthly` |
| History available | `6 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Older snapshot 502260945 is also in the public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `1e10e866ace60e02a4d551516034b3eb6542e51009f0414c4c7cb17285c4f052` | https://www.okx.com/proof-of-reserves/detail |
| `walletZip` | `fa84690d06e7b5c883f3e6f24c635781a74321a67a76f1d0a7aab45528364641` | https://static.okx.com/cdn/okx/por/chain/por_csv_2026090800_V1.zip |
| `globalProofBundle` | `9b0508cab98f21bdbf391bc8e5b4aae8fd9acff6ef09eb182bacaa171c8f42d1` | https://static.okx.com/cdn/okx/por/merkel/por_502299735_proof_data.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x808ffacc9f2e261ce92e8cd824d274ebafce53028d31ed357ef85fc9a6ae45c3` |
| Verification bundle root | `0x9c0f5258fdd32c829c259bf71c16c4c45ffbe8ad3683d9c20caf8584294a7873` |
| Artifact bundle SHA-256 | `d329da9901b57664bb2e78bceb5721d3e37f0c0d01e154be138b40ff1844f5d9` |
| Verification bundle SHA-256 | `22d13ebf90e52726a370ce7e0d486b4d403eca5d05f70fdd7f9b02cc2dbf62bd` |

Local bundle paths: [502299735.artifact-bundle.json](../../../artifacts/okx/502299735/bundles/502299735.artifact-bundle.json), [502299735.verification-bundle.v2.json](../../../artifacts/okx/502299735/bundles/502299735.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `FAIL` | `1.0000` | okx summary has no reserve currencies |
| `internal-consistency` | `1.0` | `PASS` | `1.0000` |  |
| `address-ownership` | `okx-1.1` | `PASS` | `1.0000` | verified 496948/496948 address signatures (57 Aptos resource accounts via authentication_key) |
| `onchain-balance-hot` | `2.1` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-token` | `2.0` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-ledger` | `1.2` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `global-zk-proof` | `okx-1` | `PASS` | `1.0000` | zkSTARKValidator verify-global succeeded; summary merkle root bound to zk proof |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | BTC block time anchor verifier not implemented |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `502299735.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `onchain-balance-hot` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-token` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-ledger` | RPC queries skipped (--skip-rpc) |
| `btc-anchor` | BTC block time anchor verifier not implemented |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `solvency-claim` (`FAIL`)

Summary: okx summary has no reserve currencies

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 502299735 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-09-07T16:00:00Z |
| Is there reserve summary or wallet_address_list? (S0-3) | `pass` | Summary and wallet_address_list present. |
| Is wallet_address_list public and on-chain verifiable? (S1-1) | `pass` | Wallet list present; on-chain replay skipped or incomplete. |
| Is wallet_ownership_proof public and batch-verifiable? (S1-2) | `pass` | address-ownership verifier PASS. |
| Is global_proof public and independently reproducible? (S1-3) | `pass` | global-zk-proof verifier PASS. |
| Is independent third-party review available with stated scope? (S1-4) | `unverifiable` | third-party-attestation verifier UNVERIFIABLE. |
| If trusted setup is required, is transcript public; otherwise is the proof system transparent setup? (S1-5) | `pass` | Plonky2 / zk-STARK uses transparent setup; no trusted setup ceremony required. |
| Is root/proof/vk canonically anchored outside exchange servers? (S2-1) | `fail` | No official exchange canonical on-chain/DA anchor established. |
| Is publication frequency sufficient for Stage 2 (weekly full PoR / daily anchor)? (S2-2) | `fail` | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

## 8. Recommendations

| Priority | Recommendation | Related Risk |
| --- | --- | --- |
| `P1` | Anchor root/proof/commitment on-chain or publish immutable DA/archive records. | `NO_CANONICAL_ANCHOR` |
| `P1` | Add daily anchoring, weekly full PoR, or event-triggered updates to reduce timing risk. | `HIGH_FREQUENCY_GAP` |

## 9. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

# OKX 507918525 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [507918525-assessment.json](./507918525-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `okx` |
| Snapshot | `507918525` |
| Snapshot time | `2026-04-19T16:00:00Z` |
| PoR Stage | `Stage 1 — Verifiable Disclosure` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `true` |

OKX publishes public summary, wallet_address_list, wallet_ownership_proof, and global zk proof bundle for audit 507918525. ardmere verifies address ownership and global zk binding, so this snapshot reaches Stage 1.

Stage 1 is supported by public wallet_address_list, wallet_ownership_proof, global_proof, and parameter/proof artifacts. Stage 2 is not reached because canonical official anchoring, stable DA, low-friction user inclusion proof, stronger publication frequency, and full business-consistent constraints are not established in this report.

## 2. Stage Decision

### Stage 1, Effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| — | `NO_CANONICAL_ANCHOR` | Stage 1 | No official exchange canonical on-chain/DA anchor established. |
| — | `HIGH_FREQUENCY_GAP` | Stage 1 | Monthly cadence does not satisfy Stage 2 frequency expectations. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-06-18T16:00:00Z` |
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `~monthly` |
| History available | `3 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Older snapshot 508399035 is also in the public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `eb8d2a4279f148fc425e11c6349e971697242b3558a2db73cbfd0b223bd836cb` | https://www.okx.com/proof-of-reserves/detail |
| `walletZip` | `fb55af0ad9bd7a2c70d72c0e42b8cb358eab412fa336d5a68e1a2c509cd719b0` | https://static.okx.com/cdn/okx/por/chain/por_csv_2026042000_V3.zip |
| `globalProofBundle` | `afd141a8ec0a5ccae1ef6b82aa4a458260c99fb11eedd7e36d334d0e34445d1c` | https://static.okx.com/cdn/okx/por/merkel/por_507918525_proof_data.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x54ef34dcf7257e7b2e1689f44fd8d3a2a22336a6b684258c8fb98b06ab264841` |
| Verification bundle root | `0xe69418e85af6d2bee11f345d9851eda8585a56485e26185dfc7dbaf330c3f6e1` |
| Artifact bundle SHA-256 | `e91708984155de3ded734a3f4412a8d43f34938f9f9d76f349307e5edab107b6` |
| Verification bundle SHA-256 | `550a0d97ce8bf17cb444495c44922c3e606b1c73d395db94775fbd33bab8e61c` |

Local bundle paths: [507918525.artifact-bundle.json](../../../artifacts/okx/507918525/bundles/507918525.artifact-bundle.json), [507918525.verification-bundle.v2.json](../../../artifacts/okx/507918525/bundles/507918525.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `FAIL` | `1.0000` | okx summary has no reserve currencies |
| `internal-consistency` | `1.0` | `PASS` | `1.0000` |  |
| `address-ownership` | `okx-1` | `PASS` | `1.0000` | verified 310353/310353 address signatures |
| `onchain-balance-hot` | `2.1` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-token` | `2.0` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-ledger` | `1.2` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `global-zk-proof` | `okx-1` | `PARTIAL` | `0.7500` | structure and summary binding OK; set OKX_ZK_STARK_VALIDATOR for cryptographic verification |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | BTC block time anchor verifier not implemented |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `507918525.verification-bundle.v2.json` in the local artifact bundle.

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

#### `global-zk-proof` (`PARTIAL`)

Summary: structure and summary binding OK; set OKX_ZK_STARK_VALIDATOR for cryptographic verification

Finding counts: PASS 6

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 507918525 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-04-19T16:00:00Z |
| Is there reserve summary or wallet_address_list? (S0-3) | `pass` | Summary and wallet_address_list present. |
| Is wallet_address_list public and on-chain verifiable? (S1-1) | `pass` | Wallet list present; on-chain replay skipped or incomplete. |
| Is wallet_ownership_proof public and batch-verifiable? (S1-2) | `pass` | address-ownership verifier PASS. |
| Is global_proof public and independently reproducible? (S1-3) | `pass` | global-zk-proof verifier PARTIAL. |
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

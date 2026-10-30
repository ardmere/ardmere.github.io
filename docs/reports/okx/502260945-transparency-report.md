# OKX 502260945 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [502260945-assessment.json](./502260945-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `okx` |
| Snapshot | `502260945` |
| Snapshot time | `2026-10-02T16:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

OKX audit 502260945 does not meet all Stage 1 public artifact thresholds in this pipeline.

Missing or failing verifiers block Stage 1 effective PoR for this snapshot.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `canonical on-chain/DA anchor` | `NO_CANONICAL_ANCHOR` | Stage 1 | No official exchange canonical on-chain/DA anchor established. |
| `weekly full PoR / daily anchor` | `HIGH_FREQUENCY_GAP` | Stage 1 | Monthly cadence does not satisfy Stage 2 frequency expectations. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-10-02T16:00:00Z` |
| Previous snapshot | `2026-06-18T16:00:00Z` |
| Observed cadence | `~55 days between snapshots` |
| History available | `4 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

This is the newest snapshot in the ardmere public evaluation set for okx.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `a5d9e493edba5c8ac6b17a81672cc4e439a1542743cf1671268e57787fb3aff5` | https://www.okx.com/proof-of-reserves/detail |
| `walletZip` | `5230bd14fe64f688c5fd480006f714e746a12b6065ac9d683ca02b89852e9e07` | https://static.okx.com/cdn/okx/por/chain/por_csv_2026100300_V4.zip |
| `globalProofBundle` | `6a6da5f8afec8d93b922cca48d04c70b10f5ca2a28944487edf66d8f2b23283d` | https://static.okx.com/cdn/okx/por/merkel/por_502260945_proof_data.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x2ad5edb3d7796238ec5f464f034aca0abeafeb331d5a52ec02db4323788c39b8` |
| Verification bundle root | `0xe0d49f7eccf3e1d5c45fc045189cc44f741228aadb93e9f6d930ed00b3b1a911` |
| Artifact bundle SHA-256 | `aca7e3b255b8f705e7c06b6d85bd390380884744e50c3c4156b0d7bce1400aac` |
| Verification bundle SHA-256 | `ba61d5b275ccff099407ea8ebc1d9060173749136184137c9aa0e0f60889df98` |

Local bundle paths: [502260945.artifact-bundle.json](../../../artifacts/okx/502260945/bundles/502260945.artifact-bundle.json), [502260945.verification-bundle.v2.json](../../../artifacts/okx/502260945/bundles/502260945.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.0` | `PASS` | `1.0000` |  |
| `address-ownership` | `okx-1` | `FAIL` | `0.9994` | 371/588673 address signatures failed |
| `onchain-balance-hot` | `2.1` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-token` | `2.0` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-ledger` | `1.2` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `global-zk-proof` | `okx-1` | `PARTIAL` | `0.7500` | structure and summary binding OK; set OKX_ZK_STARK_VALIDATOR for cryptographic verification |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | BTC block time anchor verifier not implemented |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `502260945.verification-bundle.v2.json` in the local artifact bundle.

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

#### `internal-consistency` (`PASS`)

Finding counts: WARN 12, PASS 22

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `BTC` | `custodyReserveBalances` | 6212.6345421722722054 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `ETH` | `custodyReserveBalances` | 76435.0842110844594599 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `USDT` | `custodyReserveBalances` | 140455480.2233823973766545 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `USDC` | `custodyReserveBalances` | 170948528.8094416760670415 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `XRP` | `custodyReserveBalances` | 9987386.9676800585949367 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `DOGE` | `custodyReserveBalances` | 71459950.2071206516013964 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `SOL` | `custodyReserveBalances` | 83913.4943315679429555 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `OKB` | `custodyReserveBalances` | 12900.49316739333 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `LINK` | `custodyReserveBalances` | 134254.8996063100999999 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `LTC` | `custodyReserveBalances` | 49433.5384791958234293 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `UNI` | `custodyReserveBalances` | 49994.1428860879014175 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |
| `USDG` | `custodyReserveBalances` | 835031.78834801 |  | custody balances are summary-only; public wallet CSV does not break out third-party addresses |

#### `address-ownership` (`FAIL`)

Summary: 371/588673 address signatures failed

Finding counts: FAIL 5

**FAIL**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `0x80299387783b6b8c0645a17ea72fdc3cd32e0c9367445fdbe7344742af521fc39afc1fbc3917efc3a4bd72c48814e9b0` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x8063539f96661a8b8b86222b7170ee312c0ecc523b89e3f938698f69d0f82908d3db19259e047afc5c780b12eda785bc` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x806e4c086bc3aff5b9940868fac93e11b139aa04cf749a1adf13e8c6058e861df0426488db4c04d249d26d7d55e9a650` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x8082858d8486e36a4102dffb38cd65ad34aa93b7e81ecca681004508872f4df261ecf7b9ac3e741288f9a0ccbc3a9d0b` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |
| `0x808c89165c54339155a9447c617c363a6c04fcb31716030e32d4db483ca1698f566c3b87d1f38d46e33e4742f0ef7b89` | `signature` | I am an OKX address |  | missing signature2 for dual-EOA row |

#### `global-zk-proof` (`PARTIAL`)

Summary: structure and summary binding OK; set OKX_ZK_STARK_VALIDATOR for cryptographic verification

Finding counts: PASS 6

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 502260945 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-10-02T16:00:00Z |
| Is there reserve summary or wallet_address_list? (S0-3) | `pass` | Summary and wallet_address_list present. |
| Is wallet_address_list public and on-chain verifiable? (S1-1) | `pass` | Wallet list present; on-chain replay skipped or incomplete. |
| Is wallet_ownership_proof public and batch-verifiable? (S1-2) | `fail` | address-ownership verifier FAIL. |
| Is global_proof public and independently reproducible? (S1-3) | `pass` | global-zk-proof verifier PARTIAL. |
| Is independent third-party review available with stated scope? (S1-4) | `unverifiable` | third-party-attestation verifier UNVERIFIABLE. |
| If trusted setup is required, is transcript public; otherwise is the proof system transparent setup? (S1-5) | `pass` | Plonky2 / zk-STARK uses transparent setup; no trusted setup ceremony required. |
| Is root/proof/vk canonically anchored outside exchange servers? (S2-1) | `not_applicable` | Stage 2 is not evaluated until Stage 1 blockers are resolved. |
| Is publication frequency sufficient for Stage 2 (weekly full PoR / daily anchor)? (S2-2) | `not_applicable` | Stage 2 is not evaluated until Stage 1 blockers are resolved. |

## 8. Recommendations

| Priority | Recommendation | Related Risk |
| --- | --- | --- |
| `P1` | Anchor root/proof/commitment on-chain or publish immutable DA/archive records. | `NO_CANONICAL_ANCHOR` |
| `P1` | Add daily anchoring, weekly full PoR, or event-triggered updates to reduce timing risk. | `HIGH_FREQUENCY_GAP` |

## 9. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

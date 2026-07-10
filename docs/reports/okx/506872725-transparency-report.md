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

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-06-18T16:00:00Z` |
| Previous snapshot | `2026-04-19T16:00:00Z` |
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
| `summarySnapshot` | `fd614924387e4e3a49265a6dee234b4ea3010382be7f49aa778c5c20bb8c5467` | https://www.okx.com/proof-of-reserves/detail |
| `walletZip` | `bf4eb2b7de9c88587166a4ead2e2ee1617381877875b00f2a775462bb9681009` | https://static.okx.com/cdn/okx/por/chain/por_csv_2026050700_V4.zip |
| `globalProofBundle` | `319a01e352db65d17c3a6fd32b89ef35501c50e9a6b8714608c65112c36eeaee` | https://static.okx.com/cdn/okx/por/merkel/por_506872725_proof_data.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x3167e8ed73c5dd0ea7edba9b4438a7acad31cf212b976b12064087584065358c` |
| Verification bundle root | `0x2c76911a8e581ebc96ae0c28b77b15bb020faabe871431afdd18b5344e0e8b33` |
| Artifact bundle SHA-256 | `a05d442bd4cc36c4bfc20c8c3fe64f623e350c7a43d4388b182b0dba75cfaec7` |
| Verification bundle SHA-256 | `ddd3c7a9e2a6631dbca6555b1e326e98ab1007fa97ca3be020f4f784e307de2b` |

Local bundle paths: [506872725.artifact-bundle.json](../../../artifacts/okx/506872725/bundles/506872725.artifact-bundle.json), [506872725.verification-bundle.v2.json](../../../artifacts/okx/506872725/bundles/506872725.verification-bundle.v2.json)

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
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `506872725.verification-bundle.v2.json` in the local artifact bundle.

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
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 506872725 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-05-06T16:00:00Z |
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

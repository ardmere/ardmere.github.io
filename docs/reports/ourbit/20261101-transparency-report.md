# OURBIT 20261101 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [20261101-assessment.json](./20261101-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `ourbit` |
| Snapshot | `20261101` |
| Snapshot time | `2026-10-31T16:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

Ourbit snapshot 20261101 is Stage 0: public stock_info ratios + live multi-chain wallet table; no snapshot-bound walletZip/ownership; user Merkle login-gated; hash truncated to 128 bits.

solvency-claim PASS means self-reported stock_info reserve ratios ≥ 100% only. wallet_detail is archived for transparency but updates live and lacks ownership signatures, so it is not Stage 1 wallet_address_list.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `wallet_address_list` | `UNVERIFIABLE` | Stage 0 | Public wallet_detail is live/unsigned — not a snapshot-bound wallet_address_list. |
| `merkle_root` | `NO_PUBLIC_MERKLE_ROOT` | Stage 0 | stock_info has no public period merkle root to bind user proofs. |
| `sha256_truncation_128` | `TRUNCATED_HASH` | Stage 0 | Ourbit Merkle hashes truncate SHA-256 to 32 hex chars (128-bit). |
| `userMerkleProof` | `LOGIN_GATED_USER_PROOF` | Stage 0 | User merkel_tree_ourbit.json not present in the public evaluation package (login-gated). |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-10-31T16:00:00Z` |
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `UNVERIFIABLE` |
| History available | `1 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Only one ourbit snapshot is in the ardmere public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `c2236f19005f50b1874673a1bffc4dc47c685349b753bc57dc56f118867b066f` | https://www.ourbit.com/api/assetbussiness/merkle/stock_info |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x66a5f7df5ba3af5195be6e241a5b7e73ed23b0f474e375d6d9a7507eac1d6383` |
| Verification bundle root | `0x0feec92d073d9ea393c35c823d0de26b515487321e66df3796f86068a353e747` |
| Artifact bundle SHA-256 | `74c5ddf13192b2692e254d811de564e3f94659c3359e7536556e28989b7af3c5` |
| Verification bundle SHA-256 | `659b3f19df2935b5dbb15240ddc875a257cb6e0eea90625faa5ba976e82bc138` |

Local bundle paths: [20261101.artifact-bundle.json](../../../artifacts/ourbit/20261101/bundles/20261101.artifact-bundle.json), [20261101.verification-bundle.v2.json](../../../artifacts/ourbit/20261101/bundles/20261101.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `user-merkle-proof` | `ourbit-1` | `UNVERIFIABLE` | `0.0000` | merkel_tree_ourbit.json not present; download from Ourbit PoR page while logged in and import via -user-proof |
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

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `20261101.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `user-merkle-proof` | merkel_tree_ourbit.json not present; download from Ourbit PoR page while logged in and import via -user-proof |
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
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 20261101 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-10-31T16:00:00Z |
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

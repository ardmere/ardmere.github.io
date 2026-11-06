# COINEX 20260914 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [20260914-assessment.json](./20260914-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `coinex` |
| Snapshot | `20260914` |
| Snapshot time | `2026-09-14T07:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

CoinEx snapshot 20260914 is Stage 0: public reserve ratios + period merkle root + XLSX address list; user Merkle login-gated; no ownership proofs or global ZK.

solvency-claim PASS means self-reported reserve_rate ≥ 100% only. Public wallet_proof_file is an XLSX address dump without ownership signatures, so it is not Stage 1 wallet_address_list/walletZip. Official Go OSS (coinexcom/proof-of-reserves) is ported as user-merkle-proof@coinex-1 for imported proofs.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `wallet_address_list` | `UNVERIFIABLE` | Stage 0 | Public wallet_proof_file is XLSX without ownership signatures — not Stage 1 walletZip. |
| `userMerkleProof` | `LOGIN_GATED_USER_PROOF` | Stage 0 | User merkle proof JSON not present in the public evaluation package (login-gated). |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-09-14T07:00:00Z` |
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `UNVERIFIABLE` |
| History available | `1 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Only one coinex snapshot is in the ardmere public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `14e0f6562ab0279d36133129852210ff2850157855fa50ee2a2cfd3162535799` | https://www.coinex.com/res/account/balance/merkle/reserve |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x9ad38e2c5027bcf9e35a4ae364f55db56629efe1510f9f8654201f91b410abac` |
| Verification bundle root | `0x6a7f4563e0ef2189f17b00d3b8acec49bc016dc5199da018c65b6f34a2c7dd8b` |
| Artifact bundle SHA-256 | `baaa77ce9c0b671b437b84f336da944ad68b5bd58c66bd25bd27cbd7b36c7ded` |
| Verification bundle SHA-256 | `996623c7b54a363b392375b19fd01abda539458fc04718b92648f8e60a17a01a` |

Local bundle paths: [20260914.artifact-bundle.json](../../../artifacts/coinex/20260914/bundles/20260914.artifact-bundle.json), [20260914.verification-bundle.v2.json](../../../artifacts/coinex/20260914/bundles/20260914.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `user-merkle-proof` | `coinex-1` | `UNVERIFIABLE` | `0.0000` | CoinEx user Merkle proof not present; copy audit JSON while logged in on reserve-proof and import via -user-proof |
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

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `20260914.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `user-merkle-proof` | CoinEx user Merkle proof not present; copy audit JSON while logged in on reserve-proof and import via -user-proof |
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
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 20260914 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-09-14T07:00:00Z |
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

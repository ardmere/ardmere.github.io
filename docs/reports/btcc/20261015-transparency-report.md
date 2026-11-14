# BTCC 20261015 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [20261015-assessment.json](./20261015-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `btcc` |
| Snapshot | `20261015` |
| Snapshot time | `2026-10-15T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |
| Public PoR API | `https://www.btcc.com/social-api/contract/reserve/reserveInfo` |

BTCC snapshot 20261015 is Stage 0 hard-capped: permissionless reserveInfo exposes per-coin ratios and a truncated (~16 hex) merkle root only — no wallet amounts/addresses, no OSS verifier.

solvency-claim PASS means self-reported reserve ratios ≥ 100% from https://www.btcc.com/social-api/contract/reserve/reserveInfo only. Cadence is monthly ~day-15. Truncated root cannot bind a full Merkle path offline; leaf validate and history APIs need login; Stage 1 walletZip/ownership/ZK stay unwired by design of the public disclosure.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `wallet_address_list` | `UNVERIFIABLE` | Stage 0 | BTCC publishes ratios only — no wallet addresses or reserve/liability amounts. |
| `merkle_root` | `TRUNCATED_MERKLE_ROOT` | Stage 0 | Published merkle root is truncated to ~16 hex characters; insufficient to bind a full Merkle proof offline. |
| `user_merkle_verifier` | `UNVERIFIABLE` | Stage 0 | No official BTCC OSS Merkle verifier repo found; user-merkle-proof stays stubbed. |
| `userMerkleProof` | `LOGIN_GATED_USER_PROOF` | Stage 0 | Leaf hash self-check and /reserve/validate require login; no public user proof artifact in the evaluation package. |

## Self-reported liabilities, reserves, and major-coin ratios

Label: **self-reported**. These rows repeat the exchange snapshot. USD totals sum per-coin customer liability and exchange reserve (self-custody + third-party custody) using CoinGecko prices at `snapshotTime`, the same price path as `reserveCoverageBps`. A coin with no price is left out of the USD total and listed as unpriced. `walletCoverageBps` and `reserveCoverageBps` are the separate **checked** coverage ratios. This section is not solvency advice.

| Figure | Value | Label |
| --- | --- | --- |
| Total liabilities (USD) | Unavailable (snapshot did not publish amounts) | self-reported |
| Total reserves (USD) | Unavailable (snapshot did not publish amounts) | self-reported |
| Unpriced coins | — | — |

The snapshot did not publish amounts, so USD totals are Unavailable. Major-coin rows below use the exchange-reported reserve rate when that coin is in the snapshot.

| Coin | Reserve ratio | Basis | Label |
| --- | --- | --- | --- |
| `BTC` | `1.47` | `reported` | self-reported |
| `ETH` | `1.23` | `reported` | self-reported |
| `USDT` | `1.15` | `reported` | self-reported |
| `USDC` | `1.11` | `reported` | self-reported |

Reserve ratio is exchange reserve ÷ customer liability when both amounts exist (`basis=amounts`); otherwise it is the exchange-reported rate already parsed into the coin summary (`basis=reported`, percentage ÷ 100 so 1 = 100%). Major coins are BTC, ETH, USDT, and USDC. A coin is omitted when the snapshot has no row for it, or when neither amounts nor a reported rate were published (Unavailable). On-chain schema v7 stores the same USD totals as uint128 whole dollars (`0` = unavailable) and the four ratios as uint16 basis points (`10000` = 100%, `65535` = unavailable).

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-10-15T00:00:00Z` |
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `UNVERIFIABLE` |
| History available | `1 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Only one btcc snapshot is in the ardmere public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `24feebb7d353810049d37e7ed5f3bcdc0e79c0d961763ffa8e10644ad46776db` | https://www.btcc.com/social-api/contract/reserve/reserveInfo |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x879fca032a641d9b1e5d6896804f099fd2ae5ea30b20e362e0f8f801ff5832d4` |
| Verification bundle root | `0x73d114b21143ec838f4517f6d4b825e86acbb362a1801e1d9b8ad5b2ab6b7e3b` |
| Artifact bundle SHA-256 | `792ff82bdc6b11774e2c5524ff4a782a2bccc25269031efbf4b64cd4806e93d3` |
| Verification bundle SHA-256 | `4f3c598a6d9d83c645a92855885804839e6cb3f49ef58bfde472fd9d74f87654` |

Local bundle paths: [20261015.artifact-bundle.json](../../../artifacts/btcc/20261015/bundles/20261015.artifact-bundle.json), [20261015.verification-bundle.v2.json](../../../artifacts/btcc/20261015/bundles/20261015.verification-bundle.v2.json)

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

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `20261015.verification-bundle.v2.json` in the local artifact bundle.

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
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 20261015 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-10-15T00:00:00Z |
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

# KUCOIN 20260731 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [20260731-assessment.json](./20260731-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `kucoin` |
| Snapshot | `20260731` |
| Snapshot time | `2026-07-31T15:59:59Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |
| Public PoR API | `null` |

kucoin snapshot 20260731 remains Stage 0 in the ardmere public evaluation set.

Public artifacts do not support independently reproducible global solvency verification for this snapshot.

No documented public PoR API. Snapshot data depends on web pages and/or press materials; page scraping can break when resource paths migrate or the front end is a SPA that ships an HTML shell without usable totals.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `wallet_address_list` | `UNVERIFIABLE` | Stage 0 | No public wallet_address_list. |

## Self-reported liabilities, reserves, and major-coin ratios

Label: **self-reported**. These rows repeat the exchange snapshot. USD totals sum per-coin customer liability and exchange reserve (self-custody + third-party custody) using CoinGecko prices at `snapshotTime`, the same price path as `reserveCoverageBps`. A coin with no price is left out of the USD total and listed as unpriced. `walletCoverageBps` and `reserveCoverageBps` are the separate **checked** coverage ratios. This section is not solvency advice.

| Figure | Value | Label |
| --- | --- | --- |
| Total liabilities (USD) | `1702898670.56` | self-reported |
| Total reserves (USD) | `1967837143.41` | self-reported |
| Unpriced coins | — | — |

| Coin | Reserve ratio | Basis | Label |
| --- | --- | --- | --- |
| `BTC` | `1.11616252` | `amounts` | self-reported |
| `ETH` | `1.20582592` | `amounts` | self-reported |
| `USDT` | `1.16066354` | `amounts` | self-reported |
| `USDC` | `1.21645985` | `amounts` | self-reported |

Reserve ratio is exchange reserve ÷ customer liability when both amounts exist (`basis=amounts`); otherwise it is the exchange-reported rate already parsed into the coin summary (`basis=reported`, percentage ÷ 100 so 1 = 100%). Major coins are BTC, ETH, USDT, and USDC. A coin is omitted when the snapshot has no row for it, or when neither amounts nor a reported rate were published (Unavailable). On-chain schema v7 stores the same USD totals as uint128 whole dollars (`0` = unavailable) and the four ratios as uint16 basis points (`10000` = 100%, `65535` = unavailable).

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-10-31T15:59:59Z` |
| Previous snapshot | `2026-06-30T15:59:59Z` |
| Observed cadence | `~41 days between snapshots` |
| History available | `4 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Newer snapshot 20260930 is also in the public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `f927135971120a645e20c66810080de8348e2f733c09d0dc4bb958329108c5f5` | https://www.kucoin.com/proof-of-reserves |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x9525e71eb94b4f9de0e90c94a16b207b6d39010af985570522f7ed6597133f1b` |
| Verification bundle root | `0xb8b1865c49c557e74549fa031e8a0b300f3ec57630a40d33889fc037dd2bbd59` |
| Artifact bundle SHA-256 | `5cb3a0e75da7bcf9c465f85d593914cfcf6f86ce67d0eb2853f7f5fe4a09b9de` |
| Verification bundle SHA-256 | `549fea558129da799b871344a598ded327147196494fa1ee7b38cbb783751d36` |

Local bundle paths: [20260731.artifact-bundle.json](../../../artifacts/kucoin/20260731/bundles/20260731.artifact-bundle.json), [20260731.verification-bundle.v2.json](../../../artifacts/kucoin/20260731/bundles/20260731.verification-bundle.v2.json)

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

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `20260731.verification-bundle.v2.json` in the local artifact bundle.

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
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 20260731 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-07-31T15:59:59Z |
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
| `P1` | Publish a documented, permissionless PoR HTTP API so third parties need not scrape marketing pages or SPA shells. | `NO_PUBLIC_POR_API` |

## 9. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

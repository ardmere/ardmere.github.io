# BITGET 20260929 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [20260929-assessment.json](./20260929-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `bitget` |
| Snapshot | `20260929` |
| Snapshot time | `2026-09-29T09:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

bitget snapshot 20260929 remains Stage 0 in the ardmere public evaluation set.

Public artifacts do not support independently reproducible global solvency verification for this snapshot.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `wallet_address_list` | `UNVERIFIABLE` | Stage 0 | No public wallet_address_list. |

## Self-reported liabilities, reserves, and major-coin ratios

Label: **self-reported**. These rows repeat the exchange snapshot. USD totals sum per-coin customer liability and exchange reserve (self-custody + third-party custody) using CoinGecko prices at `snapshotTime`, the same price path as `reserveCoverageBps`. A coin with no price is left out of the USD total and listed as unpriced. `walletCoverageBps` and `reserveCoverageBps` are the separate **checked** coverage ratios. This section is not solvency advice.

| Figure | Value | Label |
| --- | --- | --- |
| Total liabilities (USD) | `3851528557.77` | self-reported |
| Total reserves (USD) | `4884017106.24` | self-reported |
| Unpriced coins | `HYPE`, `ONDO`, `PI`, `TAO`, `USDGO`, `XAUT` | — |

| Coin | Reserve ratio | Basis | Label |
| --- | --- | --- | --- |
| `BTC` | `1.41729977` | `amounts` | self-reported |
| `ETH` | `1.09840966` | `amounts` | self-reported |
| `USDT` | `1.07451562` | `amounts` | self-reported |
| `USDC` | `1.54001304` | `amounts` | self-reported |

Reserve ratio is exchange reserve ÷ customer liability when both amounts exist (`basis=amounts`); otherwise it is the exchange-reported rate already parsed into the coin summary (`basis=reported`, percentage ÷ 100 so 1 = 100%). Major coins are BTC, ETH, USDT, and USDC. A coin is omitted when the snapshot has no row for it. On-chain schema v7 stores the same USD totals as uint128 whole dollars (`0` = unavailable) and the four ratios as uint16 basis points (`10000` = 100%, `65535` = unavailable).

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-10-21T04:00:00Z` |
| Previous snapshot | `2026-08-20T14:00:00Z` |
| Observed cadence | `~37 days between snapshots` |
| History available | `5 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Newer snapshot 20261021 is also in the public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `9659dd0b143114984d045261826c6b9a12207432ce2d8c25e48f4175eedad684` | https://www.bitget.com/proof-of-reserves |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0xdda9edc69fbfd5fc984e8adbb16667b4e1dffd1c3bb26165b3830fa5e0e672e5` |
| Verification bundle root | `0xf08f1c5c714656413206cfd3fa164f60010f2e32fd2ab85e801acef09153d6de` |
| Artifact bundle SHA-256 | `83ce431104703020d282b3b3934a4a3f35bb06d9877564d3dd33fb51aef8d07c` |
| Verification bundle SHA-256 | `70aa0f030f65a0c0b41c7359e766e46a68b513801faf3dbfeb6659100c86e3b2` |

Local bundle paths: [20260929.artifact-bundle.json](../../../artifacts/bitget/20260929/bundles/20260929.artifact-bundle.json), [20260929.verification-bundle.v2.json](../../../artifacts/bitget/20260929/bundles/20260929.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `user-merkle-proof` | `bitget-1` | `UNVERIFIABLE` | `0.0000` | merkel_tree_bg.json not present; download from Bitget while logged in and import via -user-proof |
| `internal-consistency` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | `0` | `UNVERIFIABLE` | `0.0000` | No public HotCold wallet address list |
| `onchain-balance-token` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address list for ERC20/BEP20 audit |
| `onchain-balance-ledger` | `0` | `UNVERIFIABLE` | `0.0000` | Verifier not available for this exchange snapshot |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Bitget uses pure Merkle Tree with 64-bit truncated SHA-256 (no ZK) |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No clearly identified independent third-party PoR attestation in machine-readable form |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `20260929.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `user-merkle-proof` | merkel_tree_bg.json not present; download from Bitget while logged in and import via -user-proof |
| `internal-consistency` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | No public HotCold wallet address list |
| `onchain-balance-token` | No public wallet address list for ERC20/BEP20 audit |
| `onchain-balance-ledger` | Verifier not available for this exchange snapshot |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Bitget uses pure Merkle Tree with 64-bit truncated SHA-256 (no ZK) |
| `third-party-attestation` | No clearly identified independent third-party PoR attestation in machine-readable form |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

No additional row-level `FAIL` or `WARN` findings beyond the summary table in §5.

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 20260929 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-09-29T09:00:00Z |
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

## Capture notes / gaps (2026-10-31)

Browser public bill/proof APIs; user `merkel_tree_bg.json` login-gated; no public wallet ZIP.

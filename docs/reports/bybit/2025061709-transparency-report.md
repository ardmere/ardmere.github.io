# BYBIT 2025061709 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [2025061709-assessment.json](./2025061709-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `bybit` |
| Snapshot | `2025061709` |
| Snapshot time | `2025-06-17T09:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |
| Public PoR API | `null` |

bybit snapshot 2025061709 remains Stage 0 in the ardmere public evaluation set.

Public artifacts do not support independently reproducible global solvency verification for this snapshot.

No documented public PoR API. Snapshot data depends on web pages and/or press materials; page scraping can break when resource paths migrate or the front end is a SPA that ships an HTML shell without usable totals.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| — | `UNVERIFIABLE` | Stage 0 | No public wallet_address_list. |

## Self-reported liabilities, reserves, and major-coin ratios

Label: **self-reported**. These rows repeat the exchange snapshot. USD totals sum per-coin customer liability and exchange reserve (self-custody + third-party custody) using CoinGecko prices at `snapshotTime`, the same price path as `reserveCoverageBps`. A coin with no price is left out of the USD total and listed as unpriced. `walletCoverageBps` and `reserveCoverageBps` are the separate **checked** coverage ratios. This section is not solvency advice.

| Figure | Value | Label |
| --- | --- | --- |
| Total liabilities (USD) | Unavailable | self-reported |
| Total reserves (USD) | Unavailable | self-reported |
| Unpriced coins | `BTC (price error: coingecko HTTP 401: {"error":{"status":{"timestamp":"2026-11-12T15:23:50.655+00:00","error_code":10012,"error_message":"Your request exceeds the allowed time range. Public API users are limited to querying historical dat…)`, `ETH (price error: coingecko HTTP 401: {"error":{"status":{"timestamp":"2026-11-12T15:23:50.798+00:00","error_code":10012,"error_message":"Your request exceeds the allowed time range. Public API users are limited to querying historical dat…)`, `USDT (price error: coingecko HTTP 401: {"error":{"status":{"timestamp":"2026-11-12T15:24:51.005+00:00","error_code":10012,"error_message":"Your request exceeds the allowed time range. Public API users are limited to querying historical dat…)` | — |

Amount-bearing coins were all unpriced, so USD totals are Unavailable.

| Coin | Reserve ratio | Basis | Label |
| --- | --- | --- | --- |
| `BTC` | `1.05263158` | `amounts` | self-reported |
| `ETH` | `1.01010101` | `amounts` | self-reported |
| `USDT` | `1.00502513` | `amounts` | self-reported |

Reserve ratio is exchange reserve ÷ customer liability when both amounts exist (`basis=amounts`); otherwise it is the exchange-reported rate already parsed into the coin summary (`basis=reported`, percentage ÷ 100 so 1 = 100%). Major coins are BTC, ETH, USDT, and USDC. A coin is omitted when the snapshot has no row for it, or when neither amounts nor a reported rate were published (Unavailable). On-chain schema v7 stores the same USD totals as uint128 whole dollars (`0` = unavailable) and the four ratios as uint16 basis points (`10000` = 100%, `65535` = unavailable).

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-10-27T06:00:00Z` |
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `~51 days between snapshots` |
| History available | `5 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Newer snapshot 2026052709 is also in the public evaluation set. Observed cadence averages 3 of 4 adjacent gaps; excluded as an older span superseded by the current rhythm: `2025061709` → `2026052709` (344 days). Retained gaps are each at most 3x the 34-day median gap.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `110282f57984ebd218a3ba0ee510b8e08bf6a2bb24de3cad9c0af22d528d422a` | https://www.bybit.com/x-api/por/public/v1/reserve-ratio/latest |
| `userMerkleProof` | `d095eed7e7dde2625c7311d24dbe307e9b12a25574e29ac4e1c073f14fdc2ab8` | file://./fixtures/bybit/mock_user_merkle_tree_path_40_v5.json |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x1fe7bf6ea18a09a33bb05c0f1de1504a0e5e1e02736d74ad1c76015c3bf3783e` |
| Verification bundle root | `0x4e2d01591b143c0773360b375ace38350bf60f167fb307b012cb7163d4918755` |
| Artifact bundle SHA-256 | `62e189e04471a4d274747aeae68a2d3f33f9f224e4e60632aa849a373c960b99` |
| Verification bundle SHA-256 | `88e9ce99bb083de5c10921ab119a1e99c0a73520215b9fec66b40d56b193f900` |

Local bundle paths: [2025061709.artifact-bundle.json](../../../artifacts/bybit/2025061709/bundles/2025061709.artifact-bundle.json), [2025061709.verification-bundle.v2.json](../../../artifacts/bybit/2025061709/bundles/2025061709.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `user-merkle-proof` | `bybit-1` | `PASS` | `1.0000` | user Merkle path valid (Gen 1.5 SHA-256); proves inclusion only for supplied account |
| `internal-consistency` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | `0` | `UNVERIFIABLE` | `0.0000` | No public HotCold wallet address list |
| `onchain-balance-token` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address list for ERC20/BEP20 audit |
| `onchain-balance-ledger` | `0` | `UNVERIFIABLE` | `0.0000` | Verifier not available for this exchange snapshot |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Bybit uses pure Merkle Tree (no ZK); user myProof.json proves inclusion only |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | Hacken monthly PoR reports are public but not machine-readable in this pipeline |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `2025061709.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `internal-consistency` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | No public HotCold wallet address list |
| `onchain-balance-token` | No public wallet address list for ERC20/BEP20 audit |
| `onchain-balance-ledger` | Verifier not available for this exchange snapshot |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Bybit uses pure Merkle Tree (no ZK); user myProof.json proves inclusion only |
| `third-party-attestation` | Hacken monthly PoR reports are public but not machine-readable in this pipeline |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

No additional row-level `FAIL` or `WARN` findings beyond the summary table in §5.

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 2025061709 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2025-06-17T09:00:00Z |
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

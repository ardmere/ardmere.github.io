# HTX 20261101 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [20261101-assessment.json](./20261101-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `htx` |
| Snapshot | `20261101` |
| Snapshot time | `2026-11-01T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

htx snapshot 20261101 remains Stage 0 in the ardmere public evaluation set.

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
| Total liabilities (USD) | `5182100386.54` | self-reported |
| Total reserves (USD) | `5305874609.16` | self-reported |
| Unpriced coins | `HTX`, `U`, `USDD`, `USDS` | — |

| Coin | Reserve ratio | Basis | Label |
| --- | --- | --- | --- |
| `BTC` | `1.0218904` | `amounts` | self-reported |
| `ETH` | `1.0144773` | `amounts` | self-reported |
| `USDT` | `0.35020709` | `amounts` | self-reported |
| `USDC` | `13.47805109` | `amounts` | self-reported |

Reserve ratio is exchange reserve ÷ customer liability when both amounts exist (`basis=amounts`); otherwise it is the exchange-reported rate already parsed into the coin summary (`basis=reported`, percentage ÷ 100 so 1 = 100%). Major coins are BTC, ETH, USDT, and USDC. A coin is omitted when the snapshot has no row for it. On-chain schema v7 stores the same USD totals as uint128 whole dollars (`0` = unavailable) and the four ratios as uint16 basis points (`10000` = 100%, `65535` = unavailable).

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-11-01T00:00:00Z` |
| Previous snapshot | `2026-10-01T00:00:00Z` |
| Observed cadence | `~monthly` |
| History available | `5 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

This is the newest snapshot in the ardmere public evaluation set for htx. Observed cadence averages 3 of 4 adjacent gaps; excluded as an older span superseded by the current rhythm: `20230910` → `20260801` (1055 days). Retained gaps are each at most 3x the 31-day median gap.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `bc36e587f898c596926ca101316481ed78e38e8943d90bc96c9b888bec5021be` | https://www.htx.com/zh-cn/finance/merkle/ |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x7537ef23f9f67f4c0b6748eb122ed0655cc197459e75204513e594edbaa101f6` |
| Verification bundle root | `0x7c9d4c46ef5f43cb0e48f2748d650db0d478c884616d9c7d2b60f699f25ab4ca` |
| Artifact bundle SHA-256 | `ff1f8eb063c280cbf2d58ed50dc763f98f985f048a83293e22b9e454f83ea143` |
| Verification bundle SHA-256 | `0a76707b4d3db1e94b54553257c29fdb5ae49bc040857b79f468c4609b5d43e9` |

Local bundle paths: [20261101.artifact-bundle.json](../../../artifacts/htx/20261101/bundles/20261101.artifact-bundle.json), [20261101.verification-bundle.v2.json](../../../artifacts/htx/20261101/bundles/20261101.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `FAIL` | `1.0000` | gate summary missing totalReserveRate |
| `global-zk-proof` | `htx-1` | `UNVERIFIABLE` | `0.0000` | HTX public-data zk bundle not present; download from PoR Reports or GitHub release public-data.zip |
| `internal-consistency` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | `0` | `UNVERIFIABLE` | `0.0000` | No public HotCold wallet address list |
| `onchain-balance-token` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address list for ERC20/BEP20 audit |
| `onchain-balance-ledger` | `0` | `UNVERIFIABLE` | `0.0000` | Verifier not available for this exchange snapshot |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `20261101.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `global-zk-proof` | HTX public-data zk bundle not present; download from PoR Reports or GitHub release public-data.zip |
| `internal-consistency` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | No public HotCold wallet address list |
| `onchain-balance-token` | No public wallet address list for ERC20/BEP20 audit |
| `onchain-balance-ledger` | Verifier not available for this exchange snapshot |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `solvency-claim` (`FAIL`)

Summary: gate summary missing totalReserveRate

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 20261101 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-11-01T00:00:00Z |
| Is there reserve summary or wallet_address_list? (S0-3) | `pass` | Summary present; wallet_address_list absent. |
| Is wallet_address_list public and on-chain verifiable? (S1-1) | `fail` | No public wallet_address_list artifact. |
| Is wallet_ownership_proof public and batch-verifiable? (S1-2) | `unverifiable` | address-ownership verifier UNVERIFIABLE. |
| Is global_proof public and independently reproducible? (S1-3) | `unverifiable` | global-zk-proof verifier UNVERIFIABLE. |
| Is independent third-party review available with stated scope? (S1-4) | `unverifiable` | third-party-attestation verifier UNVERIFIABLE. |
| If trusted setup is required, is transcript public; otherwise is the proof system transparent setup? (S1-5) | `unverifiable` | No public trusted setup transcript observed for zk-SNARK global proof. |
| Is root/proof/vk canonically anchored outside exchange servers? (S2-1) | `not_applicable` | Stage 2 is not evaluated until Stage 1 blockers are resolved. |
| Is publication frequency sufficient for Stage 2 (weekly full PoR / daily anchor)? (S2-2) | `not_applicable` | Stage 2 is not evaluated until Stage 1 blockers are resolved. |

## 8. Recommendations

| Priority | Recommendation | Related Risk |
| --- | --- | --- |
| `P0` | Publish missing public artifacts required for independent reproduction. | `UNVERIFIABLE` |
| `P0` | Publish trusted setup transcript or migrate to a transparent-setup proof system. | `OPAQUE_TRUSTED_SETUP` |

## 9. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

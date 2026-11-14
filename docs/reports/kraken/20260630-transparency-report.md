# KRAKEN 20260630 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [20260630-assessment.json](./20260630-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `kraken` |
| Snapshot | `20260630` |
| Snapshot time | `2026-06-30T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |
| Public PoR API | `https://widget.pro.ledgerlens.io/api/kraken-platform-spot-futures` |

Kraken snapshot 20260630 is Stage 0 attestation-model: quarterly accountant Merkle (LedgerLens / The Network Firm) plus CMS reserve ratios — no permissionless wallet ZIP or machine-verified attestation PDF.

solvency-claim PASS means self-reported CMS reserve ratios ≥ 100% (capped stables treated as published 105%+ floor). LedgerLens exposes a full combined Merkle root and names The Network Firm LLP, but user leaf/path remains login-gated, wallet ownership evidence stays inside the accountant workflow, and third-party-attestation is not machine-parsed — Stage 1 stays out of reach by disclosure design.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `wallet_address_list` | `UNVERIFIABLE` | Stage 0 | Kraken publishes ratios + accountant Merkle root only — no permissionless wallet addresses or reserve/liability amounts. |
| `third_party_attestation` | `UNVERIFIABLE` | Stage 0 | Accountant is named (The Network Firm LLP / LedgerLens) but the CPA attestation report is not ingested as a machine-verifiable public artifact. |
| `publication_frequency` | `HIGH_FREQUENCY_GAP` | Stage 0 | Quarterly accountant reviews are too infrequent for Stage 2 frequency expectations and are a weak monitor-primary path. |
| `userMerkleProof` | `LOGIN_GATED_USER_PROOF` | Stage 0 | User Merkle Leaf ID / path self-check requires Kraken Pro login (or personalized accountant tool input); no public user proof in the evaluation package. |

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
| `BTC` | `1.029` | `reported` | self-reported |
| `ETH` | `1.005` | `reported` | self-reported |
| `USDT` | `1.05` | `reported` | self-reported |
| `USDC` | `1.05` | `reported` | self-reported |

Reserve ratio is exchange reserve ÷ customer liability when both amounts exist (`basis=amounts`); otherwise it is the exchange-reported rate already parsed into the coin summary (`basis=reported`, percentage ÷ 100 so 1 = 100%). Major coins are BTC, ETH, USDT, and USDC. A coin is omitted when the snapshot has no row for it, or when neither amounts nor a reported rate were published (Unavailable). On-chain schema v7 stores the same USD totals as uint128 whole dollars (`0` = unavailable) and the four ratios as uint16 basis points (`10000` = 100%, `65535` = unavailable).

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-09-30T00:00:00Z` |
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `~92 days between snapshots` |
| History available | `2 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Newer snapshot 20260930 is also in the public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `88c9c21431a84732adf006b0a7e304092b95e9bc39874e2f309dffe005559979` | https://widget.pro.ledgerlens.io/api/kraken-platform-spot-futures |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x1e12dd3c6cabad5d6d49a1b991563d5a454ddb393ba8d2e265d31efb6671dc03` |
| Verification bundle root | `0x4eecc81d4694314dc1ce971e460bd02ac617e65e89170690c6c788ece84e892b` |
| Artifact bundle SHA-256 | `c9ffe395f0ae35672f5629cbe83b293f81f7a3b09e2c8e74b94b31f9a4e03a83` |
| Verification bundle SHA-256 | `29f2039e0ca11565ce5564b4dd30103b3dbe71675a1e617e177a1ccd33b84248` |

Local bundle paths: [20260630.artifact-bundle.json](../../../artifacts/kraken/20260630/bundles/20260630.artifact-bundle.json), [20260630.verification-bundle.v2.json](../../../artifacts/kraken/20260630/bundles/20260630.verification-bundle.v2.json)

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

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `20260630.verification-bundle.v2.json` in the local artifact bundle.

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
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 20260630 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-06-30T00:00:00Z |
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
| `P1` | Add daily anchoring, weekly full PoR, or event-triggered updates to reduce timing risk. | `HIGH_FREQUENCY_GAP` |

## 9. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

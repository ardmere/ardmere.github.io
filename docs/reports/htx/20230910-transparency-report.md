# HTX 20230910 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [20230910-assessment.json](./20230910-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `htx` |
| Snapshot | `20230910` |
| Snapshot time | `2023-09-10T18:10:33Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

htx snapshot 20230910 remains Stage 0 in the ardmere public evaluation set.

Public artifacts do not support independently reproducible global solvency verification for this snapshot.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| — | `UNVERIFIABLE` | Stage 0 | No public wallet_address_list. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2023-09-10T18:10:33Z` |
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `UNVERIFIABLE` |
| History available | `1 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Only one htx snapshot is in the ardmere public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `93feed459cfc833bd5c29fc04fdccc25a2a393d2938cf6037f1ef0789d94d47d` | https://www.htx.com/zh-cn/finance/merkle/ |
| `globalProofBundle` | `16954049006272c0cfead838b8098fff3b5a73fddbd1745ed1953113f851a182` | file://./fixtures/htx/public-data.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x9db05b45ef048abc974606d9945ad987d31c2df0e10ab8e1b960ae4e24d827f7` |
| Verification bundle root | `0x5758ec549f7d81b3b9a0ad3aeec0afda5670366b101c94e1967ceecdacebf271` |
| Artifact bundle SHA-256 | `ff5988aff0308f0746fd140293d468da53f44f44a526578fb65b860a2c7a39e8` |
| Verification bundle SHA-256 | `9b2ab6b4ab6f36731b3435fad68bfefe24be145d896d2c7aab8f9d20e96af5de` |

Local bundle paths: [20230910.artifact-bundle.json](../../../artifacts/htx/20230910/bundles/20230910.artifact-bundle.json), [20230910.verification-bundle.v2.json](../../../artifacts/htx/20230910/bundles/20230910.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `UNVERIFIABLE` | `0.0000` | zk-derived summary has no public reserve ratios; import browser-captured ratio JSON via -summary-path |
| `global-zk-proof` | `htx-1` | `PARTIAL` | `0.7000` | structure and merkle bind OK; set HTX_ZK_VERIFIER to zkverifiermac path for Groth16 cryptographic verification |
| `internal-consistency` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | `0` | `UNVERIFIABLE` | `0.0000` | No public HotCold wallet address list |
| `onchain-balance-token` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address list for ERC20/BEP20 audit |
| `onchain-balance-ledger` | `0` | `UNVERIFIABLE` | `0.0000` | Verifier not available for this exchange snapshot |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `20230910.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `solvency-claim` | zk-derived summary has no public reserve ratios; import browser-captured ratio JSON via -summary-path |
| `internal-consistency` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | No public HotCold wallet address list |
| `onchain-balance-token` | No public wallet address list for ERC20/BEP20 audit |
| `onchain-balance-ledger` | Verifier not available for this exchange snapshot |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `global-zk-proof` (`PARTIAL`)

Summary: structure and merkle bind OK; set HTX_ZK_VERIFIER to zkverifiermac path for Groth16 cryptographic verification

Finding counts: PASS 3

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 20230910 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2023-09-10T18:10:33Z |
| Is there reserve summary or wallet_address_list? (S0-3) | `pass` | Summary present; wallet_address_list absent. |
| Is wallet_address_list public and on-chain verifiable? (S1-1) | `fail` | No public wallet_address_list artifact. |
| Is wallet_ownership_proof public and batch-verifiable? (S1-2) | `unverifiable` | address-ownership verifier UNVERIFIABLE. |
| Is global_proof public and independently reproducible? (S1-3) | `pass` | global-zk-proof verifier PARTIAL. |
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

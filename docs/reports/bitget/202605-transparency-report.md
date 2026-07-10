# BITGET 202605 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [202605-assessment.json](./202605-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `bitget` |
| Snapshot | `202605` |
| Snapshot time | `2026-05-27T10:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

bitget snapshot 202605 remains Stage 0 in the ardmere public evaluation set.

Public artifacts do not support independently reproducible global solvency verification for this snapshot.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| — | `UNVERIFIABLE` | Stage 0 | No public wallet_address_list. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-05-27T10:00:00Z` |
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `UNVERIFIABLE` |
| History available | `1 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Only one bitget snapshot is in the ardmere public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `b17df25a7711b4e7c6c4186ad8fa4b36d3c89806f3fc9371b47d6030d806ed00` | https://www.bitget.com/proof-of-reserves |
| `userMerkleProof` | `87cbcf7dee34b834e379f1d3e5ace4248fe19ad728741e6c59eac2af9200a48f` | file://./fixtures/bitget/merkel_tree_bg.json |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x482e0a28d40c3f480c1296d4417f8b17ff7d4fca878bbf65641aea1c63d2a05a` |
| Verification bundle root | `0x0015676e2da406f1b0c30fedc72a8aaad22770e4b961b780bc1537e194d89de4` |
| Artifact bundle SHA-256 | `dd588088ef6d12054136407a5b5b831528851e900ac41dcca22477415b543bf9` |
| Verification bundle SHA-256 | `8950b4c8e076463b56b387f2167e442cd9028155e607338bd5e51fa5c5b4c397` |

Local bundle paths: [202605.artifact-bundle.json](../../../artifacts/bitget/202605/bundles/202605.artifact-bundle.json), [202605.verification-bundle.v2.json](../../../artifacts/bitget/202605/bundles/202605.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `user-merkle-proof` | `bitget-1` | `PASS` | `1.0000` | Bitget user Merkle path valid under weak 64-bit truncated-hash scheme; proves inclusion only for supplied account |
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

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `202605.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `internal-consistency` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | No public HotCold wallet address list |
| `onchain-balance-token` | No public wallet address list for ERC20/BEP20 audit |
| `onchain-balance-ledger` | Verifier not available for this exchange snapshot |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Bitget uses pure Merkle Tree with 64-bit truncated SHA-256 (no ZK) |
| `third-party-attestation` | No clearly identified independent third-party PoR attestation in machine-readable form |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `user-merkle-proof` (`PASS`)

Summary: Bitget user Merkle path valid under weak 64-bit truncated-hash scheme; proves inclusion only for supplied account

Finding counts: WARN 1, PASS 2

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `hashSecurity` | `sha256Truncation` | 64-bit hash output | collision resistance about 32 bits | Bitget truncates SHA-256 to 16 hex chars; Merkle PASS is weak evidence and should not be treated like full SHA-256 |

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 202605 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-05-27T10:00:00Z |
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

# GATEIO 20260316 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [20260316-assessment.json](./20260316-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `gateio` |
| Snapshot | `20260316` |
| Snapshot time | `2026-03-16T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

gateio snapshot 20260316 remains Stage 0 in the ardmere public evaluation set.

Public artifacts do not support independently reproducible global solvency verification for this snapshot.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| — | `UNVERIFIABLE` | Stage 0 | No public wallet_address_list. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-03-16T00:00:00Z` |
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `UNVERIFIABLE` |
| History available | `1 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Only one gateio snapshot is in the ardmere public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `ed9a64cb48f33df8c524dea226c06d4360e79a23f0dd4b5aa6d470aeb2b2de71` | https://www.gate.com/api/web/v1/proof-of-reserves/getProofOfReservesInfo |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0xef120b448c411c8ab6f03cdd705cc973bb642aa849a8d916ea3af5e0038d33eb` |
| Verification bundle root | `0x9cdd55018571ba1f1c09c4960b8bd56a55199695278c24a94bedb875a44f8b75` |
| Artifact bundle SHA-256 | `8e10dfcb1e5d56f07ea04884df60ac787b429a1154b20efb14b1c4760c14fbb5` |
| Verification bundle SHA-256 | `fa3bd1667afabade5e43ba0d89101028e83f72ba1e0a78f436c6ab1eea55a990` |

Local bundle paths: [20260316.artifact-bundle.json](../../../artifacts/gateio/20260316/bundles/20260316.artifact-bundle.json), [20260316.verification-bundle.v2.json](../../../artifacts/gateio/20260316/bundles/20260316.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.0000` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | `0` | `UNVERIFIABLE` | `0.0000` | No public HotCold wallet address list |
| `onchain-balance-token` | `0` | `UNVERIFIABLE` | `0.0000` | No public wallet address list for ERC20/BEP20 audit |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `gateio-0` | `UNVERIFIABLE` | `0.0000` | Global zkmerkle_cex tar.gz requires login; not available from public API |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `20260316.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `internal-consistency` | No public wallet address CSV; cannot reconcile summary vs address aggregates |
| `btc-anchor` | Summary does not declare a BTC block height time anchor |
| `onchain-balance-hot` | No public HotCold wallet address list |
| `onchain-balance-token` | No public wallet address list for ERC20/BEP20 audit |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Global zkmerkle_cex tar.gz requires login; not available from public API |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

No additional row-level `FAIL` or `WARN` findings beyond the summary table in §5.

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 20260316 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-03-16T00:00:00Z |
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

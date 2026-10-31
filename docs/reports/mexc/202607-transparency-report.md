# MEXC 202607 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [202607-assessment.json](./202607-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `mexc` |
| Snapshot | `202607` |
| Snapshot time | `2026-07-15T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

MEXC snapshot 202607 is a weak Stage 0 (press-release / import summary): self-reported ratios only; no public period merkle root in the archived package; user Merkle remains login-gated.

solvency-claim PASS means self-reported reserve ratios ≥ 100% only — not independently reproduced reserves. Datacenter fetches of www.mexc.com are WAF-blocked; Hacken PDF attestation is not machine-parsed here.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `wallet_address_list` | `UNVERIFIABLE` | Stage 0 | No public wallet_address_list. |
| `merkle_root` | `NO_PUBLIC_MERKLE_ROOT` | Stage 0 | Press-release summary has no public period merkle root to bind user proofs. |
| `por_dashboard_api` | `WAF_GATED` | Stage 0 | www.mexc.com PoR JSON is WAF-gated from datacenters; summary imported from press materials. |
| `userMerkleProof` | `LOGIN_GATED_USER_PROOF` | Stage 0 | User myProof.json not present in the public evaluation package (login-gated). |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-10-15T00:00:00Z` |
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `~monthly` |
| History available | `4 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Older snapshot 202610 is also in the public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `93f631b08db2d26e5c7c9fa1fc863b070a372fbd9908ba6fd5d69769c940aaba` | https://www.mexc.com/support/article/how-to-verify-user-assets-on-mexc-17827791528836 |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x9c3413cf09eaf6afe44141d4352138a915a97335bb8f3945090a9b5c3a6df512` |
| Verification bundle root | `0xbd970786476407076940c9110f7f5a47626f1731c71352f16843978290334bd0` |
| Artifact bundle SHA-256 | `b6e0308a64e456619e857dd871acfdd2ef467cb78fd424d62504fbcdd3c64141` |
| Verification bundle SHA-256 | `580fee926fc578396d1fa22f49096feba0d1d75e195545afbc04fb348aafdf74` |

Local bundle paths: [202607.artifact-bundle.json](../../../artifacts/mexc/202607/bundles/202607.artifact-bundle.json), [202607.verification-bundle.v2.json](../../../artifacts/mexc/202607/bundles/202607.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` |  |
| `user-merkle-proof` | `mexc-1` | `UNVERIFIABLE` | `0.0000` | myProof.json not present; download from MEXC PoR page while logged in and import via -user-proof |
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

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `202607.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `user-merkle-proof` | myProof.json not present; download from MEXC PoR page while logged in and import via -user-proof |
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
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 202607 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-07-15T00:00:00Z |
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

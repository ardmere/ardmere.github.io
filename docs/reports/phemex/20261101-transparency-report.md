# PHEMEX 20261101 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [20261101-assessment.json](./20261101-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `phemex` |
| Snapshot | `20261101` |
| Snapshot time | `2026-11-01T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

Phemex snapshot 20261101 is Stage 0: permissionless porVersions balances → computed reserve ratios; monthly Merkle inclusion and cold-wallet addresses appear on the public PoR UI, but merkle root / wallet-addresses APIs are CloudFront-gated and user proofs stay login-gated.

solvency-claim PASS means chainBalance/cexBalance ratios from https://por.phemex.com/porVersions only. UI shows Part of Our Cold Wallets and a Merkle self-check, but there is no permissionless wallet zip, no public period merkle root in the archived package, and no official OSS verifier — Stage 1 walletZip/ownership/ZK stay unwired.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| `wallet_address_list` | `UNVERIFIABLE` | Stage 0 | Cold-wallet addresses are UI-visible only; api10 wallet-addresses is CloudFront/auth gated and there is no downloadable wallet zip. |
| `merkle_root` | `NO_PUBLIC_MERKLE_ROOT` | Stage 0 | Period merkle root is not in public porVersions; api10 merkle-tree/root is CloudFront/auth gated from datacenters. |
| `user_merkle_verifier` | `UNVERIFIABLE` | Stage 0 | No official Phemex OSS Merkle verifier repo found; user-merkle-proof stays stubbed. |
| `userMerkleProof` | `LOGIN_GATED_USER_PROOF` | Stage 0 | Hashed-client-ID Merkle inclusion tool requires login; no public user proof artifact in the evaluation package. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-11-01T00:00:00Z` |
| Previous snapshot | `UNVERIFIABLE` |
| Observed cadence | `UNVERIFIABLE` |
| History available | `UNVERIFIABLE` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

No sibling assessments found in the ardmere public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summarySnapshot` | `aeea5ca2a488d6ae771e92be363a5297e9d7373a9ece7b476a9b84b21b2d9d92` | https://por.phemex.com/porVersions |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x636e82acb98b17ee07be5a8a18d9713213c0cb8a64ae38142dfdc1da139ca027` |
| Verification bundle root | `0x2f30ae43343ecfdd6ca6d356991f9aee5aef7617c3f95cc162421349c9870a71` |
| Artifact bundle SHA-256 | `8b1d9a023c2c480c98abf45c1e14810305ea4f845a4c9e08be63f445cb997d24` |
| Verification bundle SHA-256 | `6b4a454c2d6c2dc05de8dcc7f7bdaa5c875e7232c1cea4af7a5b3254b5f3bb31` |

Local bundle paths: [20261101.artifact-bundle.json](../../../artifacts/phemex/20261101/bundles/20261101.artifact-bundle.json), [20261101.verification-bundle.v2.json](../../../artifacts/phemex/20261101/bundles/20261101.verification-bundle.v2.json)

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

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `20261101.verification-bundle.v2.json` in the local artifact bundle.

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
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot 20261101 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-11-01T00:00:00Z |
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

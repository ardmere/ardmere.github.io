# Binance PR01JUN26 PoR Transparency Report

> Methodology: [`por-transparency-framework.md`](../../por-transparency-framework.md)
> Assessment JSON: [`PR01JUN26-assessment.json`](./PR01JUN26-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `binance` |
| Snapshot | `PR01JUN26` |
| Period | `43` |
| Wallet ZIP date | `2026-06-01` |
| BAPI snapshot time | `2026-01-06T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

Binance publishes a public reserve summary and `wallet_address_list` for `PR01JUN26`, and ardmere can reproduce several reserve-side checks. The available evidence supports `Gen 2 / E2`, but the snapshot remains **Stage 0** because key Stage 1 artifacts are missing or unavailable to the public pipeline.

Under the ardmere framework, this snapshot should **not** be described as `Stage 1`, `Stage 2`, or effective PoR.

## 2. Stage Decision

### Stage 0, Not Effective PoR

Binance has minimum PoR disclosure: the BAPI summary and wallet ZIP are public, artifact integrity passes, internal consistency passes, and solvency-claim checks pass.

Stage 1 is blocked by the following gaps:

| Missing / Blocked Evidence | Risk Flag | Stage Effect | Why It Matters |
| --- | --- | --- | --- |
| `wallet_ownership_proof` | `NO_WALLET_OWNERSHIP_PROOF` | max `Stage 0` | Users must trust that Binance controls the listed wallet addresses. |
| `trusted_setup_transcript` or transparent setup evidence | `OPAQUE_TRUSTED_SETUP` | max `Stage 0` | Users must trust setup honesty and absence of toxic waste. |
| Public `global_proof` / `vk` / proof schema | `UNVERIFIABLE` | max `Stage 0` | The full global ZK proof stack is not publicly reproducible by this pipeline. |
| Machine-readable independent review | `UNVERIFIABLE` | confidence downgrade | The review scope cannot be automatically checked against this snapshot. |

Additional Stage 2 blockers remain even if Stage 1 gaps are fixed: Binance would still need official canonical anchoring, data availability, low-friction inclusion proof, permissionless replay, and public business-consistent proof constraints.

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest Snapshot | `2026-06-01T00:00:00Z` |
| Previous Snapshot | `2026-05-01T00:00:00Z` |
| Observed Cadence | `monthly` |
| History Available | `true` |
| Event-triggered Updates | `unknown` |
| Daily Root / Commitment Anchor | `no` |
| Stage Impact | Monthly PoR is a transitional baseline, but insufficient for Stage 2 without daily root/commitment anchoring, event-triggered updates, or stronger cadence. |

Binance publishes recurring monthly PoR snapshots. This report does not establish event-triggered publication rules or official daily canonical anchoring.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summary` / `bapiSnapshot` | `7898c1147f470b08404618c7f485664a4cb91740be3415d469d52f23c64f4202` | `https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot` |
| `wallet_address_list` / `walletZip` | `722466eec2376e28cfeb825eeee1a08a50b7459e640d63f78a506386de65c8b2` | `https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_20260601.zip` |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0xf452e47dd22ed63dc4a905fe79da6c6f7a6975cc0d775f50d01879e97616671f` |
| Verification bundle v2 root | `0x84e7f008460dcfa4ae4969c9a6e2a8b2585d5dcbe6a33cfead74b175782f6f42` |
| Artifact bundle SHA-256 | `ffc546892f40613f04434624c3ae14e3d5bb9f23e38bc971731c071b656ade50` |
| Verification bundle v2 SHA-256 | `b3b5a3dc2f08f198110131c787c52a666b3af5b5c3fbf4eddff4dff6e7c1e709` |

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` | BAPI snapshot and wallet ZIP hashes match archived artifacts. |
| `internal-consistency` | `1.1` | `PASS` | `1.0000` | `HotCold=1262`, `Deposit=1980326`, total `1981588`; BAPI reconciliation passed. |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` | 51 solvency-claim findings passed. |
| `onchain-balance-hot` | `2.1` | `FAIL` | `0.0396` | 1 FAIL, 4 WARN, 155 UNVERIFIABLE, 1 PASS. |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.5357` | 4 FAIL, 202 WARN, 38 UNVERIFIABLE, 1 PASS. |
| `onchain-balance-ledger` | `1.4` | `FAIL` | `0.1339` | 9 FAIL, 151 WARN, 68 UNVERIFIABLE, 1 PASS. |
| `onchain-balance-deposit` | `1.2` | `WARN` | `0.1331` | Value-weighted deposit sample warned due to RPC failures. |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0` | No public wallet ownership signatures / proofs. |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0` | `global_proof` / verifying key not publicly distributed. |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0` | No public machine-readable attestation artifact. |

The reserve-side FAIL verdicts should not be read as direct evidence of Binance insolvency. The detailed ardmere audit report attributes many findings to RPC/archive availability, block-height boundary ambiguity, chain-specific accounting semantics, staking accounting differences, and wrapped/cross-chain reconciliation limits. The Stage decision is driven primarily by missing Stage 1 artifacts.

## 6. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Public PoR summary | pass | BAPI summary is archived. |
| Public `wallet_address_list` | pass | Wallet ZIP is public and partially checked by ardmere. |
| Public `wallet_ownership_proof` | fail | Blocks Stage 1. |
| Public `global_proof` + `vk/config` | fail | Blocks Stage 1. |
| Trusted setup transparency | fail | Blocks Stage 1. |
| Low-friction inclusion proof bound to public root / `global_proof` / summary | not evaluated | Stage 2 is not evaluated until Stage 1 blockers are resolved. |

## 7. Recommendations

| Priority | Recommendation | Related Risk |
| --- | --- | --- |
| `P0` | Publish batch-verifiable `wallet_ownership_proof` for the PR01JUN26 `wallet_address_list`. | `NO_WALLET_OWNERSHIP_PROOF` |
| `P0` | Publish trusted setup transcript / MPC ceremony, or migrate to a transparent setup proof system. | `OPAQUE_TRUSTED_SETUP` |
| `P0` | Publish `global_proof`, verifying keys, `vk/config`, and proof schema without login or manual request. | `UNVERIFIABLE` |
| `P1` | Add official canonical anchoring for proof/root/vk commitments and historical artifacts. | `NO_CANONICAL_ANCHOR` |
| `P1` | Publish proof constraints that show coverage of lending, margin, collateral haircuts, negative balances, and dummy-user protections. | `BUSINESS_CONSTRAINT_GAP` |

## 8. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate Binance's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities. Traditional audit and PoR should be treated as complementary evidence tracks.

## 9. References

- Assessment JSON: [`PR01JUN26-assessment.json`](./PR01JUN26-assessment.json)
- ardmere audit report: [`AUDIT-REPORT.md`](../../../artifacts/binance/PR01JUN26-20260601-period43/AUDIT-REPORT.md)
- Artifact bundle: [`PR01JUN26.artifact-bundle.json`](../../../artifacts/binance/PR01JUN26-20260601-period43/bundles/PR01JUN26.artifact-bundle.json)
- Verification bundle v2: [`PR01JUN26.verification-bundle.v2.json`](../../../artifacts/binance/PR01JUN26-20260601-period43/bundles/PR01JUN26.verification-bundle.v2.json)

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

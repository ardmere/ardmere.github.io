# OKX 506872725 PoR Transparency Report

> Methodology: [`por-transparency-framework.md`](../../por-transparency-framework.md)
> Assessment JSON: [`506872725-assessment.json`](./506872725-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `okx` |
| Snapshot | `506872725` |
| Snapshot time | `2026-05-06T16:00:00Z` |
| PoR Stage | `Stage 1 — Verifiable Disclosure` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `true` |

OKX publishes public summary, `wallet_address_list`, address signatures, and a public global zk proof bundle for audit `506872725`. ardmere verifies address ownership and global zk binding, so this snapshot reaches **Stage 1**.

The snapshot should not be described as `Stage 2` because official canonical anchoring, stable DA, low-friction user inclusion proof, permissionless replay, and public business-consistent proof constraints are not established in this report.

## 2. Stage Decision

### Stage 1, Effective PoR

OKX satisfies the core Stage 1 evidence requirements: `wallet_address_list` is public, `wallet_ownership_proof` is batch-verifiable, `global_proof` is public, and the global zk proof binds to the summary merkle root.

Stage 2 is blocked by the following gaps:

| Missing / Blocked Evidence | Risk Flag | Stage Effect | Why It Matters |
| --- | --- | --- | --- |
| Official canonical anchor / DA | `NO_CANONICAL_ANCHOR` | max `Stage 1` | Public files are available, but this report does not establish an official exchange canonical on-chain/DA record for root/proof/vk history. |
| Low-friction user inclusion proof | `HIGH_FRICTION_INCLUSION_PROOF` | max `Stage 1` | Public liability zip contains `sum_proof_data.json` only; no low-friction user proof path is archived. |
| Public business-consistent constraints | `BUSINESS_CONSTRAINT_GAP` | max `Stage 1` | ardmere verifies global zk binding, but not full coverage of all OKX business risks and risk controls. |
| Frequency and event-triggered updates | `HIGH_FREQUENCY_GAP` | max `Stage 1` | Observed cadence is monthly; daily root/commitment anchoring and event-triggered updates are not established in this report. |
| Machine-readable independent review | `UNVERIFIABLE` | confidence downgrade | No machine-readable independent third-party attestation artifact is available to this pipeline. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest Snapshot | `2026-05-06T16:00:00Z` |
| Previous Snapshot | `2026-04-06T16:00:00Z` |
| Observed Cadence | `monthly` |
| History Available | `true` |
| Event-triggered Updates | `unknown` |
| Daily Root / Commitment Anchor | `no` |
| Stage Impact | Monthly PoR supports Stage 1 disclosure, but does not satisfy Stage 2 frequency expectations without daily root/commitment anchoring, event-triggered updates, or stronger cadence. |

OKX publishes recurring PoR audits and historical public artifacts. This report does not establish event-triggered publication rules or official daily canonical anchoring.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summary` / `summarySnapshot` | `ee7f35e968ec97a38de84d4a60ad96b20379ef7a2c360a976473d36494dd196a` | `https://www.okx.com/proof-of-reserves/detail` |
| `wallet_address_list` / `walletZip` | `f7fad2b96a64d92c9dff4a94f5a8ea2ae8afa10de0bed21c4e196b71e2132b3b` | `https://static.okx.com/cdn/okx/por/chain/por_csv_2026050700_V2.zip` |
| `global_proof` / `globalProofBundle` | `319a01e352db65d17c3a6fd32b89ef35501c50e9a6b8714608c65112c36eeaee` | `https://www.okx.com/proof-of-reserves/download` |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x8c6ce11cf467a4c74a971c8008d487be586ec9988ab8dde910c5cf110c312921` |
| Verification bundle v2 root | `0x71a08d13032475f58b2c306cfe2807ca2483912bc16be298b0bf8626a07b15c4` |
| Artifact bundle SHA-256 | `0ba21485473840cbc1f25a6f4906f7cb9d15eeb645c608f526d05c3a3f971023` |
| Verification bundle v2 SHA-256 | `a0d627a67e0e2497dd788a8f01da5bf3fdbb82edd51eb8bc990269b4c6773976` |

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` | Summary, wallet ZIP, and global proof bundle hashes match archived artifacts. |
| `solvency-claim` | `1.0` | `PASS` | `1.0000` | BTC, ETH, and USDT self-reported capital ratios are >= 100%. |
| `internal-consistency` | `1.0` | `PASS` | `1.0000` | Wallet ZIP exchangeReserve balances match summary; custody balances are summary-only WARN. |
| `address-ownership` | `okx-1` | `PASS` | `1.0000` | Verified 320,917 / 320,917 address signatures. |
| `global-zk-proof` | `okx-1` | `PASS` | `1.0000` | `verify-global` succeeded; summary merkle root is bound to zk proof. |
| `onchain-balance-hot` | `2.1` | `WARN` | `0.0012` | ETH hot rows include 6 WARN and 61 UNVERIFIABLE findings; no FAIL. |
| `onchain-balance-token` | `2.0` | `FAIL` | `0.0014` | 800-row token sample includes 44 FAIL, 349 WARN, 5 UNVERIFIABLE, and 2 PASS findings. |
| `onchain-balance-ledger` | `1.2` | `FAIL` | `0.0008` | Ledger sample includes 76 FAIL, 299 WARN, 7 UNVERIFIABLE, and 2 PASS findings. |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0` | No public machine-readable attestation artifact. |

The reserve-side FAIL verdicts should not be read as direct evidence of OKX insolvency. The detailed ardmere audit report attributes many findings to omnibus/internal allocation labels, live-vs-snapshot infrastructure limits, public RPC/archive gaps, Aptos FA accounting, and UTXO partial matching. The Stage decision is driven by public Stage 1 artifacts that pass, while Stage 2 remains blocked by trust-minimization requirements.

## 6. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Public PoR summary | pass | OKX detail page summary is archived. |
| Public `wallet_address_list` | pass | Wallet ZIP is public. |
| Public `wallet_ownership_proof` | pass | 320,917 address signatures verified. |
| Public `global_proof` + `vk/config` | pass | `global-zk-proof@okx-1` passes. |
| Trusted setup transparency / transparent setup | pass | OKX uses public zk-STARK / Plonky2-style global proof verification path. |
| Low-friction inclusion proof bound to public root / `global_proof` / summary | fail | Public liability zip contains `sum_proof_data.json` only; user inclusion proof path is not archived. |
| Official canonical anchor / DA | fail | No official exchange canonical on-chain/DA anchor established in this report. |
| Stage 2 frequency and event-triggered updates | fail | Monthly cadence is observed; daily root/commitment anchor and event-triggered updates are not established. |

## 7. Recommendations

| Priority | Recommendation | Related Risk |
| --- | --- | --- |
| `P1` | Publish an official canonical root/proof/vk anchor and durable DA record for every snapshot. | `NO_CANONICAL_ANCHOR` |
| `P1` | Provide low-friction user inclusion proof export and Web/WASM/GUI verification bound to public root/global_proof/summary. | `HIGH_FRICTION_INCLUSION_PROOF` |
| `P1` | Publish proof constraint documentation covering business risks such as lending, margin, collateral haircuts, and negative-balance protections. | `BUSINESS_CONSTRAINT_GAP` |
| `P1` | Add daily root/commitment anchoring and event-triggered updates, or otherwise demonstrate a stronger cadence than monthly snapshots. | `HIGH_FREQUENCY_GAP` |

## 8. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate OKX's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities. Traditional audit and PoR should be treated as complementary evidence tracks.

## 9. References

- Assessment JSON: [`506872725-assessment.json`](./506872725-assessment.json)
- ardmere audit report: [`AUDIT-REPORT.md`](../../../artifacts/okx/506872725/AUDIT-REPORT.md)
- Artifact bundle: [`506872725.artifact-bundle.json`](../../../artifacts/okx/506872725/bundles/506872725.artifact-bundle.json)
- Verification bundle v2: [`506872725.verification-bundle.v2.json`](../../../artifacts/okx/506872725/bundles/506872725.verification-bundle.v2.json)

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

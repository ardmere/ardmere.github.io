# Bybit PoR Transparency Report

> Template version: `ardmere/exchange-transparency-report@1`
> Methodology: `docs/por-transparency-framework.md`
> Assessment JSON: `docs/reports/bybit/2025061709-assessment.json`

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `bybit` |
| Snapshot | `2025061709` |
| Period | `0` |
| Snapshot Time | `2025-06-17T09:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 1 / E1` |
| Confidence | `low` |
| Effective PoR | `false` |

Bybit exposes a reserve ratio summary and a user Merkle proof flow, and ardmere can replay a supplied SHA-256 Merkle path. However, the available public evidence does not include `wallet_address_list`, `wallet_ownership_proof`, or a `global_proof`, so this snapshot remains Stage 0 and should not be treated as effective PoR.

The current archived bundle supports limited Gen 1 / E1 evidence: summary plus optional user inclusion proof. `user-merkle-proof@bybit-1` passes for the supplied v5 test vector and binds to the summary Merkle root, but that proves only single-account inclusion for that JSON input. It does not prove exchange-wide liabilities, reserve ownership, on-chain balances, or business-consistent constraints.

Under the ardmere framework, this snapshot should **not** be described as Stage 1, Stage 2, effective PoR, or trust-minimized PoR unless the corresponding evidence is public, reproducible, and listed in this report.

## 2. Stage Decision

### Stage 0, not effective PoR

Bybit's public design is closer to a Merkle inclusion disclosure than a full public PoR system. The user proof verifier can validate one path, but Stage 1 requires publicly verifiable reserve-side ownership evidence and global liability evidence.

Stage 1 is blocked by the following evidence gaps:

| Missing / Blocked Evidence | Risk Flag | Stage Effect | Why It Matters |
| --- | --- | --- | --- |
| `wallet_address_list` | `UNVERIFIABLE` | Max Stage 0 | Without a public address list, reserve-side balances cannot be independently reconciled. |
| `wallet_ownership_proof` | `NO_WALLET_OWNERSHIP_PROOF` | Max Stage 0 | Users must trust that Bybit controls the claimed reserve addresses. |
| `global_proof` | `UNVERIFIABLE` | Max Stage 0 | A user Merkle path does not prove exchange-wide liability correctness. |
| Business-consistent proof constraints | `BUSINESS_CONSTRAINT_GAP` | Max Stage 0 | The public proof path does not establish protections against negative-net-worth dummy users or incomplete liability modeling. |

Stage 2 is also blocked by the absence of an official canonical anchor, low-friction production inclusion proof verification, stronger publication cadence, and stable data availability.

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest Snapshot | `2025-06-17T09:00:00Z` |
| Previous Snapshot | `unknown` |
| Observed Cadence | `monthly` |
| History Available | `true` |
| Event-triggered Updates | `unknown` |
| Daily Root / Commitment Anchor | `no` |
| Stage Impact | Monthly publication and login-gated user proof can support limited Stage 0 disclosure, but they do not satisfy Stage 1 or Stage 2 without public wallet/proof artifacts and stronger canonical anchoring. |

Bybit references monthly third-party PoR/PoL reporting. This ardmere snapshot uses a fixture summary and official Merkle proof test vector, so cadence is treated as contextual rather than independently verified from production artifacts.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `summary` | `110282f57984ebd218a3ba0ee510b8e08bf6a2bb24de3cad9c0af22d528d422a` | `https://www.bybit.com/x-api/por/public/v1/reserve-ratio/latest` |
| `user_merkle_proof` | `d095eed7e7dde2625c7311d24dbe307e9b12a25574e29ac4e1c073f14fdc2ab8` | `file://./fixtures/bybit/mock_user_merkle_tree_path_40_v5.json` |

The summary is an ardmere fixture imported from a Bybit-style response because datacenter access to the public x-api endpoint is WAF-gated. The user proof is the official `bybit-exchange/merkle-proof` v5 test vector, not a real production account proof.

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x7e2710dde862676f03cdde12e4590d451efca54b597049c5a732884c9fd9e69d` |
| Verification bundle root | `0xd9079261cfa47c44357e30119fdf95a9d077fda29a2a2026cc4a1bd067872c21` |
| Artifact bundle SHA-256 | `c00cc424d1218befb358b0fecb6046499dc6175b79b013079f99b347df989550` |
| Verification bundle SHA-256 | `7d7f949786a4c1d48be58b10858a32e51a7930b2bf71102de2c2d9b1291c83d6` |

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.00` | Summary and user proof hashes match archived files. |
| `solvency-claim` | `1.0` | `PASS` | `1.00` | Self-reported total reserve rate 102.50%; BTC 105.20%, ETH 101.00%, USDT 100.50%. |
| `user-merkle-proof` | `bybit-1` | `PASS` | `1.00` | Supplied v5 SHA-256 Merkle path is valid and binds to the summary root. |
| `internal-consistency` | `0` | `UNVERIFIABLE` | `0.00` | No public wallet address CSV. |
| `btc-anchor` | `0` | `UNVERIFIABLE` | `0.00` | Summary does not declare a BTC block height time anchor. |
| `onchain-balance-hot` | `0` | `UNVERIFIABLE` | `0.00` | No public HotCold wallet address list. |
| `onchain-balance-token` | `0` | `UNVERIFIABLE` | `0.00` | No public wallet address list for token audit. |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.00` | No public wallet ownership signatures/proofs. |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.00` | Pure Merkle flow; no public global ZK proof. |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.00` | Public third-party reports are not machine-readable in this pipeline. |

`PASS` on `user-merkle-proof` means the supplied account path is internally consistent. It does **not** mean that the whole exchange is solvent, that the reserve addresses are controlled by Bybit, or that every liability and business constraint is covered.

## 6. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Public summary / reserve ratio snapshot | Partial | Archived summary fixture is available; production endpoint is WAF-gated for datacenter fetches. |
| Public `wallet_address_list` | Missing | No public wallet CSV/ZIP is available. |
| Public `wallet_ownership_proof` | Missing | `address-ownership` is `UNVERIFIABLE`. |
| Public `global_proof` or equivalent liability proof | Missing | Merkle inclusion proof only; no global ZK proof. |
| Low-friction independently reproducible inclusion proof | Partial | Verifier exists, but production proof is login-gated and archived sample is a test vector. |
| Official canonical anchor and high-frequency publication | Missing | No official on-chain/DA anchor or daily commitment cadence is established. |

## 7. Recommendations

| Priority | Recommendation | Related Risk |
| --- | --- | --- |
| `P0` | Publish a complete `wallet_address_list` and batch-verifiable `wallet_ownership_proof` for every PoR snapshot. | `NO_WALLET_OWNERSHIP_PROOF` |
| `P0` | Publish a public `global_proof` or equivalent exchange-wide liability proof with documented constraints for negative balances, lending, margin, and collateral treatment. | `BUSINESS_CONSTRAINT_GAP` |
| `P1` | Provide machine-readable third-party attestation artifacts and stable public download URLs for summary, wallet, and liability evidence. | `UNVERIFIABLE` |
| `P1` | Move user inclusion proof verification into a low-friction web/WASM or one-click local verifier, and archive anonymized reproducible proof fixtures for each snapshot. | `HIGH_FRICTION_INCLUSION_PROOF` |
| `P2` | Anchor each snapshot root/proof bundle to an official canonical on-chain or DA record and increase publication cadence beyond monthly snapshots. | `NO_CANONICAL_ANCHOR` |

## 8. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate Bybit's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities. Traditional audit and PoR should be treated as complementary evidence tracks.

## 9. References

- Assessment JSON: `docs/reports/bybit/2025061709-assessment.json`
- Artifact bundle: `artifacts/bybit/2025061709/bundles/2025061709.artifact-bundle.json`
- Verification bundle: `artifacts/bybit/2025061709/bundles/2025061709.verification-bundle.json`
- Supporting report: `artifacts/bybit/2025061709/AUDIT-REPORT.md`
- Bybit PoR page: `https://www.bybit.com/en/proof-of-reserves`
- Bybit reserve ratio announcement: `https://www.bybit.com/en/announcement-info/reserve-ratio`
- Bybit Merkle validator: `https://github.com/bybit-exchange/merkle-proof`

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

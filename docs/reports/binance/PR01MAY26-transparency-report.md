# BINANCE PR01MAY26 PoR Transparency Report

> Methodology: [PoR Stage Framework](../../por-transparency-framework.md)  
> Assessment JSON: [PR01MAY26-assessment.json](./PR01MAY26-assessment.json)

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `binance` |
| Snapshot | `PR01MAY26` |
| Snapshot time | `2026-05-01T00:00:00Z` |
| PoR Stage | `Stage 0 — Trust the Exchange` |
| Gen / Evidence | `Gen 2 / E2` |
| Confidence | `medium` |
| Effective PoR | `false` |

Binance snapshot PR01MAY26 provides public reserve summary and wallet_address_list where available, but Stage 1 is blocked by missing wallet_ownership_proof, opaque trusted setup, and unavailable public global proof/vk artifacts.

The available artifacts support Gen 2 / E2 classification and Stage 0 PoR disclosure. Users still need to trust Binance for wallet control, trusted setup honesty, and public availability of the full global proof stack.

## 2. Stage Decision

### Stage 0, not effective PoR

| Missing / Blocked Evidence | Risk Flag | Max Stage | Reason |
| --- | --- | --- | --- |
| — | `NO_WALLET_OWNERSHIP_PROOF` | Stage 0 | No public batch-verifiable wallet_ownership_proof. |
| — | `UNVERIFIABLE` | Stage 0 | Public global proof/vk not available. |
| — | `OPAQUE_TRUSTED_SETUP` | Stage 0 | Trusted setup transcript is not public. |

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest snapshot (evaluation set) | `2026-06-01T00:00:00Z` |
| Previous snapshot | `2026-04-01T00:00:00Z` |
| Observed cadence | `~monthly` |
| History available | `3 snapshot(s) in public evaluation set` |
| Event-triggered updates | `UNVERIFIABLE` |
| Daily root / commitment anchor | `UNVERIFIABLE` |
| Stage impact | Monthly or slower cadence; does not meet Stage 2 frequency expectations. |

Older snapshot PR01JUN26 is also in the public evaluation set.

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `bapiSnapshot` | `7724bc46ca584e53aedf03c08d1248348f67f27670577dae9c4eeb4daf23bafb` | https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot |
| `walletZip` | `6441024de77cd1ca52dfbfbc5f9d2a01a5db12648fc616dbbb76d53e91b1455f` | https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_20260501.zip |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `0x399ca21b9a36fcf0cb31fc1e6bd9f3afeb356a9e3d21c60aabf129b70e4d1e2f` |
| Verification bundle root | `0x48a7a9cd6a8902bd00b34ca93117db0e8fbf4b84085871b2d5774313cc4d9f9c` |
| Artifact bundle SHA-256 | `02da624f0504c305616cc8c58b71557e37fe03043fcf4cb1f6d7c53fcb71f114` |
| Verification bundle SHA-256 | `384a994f5f4abe4e86edd8659e61e9e58ac77adba70638738d96d541cd60d57e` |

Local bundle paths: [PR01MAY26.artifact-bundle.json](../../../artifacts/binance/PR01MAY26/bundles/PR01MAY26.artifact-bundle.json), [PR01MAY26.verification-bundle.v2.json](../../../artifacts/binance/PR01MAY26/bundles/PR01MAY26.verification-bundle.v2.json)

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `artifact-integrity` | `1.0` | `PASS` | `1.0000` |  |
| `internal-consistency` | `1.1` | `PASS` | `1.0000` |  |
| `btc-anchor` | `1` | `UNVERIFIABLE` | `0.0000` | BTC block time anchor verifier not implemented |
| `solvency-claim` | `1.0` | `FAIL` | `1.0000` | summary has no coin rows |
| `onchain-balance-hot` | `2.1` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-token` | `2.0` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-ledger` | `1.4` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-deposit` | `1.2` | `UNVERIFIABLE` | `0.0000` | RPC queries skipped (--skip-rpc) |
| `address-ownership` | `0` | `UNVERIFIABLE` | `0.0000` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | `0` | `UNVERIFIABLE` | `0.0000` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | `0` | `UNVERIFIABLE` | `0.0000` | No public third-party attestation report available |
| `cross-chain-wrapped` | `0` | `UNVERIFIABLE` | `0.0000` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

## 6. Verifier Finding Details

This section explains `FAIL`, `WARN`, and `UNVERIFIABLE` outcomes from the verification bundle. Row-level `UNVERIFIABLE` entries often mean the verifier does not yet support that (coin, network) pair, not that the exchange failed the check. Full machine-readable output: `PR01MAY26.verification-bundle.v2.json` in the local artifact bundle.

### Capability and artifact gaps (`UNVERIFIABLE`)

| Verifier | Explanation |
| --- | --- |
| `btc-anchor` | BTC block time anchor verifier not implemented |
| `onchain-balance-hot` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-token` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-ledger` | RPC queries skipped (--skip-rpc) |
| `onchain-balance-deposit` | RPC queries skipped (--skip-rpc) |
| `address-ownership` | No public download channel for wallet ownership signatures / proofs |
| `global-zk-proof` | Global proof.csv / verifying key not publicly distributed |
| `third-party-attestation` | No public third-party attestation report available |
| `cross-chain-wrapped` | No per-row wrapped-asset metadata in public PoR artifacts (token contract, representation type e.g. WBTC/cbBTC, canonical asset, custody mode) |

### Row-level findings

#### `internal-consistency` (`PASS`)

Finding counts: WARN 1, PASS 15

**WARN**

| Subject | Field | Claim | Actual | Note |
| --- | --- | --- | --- | --- |
| `csv-extras` | `longTailCoins` |  | 50 coins not in summary snapshot: [1INCH AAVE APT ARB ASTER BCH BNB BOME BTC BUSD CAKE CHR CHZ CRV DOGE DOT ENA ENJ ETH FDUSD FORM GRT HBAR HFT LINK LTC MASK NEAR OP PAXG PENDLE PEPE POL RLUSD S SHIB SOL SSV SUI TRUMP TUSD U UNI USD1 USDC USDE USDT WIF WLFI XRP] | informational — binance summary lists top coins only |

#### `solvency-claim` (`FAIL`)

Summary: summary has no coin rows

## 7. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| Has public PoR been published during the last 12 months? (S0-1) | `pass` | Snapshot PR01MAY26 is in the public evaluation set. |
| Is snapshot time explicit? (S0-2) | `pass` | 2026-05-01T00:00:00Z |
| Is there reserve summary or wallet_address_list? (S0-3) | `pass` | Summary and wallet_address_list present. |
| Is wallet_address_list public and on-chain verifiable? (S1-1) | `pass` | Wallet list present; on-chain replay skipped or incomplete. |
| Is wallet_ownership_proof public and batch-verifiable? (S1-2) | `unverifiable` | address-ownership verifier UNVERIFIABLE. |
| Is global_proof public and independently reproducible? (S1-3) | `unverifiable` | global-zk-proof verifier UNVERIFIABLE. |
| Is independent third-party review available with stated scope? (S1-4) | `unverifiable` | third-party-attestation verifier UNVERIFIABLE. |
| If trusted setup is required, is transcript public; otherwise is the proof system transparent setup? (S1-5) | `unverifiable` | No public trusted setup transcript observed for zk-SNARK PoS. |
| Is root/proof/vk canonically anchored outside exchange servers? (S2-1) | `not_applicable` | Stage 2 is not evaluated until Stage 1 blockers are resolved. |
| Is publication frequency sufficient for Stage 2 (weekly full PoR / daily anchor)? (S2-2) | `not_applicable` | Stage 2 is not evaluated until Stage 1 blockers are resolved. |

## 8. Recommendations

| Priority | Recommendation | Related Risk |
| --- | --- | --- |
| `P0` | Publish batch-verifiable wallet_ownership_proof for wallet_address_list. | `NO_WALLET_OWNERSHIP_PROOF` |
| `P0` | Make global proof.csv, verifying key, and parameter sources publicly downloadable. | `UNVERIFIABLE` |
| `P0` | Publish trusted setup transcript or migrate to a transparent-setup proof system. | `OPAQUE_TRUSTED_SETUP` |

## 9. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities.

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

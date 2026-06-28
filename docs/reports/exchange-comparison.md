# Exchange PoR Transparency Comparison

> Methodology: [`por-transparency-framework.md`](../por-transparency-framework.md)
> Scope: first public report set for Binance, OKX, and Bybit.

## Summary

This page compares the first ardmere exchange transparency reports using the same Stage-centered methodology. `Stage 1` is the minimum threshold for effective PoR under the framework; `Stage 0` remains disclosure that still requires users to trust the exchange for key claims.

| Exchange | Snapshot | PoR Stage | Gen / Evidence | Effective PoR | Confidence | Last Verified Snapshot | Report |
| --- | --- | --- | --- | --- | --- | --- | --- |
| OKX | `506872725` | `Stage 1` | `Gen 2 / E2` | `true` | `medium` | `2026-05-06T16:00:00Z` | [`report`](./okx/506872725-transparency-report.md) |
| Binance | `PR01JUN26` | `Stage 0` | `Gen 2 / E2` | `false` | `medium` | `2026-06-01T00:00:00Z` | [`report`](./binance/PR01JUN26-transparency-report.md) |
| Bybit | `2025061709` | `Stage 0` | `Gen 1 / E1` | `false` | `low` | `2025-06-17T09:00:00Z` | [`report`](./bybit/2025061709-transparency-report.md) |

## Key Differences

| Exchange | Main Public Evidence | Missing / Blocked Evidence | Key Risk Flags |
| --- | --- | --- | --- |
| OKX | Public summary, `wallet_address_list`, `wallet_ownership_proof`, and `global_proof` are available; ardmere verifies address ownership and global zk binding. | Stage 2 is blocked by missing official canonical anchor / DA, low-friction user inclusion proof, business-consistent constraints, and stronger publication cadence. | `NO_CANONICAL_ANCHOR`, `HIGH_FRICTION_INCLUSION_PROOF`, `BUSINESS_CONSTRAINT_GAP`, `HIGH_FREQUENCY_GAP`, `UNVERIFIABLE` |
| Binance | Public reserve summary and `wallet_address_list` are available; ardmere can reproduce several reserve-side checks. | Stage 1 is blocked by missing `wallet_ownership_proof`, opaque trusted setup evidence, and unavailable public `global_proof` / verifying key stack. | `NO_WALLET_OWNERSHIP_PROOF`, `OPAQUE_TRUSTED_SETUP`, `UNVERIFIABLE`, `NO_CANONICAL_ANCHOR`, `BUSINESS_CONSTRAINT_GAP` |
| Bybit | Summary plus optional user Merkle inclusion proof flow; ardmere can replay a supplied SHA-256 Merkle path. | Stage 1 is blocked by missing `wallet_address_list`, `wallet_ownership_proof`, `global_proof`, and business-consistent proof constraints. Current archived sample uses fixture/test-vector evidence. | `UNVERIFIABLE`, `NO_WALLET_OWNERSHIP_PROOF`, `BUSINESS_CONSTRAINT_GAP`, `HIGH_FRICTION_INCLUSION_PROOF`, `NO_CANONICAL_ANCHOR`, `HIGH_FREQUENCY_GAP` |

## Interpretation

OKX is the current Stage 1 reference sample because the public artifacts are sufficient to verify the core reserve ownership and global liability proof path. Binance has stronger technology and evidence availability than a pure Merkle-only flow, but missing wallet ownership and trusted setup transparency keep it at Stage 0. Bybit demonstrates the boundary of Merkle inclusion proof: a valid user path can prove one account's inclusion, but it does not prove exchange-wide solvency or reserve ownership.

Missing data is marked as `UNVERIFIABLE`, not treated as `PASS`.

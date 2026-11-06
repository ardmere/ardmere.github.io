# Toobit Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: toobit`).

Toobit uses a **Summation Merkle Tree (SMT)** with **full SHA-256** and domain separators (`SMT_LEAF` / `SMT_NODE`). ardmere tier: **1** (summary + optional user inclusion proof). The official Go offline verifier is [`toobit-docs/por-verifier`](https://github.com/toobit-docs/por-verifier).

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Reserve ratios + merkle root + history | `GET https://bapi.toobit.com/uapi/v1/por/latest-audit` (`?auditId=` optional) | Yes (JSON) | `summarySnapshot` |
| Public PoR overview page | https://www.toobit.com/proof-of-reserves | Yes (SPA) | — |
| User SMT proof JSON | Logged-in download `proof_of_reserve_{auditId}.json` | Login | `userMerkleProof` |
| Wallet addresses | Not observed on public API/page | No | — |
| Global ZK | Not used | — | — |

## Official references

| Resource | URL |
|----------|-----|
| PoR page | https://www.toobit.com/proof-of-reserves |
| User verification help | https://www.toobit.com/support/how-do-i-view-and-verify-my-toobit-assets |
| Latest audit API | https://bapi.toobit.com/uapi/v1/por/latest-audit |
| SMT offline verifier (Go) | https://github.com/toobit-docs/por-verifier |

## Fetch

```bash
# Live public API (latest audit)
por fetch toobit

# Curated fixture (Oct 2026 / Toobit_POR_202610)
por fetch toobit -summary-path ./fixtures/toobit/20261001-summary.json

# Optional: user inclusion proof (login-gated; demo schema in fixtures)
por fetch toobit \
  -summary-path ./fixtures/toobit/20261001-summary.json \
  -user-proof ./fixtures/toobit/user_1_demo_proof.json
```

## Verify

```bash
por verify -exchange toobit -snapshot 20261001 -artifacts ./artifacts/toobit/20261001
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported BTC/ETH/USDT/USDC reserve rates ≥ 100% |
| `user-merkle-proof@toobit-1` | `PASS` if proof JSON present and SMT path validates; `UNVERIFIABLE` if absent |

### Stubbed (UNVERIFIABLE)

Wallet ZIP, on-chain balances, address ownership, global ZK, third-party attestation.

## User proof validation

ardmere ports `toobit-docs/por-verifier`:

- Leaf: `SHA256(SMT_LEAF\|nonce\|ASSET:amount…)` (assets alpha-sorted)
- Internal: `SHA256(SMT_NODE\|left\|right\|ASSET:sum…)`
- Negative sibling balances and (optional) root total-balance consistency are checked
- Root check: recomputed parent hash equals published `root.hash`

Fixture: `fixtures/toobit/user_1_demo_proof.json` (from the official repo `testdata/proofs/` — **demo only**, not bound to live `Toobit_POR_*` roots).

## Boundaries

- **PASS on user proof** = cryptographic inclusion for one account under SHA256-SMT; not exchange-wide solvency.
- **PASS on solvency-claim** = self-reported percent integers from the public API (UI appends `%`).
- Live user proofs require login → cannot Stage-0 bind a real leaf without a logged-in capture.
- Sept 2026 announcement USDT/USDC ratios are swapped vs API — prefer API numbers/root.

# KuCoin Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: kucoin`).

KuCoin publishes monthly Merkle-style PoR with public reserve ratios. ardmere tier: **1** (summary only). There is **no official OSS Merkle verifier** in ardmere's upstream set, so user inclusion stays stubbed.

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Reserve ratio summary | PoR page SSR `__NEXT_DATA__.initialReserve` | Yes | `summarySnapshot` |
| User Merkle path | Logged-in “Verify Now” export | Login | stored only if imported; **not verified** |
| Wallet addresses | Page claims wallet disclosure; no stable permissionless ZIP | Not integrated | — |
| Global ZK | Not used | — | — |

## Official references

| Resource | URL |
|----------|-----|
| PoR page | https://www.kucoin.com/proof-of-reserves |
| How to verify | https://www.kucoin.com/blog/how-to-verify-your-funds-with-kucoin-proof-of-reserve |
| Auditor | Page often unnamed for the latest snapshot; Hacken historically |

## Fetch

```bash
# Live public scrape (no login)
por fetch kucoin

# Or import a saved SSR/bundle fixture
por fetch kucoin -summary-path ./fixtures/kucoin/20260930-summary.json
```

## Verify

```bash
por verify -exchange kucoin -snapshot 20260930 -artifacts ./artifacts/kucoin/20260930
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported per-coin (and derived min total) reserve rates ≥ 100% |

### Stubbed (UNVERIFIABLE)

`user-merkle-proof@0` (no published OSS hash scheme), wallet CSV, on-chain balances, address ownership, global ZK, attestation PDF parser.

## Boundaries

- **PASS on solvency-claim** = self-reported page ratios only.
- Public SSR does **not** expose a merkle root in `initialReserve`.
- Do not treat a downloaded personal Merkle export as verified until an official algorithm/reference implementation is published and ported.

# BingX Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: bingx`).

BingX publishes monthly Merkle Tree PoR with public reserve ratios. ardmere tier: **1** (summary + optional user inclusion).

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Reserve ratio summary | Balance-reserve page / monthly support articles | Browser / article scrape | `summarySnapshot` |
| User Merkle leaf + path | Logged-in verifier page | Login | `userMerkleProof` |
| Wallet addresses | Hacken PoR ownership audit (PDF) | Attestation PDF only | — |
| Global ZK | Not used | — | — |

## Official references

| Resource | URL |
|----------|-----|
| PoR page | https://bingx.com/en/balance-reserve |
| Oct 2026 PoR article | https://bingx.com/en/support/articles/18020766709903 |
| Python verifier | https://github.com/BingX-API/BingX-proof-of-reserves |
| Auditor | Hacken |

## Fetch (import-first)

Reserve JSON endpoints return **device-time / security-policy** errors from many datacenter clients. Prefer support-article scrape or browser capture:

```bash
por fetch bingx -summary-path ./fixtures/bingx/20261015-summary.json

# Optional user proof JSON (leaf + path + root)
por fetch bingx \
  -summary-path ./fixtures/bingx/20261015-summary.json \
  -user-proof ./fixtures/bingx/mock_user_proof.json
```

## Verify

```bash
por verify -exchange bingx -snapshot 20261015 -artifacts ./artifacts/bingx/20261015
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported total/per-coin reserve rates ≥ 100% |
| `user-merkle-proof@bingx-1` | `PASS` if leaf+path recompute to published root; `UNVERIFIABLE` if absent |

### Stubbed (UNVERIFIABLE)

Wallet CSV, on-chain balances, address ownership, global ZK, Hacken PDF attestation parser.

## User proof validation

Port of `BingX-API/BingX-proof-of-reserves` `por.py`:

- Leaf: `SHA256(utf8(assetLine))` where assetLine is `hashedUID,BTC:…,ETH:…,USDT:…,USDC:…,Total:…`
- Path: successive `SHA256(left||right)` with side markers `l:` / `r:` (also accepts `hash:l` / `hash:r`)

Fixture: `fixtures/bingx/mock_user_proof.json` (synthetic self-consistent path for unit tests).

## Boundaries

- **PASS on user proof** = inclusion for one leaf; not exchange-wide solvency.
- Support-article fixtures may omit absolute reserve/liability amounts and merkle root.
- No public wallet ZIP → no on-chain verifiers.

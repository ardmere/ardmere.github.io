# Bitunix Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: bitunix`).

Bitunix uses a **Gen 1** Merkle Tree with **full SHA-256** (64 hex). ardmere tier: **1** (summary + optional user inclusion proof). The official Java verifier is structurally isomorphic to MEXC's Gen-1 path validator (self + sibling path TreeNode), but uses full SHA-256 (not truncated) and a different leaf formula.

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Reserve ratio + merkle root | `queryAuditTimeList` / `queryAuditInfoByTime` | Yes (JSON) | `summarySnapshot` |
| User Merkle path | Logged-in `merkle_tree_bitunix.json` download | Login | `userMerkleProof` |
| Wallet addresses | `walletAddress` URL field (often empty for recent audits) | Partial / empty | — |
| Global ZK | Not used | — | — |

## Official references

See [registry `bitunix`](../config/exchanges/registry.yaml) for canonical URLs. Integration uses:

| Resource | URL |
|----------|-----|
| PoR page | https://www.bitunix.com/proof-of-reserve |
| User verification help | https://www.bitunix.com/hub/helpcenter/article/how-to-verify-my-assets-on-bitunix?id=166 |
| Audit time list | https://api.bitunix.com/asset-hub-api/api/public/v1/reserve/certificate/queryAuditTimeList |
| Audit info by time | https://api.bitunix.com/asset-hub-api/api/public/v1/reserve/certificate/queryAuditInfoByTime |
| Merkle validator (Java) | https://github.com/BitunixOfficial/proof-of-reserves |

## Fetch

```bash
# Live public API (latest audit)
por fetch bitunix

# Curated fixture
por fetch bitunix -summary-path ./fixtures/bitunix/20261103-summary.json

# Optional: user inclusion proof (login-gated; official sample schema)
por fetch bitunix \
  -summary-path ./fixtures/bitunix/20261103-summary.json \
  -user-proof ./fixtures/bitunix/merkle_tree_bitunix.json
```

## Verify

```bash
por verify -exchange bitunix -snapshot 20261103 -artifacts ./artifacts/bitunix/20261103
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported BTC/ETH/USDT reserve rates ≥ 100% |
| `user-merkle-proof@bitunix-1` | `PASS` if `merkle_tree_bitunix.json` present and path validates; `UNVERIFIABLE` if absent |

### Stubbed (UNVERIFIABLE)

Wallet ZIP, on-chain balances, address ownership, global ZK, third-party attestation.

## User proof validation

ardmere ports `BitunixOfficial/proof-of-reserves`:

- Leaf: `SHA256(encryptUid + "," + nonce + "," + compactBalancesJSON)`
- Internal: `SHA256(leftLeaf + rightLeaf + "," + compactBalancesJSON + "," + height)` with balances summed using per-asset precision (`RoundingMode.DOWN`)
- Root check: recomputed parent merkleLeaf equals last path node

Fixture: `fixtures/bitunix/merkle_tree_bitunix.json` (from the official Java repo).

## Boundaries

- **PASS on user proof** = cryptographic inclusion for one account under full SHA-256; not exchange-wide solvency.
- **PASS on solvency-claim** = self-reported ratios only (percent strings without `%` from the public API).
- Empty `walletAddress` on recent periods → no Stage 1 wallet list / ownership path.
- GitHub sample `merkle_tree_bitunix.json` is a **demo leaf**, not bound to the live published merkle root for a given audit.

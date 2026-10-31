# MEXC Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: mexc`).

MEXC uses a **Gen 1** Merkle Tree with SHA-256 truncated to 32 hex characters (128 bits). ardmere tier: **1** (summary + optional user inclusion proof).

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Reserve ratio summary | PoR dashboard / monthly Hacken press release | Browser / press scrape | `summarySnapshot` |
| User Merkle path | Logged-in user downloads `myProof.json` | Login | `userMerkleProof` |
| Wallet addresses | Historical disclosures; no stable permissionless ZIP | Not integrated | — |
| Global ZK | Not used | — | — |

## Official references

See [registry `mexc`](../config/exchanges/registry.yaml) for canonical URLs. Integration uses:

| Resource | URL |
|----------|-----|
| User verification help | https://www.mexc.com/support/article/how-to-verify-user-assets-on-mexc-17827791528836 |
| July 2026 PoR press release | https://blog.mexc.com/press-release/mexcs-july-2026-proof-of-reserves-report/ |
| Merkle validator (Java) | https://github.com/mexcmerkle/mexc-merkle-proof-of-reserves |
| Auditor | Hacken (monthly PoR) |

## Fetch (import-first)

Datacenter IPs commonly receive **HTTP 403 / Access Denied** from `www.mexc.com`. Prefer press-release scrape or browser DevTools capture:

```bash
# Press-release / curated fixture
por fetch mexc -summary-path ./fixtures/mexc/202607-summary.json

# Optional: user inclusion proof (login-gated; official sample schema)
por fetch mexc \
  -summary-path ./fixtures/mexc/202607-summary.json \
  -user-proof ./fixtures/mexc/merkle_proof.json
```

## Verify

```bash
por verify -exchange mexc -snapshot 202607 -artifacts ./artifacts/mexc/202607
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported total/per-coin reserve rates ≥ 100% |
| `user-merkle-proof@mexc-1` | `PASS` if `myProof.json` present and path validates; `UNVERIFIABLE` if absent |

### Stubbed (UNVERIFIABLE)

Wallet CSV, on-chain balances, address ownership, global ZK, Hacken PDF attestation parser.

## User proof validation

ardmere ports `mexcmerkle/mexc-merkle-proof-of-reserves`:

- Leaf: `SHA256(proofId + "," + accountId + "," + "ASSET:amount,...")[0:32]`
- Internal: `SHA256(leftHash + "," + rightHash + "," + balanceString)[0:32]`
- Root check: **hash-only** (does not re-bind aggregated balances at the root like Bybit v5)

Fixture: `fixtures/mexc/merkle_proof.json` (from the official Java repo resources).

## Boundaries

- **PASS on user proof** = cryptographic inclusion for one account under a 128-bit truncated hash; not exchange-wide solvency.
- **PASS on solvency-claim** = self-reported ratios only; press-release fixtures have no public merkle root bind.
- No public wallet ZIP → no `internal-consistency` or on-chain verifiers.

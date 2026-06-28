# Bybit Proof-of-Reserves — data guide

Bybit uses **Gen 1.5** pure SHA-256 binary Merkle Tree (no ZK). ardmere tier: **1** (summary + optional user inclusion proof).

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Reserve ratio summary | PoR dashboard x-api (WAF-gated) | Browser only | `summarySnapshot` |
| User Merkle path | Logged-in user downloads `myProof.json` | Login | `userMerkleProof` |
| Wallet addresses | Hacken audit only | No CSV | — |
| Global ZK | Not used | — | — |

## Official references

| Resource | URL |
|----------|-----|
| PoR page | https://www.bybit.com/en/proof-of-reserves |
| Reserve ratio announcement | https://www.bybit.com/en/announcement-info/reserve-ratio |
| Merkle validator (Java) | https://github.com/bybit-exchange/merkle-proof |
| Auditor | Hacken (monthly PoR + PoL) |

## Fetch (import-first)

Datacenter IPs receive **HTTP 403** from `www.bybit.com/x-api/*`. Capture JSON in browser DevTools:

```bash
# Browser capture → import
por fetch bybit -info-file ./info.json -coins-file ./coins.json

# Or pre-merged bundle
por fetch bybit -summary-path ./fixtures/bybit/2025061709-summary.json

# Optional: user inclusion proof (login-gated)
por fetch bybit -summary-path ./summary.json -user-proof ./myProof.json
```

## Verify

```bash
por verify -exchange bybit -snapshot 2025061709 -artifacts ./artifacts/bybit/2025061709
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported `total_reserve_rate` and per-coin rates ≥ 100% |
| `user-merkle-proof@bybit-1` | `PASS` if `myProof.json` present and v5 path validates; `UNVERIFIABLE` if absent |

### Stubbed (UNVERIFIABLE)

Wallet CSV, on-chain balances, address ownership, global ZK, Hacken PDF attestation parser.

## User proof validation

ardmere ports `bybit-exchange/merkle-proof` v5 schema (40 online assets):

- Leaf: `SHA256(userHash + concat(balances))`
- Internal: `SHA256(leftHash + rightHash + concatBalances + height)`
- Root check: hash + balances + height (stricter than MEXC hash-only)

Fixture: `fixtures/bybit/mock_user_merkle_tree_path_40_v5.json` (from official repo tests).

## Boundaries

- **PASS on user proof** = cryptographic inclusion for one account; not exchange-wide solvency.
- **PASS on solvency-claim** = self-reported ratios only; no wallet or ZK backing.
- No public wallet ZIP → no `internal-consistency` or on-chain verifiers.

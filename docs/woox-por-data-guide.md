# WOO X Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: woox`).

WOO X publishes a **~15-minute live balance-sheet dashboard** plus **monthly Merkle audits** with a permissionless full-tree download. ardmere tier: **1** (summary / solvency claim). Official Java validator: [`woo-x/merkle-tree`](https://github.com/woo-x/merkle-tree) (protobuf binary tree + login user hash).

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Live assets / liabilities / venues | `/asset/porl/assets`, `/asset/porl/liabilities` | Yes (JSON) | fields inside `summarySnapshot.live` |
| Audit list + merkle URL | `/asset/porl/v2/audits` | Yes (JSON) | `summarySnapshot.audit` |
| Audit asset/liability overview | `/asset/porl/v2/snapshot_overview?auditId=` | Yes (JSON) | `summarySnapshot.overview` |
| Full Merkle tree (protobuf) | `https://oss.woo.network/merkle/{auditId}` | Yes (~70–216MB) | URL + size/etag recorded; **not** auto-downloaded |
| User Merkle hash | Logged-in Audit page | Login | — |
| Wallet ownership / global ZK | Not used | — | — |

## Official references

See [registry `woox`](../config/exchanges/registry.yaml) for canonical URLs. Integration uses:

| Resource | URL |
|----------|-----|
| PoR page | https://woox.io/proof-of-reserves |
| Audits list | https://api.woox.io/asset/porl/v2/audits |
| Snapshot overview | https://api.woox.io/asset/porl/v2/snapshot_overview |
| Live assets | https://api.woox.io/asset/porl/assets |
| Live liabilities | https://api.woox.io/asset/porl/liabilities |
| Merkle validator (Java) | https://github.com/woo-x/merkle-tree |

## Fetch

```bash
# Live public API (latest monthly merkle audit + live dashboard totals)
por fetch woox

# Curated fixture
por fetch woox -summary-path ./fixtures/woox/20261015-summary.json
```

Reserve ratio in the summary is:

`sum(overview.asset.usdt_notional) / sum(overview.liabilities.usdt_notional) * 100`

(same formula as the woox.io audit-record UI `.toFixed(1)`).

## Verify

```bash
por verify -exchange woox -snapshot 20261015 -artifacts ./artifacts/woox/20261015
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported total reserve ratio ≥ 100% |

### Stubbed (UNVERIFIABLE)

`user-merkle-proof@0` (OSS exists but needs ~200MB tree + login user hash; not ported), wallet ZIP / ownership, on-chain balances, global ZK, third-party attestation.

## Boundaries

- **PASS on solvency-claim** = self-reported USD-notional assets ≥ liabilities for the **audit overview** only.
- Latest public audit `20261015` can be **under 100%** while the live ~15-min dashboard is above — both are archived; the assessment binds to the merkle audit day.
- Full tree download is public; **period merkle root hex is not** published separately.
- Live dashboard assets include **custody / third-party venue** balances → Stage 1 on-chain ownership is hard.
- Do not commit merkle tree binaries into git (`artifacts/` is gitignored; trees are too large for fixtures).

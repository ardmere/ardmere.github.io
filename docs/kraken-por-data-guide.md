# Kraken Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: kraken`).

Kraken publishes **quarterly** accountant-led Merkle PoR reviews (The Network Firm LLP via LedgerLens). The public PoR page shows per-asset reserve ratios; LedgerLens exposes the combined Merkle root. ardmere tier: **1** (summary only). **Stage 0 attestation-model** — institutional control row; **no Stage 1** (no permissionless wallet ZIP / ownership / machine-parsed CPA report).

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Combined Merkle root + dataset metadata | `GET https://widget.pro.ledgerlens.io/api/kraken-platform-spot-futures` | Yes (JSON) | `summarySnapshot` (ledgerlens block) |
| Per-period reserve ratios | https://www.kraken.com/proof-of-reserves (`window.__INITIAL_PROPS__` CMS) | Yes (HTML/SSR) | `summarySnapshot` (cms block) |
| Accountant widget UI | https://widget.pro.ledgerlens.io/kraken-platform-spot-futures | Yes | observational |
| User Merkle Leaf ID / path | Kraken Pro → Proof of Reserves → Verify yourself | **Login** | `userMerkleProof` if imported; **not verified** |
| Wallet amounts / addresses | Not permissionlessly published | — | — |
| CPA attestation PDF | Named on LedgerLens; not ingested as machine-verifiable artifact | Partial | `third-party-attestation@0` stub |

## Cadence

| | |
|--|--|
| Typical publish | **Quarterly** (quarter-end snapshots; e.g. 2026-06-30, 2026-03-31, 2025-12-31) |
| Monitor window (suggested) | days around quarter-end **+ ~2–3 weeks** publish lag (LedgerLens `updatedAt` often weeks after snapshot) |
| Latest assessed | `20260630` (June 30, 2026); as of 2026-11-06 no later quarter published on the public page |

`config/por-monitor/schedule.yaml` may still be on an open cadence PR — document cadence here + registry until that lands.

## Official references

| Resource | URL |
|----------|-----|
| PoR page | https://www.kraken.com/proof-of-reserves |
| LedgerLens widget | https://widget.pro.ledgerlens.io/kraken-platform-spot-futures |
| Dec 2025 blog | https://blog.kraken.com/news/december-2025-proof-of-reserves |
| The Network Firm | https://www.thenetworkfirm.com/ |

## Fetch

```bash
# Live public sources (latest LedgerLens dataset + matching CMS ratios)
por fetch kraken
por fetch kraken -audit-id 2026-06-30

# Or import a saved bundle fixture
por fetch kraken -summary-path ./fixtures/kraken/20260630-summary.json
```

Snapshot id is the compact calendar form of the CMS audit label date (e.g. June 30, 2026 → `20260630`). CMS may show stables as **105%+** (`capPercentage=true`); ardmere stores an effective floor of **105** for solvency-claim while preserving the display string.

## Verify

```bash
por verify -exchange kraken -snapshot 20260630 -artifacts ./artifacts/kraken/20260630
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported per-coin effective reserve rates ≥ 100% (min rate also stored as total) |

### Stubbed (UNVERIFIABLE) — honest blockers

| Gap | Why |
|-----|-----|
| No wallet list / amounts | Ratios + Merkle root only → Stage 1 `walletZip` / on-chain unwired by design. |
| `user-merkle-proof@0` | Merkle Leaf ID + path need Kraken Pro login / personalized accountant tool input. |
| `third-party-attestation@0` | Accountant named (TNF / LedgerLens) but CPA report not machine-parsed into a verifier. |
| Quarterly cadence | Too infrequent for Stage 2; weak as a high-frequency monitor primary. |
| Global ZK | Not part of this attestation model. |

## Boundaries

- **PASS on solvency-claim** = self-reported CMS ratios (with capped `105%+` floor) only.
- LedgerLens `combinedRoot` is a **full** hex root (unlike truncated CEX dashboards) but still does **not** unlock Stage 1 without permissionless wallets + machine-verified attestation/ownership.
- Do **not** treat the accountant label alone as a passed `third-party-attestation` verifier.

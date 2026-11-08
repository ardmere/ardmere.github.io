# Poloniex Proof-of-Reserves — data guide

> **Upstream registry**: [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml) (`id: poloniex`).

Poloniex publishes monthly Merkle-tree PoR with public reserve ratios and a full Merkle root hash. ardmere tier: **1** (summary only for Stage 0). Address-ownership OSS exists, but the current-period wallet CSV is not on a confirmed permissionless channel.

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Reserve ratio + Merkle root | `proxy/sapi/merkleProof/public/versionInfo` + `snapshotBalance` | Yes | `summarySnapshot` |
| User Merkle leaf + path | Logged-in Wallet → Verification (`merkleProof.json`) | Login | stored only if imported; **not verified** |
| Wallet addresses + ownership sigs | Claimed PoR CSV + [`tools-nodejs-address-verify`](https://github.com/poloniex/tools-nodejs-address-verify) | Sample CSV in OSS; **current CSV URL unconfirmed** | — |
| Global ZK | Marketing copy mentions ZK; no public proof/vk channel | — | — |

## Cadence

| | |
|--|--|
| Typical publish | 月更，1 日 |
| Monitor window (UTC) | 每月 1–3 日 |

`proxy/sapi/merkleProof/public/versionInfo` 2026-11-07 实探 `versionList` 6 期（20260501–20261001）全部 1 日；
`docs/reports/` 三份报告也都是 1 日。

Monitor entry: [`schedule.yaml`](../config/por-monitor/schedule.yaml) `id: poloniex`, `confidence: high`.
See [por-monitor-schedule.md](./por-monitor-schedule.md).

## Official references

| Resource | URL |
|----------|-----|
| PoR page | https://www.poloniex.com/topic/security/proof-of-reserves |
| versionInfo API | https://www.poloniex.com/proxy/sapi/merkleProof/public/versionInfo |
| snapshotBalance API | https://www.poloniex.com/proxy/sapi/merkleProof/public/snapshotBalance?version=YYYYMMDD |
| Address ownership OSS | https://github.com/poloniex/tools-nodejs-address-verify |
| Manual verify announcement | https://www.poloniex.com/announcement/25603738977559 |

## Fetch

```bash
# Live public APIs (no login)
por fetch poloniex
por fetch poloniex -audit-id 20260901

# Or import a saved bundle fixture
por fetch poloniex -summary-path ./fixtures/poloniex/20261001-summary.json
```

API rates arrive as fractions (`1.00` = 100%). The adapter normalizes them to percentage units for `solvency-claim@1`.

## Verify

```bash
por verify -exchange poloniex -snapshot 20261001 -artifacts ./artifacts/poloniex/20261001
```

### Active verifiers

| Verifier | Verdict when |
|----------|----------------|
| `artifact-integrity@1` | Summary SHA256 matches archived file |
| `solvency-claim@1` | Self-reported total/per-coin reserve rates ≥ 100% |

### Stubbed (UNVERIFIABLE) — honest blockers

| Gap | Why |
|-----|-----|
| `user-merkle-proof@0` | Leaf/path download is login-gated. Announced `MerkleVerify` binary has **no OSS source repo** under the poloniex GitHub org. |
| `address-ownership@0` | OSS `VerifyAddress` exists (signed message `"poloniex"`), but README “download PoR file” links are empty; GitHub Releases ship binaries only; guessed CDN paths return AccessDenied/404. Sample `snapshot/poloniex_por.csv` is not a current-period permissionless channel. |
| Wallet ZIP / on-chain | Same CSV gap — no machine-readable public wallet list for this period. |
| Global ZK / attestation | Not published as reproducible public artifacts. |

## Boundaries

- **PASS on solvency-claim** = self-reported page/API ratios only.
- Public `merkleHash` is archived in the summary; ardmere does **not** recompute the liability tree without a public leaf set.
- Do not treat ownership OSS availability as Stage 1 until a current-period wallet CSV is permissionlessly fetchable and wired into `address-ownership`.

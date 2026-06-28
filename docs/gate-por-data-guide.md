# Gate.com Proof-of-Reserves — Data Availability Guide

> Companion to [`verifier-architecture.md`](./verifier-architecture.md) and the [`gateio` adapter](../internal/exchanges/gateio/adapter.go).

## Key differences from Binance

| Dimension | Binance | Gate.com |
|---|---|---|
| Aggregate data | BAPI JSON | Web API `getProofOfReservesInfo` + coin list |
| Wallet address CSV | Public ZIP (HotCold + Deposit) | **Not published** |
| Global zk proof bundle | Logged-in user download only | Logged-in user download from [My Audit](https://www.gate.com/myaccount/myavailableproof) as `zkmerkle_cex_*.tar.gz` |
| User inclusion | User zip + WASM | `user_config.json` + `./main verify user` |
| On-chain balance audit | ardmere can reuse wallet CSV | **No public address list → UNVERIFIABLE** |

Gate PoR page: [https://www.gate.com/zh/proof-of-reserves](https://www.gate.com/zh/proof-of-reserves)

Open-source verifier: [gateio/proof-of-reserves](https://github.com/gateio/proof-of-reserves)

## Data planes

### 1. Public aggregate (dashboard API)

Page content comes from Gate Web API (browser requests; Akamai may block datacenter IPs):

| Endpoint | Purpose |
|---|---|
| `GET /api/web/v1/proof-of-reserves/getProofOfReservesInfo` | Latest audit: Merkle root, total reserve ratio, customer net balance, excess reserves |
| `GET /api/web/v1/proof-of-reserves/getProofOfReservesCoinList` | Per-coin reserve ratio list |
| `GET /api/web/v1/proof-of-reserves/getProofOfReservesList` | Historical batch pagination |

`gateio.Adapter` merges info + coinList into a `summarySnapshot` artifact stored at:

```text
artifacts/gateio/<auditId>/
  raw/<sha256>.json          # raw summary (API merge bundle or import as-is)
  fetch.json                 # fetch metadata
  bundles/                   # written by por-run anchor / verify
    <auditId>.artifact-bundle.json
    <auditId>.verification-bundle.json
    <auditId>.anchor.json
```

### Fetch raw data (archive only, no verification)

```bash
# Recommended: try API first; on failure auto-import fixtures to artifacts/gateio/<auditId>/raw/
./scripts/gateio/gate-save-local.sh

# Direct Gate public API (Akamai may block datacenter IP)
go run ./cmd/por fetch gateio
# or
./scripts/gateio/gate-fetch.sh

# Browser DevTools: copy API response
go run ./cmd/por fetch gateio -info-file ./info.json -coins-file ./coinList.json
./scripts/gateio/gate-import-browser.sh ./info.json ./coinList.json
```

### Fetch + verify + generate bundle

```bash
# API or -summary-path; raw data auto-written to artifacts/gateio/<auditId>/raw/
go run ./cmd/por anchor -exchange gateio -skip-zip

# Manual import
go run ./cmd/por anchor -exchange gateio -summary-path ./summary.json -skip-zip

# Re-run verification from archived snapshot
go run ./cmd/por verify -exchange gateio -snapshot <auditId> \
  -artifacts ./artifacts/gateio/<auditId>
```

**Local workaround for API blocks (legacy pattern; now defaults to artifacts/gateio/):**

### 2. zk global proof bundle (login required)

Users download from [My Audit](https://www.gate.com/myaccount/myavailableproof):

- **Download Merkle Tree** → `zkmerkle_cex_xxx.tar.gz`
- **Download User Config** → `user_config.json` (place in `config/`)

Extracted structure (from [Gate official docs](https://github.com/gateio/proof-of-reserves)):

```text
config/
  cex_config.json    # CexAssetsInfo + proof.csv path + vk prefix
  user_config.json   # user inclusion (optional)
proof.csv
zkpor864.vk.save
main                 # verifier binary (GitHub Releases)
```

**Exchange asset verification (global zk):**

```bash
./main verify cex
# success output: All proofs verify passed!!!
```

**User inclusion verification:**

```bash
./main verify user
# success output: verify pass!!!
```

Import tar.gz manually into artifact bundle (archived to `artifacts/gateio/<auditId>/raw/`):

```bash
go run ./cmd/por anchor -exchange gateio \
  -summary-path ./summary.json \
  -zk-bundle ./zkmerkle_cex_xxx.tar.gz \
  -skip-zip
```

> ardmere currently archives zk bundle as `globalProofBundle` artifact; `global-zk-proof@gateio-1` verifier not yet implemented (stub `@gateio-0`).

### 3. Missing public data → UNVERIFIABLE dimensions

| Verifier | Gate status | Reason |
|---|---|---|
| `internal-consistency` | UNVERIFIABLE | No wallet CSV for per-coin reconciliation with summary |
| `onchain-balance-*` | UNVERIFIABLE | No public HotCold address list |
| `btc-anchor` | UNVERIFIABLE | Gate summary does not bind BTC block height |
| `global-zk-proof@gateio` | UNVERIFIABLE (default) | tar.gz requires login; public API does not provide |
| `solvency-claim` | Partially verifiable | Public `total_reserve_rate ≥ 100` (self-reported) |

## ardmere verification path (current)

```text
Public API / manual summary.json
        ↓
gateio.Adapter → por.Snapshot
        ↓
solvency-claim@1 (total reserve ratio ≥ 100%, labeled self-reported)
        ↓
stub verifiers (honestly labeled UNVERIFIABLE + reason)
```

Future activation of `global-zk-proof@gateio-1`: parse `cex_config.json` + `proof.csv` + vk; reuse Go verifier from [gateio/proof-of-reserves](https://github.com/gateio/proof-of-reserves) or exec `./main verify cex`.

## References

- [Gate PoR site](https://www.gate.com/zh/proof-of-reserves)
- [Gate Learn — how to verify](https://www.gate.com/learn/articles/how-to-use-gate-io-proof-of-reserves-to-verify-your-assets-security/1017)
- [gateio/proof-of-reserves README](https://github.com/gateio/proof-of-reserves)

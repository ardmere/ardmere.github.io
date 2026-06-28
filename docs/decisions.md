# ardmere zkPoR Verifier — Architecture Decision Records

> This file records irreversible architecture-level decisions (ADRs).
> Companion: [`verifier-architecture.md`](./verifier-architecture.md)

| Field | Value |
|---|---|
| Document version | v1.0 |
| Decision date | 2026-05-05 |

---

## ADR-001 Backend language: **Go**

**Options**: Go / TypeScript

**Decision**: Go

**Rationale**:
- Reuse the verifier package from [`binance/zkmerkle-proof-of-solvency`](https://github.com/binance/zkmerkle-proof-of-solvency) directly; future `global-zk-proof@1` activation is nearly zero cost
- Native-friendly streaming for large files (94 MB walletZip / 275 MB CSV)
- Simple goroutine model for multi-chain RPC concurrency
- Static binary, minimal deployment

**Impact**:
- React island frontend decoupled from backend via REST + JSON
- Browser WASM user verification (decision 4) requires compiling Go to WASM (`GOOS=js GOARCH=wasm`, ~5–10 MB binary), acceptable

---

## ADR-002 Primary anchor chain: **Base**

**Options**: Base / Arbitrum / Optimism / Ethereum L1

**Decision**: Base

**Rationale**:
- Extremely low gas (~$0.001 per anchor tx), supports high-frequency heartbeat
- EVM compatible; mature Foundry / viem toolchain
- Data availability inherited from Ethereum L1 (OP Stack via batch posting)
- Multiple public RPC providers, aligned with ADR-005 zero self-operated principle
- Coinbase backing; strong long-term availability expectations

**Impact**:
- Deploy one minimal event-only contract `ArdmerePoRAnchor`
- Service holds one secp256k1 private key (KMS-managed) as anchor signer
- Chain ID fixed in config: `8453` (Base mainnet)

**Future option**:
- Annual audit snapshots could additionally anchor to Bitcoin OP_RETURN (strongest but expensive); out of MVP scope

---

## ADR-003 Cold storage strategy: **No Arweave / IPFS; on-chain digest only**

**Options**:
- A) Arweave permanent storage for raw verification bundle JSON
- B) IPFS pinning (multiple pinners)
- C) **No independent cold storage; on-chain digest only**
- D) S3 + self-maintained backup

**Decision**: **C — No independent cold storage; on-chain digest only**

**Rationale**:
- Simpler architecture; lower ops cost and external dependencies
- On-chain anchored Merkle root already provides core "post-hoc immutability" guarantee
- Raw data retained in our Postgres + S3 (operational redundancy sufficient)
- Truly unrecoverable data is "Binance walletZip published for that period" — outside our control; we can only retain sha256 on-chain for post-hoc comparison
- Arweave one-time payment is cheap but adds dependency and ops surface unnecessary for MVP

**Impact**:

| Data | Storage | Immutability guarantee |
|---|---|---|
| Merkle root (artifact bundle / verification bundle) | **Base on-chain event** | ✅ On-chain irreversible |
| Raw artifact (bapiSnapshot JSON / walletZip) | **S3 + Postgres** | ⚠ Our operational guarantee + on-chain sha256 comparison |
| Raw verification (findings JSON) | **S3 + Postgres** | ⚠ Our operational guarantee + on-chain sha256 comparison |
| Our ed25519 public key | **Base on-chain + domain DNSSEC TXT** | ✅ Dual-source verifiable |

**Revised §6 anchor objects (final, schema v2)**:

```solidity
event SnapshotAnchored(
    bytes32 indexed snapshotId,
    bytes32 indexed exchangeTag,       // keccak256("binance")
    string exchange,
    bytes32 artifactBundleRoot,        // raw artifact digest Merkle root
    bytes32 verificationBundleRoot,    // verifier conclusion Merkle root
    uint8 verdictSummary,
    uint16 coverageBps,
    // … periodSeq, snapshotTime, btcBlockHeight, exchangeMerkleRoot, schemaVersion, anchoredAt
);
```

> Each snapshot period uses **1 tx** to anchor both artifact and verification roots. Raw data retrieval goes through ardmere API (`GET /artifacts/:sha256`); users verify against on-chain roots locally.

**Trade-off acknowledged**:
- Lost: strong guarantee that "historical raw data remains accessible if ardmere disappears"
- Gained: simplified architecture, converged ops, faster launch
- Mitigation: if demand is strong post-MVP, add independent Arweave mirror service in V2 (anchor roots unchanged)

---

## ADR-004 User verification form: **Fully local browser WASM**

**Options**:
- A) User uploads zip to our backend for verification
- B) **Fully local in-browser WASM verification (zero upload)**
- C) Desktop CLI for users to run themselves

**Decision**: B

**Rationale**:
- User-downloaded Merkle Tree zip contains sensitive account identity (account index, balance structure); **any upload is a privacy leak risk**
- In-browser WASM fully local:
  - Near one-click UX
  - ardmere never touches user data; zero liability, zero compliance exposure
  - Can embed on ardmere.org integrated with main site visuals
- Go-to-WASM size acceptable (5–10 MB, first load + cache)

**Implementation path**:
- Compile user verifier portion of `binance/zkmerkle-proof-of-solvency` to WASM
- Frontend `<input type="file">` accepts zip → File API read → WASM memory → PASS/FAIL output
- ardmere backend **does not participate** in this path; do not even log

**Impact**:
- No user-data APIs, database tables, or compliance flows needed
- Independent React island frontend (integrated with main site but state-isolated)

---

## ADR-005 Multi-exchange roadmap: **OKX second**

**Options**: OKX / Coinbase / Bitget / Kraken / Bybit / others

**Decision**: After MVP, **prioritize OKX**

**Rationale**:
- OKX also open-sourced PoR tooling ([`okx/proof-of-reserves-v2`](https://github.com/okx/proof-of-reserves-v2)); data spec similar to Binance
- OKX is second-largest exchange by volume; coverage immediately expands user base
- Same zk-SNARK + Merkle Tree approach; most verifiers reusable
- Strong narrative comparing Binance and OKX under one independent standard

**Impact**:
- Domain model [`Snapshot`](./verifier-architecture.md#3-domain-model) must include `exchange` field from day 1; cannot hardcode Binance
- Config layer abstracts `ExchangeAdapter`, one per exchange:
  ```go
  type ExchangeAdapter interface {
      Id() string                            // "binance" | "okx"
      Fetchers() []Fetcher
      DefaultVerifiers() []Verifier
      ParseSnapshotId(s string) (SnapshotId, error)
  }
  ```
- MVP registers only BinanceAdapter; interface reserves OKX slot

---

## Decision timeline

```
2026-05-05  ADR-001..005 finalized (v1.0)
   Future decisions (if any) append ADR-006…; old ADRs are not edited, only superseded
```

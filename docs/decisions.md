# Architecture Decision Records

> Companion to [`verifier-architecture.md`](./verifier-architecture.md).  
> Five ADRs finalized **2026-05-05**. Status column reflects implementation reality as of the **2026-10** public-site refresh.

| ADR | Decision | Status |
| --- | --- | --- |
| [ADR-001](#adr-001-backend-language-go) | Backend language: **Go** | Accepted · implemented |
| [ADR-002](#adr-002-primary-anchor-chain-base) | Primary anchor chain: **Base** (8453) | Accepted · **Sepolia live; Base mainnet pending** |
| [ADR-003](#adr-003-cold-storage-no-arweave-ipfs-on-chain-digest-only) | Cold storage: **no Arweave/IPFS; digests on-chain only** | Accepted · implemented (local/`artifacts/` + future S3) |
| [ADR-004](#adr-004-user-verification-form-local-wasm) | User verification: **fully local in-browser WASM** | Accepted · roadmap (deposit spot-check page exists; WASM L3 pending) |
| [ADR-005](#adr-005-multi-exchange-roadmap-okx-after-mvp) | Multi-exchange: **OKX prioritized after Binance MVP** | Accepted · six adapters live (OKX / Binance / Gate / Bitget / Bybit / HTX) |

---

## ADR-001 Backend language: Go

**Date:** 2026-05-05 · **Status:** Accepted

### Context

The verifier must reuse exchange-published proof tooling (notably Binance `zkmerkle-proof-of-solvency` and OKX zk-STARK CLI bindings), run concurrent on-chain sampling, and ship as a single CLI (`por`) operators can run locally.

### Decision

Implement the verification backend and CLI in **Go**.

### Consequences

- Direct reuse of Go verifier packages and straightforward FFI to external CLIs.
- Static binary distribution; matches GitHub Pages + operator-laptop workflow.
- Rejected alternatives: TypeScript-only backend (weaker concurrency / zk package reuse), Rust monolith (higher integration cost for MVP).

---

## ADR-002 Primary anchor chain: Base

**Date:** 2026-05-05 · **Status:** Accepted (target) · **Runtime today: Ethereum Sepolia only**

### Context

Root digests (`artifactBundleRoot`, `verificationBundleRoot`) must be cheap enough for per-snapshot anchors and weekly heartbeats, while remaining publicly verifiable on an EVM explorer.

### Decision

**Primary production anchor chain = Base** (OP Stack L2, chain id `8453`). Testnet development and the first live deployments use **Ethereum Sepolia** (`11155111`).

### Current deployment reality

| Network | Role | Contract | Notes |
| --- | --- | --- | --- |
| Ethereum Sepolia | Active testnet primary | See [`deployments.md`](./deployments.md) | UUPS proxy + impl v3 verified |
| Base mainnet (`8453`) | **Intended production primary** | *Not deployed* | Blocked on operator keys / Base RPC / funded deployer |
| Base Sepolia | Optional L2 testnet | *Not deployed* | May precede mainnet |

Until Base mainnet is deployed, CLI defaults (`.env.example`, `docs/ledger-rpc-runbook.md`) keep **Sepolia** as the operable anchor chain. Architecture text that says “Primary = Base” refers to the **product target**, not the currently live chain.

### Consequences

- Low gas for frequent anchors; EVM toolchain (Foundry) unchanged.
- Operators must not treat Sepolia anchors as production integrity proofs.
- Base mainnet deploy checklist: funded deployer, `BASE_RPC_URL`, `PRIVATE_KEY` / `ANCHOR_SIGNER` in operator env, verify on BaseScan/Sourcify, then update [`deployments.md`](./deployments.md) and `.env.example`.

---

## ADR-003 Cold storage: no Arweave/IPFS; on-chain digest only

**Date:** 2026-05-05 · **Status:** Accepted

### Context

Raw PoR artifacts (wallet zips, summaries, zk bundles) are large. Putting bulk bytes on-chain or on permanent public storage (Arweave / IPFS) adds cost and availability coupling without improving digest integrity.

### Decision

- Store raw artifacts in operator-controlled storage (today: local gitignored `artifacts/`; target: S3 + Postgres metadata).
- Anchor **Merkle roots / digests only** on-chain.
- Do **not** require Arweave or IPFS as a cold tier for MVP.

### Consequences

- Public GitHub Pages **does not** host `artifacts/` (directory is gitignored). Archive index links may 404 for casual visitors; reproducibility requires a local `artifacts/` tree or a future release mirror.
- On-chain records remain tamper-evident for roots even when bulk files live off-site.

---

## ADR-004 User verification form: local WASM

**Date:** 2026-05-05 · **Status:** Accepted · partially delivered

### Context

User inclusion proofs must not require uploading private account trees to ardmere servers.

### Decision

User inclusion verification is **fully local** (in-browser WASM / local CLI). The hosted site may guide the flow but must not receive the user zip contents.

### Consequences

- Privacy-preserving L3 path; zero upload of user proofs.
- Public deposit spot-check page (`verify/deposit/`) is a related, lighter-weight public check—not a substitute for full WASM inclusion verify.
- Full WASM productization remains on the roadmap (architecture W4).

---

## ADR-005 Multi-exchange roadmap: OKX after MVP

**Date:** 2026-05-05 · **Status:** Accepted · expanded

### Context

Binance was the MVP evidence surface. Other CEXs differ in public API depth, login walls, and zk packaging.

### Decision

After Binance MVP, prioritize **OKX** (full public artifact stack → Stage 1 reference), then broaden adapters.

### Current adapter set

| Exchange | Public depth (typical) | Notes |
| --- | --- | --- |
| OKX | High | Often Stage 1 when wallet + ownership + global zk are public |
| Binance | Medium | Summary + wallet zip; ownership / global zk often stub |
| Gate.io / Bitget / Bybit | Low–medium | Scrape / browser capture common → Stage 0 acceptable when honest |
| HTX | Sample / login-gated | May remain on dated sample until `public-data.zip` is obtainable |

### Consequences

- Registry + CLI `por exchanges` is the canonical upstream map.
- Scrape-only Stage 0 publications are allowed; gaps stay `UNVERIFIABLE`.

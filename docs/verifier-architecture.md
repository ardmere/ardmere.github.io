# ardmere zkPoR Verifier — Architecture

> An **independent, replayable, extensible** verification service for Binance Proof-of-Reserves (and future exchanges).
> It also models dimensions that **cannot be verified today** as first-class citizens, so they can be activated seamlessly when external conditions mature.

| Field | Value |
|---|---|
| Document version | v0.1 |
| Last updated | 2026-06-07 |
| Companion doc | [`binance-por-data-guide.md`](./binance-por-data-guide.md) |
| Status | Draft — suitable as MVP implementation reference |

---

## 1. Design Philosophy

Three core principles that drive all subsequent structure:

1. **Evidence-first, not Boolean-first**. Every verification run outputs more than `PASS / FAIL` — it produces a **replayable evidence bundle**: input data hashes + verification logic version + verification result + timestamp + our signature + **key digests anchored on-chain**. Even if Binance changes APIs, deletes files, or swaps algorithms later, historical conclusions remain traceable and auditable.
2. **Verifier is a plugin, not if/else**. Each verification dimension (internal consistency, on-chain balance, zk proof, address signatures, …) is an independent `Verifier` connected through a unified interface. **Adding a dimension = adding one file + registration** — the core stays untouched.
3. **Honesty over completeness**. Any dimension that cannot be verified today **explicitly outputs `UNVERIFIABLE` with a reason** — never degrading to `PASS` or `PENDING`. This is the most important trust asset of a PoR service.

---

## 2. Current Data Availability Boundaries (Determines Verifiable Dimensions)

See [`binance-por-data-guide.md`](./binance-por-data-guide.md) for full detail; summary here:

| Data | Public channel | Affected verifier types |
|---|---|---|
| Current Merkle Root + asset summary | ✅ BAPI | `internal-consistency`, `solvency-claim` |
| Historical snapshot list + BTC block height | ✅ BAPI | `btc-anchor` |
| Wallet addresses + balances (plaintext CSV) | ✅ Public ZIP | `internal-consistency`, `onchain-balance-*` |
| Global zk-SNARK proof.csv / vk | ❌ Logged-in users only | `global-zk-proof` (stub) |
| Address signatures / ownership proofs | ❌ No official download channel | `address-ownership` (stub) |
| Third-party audit attestation | ❌ No official download channel | `third-party-attestation` (stub) |

---

## 3. Domain Model

```text
Snapshot                     ── A Binance PoR snapshot for one period
├── id            : "PR01APR26"
├── snapshotTime  : 2026-04-01T00:00:00Z
├── btcAnchor     : { height: 943129, hash, timestamp }
├── claims        : Claim[]                    ← Binance self-reported "fact claims"
└── artifactRefs  : sha256[]                   ← content-addressed Artifact references

Claim                        ── A single verifiable assertion
├── kind          : "merkleRoot" | "coinReserve" | "addressBalance" | ...
├── subject       : subject (coin / address / batch)
├── value         : self-reported value
└── source        : (artifactSha256, locator)

Artifact                     ── Raw data snapshot (content-addressed, immutable)
├── kind          : "bapiSnapshot" | "walletZip" | "globalProof" | "addressSig"
├── sha256        : content hash (primary key)
├── source        : { url, fetcherId, fetcherVersion }
├── fetchedAt     : timestamp
├── ourSignature  : ed25519(sha256 ‖ fetchedAt ‖ url)
├── onchainTx     : on-chain anchor tx hash (see §6)
└── storage       : { s3: key }                ← S3 only (ADR-003)

Verification                 ── One verification execution
├── id            : ULID
├── snapshotId    : "PR01APR26"
├── verifierId    : "internal-consistency@1.2"
├── verifiedAt    : timestamp
├── inputHashes   : sha256[]                   ← which Artifacts were read during verification
├── verdict       : PASS | FAIL | UNVERIFIABLE | PARTIAL
├── findings      : Finding[]                  ← line-by-line details
├── coverage      : 0.0 ~ 1.0                  ← sampling ratio (required for on-chain audit)
├── ourSignature  : ed25519(canonical JSON of entire Verification)
└── onchainTx?    : on-chain anchor for key verifications (root-level conclusions only)
```

> **`coverage` field**: L2 on-chain audit is always sampled — it must explicitly declare "how much value / how many addresses we covered." The frontend must display coverage; it cannot be hidden.

---

## 4. Verifier Plugin Interface

```ts
interface Verifier<C extends Claim = Claim> {
  readonly id: string;              // "internal-consistency@1.2"
  readonly version: string;         // semver; logic changes must bump version
  readonly claimKinds: Claim["kind"][];
  readonly requires: Artifact["kind"][];

  // Core method: pure function; input = all artifacts from ctx; output = verification conclusion
  verify(ctx: VerifyCtx): Promise<Verification>;

  // Self-description (for frontend to dynamically render "why this item is UNVERIFIABLE")
  capability(): {
    canVerify: Claim["kind"][];
    cannotVerify: { kind: Claim["kind"]; reason: string }[];
  };
}

interface VerifyCtx {
  snapshot: Snapshot;
  artifacts: Map<sha256, Artifact>;
  // Controlled external dependencies (injected, not self-fetched — easy to mock and replay)
  rpc: ChainRpcPool;             // §7 public node pool
  clock: () => Date;
}
```

**Key design points**:

- `verify()` is a **pure function** — input is fully provided by ctx; direct `fetch()` / `fs.readFile` is forbidden — enabling **replay, unit testing, and future on-chain zk**.
- Fetching data is the job of a **separate `Fetcher` component**; verifiers only consume persisted artifacts.
- Bumping `version` by one → old verifications are automatically marked `STALE`, triggering re-runs (old records kept for comparing logic-change diffs).
- The `rpc` pool is injected externally (§7); verifiers do not know which node provider is used.

---

## 5. Verifier Matrix

### 5.1 Day 1 Launch (100% feasible on public data)

| Verifier | Input artifacts | Verification logic | Output | coverage |
|---|---|---|---|---|
| `artifact-integrity@1` | All | Recorded sha256 at ingest matches recompute; our signature valid; onchainTx confirmed | PASS / FAIL | n/a |
| `internal-consistency@1` | bapiSnapshot + walletZip | For 14 major coins, CSV aggregation = BAPI `binanceLiability/exchangeBalance/thirdPartyCustody` | PASS / FAIL | 100% |
| `btc-anchor@1` | bapiSnapshotList + public chain query | Snapshot-declared BTC block height timestamp ≈ snapshotTime (error < 30 min) | PASS / FAIL | 100% |
| `solvency-claim@1` | bapiSnapshot | 14 major coins: `binanceLiability ≥ customerLiability` (note: only checks relationships among Binance self-reported numbers — **does not prove reserve authenticity** — labeled separately as "self-reported") | PASS / WARN | 100% |
| `onchain-balance-hot@1` | walletZip + public RPC | HotCold rows: at CSV `Height` block, compare **observable balance** with CSV `balance` (see §5.3) | PASS / FAIL / WARN | per `(coin,network)` pair |
| `onchain-balance-deposit@1` | walletZip + public RPC | Deposit value-weighted sampling audit (~10⁵ addresses for ~99% value coverage) | PARTIAL | 0~99% |

### 5.2 Placeholder stubs (ship today, activate later)

| Verifier | Stub output | Activation condition |
|---|---|---|
| `address-ownership@0` | `UNVERIFIABLE` reason: *"No public download channel for wallet ownership signatures / proofs"* | Binance provides official download channel |
| `global-zk-proof@0` | `UNVERIFIABLE` reason: *"Global proof.csv / verifying key not publicly distributed; only available via logged-in user download"* | Binance public release / third-party mirror |
| `third-party-attestation@0` | `UNVERIFIABLE` reason: *"No public third-party attestation report available"* | Any institution publishes publicly downloadable PoR attestation |
| `cross-chain-wrapped@0` | `UNVERIFIABLE` reason: *"Wrapped tokens (wBTC/cbBTC/...) reconciliation rules not finalized"* | We finalize wrapped-asset reconciliation spec |

> **Value of placeholder verifiers**: The frontend UI **today** shows all 10 verification dimensions — 6 with real ✅/❓ verification + 4 explicitly marked ❌ "cannot verify + reason." This is more honest and more compelling than showing only the 6 real dimensions.

### 5.3 `onchain-balance-hot`: staking-aware evolution

#### 5.3.1 Problem: v1 only checks liquid native balance

Current implementation (`onchain-balance-hot@1.0`) calls `eth_getBalance(address, Height)` for each HotCold row and compares directly with the CSV `balance` field.

Binance walletZip `balance` is **book balance** — the same address may include:

| Form | In EOA `balance`? | Visible in v1? |
|---|---|---|
| Liquid native (ETH/BNB) | ✅ | ✅ |
| BNB Stake Hub staking (`getPooledBNB`) | ❌ credited in contract ledger | ❌ |
| BNB unbonding queue (`lockedBNBs`) | ❌ | ❌ |
| Ethereum 2.0 deposits (`DepositContract`) | ❌ moved to beacon chain | ❌ |
| ERC20 / BEP20 | ❌ in token contracts | ❌ (separate `onchain-balance-erc20`) |

**PR01JUN26 empirical data**: Of 3 native-only FAILs, 2 were staking/deposit accounting gaps (`0xbf83…` BNB staking, `0x32e11…` 80k ETH deposited), and 1 (`0x86523…`) had ~**75k / 82k BNB** gap attributable to Stake Hub delegation (see §5.3.6). v1 marking these as FAIL **overstates true deviation**.

Design principle: **Do not treat "queries we haven't implemented yet" as Binance fraud** — either extend observable balance, or explicitly `WARN` and explain uncovered forms.

#### 5.3.2 Target balance model: `totalAccounted`

For each HotCold row `(coin, network, address, Height)`:

```text
totalAccounted = liquidNative
               + stakedNative      // on-chain verifiable staked ledger
               + unbondingNative   // unbonding portion not yet back in EOA
               + (future) erc20Balances[coin]
```

Compare with CSV:

```text
diff = totalAccounted - csvBalance
|diff| ≤ tolerance  →  PASS (or WARN if diff non-zero but within tolerance)
|diff| > tolerance  →  FAIL (Finding must break down components for third-party replay)
```

Tolerance follows v1: `abs ≤ 1e-4` native unit **or** `rel ≤ 1e-7` (see `docs/por-cli.md`).

**Extended Finding fields** (from `onchain-balance-hot@2.0`):

```json
{
  "subject": "0x86523c87…",
  "field": "BNB@BSC#101590091",
  "claim": "172797.6409",
  "actual": "166398.0",
  "status": "PASS",
  "components": {
    "liquid": "90932.07",
    "staked": "75466.0",
    "unbonding": "0",
    "unsupported": "0"
  },
  "note": "staked via StakeHub credit contracts; 53 validators scanned"
}
```

If a form is not yet implemented (e.g. staking on some chain), leave `components.unsupported` empty and explain in `note` → verdict downgrades to **`WARN`** ("incomplete observation") — **must not** mark FAIL.

#### 5.3.3 BSC — BNB Stake Hub (`onchain-balance-hot@2.0`)

| Item | Value |
|---|---|
| StakeHub contract | `0x0000000000000000000000000000000000002002` |
| Snapshot block height | CSV row `Height` (same as v1) |
| Validator list | `getValidators(offset, limit)` → `(operator[], credit[])` |
| Per-validator query | On `credit`: `getPooledBNB(delegator)`, `lockedBNBs(delegator, 0)` |
| Aggregation | `stakedNative = Σ getPooledBNB`; `unbondingNative = Σ lockedBNBs` |

```text
┌─ HotCold row: BNB|BSC, address A, height H ─────────────────────┐
│  liquidNative  ← eth_getBalance(A, H)                           │
│  stakedNative  ← Σ_credit getPooledBNB(A) @ H                   │
│  unbonding     ← Σ_credit lockedBNBs(A, 0) @ H                  │
│  totalAccounted = liquid + staked + unbonding                   │
│  compare ↔ CSV balance                                          │
└─────────────────────────────────────────────────────────────────┘
```

**Performance optimizations** (avoid 53 validators × 2 calls × 10³ rows):

1. **Validator list cache**: For the same `Height`, fetch `getValidators` once; store in Postgres `(height → credit[])`.
2. **Sparse scan**: For a `delegator`, read StakeHub `Delegated` / `Redelegated` logs first (`topics[2] = delegator`) to get **non-zero operator set**; query only relevant credit contracts (usually ≪ 53).
3. **Zero skip**: Credits where `getPooledBNB == 0 && lockedBNBs == 0` are not re-queried (reuse results across rows for same delegator).

**Archive requirement**: BSC snapshot block height is typically days behind tip; most publicnode instances **lack** state at that height (`historical state not available`). `rpc-providers.json` must annotate `archiveDepth`; Stake Hub queries and `eth_getBalance` **share** the same archive-capable provider (tested: `bsc.drpc.org` works, but free tier hits 429 easily — needs failover + cache).

#### 5.3.4 Ethereum — Beacon deposits (`onchain-balance-hot@2.1`)

| Item | Description |
|---|---|
| Scenario | HotCold row `ETH|ETH`, CSV balance is large integer (e.g. 80,000 ETH), EOA `balance ≈ 0` |
| On-chain evidence | Same address has multiple `DepositContract.deposit()` calls (Etherscan Deposits tab) |
| Query | Beacon chain API / execution-layer deposit events + sum ETH deposited by address (32 ETH multiples) |
| Accounting | `stakedNative += Σ depositAmount` (deposits confirmed before `Height`) |

**PR01JUN26 empirical data**: `0x32e11a20337ebc79abd0eeab2d91bafbd9591149` — CSV **80,000.017 ETH**, liquid **0.017 ETH**, Etherscan **40 × 2,000 ETH deposits**; book balance matches on-chain staking records; native-only FAIL was a false positive.

Implementation can be phased:

- **2.1 heuristic**: If `liquid << claim` and address has deposit events and `Σ deposits ≈ claim` → PASS + `components.staked`
- **2.2 full**: Integrate beacon API (e.g. `beaconcha.in` / self-hosted lighthouse) to query effective balance by validator pubkey

#### 5.3.5 Verdict and coverage rules

| Condition | Verdict |
|---|---|
| All supported `(coin,network)` pairs: `totalAccounted` within tolerance | **PASS** |
| Row `totalAccounted` exceeds tolerance | **FAIL** (must include `components`) |
| Row only has liquid implemented, but heuristic detects likely staking (large claim + near-zero liquid + StakeHub/deposit activity) | **WARN** (v1 behavior; v2 should reduce) |
| RPC archive unavailable, cannot query `Height` | **WARN** (`note: rpc archive unavailable`) |
| `(coin,network)` has no native/staking handler yet | **UNVERIFIABLE** (aggregate row count by coin\|network, same as v1) |

**Coverage definition adjustment**:

```text
coverage = (# HotCold rows with full totalAccounted) / (# HotCold rows total)
```

"v1's 4%" is because only `ETH|ETH` + `BNB|BSC` are supported and only liquid is checked; v2 goal is to raise **observation depth** from liquid-only to liquid+stake on the same two network pairs, not expand coin types.

#### 5.3.6 PR01JUN26 three addresses (appendix-level empirical evidence)

Snapshot `PR01JUN26`, BSC `#101590091`, ETH `#25218797` (2026-05-31 23:59:59 UTC).

| Address | CSV | liquid | Staking explanation | v1 | v2 expected |
|---|---|---|---|---|---|
| `0xbf83…ed7a` | 8,657,966 BNB | 0.087 BNB | BNB delegated/claimed out of Stake Hub; EOA near zero long-term | FAIL | PASS (staked+liquid or WARN→PASS) |
| `0x32e11…1149` | 80,000 ETH | 0.017 ETH | 80k fully Eth2 deposited (40×2k) | FAIL | PASS (staked=80k) |
| `0x86523…08d96` | 172,798 BNB | 90,932 BNB | Stake Hub pooled ≈ **75,466 BNB** (30/53 validators scanned); total ≈ **166,398 BNB**, residual ≈ **6,400 BNB** pending full scan | FAIL | PASS or minor FAIL |

> Numbers above from bootstrap CLI multi-RPC cross-check (2026-06); can serve as golden fixtures when writing verification bundles.

#### 5.3.7 Version roadmap

| Version | Scope | Milestone |
|---|---|---|
| `@1.0` (current) | `ETH\|ETH`, `BNB\|BSC` liquid only | W2 bootstrap; known false-positive FAILs |
| `@2.0` | + BSC Stake Hub (pooled + unbonding) | Eliminate BNB staking false positives; requires archive RPC + validator cache |
| `@2.1` | + ETH beacon deposits | Eliminate Eth2 deposit false positives |
| `@3.0` | + ERC20 `balanceOf` (USDT/USDC/…) | Significantly raise HotCold value coverage |

`onchain-balance-hot@2.x` **does not replace** verifier id semantics — still row-by-row on-chain audit of HotCold; only extends `components` and RPC dependencies. After version bump, old snapshots should be **re-run** with new version and diff compared (§4 `version` rule).

---

## 6. On-Chain Anchoring of Key Data

So that ardmere's own "honesty" is also verifiable (not "trust me bro"), all **root-level evidence** is anchored on-chain.

> **Decision record**: This section follows [ADR-003](./decisions.md#adr-003-cold-storage-no-arweave-ipfs-on-chain-digest-only) — no separate cold storage (Arweave / IPFS); all raw data stays in our S3 + Postgres; chain writes digest roots only.

### 6.1 Anchored objects and frequency

| Anchored object | Digest content | On-chain frequency | Size |
|---|---|---|---|
| **Single-period snapshot anchor** | One tx carries `artifactBundleRoot` + `verificationBundleRoot` Merkle roots, plus exchange self-reported root and verifier verdict summary | After all verifiers finish for each new snapshot (~1×/month) | 2×32 B (roots) + metadata |
| **Weekly heartbeat** | Cumulative root of all artifacts / verifications that week + our public key | Every Monday 00:00 UTC | 64 B |

> What goes on-chain is not raw bulk data but **digest Merkle roots** — tamper-evident and cheap. Raw data is served via ardmere API (`GET /artifacts/:sha256`); users self-verify against on-chain roots.
>
> **Design decision (v2)**: artifact and verification roots **merge into 1 on-chain transaction**, avoiding frontend/index complexity from two independent anchors per `snapshotId`, with negligible gas cost change (~$0.001/tx on Base).

### 6.2 Chain selection

**Primary anchor chain: Base** (OP Stack L2, Chain ID `8453`) — [ADR-002](./decisions.md#adr-002-primary-anchor-chain-base)

Rationale:
- Extremely low gas (~$0.001 per tx), supports frequent heartbeat
- EVM-compatible, mature toolchain
- Data availability inherited from Ethereum L1 (via batch posting)
- Multiple public RPC providers (not tightly coupled to §7 verification RPC)

**Future option**: Bitcoin OP_RETURN (strongest but expensive) as annual audit anchor — out of MVP scope.

### 6.3 Anchor contract (minimal)

```solidity
// Deployed on Base mainnet (chain id 8453)
contract ArdmerePoRAnchor {
    uint8 public constant SCHEMA_VERSION = 2;

    event SnapshotAnchored(
        bytes32 indexed snapshotId,           // keccak256("PR01APR26")
        bytes32 indexed exchangeTag,          // keccak256("binance") — filter by exchange
        string exchange,                      // plaintext, e.g. "binance" / "okx" (BaseScan readable)
        uint32 periodSeq,                     // period number for this exchange (1 = earliest)
        uint64 snapshotTime,                  // exchange snapshot UTC unix
        uint32 btcBlockHeight,                // BTC time anchor
        bytes32 exchangeMerkleRoot,           // exchange self-reported Merkle root
        bytes32 artifactBundleRoot,           // ardmere artifact bundle Merkle root
        bytes32 verificationBundleRoot,       // ardmere verification bundle Merkle root
        uint8 verdictSummary,                 // verifier verdict bitfield
        uint16 coverageBps,                   // on-chain audit coverage × 10_000
        uint8 schemaVersion,
        uint256 anchoredAt
    );

    function anchorSnapshot(
        string calldata exchange,
        bytes32 snapshotId,
        uint32 periodSeq,
        uint64 snapshotTime,
        uint32 btcBlockHeight,
        bytes32 exchangeMerkleRoot,
        bytes32 artifactBundleRoot,
        bytes32 verificationBundleRoot,
        uint8 verdictSummary,
        uint16 coverageBps
    ) external;
}
```

> Minimal trusted kernel — permissioned `signer` plus **immutable per-snapshot storage** (schema v3). Each `anchorSnapshot` writes a `SnapshotRecord` and emits `SnapshotAnchored`. The contract is **UUPS upgradeable** so the anchor schema can evolve without migrating the proxy address. Neither `artifactBundleRoot` nor `verificationBundleRoot` may be zero. See §6.5 for query API.

**Deployed instance (Ethereum Sepolia testnet)** — see [`deployments.md`](./deployments.md) for proxy + implementation addresses.

### 6.4 Verification loop closure

```text
How any third party independently verifies ardmere:

1. Obtain snapshotId (e.g. PR01APR26)
2. On Base, query ArdmerePoRAnchor SnapshotAnchored event → get artifactBundleRoot + verificationBundleRoot
3. Fetch raw bundles from ardmere API (GET /artifacts/:sha256, GET /verifications/:id)
4. Recompute Merkle roots for both bundles; compare with on-chain fields → ✅ data not tampered
5. Re-run any verifier (open-source on GitHub + matching version) → conclusion matches bundle → ✅ logic replayable
```

> **Acknowledged trade-off**: If ardmere service goes offline, raw verification data becomes inaccessible; but **on-chain roots + any bundle cached by third parties** still enable post-hoc verification. Stronger guarantees can add independent Arweave mirror service in V2 (anchor roots unchanged).

### 6.5 On-chain query API (schema v3)

Schema v2 anchors were **event-only** — history lives in `eth_getLogs`. Schema v3 adds **immutable on-chain storage** per snapshot while keeping the same `SnapshotAnchored` event shape (only `schemaVersion` bumps to `3`).

#### Storage model

All snapshot state lives in an ERC-7201 namespaced slot (`ardmere.storage.ArdmerePoRAnchor`) so UUPS upgrades do not collide with OpenZeppelin `Ownable` / `Initializable` namespaces or the legacy `signer` slot.

| Mapping | Purpose |
|---|---|
| `snapshots[snapshotId]` | Full `SnapshotRecord` (all anchor fields + `anchoredAt`) |
| `snapshotIds[exchangeTag]` | Append-only chronological id list |
| `snapshotByPeriod[exchangeTag][periodSeq]` | O(1) lookup by exchange period |
| `latestSnapshotId[exchangeTag]` | Most recent anchor for dashboards |
| `snapshotCount[exchangeTag]` | Total anchors per exchange |

`anchorSnapshot` **reverts** if `snapshotId` already exists or if `periodSeq` is already anchored for that exchange — records are immutable once written.

#### Query functions

| Function | Use case |
|---|---|
| `snapshotExists(snapshotId)` | Cheap existence check before fetching bundles off-chain |
| `getSnapshot(snapshotId)` | Full record by native audit id hash |
| `getLatestSnapshot(exchangeTag)` | “Current” PoR verdict for an exchange widget |
| `getSnapshotByPeriod(exchangeTag, periodSeq)` | Historical period without scanning logs |
| `getSnapshotCount(exchangeTag)` | Pagination / progress |
| `getSnapshotIdAt(exchangeTag, index)` | Walk history in anchor order |

`exchangeTag = keccak256(bytes(exchange))` — e.g. `cast keccak "binance"`.

#### Upgrade path (v2 → v3)

1. Deploy new implementation (`forge script script/Upgrade.s.sol:Upgrade`).
2. Owner calls `upgradeToAndCall` on the existing proxy (see [`contracts/README.md`](../contracts/README.md)).
3. `signer`, `owner`, and pre-upgrade event history are unchanged.
4. **Pre-v3 anchors remain log-only** — storage is empty for those `snapshotId`s until re-anchored (not planned). New `anchorSnapshot` calls populate storage + emit events.

`STORAGE_VERSION` (currently `1`) tracks the on-chain layout independently of `SCHEMA_VERSION`.

#### Gas tradeoffs

| Approach | Read cost | Write cost | Notes |
|---|---|---|---|
| Events only (v2) | High off-chain (`eth_getLogs`) | ~lowest | No contract storage growth |
| Storage + events (v3) | O(1) `eth_call` | +~200k gas / anchor (string + 5 mappings) | Better UX for wallets / block explorers |

We keep emitting events so indexers and v2 integrations keep working.

#### Backward compatibility

- Event topic layout and field order are **unchanged**; `schemaVersion` in the event becomes `3` for new anchors.
- v2 log queries still return historical anchors.
- v3 view functions return `SnapshotNotFound` for snapshots anchored before the upgrade.

See also [`anchor-query-api.md`](./anchor-query-api.md) for `cast call` examples.

---

## 7. Public Node Pool Strategy (No Self-Hosted Archive)

**Design constraint**: All on-chain queries use third-party public RPC — **no self-hosted Erigon / electrs / archive nodes**.

### 7.1 Node pool abstraction

```ts
interface ChainRpcPool {
  // Select endpoint by chain; automatic failover + rate limiting + cache
  getBalance(network: Network, address: string, blockHeight: number): Promise<bigint>;
  getErc20Balance(network: Network, token: string, address: string, blockHeight: number): Promise<bigint>;
  getBlockMeta(network: Network, height: number): Promise<{ hash: string; timestamp: number }>;
}
```

### 7.2 Provider selection by chain

| Network | Primary provider | Fallback 1 | Fallback 2 | Archive support |
|---|---|---|---|---|
| Ethereum | Ankr Public | Llamanodes | publicnode.com | ✅ All supported |
| BSC | publicnode.com | bsc-dataseed.binance.org | Ankr | ✅ |
| Polygon | Ankr | publicnode | Polygon Foundation RPC | ✅ |
| Arbitrum | Ankr | Arbitrum Foundation | publicnode | ✅ |
| Base | Base public RPC | Ankr | publicnode | ✅ |
| BTC | mempool.space API | blockstream.info API | btc.getblock.io | ✅ |
| Solana | Helius free tier | publicnode | Triton | ⚠ Epoch-limited |
| Tron | TronGrid | Nile public | publicnode | ⚠ Historical balance API limited |
| Algorand | AlgoNode | PureStake free | Algorand Foundation | ✅ indexer |
| Aptos | Aptos Labs full node | publicnode | - | ✅ |

> **Full list**: Maintained in `config/rpc-providers.json`, grouped by `(network, capability)`; each entry records `(url, weight, rateLimit, lastFailedAt, costModel)`.

### 7.3 Node pool policies

| Policy | Implementation |
|---|---|
| **Failover** | Primary 5xx / timeout / rate limit → auto-switch to fallback; no retry for 5 minutes |
| **Rate limiting** | Token bucket per provider quota; overflow degrades to next provider |
| **Result cache** | `(network, address, blockHeight)` is invariant → permanent Postgres cache (high hit rate — same snapshot is not recomputed) |
| **Redundant verification** | Critical queries (e.g. large HotCold addresses): **query 2 providers in parallel; mismatch → mark `WARN`** |
| **Observability** | Real-time dashboard per provider success rate / latency / rate-limit count; monthly pool review |
| **Sandbox validation** | On deploy, run baseline tests (e.g. query known BTC genesis block balance) to confirm provider behavior |

### 7.4 Public node limitations and mitigations

| Limitation | Mitigation |
|---|---|
| **Historical balance query limits** (some free tiers only last 128 blocks) | Annotate `archiveDepth` in `rpc-providers.json`; auto-select archive-capable provider when exceeded |
| **BSC snapshot block archive scarcity** (PR01JUN26 `#101590091`: publicnode returns `historical state not available`; only some providers like drpc work, prone to 429) | Bind Stake Hub queries and `eth_getBalance` to same archive provider; permanently cache results `(network, address, height, component)`; cache StakeHub validator list by height |
| **Low QPS** (public nodes typically 5~20 QPS) | Deposit full audit runs slowly — hours OK; HotCold prioritized |
| **IP rate limiting** | Multi-provider parallel + egress IP pool (optional Cloudflare Workers distributed calls) |
| **Result tampering** | Dual-provider redundant comparison (§7.3) |
| **Solana / Tron historical API limits** | These chains start with `coverage < 100%`; UI labels "this network only checks recent balance" |

### 7.5 Are public nodes alone sufficient?

| Scenario | Public nodes sufficient? |
|---|---|
| HotCold 100% verification (10³ addresses × multi-chain) | ✅ One night |
| Deposit value sampling 99% coverage (~10⁵ addresses) | ✅ One week |
| Deposit full 100% (10⁷ addresses) | ⚠ Public nodes overwhelmed — months of queue or partial paid tier |

> Our SLA commits to "99% value coverage," not "100% address coverage." This is a **rational, achievable** target under the public-node strategy.

---

## 8. Physical Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│  Scheduler (cron)                                                    │
│  ├── Poll BAPI snapshot list every 6h (discover new snapshots)       │
│  └── Weekly heartbeat on-chain every Monday 00:00 UTC                │
└──────────────┬──────────────────────────────────────────────────────┘
               ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Fetcher Pool                                                        │
│  ├── BinanceBapiFetcher       → Artifact{kind:"bapiSnapshot"}        │
│  ├── BinanceWalletZipFetcher  → Artifact{kind:"walletZip"}           │
│  ├── BtcAnchorFetcher         → Artifact{kind:"btcBlockMeta"}        │
│  └── (future) AddressSigFetcher / ProofCsvFetcher / AttestationFetcher│
└──────────────┬──────────────────────────────────────────────────────┘
               ▼  Write to Artifact Store + our signature
┌─────────────────────────────────────────────────────────────────────┐
│  Artifact Store                                                      │
│  ├── Content-addressed: S3 (only location, ADR-003)                  │
│  └── Metadata: Postgres                                              │
└──────────────┬──────────────────────────────────────────────────────┘
               ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Verifier Runtime                                                    │
│  ├── Load registered verifiers (including placeholder stubs)         │
│  ├── Schedule: prepare ctx per verifier.requires                     │
│  ├── Execute: pure function verify(), access RpcPool (§7)            │
│  └── Persist: Verification + our signature                           │
└──────────────┬──────────────────────────────────────────────────────┘
               ▼  artifact + verification bundle roots → on-chain (1 tx)
┌─────────────────────────────────────────────────────────────────────┐
│  On-chain Anchor (Base, chain id 8453)                              │
│  └── ArdmerePoRAnchor.anchorSnapshot(exchange, snapshotId, …)       │
└──────────────┬──────────────────────────────────────────────────────┘
               ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Public API + static site                                              │
│  GET /snapshots                  → history + current                 │
│  GET /snapshots/:id              → all verification results for period│
│  GET /verifications/:id          → single verification + replay evidence + on-chain tx│
│  GET /artifacts/:sha256          → raw file (our signature + on-chain anchor)│
│  GET /anchors                    → on-chain anchor history (incl. Arweave refs)│
│  POST /user-inclusion/verify     → user zip in-browser WASM verify (L3)│
└─────────────────────────────────────────────────────────────────────┘
```

Suggested tech stack:

| Layer | Choice | Rationale |
|---|---|---|
| Backend language | **Go** | Reuse `binance/zkmerkle-proof-of-solvency` verifier package directly; concurrency-friendly |
| Database | Postgres | Clear relations; strong artifact ↔ verification linkage |
| Object storage | S3 (only location) | ADR-003 simplifies architecture |
| On-chain anchor | Base + Foundry | Cheap gas + mature toolchain |
| Signing | ed25519 (off-chain) + secp256k1 (on-chain) | Dual-key separation, tiered key management |
| Frontend | Existing ardmere.org static site + React island for interactive verification view | Consistent with [`DESIGN.md`](./DESIGN.md) terminal style |
| Deployment | GitHub Actions + Fly.io / Railway (backend) + GitHub Pages (frontend) | Zero ops |

---

## 9. Frontend Presentation Model

Make "verification service trustworthiness" the top visual principle. What end users see:

```
┌─ Snapshot PR01APR26 (2026-04-01 UTC, BTC#943129) ──────────────────┐
│                                                                     │
│  ardmere verdict:    ✅ Reserves Verified (with caveats)            │
│  Coverage:           ████████░░  6 / 10 dimensions                  │
│                                                                     │
│  ✅ Artifact Integrity         sha256 ✓, anchored on Base #...   [↗]│
│  ✅ Internal Consistency       14 coins reconcile                [↗]│
│  ✅ BTC Anchor                 #943129 @ 2026-04-01T00:01 UTC    [↗]│
│  ⚠ Solvency Claim              Self-reported by Binance, ratio≥1[↗]│
│  ✅ On-chain Balance (Hot)     100% verified, $42.3B             [↗]│
│  ⏳ On-chain Balance (Deposit) 87.4% by value (in progress…)     [↗]│
│  ❌ Address Ownership          NOT VERIFIABLE                    [?]│
│      └─ No public download channel                                  │
│  ❌ Global zk-SNARK Proof      NOT VERIFIABLE                    [?]│
│      └─ No public download channel                                  │
│  ❌ Third-Party Attestation    NOT VERIFIABLE                    [?]│
│      └─ No public attestation report available                      │
│  ❌ Wrapped Token Reconciliation  NOT VERIFIABLE                 [?]│
│      └─ Reconciliation rules pending                                │
│                                                                     │
│  Verified by ardmere @ 2026-05-05T14:32 UTC                         │
│  Bundle root:    0xa3f2…  →  Base tx 0x7c1e…                       │
│  Signing key:    ed25519:abc…  (verify via ardmere.org/.well-known) │
└─────────────────────────────────────────────────────────────────────┘
```

Each icon is clickable:

- ✅ → Expand finding list + on-chain anchor tx
- ⚠ → Explain why warn ("self-reported," etc.), not fail
- ❌ → Expand "why unverifiable + what events we're tracking for activation"
- ⏳ → Show live coverage percentage + estimated completion

---

## 10. Extension Scenario Walkthroughs

### Case A: Future Binance PoR announcement includes address signatures

```text
1. Write fetcher: parse announcement → download sig files → persist Artifact{kind:"addressSig"}
2. Upgrade verifier: address-ownership@0 → address-ownership@1
   - requires: ["walletZip", "addressSig"]
   - verify(): verify each sig with secp256k1/ed25519
3. Register and activate; frontend auto-updates ❌ → ✅
4. Historical snapshots re-run automatically (version upgrade); if no sig artifact, output PARTIAL
5. New conclusions anchored on-chain; establish "activation timeline"
```

**Core code change = one new fetcher + one new verifier + one registration line**.

### Case B: Future publicly downloadable third-party PoR attestation

```text
1. Write fetcher: periodically pull attestation PDF / JSON from third-party site
2. Write verifier: third-party-attestation@1
   - Verify signer is on allowlist
   - Verify attestation Merkle Root matches our BAPI-fetched root
3. Register and activate
```

### Case C: Future Binance public global zk proof

```text
1. Fetcher grabs proof.csv / vk
2. global-zk-proof@1 uses binance/zkmerkle-proof-of-solvency Go verifier
   (can compile to WASM at edge, isolated from §8 main process)
3. Register and activate
```

### Case D: Add new exchange (OKX / Coinbase)

```text
1. Abstract Exchange dimension:
   - Snapshot gains exchange field
   - Each Exchange has its own Fetcher set + default Verifier set
2. Verifier interface unchanged (only reads artifact + claim), so most verifiers reusable cross-exchange:
   - internal-consistency fully generic
   - onchain-balance-* fully generic
   - address-ownership: per-exchange signature formats → one version each
3. Frontend tabs by Exchange; compare public verifier ordering
```

---

## 11. Security and Threat Model

| Threat | Mitigation |
|---|---|
| Binance temporarily swaps source data (CDN poisoning) | sha256 + on-chain anchor; post-hoc comparison |
| Public RPC node lies | Dual-provider redundant comparison; mismatch → WARN |
| We are compromised and tamper conclusions | Hardware-isolated keys + on-chain anchor irreversible; cannot backfill history after tamper |
| Verifier logic bug | Forced versioning + periodic historical snapshot regression re-runs |
| Scraping blocked | Multi-entry fetcher + client browser self-service zip verification (L3) as backup |
| User privacy (zip upload) | In-browser WASM fully local verification; zero upload |

Signing key tiers:

| Purpose | Key | Storage |
|---|---|---|
| Artifact / Verification offline sign | ed25519 #1 | In service process (rotate quarterly) |
| On-chain anchor (contract call) | secp256k1 #2 | KMS / HSM; anchor process only |
| Site identity / domain DNSSEC | secp256k1 #3 | Offline; rotate annually |

---

## 12. MVP Plan

| Week | Work | Deliverable |
|---|---|---|
| W1 | Domain model + Verifier interface + Fetcher (BAPI + walletZip) + Artifact Store + 4 core verifiers (`integrity` / `internal-consistency` / `btc-anchor` / `solvency-claim`) + 4 placeholder stubs | CLI works: `ardmere verify PR01APR26` outputs 8-dimension verdict |
| W2 | `onchain-balance-hot@1` (liquid native) + Anchor contract + first anchor tx; **@2.0 staking-aware design** see §5.3 | Bootstrap works; BNB/ETH staking false positives eliminated in @2.0 |
| W3 | `onchain-balance-deposit@1` (sampling + continuous run) + Public API + frontend verdict card | Site live with 6/10 dimension verdict |
| W4 | In-browser WASM user inclusion verification (L3 self-service) + documentation + launch | Full MVP, public release |

> `onchain-balance-deposit` full run takes weeks — but W3 ships "live progress bar" — **do not wait for completion before launch**.

---

## 13. Trust Accumulation Path

ardmere's core asset = a credible record of "under public data constraints, we did the strictest verification possible."

| Phase | Trust source |
|---|---|
| Month 1 | Launch + first anchor — proves the system is live |
| Month 3 | Three consecutive periods with consistent verdict; frontend comparison charts |
| Month 6 | When Binance's own ratio fluctuates, our L1+L2 reproduce and publish first — establish credibility as "independent voice from Binance" |
| Year 1 | Integrate OKX / Coinbase / Bitget; become PoR industry reference verifier |
| Long term | When Binance finally publishes signatures / global proof, immediately activate corresponding verifiers; highlight "activation timeline" on frontend |

---

## 14. Answering the Core Question

> **How to preserve extensibility for information we cannot verify today?**

> **Model "cannot verify" as its own kind of verification** — every future dimension exists today as a "placeholder Verifier" in the registry, in API responses, and on the frontend UI. When external conditions mature (Binance official download channel, third-party publicly downloadable attestation, …), **upgrade that Verifier from v0 stub to v1 real implementation** — no upstream schema changes, no frontend code changes, no database redesign — all historical snapshots automatically re-run with the new Verifier and compare.

> Correspondingly, **anchor all key digests on-chain** so ardmere's own honesty is independently verifiable; **full public node pool strategy** keeps infrastructure zero-trust, zero self-operation — avoiding the "you can only trust us" paradox.

---

## Appendix A — Relationship to Existing ardmere Site

This service's frontend module embeds into ardmere.org as described in [`DESIGN.md`](./DESIGN.md), as a new section:

- Visual continuity with terminal / matrix / glitch system
- Verdict cards as "scanning" ASCII progress bars
- Each ❌ UNVERIFIABLE expand uses pseudo-command output in `cat /dev/binance | grep signature` style
- Colors reuse `--green / --cyan / --red`; no new palette

## Appendix B — Decision Records

See [`decisions.md`](./decisions.md); five ADRs finalized 2026-05-05:

| ADR | Decision |
|---|---|
| ADR-001 Backend language | **Go** |
| ADR-002 Primary anchor chain | **Base** (chain id 8453) |
| ADR-003 Cold storage | **No separate cold storage; on-chain digests only** |
| ADR-004 User verification form | **Fully local in-browser WASM** |
| ADR-005 Multi-exchange roadmap | **OKX prioritized after MVP** |

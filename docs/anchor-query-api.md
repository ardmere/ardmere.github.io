# ArdmerePoRAnchor On-Chain Query Guide

**User manual** for third-party integrators, wallets, block explorers, and independent auditors: how to read exchange PoR (Proof of Reserves) anchor data from the `ArdmerePoRAnchor` contract.

Related docs: [verifier-architecture.md §6.5](./verifier-architecture.md#65-on-chain-query-api-schema-v3) (architecture background) · [deployments.md](./deployments.md) (deployment history) · [por-cli.md](./por-cli.md) (`por anchor` operations)

---

## Table of Contents

1. [Overview: On-Chain vs Off-Chain Verification](#1-overview-on-chain-vs-off-chain-verification)
2. [Contract Address and Network](#2-contract-address-and-network)
3. [Identifiers: `exchangeTag` and `snapshotId`](#3-identifiers-exchangetag-and-snapshotid)
4. [`SnapshotRecord` Field Reference](#4-snapshotrecord-field-reference)
5. [Constants and Public State](#5-constants-and-public-state)
6. [View Function Reference](#6-view-function-reference)
7. [Practical Examples: `cast` and `ethers.js`](#7-practical-examples-cast-and-ethersjs)
8. [Common Query Patterns](#8-common-query-patterns)
9. [Pre-v3 Anchors: Event Logs Only](#9-pre-v3-anchors-event-logs-only)
10. [On-Chain Query + Off-Chain Bundle Verification Loop](#10-on-chain-query--off-chain-bundle-verification-loop)
11. [Error Codes and Troubleshooting](#11-error-codes-and-troubleshooting)
12. [Related Links](#12-related-links)

---

## 1. Overview: On-Chain vs Off-Chain Verification

`ArdmerePoRAnchor` is ardmere's **minimal trusted kernel**: each exchange snapshot period, an authorized `signer` calls `anchorSnapshot` to write the **digest roots** of a PoR audit on-chain.

| Capability | On-chain (this contract) | Off-chain (ardmere repo / third-party cache) |
| --- | --- | --- |
| Snapshot existence | ✅ `snapshotExists` / events | — |
| Exchange name, period, time anchor | ✅ `SnapshotRecord` | assessment JSON |
| artifact / verification **Merkle root** | ✅ fields + events | recompute bundle and compare |
| Exchange self-reported Merkle root | ✅ `exchangeMerkleRoot` | exchange original disclosure |
| Verification verdict summary | ✅ `verdictSummary`, `coverageBps` | full verifier output |
| Raw artifact files, proof data | ❌ | [artifact archive](./reports/artifact-archive-index.md) |
| Replay verifier logic | ❌ | `por verify` / open-source verifier |

**Schema v3** (current implementation) adds **immutable on-chain storage** while retaining the `SnapshotAnchored` event, enabling O(1) reads via `eth_call` without scanning historical logs.

**Schema v2 and earlier** anchors exist only in event logs; v3 view functions return `SnapshotNotFound` for those `snapshotId` values. See [§9](#9-pre-v3-anchors-event-logs-only).

---

## 2. Contract Address and Network

| Item | Value |
| --- | --- |
| **Network** | Ethereum **Sepolia** testnet |
| **Chain ID** | `11155111` |
| **Proxy contract (user-facing address)** | [`0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7`](https://sepolia.etherscan.io/address/0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7#code) |
| **Implementation v3** | [`0x62050405283f222EFEdF3BF0d6Cb16541cd1327a`](https://sepolia.etherscan.io/address/0x62050405283f222EFEdF3BF0d6Cb16541cd1327a#code) |
| **Upgrade pattern** | UUPS (`ERC1967Proxy` + upgradeable implementation) |
| **`SCHEMA_VERSION`** | `3` |
| **`STORAGE_VERSION`** | `1` |
| **Anchor signer** | `0xf11AcEcBB54bf72Db69f6BaC4f16FC6491cC670F` |

> Always issue `eth_call` and read events against the **proxy address**, not the implementation address directly.

Full deployment timeline, upgrade transactions, and deprecated instances are in [deployments.md](./deployments.md).

### Environment Variables (example)

```bash
export ANCHOR_CONTRACT=0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7
export RPC=sepolia                    # or export RPC_URL=https://...
export CHAIN_ID=11155111
```

---

## 3. Identifiers: `exchangeTag` and `snapshotId`

The contract does **not** accept plaintext exchange names or audit ids as mapping keys; compute keccak256 before querying:

| Name | Computation | Example input | Example output |
| --- | --- | --- | --- |
| **`exchangeTag`** | `keccak256(bytes(exchange))` | `"binance"` | `0x099e8cf4d817a6e4eec62bff4cdef05faa4a00fcde8d7e99f5090708d23ad9b2` |
| **`exchangeTag`** | same | `"okx"` | `0xd6b6d5e0aacce0469a313983d889ed10d0bb7c9545af0285a19b4ff094b4041d` |
| **`snapshotId`** | `keccak256(bytes(auditId))` | `"PR01JUN26"` | `0x8eebf94be87651156d5d4b7eed3d99daed7fe308ba9ab0771cd5f6efb417c638` |
| **`snapshotId`** | same | `"506872725"` | `0x88875070b19fd65cafbdd4674a55d835e0f24a96091375b001aef45e4248b2f0` |

```bash
# Foundry cast
export EXCHANGE_TAG=$(cast keccak "binance")
export OKX_TAG=$(cast keccak "okx")
export SNAPSHOT_ID=$(cast keccak "PR01JUN26")
export OKX_SNAPSHOT_ID=$(cast keccak "506872725")

# Verify
cast keccak "binance"
# 0x099e8cf4d817a6e4eec62bff4cdef05faa4a00fcde8d7e99f5090708d23ad9b2
```

The `exchange` string is stored in **plaintext** in `SnapshotRecord` and events (for block explorer readability), but all exchange-aggregating view functions use **`exchangeTag` (bytes32)**.

`periodSeq` is a **1-based** period index (the Nth anchor for that exchange), independent of `snapshotId`.

---

## 4. `SnapshotRecord` Field Reference

Struct returned by `getSnapshot` / `getLatestSnapshot` / `getSnapshotByPeriod`:

```solidity
struct SnapshotRecord {
    bytes32 snapshotId;
    bytes32 exchangeTag;
    string  exchange;
    uint32  periodSeq;
    uint64  snapshotTime;
    uint32  btcBlockHeight;
    bytes32 exchangeMerkleRoot;
    bytes32 artifactBundleRoot;
    bytes32 verificationBundleRoot;
    uint8   verdictSummary;
    uint16  coverageBps;
    uint8   schemaVersion;
    uint256 anchoredAt;
}
```

| Field | Solidity type | Meaning |
| --- | --- | --- |
| `snapshotId` | `bytes32` | `keccak256(bytes(auditId))`, e.g. `keccak256("PR01JUN26")` |
| `exchangeTag` | `bytes32` | `keccak256(bytes(exchange))` |
| `exchange` | `string` | Plaintext exchange id, e.g. `"binance"`, `"okx"` |
| `periodSeq` | `uint32` | Anchor sequence for this exchange; **1 = earliest period** |
| `snapshotTime` | `uint64` | Exchange snapshot UTC **Unix timestamp** (seconds) |
| `btcBlockHeight` | `uint32` | BTC time-anchor block height; `0` if none |
| `exchangeMerkleRoot` | `bytes32` | Exchange self-reported liability/user Merkle root; may be `0x00…00` if unknown |
| `artifactBundleRoot` | `bytes32` | ardmere **artifact bundle** Merkle root (must be non-zero) |
| `verificationBundleRoot` | `bytes32` | ardmere **verification bundle** Merkle root (must be non-zero) |
| `verdictSummary` | `uint8` | Verifier conclusion **bitfield summary** (see table below) |
| `coverageBps` | `uint16` | On-chain audit coverage × 10,000 (e.g. `8500` ≈ 85.00%) |
| `schemaVersion` | `uint8` | `SCHEMA_VERSION` at write time (new anchors use `3`) |
| `anchoredAt` | `uint256` | `block.timestamp` at anchor time (seconds); **`0` means record does not exist** |

### `verdictSummary` Bitfield (ardmere pipeline)

| Bit | Verifier ID | Meaning |
| --- | --- | --- |
| bit 0 | `internal-consistency` | Set to 1 if passed |
| bit 1 | `onchain-balance-hot` | Set to 1 if passed |
| bit 2 | `solvency-claim` | Set to 1 if passed |
| bit 3 | (derived) | Automatically set to 1 when bit 0 is 1 |

Full verifier matrix and Stage determination: [verifier-architecture.md](./verifier-architecture.md).

---

## 5. Constants and Public State

| Function / variable | Signature | Return value | Revert |
| --- | --- | --- | --- |
| `SCHEMA_VERSION` | `SCHEMA_VERSION()(uint8)` | `3` | — |
| `STORAGE_VERSION` | `STORAGE_VERSION()(uint8)` | `1` | — |
| `signer` | `signer()(address)` | Current authorized anchor address | — |

```bash
cast call $ANCHOR_CONTRACT "SCHEMA_VERSION()(uint8)" --rpc-url $RPC
cast call $ANCHOR_CONTRACT "STORAGE_VERSION()(uint8)" --rpc-url $RPC
cast call $ANCHOR_CONTRACT "signer()(address)" --rpc-url $RPC
```

---

## 6. View Function Reference

All functions below are `view` and must be called via the proxy. The write function `anchorSnapshot` is restricted to `signer` and is **out of scope for this query guide**.

### `snapshotExists`

```solidity
function snapshotExists(bytes32 snapshotId) external view returns (bool);
```

| Parameter | Description |
| --- | --- |
| `snapshotId` | `keccak256(bytes(auditId))` |

| Return | Description |
| --- | --- |
| `bool` | `true` if and only if the id has an on-chain storage record (`anchoredAt != 0`) |

| Revert | None |

---

### `getSnapshot`

```solidity
function getSnapshot(bytes32 snapshotId) external view returns (SnapshotRecord memory);
```

| Parameter | Description |
| --- | --- |
| `snapshotId` | Target snapshot id |

| Return | Full `SnapshotRecord` |

| Revert | `SnapshotNotFound(bytes32 snapshotId)` — record does not exist (includes pre-v3 event-only anchors) |

---

### `getLatestSnapshot`

```solidity
function getLatestSnapshot(bytes32 exchangeTag) external view returns (SnapshotRecord memory);
```

| Parameter | Description |
| --- | --- |
| `exchangeTag` | `keccak256(bytes(exchange))` |

| Return | **Most recent** anchored `SnapshotRecord` for that exchange |

| Revert | `SnapshotNotFound(bytes32)` — this `exchangeTag` has never been anchored in v3 storage (`snapshotId == 0`) |

---

### `getSnapshotByPeriod`

```solidity
function getSnapshotByPeriod(bytes32 exchangeTag, uint32 periodSeq) external view returns (SnapshotRecord memory);
```

| Parameter | Description |
| --- | --- |
| `exchangeTag` | Exchange tag |
| `periodSeq` | Period index (≥ 1) |

| Return | `SnapshotRecord` for the specified period |

| Revert | `SnapshotNotFound(bytes32)` — that period was not anchored |

---

### `getSnapshotCount`

```solidity
function getSnapshotCount(bytes32 exchangeTag) external view returns (uint256);
```

| Parameter | `exchangeTag` |

| Return | Total anchors for that exchange in v3 storage; `0` if none |

| Revert | None |

---

### `getSnapshotIdAt`

```solidity
function getSnapshotIdAt(bytes32 exchangeTag, uint256 index) external view returns (bytes32);
```

| Parameter | Description |
| --- | --- |
| `exchangeTag` | Exchange tag |
| `index` | **0-based**; `0` = earliest period |

| Return | `snapshotId` at position `index` in chronological order |

| Revert | `SnapshotIndexOutOfBounds(bytes32 exchangeTag, uint256 index)` — `index >= getSnapshotCount(exchangeTag)` |

---

### Event `SnapshotAnchored` (emitted in sync with storage)

```solidity
event SnapshotAnchored(
    bytes32 indexed snapshotId,
    bytes32 indexed exchangeTag,
    string  exchange,
    uint32  periodSeq,
    uint64  snapshotTime,
    uint32  btcBlockHeight,
    bytes32 exchangeMerkleRoot,
    bytes32 artifactBundleRoot,
    bytes32 verificationBundleRoot,
    uint8   verdictSummary,
    uint16  coverageBps,
    uint8   schemaVersion,
    uint256 anchoredAt
);
```

- **Topic0** (event signature): `0x48e87afd431a5a42eecd1c689a4eeb946147458e85bb8ad30120047f73cd2de1`
- **Topic1**: `snapshotId` (indexed)
- **Topic2**: `exchangeTag` (indexed)

After v3, each `anchorSnapshot` writes storage and emits an event; v2 only emitted events.

---

## 7. Practical Examples: `cast` and `ethers.js`

### 7.1 Common Tuple Return Type

`cast call` return type for `SnapshotRecord` (matches ABI):

```text
(bytes32,bytes32,string,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16,uint8,uint256)
```

Referred to below as `$RECORD_TUPLE`.

```bash
RECORD_TUPLE="(bytes32,bytes32,string,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16,uint8,uint256)"
```

### 7.2 Binance Example (`PR01JUN26`, period 43)

```bash
export EXCHANGE_TAG=$(cast keccak "binance")
export SNAPSHOT_ID=$(cast keccak "PR01JUN26")

# Existence check
cast call $ANCHOR_CONTRACT "snapshotExists(bytes32)(bool)" $SNAPSHOT_ID --rpc-url $RPC

# Read by snapshotId
cast call $ANCHOR_CONTRACT \
  "getSnapshot(bytes32)${RECORD_TUPLE}" \
  $SNAPSHOT_ID --rpc-url $RPC

# Latest period for this exchange
cast call $ANCHOR_CONTRACT \
  "getLatestSnapshot(bytes32)${RECORD_TUPLE}" \
  $EXCHANGE_TAG --rpc-url $RPC

# By period (Binance PR01JUN26 = period 43)
cast call $ANCHOR_CONTRACT \
  "getSnapshotByPeriod(bytes32,uint32)${RECORD_TUPLE}" \
  $EXCHANGE_TAG 43 --rpc-url $RPC

# Historical enumeration
cast call $ANCHOR_CONTRACT "getSnapshotCount(bytes32)(uint256)" $EXCHANGE_TAG --rpc-url $RPC
cast call $ANCHOR_CONTRACT "getSnapshotIdAt(bytes32,uint256)(bytes32)" $EXCHANGE_TAG 0 --rpc-url $RPC
```

### 7.3 OKX Example (audit id `506872725`)

```bash
export OKX_TAG=$(cast keccak "okx")
export OKX_SNAPSHOT_ID=$(cast keccak "506872725")

cast call $ANCHOR_CONTRACT "snapshotExists(bytes32)(bool)" $OKX_SNAPSHOT_ID --rpc-url $RPC

cast call $ANCHOR_CONTRACT \
  "getSnapshot(bytes32)${RECORD_TUPLE}" \
  $OKX_SNAPSHOT_ID --rpc-url $RPC

cast call $ANCHOR_CONTRACT \
  "getLatestSnapshot(bytes32)${RECORD_TUPLE}" \
  $OKX_TAG --rpc-url $RPC
```

### 7.4 `ethers.js` v6 Snippet

```javascript
import { Contract, JsonRpcProvider, id, keccak256, toUtf8Bytes } from "ethers";

const ANCHOR = "0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7";
const RPC = process.env.RPC_URL ?? "https://rpc.sepolia.org";

const abi = [
  "function SCHEMA_VERSION() view returns (uint8)",
  "function snapshotExists(bytes32 snapshotId) view returns (bool)",
  "function getSnapshot(bytes32 snapshotId) view returns (tuple(bytes32 snapshotId, bytes32 exchangeTag, string exchange, uint32 periodSeq, uint64 snapshotTime, uint32 btcBlockHeight, bytes32 exchangeMerkleRoot, bytes32 artifactBundleRoot, bytes32 verificationBundleRoot, uint8 verdictSummary, uint16 coverageBps, uint8 schemaVersion, uint256 anchoredAt))",
  "function getLatestSnapshot(bytes32 exchangeTag) view returns (tuple(bytes32 snapshotId, bytes32 exchangeTag, string exchange, uint32 periodSeq, uint64 snapshotTime, uint32 btcBlockHeight, bytes32 exchangeMerkleRoot, bytes32 artifactBundleRoot, bytes32 verificationBundleRoot, uint8 verdictSummary, uint16 coverageBps, uint8 schemaVersion, uint256 anchoredAt))",
  "function getSnapshotByPeriod(bytes32 exchangeTag, uint32 periodSeq) view returns (tuple(bytes32 snapshotId, bytes32 exchangeTag, string exchange, uint32 periodSeq, uint64 snapshotTime, uint32 btcBlockHeight, bytes32 exchangeMerkleRoot, bytes32 artifactBundleRoot, bytes32 verificationBundleRoot, uint8 verdictSummary, uint16 coverageBps, uint8 schemaVersion, uint256 anchoredAt))",
  "function getSnapshotCount(bytes32 exchangeTag) view returns (uint256)",
  "function getSnapshotIdAt(bytes32 exchangeTag, uint256 index) view returns (bytes32)",
  "event SnapshotAnchored(bytes32 indexed snapshotId, bytes32 indexed exchangeTag, string exchange, uint32 periodSeq, uint64 snapshotTime, uint32 btcBlockHeight, bytes32 exchangeMerkleRoot, bytes32 artifactBundleRoot, bytes32 verificationBundleRoot, uint8 verdictSummary, uint16 coverageBps, uint8 schemaVersion, uint256 anchoredAt)",
];

const provider = new JsonRpcProvider(RPC, 11155111);
const anchor = new Contract(ANCHOR, abi, provider);

// keccak256(bytes("binance")) — matches cast keccak
const exchangeTag = keccak256(toUtf8Bytes("binance"));
const snapshotId = keccak256(toUtf8Bytes("PR01JUN26"));

const exists = await anchor.snapshotExists(snapshotId);
const record = await anchor.getSnapshot(snapshotId);
const latest = await anchor.getLatestSnapshot(exchangeTag);

// Enumerate all v3 anchors for an exchange
const count = await anchor.getSnapshotCount(exchangeTag);
const ids = [];
for (let i = 0n; i < count; i++) {
  ids.push(await anchor.getSnapshotIdAt(exchangeTag, i));
}

// Event filtering (covers full history including pre-v3)
const filter = anchor.filters.SnapshotAnchored(null, exchangeTag);
const logs = await anchor.queryFilter(filter, 11_182_041); // proxy deploy block, see deployments.md
```

> `id("PR01JUN26")` in ethers v6 is equivalent to `keccak256(toUtf8Bytes("PR01JUN26"))`, matching the contract's `keccak256(bytes(...))`.

---

## 8. Common Query Patterns

### 8.1 Check Whether an Audit Id Is Anchored (v3 storage)

```bash
cast call $ANCHOR_CONTRACT "snapshotExists(bytes32)(bool)" $(cast keccak "PR01JUN26") --rpc-url $RPC
```

- `true` → use `getSnapshot`
- `false` → may be unanchored, or v2 event-only anchor (see §9)

### 8.2 Read an Exchange's "Current" PoR Summary (widget / dashboard)

```bash
cast call $ANCHOR_CONTRACT \
  "getLatestSnapshot(bytes32)${RECORD_TUPLE}" \
  $(cast keccak "okx") --rpc-url $RPC
```

Suggested display fields: `exchange`, `periodSeq`, `snapshotTime`, `verdictSummary`, `coverageBps`, `artifactBundleRoot`, `verificationBundleRoot`.

### 8.3 Query by Historical Period (no log scan)

```bash
# Binance period 43
cast call $ANCHOR_CONTRACT \
  "getSnapshotByPeriod(bytes32,uint32)${RECORD_TUPLE}" \
  $(cast keccak "binance") 43 --rpc-url $RPC
```

### 8.4 List All v3 Anchored Snapshots for an Exchange

```bash
TAG=$(cast keccak "binance")
N=$(cast call $ANCHOR_CONTRACT "getSnapshotCount(bytes32)(uint256)" $TAG --rpc-url $RPC)
# In shell, loop 0 .. N-1 calling getSnapshotIdAt, then getSnapshot for each id
```

`getSnapshotIdAt` is ordered by **anchor chronology** (`0` = earliest), not by `periodSeq` value; use `getSnapshotByPeriod` if you need lookup by `periodSeq`.

### 8.5 Look Up Full Record from `snapshotId`

Given a known audit id string (e.g. from operator `por anchor` output):

```bash
SID=$(cast keccak "506872725")
cast call $ANCHOR_CONTRACT "getSnapshot(bytes32)${RECORD_TUPLE}" $SID --rpc-url $RPC
```

---

## 9. Pre-v3 Anchors: Event Logs Only

Before the 2026 v2→v3 upgrade, `anchorSnapshot` **only emitted `SnapshotAnchored` events** and did not write storage. Therefore:

- `snapshotExists(id)` → `false`
- `getSnapshot(id)` → `SnapshotNotFound`
- Data is still readable via **`eth_getLogs` / `cast logs`** from the same proxy address

### 9.1 `cast logs` Example

```bash
export DEPLOY_BLOCK=11182041   # Sepolia proxy deploy block, see deployments.md
export EXCHANGE_TAG=$(cast keccak "binance")

# Filter by exchange (topic2 = exchangeTag)
cast logs \
  --rpc-url $RPC \
  --address $ANCHOR_CONTRACT \
  --from-block $DEPLOY_BLOCK \
  "SnapshotAnchored(bytes32,bytes32,string,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16,uint8,uint256)" \
  --topic2 $EXCHANGE_TAG

# Filter by snapshotId (topic1)
cast logs \
  --rpc-url $RPC \
  --address $ANCHOR_CONTRACT \
  --from-block $DEPLOY_BLOCK \
  "SnapshotAnchored(bytes32,bytes32,string,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16,uint8,uint256)" \
  --topic1 $(cast keccak "PR01JUN26")
```

### 9.2 Distinguishing v2 vs v3 Anchors

| Characteristic | v2 (events only) | v3 (storage + events) |
| --- | --- | --- |
| `snapshotExists` | `false` | `true` |
| `schemaVersion` in event | `2` | `3` |
| Read method | `cast logs` / indexer | `getSnapshot*` views |

Historical v2 anchors are **not** backfilled into storage; new anchors automatically write to v3 storage.

---

## 10. On-Chain Query + Off-Chain Bundle Verification Loop

Typical flow for independent third-party verification (aligned with [verifier-architecture.md §6.4](./verifier-architecture.md#64-verification-loop-closure)):

```text
1. Determine auditId (e.g. PR01JUN26) → compute snapshotId
2. Read artifactBundleRoot, verificationBundleRoot on-chain
      · v3: getSnapshot(snapshotId)
      · v2: SnapshotAnchored event
3. Fetch bundle JSON from [artifact archive](./reports/artifact-archive-index.md) or local clone
4. Recompute both Merkle roots per ardmere rules and compare to on-chain fields
5. Run `por verify -snapshot <id>` or replay open-source verifier → conclusion should match verification bundle
```

Example (Binance PR01JUN26, from public assessment / archive):

| Root | On-chain field | Public archive reference |
| --- | --- | --- |
| artifact bundle | `artifactBundleRoot` | `0xf452e47dd22ed63dc4a905fe79da6c6f7a6975cc0d775f50d01879e97616671f` |
| verification bundle | `verificationBundleRoot` | `0x84e7f008460dcfa4ae4969c9a6e2a8b2585d5dcbe6a33cfead74b175782f6f42` |

Matching on-chain roots proves only that **ardmere's anchored digests have not been tampered with**; a full PoR conclusion still requires downloading artifacts and replaying the verifier. If ardmere services are offline, **on-chain roots + any third-party cached bundle** still enable retrospective verification.

---

## 11. Error Codes and Troubleshooting

Contract custom errors (view-related):

| Error | Trigger condition |
| --- | --- |
| `SnapshotNotFound(bytes32 snapshotId)` | Target does not exist for `getSnapshot` / `getLatestSnapshot` / `getSnapshotByPeriod` |
| `SnapshotIndexOutOfBounds(bytes32 exchangeTag, uint256 index)` | `index` out of range for `getSnapshotIdAt` |

Write-path errors (for reference; query callers typically do not encounter these): `Unauthorized`, `InvalidRoot`, `EmptyExchange`, `ZeroAddress`, `SnapshotAlreadyExists`, `PeriodAlreadyAnchored`.

| Symptom | Likely cause | Suggestion |
| --- | --- | --- |
| `SnapshotNotFound` but Etherscan shows a transaction | v2 event-only anchor | Use §9 `cast logs` |
| `getLatestSnapshot` fails | Exchange has no v3 storage record yet | Query events or wait for new anchor |
| `exchangeTag` mismatch | Used plaintext `"binance"` instead of `cast keccak` | Recompute tag |
| RPC returns empty | Wrong network / not synced | Confirm Sepolia, `CHAIN_ID=11155111` |

---

## 12. Related Links

| Resource | Link |
| --- | --- |
| Proxy (Sepolia Etherscan) | https://sepolia.etherscan.io/address/0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7#code |
| Implementation v3 (verified source) | https://sepolia.etherscan.io/address/0x62050405283f222EFEdF3BF0d6Cb16541cd1327a#code |
| Deployment and upgrade history | [deployments.md](./deployments.md) |
| Contract source | [`contracts/src/ArdmerePoRAnchor.sol`](../contracts/src/ArdmerePoRAnchor.sol) |
| CLI anchor workflow | [por-cli.md](./por-cli.md) |
| Evidence archive | [artifact-archive-index.md](./reports/artifact-archive-index.md) |
| Architecture and trust boundaries | [verifier-architecture.md](./verifier-architecture.md) |

---

*This document corresponds to implementation v3 (`SCHEMA_VERSION=3`, `STORAGE_VERSION=1`). After contract upgrades, rely on on-chain `SCHEMA_VERSION` / `STORAGE_VERSION` and the Etherscan verified ABI.*

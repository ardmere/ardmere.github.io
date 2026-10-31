# Query ArdmerePoRAnchor

How to read exchange PoR audit / anchor records from the on-chain `ArdmerePoRAnchor` contract. Use this when you want to confirm that a homepage report snapshot was anchored, or to pull the digest roots that bind an assessment.

Related: [Deployments](https://github.com/ardmere/ardmere.github.io/blob/main/docs/deployments.md) · [Full ABI / query API](https://github.com/ardmere/ardmere.github.io/blob/main/docs/anchor-query-api.md) · [Evidence archive](./reports/artifact-archive-index.md)

---

## What is on-chain

Each successful `anchorSnapshot` stores one **immutable** `SnapshotRecord` (schema v3) and emits `SnapshotAnchored`. The record holds exchange id, period, snapshot time, and the Merkle roots of the artifact and verification bundles—not the raw PoR files.

| You can read on-chain | You still need off-chain |
| --- | --- |
| Whether a snapshot id is anchored | Raw wallets, proofs, ZK bundles |
| Artifact / verification Merkle roots | Recomputing those roots from files |
| Compact `verdictSummary` / `coverageBps` | Full verifier matrix and Stage narrative |

Always call the **proxy** address (UUPS). Do not send reads to the implementation address.

---

## Network and proxy

Operable deployment today is **Ethereum Sepolia** only. Base mainnet is the product target but is not deployed yet—see [deployments.md](https://github.com/ardmere/ardmere.github.io/blob/main/docs/deployments.md).

Copy the live **Proxy (`ANCHOR_CONTRACT`)** row from that file (and the matching Etherscan / Sourcify links). Summary of the current testnet instance:

| Item | Value |
| --- | --- |
| Network | Ethereum Sepolia |
| Chain ID | see `CHAIN_ID` / deployments.md |
| Proxy (`ANCHOR_CONTRACT`) | see [deployments.md](https://github.com/ardmere/ardmere.github.io/blob/main/docs/deployments.md) — use the Proxy row |
| Implementation v3 | [`0x62050405283f222EFEdF3BF0d6Cb16541cd1327a`](https://sepolia.etherscan.io/address/0x62050405283f222EFEdF3BF0d6Cb16541cd1327a#code) |
| `SCHEMA_VERSION` | `3` |
| `STORAGE_VERSION` | `1` |
| Anchor signer | `0xf11AcEcBB54bf72Db69f6BaC4f16FC6491cC670F` |

```bash
# Proxy address: copy from docs/deployments.md (Proxy / ANCHOR_CONTRACT row)
export ANCHOR_CONTRACT=0x…   # Sepolia proxy
export RPC_URL=…             # any Sepolia HTTPS RPC
```

Sepolia anchors are **testnet evidence**. Treat them as reproducible digests, not mainnet production finality.

---

## Identifiers

| Name | How to compute | Homepage example |
| --- | --- | --- |
| `snapshotId` | `keccak256(bytes(auditId))` | audit id `502260945` (OKX) |
| `exchangeTag` | `keccak256(bytes(exchange))` | `"okx"`, `"binance"`, `"bingx"`, … |

```bash
cast keccak "502260945"
# 0xadc693dc11246c6b6472251c2b8aa6a79b2e84d12d9635db7b123b3049bd36c4

cast keccak "okx"
# 0xd6b6d5e0aacce0469a313983d889ed10d0bb7c9545af0285a19b4ff094b4041d
```

Prefer lookup by **`snapshotId` from the homepage Snapshot column** (same string as the report filename prefix). `periodSeq` on-chain may differ from the `periodSeq` field inside assessment JSON when the operator allocated the next free period for that exchange.

---

## Read methods (ABI)

Source: [`contracts/src/ArdmerePoRAnchor.sol`](https://github.com/ardmere/ardmere.github.io/blob/main/contracts/src/ArdmerePoRAnchor.sol). Full field notes: [anchor-query-api.md](https://github.com/ardmere/ardmere.github.io/blob/main/docs/anchor-query-api.md).

| Method | Returns | Notes |
| --- | --- | --- |
| `snapshotExists(bytes32 snapshotId)` | `bool` | `true` iff v3 storage has the id |
| `getSnapshot(bytes32 snapshotId)` | `SnapshotRecord` | Reverts `SnapshotNotFound` if missing |
| `getLatestSnapshot(bytes32 exchangeTag)` | `SnapshotRecord` | Most recently anchored for that exchange |
| `getSnapshotByPeriod(bytes32 exchangeTag, uint32 periodSeq)` | `SnapshotRecord` | By on-chain period index |
| `getSnapshotCount(bytes32 exchangeTag)` | `uint256` | Count of v3 anchors |
| `getSnapshotIdAt(bytes32 exchangeTag, uint256 index)` | `bytes32` | Chronological id; `0` = earliest |
| `SCHEMA_VERSION()` / `STORAGE_VERSION()` / `signer()` | constants / address | Sanity checks |

`SnapshotRecord` fields used most often when auditing a homepage report:

- `exchange`, `periodSeq`, `snapshotTime`, `btcBlockHeight`
- `artifactBundleRoot`, `verificationBundleRoot`
- `exchangeMerkleRoot`, `verdictSummary`, `coverageBps`, `anchoredAt`

Minimal ABI for `cast` / ethers:

```text
snapshotExists(bytes32)(bool)
getSnapshot(bytes32)((bytes32,bytes32,string,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16,uint8,uint256))
getLatestSnapshot(bytes32)((bytes32,bytes32,string,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16,uint8,uint256))
getSnapshotCount(bytes32)(uint256)
```

---

## Homepage snapshot examples

These audit ids match the current public homepage report set and are anchored on Sepolia (see the homepage table in [deployments.md](https://github.com/ardmere/ardmere.github.io/blob/main/docs/deployments.md)).

| Exchange | Homepage Snapshot (`auditId`) | `snapshotId` = `keccak256(auditId)` |
| --- | --- | --- |
| BingX | `20261015` | `0x770ca1dc46e92652a9eeb4887dcd7ee9d06d93e05e25abe55823609bdb534ffa` |
| MEXC | `202610` | `0xde6d943879fa24f6c45ad35a6dd1f1f4badcdb78acf82b7271ffa4ef94a78b90` |
| OKX | `502260945` | `0xadc693dc11246c6b6472251c2b8aa6a79b2e84d12d9635db7b123b3049bd36c4` |
| Binance | `PR01OCT26` | `0x119a99ffa3d507d3acf3beec2c5503ad2588d8cd2d663b061f7eaddb31ec5f80` |
| KuCoin | `20260930` | `0x723d99e9297f29529761b75bdf9e45d5743050e64aed2415424b2e7158cf4a84` |
| BingX | `20260915` | `0x0906c8899d82ae5ab245b26e34d6ea43275fc5aacd2ddbe74cae2f71af5e667a` |

### Check that a homepage row is anchored

```bash
# ANCHOR_CONTRACT = Proxy row from docs/deployments.md
export ANCHOR_CONTRACT=0x…
export RPC_URL=${RPC_URL:-https://rpc.sepolia.org}
RECORD_TUPLE="(bytes32,bytes32,string,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16,uint8,uint256)"

# OKX Stage 1 row on the homepage
cast call $ANCHOR_CONTRACT "snapshotExists(bytes32)(bool)" \
  $(cast keccak "502260945") --rpc-url $RPC_URL

cast call $ANCHOR_CONTRACT "getSnapshot(bytes32)${RECORD_TUPLE}" \
  $(cast keccak "502260945") --rpc-url $RPC_URL
```

### Binance `PR01OCT26`

```bash
cast call $ANCHOR_CONTRACT "snapshotExists(bytes32)(bool)" \
  $(cast keccak "PR01OCT26") --rpc-url $RPC_URL

cast call $ANCHOR_CONTRACT "getSnapshot(bytes32)${RECORD_TUPLE}" \
  $(cast keccak "PR01OCT26") --rpc-url $RPC_URL

# Latest Binance anchor (may be this id after the homepage batch)
cast call $ANCHOR_CONTRACT "getLatestSnapshot(bytes32)${RECORD_TUPLE}" \
  $(cast keccak "binance") --rpc-url $RPC_URL
```

### Count anchors per exchange

```bash
for x in binance okx bingx mexc kucoin; do
  echo -n "$x "
  cast call $ANCHOR_CONTRACT "getSnapshotCount(bytes32)(uint256)" \
    $(cast keccak "$x") --rpc-url $RPC_URL
done
```

After the homepage assessment batch documented in deployments, counts include at least: binance `4`, okx `4`, bingx `2`, kucoin `1`, mexc `1`.

### Match roots to the public assessment JSON

1. Read `artifactBundleRoot` and `verificationBundleRoot` from `getSnapshot`.
2. Open `docs/reports/<exchange>/<auditId>-assessment.json`.
3. Compare to `artifacts.artifactBundleRoot` and `artifacts.verificationBundleRoot`.

Example paths for homepage rows:

- [`docs/reports/okx/502260945-assessment.json`](./reports/okx/502260945-assessment.json)
- [`docs/reports/binance/PR01OCT26-assessment.json`](./reports/binance/PR01OCT26-assessment.json)
- [`docs/reports/mexc/202610-assessment.json`](./reports/mexc/202610-assessment.json)
- [`docs/reports/bingx/20261015-assessment.json`](./reports/bingx/20261015-assessment.json)

Matching roots means the on-chain digest agrees with that assessment file. Full PoR replay still needs the artifact archive and `por verify`.

---

## Etherscan

Open the Sepolia proxy’s **Read as Proxy** tab (address in [deployments.md](https://github.com/ardmere/ardmere.github.io/blob/main/docs/deployments.md)), use the verified ABI, and call `snapshotExists` / `getSnapshot` with a `snapshotId` from the table above (as `bytes32`).

---

## Further reading

| Topic | Doc |
| --- | --- |
| Full view ABI, events, pre-v3 logs, ethers.js | [anchor-query-api.md](https://github.com/ardmere/ardmere.github.io/blob/main/docs/anchor-query-api.md) |
| Deploy / upgrade history and homepage anchor txs | [deployments.md](https://github.com/ardmere/ardmere.github.io/blob/main/docs/deployments.md) |
| Trust boundary and verification loop | [verifier-architecture.md](https://github.com/ardmere/ardmere.github.io/blob/main/docs/verifier-architecture.md) |
| Operator CLI | [por-cli.md](https://github.com/ardmere/ardmere.github.io/blob/main/docs/por-cli.md) |

# Binance Proof-of-Reserves — Data Guide

> Quick reference for zkPoR verification services: **data sources, APIs, and fields**. All endpoints are **public, unauthenticated, and require no cookies**—the same endpoints used by the `https://www.binance.com/en/proof-of-reserves` landing page itself.

Last verified: 2026-05-05 (current `auditId` = `PR01APR26`)

---

## 0. TL;DR

Layer 1 data collection for a PoR verification service requires only **3 HTTPS endpoints + 1 static object store**:

| # | Type | Path | Purpose |
|---|---|---|---|
| 1 | `GET` | `/bapi/apex/v1/public/apex/market/query/auditProofSnapshotCondition` | List all historical snapshots (including BTC block height anchors) |
| 2 | `POST` | `/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot` | Fetch the **latest** snapshot: Merkle Root + full asset/liability breakdown by coin |
| 3 | `GET` | `/bapi/apex/v1/public/apex/market/por/getDownloadUrl?auditId=...` | Obtain the wallet-address ZIP download URL for any historical snapshot |
| 4 | static | `https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_YYYYMMDD.zip` | Reserve-side plaintext inventory (HotCold + Deposit CSVs; see §6) |

The key to bypassing the Cloudflare landing page (`HTTP 202`): **do not scrape HTML—call the BAPI endpoints above directly**.

---

## 1. Endpoint Reference

Base URL: `https://www.binance.com`

### 1.1 Historical Snapshot List

```
GET /bapi/apex/v1/public/apex/market/query/auditProofSnapshotCondition
```

No parameters. Response:

```json
{
  "code": "000000",
  "success": true,
  "data": [
    "01/04/26 00:00:00 UTC | BTC Block Height 943129",
    "01/03/26 00:00:00 UTC | BTC Block Height 938780",
    "01/02/26 00:00:00 UTC | BTC Block Height 934541",
    "..."
  ]
}
```

Field semantics:

- Each string = `<MM/DD/YY HH:mm:ss UTC> | BTC Block Height <int>`
- The BTC block height is a **public, independently verifiable time anchor** (any BTC node or mempool.space can look it up), used to lock in an immutable snapshot timestamp.
- The list is in reverse chronological order, dating back to November 2022 (the first PoR release).

### 1.2 Current Snapshot (Including Merkle Root)

```
POST /bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot
Content-Type: application/json
Body: {"time":"","pageIndex":0,"pageSize":0}
```

> ⚠️ **Important finding**: The `time` field is currently **ignored by the server**—any value returns the latest snapshot; a non-empty value that does not match returns an empty object with all `null` fields. **Historical snapshot Merkle Roots cannot be retrieved via this endpoint** (you must consult the announcement for that period or a local cache).

Response (excerpt):

```json
{
  "code": "000000",
  "success": true,
  "data": {
    "snapshotTime": "01/04/26 00:00:00 UTC",
    "merkleRootHash": "1906250e2d94afd82c37d219ef823ff8852cf2321f1be438af07cd18b4f63f48",
    "auditor": "zk-SNARKs",
    "auditorLink": "https://www.binance.com/en/support/faq/815b25f0cb054bdd9d35eccc408fe981",
    "auditId": "PR01APR26",
    "auditDate": "01/04/26",
    "snapshotDataList": [
      {
        "coin": "BTC",
        "ratio": "1.0002990000000000",
        "customerLiability": 618951.27130074,
        "binanceLiability": 619136.126,
        "exchangeBalance": 610597.564,
        "thirdPartyCustody": 8538.562,
        "marginInsurance": null,
        "futureInsurance": null,
        "customerLiabilityUsdt": 42259690305.68295822
      }
    ]
  }
}
```

Field semantics:

| Field | Meaning |
|---|---|
| `snapshotTime` | Snapshot timestamp (UTC) |
| `merkleRootHash` | Poseidon hash root of the current account tree—the **core anchor for zkPoR verification** |
| `auditId` | Format `PR01APR26`; used as input to `getDownloadUrl` |
| `auditor` | `zk-SNARKs` (self-verified) or `Mazars` / third-party firm name |
| `auditorLink` | Link to third-party report or verification documentation |
| `coin.ratio` | `binanceLiability / customerLiability`; ≥ 1 indicates full reserves |
| `customerLiability` | Net user liabilities (after margin/contract collateral offsets) |
| `binanceLiability` | On-chain wallet balances + third-party custody balances |
| `exchangeBalance` | Binance-owned hot/cold wallet balances |
| `thirdPartyCustody` | Third-party custody balances |
| `customerLiabilityUsdt` | USDT-equivalent valuation (display/sorting only; not used in proofs) |

For the current period (`PR01APR26`), `snapshotDataList` contains 14 coins in practice: BTC / ETH / USDT / BNB / SOL / USDC / USD1 / XRP / DOGE / PAXG / LINK / U / SUI / LTC.

### 1.3 Historical Snapshot Download URL

```
GET /bapi/apex/v1/public/apex/market/por/getDownloadUrl?auditId=PR01APR26
```

Response:

```json
{
  "code": "000000",
  "success": true,
  "data": "https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_20260401.zip"
}
```

`auditId` naming convention: `PR<DD><MMM><YY>`, e.g.:

| Snapshot Time | auditId | Wallet ZIP filename |
|---|---|---|
| 01/04/26 | `PR01APR26` | `wallet_address_20260401.zip` |
| 01/03/26 | `PR01MAR26` | `wallet_address_20260301.zip` |
| 01/01/26 | `PR01JAN26` | `wallet_address_20260101.zip` |
| 01/12/25 | `PR01DEC25` | `wallet_address_20251201.zip` |

**Direct URL construction also works** (bypassing this API):

```
https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_<YYYYMMDD>.zip
```

Verified for 2025-12 through 2026-04: all return `HTTP 200`, served directly via CloudFront + S3; each file is ~94 MB (`application/zip`).

---

## 2. What You Can and Cannot Obtain

### ✅ Available via Public API

| Information | Source | Use |
|---|---|---|
| All historical snapshot times + BTC block heights | 1.1 | Time anchor / snapshot selection |
| Current Merkle Root Hash | 1.2 | Global zk proof anchor |
| Current auditor + auditorLink | 1.2 | Third-party attestation traceability |
| Current ~14 major coins: customer liabilities / on-chain balances / reserve ratio | 1.2 | Layer 1 solvency display |
| Wallet-address inventory ZIP for any historical snapshot (~94 MB per period) | 1.3 / direct link | Self-verify on-chain balances, address clustering |

### ❌ Not Available via Public API (Other Sources Required)

| Information | Actual source |
|---|---|
| **Historical** Merkle Root | Binance announcement pages + local archives / third-party mirrors (e.g. DeFiLlama PoR history) |
| Per-period zk-SNARK `proof.csv` / `cex_assets_info.json` / verifying key | Download links in Binance announcements + GitHub Releases ([`binance/zkmerkle-proof-of-solvency`](https://github.com/binance/zkmerkle-proof-of-solvency)) |
| A user's Merkle inclusion proof JSON | User must **log in** → Wallet → Verification → "Download Merkle Tree" |
| Asset details for non-Top-14 coins | Same as above (full list only in `cex_assets_info.json` inside the zip) |

> Rule of thumb: **global/aggregate data** via BAPI; **user-level/raw proof artifacts** via static object storage or the user dashboard.

---

## 3. Practical Code Snippets

### 3.1 Fetch Current Root Hash (Node 18+ / browser)

```ts
const r = await fetch(
  "https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot",
  {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ time: "", pageIndex: 0, pageSize: 0 }),
  },
).then((x) => x.json());

console.log(r.data.merkleRootHash); // -> 1906250e2d94...
console.log(r.data.snapshotTime);   // -> 01/04/26 00:00:00 UTC
console.log(r.data.auditId);        // -> PR01APR26
```

### 3.2 List All Historical Snapshots

```ts
const list = await fetch(
  "https://www.binance.com/bapi/apex/v1/public/apex/market/query/auditProofSnapshotCondition",
).then((x) => x.json());

const parsed = list.data.map((s: string) => {
  const [time, btc] = s.split(" | ");
  return {
    time,
    btcBlockHeight: Number(btc.replace(/\D/g, "")),
  };
});
```

### 3.3 Derive `auditId` and ZIP URL from `snapshotTime`

```ts
const MONTHS = ["JAN","FEB","MAR","APR","MAY","JUN","JUL","AUG","SEP","OCT","NOV","DEC"];

function toAuditId(snapshotTime: string) {
  // "01/04/26 00:00:00 UTC" -> "PR01APR26"
  const [mdy] = snapshotTime.split(" ");
  const [mm, dd, yy] = mdy.split("/");
  return `PR${dd}${MONTHS[Number(mm) - 1]}${yy}`;
}

function toWalletZipUrl(snapshotTime: string) {
  const [mm, dd, yy] = snapshotTime.split(" ")[0].split("/");
  return `https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_20${yy}${mm}${dd}.zip`;
}
```

### 3.4 One-Line curl Reference

```bash
# Current Root
curl -s -X POST -H 'Content-Type: application/json' \
  -d '{"time":"","pageIndex":0,"pageSize":0}' \
  https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot \
  | jq '.data | {snapshotTime, merkleRootHash, auditId}'

# Current wallet ZIP URL
curl -s 'https://www.binance.com/bapi/apex/v1/public/apex/market/por/getDownloadUrl?auditId=PR01APR26' \
  | jq -r '.data'
```

---

## 4. Integration Recommendations for zkPoR Verification Services

1. **Polling frequency**: Poll endpoint 1.1 every 6 hours; trigger 1.2 when `data[0]` shows a new entry. Binance publishes a new snapshot at the start of each month—typical cadence is once per month.
2. **Archival strategy**: On each new snapshot, immediately persist the full JSON from 1.2 and the complete zip from 1.3 locally with content addressing (IPFS / S3 + sha256), because:
   - The Merkle Root is only readable from BAPI while it is the "current" snapshot; once it expires, it is no longer available.
   - Historical zip filenames can be guessed, but Binance has not committed to permanent retention.
3. **Triple cross-validation**:
   - (a) `merkleRootHash` ↔ `account_tree_root` in `proof.csv` from the GitHub Release
   - (b) `snapshotTime` ↔ `BTC Block Height` ↔ block timestamp on mempool.space
   - (c) `binanceLiability` ↔ aggregated on-chain balances of wallet addresses in the zip
4. **External exposure**: Package the (a)(b)(c) comparisons into a one-click reproducible `verify.sh` that outputs `PASS / FAIL`.

---

## 5. Known Gotchas

- The landing page `https://www.binance.com/en/proof-of-reserves` returns `HTTP 202` + empty body on direct curl (CF challenge), but BAPI paths **do not** go through CF anti-bot.
- The `time` field in `userReserveAuditProofSnapshot` is a legacy parameter; the server currently only supports an empty string = latest. Passing any specific timestamp returns `data` with all `null`. **Do not mistake this for a broken endpoint.**
- `getDownloadUrl` must use a query string (`?auditId=...`); a JSON body POST returns `code: 000002 illegal parameter`.
- The wallet ZIP is a **wallet address inventory**, not the zk proof file itself. Complete zk proofs must still be obtained from each period's Binance announcement + GitHub Release.

---

## 6. Reserve-Side Data: Wallet Address ZIP Reference

The `wallet_address_<YYYYMMDD>.zip` obtained via §1.3 is not a zk proof file—it is **Binance's plaintext inventory of all reserve addresses for that period**. This is the only public, independently auditable hard evidence on the PoR "reserve side."

### 6.1 ZIP Structure

Using `PR01APR26` as an example (downloaded and verified 2026-05-05):

| File | Size | Row count (order of magnitude) | Contents |
|---|---:|---:|---|
| `PR01APR26_HotCold.csv` | 110 KB | ~10³ | Binance-owned hot/cold wallets |
| `PR01APR26_Deposit.csv` | 275 MB | ~10⁷ | All user deposit addresses |

> The ZIP itself is ~94 MB (gzip compressed); ~275 MB uncompressed. HotCold covers the high-concentration dimension; Deposit covers the long tail.

### 6.2 CSV Schema (Identical for Both Files)

```csv
coin,network,address,balance,Height,Third party custodian name
USDC,ALGO,QYXDGS2XJJT7QNR6EJ2YHNZFONU6ROFM6BKTBNVT63ZXQ5OC6IYSPNDJ4U,27405800.080000000000000000,59804916,""
ETH,ETH,0x57a0dfd29d8aa63a34acddb8dce2910b7e98a646,0.031817860000000000,24781026,""
```

| Column | Meaning | Notes |
|---|---|---|
| `coin` | Asset symbol (BTC/ETH/USDT/...) | Matches §1.2 `snapshotDataList[].coin` |
| `network` | Chain for this address (BTC/ETH/BSC/TRX/SOL/APT/ALGO/...) | Same `coin` may span multiple chains (especially USDT) |
| `address` | Plaintext on-chain address | Can be pasted directly into a block explorer |
| `balance` | Balance at block `Height` | **Binance-reported value**—requires independent verification |
| `Height` | Block height at which balance was taken | Per-`network` height; enables archive-node lookup |
| `Third party custodian name` | Third-party custodian name (if address is custodied externally) | Reconcile against §1.2 `thirdPartyCustody` |

### 6.3 Reconciliation with BAPI Aggregate Data

```text
Σ_{rows where coin=X}  balance   =   binanceLiability[X]   (from §1.2)
                                  =   exchangeBalance[X] + thirdPartyCustody[X]
```

Where:

- Rows with `Third party custodian name == ""` → count toward `exchangeBalance`
- Rows with `Third party custodian name != ""` → count toward `thirdPartyCustody`

This yields a **fully public, one-click reproducible** internal consistency check.

### 6.4 Independent On-Chain Audit Workflow ("Reserve Side" Hard Evidence)

This is the strongest PoR verification achievable through public channels, with no dependency on any zk toolchain:

```text
for each row (coin, network, address, balance_claim, height) in CSV:
    real_balance := chain_archive_node(network).getBalance(address, height)
    if abs(real_balance - balance_claim) > ε:
        flag(row)  # Binance-reported value does not match on-chain

aggregated_real[coin] := Σ real_balance per coin
aggregated_real[coin] >= customerLiability[coin]  ?  # PASS / FAIL
```

Implementation notes:

1. **Group by network and run concurrently**: Reuse one archive connection per chain. ETH/BSC/TRX dominate Deposit.csv—apply rate limiting.
2. **Historical balance queries**: Must use archive nodes (`eth_getBalance` with `blockNumber`). For EVM chains, self-host Erigon/Reth; for BTC use `electrs` + `getblock`; for Solana SPL use a **slot-indexed API** (`SOLANA_INDEX_API_KEY` → [Solana Index](https://solanaindex.top) `GET /api/v1/solana/token-balance/{address}/{mint}/{slot}`); standard `getTokenAccountsByOwner` **has no historical state**.
3. **Price-agnostic**: This layer validates native coin quantities only—no price feed required.
4. **Sampling vs. full audit**: Full audit of 275 MB / tens of millions of addresses is extremely costly. Prioritize 100% verification of HotCold (~10³ rows), then weighted random sampling of Deposit (covering 99% of value by `balance` head is sufficient).

### 6.5 One-Click Script Skeleton

```bash
# Download + verify + extract
URL=$(curl -s "https://www.binance.com/bapi/apex/v1/public/apex/market/por/getDownloadUrl?auditId=PR01APR26" | jq -r '.data')
curl -L "$URL" -o wallet.zip
shasum -a 256 wallet.zip            # Archive: publish SHA256 + download timestamp signature
unzip wallet.zip -d ./wallet/

# Internal consistency: HotCold + Deposit sums vs BAPI binanceLiability
python3 - <<'PY'
import csv, json, urllib.request
from collections import defaultdict
from decimal import Decimal

agg_self, agg_3p = defaultdict(Decimal), defaultdict(Decimal)
for fn in ['wallet/PR01APR26_HotCold.csv', 'wallet/PR01APR26_Deposit.csv']:
    with open(fn) as f:
        for row in csv.DictReader(f):
            bal = Decimal(row['balance'])
            (agg_3p if row['Third party custodian name'] else agg_self)[row['coin']] += bal

api = json.loads(urllib.request.urlopen(urllib.request.Request(
    "https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot",
    data=b'{"time":"","pageIndex":0,"pageSize":0}',
    headers={"Content-Type":"application/json"}
)).read())['data']['snapshotDataList']

for c in api:
    coin = c['coin']
    print(f"{coin:8s}  csv_self={agg_self[coin]:>20}  api_self={c['exchangeBalance']:>20}  "
          f"csv_3p={agg_3p[coin]:>20}  api_3p={c['thirdPartyCustody']:>20}")
PY
```

After running, CSV aggregates should **align one-to-one** with BAPI-reported values across all 14 major coins (decimal error < 1e-8). Any discrepancy means Binance's two data sources (frontend aggregation vs. wallet export) are already inconsistent—a very strong red flag.

### 6.6 Three-Layer Verification Overview (Updated)

| Layer | Data source | What it proves | Depends on Binance private data? |
|---|---|---|---|
| L1 Internal consistency | §1.2 BAPI + §1.3 CSV | Whether Binance's two self-reported datasets are mutually consistent | No |
| L2 Independent on-chain audit | §1.3 CSV + public chain archives | Reserve addresses hold real balances at snapshot height ≥ user liabilities | No |
| L2b Deposit user spot check | [deposit-sample-manifest](./deposit-spot-check.md) + block explorer | Independently reproduce any row from ardmere's sampling conclusions | No |
| L3 User-level zk inclusion | User-downloaded zip after login + GitHub verifier | Individual user balance is included in the tree + full tree is non-negative | **Yes** (private zip) |

> Public services can deliver L1+L2; L3 must be implemented as **user-local in-browser WASM verification**.

---

## 7. Reverse-Engineering Provenance (Reproducible)

None of the APIs documented here appear in Binance's public documentation—they were discovered by reverse-engineering the landing page JS chunks. Reproduction steps:

1. `curl https://www.binance.info/en/proof-of-reserves -o por.html` (`.info` domain uses SSR, no CF blocking)
2. In the HTML, locate `webpack-runtime.*.js` to obtain the chunk id → hash mapping
3. In the main entry `main.*.js`, the `proof-of-reserves` route maps to chunk id `9a3a`
4. Download `https://bin.bnbstatic.com/static/chunks/page-9a3a.<hash>.js`
5. Grep `/bapi/` inside that chunk to find the three endpoints in §1 and their call patterns

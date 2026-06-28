# Binance Proof-of-Reserves 数据获取指南

> 给 zkPoR 验证服务用的「数据源 + API + 字段」速查手册。所有接口均为 **公开、无需鉴权、无需 cookie**，调用方为 `https://www.binance.com/en/proof-of-reserves` 落地页本身。

最后核对：2026-05-05（当期 auditId = `PR01APR26`）

---

## 0. TL;DR

只需要 3 个 HTTPS 端点 + 1 个静态对象存储，就能完成 PoR 验证服务 Layer 1 的全部数据采集：

| # | 类型 | 路径 | 用途 |
|---|---|---|---|
| 1 | `GET` | `/bapi/apex/v1/public/apex/market/query/auditProofSnapshotCondition` | 列出所有历史快照（含 BTC 区块高度锚点） |
| 2 | `POST` | `/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot` | 拉取**最新**快照：Merkle Root + 全币种资产负债 |
| 3 | `GET` | `/bapi/apex/v1/public/apex/market/por/getDownloadUrl?auditId=...` | 拿任意历史快照的钱包地址 ZIP 下载链接 |
| 4 | static | `https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_YYYYMMDD.zip` | 储备侧明文清单（HotCold + Deposit 两份 CSV，详见 §6） |

绕过 Cloudflare 落地页 (`HTTP 202`) 的关键点：**不要爬 HTML，直接调上面的 BAPI**。

---

## 1. 接口详解

Base URL：`https://www.binance.com`

### 1.1 历史快照清单

```
GET /bapi/apex/v1/public/apex/market/query/auditProofSnapshotCondition
```

无参数。响应：

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

字段含义：

- 每条字符串 = `<MM/DD/YY HH:mm:ss UTC> | BTC Block Height <int>`
- BTC 区块高度是 **公开可独立验证的时间锚**（任何 BTC 节点 / mempool.space 都能回查），用于锁死"快照时刻不可篡改"。
- 列表按时间倒序，覆盖到 2022 年 11 月（首期 PoR）。

### 1.2 当期快照（含 Merkle Root）

```
POST /bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot
Content-Type: application/json
Body: {"time":"","pageIndex":0,"pageSize":0}
```

> ⚠️ **重要发现**：`time` 字段当前**被服务端忽略**——任何取值都返回最新快照；非空且不匹配时返回全 `null` 的空对象。**历史快照的 Merkle Root 不能通过此接口取到**（只能去查当期公告或本地缓存）。

响应（节选）：

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

字段语义：

| 字段 | 含义 |
|---|---|
| `snapshotTime` | 快照时刻 (UTC) |
| `merkleRootHash` | 当期账户树 Poseidon hash root，**zkPoR 验证的核心锚点** |
| `auditId` | 形如 `PR01APR26`，对应 `getDownloadUrl` 的入参 |
| `auditor` | `zk-SNARKs`（自验）或 `Mazars` / 第三方机构名 |
| `auditorLink` | 第三方报告 / 验证说明跳转链接 |
| `coin.ratio` | `binanceLiability / customerLiability`，≥ 1 表示满储备 |
| `customerLiability` | 用户净负债（含杠杆/合约保证金抵扣） |
| `binanceLiability` | 链上钱包余额 + 第三方托管余额 |
| `exchangeBalance` | 自有冷热钱包余额 |
| `thirdPartyCustody` | 第三方托管余额 |
| `customerLiabilityUsdt` | 折算 USDT 估值（仅用于排序展示，不参与证明） |

实测当期 (`PR01APR26`) `snapshotDataList` 共 14 个币种：BTC / ETH / USDT / BNB / SOL / USDC / USD1 / XRP / DOGE / PAXG / LINK / U / SUI / LTC。

### 1.3 历史快照下载链接

```
GET /bapi/apex/v1/public/apex/market/por/getDownloadUrl?auditId=PR01APR26
```

响应：

```json
{
  "code": "000000",
  "success": true,
  "data": "https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_20260401.zip"
}
```

`auditId` 命名规则：`PR<DD><MMM><YY>`，例：

| Snapshot Time | auditId | wallet zip 文件名 |
|---|---|---|
| 01/04/26 | `PR01APR26` | `wallet_address_20260401.zip` |
| 01/03/26 | `PR01MAR26` | `wallet_address_20260301.zip` |
| 01/01/26 | `PR01JAN26` | `wallet_address_20260101.zip` |
| 01/12/25 | `PR01DEC25` | `wallet_address_20251201.zip` |

**直接拼 URL 也可用**（绕过这个 API）：

```
https://public.bnbstatic.com/static/proof-of-reserve/wallet_address_<YYYYMMDD>.zip
```

实测 2025-12 ~ 2026-04 全部 `HTTP 200`，由 CloudFront + S3 直供，单文件 ~94 MB（`application/zip`）。

---

## 2. 能拿到 / 不能拿到 的信息一览

### ✅ 公开 API 能拿到

| 信息 | 来源 | 用途 |
|---|---|---|
| 全部历史快照时刻 + BTC Block Height | 1.1 | 时间锚 / 选盘点 |
| 当期 Merkle Root Hash | 1.2 | zk 全局证明锚 |
| 当期 auditor + auditorLink | 1.2 | 第三方背书追溯 |
| 当期 ~14 大币种 客户负债 / 链上余额 / 储备率 | 1.2 | Layer 1 偿付能力展示 |
| 任意历史快照的钱包地址清单 ZIP（94 MB / 期） | 1.3 / 直链 | 自验链上余额、做地址聚类 |

### ❌ 公开 API 拿不到（需要别的途径）

| 信息 | 实际途径 |
|---|---|
| **历史**期 Merkle Root | Binance 公告页 + 本地存档 / 第三方镜像（如 DeFiLlama PoR 历史库） |
| 各期 zk-SNARK `proof.csv` / `cex_assets_info.json` / `verifying key` | Binance 公告内的下载链接 + GitHub Release（[`binance/zkmerkle-proof-of-solvency`](https://github.com/binance/zkmerkle-proof-of-solvency)） |
| 某用户的 Merkle inclusion proof JSON | 用户**登录后** Wallet → Verification → "Download Merkle Tree" |
| 非 Top14 币种的资产明细 | 同上（zip 内 `cex_assets_info.json` 才有完整列表） |

> 经验法则：**全局/聚合数据**走 BAPI，**用户级/原始证明文件**走静态对象存储或用户面板。

---

## 3. 实战代码片段

### 3.1 取当期 Root Hash（Node 18+ / 浏览器均可）

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

### 3.2 列出全部历史快照

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

### 3.3 由 `snapshotTime` 反推 `auditId` 与 zip 链接

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

### 3.4 一行 curl 速查

```bash
# 当期 Root
curl -s -X POST -H 'Content-Type: application/json' \
  -d '{"time":"","pageIndex":0,"pageSize":0}' \
  https://www.binance.com/bapi/apex/v1/public/apex/market/query/userReserveAuditProofSnapshot \
  | jq '.data | {snapshotTime, merkleRootHash, auditId}'

# 当期钱包 ZIP 地址
curl -s 'https://www.binance.com/bapi/apex/v1/public/apex/market/por/getDownloadUrl?auditId=PR01APR26' \
  | jq -r '.data'
```

---

## 4. 给 zkPoR 验证服务的接入建议

1. **轮询频率**：1.1 接口每 6 小时拉一次；当 `data[0]` 出现新条目时触发 1.2。Binance 每月初出新快照，常态新增频率 = 1 次/月。
2. **存档策略**：每出新快照时立即把 1.2 的完整 JSON、1.3 返回的 zip 全量落本地 + 内容寻址（IPFS / S3 + sha256），因为：
   - Merkle Root 只在"是当期"窗口内能从 BAPI 读，过期就读不到。
   - 历史 zip 文件名虽然能猜到，但 Binance 没有承诺永久保留。
3. **三重交叉验证**：
   - (a) `merkleRootHash` ←→ GitHub Release 提供的 `proof.csv` 中的 `account_tree_root` 比对
   - (b) `snapshotTime` ←→ `BTC Block Height` ←→ mempool.space 区块时间戳比对
   - (c) `binanceLiability` ←→ zip 内钱包地址在公链上的实际余额聚合比对
4. **对外暴露**：把 (a)(b)(c) 三个比对结果做成可一键复现的 `verify.sh`，输出 `PASS / FAIL` 即可。

---

## 5. 已知坑

- 落地页 `https://www.binance.com/en/proof-of-reserves` 直接 curl 会返回 `HTTP 202` + 空 body（CF challenge），但 BAPI 路径**不走** CF 反爬。
- `userReserveAuditProofSnapshot` 的 `time` 字段是历史遗留入参，当前服务端只支持空字符串 = 最新；传任何具体时间会得到 `data` 全 `null`。**别误判为接口坏了**。
- `getDownloadUrl` 必须用 query string（`?auditId=...`），用 JSON body POST 会得到 `code: 000002 illegal parameter`。
- 钱包 ZIP 是**钱包地址清单**，不是 zk proof 文件本身。完整 zk 证明仍须从 Binance 每期公告 + GitHub Release 取。

---

## 6. 储备侧数据：钱包地址 ZIP 详解

§1.3 拿到的 `wallet_address_<YYYYMMDD>.zip` 不是 zk proof 文件，而是 **Binance 当期所有储备地址的明文清单**——这是 PoR "储备侧" 唯一公开、可独立审计的硬证据。

### 6.1 ZIP 结构

以 `PR01APR26` 为例（实测下载于 2026-05-05）：

| 文件 | 大小 | 行数级别 | 内容 |
|---|---:|---:|---|
| `PR01APR26_HotCold.csv` | 110 KB | ~10³ | Binance 自有冷热钱包 |
| `PR01APR26_Deposit.csv` | 275 MB | ~10⁷ | 全部用户充值地址 |

> ZIP 本身 ~94 MB（gzip 压缩），解压后 ~275 MB；HotCold 是热点维度，Deposit 是长尾维度。

### 6.2 CSV Schema（两份相同）

```csv
coin,network,address,balance,Height,Third party custodian name
USDC,ALGO,QYXDGS2XJJT7QNR6EJ2YHNZFONU6ROFM6BKTBNVT63ZXQ5OC6IYSPNDJ4U,27405800.080000000000000000,59804916,""
ETH,ETH,0x57a0dfd29d8aa63a34acddb8dce2910b7e98a646,0.031817860000000000,24781026,""
```

| 列 | 含义 | 备注 |
|---|---|---|
| `coin` | 资产代号（BTC/ETH/USDT/...） | 与 §1.2 `snapshotDataList[].coin` 一致 |
| `network` | 该地址所在链（BTC/ETH/BSC/TRX/SOL/APT/ALGO/...） | 同一 `coin` 可能跨多链（USDT 尤甚） |
| `address` | 链上地址明文 | 直接可塞进区块浏览器 |
| `balance` | 该地址在 `Height` 区块的余额 | **Binance 自报值**，需独立校验 |
| `Height` | 余额所在区块高度 | 各 `network` 独立计高，便于 archive 节点回查 |
| `Third party custodian name` | 第三方托管方名称（若该地址属第三方托管） | 与 §1.2 `thirdPartyCustody` 对账 |

### 6.3 与 BAPI 聚合数据的对账关系

```text
Σ_{rows where coin=X}  balance   =   binanceLiability[X]   （来自 §1.2）
                                  =   exchangeBalance[X] + thirdPartyCustody[X]
```

其中：

- `Third party custodian name == ""` 的行 → 计入 `exchangeBalance`
- `Third party custodian name != ""` 的行 → 计入 `thirdPartyCustody`

这给出一条**完全公开、可一键复现**的内部一致性检查。

### 6.4 独立链上审计流程（"储备侧"硬证据）

这是公开渠道下能做到的最强 PoR 验证，不依赖任何 zk 工具链：

```text
for each row (coin, network, address, balance_claim, height) in CSV:
    real_balance := chain_archive_node(network).getBalance(address, height)
    if abs(real_balance - balance_claim) > ε:
        flag(row)  # Binance 自报与链上不符

aggregated_real[coin] := Σ real_balance per coin
aggregated_real[coin] >= customerLiability[coin]  ?  # PASS / FAIL
```

实施要点：

1. **按 network 分组并发**：同链复用同一个 archive 连接，Deposit.csv 里 ETH/BSC/TRX 占绝大头，做好限流。
2. **历史余额查询**：必须用 archive node（`eth_getBalance` 带 `blockNumber`），EVM 链建议自建 Erigon/Reth；BTC 用 `electrs` + `getblock`；Solana SPL 用 **slot 索引 API**（`SOLANA_INDEX_API_KEY` → [Solana Index](https://solanaindex.top) `GET /api/v1/solana/token-balance/{address}/{mint}/{slot}`）；标准 `getTokenAccountsByOwner` **无历史态**。
3. **价格无关**：这一层只校验"原生币数量"，不需要任何报价源。
4. **采样而非全量**：275 MB / 千万级地址全量审计成本极高，可优先 100% 校验 HotCold（10³ 级），然后对 Deposit 做按余额加权随机采样（按 `balance` 头部覆盖 99% 价值即可）。

### 6.5 一键脚本骨架

```bash
# 下载 + 校验 + 解压
URL=$(curl -s "https://www.binance.com/bapi/apex/v1/public/apex/market/por/getDownloadUrl?auditId=PR01APR26" | jq -r '.data')
curl -L "$URL" -o wallet.zip
shasum -a 256 wallet.zip            # 落档：SHA256 + 下载时间戳签名后公示
unzip wallet.zip -d ./wallet/

# 内部一致性：HotCold + Deposit 求和 vs BAPI binanceLiability
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

跑完应能看到 csv 聚合值与 BAPI 自报值在 14 大币种上 **一一对齐**（小数误差 < 1e-8）。如果出现偏差，说明 Binance 自己的两个数据源（前端聚合 vs 钱包导出）就已经不一致——这是非常强的红旗信号。

### 6.6 三层验证总览（更新）

| Layer | 数据源 | 能证明什么 | 是否依赖 Binance 私有数据 |
|---|---|---|---|
| L1 内部一致性 | §1.2 BAPI + §1.3 CSV | Binance 两套自报值是否互洽 | 否 |
| L2 链上独立审计 | §1.3 CSV + 公链 archive | 储备地址在快照高度的真实余额 ≥ 用户负债 | 否 |
| L2b Deposit 用户抽查 | [deposit-sample-manifest](./deposit-spot-check.md) + 区块浏览器 | 独立复现 ardmere 抽样结论中的任意一行 | 否 |
| L3 用户级 zk inclusion | 用户登录后下载的 zip + GitHub verifier | 单个用户余额已被纳入树 + 全树非负 | **是**（私有 zip） |

> 公开服务能做 L1+L2，L3 必须做成「用户本地浏览器内 WASM 验证」的形态。

---

## 7. 反查溯源（可复现）

本文所有 API 均不在 Binance 公开文档中，是通过反查落地页 JS chunk 拿到的。复现步骤：

1. `curl https://www.binance.info/en/proof-of-reserves -o por.html`（`.info` 域名走 SSR，无 CF 拦截）
2. 在 HTML 中找到 `webpack-runtime.*.js`，得到 chunk id → hash 映射
3. 主入口 `main.*.js` 中 `proof-of-reserves` 路由对应 chunk id `9a3a`
4. 下载 `https://bin.bnbstatic.com/static/chunks/page-9a3a.<hash>.js`
5. 在该 chunk 内 grep `/bapi/`，得到本文 §1 的三个 endpoint 与调用形态

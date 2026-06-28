# Gate.com Proof-of-Reserves — 数据可得性指南

> 配套 [`verifier-architecture.md`](./verifier-architecture.md) 与 [`gateio` adapter](../internal/exchanges/gateio/adapter.go)。

## 与 Binance 的关键差异

| 维度 | Binance | Gate.com |
|---|---|---|
| 汇总数据 | BAPI JSON | Web API `getProofOfReservesInfo` + 币种列表 |
| 钱包地址 CSV | 公开 ZIP（HotCold + Deposit） | **不公开** |
| 全局 zk 证明包 | 仅登录用户下载 | 登录用户从 [我的审计](https://www.gate.com/myaccount/myavailableproof) 下载 `zkmerkle_cex_*.tar.gz` |
| 用户 inclusion | 用户 zip + WASM | `user_config.json` + `./main verify user` |
| 链上余额审计 | ardmere 可复用 wallet CSV | **无公开地址清单 → UNVERIFIABLE** |

Gate 的 PoR 页面：[https://www.gate.com/zh/proof-of-reserves](https://www.gate.com/zh/proof-of-reserves)

开源验证工具：[gateio/proof-of-reserves](https://github.com/gateio/proof-of-reserves)

## 数据平面

### 1. 公开汇总（Dashboard API）

页面展示的内容来自 Gate Web API（浏览器内请求，Akamai 可能拦截 datacenter IP）：

| 端点 | 用途 |
|---|---|
| `GET /api/web/v1/proof-of-reserves/getProofOfReservesInfo` | 最新审计：Merkle root、总储备率、客户净余额、超额储备 |
| `GET /api/web/v1/proof-of-reserves/getProofOfReservesCoinList` | 分币种储备率列表 |
| `GET /api/web/v1/proof-of-reserves/getProofOfReservesList` | 历史批次分页 |

`gateio.Adapter` 将 info + coinList 合并落档为 `summarySnapshot` artifact，保存到：

```text
artifacts/gateio/<auditId>/
  raw/<sha256>.json          # 原始 summary（API 合并 bundle 或 import 原样）
  fetch.json                 # 抓取元数据
  bundles/                   # por-run anchor / verify 写入
    <auditId>.artifact-bundle.json
    <auditId>.verification-bundle.json
    <auditId>.anchor.json
```

### 抓取原始数据（仅落档，不跑验证）

```bash
# 推荐：先试 API，失败则自动导入 fixtures 到 artifacts/gateio/<auditId>/raw/
./scripts/gateio/gate-save-local.sh

# 直拉 Gate 公开 API（Akamai 可能拦截 datacenter IP）
go run ./cmd/por fetch gateio
# 或
./scripts/gateio/gate-fetch.sh

# 浏览器 DevTools 复制 API 响应
go run ./cmd/por fetch gateio -info-file ./info.json -coins-file ./coinList.json
./scripts/gateio/gate-import-browser.sh ./info.json ./coinList.json
```

### 抓取 + 验证 + 生成 bundle

```bash
# API 或 -summary-path 均可；原始数据自动写入 artifacts/gateio/<auditId>/raw/
go run ./cmd/por anchor -exchange gateio -skip-zip

# 手动 import
go run ./cmd/por anchor -exchange gateio -summary-path ./summary.json -skip-zip

# 从已落档 snapshot 重跑验证
go run ./cmd/por verify -exchange gateio -snapshot <auditId> \
  -artifacts ./artifacts/gateio/<auditId>
```

**本地绕过 API 封锁（旧写法仍可用，现已默认写入 artifacts/gateio/）：**

### 2. zk 全局证明包（需登录）

用户在 [我的审计](https://www.gate.com/myaccount/myavailableproof) 下载：

- **Download Merkle Tree** → `zkmerkle_cex_xxx.tar.gz`
- **Download User Config** → `user_config.json`（放入 `config/`）

解压后结构（来自 [Gate 官方文档](https://github.com/gateio/proof-of-reserves)）：

```text
config/
  cex_config.json    # CexAssetsInfo + proof.csv 路径 + vk 前缀
  user_config.json   # 用户 inclusion（可选）
proof.csv
zkpor864.vk.save
main                 # 验证二进制（GitHub Releases）
```

**交易所资产验证（全局 zk）：**

```bash
./main verify cex
# 成功输出: All proofs verify passed!!!
```

**用户 inclusion 验证：**

```bash
./main verify user
# 成功输出: verify pass!!!
```

将 tar.gz 手动导入 artifact bundle（会落档到 `artifacts/gateio/<auditId>/raw/`）：

```bash
go run ./cmd/por anchor -exchange gateio \
  -summary-path ./summary.json \
  -zk-bundle ./zkmerkle_cex_xxx.tar.gz \
  -skip-zip
```

> ardmere 当前将 zk bundle 落档为 `globalProofBundle` artifact；`global-zk-proof@gateio-1` verifier 尚未实现（stub `@gateio-0`）。

### 3. 无公开数据 → UNVERIFIABLE 的维度

| Verifier | Gate 状态 | 原因 |
|---|---|---|
| `internal-consistency` | UNVERIFIABLE | 无 wallet CSV 与汇总逐币对账 |
| `onchain-balance-*` | UNVERIFIABLE | 无公开 HotCold 地址列表 |
| `btc-anchor` | UNVERIFIABLE | Gate 摘要不绑 BTC block height |
| `global-zk-proof@gateio` | UNVERIFIABLE（默认） | 需登录下载 tar.gz；公开 API 不提供 |
| `solvency-claim` | 部分可验证 | 公开 `total_reserve_rate ≥ 100`（自报） |

## ardmere 验证路径（当前）

```text
公开 API / 手动 summary.json
        ↓
gateio.Adapter → por.Snapshot
        ↓
solvency-claim@1（总储备率 ≥ 100%，标注 self-reported）
        ↓
stub verifiers（诚实标注 UNVERIFIABLE + 原因）
```

未来激活 `global-zk-proof@gateio-1`：解析 `cex_config.json` + `proof.csv` + vk，复用 [gateio/proof-of-reserves](https://github.com/gateio/proof-of-reserves) 的 Go verifier 或 exec `./main verify cex`。

## 参考

- [Gate PoR 官网](https://www.gate.com/zh/proof-of-reserves)
- [Gate Learn — 如何验证](https://www.gate.com/learn/articles/how-to-use-gate-io-proof-of-reserves-to-verify-your-assets-security/1017)
- [gateio/proof-of-reserves README](https://github.com/gateio/proof-of-reserves)

# 下一波接入短名单（Top 10）

> **基准日**：2026-11-06  
> **已覆盖（勿再推荐）**：`binance` · `okx` · `gateio` · `bybit` · `bitget` · `htx` · `mexc` · `kucoin` · `bingx`（见 [`config/exchanges/registry.yaml`](../config/exchanges/registry.yaml)）  
> **旧文**：Project store 里仍把 MEXC/KuCoin/BingX 写成「下一所」的备忘录已过时，以本文为准。  
> **范围**：只读建议 + 登记；**未实现**完整适配器（除非后续任务明确要求）。

产品目标：独立 CEX PoR 透明度登记（[Stage 框架](./por-transparency-framework.md)）。优先选 **可公开核验的 artifact**（Merkle / 钱包列表 / ownership / ZK / 审计包）且能进现有 `fetch → verify → report` 管线的所。

## 评分维度

| 维 | 含义 |
| --- | --- |
| 公开 artifact | 免登录 summary / root / wallet list / ownership / OSS verifier 是否存在 |
| 管线契合 | 是否接近现有 Gen1 Merkle（Bitget/Bybit/MEXC）或 ownership+ZK（OKX）模式 |
| 活跃度 | 2026 仍在更新；陈旧附证降权 |
| 预期 Stage 天花板 | Pre-Stage / Stage 0 / Stage 1（Stage 2 需高频锚，首版不承诺） |
| 阻塞 | 登录墙、WAF、仅 scrape、截断 hash、更新稀少、交易所停运 |

## 推荐排序（registry 之外）

### 1. Poloniex — 首选下一适配器

| | |
| --- | --- |
| **Why** | 公开 PoR 页月更比率 + Merkle root（探针见 Snapshot `20261001`）；GitHub [`tools-nodejs-address-verify`](https://github.com/poloniex/tools-nodejs-address-verify) 提供 **地址 ownership**；另有 MerkleVerify 用户纳入工具——在未覆盖所里最接近「summary + ownership + user Merkle」组合。 |
| **预期 Stage 天花板** | **Stage 0** 首版；若钱包 CSV + ownership 可 permissionless 拉取，有 **Stage 1** 冲刺可能（仍缺全局 ZK / 高频则到不了 Stage 2）。 |
| **依赖 / 阻塞** | 用户 leaf 下载偏登录；需确认公开钱包包 URL 与签名消息格式后才能写 ownership verifier。 |

### 2. Ourbit

| | |
| --- | --- |
| **Why** | [`/assets/reserve-proof`](https://www.ourbit.com/assets/reserve-proof) 宣称公开多链钱包地址/余额/块高；OSS [`OurbitOfficial/merkle_validate`](https://github.com/OurbitOfficial/merkle_validate)（Java，与 MEXC/Bitget 同构）。 |
| **预期 Stage 天花板** | **Stage 0**；钱包列表若可机器拉取且可对账 → 向 Stage 1 证据靠拢（ownership / 全局证明仍可能缺）。 |
| **依赖 / 阻塞** | 用户 path JSON（`merkel_tree_ourbit.json`）登录门控；hash 截断至 32 hex（128-bit）需在报告里标风险。 |

### 3. CoinEx

| | |
| --- | --- |
| **Why** | 公开 [reserve-proof](https://www.coinex.com/en/reserve-proof) 页；OSS Go [`coinexcom/proof-of-reserves`](https://github.com/coinexcom/proof-of-reserves) ——与 ardmere Go 栈最省摩擦。 |
| **预期 Stage 天花板** | **Stage 0**（用户证明 + 储备对比）；缺公开 ownership / 全局 ZK 则无 Stage 1。 |
| **依赖 / 阻塞** | 仓库活跃度低、issue 较多；需探针确认当前期 merkle 文件是否仍无登录可下。 |

### 4. Bitunix

| | |
| --- | --- |
| **Why** | 2026 仍更新（公开材料含 Apr/Jul 比率）；[`BitunixOfficial/proof-of-reserves`](https://github.com/BitunixOfficial/proof-of-reserves) Java verifier + `merkle_tree_bitunix.json` 下载路径，适配器可抄 MEXC。 |
| **预期 Stage 天花板** | **Stage 0**。 |
| **依赖 / 阻塞** | 叶子/整树下载偏登录；云端是否可抓最新 root 需 WAF 探针。 |

### 5. WOO X

| | |
| --- | --- |
| **Why** | [Live PoR 看板](https://woox.io/proof-of-reserves)（自称 ~15 分钟刷新）+ 可下载完整 Merkle 树；OSS [`woo-x/merkle-tree`](https://github.com/woo-x/merkle-tree)（Protobuf 树 + Java `MerkleCheck`）。对「高频披露」叙事有对照价值。 |
| **预期 Stage 天花板** | **Stage 0**（负债 Merkle）；资产侧大量依赖第三方流动性源 API——链上钱包全集未必 permissionless，**Stage 1 难**。 |
| **依赖 / 阻塞** | SPA/API 抓取；用户 hash 登录；资产侧非纯自托管钱包列表。 |

### 6. Toobit

| | |
| --- | --- |
| **Why** | 公开 PoR 页 + 较新的 SMT 离线验证器 [`toobit-docs/por-verifier`](https://github.com/toobit-docs/por-verifier)（SHA-256 summation Merkle，2026 仓）。 |
| **预期 Stage 天花板** | **Stage 0**。 |
| **依赖 / 阻塞** | 份额小于头部；需确认公开 root/审计包是否免登录；用户 proof 文件多半登录。 |

### 7. Phemex

| | |
| --- | --- |
| **Why** | 长期月更 Merkle（CEX Compare 录到 2026-11）；PoR UI 展示钱包列表与负债对比，适合 Stage 0 登记行。 |
| **预期 Stage 天花板** | **Stage 0**（钱包列表若可导出可增强证据，仍缺 ownership/ZK）。 |
| **依赖 / 阻塞** | 用户自证登录；官方页常被营销壳包裹，summary 需稳定 API/浏览器抓取。 |

### 8. BTCC

| | |
| --- | --- |
| **Why** | 节奏极稳：公开材料显示 **每月 15 日**快照（自 2025-04 起连续多期）；比率 + Merkle root 公开，便于写进 cadence 窗。 |
| **预期 Stage 天花板** | **Stage 0 硬顶**：无钱包金额/地址、root 仅约 16 hex 截断、叶子登录——管线可做 summary 行，冲高 Stage 不划算。 |
| **依赖 / 阻塞** | 截断 hash；无 permissionless 资产侧。 |

### 9. Kraken

| | |
| --- | --- |
| **Why** | 自 2014 起的会计师主导 Merkle；机构/监管叙事对照面（非成交驱动）。最新公开窗约 **季度**（如 2026-06-30）。 |
| **预期 Stage 天花板** | **Stage 0**「attestation-model」轻量行；无交易所侧 permissionless 全量钱包 ZIP / 全局 ZK → **无 Stage 1**。 |
| **依赖 / 阻塞** | 季度低频；依赖会计师工具链与用户个性化证明；不适合作高频 monitor 主路径。 |

### 10. CoinW

| | |
| --- | --- |
| **Why** | 公开宣称 Merkle PoR 页；补齐二线所覆盖表。 |
| **预期 Stage 天花板** | **Stage 0**（若只能看到营销页/无稳定 artifact → 实为 Pre-Stage 行）。 |
| **依赖 / 阻塞** | **Cloudflare / WAF**（本环境 2026-11-06 探针 HTTP 403）；缺确认 OSS verifier；落地前必须浏览器探针。 |

## 明确后置（不要进本轮 Top 10）

| 所 | 原因 |
| --- | --- |
| **Crypto.com** | Mazars 等公开附证停在 ~2022-12（陈旧）。 |
| **WhiteBIT** | Hacken 附证停在 ~2024-11。 |
| **BitMEX** | 2026-09-23 起停止交易/风停页；历史 OSS 不再作为活跃接入目标。 |
| **Deribit** | 2026-09-01 起停止日更 PoR 页（并入 Coinbase 托管叙事）。 |
| **Bitfinex / Gemini / Bitstamp / Upbit / Coinbase** | 无可用零售 Merkle 管线或仅监管披露 → Pre-Stage 对照，非本轮适配器。 |
| **XT.com** | 第三方整理显示仍有 2026-10 快照，但官方 PoR URL 本环境 404；复核后再升榜。 |

## 给协调员的可粘贴短榜

```
已覆盖：Binance OKX Gate Bybit Bitget HTX MEXC KuCoin BingX
下一波 Top 10（2026-11-06）：
1. Poloniex — ownership OSS + 月更 Merkle；Stage 0→1 潜力最高
2. Ourbit — 公开多链钱包页 + Java merkle_validate
3. CoinEx — Go OSS verifier，栈契合
4. Bitunix — 2026 仍更；Java verifier 同构 MEXC
5. WOO X — 实时看板 + 全树下载 + woo-x/merkle-tree
6. Toobit — 新 SMT 离线 verifier
7. Phemex — 月更 + UI 钱包列表
8. BTCC — 每月 15 日极稳；Stage 0 硬顶（截断 root/无钱包）
9. Kraken — 季度会计师对照；轻量行
10. CoinW — Merkle 宣称；WAF 403 阻塞需浏览器
后置：Crypto.com/WhiteBIT 陈旧；BitMEX 停运；Deribit 停日更；XT URL 待复核
```

## 来源

- 已覆盖集合：`config/exchanges/registry.yaml` + `docs/reports/snapshot-history.md`
- 景观对照：[CEX Compare PoR landscape](https://cex-compare.com/en/insights/proof-of-reserves-landscape/)（checked ~2026-11-05）
- 各所官方 PoR 页 / GitHub verifier（2026-11-06 探针）
- 落地前对目标所再跑公开面探针；数字与入口会变

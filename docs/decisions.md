# ardmere zkPoR Verifier — 关键决策记录

> 此文件记录架构层面的不可逆决策（ADR — Architecture Decision Records）。
> 配套：[`verifier-architecture.md`](./verifier-architecture.md)

| 字段 | 值 |
|---|---|
| 文档版本 | v1.0 |
| 决策日期 | 2026-05-05 |

---

## ADR-001 后端语言：**Go**

**选项**：Go / TypeScript

**决策**：Go

**理由**：
- 直接复用 [`binance/zkmerkle-proof-of-solvency`](https://github.com/binance/zkmerkle-proof-of-solvency) 的 verifier 包，未来 `global-zk-proof@1` 激活几乎零成本
- 大文件流式处理（94 MB walletZip / 275 MB CSV）原生友好
- 多链 RPC 并发调用 goroutine 模型简单
- 静态二进制，部署极简

**影响**：
- 前端 React island 与后端通过 REST + JSON 解耦
- 浏览器 WASM 用户验证（决策 4）需要把 Go 编译成 WASM（`GOOS=js GOARCH=wasm`，二进制约 5~10 MB），可接受

---

## ADR-002 主锚定链：**Base**

**选项**：Base / Arbitrum / Optimism / Ethereum L1

**决策**：Base

**理由**：
- Gas 极低（~$0.001 per anchor tx），可承受高频心跳
- EVM 兼容，Foundry / viem 工具链成熟
- 数据可用性继承自 Ethereum L1（OP Stack 通过 batch posting）
- 公共 RPC 多家提供，符合 ADR-005 的零自营原则
- 由 Coinbase 背书，长期可用性预期高

**影响**：
- 部署一份极简事件型合约 `ArdmerePoRAnchor`
- 服务持有一个 secp256k1 私钥（由 KMS 管理）作为 anchor signer
- 链 ID 写死到配置：`8453`（Base mainnet）

**未来可选**：
- 年度审计性快照可考虑额外锚到 Bitcoin OP_RETURN（最高强度但贵），不在 MVP 范围

---

## ADR-003 冷存储策略：**不引入 Arweave / IPFS，仅关键摘要上链**

**选项**：
- A) Arweave 永久存原始 verification bundle JSON
- B) IPFS pinning（多 pinner）
- C) **无独立冷存储，仅链上写摘要**
- D) S3 + 自维护备份

**决策**：**C — 无独立冷存储，仅关键摘要上链**

**理由**：
- 简化架构，降低运维成本和外部依赖
- 链上锚定的 Merkle Root 已经提供了"事后不可篡改"的核心保证
- 原始数据保留在我方 Postgres + S3（运营级别冗余足够）
- 真正不可恢复的数据是"Binance 当期发布的 walletZip"——这本来就不在我们控制范围，只能保留 sha256 在链上做事后比对
- Arweave 一次性付费虽然便宜，但引入新的依赖与运维面，对 MVP 不必要

**影响**：

| 数据 | 存哪 | 不可篡改保证 |
|---|---|---|
| Merkle Root（artifact bundle / verification bundle） | **Base 链上 event** | ✅ 链上不可逆 |
| 原始 Artifact（bapiSnapshot JSON / walletZip） | **S3 + Postgres** | ⚠ 我方运营保证 + 链上 sha256 比对 |
| 原始 Verification（findings JSON） | **S3 + Postgres** | ⚠ 我方运营保证 + 链上 sha256 比对 |
| 我方 ed25519 公钥 | **Base 链上 + 域名 DNSSEC TXT** | ✅ 双源可验 |

**修订的 §6 锚定对象（最终版，schema v2）**：

```solidity
event SnapshotAnchored(
    bytes32 indexed snapshotId,
    bytes32 indexed exchangeTag,       // keccak256("binance")
    string exchange,
    bytes32 artifactBundleRoot,        // 原始 artifact 摘要 Merkle Root
    bytes32 verificationBundleRoot,    // verifier 结论 Merkle Root
    uint8 verdictSummary,
    uint16 coverageBps,
    // … periodSeq, snapshotTime, btcBlockHeight, exchangeMerkleRoot, schemaVersion, anchoredAt
);
```

> 每期快照 **1 笔 tx** 同时锚定 artifact 与 verification 两个 root。原始数据获取统一走 ardmere 自家 API（`GET /artifacts/:sha256`），用户拿到后用链上 root 自验。

**取舍承认**：
- 失去了"ardmere 跑路后历史原始数据仍然可访问"的强保证
- 换得了：架构简化、运维收敛、上线提速
- 缓解：MVP 上线后若需求强烈，可在 V2 加一个独立的 Arweave 镜像服务（仅锚定 root 不变）

---

## ADR-004 用户验证形态：**全本地浏览器 WASM**

**选项**：
- A) 用户上传 zip 到我方后端验证
- B) **全本地浏览器内 WASM 验证（零上传）**
- C) 提供桌面 CLI 让用户自己跑

**决策**：B

**理由**：
- 用户自己下载的 Merkle Tree zip 包含敏感账户身份信息（账户索引、余额结构），**任何上传都是隐私泄露风险**
- 浏览器内 WASM 全本地：
  - 用户体验接近一键
  - ardmere 不接触任何用户数据，零责任、零合规暴露
  - 可在 ardmere.org 直接嵌入，与主站视觉一体
- Go 编 WASM 体积可接受（5~10 MB，首次加载 + 缓存）

**实现路径**：
- 把 `binance/zkmerkle-proof-of-solvency` 的 user verifier 部分编 WASM
- 前端 `<input type="file">` 接受 zip → File API 读取 → 传 WASM 内存 → 输出 PASS/FAIL
- ardmere 后端**完全不参与**这条路径，连日志都不要记

**影响**：
- 不需要任何用户数据相关的 API、数据库表、合规流程
- 前端要做一个独立的 React island（与主站集成但状态隔离）

---

## ADR-005 多交易所路线图：**OKX 第二个**

**选项**：OKX / Coinbase / Bitget / Kraken / Bybit / 其他

**决策**：MVP 之后，**优先支持 OKX**

**理由**：
- OKX 也开源了 PoR 工具（[`okx/proof-of-reserves-v2`](https://github.com/okx/proof-of-reserves-v2)），数据规范度类似 Binance
- OKX 是体量第二大的交易所，覆盖它能立即扩大用户基数
- 同样使用 zk-SNARK + Merkle Tree 方案，大部分 verifier 可复用
- 与 Binance 形成横向对比的故事性强（"同一套独立标准下两家的表现"）

**影响**：
- 领域模型 [`Snapshot`](./verifier-architecture.md#3-领域模型) 必须从 Day 1 就带 `exchange` 字段，不能写死 Binance
- 配置层抽象出 `ExchangeAdapter`，每个交易所一个：
  ```go
  type ExchangeAdapter interface {
      Id() string                            // "binance" | "okx"
      Fetchers() []Fetcher
      DefaultVerifiers() []Verifier
      ParseSnapshotId(s string) (SnapshotId, error)
  }
  ```
- MVP 只注册 BinanceAdapter，但接口预留 OKX 槽位

---

## 决策时间线

```
2026-05-05  ADR-001..005 拍板（v1.0）
   后续决策（如出现）追加 ADR-006...，旧 ADR 不修改，仅 supersede
```

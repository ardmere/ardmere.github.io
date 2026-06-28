# ardmere zkPoR Verifier — 架构设计文档

> 一个面向 Binance Proof-of-Reserves（及未来其他交易所）的**独立、可重放、可扩展**的验证服务。
> 同时把"今天验证不了的维度"也作为一等公民进行建模，以便外部条件成熟时无缝激活。

| 字段 | 值 |
|---|---|
| 文档版本 | v0.1 |
| 最近更新 | 2026-06-07 |
| 配套文档 | [`binance-por-data-guide.md`](./binance-por-data-guide.md) |
| 状态 | Draft，可作为 MVP 实施依据 |

---

## 1. 设计哲学

三条核心原则，决定后面所有结构：

1. **Evidence-first，不是 Boolean-first**。每一次验证不只输出 `PASS / FAIL`，而是输出一份**可重放的证据包**：输入数据哈希 + 验证逻辑版本 + 验证结果 + 时间戳 + 我方签名 + **关键摘要上链**。即使将来 Binance 改了 API、删了文件、换了算法，历史结论依然可追溯、可审计。
2. **Verifier 是插件，不是 if/else**。每个验证维度（内部一致性、链上余额、zk 证明、地址签名……）都是一个独立的 `Verifier`，通过统一接口接入。**新增维度 = 新增一个文件 + 注册**，不动核心。
3. **诚实优于完整**。任何当前无法验证的维度，**显式输出 `UNVERIFIABLE` 并说明原因**，绝不退化成 `PASS` 或 `PENDING`。这是 PoR 服务最重要的信用资产。

---

## 2. 当前数据可得性边界（决定可验证维度）

详见 [`binance-por-data-guide.md`](./binance-por-data-guide.md)，此处摘要：

| 数据 | 公开渠道 | 影响哪类 verifier |
|---|---|---|
| 当期 Merkle Root + 资产汇总 | ✅ BAPI | `internal-consistency`, `solvency-claim` |
| 历史快照清单 + BTC 区块高度 | ✅ BAPI | `btc-anchor` |
| 钱包地址 + 余额（明文 CSV） | ✅ 公开 ZIP | `internal-consistency`, `onchain-balance-*` |
| 全局 zk-SNARK proof.csv / vk | ❌ 仅登录用户 | `global-zk-proof`（占位）|
| 地址签名 / 归属证明 | ❌ 无官方下载渠道 | `address-ownership`（占位）|
| 第三方审计 attestation | ❌ 无官方下载渠道 | `third-party-attestation`（占位）|

---

## 3. 领域模型

```text
Snapshot                     ── Binance 的某一期 PoR 快照
├── id            : "PR01APR26"
├── snapshotTime  : 2026-04-01T00:00:00Z
├── btcAnchor     : { height: 943129, hash, timestamp }
├── claims        : Claim[]                    ← Binance 自报的"事实声明"
└── artifactRefs  : sha256[]                   ← 关联 Artifact 内容寻址

Claim                        ── 一条可验证的断言
├── kind          : "merkleRoot" | "coinReserve" | "addressBalance" | ...
├── subject       : 描述对象（coin / address / batch）
├── value         : 自报值
└── source        : (artifactSha256, locator)

Artifact                     ── 原始数据快照（内容寻址，不可变）
├── kind          : "bapiSnapshot" | "walletZip" | "globalProof" | "addressSig"
├── sha256        : 内容哈希（主键）
├── source        : { url, fetcherId, fetcherVersion }
├── fetchedAt     : timestamp
├── ourSignature  : ed25519(sha256 ‖ fetchedAt ‖ url)
├── onchainTx     : 上链锚定 tx hash（见 §6）
└── storage       : { s3: key }                ← 仅 S3（决策 ADR-003）

Verification                 ── 一次验证执行
├── id            : ULID
├── snapshotId    : "PR01APR26"
├── verifierId    : "internal-consistency@1.2"
├── verifiedAt    : timestamp
├── inputHashes   : sha256[]                   ← 验证时读了哪些 Artifact
├── verdict       : PASS | FAIL | UNVERIFIABLE | PARTIAL
├── findings      : Finding[]                  ← 逐条细节
├── coverage      : 0.0 ~ 1.0                  ← 抽样比例（链上审计必填）
├── ourSignature  : ed25519(整份 Verification 的 canonical JSON)
└── onchainTx?    : 关键 verification 的上链锚定（仅对 root-level 结论）
```

> **`coverage` 字段**：L2 链上审计永远是采样的，必须显式声明"我覆盖了多少价值 / 多少地址"。前端必须把 coverage 显示出来，不可隐藏。

---

## 4. Verifier 插件接口

```ts
interface Verifier<C extends Claim = Claim> {
  readonly id: string;              // "internal-consistency@1.2"
  readonly version: string;         // semver；逻辑变更必须升版
  readonly claimKinds: Claim["kind"][];
  readonly requires: Artifact["kind"][];

  // 核心方法：纯函数，输入 = ctx 提供的所有 artifact，输出 = 验证结论
  verify(ctx: VerifyCtx): Promise<Verification>;

  // 自描述（用于前端动态渲染"为什么这一项是 UNVERIFIABLE"）
  capability(): {
    canVerify: Claim["kind"][];
    cannotVerify: { kind: Claim["kind"]; reason: string }[];
  };
}

interface VerifyCtx {
  snapshot: Snapshot;
  artifacts: Map<sha256, Artifact>;
  // 受控外部依赖（注入而非自取，便于 mock 与重放）
  rpc: ChainRpcPool;             // §7 公共节点池
  clock: () => Date;
}
```

**关键设计**：

- `verify()` 是**纯函数**，输入完全由 ctx 提供，禁止直接 `fetch()` / `fs.readFile`——这样**可重放、可单测、可链上 zk 化**。
- 抓数据是**另一类组件 `Fetcher`** 的事，verifier 只消费已落档的 artifact。
- `version` 升一位 → 旧 verification 自动标记 `STALE`，触发重跑（旧记录保留，方便对比逻辑变更带来的差异）。
- `rpc` 池由外部注入（§7），verifier 不感知节点提供商。

---

## 5. Verifier 矩阵

### 5.1 Day 1 上线（基于公开数据 100% 可做）

| Verifier | 输入 artifact | 验证逻辑 | 输出 | coverage |
|---|---|---|---|---|
| `artifact-integrity@1` | 所有 | 校验落档时记录的 sha256 与重算一致；我方签名有效；onchainTx 已确认 | PASS / FAIL | n/a |
| `internal-consistency@1` | bapiSnapshot + walletZip | 对 14 大币种，CSV 聚合 = BAPI `binanceLiability/exchangeBalance/thirdPartyCustody` | PASS / FAIL | 100% |
| `btc-anchor@1` | bapiSnapshotList + 公链查询 | 快照声明的 BTC Block Height 区块时间戳 ≈ snapshotTime（误差 < 30 min）| PASS / FAIL | 100% |
| `solvency-claim@1` | bapiSnapshot | 14 大币种 `binanceLiability ≥ customerLiability`（注：仅校验 Binance 自报数字之间的关系，**不证明储备真实性**——这条单独标注 "self-reported"）| PASS / WARN | 100% |
| `onchain-balance-hot@1` | walletZip + 公共 RPC | HotCold 行：在 CSV `Height` 区块比对**可观测余额**与 CSV `balance`（见 §5.3） | PASS / FAIL / WARN | 按 `(coin,network)` 对 |
| `onchain-balance-deposit@1` | walletZip + 公共 RPC | Deposit 按价值加权采样审计（覆盖 99% 价值约 10⁵ 地址） | PARTIAL | 0~99% |

### 5.2 占位 stub（今天发布，未来直接激活）

| Verifier | 占位输出 | 激活条件 |
|---|---|---|
| `address-ownership@0` | `UNVERIFIABLE` reason: *"No public download channel for wallet ownership signatures / proofs"* | Binance 提供官方下载渠道 |
| `global-zk-proof@0` | `UNVERIFIABLE` reason: *"Global proof.csv / verifying key not publicly distributed; only available via logged-in user download"* | Binance 公开发布 / 第三方镜像 |
| `third-party-attestation@0` | `UNVERIFIABLE` reason: *"No public third-party attestation report available"* | 任意机构发布可公开下载的 PoR attestation |
| `cross-chain-wrapped@0` | `UNVERIFIABLE` reason: *"Wrapped tokens (wBTC/cbBTC/...) reconciliation rules not finalized"* | 我方完成 wrapped 资产对账规范 |

> **占位 verifier 的价值**：让前端 UI **今天**就显示出 10 个验证维度——其中 6 个 ✅/❓ 真实验证 + 4 个 ❌ 显式标注"无法验证 + 原因"，比"只展示 6 个真实维度"更诚实，也更有故事。

### 5.3 `onchain-balance-hot`：staking-aware 演进

#### 5.3.1 问题：v1 只查 liquid native balance

当前实现（`onchain-balance-hot@1.0`）对 HotCold 每一行调用 `eth_getBalance(address, Height)`，与 CSV 的 `balance` 字段直接比对。

Binance walletZip 的 `balance` 是**账面持仓**——同一地址上可能包含：

| 形态 | 是否在 EOA `balance` 里 | v1 能否看到 |
|---|---|---|
| Liquid native（ETH/BNB） | ✅ | ✅ |
| BNB Stake Hub 质押（`getPooledBNB`） | ❌ 在 credit 合约记账 | ❌ |
| BNB 解绑队列（`lockedBNBs`） | ❌ | ❌ |
| Ethereum 2.0 存款（`DepositContract`） | ❌ 已进信标链 | ❌ |
| ERC20 / BEP20 | ❌ 在 token 合约 | ❌（另开 `onchain-balance-erc20`） |

**PR01JUN26 实证**：3 个 native-only FAIL 中，2 个为质押/存款口径差（`0xbf83…` BNB 质押、`0x32e11…` 80k ETH 已 deposit），1 个（`0x86523…`）约 **75k / 82k BNB 差额**可归因于 Stake Hub delegation（见 §5.3.6）。v1 把这类行标为 FAIL 会**高估真实偏差**。

设计原则：**不把"我们还没实现的查询"当成 Binance 造假**——要么扩展可观测余额，要么显式 `WARN` 并说明未覆盖形态。

#### 5.3.2 目标余额模型：`totalAccounted`

对每一 HotCold 行 `(coin, network, address, Height)`，定义：

```text
totalAccounted = liquidNative
               + stakedNative      // 链上可验证的质押账面
               + unbondingNative   // 解绑中、尚未回到 EOA 的部分
               + (future) erc20Balances[coin]
```

与 CSV 比对：

```text
diff = totalAccounted - csvBalance
|diff| ≤ tolerance  →  PASS（或 WARN 若 diff 非零但在容差内）
|diff| > tolerance  →  FAIL（Finding 须拆分各分量，便于第三方复现）
```

容差沿用 v1：`abs ≤ 1e-4` native unit **或** `rel ≤ 1e-7`（见 `docs/por-cli.md`）。

**Finding 扩展字段**（`onchain-balance-hot@2.0` 起）：

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

若某形态尚未实现（如某链 staking），`components.unsupported` 留空并在 `note` 说明 →  verdict 降为 **`WARN`**（"incomplete observation"），**不得**标 FAIL。

#### 5.3.3 BSC — BNB Stake Hub（`onchain-balance-hot@2.0`）

| 项 | 值 |
|---|---|
| StakeHub 合约 | `0x0000000000000000000000000000000000002002` |
| 快照块高 | 取 CSV 行内 `Height`（与 v1 相同） |
| Validator 列表 | `getValidators(offset, limit)` → `(operator[], credit[])` |
| 每 validator 查询 | 对 `credit` 调用 `getPooledBNB(delegator)`、`lockedBNBs(delegator, 0)` |
| 聚合 | `stakedNative = Σ getPooledBNB`；`unbondingNative = Σ lockedBNBs` |

```text
┌─ HotCold row: BNB|BSC, address A, height H ─────────────────────┐
│  liquidNative  ← eth_getBalance(A, H)                           │
│  stakedNative  ← Σ_credit getPooledBNB(A) @ H                   │
│  unbonding     ← Σ_credit lockedBNBs(A, 0) @ H                  │
│  totalAccounted = liquid + staked + unbonding                   │
│  compare ↔ CSV balance                                          │
└─────────────────────────────────────────────────────────────────┘
```

**性能优化**（避免 53 validators × 2 calls × 10³ 行）：

1. **Validator 列表缓存**：同一 `Height` 只拉一次 `getValidators`，存 Postgres `(height → credit[])`。
2. **稀疏扫描**：对某 `delegator`，可先读 StakeHub `Delegated` / `Redelegated` logs（`topics[2] = delegator`）得到**非零 operator 集合**，只查相关 credit 合约（通常 ≪ 53）。
3. **零跳过**：`getPooledBNB == 0 && lockedBNBs == 0` 的 credit 不再重复查询（同一 delegator 跨行复用结果）。

**Archive 要求**：BSC 快照块高通常比 tip 早数天；多数 publicnode 实例 **无** 该高度的 state（`historical state not available`）。`rpc-providers.json` 须标注 `archiveDepth`；Stake Hub 查询与 `eth_getBalance` **共用** 同一 archive-capable provider（实测 `bsc.drpc.org` 可用，但免费 tier 易 429——需 failover + 缓存）。

#### 5.3.4 Ethereum — Beacon 存款（`onchain-balance-hot@2.1`）

| 项 | 说明 |
|---|---|
| 场景 | HotCold 行 `ETH|ETH`，CSV 余额为大额整数（如 80,000 ETH），EOA `balance ≈ 0` |
| 链上证据 | 同一地址有多笔 `DepositContract.deposit()`（Etherscan Deposits tab） |
| 查询 | 信标链 API / 执行层 deposit 事件 + 汇总该地址已 deposit 的 ETH（32 ETH 整数倍） |
| 计入 | `stakedNative += Σ depositAmount`（在 `Height` 之前已确认的 deposit） |

**PR01JUN26 实证**：`0x32e11a20337ebc79abd0eeab2d91bafbd9591149` — CSV **80,000.017 ETH**，liquid **0.017 ETH**，Etherscan **40 × 2,000 ETH deposit**；账面与链上质押记录一致，native-only FAIL 为误报。

实现可分两阶段：

- **2.1 启发式**：若 `liquid << claim` 且地址存在 deposit 事件且 `Σ deposits ≈ claim` → PASS + `components.staked`
- **2.2 完整**：对接 beacon API（如 `beaconcha.in` / 自建 lighthouse）按 validator pubkey 查 effective balance

#### 5.3.5 Verdict 与 coverage 规则

| 条件 | Verdict |
|---|---|
| 所有支持 `(coin,network)` 对，`totalAccounted` 在容差内 | **PASS** |
| 某行 `totalAccounted` 超容差 | **FAIL**（须含 `components`） |
| 某行仅实现了 liquid，但 heuristic 检测到 likely staking（大额 claim + 近零 liquid + StakeHub/deposit 活动） | **WARN**（v1 行为；v2 应消减） |
| RPC archive 不可用，无法查 `Height` | **WARN**（`note: rpc archive unavailable`） |
| `(coin,network)` 尚无 native/staking 处理器 | **UNVERIFIABLE**（按 coin\|network 聚合行数，与 v1 相同） |

**Coverage 定义调整**：

```text
coverage = (# HotCold rows with full totalAccounted) / (# HotCold rows total)
```

"v1 的 4%"是因为只支持 `ETH|ETH` + `BNB|BSC` 且只查 liquid；v2 目标是在同样两对网络上把 **observation depth** 从 liquid-only 提升到 liquid+stake，而不是扩大 coin 种类。

#### 5.3.6 PR01JUN26 三个地址（设计文档附录级实证）

快照 `PR01JUN26`，BSC `#101590091`，ETH `#25218797`（2026-05-31 23:59:59 UTC）。

| 地址 | CSV | liquid | staking 解释 | v1 | v2 预期 |
|---|---|---|---|---|---|
| `0xbf83…ed7a` | 8,657,966 BNB | 0.087 BNB | BNB 已 delegate/claim 出 Stake Hub；EOA 长期近零 | FAIL | PASS（staked+liquid 或 WARN→PASS） |
| `0x32e11…1149` | 80,000 ETH | 0.017 ETH | 80k 已全部 Eth2 deposit（40×2k） | FAIL | PASS（staked=80k） |
| `0x86523…08d96` | 172,798 BNB | 90,932 BNB | Stake Hub pooled ≈ **75,466 BNB**（30/53 validators 已扫）；合计 ≈ **166,398 BNB**，残差 ≈ **6,400 BNB** 待扫完 | FAIL | PASS 或小幅 FAIL |

> 上述数字来自 bootstrap CLI 多 RPC 复核（2026-06）；写入 verification bundle 时可作为 golden fixture。

#### 5.3.7 版本路线图

| 版本 | 范围 | 里程碑 |
|---|---|---|
| `@1.0`（当前） | `ETH\|ETH`、`BNB\|BSC` liquid only | W2 bootstrap；已知误报 FAIL |
| `@2.0` | + BSC Stake Hub（pooled + unbonding） | 消除 BNB 质押误报；需 archive RPC + validator 缓存 |
| `@2.1` | + ETH beacon deposits | 消除 Eth2 存款误报 |
| `@3.0` | + ERC20 `balanceOf`（USDT/USDC/…） | 显著提升 HotCold 价值覆盖率 |

`onchain-balance-hot@2.x` **不替换** verifier id 的语义——仍是对 HotCold 逐行链上审计；仅扩展 `components` 与 RPC 依赖。升版后旧 snapshot 应用新版本**重跑**并对比 diff（§4 `version` 规则）。

---

## 6. 关键数据上链锚定

为了让 ardmere 自身的"诚实"也是可验证的（不是"trust me bro"），所有**根级证据**都做链上锚定。

> **决策记录**：本节方案遵循 [ADR-003](./decisions.md#adr-003-冷存储策略不引入-arweave--ipfs仅关键摘要上链)——不引入 Arweave / IPFS 等独立冷存储，所有原始数据保留在我方 S3 + Postgres，链上仅写摘要 Root。

### 6.1 上链对象与频率

| 锚定对象 | 摘要内容 | 上链频率 | 大小 |
|---|---|---|---|
| **单期快照锚定** | 同一笔 tx 携带 `artifactBundleRoot` + `verificationBundleRoot` 两个 Merkle Root，外加交易所自报 root、verifier 结论摘要 | 每期新快照跑完全部 verifier 后（约 1 次/月） | 2×32 B (roots) + 元数据 |
| **每周心跳** | 当周所有 artifact / verification 的累积 root + 我方公钥 | 每周一 00:00 UTC | 64 B |

> 上链的不是原始大数据，而是**摘要 Merkle Root**——既不可篡改，又便宜。原始数据通过 ardmere 自家 API（`GET /artifacts/:sha256`）提供，用户拿到后用链上 root 自验。
>
> **设计决策（v2）**：artifact 与 verification 两个 root **合并为 1 笔链上交易**，避免同一 `snapshotId` 出现两笔独立 anchor 导致的前端/索引复杂度，且 gas 成本几乎不变（~$0.001/tx on Base）。

### 6.2 链选型

**主锚定链：Base**（OP Stack L2，Chain ID `8453`）—— [ADR-002](./decisions.md#adr-002-主锚定链base)

理由：
- Gas 极低（~$0.001 per tx），可承受高频心跳
- EVM 兼容，工具链成熟
- 数据可用性继承自 Ethereum L1（通过 batch posting）
- 公共 RPC 多家提供（不与 §7 验证用 RPC 强耦合）

**未来可选**：Bitcoin OP_RETURN（最高强度但贵）作为年度审计锚，不在 MVP 范围。

### 6.3 锚定合约（极简）

```solidity
// 部署在 Base mainnet (chain id 8453)
contract ArdmerePoRAnchor {
    uint8 public constant SCHEMA_VERSION = 2;

    event SnapshotAnchored(
        bytes32 indexed snapshotId,           // keccak256("PR01APR26")
        bytes32 indexed exchangeTag,          // keccak256("binance") — 便于按交易所过滤
        string exchange,                      // 明文，如 "binance" / "okx"（BaseScan 直接可读）
        uint32 periodSeq,                     // 该交易所第几期 (1 = 最早一期)
        uint64 snapshotTime,                  // 交易所快照 UTC unix
        uint32 btcBlockHeight,                // BTC 时间锚
        bytes32 exchangeMerkleRoot,           // 交易所自报 Merkle root
        bytes32 artifactBundleRoot,           // ardmere artifact bundle Merkle root
        bytes32 verificationBundleRoot,       // ardmere verification bundle Merkle root
        uint8 verdictSummary,                 // verifier 结论 bitfield
        uint16 coverageBps,                   // 链上审计覆盖率 × 10_000
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

> 极简刻意为之——**合约本身没有状态、没有升级、没有外部存储引用**，只发 event。前端通过 `eth_getLogs` 查询历史；每个 `(exchange, snapshotId)` 对应**唯一一条** `SnapshotAnchored` 记录（1 tx = 1 period）。`artifactBundleRoot` 与 `verificationBundleRoot` 均不可为零。

**已部署实例（Base Sepolia testnet）** — 详见 [`deployments.md`](./deployments.md)：

| | |
|---|---|
| Contract | [`0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9`](https://sepolia.basescan.org/address/0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9) |
| Deploy tx | [`0x3e844a…2464c`](https://sepolia.basescan.org/tx/0x3e844a1099f932638f695ce0d3045f78e5cc4e5def63d0edfb40bb603bb2464c) |
| First anchor tx (`PR01JUN26`) | [`0x3ce248…32bb`](https://sepolia.basescan.org/tx/0x3ce248b76d7638ea3326b93b6ef731fa40eb07f52c8397ab00633079614932bb) |

### 6.4 验证流程闭环

```text
任何第三方独立验证 ardmere 的方法：

1. 拿 snapshotId（如 PR01APR26）
2. 在 Base 上查 ArdmerePoRAnchor 的 SnapshotAnchored event → 得到 artifactBundleRoot + verificationBundleRoot
3. 从 ardmere API 拿原始 bundle（GET /artifacts/:sha256、GET /verifications/:id）
4. 分别重算两个 bundle 的 Merkle Root，比对链上对应字段 → ✅ 数据未被篡改
5. 重新执行任一 verifier（代码 GitHub 开源 + 版本号匹配） → 得到与 bundle 一致的结论 → ✅ 逻辑可重放
```

> **承认的取舍**：此方案下，若 ardmere 服务下线，原始 verification 数据将不再可访问；但**链上的 root + 任何已被第三方下载缓存的 bundle**仍可用于事后比对验证。如未来需更强保证，可在 V2 增加独立的 Arweave 镜像服务（仅锚定 root 不变）。

---

## 7. 公共节点池策略（无自建 archive）

**设计约束**：所有链上查询走第三方公共 RPC，**不自建 Erigon / electrs / 任何 archive node**。

### 7.1 节点池抽象

```ts
interface ChainRpcPool {
  // 按链选 endpoint，自动 failover + 限流 + 缓存
  getBalance(network: Network, address: string, blockHeight: number): Promise<bigint>;
  getErc20Balance(network: Network, token: string, address: string, blockHeight: number): Promise<bigint>;
  getBlockMeta(network: Network, height: number): Promise<{ hash: string; timestamp: number }>;
}
```

### 7.2 各链 provider 选型

| 网络 | 主 provider | 备用 1 | 备用 2 | archive 支持 |
|---|---|---|---|---|
| Ethereum | Ankr Public | Llamanodes | publicnode.com | ✅ 全部支持 |
| BSC | publicnode.com | bsc-dataseed.binance.org | Ankr | ✅ |
| Polygon | Ankr | publicnode | Polygon Foundation RPC | ✅ |
| Arbitrum | Ankr | Arbitrum Foundation | publicnode | ✅ |
| Base | Base public RPC | Ankr | publicnode | ✅ |
| BTC | mempool.space API | blockstream.info API | btc.getblock.io | ✅ |
| Solana | Helius free tier | publicnode | Triton | ⚠ 受 epoch 限制 |
| Tron | TronGrid | Nile public | publicnode | ⚠ 历史余额 API 受限 |
| Algorand | AlgoNode | PureStake free | Algorand Foundation | ✅ indexer |
| Aptos | Aptos Labs full node | publicnode | - | ✅ |

> **完整列表**：维护在 `config/rpc-providers.json`，按 `(network, capability)` 分组，每条记录 `(url, weight, rateLimit, lastFailedAt, costModel)`。

### 7.3 节点池策略

| 策略 | 实现 |
|---|---|
| **Failover** | 主节点 5xx / timeout / 限流 → 自动切到备用，5 分钟内不再尝试 |
| **限流** | 按 provider 配额令牌桶，超限自动降级到下家 |
| **结果缓存** | `(network, address, blockHeight)` 是不变量 → 永久缓存到 Postgres（命中率会很高，因为同一快照不会重算） |
| **冗余校验** | 关键查询（如 HotCold 大额地址）**对同一条同时问 2 个 provider，结果不一致 → 标记 `WARN`** |
| **可观测性** | 每个 provider 的成功率 / 延迟 / 限流次数实时面板，月度复盘换池 |
| **Sandbox 验证** | 部署时跑一遍"对历史已知 BTC 创世块查余额"等基线测试，确认 provider 行为符合预期 |

### 7.4 公共节点的局限与应对

| 局限 | 应对 |
|---|---|
| **历史余额查询有限制**（部分免费 tier 只支持最近 128 个区块）| 在 `rpc-providers.json` 标注 `archiveDepth`，超出时自动选 archive 支持 provider |
| **BSC 快照块 archive 稀缺**（PR01JUN26 `#101590091`：publicnode 返回 `historical state not available`；仅个别 provider 如 drpc 可用，且易 429） | Stake Hub 查询与 `eth_getBalance` 绑定同一 archive provider；结果永久缓存 `(network, address, height, component)`；StakeHub validator 列表按 height 缓存 |
| **每秒请求数低**（公共节点典型 5~20 QPS）| Deposit 全量审计放慢节奏，跑几小时也无所谓；HotCold 优先 |
| **可能被限流到 IP** | 多 provider 并行 + 出口 IP 池（可选 Cloudflare Workers 分布式调用）|
| **结果可能被劫持** | 关键查询双 provider 冗余比对（§7.3）|
| **Solana / Tron 等链历史 API 受限** | 这些链先做 `coverage < 100%` 验证，UI 标注"该网络仅校验最近余额" |

### 7.5 完全公开节点是否够用？

| 场景 | 公共节点是否够 |
|---|---|
| HotCold 100% 校验（10³ 地址 × 多链） | ✅ 一晚跑完 |
| Deposit 按价值采样 99% 覆盖（约 10⁵ 地址）| ✅ 一周跑完 |
| Deposit 全量 100%（10⁷ 地址） | ⚠ 公共节点会被打爆，需排队几个月或部分用付费 tier |

> 我们的 SLA 承诺到"按价值覆盖 99%"，不承诺"100% 地址覆盖"。这是基于公共节点策略下**理性、可达成**的目标。

---

## 8. 物理架构

```
┌─────────────────────────────────────────────────────────────────────┐
│  Scheduler (cron)                                                    │
│  ├── 6h 轮询 BAPI snapshot list（发现新快照）                          │
│  └── 每周一 00:00 UTC 心跳上链                                         │
└──────────────┬──────────────────────────────────────────────────────┘
               ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Fetcher Pool                                                        │
│  ├── BinanceBapiFetcher       → Artifact{kind:"bapiSnapshot"}        │
│  ├── BinanceWalletZipFetcher  → Artifact{kind:"walletZip"}           │
│  ├── BtcAnchorFetcher         → Artifact{kind:"btcBlockMeta"}        │
│  └── (future) AddressSigFetcher / ProofCsvFetcher / AttestationFetcher│
└──────────────┬──────────────────────────────────────────────────────┘
               ▼  写入 Artifact Store + 我方签名
┌─────────────────────────────────────────────────────────────────────┐
│  Artifact Store                                                      │
│  ├── 内容寻址：S3（仅此一处，决策 ADR-003）                           │
│  └── 元数据：Postgres                                                 │
└──────────────┬──────────────────────────────────────────────────────┘
               ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Verifier Runtime                                                    │
│  ├── 加载已注册 verifier（含占位 stub）                               │
│  ├── 调度：按 verifier.requires 准备 ctx                              │
│  ├── 执行：纯函数 verify()，访问 RpcPool (§7)                         │
│  └── 落档：Verification + 我方签名                                    │
└──────────────┬──────────────────────────────────────────────────────┘
               ▼  artifact + verification bundle roots → 上链（1 tx）
┌─────────────────────────────────────────────────────────────────────┐
│  On-chain Anchor (Base, chain id 8453)                              │
│  └── ArdmerePoRAnchor.anchorSnapshot(exchange, snapshotId, …)       │
└──────────────┬──────────────────────────────────────────────────────┘
               ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Public API + 静态站点                                                │
│  GET /snapshots                  → 历史 + 当期                        │
│  GET /snapshots/:id              → 当期所有 verification 结果         │
│  GET /verifications/:id          → 单次验证 + 可重放证据 + 链上 tx     │
│  GET /artifacts/:sha256          → 原始文件（带我方签名 + 链上锚定）   │
│  GET /anchors                    → 链上锚定历史（含 Arweave 引用）     │
│  POST /user-inclusion/verify     → 用户 zip 浏览器内 WASM 验证（L3）   │
└─────────────────────────────────────────────────────────────────────┘
```

技术栈建议：

| 层 | 选型 | 理由 |
|---|---|---|
| 后端语言 | **Go** | 直接复用 `binance/zkmerkle-proof-of-solvency` 的 verifier 包；并发友好 |
| 数据库 | Postgres | 关系明确，verification 与 artifact 关联强 |
| 对象存储 | S3（仅此一处） | 决策 ADR-003 简化架构 |
| 链上锚定 | Base + Foundry | gas 便宜 + 工具链成熟 |
| 签名 | ed25519（off-chain）+ secp256k1（on-chain） | 双密钥分离，私钥分级管理 |
| 前端 | 现有 ardmere.org 静态站 + 一个 React island 用于交互式验证视图 | 与 [`DESIGN.md`](./DESIGN.md) 的 terminal 风格一致 |
| 部署 | GitHub Actions + Fly.io / Railway（后端）+ GitHub Pages（前端） | 零运维 |

---

## 9. 前端展示模型

把"验证服务的可信度"做成视觉层第一原则。最终用户看到的样子：

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

每个图标都点得开：

- ✅ → 展开 finding 列表 + 链上 anchor tx
- ⚠ → 解释为什么是 warn（"自报"等），不是 fail
- ❌ → 展开"为什么验证不了 + 我们在追踪什么样的事件触发激活"
- ⏳ → 展示当前 coverage 实时百分比 + 预计完成时间

---

## 10. 扩展场景演练

### Case A：未来 Binance 在 PoR 公告里给出地址签名

```text
1. 写 fetcher：解析公告 → 下载 sig 文件 → 落档 Artifact{kind:"addressSig"}
2. 升级 verifier：address-ownership@0 → address-ownership@1
   - requires: ["walletZip", "addressSig"]
   - verify(): 对每条 sig 用 secp256k1/ed25519 验签
3. 注册即生效；前端自动从 ❌ → ✅
4. 历史快照自动重跑（因 version 升级），如旧快照无 sig artifact，输出 PARTIAL
5. 新结论上链锚定，建立"激活时间线"
```

**核心代码改动 = 一个新 fetcher + 一个新 verifier + 注册一行**。

### Case B：未来出现可公开下载的第三方 PoR attestation

```text
1. 写 fetcher：定时拉第三方机构官网 attestation PDF / JSON
2. 写 verifier：third-party-attestation@1
   - 验证签名机构是否在白名单
   - 验证 attestation 中的 Merkle Root 与我们 BAPI 拉到的 root 一致
3. 注册即生效
```

### Case C：未来 Binance 公开 zk 全局 proof

```text
1. fetcher 抓 proof.csv / vk
2. global-zk-proof@1 用 binance/zkmerkle-proof-of-solvency 的 Go verifier
   （可编 WASM 跑在边缘，与 §8 主进程隔离）
3. 注册即生效
```

### Case D：增加新交易所（OKX / Coinbase）

```text
1. 抽象 Exchange 维度：
   - Snapshot 增加 exchange 字段
   - 每个 Exchange 有自己的 Fetcher 集合 + 默认 Verifier 集合
2. Verifier 接口本身不变（它只看 artifact + claim），所以大多数 verifier 可跨交易所复用：
   - internal-consistency 完全通用
   - onchain-balance-* 完全通用
   - address-ownership 各家签名格式不同 → 每家一个版本
3. 前端按 Exchange 分 tab；公共 verifier 排序对比
```

---

## 11. 安全与威胁模型

| 威胁 | 缓解 |
|---|---|
| Binance 临时换源数据（CDN 投毒） | sha256 + 上链锚定，事后可比对 |
| 公共 RPC 节点撒谎 | 双 provider 冗余比对，不一致即 WARN |
| 我们自己被攻陷篡改结论 | 私钥硬件隔离 + 链上 anchor 不可逆，篡改后无法回填历史 |
| Verifier 逻辑 bug | version 强制 + 历史快照定期重跑回归 |
| 抓取被反爬 | 多入口 fetcher + 客户端浏览器自助 zip 验证（L3）作为 backup |
| 用户隐私（zip 上传）| 浏览器内 WASM 全本地验证，零上传 |

签名密钥分级：

| 用途 | 密钥 | 存放 |
|---|---|---|
| Artifact / Verification 离线签 | ed25519 #1 | 服务进程内（rotate 季度）|
| 上链锚定（合约调用） | secp256k1 #2 | KMS / HSM，仅 anchor 进程可访问 |
| 站点身份 / 域名 DNSSEC | secp256k1 #3 | 离线，年度 rotate |

---

## 12. MVP 计划

| 周 | 工作 | 产出 |
|---|---|---|
| W1 | Domain model + Verifier 接口 + Fetcher (BAPI + walletZip) + Artifact Store + 4 个核心 verifier (`integrity` / `internal-consistency` / `btc-anchor` / `solvency-claim`) + 4 个占位 stub | CLI 跑通：`ardmere verify PR01APR26` 输出 8 维 verdict |
| W2 | `onchain-balance-hot@1`（liquid native）+ Anchor 合约 + 第一笔 anchor tx；**@2.0 staking-aware 设计**见 §5.3 | bootstrap 跑通；BNB/ETH 质押误报在 @2.0 消除 |
| W3 | `onchain-balance-deposit@1`（采样 + 持续运行）+ Public API + 前端 verdict 卡片 | 网站可访问，看到 6/10 维 verdict |
| W4 | 浏览器 WASM 用户 inclusion 验证（L3 self-service） + 文档完整化 + 上线 | 完整 MVP，可对外发布 |

> `onchain-balance-deposit` 跑全量需数周，但在 W3 就上线"实时进度条"——**不要等它跑完才上线**。

---

## 13. 信用积累路径

ardmere 的核心资产 = 「在公开数据条件下，我们做到了能做的最严格」的可信记录。

| 阶段 | 信用来源 |
|---|---|
| Month 1 | 上线 + 第一笔 anchor，证明系统是活的 |
| Month 3 | 三期连续 verdict 一致；前端展示对比图 |
| Month 6 | 当 Binance 自己的 ratio 出现波动时，我们的 L1+L2 第一时间复现并公示，建立"独立于 Binance 的发声渠道"的可信度 |
| Year 1 | 接入 OKX / Coinbase / Bitget；逐渐成为 PoR 行业的 reference verifier |
| 长期 | 当 Binance 终于公开签名 / 全局 proof 时，我们立即激活对应 verifier，并在前端高亮"激活时间线" |

---

## 14. 回答核心问题

> **如何保留扩展性以便将来验证目前暂时验证不了的信息？**

> **把"无法验证"也建模成"一种验证"**——每个未来想做的验证维度，今天就以"占位 verifier"的形式存在于注册表里、出现在 API 响应里、显示在前端 UI 上。当外部条件成熟（Binance 提供官方下载渠道、第三方机构发布可公开下载的 attestation……）时，**把那个 verifier 从 v0 stub 升级到 v1 真实实现**，不需要动任何上游 schema、不需要改前端代码、不需要重设计数据库——所有历史快照会自动被新 verifier 重跑一遍并对比。

> 对应的，**关键摘要全部上链**让 ardmere 自身的"诚实"也可独立验证；**全公共节点池策略**让基础设施零信任、零自营，避免"只能信我"的悖论。

---

## 附录 A — 与现有 ardmere 站点的关系

本服务的前端模块嵌入 [`DESIGN.md`](./DESIGN.md) 描述的 ardmere.org，作为新增 section：

- 视觉延续 terminal / matrix / glitch 体系
- verdict 卡片做成"扫描中"的 ASCII 进度条
- 每个 ❌ UNVERIFIABLE 的展开做成 `cat /dev/binance | grep signature` 风格的伪命令输出
- 配色沿用 `--green / --cyan / --red`，不引入新色

## 附录 B — 决策记录

详见 [`decisions.md`](./decisions.md)，五项 ADR 已于 2026-05-05 拍板：

| ADR | 决策 |
|---|---|
| ADR-001 后端语言 | **Go** |
| ADR-002 主锚定链 | **Base** (chain id 8453) |
| ADR-003 冷存储 | **不引入独立冷存储，仅关键摘要上链** |
| ADR-004 用户验证形态 | **全本地浏览器 WASM** |
| ADR-005 多交易所路线 | **MVP 后优先支持 OKX** |

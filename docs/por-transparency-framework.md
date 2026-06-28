# PoR Stage 分级框架

本文定义 ardmere 对交易所 Proof-of-Reserves（PoR）的透明度评估体系。对具备最低可用 PoR 的交易所，**最终评级是 PoR Stage（0 / 1 / 2）**，衡量其 PoR 在信任最小化上的成熟度。没有最低可用 PoR 的交易所标记为 **Pre-Stage — No usable PoR**，不进入 Stage 0。

本框架参考 [L2BEAT Stages Framework](https://l2beat.com/stages) 的设计思想：Stage 不是安全认证，也不是技术先进程度排名，而是回答一个更基础的问题：

> 用户最终还要信任谁？

`Technology Generation（Gen）` 和 `Evidence Level（E）` 不是与 Stage 平级的评级。它们是解释 Stage 判定的两个辅助视角：

- **Gen**：说明交易所使用什么证明技术。
- **E**：说明外部观察者能取得、复现、交叉验证多少证据。

---

## 1. 框架目标

### 1.1 本框架评估什么

PoR Stage 评估交易所 PoR 的**信任最小化程度**：


| 问题        | 含义                                                                                |
| --------- | --------------------------------------------------------------------------------- |
| 证据是否公开    | 第三方是否能取得 summary、wallet_address_list、global_proof、wallet_ownership_proof、vk/config 等 artifacts |
| 证据是否可复验   | 第三方是否能无需交易所授权独立重跑核心验证                                                             |
| 证明系统是否免信任 | 证明系统本身是否依赖可信设置、不可验证参数、专有实现或交易所掌握的 toxic waste；若有信任假设，是否公开可验证                      |
| 证明约束是否充分  | 证明系统是否能自证关键安全性质，例如没有插入负净值虚假用户、没有用负余额抵消真实负债、用户余额承诺与全局负债汇总一致                        |
| 历史是否可改写   | root/proof/commitment 是否有 canonical 记录和不可篡改发布层                                    |
| 信任假设是否收敛  | 用户最终信交易所叙事、审查流程，还是数学与公开数据                                                         |


### 1.2 本框架不评估什么

Stage 不直接评估：

- 单个 verifier 工具当前实现了多少检查。
- 某个 RPC、链、token 的工程覆盖程度。
- 代码是否完全没有 bug。
- 交易所整体财务健康、公司治理或监管合规。

Stage 高不代表“无安全风险”；Stage 低也不代表“一定作恶”。Stage 只说明：**PoR 披露把用户从多少信任假设中解放出来。**

---

## 2. Stage 分级总览


| Stage         | 名称                    | 一句话定义                                            | 典型证据状态                                          | 核心信任边界                                                              |
| ------------- | --------------------- | ------------------------------------------------ | ----------------------------------------------- | ------------------------------------------------------------------- |
| **Pre-Stage** | No usable PoR         | 没有最低可用 PoR，不进入 Stage                             | 无 snapshot / summary / proof，或 PoR 已停止          | 只能信交易所口头声明、品牌、财报或监管披露                                               |
| **Stage 0**   | Trust the Exchange    | 有 PoR，但第三方无法独立复现全局偿付关系                           | summary 或用户 inclusion 为主                        | 仍必须信交易所的资产归属、负债完整性或证明参数声明                                           |
| **Stage 1**   | Verifiable Disclosure | 关键 artifacts 公开、可复验，但 proof 语义仍主要以披露边界为准         | wallet_address_list、wallet_ownership_proof、global_proof 与参数来源均公开；通常为周期性快照 | 不再信交易所单方声明 wallet_address_list 控制权、global_proof 可得性和可信设置诚实性，但仍信 proof 边界足够反映业务风险 |
| **Stage 2**   | Trust-minimized PoR   | 证明与数据 permissionless 可验证，证明约束与业务风险一致，风险窗口被高频锚定压缩 | 双侧公开 + 低门槛 inclusion proof + 业务一致约束 + 链上锚定 + DA + 高频/事件触发更新 | 进一步不再信交易所自行解释业务风险、维护唯一历史版本、持续提供数据或任意拉长未证明窗口                         |


**当前行业状态**：主流 CEX 最高约为 **Stage 1**；尚无完整 **Stage 2**。

**Stage 0 与 Stage 1 的核心区别**：Stage 0 虽然已有 PoR 工件，但关键断言仍由交易所单方提供；Stage 1 则要求交易所公开 wallet_address_list、wallet_ownership_proof、global_proof 和可信设置 / 参数来源，使第三方可以独立复验这些核心断言。换言之，Stage 1 的跃迁不是“更多披露”，而是交易所单方叙事被公开证据约束。

这里的“关键断言”包括：交易所确实控制 wallet_address_list 中的地址；公开的链上余额确实属于本期 reserves；global_proof 确实对应本期 liability summary；proof 使用的 vk/config 与公开参数一致；若使用可信设置，交易所没有单方掌握 toxic waste。

**Stage 1 与 Stage 2 的核心区别**：Stage 1 解决“这一期 artifacts 是否可公开复验”，并要求 proof 边界清楚；Stage 2 进一步要求证明约束与交易所实际业务一致，并能自证没有插入负净值虚假用户、没有用负余额抵消真实负债。例如涉及借贷、保证金、组合保证金或抵押品折扣的交易所，PoR 必须在电路或等价证明中约束相应的风险控制逻辑。同时，Stage 2 还要求用户侧 inclusion proof 具备低门槛验证路径，root/proof/commitment 有高频锚定、事件触发增发或准实时证明，把快照间隔中的择时、替换和历史改写风险降到更低。

---

## 3. Stage 判定方法

### 3.1 判定原则


| 原则                              | 含义                                                |
| ------------------------------- | ------------------------------------------------- |
| Stage-first                     | 对外评级以 Stage 为主标签；Gen / E 只解释 Stage                |
| 门槛制                             | 先按硬性 checklist 判 Stage，不用平均分升档                    |
| Shortest-board model            | PoR 是短板系统；关键短板会阻断升 Stage 或下调置信度                   |
| Evidence-first                  | 评级必须绑定可取得的原始工件、哈希、URL、报告或链上数据                     |
| Public-by-default               | 无登录、无需 cookie、无需人工申请的数据权重大于登录后用户数据                |
| Reserve 与 liability 分离          | wallet_address_list、负债侧 Merkle/ZK、summary ratio 是不同数据平面 |
| UNVERIFIABLE 显式化                | 缺数据不是 PASS；必须写清无法验证的维度                            |
| Snapshot-scoped                 | 每个评级针对某一期快照，不能默认外推到所有历史/未来期                       |
| No overclaiming                 | Merkle inclusion PASS 不能表述为全所偿付能力 PASS            |
| Business-consistent constraints | PoR 证明约束必须与交易所实际业务一致；负净值虚假用户、负余额抵消、借贷、保证金、抵押品折扣等风险不能只靠文字说明 |
| User-verification usability     | inclusion proof 的验证流程应低门槛、可解释、可导出，不能要求普通用户执行复杂命令或理解底层电路 |
| Machine-readable first          | JSON/CSV/ZIP/schema/API 等机器可读工件优先于不可解析的 PDF 或营销页面 |
| Audit taxonomy                  | 必须区分技术核验、AUP/商定程序、有限保证、合理保证/完整财务审计                |
| Frequency-aware                 | 透明度不仅看单期快照，也看发布频率、事件触发披露、历史版本可归档性和抗择时能力           |


### 3.2 Pre-Stage — No usable PoR

**定义**：交易所没有最低可用 PoR，或历史上有 PoR 但当前已经停止。用户可以参考交易所品牌、传统财报或监管披露等其他信任来源，但无法取得最基本的 PoR 工件进行复核。

**满足任一即为 Pre-Stage**：


| #   | 情况       | 判定标准                                                       |
| --- | -------- | ---------------------------------------------------------- |
| P-1 | 无持续 PoR  | 近 12 个月无公开 PoR                                             |
| P-2 | PoR 已停止  | 历史有 PoR，但当前不再发布                                            |
| P-3 | 只有营销声明   | 只有 “100% backed” 等笼统表述，无 snapshot、summary、root、proof 或审查报告 |
| P-4 | 只有传统财报   | 有财务审计或上市公司披露，但未提供可映射用户资产负债侧的 PoR 工件                        |
| P-5 | 只有链上余额展示 | 只有钱包余额看板，但无用户负债侧 summary / inclusion / proof               |


**Pre-Stage 与 Stage 0 的区别**：


| 状态            | 用户信任什么                  | 最低 PoR 证据                                    |
| ------------- | ----------------------- | -------------------------------------------- |
| **Pre-Stage** | 品牌、财报、监管披露或其他非 PoR 信任来源 | 无                                            |
| **Stage 0**   | 交易所叙事 + 最低 PoR 工件       | 有 snapshot、summary、用户 inclusion 或基础 verifier |


**代表**：

- Coinbase（交易所级）：Pre-Stage（有上市公司财务披露，但未提供交易所级 PoR 工件）。
- KuCoin：Pre-Stage（PoR 停止）。

### 3.3 Stage 0 — Trust the Exchange

**定义**：交易所发布 PoR，用户最多能验证自身 inclusion 或阅读 summary，但独立第三方无法复现全局资产—负债关系。

**硬性门槛（全部满足）**：


| #    | 要求     | 判定标准                                                       |
| ---- | ------ | ---------------------------------------------------------- |
| S0-1 | 持续 PoR | 近 12 个月有公开 PoR，且未停止                                        |
| S0-2 | 快照边界   | 有明确 snapshot 时间、资产覆盖范围                                     |
| S0-3 | 最低披露   | 至少满足以下之一：自报 summary（储备率 + 覆盖资产）；或用户 Merkle inclusion proof |
| S0-4 | 可复验性下限 | 若提供用户 proof：leaf/hash 规则公开，proof 可导出，第三方可用用户提供的 proof 本地复验 |
| S0-5 | 工具或规范  | 有开源 verifier，或足够详细的算法说明                                    |


**允许缺失（因此无法升到 Stage 1）**：

- 无公开 wallet_address_list。
- 无公开 global_proof。
- 无公开 wallet_ownership_proof。
- 无独立第三方审查，或审查类型/范围未明示。
- 无链上锚定。

**Stage 0 信任假设**：用户信任“交易所说多少储备就是多少”。

**代表**：

- Bitget：weak Stage 0（64-bit 截断 hash，无独立 attestation）。
- Bybit：Stage 0（用户 Merkle，但无公开 wallet_address_list / ZK）。
- Gate.io 生产环境：Stage 0（summary 公开，ZK 包登录，无 wallet_address_list）。

### 3.4 Stage 1 — Verifiable Disclosure

**定义**：资产侧 wallet_address_list 及 wallet_ownership_proof、负债侧 global_proof 均公开可得，证明参数 / 可信设置来源可验证，任何第三方原则上都能独立复验；但 proof/root/vk 的 canonical 记录和历史版本仍依赖交易所或少数审查方。这里的“公开”指无需登录、注册、cookie、API key 或人工申请即可取得。

**硬性门槛（全部满足）**：


| #     | 要求           | 判定标准                                                                                                 |
| ----- | ------------ | ---------------------------------------------------------------------------------------------------- |
| S1-1  | 满足 Stage 0   | 全部 Stage 0 门槛                                                                                        |
| S1-2  | 资产侧公开        | 公开 wallet_address_list + 块高，可链上复核                                                                      |
| S1-3  | 地址控制权证明      | 公开 wallet_ownership_proof 或其他可验证控制证明，且可批量验证                                                              |
| S1-4  | 负债侧公开        | 公开 global_proof + root/vk/config，可本地重跑或复验结构；不得要求登录、注册、cookie、API key 或人工申请                         |
| S1-5  | 无许可核心工件      | summary、wallet_address_list、wallet_ownership_proof、global_proof、vk/config 等核心 artifacts 无需登录、注册、cookie、API key 或人工申请即可下载 |
| S1-6  | 独立审查         | 有独立第三方审查，且类型明示（技术核验 / AUP / 有限保证 / 合理保证）                                                             |
| S1-7  | 双向交叉绑定       | summary ↔ wallet_address_list aggregate，且 summary ↔ proof root / commitment 可比对                       |
| S1-8  | 参数透明         | 若使用可信设置，必须公开可验证 transcript / MPC ceremony；或采用无需可信设置的证明系统                                             |
| S1-9  | 历史可归档        | 历史 snapshot artifacts 可落档，不是仅当期网页展示                                                                  |
| S1-10 | 证明边界         | 明确披露 proof 覆盖哪些负债、不覆盖哪些表外/链外项目                                                                       |


**Stage 1 等价表述**：

> 除实现 bug 外，要让独立第三方**长期无法发现**偿付关系被伪造，必须依赖交易所与审查方合谋，或依赖交易所对 server 托管 artifacts 的单边替换/删除能力；不能再依赖“交易所确实控制 wallet_address_list 中的地址”或“交易所没有掌握 toxic waste”这类不可验证假设。

**仍允许存在的关键短板（因此无法升到 Stage 2）**：

- proof / root / vk 无链上锚定。
- global_proof 需要登录、注册、cookie、API key 或人工申请才能取得，只有 sample-only 包，或无法与当期 wallet_address_list / summary 绑定。
- 仅月度快照，无事件触发增发。
- 审查方与技术提供方存在未披露利益冲突。

**Stage 1 信任假设**：用户不再只信交易所一句话，不再信交易所确实控制其 wallet_address_list，也不再信交易所没有掌握可信设置秘密；但仍要信交易所不会替换当期 artifacts、不会对不同观察者 split-view、不会隐瞒 liability 口径。

**代表**：

- OKX：**Stage 1**（wallet_address_list + wallet_ownership_proof + global_proof，proof 仍链下分发）。

### 3.5 Stage 2 — Trust-minimized PoR

**定义**：在 Stage 1 双侧公开基础上，用户侧 inclusion proof 具备低门槛验证路径，证明约束与交易所实际业务一致，root / proof / commitment 有 canonical 链上或 DA 记录；任何第三方无需交易所授权即可 permissionless 复验；已锚定历史不可单方改写；并通过高频锚定、事件触发增发或准实时证明压缩两次完整 PoR 之间的风险窗口。

**硬性门槛（全部满足）**：


| #    | 要求                | 判定标准                                                                           |
| ---- | ----------------- | ------------------------------------------------------------------------------ |
| S2-1 | 满足 Stage 1        | 全部 Stage 1 门槛                                                                  |
| S2-2 | 业务一致约束            | 证明约束覆盖交易所实际业务风险，并能自证没有插入负净值虚假用户、没有用负余额抵消真实负债；涉及借贷、保证金、组合保证金、抵押品折扣、价格参数或风险限额的业务，必须在电路或等价公开证明中被约束 |
| S2-3 | 链上锚定              | root / proof root / commitment 上链，且有不可删除的历史 commitment 链                       |
| S2-4 | Data Availability | proof 包、vk、config、wallet_address_list bundle 在稳定公开层（链上、IPFS、第三方镜像/DA）长期可得       |
| S2-5 | 低门槛用户验证          | inclusion proof 可导出、可本地复验，并提供 Web/WASM/GUI 一键验证、清晰错误提示；验证结果与公开 root / global_proof / summary 绑定 |
| S2-6 | Permissionless    | 第三方不登录、不申请 API key 即可重跑核心验证；若有 ZK，prefer 链上 verifier 或等价公开验证路径                 |
| S2-7 | 参数与版本锚定           | vk/config/verifier 版本的哈希或 commitment 被链上/DA 固定，变更有 audit trail                 |
| S2-8 | 抗择时频率             | 对头部 CEX：至少满足“周度完整 PoR + 日度 root/commitment 锚定 + 事件触发增发”之一；长期仅月度快照不足以维持 Stage 2 |
| S2-9 | 审查独立性             | 审查方与技术提供方利益冲突已披露；不得“运动员兼裁判员”且无说明                                               |


**Stage 2 等价表述**：

> 任何用户都能低门槛验证自己的 inclusion proof，任何外部观察者都可以独立验证“这一期 canonical root 对应的 proof 是否成立”，且该 proof 能自证关键安全性质并约束交易所实际业务风险；交易所无法单方面改写已锚定历史，也不能任意拉长未被重新证明的风险窗口。

**强加分项**：

- 链上 verifier 合约可直接验证 ZK proof。
- proof generation / verifier 版本可追踪。
- 支持连续证明或高频快照（日度或更高）。
- 历史 commitment 链完整，proof/vk/root 变更有 audit trail。
- 对借贷、保证金、衍生品净额、抵押品折扣和价格参数有标准化 proof schema。

**当前状态**：主流交易所尚无完整 Stage 2。第三方归档者可以把已抓取 artifacts 的 root 上链，但这只能证明“该第三方见过并归档过这组 artifacts”，**不能**替代交易所官方 Stage 2。

**Stage 2 信任假设**：用户主要信任密码学证明、公开数据和链上 canonical 记录。

Stage 2 是 **trust-minimized**，不是 trustless。它仍依赖证明系统实现正确、PoR 覆盖口径真实、链外资产和表外负债被如实披露，以及审查范围没有被误读。若交易所实际经营借贷、保证金或衍生品业务，但 PoR 只证明现货余额或简单总负债，则不能视为 Stage 2。

### 3.6 Stage 置信度修正

Stage 是成熟度标签，不是安全分数。可在 Stage 不变时附加置信度标记：


| 标记                | 含义                                              |
| ----------------- | ----------------------------------------------- |
| high confidence   | 双侧公开完整、审查独立、实现风险低                               |
| medium confidence | 双侧公开完整，但链下分发、用户可及性或历史归档不足                       |
| low confidence    | 存在 sample-only、登录墙、已知 HIGH 缺陷或 attestation 范围极窄 |


**短板处理规则**：


| 情况                                  | 处理                          |
| ----------------------------------- | --------------------------- |
| 只有样例 proof，非当期生产                    | 不给生产 Stage 1，标记 sample-only |
| proof 包需要登录、注册、cookie、API key 或人工申请 | 若负债侧无法公开复现，最高 Stage 0       |
| inclusion proof 不可导出、不可本地复验，或只能黑盒验证 | 最高 Stage 0                  |
| inclusion proof 可复验但验证流程高摩擦              | 最高 Stage 1                  |
| 无 wallet_ownership_proof            | 最高 Stage 0                  |
| trusted setup 不透明，或无可验证 transcript  | 最高 Stage 0                  |
| 无独立第三方审查                            | 最高 Stage 0                  |
| proof/root/vk 仅 server 托管、无归档/锚定    | 最高 Stage 1                  |
| 重大未修复实现缺陷影响 proof 语义                | Stage 不变，置信度下调              |
| 审查方与技术提供方未披露利益冲突                    | Stage 不变，置信度下调              |


---

## 4. 用于解释 Stage 的辅助视角

### 4.1 技术路线视角：Technology Generation（Gen）

Gen 说明交易所使用什么证明技术。它帮助解释 Stage，但**不直接决定 Stage**。


| Gen       | 名称                | 最低要求                               | 代表                          |
| --------- | ----------------- | ---------------------------------- | --------------------------- |
| **Gen 0** | No PoR            | 无持续 PoR，或只有传统财报/营销声明               | Coinbase（交易所级）、KuCoin（停止）   |
| **Gen 1** | Merkle inclusion  | 用户可验证自身 inclusion；无 ZK 偿付约束        | Bitget, MEXC, Bybit, Kraken |
| **Gen 2** | Merkle + ZK       | 有全局 ZK / zk-SNARK / zk-STARK 偿付证明  | Binance, OKX, Gate.io, HTX  |
| **Gen 3** | ZK + onchain + DA | 公开 proof 可链上验证，root/proof 有不可篡改发布层 | 行业尚无                        |


**Gen 如何解释 Stage**：

- Gen 0 通常对应 Pre-Stage，不进入 PoR Stage。
- Gen 1 通常只能到 Stage 0，除非同时具备公开 wallet_address_list、wallet_ownership_proof、global_proof、审查和交叉验证机制。
- Gen 2 说明负债侧可能有更强证明；只有在 wallet_address_list、wallet_ownership_proof、global_proof 与可信设置 / 参数来源均公开可验证时，才具备 Stage 1 基础。
- Gen 3 接近 Stage 2 的技术形态，但仍需满足双侧公开、低门槛 inclusion proof、参数透明、DA、permissionless、频率等门槛。

**Gen 为什么不直接决定 Stage**：

- ZK proof 可以很先进，但如果 proof 包不公开、可信设置不可验证或 root 不锚定，用户仍要信交易所分发层或 setup 诚实性。
- Merkle 技术较基础，但若资产侧、wallet_ownership_proof、审查、历史记录更公开，Stage 可能高于某些“链下 ZK 展示”。
- 第三方审查不改变 Gen，它影响 Stage 门槛和置信度。

### 4.2 证据公开视角：Evidence Level（E）

E 说明外部观察者能取得和复现多少证据。它是 Stage 判定的证据输入，**不是最终透明度评级**。E 主要描述“证据公开到哪里”，同时观察用户侧验证路径是否可用、低门槛；但它不自动判断 ownership、可信设置、业务一致约束或发布频率是否满足 Stage 门槛。


| E      | 名称                               | 一句话定义                                               | 常见 Stage 解释                  |
| ------ | -------------------------------- | --------------------------------------------------- | ---------------------------- |
| **E0** | No usable evidence               | 没有可用 PoR 工件，或 PoR 已停止                               | Pre-Stage                    |
| **E1** | Limited disclosure               | 有 summary、储备率、用户 inclusion 或基础 verifier，但无法复现全局；验证路径可能仍偏复杂 | Stage 0                      |
| **E2** | Public verifiable artifacts      | 有免许可公开 artifacts，可复验资产侧、负债侧或其关键部分                   | Stage 0 或 Stage 1，取决于硬门槛是否齐备 |
| **E3** | Anchored permissionless evidence | E2 + 低门槛 inclusion proof + 链上锚定 / DA / permissionless verification / 高频更新 | Stage 2 候选                   |


**E 如何支持 Stage 判定**：

- E0 对应 Pre-Stage，不进入 PoR Stage。
- E1 说明交易所有最低 PoR 披露，但通常不足以越过 Stage 0；若 inclusion proof 验证流程复杂、只能依赖 CLI 或人工拼接参数，应降低用户可验证性评分。
- E2 说明已有公开可复验 artifacts，但仍需检查 wallet_address_list、wallet_ownership_proof、global_proof、参数来源、业务一致约束和审查独立性，才能判断是否达到 Stage 1。
- E3 是 Stage 2 的证据形态，但仍需同时满足 Stage 2 的低门槛用户验证、业务一致约束、permissionless、DA、参数与版本锚定、频率和独立性要求。

**E 级详细判定要点**：

- **E0**：无 snapshot、summary、root、proof、wallet_address_list、verifier，或 PoR 已停止。
- **E1**：有 snapshot、summary、储备率、用户 inclusion proof 或基础 verifier；但第三方无法复现全局资产与负债关系。若用户验证需要复杂命令、环境配置或手动处理 proof，应标记为 high-friction E1。
- **E2**：wallet_address_list、wallet_ownership_proof、global_proof、root/vk/config、proof schema、trusted setup transcript 等核心 artifacts 中已有公开可复验部分；若缺关键项，仍可能只是 Stage 0。
- **E3**：E2 基础上，inclusion proof 有低门槛用户验证路径，root/proof/commitment 有 canonical 锚定，artifacts 有稳定 DA，第三方可 permissionless 重放验证，并具备更高频或事件触发更新。

### 4.3 Gen / E 到 Stage 的解释关系


| 组合           | Stage 解释                                              |
| ------------ | ----------------------------------------------------- |
| Gen 0 / E0   | Pre-Stage：没有最低可用 PoR                                  |
| Gen 1 + E1   | 通常 Stage 0：有 summary 或用户自验，但无法复现全局偿付                  |
| Gen 2 + E1   | 通常 Stage 0：技术上可能有 ZK，但生产 proof 受限公开或关键 artifacts 不完整  |
| Gen 2 + E2   | Stage 0 或 Stage 1：取决于 wallet_ownership_proof、global_proof、参数来源是否齐备 |
| Gen 2/3 + E3 | Stage 2 候选：仍需检查低门槛 inclusion proof、业务一致约束、permissionless、DA、频率、参数与版本锚定 |


**反例说明**：

- Binance：Gen 2 + E2，wallet_address_list 与 global_proof 均公开，但缺少 wallet_ownership_proof 和公开可信设置 transcript，因此不进入 Stage 1。
- OKX：Gen 2 + E2，wallet_address_list、wallet_ownership_proof 与 global_proof 均公开，因此可进入 Stage 1；但 proof 仍链下分发，因此不是 Stage 2。
- Bybit：Gen 1 + E1，有用户 Merkle 和第三方机构核验，但无公开 wallet_address_list / ZK，因此仍是 Stage 0。

---

## 5. Stage 内评分

Stage 是门槛制；同一 Stage 内可用 100 分做横向排序：


| 类别         | 权重  | 评分提示                                                                |
| ---------- | --- | ------------------------------------------------------------------- |
| 公开数据可得性    | 15  | summary、history、wallet_address_list、global_proof、URL 稳定性、机器可读性          |
| 负债侧证明强度与口径 | 20  | Merkle / Sum Tree / ZK / per-user solvency / 负净值虚假用户防御 / 表外负债边界     |
| 储备侧可验证性与质量 | 20  | wallet_address_list、块高、onchain 可复核、clean assets、自发行代币、质押/托管分类       |
| 抗篡改、时间锚与频率 | 15  | 链上 root、第三方镜像、DA、历史 commitment、周度/日度/事件触发披露                         |
| 审查独立性与保证等级 | 10  | 技术核验、AUP、有限保证、合理保证、审查方资质与利益冲突                                       |
| 用户可验证性     | 10  | 用户 proof、导出能力、Web/WASM/GUI 一键验证、开源 verifier、清晰错误提示、真实参与门槛与验证率 |
| 证明系统与实现风险  | 10  | trusted setup、透明设置、hash truncation、证明约束充分性、参数一致性、代码维护、测试、审计发现修复状态 |


**地板分规则**（Stage 不变，但分数上限受限）：


| 条件                                | 总分上限 |
| --------------------------------- | ---- |
| 无用户负债侧数据                          | 40   |
| 无资产侧公开数据                          | 65   |
| 无 wallet_address_list 且无 ZK        | 55   |
| hash 输出 <128 bit                  | 45   |
| 无近期 PoR                           | 20   |
| proof/root/vk 仅 server 托管且无归档/锚定  | 80   |
| 无 wallet_ownership_proof          | 60   |
| trusted setup 不透明或无可验证 transcript | 60   |
| 重大未修复实现缺陷影响 proof 语义              | 70   |
| 审查方与技术提供方存在未披露利益冲突                | 75   |
| 有借贷/保证金业务但 proof 未约束风险控制逻辑        | 70   |
| 不能证明无负净值虚假用户或负余额抵消真实负债          | 70   |


---

## 6. 交易所映射


| Exchange       | **PoR Stage** | Gen   | E Level | Stage 结论理由                                                                            |
| -------------- | ------------- | ----- | ------- | ------------------------------------------------------------------------------------- |
| Coinbase（交易所级） | Pre-Stage     | Gen 0 | E0      | 有上市公司财务披露，但未提供交易所级 PoR 工件                                                             |
| KuCoin         | Pre-Stage     | Gen 0 | E0      | PoR 已停止                                                                               |
| Binance        | Stage 0       | Gen 2 | E2      | wallet_address_list 与 global_proof 公开；但无 wallet_ownership_proof、无公开可信设置 transcript，Stage 1 blocked |
| OKX            | **Stage 1**   | Gen 2 | E2      | wallet_address_list、wallet_ownership_proof 与 global_proof 公开；proof 链下分发                  |
| Gate.io        | Stage 0       | Gen 2 | E1      | summary 公开；ZK 包登录；无 wallet_address_list                                                 |
| HTX            | Stage 0       | Gen 2 | E1      | 样例 ZK 包可验；生产数据和 wallet_address_list 缺口明显                                                |
| Bybit          | Stage 0       | Gen 1 | E1      | 用户 Merkle + 第三方机构核验；无 wallet_address_list / ZK                                          |
| Bitget         | Stage 0       | Gen 1 | weak E1 | 64-bit 截断 hash；无独立 attestation                                                        |


---

## 7. 实操评定清单

评定某交易所 Stage 前，按顺序回答：

### 7.1 Pre-Stage checklist

满足任一即标记为 `Pre-Stage — No usable PoR`：

1. 近 12 个月是否无公开 PoR？
2. PoR 是否已经停止？
3. 是否只有营销声明，没有 snapshot、summary、root、proof 或审查报告？
4. 是否只有传统财报，尚未提供可映射用户资产负债侧的 PoR 工件？
5. 是否只有链上余额展示，但无用户负债侧 summary / inclusion / proof？

### 7.2 Stage 0 checklist

1. 近 12 个月是否持续发布 PoR？
2. 是否有明确 snapshot 时间、覆盖资产？
3. 是否有 summary 或用户 inclusion proof？
4. 若提供用户 proof：规则是否公开、proof 是否可导出、第三方能否本地复验？
5. inclusion proof 是否提供低门槛验证路径（例如 Web/WASM/GUI 一键验证），避免普通用户必须执行复杂 CLI、配置环境或手动拼接参数？
6. 是否有开源 verifier 或足够算法说明？

### 7.3 Stage 1 checklist

1. 资产侧 wallet_address_list 是否公开、无需登录、可链上复核？
2. 是否公开 wallet_ownership_proof，且可批量验证？
3. 负债侧 global_proof 是否免登录、免注册、免 cookie、免 API key、免人工申请即可下载并独立复验？
4. 是否有独立第三方审查？类型是否明示（技术核验 / AUP / 有限保证 / 合理保证）？
5. summary ↔ wallet_address_list aggregate 与 summary ↔ proof root / commitment 是否都可 bind？
6. 若使用可信设置，是否公开可验证 transcript / MPC ceremony？若不使用可信设置，证明系统是否透明设置？
7. 历史 snapshot artifacts 是否可归档？
8. 用户侧 inclusion proof 的验证体验是否足够低门槛，且与公开 global_proof / root / summary 绑定？
9. proof 边界是否清楚（覆盖哪些负债、不覆盖哪些表外项目）？

### 7.4 Stage 2 checklist

1. 是否已满足 Stage 1 的 wallet_address_list + wallet_ownership_proof + global_proof + 参数透明门槛？
2. 证明约束是否与交易所实际业务一致？
3. 是否能自证没有插入负净值虚假用户、没有用负余额抵消真实负债、用户余额承诺与全局负债汇总一致？
4. 若有借贷、保证金、组合保证金、抵押品折扣、负余额、价格参数或风险限额，是否在电路或等价公开证明中被约束？
5. inclusion proof 是否可导出、可本地复验，并提供 Web/WASM/GUI 一键验证、清晰错误提示，且与公开 global_proof / root / summary 绑定？
6. root / proof / commitment 是否链上锚定，且有历史 commitment 链？
7. proof 包、vk、config 是否在稳定 DA 层长期可得？
8. 第三方是否 permissionless 可重跑核心验证（无需登录/API key）？
9. vk/config/verifier 版本是否被链上/DA 固定，变更是否有 audit trail？
10. 发布频率是否足够（周度完整 PoR / 日度锚定 / 事件触发增发）？
11. 审查方与技术提供方利益冲突是否披露？

### 7.5 通用风险项

1. 储备质量是否分级（clean / 质押 / 链外 / 自发行代币）？
2. 是否有重大未修复实现缺陷？
3. 报告是否清楚说明不能证明什么？

---

## 8. 对监管者的建议

PoR Stage 可以作为监管者设定交易所储备透明度要求的分层语言，但它与传统财务审计属于不同维度。监管文本应明确：PoR 证明的是某一快照下、某一披露口径内的资产与负债关系；传统审计则覆盖财务报表、内部控制、公司治理、链外资产和持续经营等公司层面问题。二者应互补使用。

### 8.1 最低有效 PoR 要求

监管者不应接受 `Pre-Stage` 或 **Stage 0** 作为有效 PoR。若交易所面向公众托管客户资产，有效 PoR 的最低要求应至少达到 **Stage 1**；Stage 0 只能作为过渡期披露、整改观察项或风险提示，不能作为满足监管 PoR 要求的依据。


| 要求        | 监管处理建议                                                     |
| --------- | ---------------------------------------------------------- |
| Pre-Stage | 不应被表述为“已有 PoR”；只能标记为无最低可用 PoR                              |
| Stage 0   | 不应被认定为有效 PoR；只能作为过渡披露或整改状态，且必须显式标记仍主要依赖交易所自述               |
| Stage 1   | 可作为有效 PoR 的最低监管标准：wallet_address_list、wallet_ownership_proof、global_proof、参数来源均公开可验证 |
| Stage 2   | 可作为长期最佳实践：低门槛 inclusion proof、链上锚定、DA、permissionless verification、持续可追溯 |


### 8.2 不应混淆的几类报告

监管者应要求交易所在对外披露中明确区分：


| 类型            | 含义                                       |
| ------------- | ---------------------------------------- |
| 技术核验          | 检查 Merkle/ZK/wallet_address_list/wallet_ownership_proof/global_proof 等技术工件是否可复验 |
| AUP / 商定程序    | 审查方执行约定检查步骤，但通常不发表整体保证意见                 |
| 有限保证          | 对特定范围提供有限程度保证                            |
| 合理保证 / 完整财务审计 | 更接近传统审计意见，但仍需说明是否覆盖 PoR 负债口径             |


“有审计”和“有 PoR”回答的问题不同。传统财务审计提供公司层面的保证；PoR 提供面向用户和第三方的公开可复验证据。监管与用户披露中应同时说明两者各自覆盖的范围。

### 8.3 传统审计与 PoR 的差异和互补

传统审计与 PoR 回答的问题不同，监管要求中应同时保留二者，并说明各自边界。


| 维度   | 传统审计                      | PoR                                                    |
| ---- | ------------------------- | ------------------------------------------------------ |
| 核心问题 | 公司财务报表是否在重大方面公允列报         | 某一快照下用户负债是否被纳入证明，储备是否足以覆盖披露口径                          |
| 时间语义 | 通常是季度、年度或特定期间             | 通常是某一 snapshot，也可扩展到高频或连续证明                            |
| 资产验证 | 覆盖现金、银行存款、应收款、链外资产、投资等    | 更擅长验证链上资产、wallet_ownership_proof、链上余额                   |
| 负债验证 | 按会计口径确认负债、或通过抽样/程序检查      | 可对全量用户余额承诺、Merkle/ZK inclusion、global_proof 做可复验承诺     |
| 验证主体 | 审计师、监管者、投资者               | 用户、第三方研究者、监管者、自动化 verifier                             |
| 输出形式 | 审计意见、AUP、保证报告、财务报表附注      | root、proof、wallet_address_list、wallet_ownership_proof、vk/config、verifier 输出 |
| 主要盲区 | 不一定提供用户级可复验证据；不一定覆盖实时托管风险 | 不能自动证明链外资产、表外负债、内部控制或经营能力                              |


二者的互补方式：

1. **传统审计覆盖链外与公司层面风险**：银行存款、应收款、关联方交易、表外安排、内部控制、持续经营能力等，属于传统审计更适合处理的范围。
2. **PoR 覆盖用户级和链上可复验风险**：用户 inclusion、global_proof、wallet_ownership_proof、链上余额、proof/root/commitment 的公开复验，是传统审计通常不给用户直接验证的部分。
3. **会计口径与 PoR 口径必须对齐**：监管者应要求披露“PoR 覆盖哪些负债”和“财务报表确认哪些负债”之间的差异，尤其是借贷、保证金、衍生品、关联方、表外项目。
4. **审计师不应只验证交易所生成的结论**：如果审查范围只覆盖交易所提供的 summary，而不验证 wallet_address_list、wallet_ownership_proof、global_proof、参数来源和历史 artifacts，则不能被表述为 PoR Stage 1 级别审查。
5. **监管应采用双轨要求**：传统审计负责财务完整性和链外风险，PoR Stage 负责公开可复验的储备透明度。二者同时满足，才更接近用户关心的“可兑付性”问题。

### 8.4 监管应强制披露的字段

监管报告至少应要求机器可读披露：

1. snapshot 时间、覆盖资产、覆盖链、覆盖用户负债口径。
2. reserve summary、liability summary、reserve ratio。
3. wallet_address_list、块高、余额、网络、资产分类。
4. wallet_ownership_proof 与验证方法。
5. global_proof、root、vk/config、proof schema。
6. 证明约束与业务模型及关键安全性质的对应关系：是否能自证无负净值虚假用户、无负余额抵消真实负债；借贷、保证金、组合保证金、抵押品折扣、价格参数和风险限额是否被约束。
7. trusted setup transcript / MPC ceremony，或透明设置证明系统说明。
8. 第三方审查机构、审查类型、范围限制、利益冲突声明。
9. proof/root/commitment 的归档位置、哈希、历史版本和锚定记录。
10. 发布频率、事件触发增发规则、历史快照保留策略。
11. 明确说明 PoR 不覆盖的资产、负债和表外项目。

### 8.5 频率与事件触发

对大型 CEX，月度 PoR 只能作为过渡期基线。监管者应考虑按风险分层提出更高频要求：


| 场景               | 建议                                        |
| ---------------- | ----------------------------------------- |
| 普通持续披露           | 至少月度完整 PoR，历史 artifacts 可归档               |
| 高体量或高杠杆平台        | 周度完整 PoR，或至少周度负债侧更新                       |
| 市场极端波动、挤兑、重大安全事件 | 事件触发增发 PoR                                |
| Stage 2 目标       | 低门槛 inclusion proof + 日度 root/commitment 锚定，逐步走向日度完整 PoR 或准实时证明 |


### 8.6 执法与用户呈现

监管者应要求交易所以统一格式展示：

```text
PoR Stage: Stage 1
Technology Generation: Gen 2
Evidence Level: E2
Scope: spot + margin liabilities, excludes off-balance-sheet exposures
Snapshot: 2026-06-01T00:00:00Z
```

若交易所缺少 wallet_address_list、wallet_ownership_proof、可信设置 transcript、公开 global_proof、关键证明约束说明或历史 artifacts，应在用户界面中显式标记 `UNVERIFIABLE`，不能用“100% backed”“audited”或“ZK verified”等营销语句覆盖关键缺口。

---

## 9. 对行业的建议

PoR 的长期目标不是让交易所发布更多 PDF 或营销页面，而是让用户、第三方研究者、审查机构和监管者都能复验同一组核心事实：交易所控制哪些资产、承诺了哪些负债、证明约束是否覆盖真实业务风险，以及这些证据是否能被长期追溯。

### 9.1 对交易所的建议

交易所应把 **Stage 1** 作为近期最低目标，把 **Stage 2** 作为长期建设方向：

1. **不要停留在 Stage 0**：summary、储备率或用户 Merkle proof 只能说明有最低披露，不能证明全局偿付关系。
2. **公开双侧核心 artifacts**：wallet_address_list、wallet_ownership_proof、global_proof、root/vk/config、proof schema 和 trusted setup transcript 应无需登录、注册、cookie、API key 或人工申请即可取得。
3. **披露 proof 覆盖边界**：明确说明 PoR 覆盖哪些资产、负债、产品线和风险场景，不覆盖哪些链外资产、表外负债、借贷、保证金、衍生品或关联方项目。
4. **让证明约束贴合业务模型**：PoR 不应只证明静态余额，还应自证没有插入负净值虚假用户、没有用负余额抵消真实负债；如果交易所有借贷、保证金、组合保证金、抵押品折扣或价格参数，还应说明这些风险控制逻辑如何被约束或为何未覆盖。
5. **降低用户验证门槛**：inclusion proof 不应只提供给开发者运行 CLI；交易所应提供 Web/WASM/GUI 一键验证、清晰错误提示、proof 导出和本地复验路径，并披露验证入口与参与率指标。
6. **优先机器可读和可归档**：JSON/CSV/ZIP/schema/API、开源 verifier 和历史 artifacts 应优先于不可解析 PDF；每期 snapshot 应有稳定 URL、哈希和版本记录。
7. **提高发布频率**：月度 PoR 只能作为过渡基线；高体量或高杠杆平台应逐步走向周度、日度锚定或事件触发增发。

### 9.2 对第三方审查机构的建议

第三方机构不应只复述交易所生成的结论，而应明确自己验证了哪些 artifacts、没有验证哪些边界：

1. 区分技术核验、AUP / 商定程序、有限保证和合理保证，不用“audited”模糊覆盖范围差异。
2. 至少检查 wallet_address_list、wallet_ownership_proof、global_proof、参数来源、proof schema、证明约束、verifier 输出和历史 artifacts 是否一致。
3. 披露审查范围、抽样方法、限制条件、利益冲突和未解决问题。
4. 若只验证 summary 或用户 inclusion，不能将结论表述为 Stage 1 级别 PoR。

### 9.3 对工具和基础设施建设者的建议

行业需要共享的验证工具和数据基础设施，而不是每家交易所各自定义不可复用格式：

1. 建立标准化 PoR schema，覆盖 snapshot、wallet_address_list、wallet_ownership_proof、global_proof、vk/config、setup transcript、proof constraints、scope 和 exclusion。
2. 建设开源 verifier、批量签名验证工具、proof 重放工具、Web/WASM 一键验证组件和可复现测试集。
3. 提供独立归档与镜像服务，使 proof/root/commitment、wallet_address_list 和审查报告不依赖交易所单一 server。
4. 推动链上锚定、DA、IPFS/Arweave/对象存储哈希索引等长期可追溯发布层。
5. 对外展示时应同时显示 Stage、Gen、E Level、scope、snapshot 和关键缺口，而不是只展示储备率。

### 9.4 行业迁移路径

行业可按以下路径推进：


| 阶段        | 目标                   | 关键动作                                        |
| --------- | -------------------- | ------------------------------------------- |
| 从无到有      | Pre-Stage -> Stage 0 | 建立持续 snapshot、summary、用户 proof 和基础 verifier |
| 从自述到可复验   | Stage 0 -> Stage 1   | 公开 wallet_address_list、wallet_ownership_proof、global_proof、参数来源和审查范围 |
| 从披露到信任最小化 | Stage 1 -> Stage 2   | 引入低门槛 inclusion proof、canonical 锚定、DA、permissionless 验证、关键安全性质自证和业务一致约束 |
| 从单点到生态协作  | 单交易所披露 -> 行业公共验证基础设施 | 标准 schema、开源 verifier、独立归档、跨机构复验            |



# ardmere 服务对象分析

`ardmere` 的产品定位不应只是“给普通用户看的 PoR 验证工具”，而应是面向监管者、机构资金方、独立审查者和交易所的 PoR 透明度基础设施。

一句话定位：

> ardmere 帮助外部观察者判断交易所 PoR 是否真实可验证、是否达到有效 PoR 标准，以及缺少哪些关键证据。

---

## 1. 核心服务对象

### 1.1 监管者 / 政策制定者

监管者是最重要的服务对象之一。

他们需要的不是亲自运行 proof，而是回答：

1. 哪些交易所有有效 PoR？
2. 哪些只是 Stage 0，不能被当作有效 PoR？
3. 交易所缺少哪些关键 artifacts？
4. 是否公开 `wallet_address_list`、`wallet_ownership_proof`、`global_proof`？
5. 证明系统是否依赖不透明可信设置？
6. 是否存在负净值虚假用户、负余额抵消真实负债、链下分发、历史可替换等风险？
7. 如何把 PoR 要求写进监管规则或牌照条件？

ardmere 对监管者的价值：

- 提供可执行的 PoR 最低标准。
- 区分 `Pre-Stage`、`Stage 0`、`Stage 1`、`Stage 2`。
- 明确指出 Stage 0 不应被接受为有效 PoR。
- 将 PoR 与传统审计的边界拆清楚，帮助形成双轨监管要求。

### 1.2 机构投资者 / 做市商 / 托管客户

机构资金方关心的是交易所透明度和托管风险。

他们需要：

1. 横向比较不同交易所 PoR 透明度。
2. 识别“100% backed”“audited”“ZK verified”等表述是否过度营销。
3. 判断某个交易所 PoR 是 Stage 0、Stage 1 还是 Stage 2 候选。
4. 将 PoR 风险信号纳入交易所准入、额度管理、做市风险和托管风险模型。
5. 追踪交易所 PoR 是否持续发布、是否历史可归档、是否发生退步。

ardmere 对机构的价值：

- 提供交易所 PoR risk rating / transparency intelligence。
- 提供可追溯的证据链接、缺口清单和风险标记。
- 将技术验证结果转化成机构风控可读的评级和报告。

### 1.3 独立审查机构 / 安全公司 / 会计师事务所

第三方审查机构可以把 ardmere 作为 PoR 审查基础设施。

他们需要：

1. 自动检查 PoR artifacts 是否完整。
2. 重跑 proof / verifier。
3. 检查 `wallet_ownership_proof`、`global_proof`、`vk/config`、trusted setup transcript。
4. 生成 AUP / 技术核验 / 有限保证所需的 checklist。
5. 标记无法验证的范围和证据缺口。

ardmere 对审查机构的价值：

- 成为 PoR technical assessment 的 tooling layer。
- 降低审查机构重复构建 verifier、抓取、归档、对账工具的成本。
- 帮助审查机构避免只复述交易所 summary，而忽略关键 artifacts。

### 1.4 交易所

交易所也是潜在服务对象，但需要特别注意独立性。

他们需要：

1. 知道自己距离 Stage 1 / Stage 2 还差什么。
2. 预审 PoR artifacts 是否完整。
3. 检查 proof、verifier、披露格式、历史归档和用户验证体验是否达标。
4. 生成面向监管、用户和机构客户的透明度改进路线图。

ardmere 对交易所的价值：

- 提供 PoR readiness assessment。
- 提供 Stage upgrade gap analysis。
- 提供 artifacts schema、verifier 输出和 disclosure checklist。

但如果服务交易所，应明确区分：

- independent public rating
- paid technical assessment
- remediation consulting

避免出现“自己审自己”的信任冲突。

### 1.5 高阶用户 / 研究者 / 媒体

这类用户通常不会每天运行 proof，但会引用评级和报告。

他们需要：

1. 简洁可信的交易所透明度排名。
2. 可追溯的证据链接。
3. 对交易所营销话术的事实核查。
4. 对 PoR 技术差异、审查范围和风险缺口的可读解释。

ardmere 对他们的价值：

- 提供公共透明度报告。
- 提供交易所 PoR 证据索引。
- 提供可引用的 Stage 结论和风险解释。

---

## 2. 谁是第一阶段直接客户

普通用户是最终受益者，但不一定是第一阶段的直接客户。

第一阶段最适合优先服务：

1. **监管者和政策研究者**：需要定义有效 PoR 的最低标准。
2. **机构资金方**：需要把 PoR 透明度纳入风控。
3. **独立审查者**：需要工具化、标准化 PoR 审查流程。

交易所可以作为第二阶段服务对象，但必须建立独立性边界。

---

## 3. 产品形态

### 3.1 Public Transparency Dashboard

面向公众、研究者、媒体和机构初筛。

核心内容：

- 交易所 PoR Stage。
- Technology Generation（Gen）。
- Evidence Level（E）。
- 缺失 artifacts。
- 风险标记。
- 历史快照。
- 最后验证时间。

### 3.2 Regulator / Institution Report

面向监管者和机构客户。

核心内容：

- 每家交易所的 Stage 结论。
- 是否达到有效 PoR 最低标准。
- Stage 0 / Stage 1 / Stage 2 blocked reasons。
- proof、wallet、setup、DA、frequency、audit taxonomy 的详细检查结果。
- 与传统审计的互补关系说明。

### 3.3 Verification API

面向审查机构、机构客户和内部系统。

核心能力：

- 输入 PoR artifacts。
- 输出 verifier result。
- 输出 missing artifacts。
- 输出 risk flags。
- 输出 suggested Stage / E Level。
- 输出 machine-readable report。

### 3.4 Exchange Readiness Assessment

面向交易所，但需保持独立性。

核心内容：

- 当前 Stage。
- 升级到 Stage 1 / Stage 2 的缺口。
- artifacts schema 建议。
- inclusion proof 用户验证体验建议。
- trusted setup、wallet ownership、DA、anchoring、frequency 改进路线。

---

## 4. 独立性原则

ardmere 如果要长期建立公信力，需要明确几条边界：

1. 公共评级不能由交易所付费购买。
2. 付费技术评估必须披露范围和利益关系。
3. remediation consulting 不能直接等同于独立认证。
4. 所有评级必须绑定 artifacts、hash、URL、验证输出和时间戳。
5. 缺数据必须标记为 `UNVERIFIABLE`，不能默认通过。

---

## 5. 推荐推进顺序

### Phase 1: Public Methodology + First Reports

目标：建立可信方法论和样例报告。

建议对象：

- OKX：Stage 1 代表样本。
- Binance：Gen 2 / E2 但 Stage 1 blocked 的代表样本。
- Bybit 或 Bitget：Stage 0 / Merkle inclusion 型样本。

输出：

- public dashboard v0
- exchange transparency report v0
- methodology page

### Phase 2: Institution / Regulator Package

目标：把评级转化为监管和机构风控语言。

输出：

- 交易所 PoR 风险矩阵。
- 有效 PoR 最低标准说明。
- Stage 0 不应作为有效 PoR 的政策建议。
- 传统审计与 PoR 互补说明。

### Phase 3: API and Tooling

目标：让验证服务产品化。

输出：

- artifacts upload / fetch API
- verification API
- report API
- rule engine
- historical archive

### Phase 4: Exchange Readiness / Remediation

目标：帮助交易所改进 PoR，但不牺牲独立性。

输出：

- readiness assessment
- Stage upgrade checklist
- artifacts schema
- user verification UX checklist

---

## 6. 结论

ardmere 的直接客户不应首先定义为普通散户，而应定义为：

> 监管者、机构资金方和独立审查者。

普通用户是最终受益者；交易所是潜在被评估对象和后续服务对象。

最稳健的产品路线是：

1. 先建立公开、独立、可追溯的 PoR Stage 评级。
2. 再服务监管者和机构客户。
3. 最后在明确独立性边界的前提下服务交易所改进。

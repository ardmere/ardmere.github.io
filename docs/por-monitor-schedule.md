# Ardmere PoR 监控日程（按所窗口）

> **生效**：2026-11-06（2026-11-07 扩展覆盖；2026-11-08 起扫描 PR 自动 merge）
> **机器配置**：[`config/por-monitor/schedule.yaml`](../config/por-monitor/schedule.yaml)
> **Due 列表**：[`scripts/por-monitor-due.sh`](../scripts/por-monitor-due.sh)
> **回归测试**：[`scripts/por-monitor-due_test.sh`](../scripts/por-monitor-due_test.sh)
> **扩所短名单**：[`exchange-onboarding-shortlist.md`](./exchange-onboarding-shortlist.md)

## 覆盖规则（2026-11-07 起）

**凡是有已合并适配器 + 已提交 `docs/reports/<id>/` 评估的所，必须在 `schedule.yaml` 有 `exchanges:` 条目。**
2026-11-07 之前十个所（ourbit, coinex, bitunix, woox, toobit, phemex, coinw, btcc, kraken, poloniex）没有条目，
`por-monitor-due.sh` 因此从不为它们排期 —— 这些所对监控完全不可见。

现由 `./scripts/por-monitor-due.sh --check` 守住：扫描 `docs/reports/*/` 中带 `*-assessment.json` 的目录，
与实际 `exchanges:` 条目比对，发现缺口即 **exit 1**。`scripts/por-monitor-due_test.sh` 在回归测试里调用它。

## 政策（2026-11-08 更新：自动 merge）

扫描出的报告 PR 不再等人工确认，由扫描 worker 自己合入。

| | 废止 | 现行 |
| --- | --- | --- |
| 合并 | 扫描 PR 留待人工 review 后合入 | **扫描 worker 自己 squash merge**：本地校验四项全过后直接合入 `main` |
| 阻断 | — | 本地校验失败、且该失败不是 `main` 上的既存失败 → **不合**，留 PR 并报出失败项 |

**PR 检查不是门禁 —— 本仓库没有 PR CI。** 没有 `.github/workflows/`；唯一的 workflow 是 `main` 上的
`pages-build-deployment`，而 `gh pr checks` 对每个 PR 都回 "no checks reported"。所以「所有检查 terminal 且无失败」
在这里是**空条件**，不能当自动合并的判据。

自动合并的真实门禁是**本地验证**，以下四项**必须全过**：

```bash
go build ./...
go test ./...
./scripts/check-skills.sh
./scripts/por-monitor-due_test.sh
```

只有**已知在 `main` 上同样失败、且已写进 PR 正文**的项可以豁免；其它任何失败都阻断合并，并须报出失败的检查名。

**自动 merge 的边界**：只适用于扫描自身的报告 PR（新快照、首页/对比/历史/归档刷新、frequency 修复）。
**代码、合约、schema、skills、schedule 结构**改动仍需人工 review；一个 PR 里同时混入这两类改动时，按人工 review 处理。

## 政策变更（2026-11-06）

| | 旧（≤2026-11-05） | 新（2026-11-06） |
| --- | --- | --- |
| 扫描对象 | 每天 2 次扫 **全部** P0+P1 | 只扫 **当日窗口到期** 的所 |
| 深度 | 日检常带全量探针 | **轻探针**；仅发现新 snapshot id 才 `fetch`→`report`→首页 |
| SLA | 公开更新后 ≤1 自然日 | **不变** |

旧「HK 10:00 / 22:00 全量 P0+P1」偏好已 **废止**（用户请求 2026-11-06）。定时器仍可保留每日一次（及可选傍晚 catch-up），但 prompt 必须先跑 `por-monitor-due.sh`。

## 为何改

`docs/reports/*/ *-assessment.json` 历史显示多数所 **约月更**（Binance/BingX/HTX/KuCoin/MEXC 日分布稳定；OKX/Bitget/Bybit 用「上次 + 中位间隔」）。每日全扫浪费配额，且对 WAF 所反复 403 无增量。

## 历史窗口摘要（评估集）

| 所 | 观察到的节奏 | 扫描窗（UTC，含 buffer） |
| --- | --- | --- |
| Binance | 多在每月 1 日；偶发月末/中途（如 10-31、PR04NOV26） | 每月 1–3 日与 30–31 日 |
| OKX | 间隔中位 ~29 天，日分散 | `last + 29d ±4`；fallback 月初/月末 |
| MEXC | ~每月 10–15 日 | 9–17 日 |
| KuCoin | 月末 30/31 | 29–31 与次月 1 日 |
| BingX | 稳定 15 日 | 14–17 日 |
| HTX | 近月 1 日 | 1–3 日 |
| Gate.io | 有公开数据时 ~16–18 日（稀疏） | 16–20 日 |
| Bybit | 近月偏 23–27 日 | `last + 30d ±4`；fallback 22–29 |
| Bitget | 间隔波动大（~3–5 周） | `last + 31d ±5`；fallback 18–29 |

完整字段与 rationale 见 YAML。

## 新增覆盖（2026-11-07）：十个 summary-tier 所

| 所 | 节奏 | 窗口（UTC） | 置信 | 探针策略 | 阻塞 |
| --- | --- | --- | --- | --- | --- |
| Ourbit | 月锚（单快照，UTC+8 月初） | 每月 1–3 日（−1/+2） | **low** | `stock_info` HTTP | 无 period root；用户证明登录 |
| CoinEx | ~34 天（中位），日分散 | `last + 34d ±5`；≥45 天静默强制补检 | **high** | `merkle/reserve` 历史 HTTP | XLSX 无 ownership |
| Bitunix | 月更，月初 1–6 日 | 每月 1–6 日 | **high** | `queryAuditTimeList` HTTP | 用户证明登录 |
| WOO X | 月更，固定 15 日 | 每月 14–17 日 | **high** | `porl/v2/audits` HTTP | 用户 hash 登录 |
| Toobit | 月更，1 日 | 每月 1–3 日 | **high** | `latest-audit` HTTP | 无公开钱包表 |
| Phemex | 月更，1 日 | 每月 1–3 日 | **high** | `porVersions` HTTP | api10 root/addresses CloudFront 闸 |
| **CoinW** | 月更，月末（两次观察） | 每月 20–31 日（−1/+3），溢出次月 1–3 日 | **low** | **`announcement_browser_capture` + `capture: browser`** | **Cloudflare WAF 403 + 地域限制** |
| BTCC | 月更，固定 15 日 | 每月 14–17 日 | **high** | `reserveInfo?reportDate=` HTTP | root 截断；历史需登录 |
| Kraken | **季度**，季末后 17–48 日发布 | 季末月 28–31 日 + 季末后 15–52 日 | **high** | LedgerLens datasets HTTP | 无钱包 ZIP；CPA 报告非机器可读 |
| Poloniex | 月更，1 日 | 每月 1–3 日 | **high** | `versionInfo` HTTP | 钱包 CSV 无稳定公开 URL |

### 置信度怎么定的

- **high**：有可机器读取的公开期次列表（`queryAuditTimeList` 23 期、`porVersions` 42 期、
  `porl/v2/audits` 51 期、CoinEx 20 期、`versionInfo` 6 期）或逐月确认（BTCC 连续 12 个月 15 日），
  且日分布收敛。
- **medium**：只有少数评估快照，间隔方差大（OKX、Bybit、Bitget）。
- **low**：**历史只有单个快照**，无法从 repo 推出真实节奏：
  - `ourbit` — `docs/reports/ourbit/` 仅 20261101 一份；`stock_info` 不提供期次列表，
    2026-11-07 实探返回的仍是同一 `snapshotDate`（仍只有一个观测）。选最宽的「月初」窗。
  - `coinw` — 仅 202609 一份（且只有月份标签）+ 一份 2026-10-24 浏览器观察（未入评估）。
    两次观察都落在下旬，故取最宽的下旬窗。
  - `gateio` — 原有条目，3 个稀疏快照，Aug/Sep 登录限制。（顺带标注，非本次新增。）

## CoinW 为什么必须 `capture: browser`

`www.coinw.com` 从数据中心返回 Cloudflare `Just a moment...` 403（2026-11-07 复验：PoR 页与
2026-09 公告页都是 403）。交互浏览器可过 WAF，但会撞「Service Restricted」地域弹窗，
且 PoR SSR 面板的 merkle root / 总储备率在会话中**动画跳变**（记录到两组互相冲突的 64-hex）。

因此该条目：

- `capture: browser` —— 不能当成普通 curl 探针；
- `waf:` 写明三项限制；
- `probe: announcement_browser_capture` —— 稳定通道是 help-center 公告，不是 dashboard。

**运维约定**：CoinW 到期日若 HTTP 探针为空，必须写 `UNVERIFIABLE` / `WAF_GATED` 缺口，
**不得**把「探针无响应」记成「无变化」。`--with-reasons` 会把 `capture=browser` 和
`waf=...` 一起打印出来，worker prompt 必须据此走浏览器路径。

## 脚本 schema 扩展（本次最小改动）

原 schema 有三个表达不了的现实模式，均已在 `scripts/por-monitor-due.sh` 内实现并有回归测试：

| 新字段 | 解决的问题 |
| --- | --- |
| `buffer_days_before` / `buffer_days_after`（**per-entry**） | 原来只能全局 1/2，CoinW 需要 −1/+3 才能覆盖溢出次月的下旬窗 |
| `max_silence_days` | `next_expected` 所在完全停更后永不再到期（CoinEx 已静默 54 天）；超期强制补检 |
| `publish_lag_min_days` / `publish_lag_days` | Kraken 是**季度**且发布滞后 17–48 天，原 `quarter_end` 只开到次月 5 日，整个窗口都错过 |
| `probe_endpoint` / `capture` / `waf` / `confidence` / `confidence_note` | 原 schema 无字段描述探针方式与 WAF，运维只能靠口头约定 |

顺带修掉一个**早已存在**的静默失效：`next_expected` 一旦错过一个窗口就永久掉出，
`OKX`/`Bybit`/`Bitget` 在跨过预期日后会永久不再到期（2027 全年 0 次命中）。
现按整周期前滚（`overdue_roll`），并纳入回归测试。

## 日检怎么跑

```bash
# 校验覆盖：任何有评估报告却无条目的所都会让命令失败
./scripts/por-monitor-due.sh --check

# 今日到期所（UTC），带探针方式/置信/WAF
./scripts/por-monitor-due.sh --with-reasons

# 指定日期复现（排障/回溯）
./scripts/por-monitor-due.sh --date 2026-10-15 --with-reasons

# 回归测试
./scripts/por-monitor-due_test.sh

# 仅轻探针这些 id 的「官方 latest snapshot id」
# 若与 docs/reports/<id>/ 最新 assessment 不同 → 当次 deep-run + 开/续 PR
# 若相同 → 记 unchanged，结束
```

`--with-reasons` 每行四列：`id <TAB> priority <TAB> reason <TAB> probe=<token> capture=<http|browser> confidence=<high|medium|low>`，
有 WAF 时追加 `waf=<note>`。

### Worker prompt 要点

1. `./scripts/por-monitor-due.sh` → 得到 due 列表（空则退出，写「无窗口」笔记即可）。
2. **只**对 due 所做公开探针（不要扫全集）；先读第 4 列的 `capture`/`waf`。
3. 新 id：`fetch` / 浏览器补抓 → `por report` → `por site` → `build-public-html` → PR。
4. **合并步**：先在本地跑齐四项校验 —— `go build ./...`、`go test ./...`、`./scripts/check-skills.sh`、
   `scripts/por-monitor-due_test.sh`。**四项全过 → 自己 squash merge 到 `main`**，再向 Project 汇报合入的 commit 与快照 id。
   **PR 检查不是门禁**（本仓库没有 PR CI，`gh pr checks` 永远 "no checks reported"），别拿它当判据。
   某项失败但**在 `main` 上同样失败**：写进 PR 正文说明后可合；其它失败 **不合**，留 PR 并报出失败的检查名。
5. WAF/登录：同日标注缺口；需要时派 self-hosted 浏览器；**空响应 ≠ 无变化**。
6. **不要**在无新 id 时重跑 OKX `verify-global` 或全量 RPC。

### 定时（cursor-subscriptions）

建议：

| 名 | Cron (UTC) | 作用 |
| --- | --- | --- |
| `por-scan-due` | `0 2 * * *` | 主路径：due 窗轻探针 |
| `por-scan-due-catchup` | `0 14 * * *` | 仅复检 **今日仍在窗内** 或上午已标 pending publish 的所 |

废止：无条件的 `por-scan-am` / `por-scan-pm` 全量 P0+P1。

timer 的 `prompt` 沿用上面「Worker prompt 要点」（含第 4 步的本地校验 + 自行 squash merge），不要退回「等人工 review」。

## 监控清单（能力，非每日必扫）

| 优先级 | 交易所 | 探针 | Cloud | 浏览器/登录 |
| --- | --- | --- | --- | --- |
| P0 | Binance | BAPI / walletZip 列表 | 通常可 | ZK/ownership 仍缺公开 |
| P0 | OKX | PoR API | 通常可 | ZK：`verify-global` 需大内存/swap |
| P0 | MEXC | 新闻/blog | 部分 WAF | 用户证明登录 |
| P1 | KuCoin | SSR `initialReserve` | 通常可 | 无 OSS Merkle |
| P1 | BingX | Support 文章 | 文章可；API 常闸 | Merkle 登录 |
| P1 | HTX | merkle SPA XHR | 裸 curl 弱 | 浏览器 summary；zip 登录 |
| P1 | Gate.io | lastSnapshot | 常 403 | 浏览器；ZK 登录 |
| P1 | Bybit | x-api / PoR | 常 403 | 浏览器；Merkle 登录 |
| P1 | Bitget | queryTotalReserveAmount | 弱 | 浏览器 scrape |
| P2 | Ourbit | `stock_info` | 可 | 用户证明登录；无 period root |
| P2 | CoinEx | `merkle/reserve` | 可 | XLSX 无 ownership |
| P2 | Bitunix | `queryAuditTimeList` | 可 | 用户证明登录 |
| P2 | WOO X | `porl/v2/audits` | 可 | 用户 hash 登录 |
| P2 | Toobit | `latest-audit` | 可 | 无公开钱包表 |
| P2 | Phemex | `porVersions` | 可 | api10 闸；用户证明登录 |
| P2 | CoinW | 公告 + 浏览器 | **403** | **WAF + 地域限制 + 面板不稳** |
| P2 | BTCC | `reserveInfo` | 可 | root 截断；历史登录 |
| P2 | Kraken | LedgerLens | 可 | 无钱包 ZIP；季度低频 |
| P2 | Poloniex | `versionInfo` | 可 | 钱包 CSV 无公开 URL |

## SLA

**交易所公开更新后 ≤ 1 个自然日**内：更新 transparency + assessment、刷新首页与 comparison/history/archive、开 PR，并在本地校验通过后**自行 squash merge**（见「政策（2026-11-08 更新：自动 merge）」）。抓不到则同日写清 UNVERIFIABLE / WAF / 登录缺口。

窗口设计已加 1–3 日 buffer，使「更新日落在窗内 → 次日 catch-up 仍能扫到」，以支撑该 SLA；**不是**靠每日全扫。

## 不在范围内

- 登录态 cookie 常驻爬虫
- 无 signer 时自动上链
- 把 gitignored `artifacts/` 整包推到 Pages
- 短名单所的完整适配器（`shortlist_preview` 现为空：原 poloniex / btcc / kraken 均已有适配器与真实条目）

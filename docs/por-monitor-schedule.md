# Ardmere PoR 监控日程（按所窗口）

> **生效**：2026-11-06  
> **机器配置**：[`config/por-monitor/schedule.yaml`](../config/por-monitor/schedule.yaml)  
> **Due 列表**：[`scripts/por-monitor-due.sh`](../scripts/por-monitor-due.sh)  
> **扩所短名单**：[`exchange-onboarding-shortlist.md`](./exchange-onboarding-shortlist.md)

## 政策变更

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

## 日检怎么跑

```bash
# 今日到期所（UTC）
./scripts/por-monitor-due.sh --with-reasons

# 仅轻探针这些 id 的「官方 latest snapshot id」
# 若与 docs/reports/<id>/ 最新 assessment 不同 → 当次 deep-run + 开/续 PR
# 若相同 → 记 unchanged，结束
```

### Worker prompt 要点

1. `./scripts/por-monitor-due.sh` → 得到 due 列表（空则退出，写「无窗口」笔记即可）。  
2. **只**对 due 所做公开探针（不要扫全集）。  
3. 新 id：`fetch` / 浏览器补抓 → `por report` → `por site` → `build-public-html` → PR。  
4. WAF/登录：同日标注缺口；需要时派 self-hosted 浏览器。  
5. **不要**在无新 id 时重跑 OKX `verify-global` 或全量 RPC。

### 定时（cursor-subscriptions）

建议：

| 名 | Cron (UTC) | 作用 |
| --- | --- | --- |
| `por-scan-due` | `0 2 * * *` | 主路径：due 窗轻探针 |
| `por-scan-due-catchup` | `0 14 * * *` | 仅复检 **今日仍在窗内** 或上午已标 pending publish 的所 |

废止：无条件的 `por-scan-am` / `por-scan-pm` 全量 P0+P1。

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

## SLA

**交易所公开更新后 ≤ 1 个自然日**内：更新 transparency + assessment、刷新首页与 comparison/history/archive，并开 PR。抓不到则同日写清 UNVERIFIABLE / WAF / 登录缺口。

窗口设计已加 1–2 日 buffer，使「更新日落在窗内 → 次日 catch-up 仍能扫到」，以支撑该 SLA；**不是**靠每日全扫。

## 不在范围内

- 登录态 cookie 常驻爬虫  
- 无 signer 时自动上链  
- 把 gitignored `artifacts/` 整包推到 Pages  
- 短名单所的完整适配器（见 onboarding 文；YAML `shortlist_preview` 仅可选窗）

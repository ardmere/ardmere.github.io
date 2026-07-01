# “Audited” is not effective PoR

> **Type:** Explainer · **Published:** 2026-06-28 (UTC) · **Author:** ardmere editorial  
> ← [Insights](./index.md)  
> **Related records:** [Binance PR01JUN26 report](../reports/binance/PR01JUN26-transparency-report.md) · [OKX 506872725 report](../reports/okx/506872725-transparency-report.md) · [Effective PoR Standard §1.3](../effective-por-standard.md#por-and-traditional-audit-answer-different-questions)

*This is an Insights piece—interpretive, not an artifact-bound rating. For Stage conclusions, read the linked reports.*

---

After every major exchange stress event, the same phrases return in press releases and social posts: **“financially audited,” “100% backed,” “ZK verified.”** Institutional and retail readers often treat these as interchangeable guarantees of solvency. They are not.

This note explains why **traditional audit and effective PoR answer different questions**, and how to read “audited” without conflating it with the [minimum effective standard](../effective-por-standard.md) ardmere applies in public reports.

---

## 1. Three sentences that sound alike

| What the exchange says | What many readers hear | What actually needs to be checked |
| --- | --- | --- |
| “Audited by a Big Four firm” | Customer assets are safe and fully reserved | Which entity was audited, which period, and whether **user custody liabilities** were in scope |
| “100% reserve ratio” | Every user can withdraw today | **Which liabilities** are in the ratio, at **which snapshot**, with **which artifacts** publicly reproducible |
| “ZK proof verified” | The whole exchange is solvent on-chain | Whether `wallet_address_list`, `wallet_ownership_proof`, and `global_proof` are **public** and bind to the same snapshot |

Confusion here is not academic. It drives licensing conversations, counterparty limits, and media narratives that lag behind what public artifacts actually show.

---

## 2. What financial audit is good at

Traditional audit and assurance work remain essential for **company-level** questions:

- Are consolidated financial statements fairly presented?
- Do internal controls over financial reporting exist and operate?
- Are off-chain bank balances, receivables, and related-party exposures addressed?
- Is going concern appropriately assessed?

For a listed parent or holding company, these outputs matter—and regulators rightly expect them on a **dual-track** basis alongside custody transparency. See [Why PoR is required](../ardmere-service-audience.md#why-por-is-required-for-digital-asset-exchanges) for the custody-specific case.

---

## 3. What audit typically does *not* give PoR readers

PoR is built for a different failure mode: **whether customer liabilities under a defined scope are backed by reserves the exchange controls, in a way third parties can reproduce without trusting marketing copy.**

Audit procedures generally do **not** provide:

1. **User-level reproducibility** — An audit opinion does not publish a Merkle or ZK path that any user or researcher can re-run against public roots.
2. **Wallet control at snapshot** — Confirming “cash and equivalents” on a balance sheet is not the same as batch-verifying `wallet_ownership_proof` for a published `wallet_address_list`.
3. **Liability-side proof scope** — Audit scope follows accounting standards; PoR scope must be **explicitly stated** (spot, margin, earn, off-balance-sheet exclusions).
4. **Timeliness between periods** — Quarterly or annual audit cycles do not close the gap between monthly—or daily—custody risk windows.

An exchange can therefore present a **clean audit narrative** while remaining **[Stage 0](../reports/exchange-comparison.md)** under a PoR framework: some artifacts exist, but independent third parties cannot reproduce the global asset–liability relationship.

---

## 4. A concrete contrast from the public report set

The first ardmere report set illustrates the boundary without naming malice—only **evidence completeness**.

### Binance (`PR01JUN26`) — Stage 0

Public reserve summary and `wallet_address_list` are available; ardmere reproduces several reserve-side checks. Stage 1 remains blocked because **`wallet_ownership_proof` is missing**, **trusted setup evidence is opaque**, and **`global_proof` / verifying keys are not publicly distributed**.

A reader who stops at “large exchange with sophisticated ZK marketing” misses the policy point: **Gen 2 technology in press materials ≠ Stage 1 effective PoR.**

Full record: [Binance transparency report](../reports/binance/PR01JUN26-transparency-report.md).

### OKX (`506872725`) — Stage 1

Public summary, `wallet_address_list`, batch-verifiable ownership proof, and public `global_proof` are obtainable without login walls. ardmere verifies address ownership and global zk binding—so this snapshot meets the **minimum effective** bar.

Stage 1 is **not** Stage 2: canonical anchoring, low-friction user inclusion proof, and business-consistent constraints remain open gaps documented in the report.

Full record: [OKX transparency report](../reports/okx/506872725-transparency-report.md).

---

## 5. How to read “audit taxonomy” in disclosures

When a third party is named, ask which bucket applies ([Standard §3](../effective-por-standard.md#report-types-that-must-not-be-conflated)):

| Type | Reasonable takeaway |
| --- | --- |
| **Technical verification** | Artifacts were run through agreed checks; read which artifacts |
| **AUP / agreed-upon procedures** | Steps were executed; usually **no** overall solvency opinion |
| **Limited assurance** | Limited scope on defined procedures |
| **Reasonable assurance / financial audit** | Company financials; **may not cover** PoR liability scope |

“Audited” without this taxonomy is not actionable for custody policy.

---

## 6. Practical checklist for regulators and allocators

1. **Separate tracks** — Require financial audit *and* PoR Stage evidence; do not accept one as proxy for the other.
2. **Demand artifact names** — `wallet_address_list`, `wallet_ownership_proof`, `global_proof`, scope statement, snapshot time.
3. **Treat Stage 0 as transitional** — Not effective PoR for licensing or institutional custody minimums ([Standard §2](../effective-por-standard.md#stage-thresholds-for-policy-use)).
4. **Follow the registry record** — Marketing upgrades; artifacts either exist or they do not. Missing items stay `UNVERIFIABLE`, not PASS.

---

## 7. Bottom line

**Audit answers “how did the company report its financial condition?”**  
**Effective PoR answers “can outsiders reproduce whether customer liabilities in scope are backed at this snapshot?”**

Until those questions are split in public discourse, “audited” will keep doing more reputational work than evidentiary work. Insights pieces like this one are meant to sharpen that distinction; **ratings stay in Reports.**

---

*Next in this series (planned): why user Merkle inclusion PASS is not exchange-wide solvency PASS—lessons from Merkle-only disclosures.*

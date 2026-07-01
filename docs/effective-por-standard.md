# Effective PoR Standard

> Full framework: [PoR Stage Framework](./por-transparency-framework.md)  
> Current assessments: [Exchange comparison](./reports/exchange-comparison.md)

This page states ardmere’s **minimum policy standard** for exchange Proof-of-Reserves (PoR). It is the public reference behind the registry’s Stage ratings—not a marketing checklist, and not a substitute for reading the bound artifacts in each report.

---

## 1. Three policy principles

### 1.1 Stage 0 is not effective PoR

An exchange may publish summaries, reserve ratios, or user Merkle inclusion proofs and still remain **Stage 0** if independent third parties cannot reproduce the global asset–liability relationship.

**Policy implication:** Stage 0 disclosures are transitional or observational. They must not be treated as meeting an effective PoR requirement for licensing, listing, or institutional custody decisions.

Common Stage 0 gaps include missing public `wallet_address_list`, `wallet_ownership_proof`, or `global_proof`; opaque trusted setup; or proof packages behind login walls.

### 1.2 Stage 1 is the minimum effective threshold

**Effective PoR** requires public, reproducible evidence on both sides of the snapshot:

- **Reserves:** `wallet_address_list`, block heights, and verifiable `wallet_ownership_proof`
- **Liabilities:** `global_proof`, root/vk/config, and clear proof scope
- **Parameters:** public trusted-setup transcript or a transparent proof system
- **Review:** independent third-party review with type and scope stated

Stage 1 does not mean “no risk.” It means key assertions are no longer accepted on exchange narrative alone.

### 1.3 PoR and traditional audit answer different questions

PoR and financial audit are **complementary**; neither replaces the other.

| Dimension | Traditional audit | PoR |
| --- | --- | --- |
| Core question | Are financial statements fairly presented? | Under a defined snapshot, are user liabilities in proof and do reserves cover that scope? |
| Best at | Off-chain assets, internal control, going concern | User-level commitments, on-chain balances, public reproducibility |
| Typical output | Audit opinion, AUP, assurance report | Artifacts, roots, proofs, verifier output |

**Policy implication:** Regulators and institutions should use a **dual-track** model—financial audit for company-level assurance, PoR Stage for publicly reproducible reserve transparency. An exchange with a clean audit opinion but Stage 0 PoR has not demonstrated effective PoR.

Gap reporting rules (`UNVERIFIABLE`, not PASS) are defined in the [methodology — §1.3](./por-transparency-framework.md#gap-reporting-discipline).

---

## 2. Stage thresholds for policy use

| Stage | Effective PoR? | Regulatory handling |
| --- | --- | --- |
| **Pre-Stage** | No | No minimally usable PoR; cannot be described as “having PoR” |
| **Stage 0** | No | Transitional disclosure only; must be labeled as relying primarily on exchange self-reporting |
| **Stage 1** | **Yes (minimum)** | Minimum standard for effective PoR in public custody contexts |
| **Stage 2** | Best practice | Long-term target: permissionless verification, anchoring, DA, business-consistent constraints |

Current Stage outcomes: [Exchange comparison](./reports/exchange-comparison.md).

Stage 1 operational requirements: [methodology §7](./por-transparency-framework.md#7-practical-assessment-checklist).

---

## 3. Report types that must not be conflated

Exchanges and reviewers should distinguish:

| Type | What it shows |
| --- | --- |
| **Technical verification** | Whether Merkle/ZK/wallet/proof artifacts are reproducible |
| **AUP / agreed-upon procedures** | Agreed steps executed; usually no overall assurance opinion |
| **Limited assurance** | Limited scope assurance on defined procedures |
| **Reasonable assurance / financial audit** | Company financial statements; may not cover PoR liability scope |

“Audited” does not automatically mean effective PoR. “ZK verified” does not automatically mean Stage 1 if critical artifacts are not public.

---

## 4. How ardmere applies this standard

Every public rating follows this **evidence chain**:

1. **Artifact** — bind to source materials (summary, wallet list, proof bundle)
2. **Hash** — record digests so the report can be reproduced and checked later
3. **Verifier output** — independent runs with explicit PASS, FAIL, WARN, or `UNVERIFIABLE`
4. **Gaps** — missing or blocked artifacts stay `UNVERIFIABLE`, not PASS

Each registry record also:

- Binds to a **specific snapshot** and disclosed scope
- Separates **Stage** (primary label) from **Gen / E** (explanatory lenses only)
- Links assessment JSON, bundles, hashes, and verifier output in the [artifact archive](./reports/artifact-archive-index.md)

Ratings are **records, not endorsements**. When marketing claims exceed public evidence, the report states the gap explicitly.

---

## 5. Related pages

- [PoR Stage Framework](./por-transparency-framework.md) — full methodology, Stage 2 path, industry recommendations
- [Exchange comparison](./reports/exchange-comparison.md) — current Stage outcomes
- [Artifact archive](./reports/artifact-archive-index.md) — evidence pointers for each report
- [Deposit spot-check](../verify/deposit/) — independent row-level reserve sampling

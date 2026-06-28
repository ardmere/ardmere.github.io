# PoR Stage Framework

This document defines ardmere’s transparency evaluation framework for exchange Proof-of-Reserves (PoR). For exchanges with a minimally usable PoR, **the final rating is PoR Stage (0 / 1 / 2)**, measuring how mature the exchange’s PoR is in trust minimization. Exchanges without a minimally usable PoR are labeled **Pre-Stage — No usable PoR** and do not enter Stage 0.

This framework draws on the design philosophy of the [L2BEAT Stages Framework](https://l2beat.com/stages): Stage is not a security certification, nor a ranking of technical sophistication. It answers a more fundamental question:

> Who must users ultimately trust?

`Technology Generation (Gen)` and `Evidence Level (E)` are not ratings on par with Stage. They are two auxiliary lenses for explaining Stage determinations:

- **Gen**: Describes what proof technology the exchange uses.
- **E**: Describes how much evidence external observers can obtain, reproduce, and cross-verify.

---

## 1. Framework Objectives

### 1.1 What This Framework Evaluates

PoR Stage evaluates the **degree of trust minimization** in an exchange’s PoR:


| Question | Meaning |
| -------- | ------- |
| Is evidence public? | Can third parties obtain artifacts such as summary, wallet_address_list, global_proof, wallet_ownership_proof, vk/config, etc.? |
| Is evidence reproducible? | Can third parties independently re-run core verification without exchange authorization? |
| Is the proof system trust-minimized? | Does the proof system rely on trusted setup, unverifiable parameters, proprietary implementations, or toxic waste held by the exchange? If trust assumptions exist, are they publicly verifiable? |
| Are proof constraints sufficient? | Can the proof system self-attest to key security properties—for example, no insertion of fake users with negative net worth, no offsetting real liabilities with negative balances, and consistency between user balance commitments and global liability aggregation? |
| Is history rewriteable? | Do root/proof/commitment have canonical records and a tamper-resistant publication layer? |
| Are trust assumptions converging? | Do users ultimately trust exchange narratives, review processes, or mathematics and public data? |


### 1.2 What This Framework Does Not Evaluate

Stage does not directly evaluate:

- How many checks a given verifier tool currently implements.
- Engineering coverage for a specific RPC, chain, or token.
- Whether code is completely bug-free.
- Overall exchange financial health, corporate governance, or regulatory compliance.

A higher Stage does not mean “no security risk”; a lower Stage does not mean “certainly malicious.” Stage only indicates: **how many trust assumptions PoR disclosure frees users from.**

---

## 2. Stage Overview


| Stage | Name | One-line Definition | Typical Evidence State | Core Trust Boundary |
| ----- | ---- | ------------------- | ---------------------- | ------------------- |
| **Pre-Stage** | No usable PoR | No minimally usable PoR; does not enter Stage | No snapshot / summary / proof, or PoR has stopped | Can only trust exchange verbal claims, brand, financial reports, or regulatory disclosures |
| **Stage 0** | Trust the Exchange | Has PoR, but third parties cannot independently reproduce global solvency | Summary or user inclusion primarily | Must still trust exchange asset ownership, liability completeness, or proof parameter claims |
| **Stage 1** | Verifiable Disclosure | Key artifacts are public and reproducible, but proof semantics remain primarily bounded by disclosure scope | wallet_address_list, wallet_ownership_proof, global_proof, and parameter sources all public; typically periodic snapshots | No longer trust exchange unilateral claims over wallet_address_list control, global_proof availability, or trusted-setup honesty, but still trust that proof boundaries adequately reflect business risk |
| **Stage 2** | Trust-minimized PoR | Proof and data are permissionless verifiable; proof constraints align with business risk; risk windows compressed by high-frequency anchoring | Both sides public + low-barrier inclusion proof + business-consistent constraints + on-chain anchoring + DA + high-frequency/event-triggered updates | Further freed from trusting exchange self-interpretation of business risk, maintenance of the sole historical version, continuous data provision, or arbitrary extension of unproven windows |


**Current industry state**: Mainstream CEXs peak at approximately **Stage 1**; no complete **Stage 2** exists yet.

**Core distinction between Stage 0 and Stage 1**: Stage 0 has PoR artifacts, but key assertions are still provided unilaterally by the exchange; Stage 1 requires the exchange to publish wallet_address_list, wallet_ownership_proof, global_proof, and trusted setup / parameter sources so third parties can independently re-verify these core assertions. In other words, the Stage 1 transition is not “more disclosure”—it is the exchange’s unilateral narrative being constrained by public evidence.

“Key assertions” here include: the exchange actually controls addresses in wallet_address_list; published on-chain balances actually belong to the current period’s reserves; global_proof actually corresponds to the current period’s liability summary; vk/config used in the proof matches public parameters; and if trusted setup is used, the exchange does not unilaterally hold toxic waste.

**Core distinction between Stage 1 and Stage 2**: Stage 1 addresses “whether this period’s artifacts are publicly reproducible” and requires clear proof boundaries; Stage 2 further requires proof constraints to align with the exchange’s actual business and to self-attest that no fake users with negative net worth were inserted and no negative balances were used to offset real liabilities. For example, exchanges involving lending, margin, portfolio margin, or collateral haircuts must constrain the corresponding risk-control logic in the circuit or equivalent proof. Stage 2 also requires low-barrier verification paths for user-side inclusion proof, high-frequency anchoring of root/proof/commitment, event-triggered supplementary issuance, or near-real-time proofs to reduce timing, substitution, and historical rewrite risk within snapshot intervals.

---

## 3. Stage Determination Method

### 3.1 Determination Principles


| Principle | Meaning |
| --------- | ------- |
| Stage-first | External rating uses Stage as the primary label; Gen / E only explain Stage |
| Threshold-based | Determine Stage first via hard checklist; do not upgrade via average scores |
| Shortest-board model | PoR is a weakest-link system; critical shortcomings block Stage upgrades or reduce confidence |
| Evidence-first | Ratings must bind to obtainable raw artifacts, hashes, URLs, reports, or on-chain data |
| Public-by-default | Data without login, cookies, or manual application carries more weight than post-login user data |
| Reserve and liability separation | wallet_address_list, liability-side Merkle/ZK, and summary ratio are different data planes |
| UNVERIFIABLE explicit | Missing data is not PASS; dimensions that cannot be verified must be stated clearly |
| Snapshot-scoped | Each rating applies to a specific snapshot period; cannot be assumed to extend to all historical/future periods |
| No overclaiming | Merkle inclusion PASS cannot be stated as full-exchange solvency PASS |
| Business-consistent constraints | PoR proof constraints must align with the exchange’s actual business; risks such as fake users with negative net worth, negative balance offsetting, lending, margin, and collateral haircuts cannot rely on prose alone |
| User-verification usability | inclusion proof verification should be low-barrier, explainable, and exportable; ordinary users must not be required to run complex commands or understand underlying circuits |
| Machine-readable first | Machine-readable artifacts such as JSON/CSV/ZIP/schema/API take priority over unparseable PDFs or marketing pages |
| Audit taxonomy | Must distinguish technical verification, AUP/agreed-upon procedures, limited assurance, reasonable assurance/full financial audit |
| Frequency-aware | Transparency considers not only single-period snapshots but also publication frequency, event-triggered disclosure, historical version archivability, and anti-timing capability |


### 3.2 Pre-Stage — No usable PoR

**Definition**: The exchange has no minimally usable PoR, or historically had PoR but has currently stopped. Users may refer to other trust sources such as exchange brand, traditional financial reports, or regulatory disclosures, but cannot obtain the most basic PoR artifacts for review.

**Any one of the following qualifies as Pre-Stage**:


| # | Situation | Determination Criteria |
| - | --------- | ---------------------- |
| P-1 | No ongoing PoR | No public PoR in the last 12 months |
| P-2 | PoR stopped | Had PoR historically, but no longer publishes |
| P-3 | Marketing claims only | Only vague statements such as “100% backed,” with no snapshot, summary, root, proof, or review report |
| P-4 | Traditional financial reports only | Has financial audit or public company disclosure, but no PoR artifacts mappable to user asset/liability side |
| P-5 | On-chain balance display only | Only wallet balance dashboard, with no user liability-side summary / inclusion / proof |


**Distinction between Pre-Stage and Stage 0**:


| State | What Users Trust | Minimum PoR Evidence |
| ----- | ---------------- | -------------------- |
| **Pre-Stage** | Brand, financial reports, regulatory disclosures, or other non-PoR trust sources | None |
| **Stage 0** | Exchange narrative + minimum PoR artifacts | Has snapshot, summary, user inclusion, or basic verifier |


**Examples**:

- Coinbase (exchange-level): Pre-Stage (has public company financial disclosure, but no exchange-level PoR artifacts).
- KuCoin: Pre-Stage (PoR stopped).

### 3.3 Stage 0 — Trust the Exchange

**Definition**: The exchange publishes PoR; users can at most verify their own inclusion or read summary, but independent third parties cannot reproduce the global asset–liability relationship.

**Hard thresholds (all must be met)**:


| # | Requirement | Determination Criteria |
| - | ----------- | ---------------------- |
| S0-1 | Ongoing PoR | Public PoR in the last 12 months, not stopped |
| S0-2 | Snapshot boundary | Clear snapshot time and asset coverage scope |
| S0-3 | Minimum disclosure | At least one of: self-reported summary (reserve ratio + covered assets); or user Merkle inclusion proof |
| S0-4 | Reproducibility floor | If user proof is provided: leaf/hash rules are public, proof is exportable, third parties can locally re-verify using user-provided proof |
| S0-5 | Tool or specification | Open-source verifier, or sufficiently detailed algorithm description |


**Allowed gaps (therefore cannot reach Stage 1)**:

- No public wallet_address_list.
- No public global_proof.
- No public wallet_ownership_proof.
- No independent third-party review, or review type/scope not disclosed.
- No on-chain anchoring.

**Stage 0 trust assumption**: Users trust “whatever reserves the exchange claims are the reserves.”

**Examples**:

- Bitget: weak Stage 0 (64-bit truncated hash, no independent attestation).
- Bybit: Stage 0 (user Merkle, but no public wallet_address_list / ZK).
- Gate.io production environment: Stage 0 (summary public, ZK package behind login, no wallet_address_list).

### 3.4 Stage 1 — Verifiable Disclosure

**Definition**: Asset-side wallet_address_list and wallet_ownership_proof, and liability-side global_proof are all publicly obtainable; proof parameters / trusted setup sources are verifiable; any third party can in principle independently re-verify. However, canonical records and historical versions of proof/root/vk still depend on the exchange or a small number of reviewers. “Public” here means obtainable without login, registration, cookies, API key, or manual application.

**Hard thresholds (all must be met)**:


| # | Requirement | Determination Criteria |
| - | ----------- | ---------------------- |
| S1-1 | Meets Stage 0 | All Stage 0 thresholds |
| S1-2 | Asset side public | Public wallet_address_list + block height, on-chain verifiable |
| S1-3 | Address control proof | Public wallet_ownership_proof or other verifiable control proof, batch-verifiable |
| S1-4 | Liability side public | Public global_proof + root/vk/config, locally re-runnable or structurally reproducible; must not require login, registration, cookies, API key, or manual application |
| S1-5 | Permissionless core artifacts | Core artifacts such as summary, wallet_address_list, wallet_ownership_proof, global_proof, vk/config downloadable without login, registration, cookies, API key, or manual application |
| S1-6 | Independent review | Independent third-party review with disclosed type (technical verification / AUP / limited assurance / reasonable assurance) |
| S1-7 | Bidirectional cross-binding | summary ↔ wallet_address_list aggregate, and summary ↔ proof root / commitment are both comparable |
| S1-8 | Parameter transparency | If trusted setup is used, must publish verifiable transcript / MPC ceremony; or use a proof system requiring no trusted setup |
| S1-9 | Historically archivable | Historical snapshot artifacts can be archived, not only current-period web display |
| S1-10 | Proof boundary | Clearly disclose which liabilities proof covers and which off-balance-sheet/off-chain items it does not |


**Stage 1 equivalent statement**:

> Aside from implementation bugs, for an independent third party to **long-term fail to detect** forged solvency relationships, collusion between the exchange and reviewers, or the exchange’s unilateral ability to replace/delete server-hosted artifacts must be relied upon; it can no longer rely on unverifiable assumptions such as “the exchange actually controls addresses in wallet_address_list” or “the exchange does not hold toxic waste.”

**Key gaps still allowed (therefore cannot reach Stage 2)**:

- proof / root / vk not on-chain anchored.
- global_proof requires login, registration, cookies, API key, or manual application; sample-only packages; or cannot bind to current-period wallet_address_list / summary.
- Monthly snapshots only, no event-triggered supplementary issuance.
- Undisclosed conflicts of interest between reviewer and technology provider.

**Stage 1 trust assumption**: Users no longer trust the exchange’s word alone, no longer trust that the exchange controls its wallet_address_list, and no longer trust that the exchange does not hold trusted-setup secrets; but still trust that the exchange will not replace current-period artifacts, split-view different observers, or conceal liability scope.

**Examples**:

- OKX: **Stage 1** (wallet_address_list + wallet_ownership_proof + global_proof; proof still distributed off-chain).

### 3.5 Stage 2 — Trust-minimized PoR

**Definition**: Building on Stage 1 bilateral public disclosure, user-side inclusion proof has a low-barrier verification path; proof constraints align with the exchange’s actual business; root / proof / commitment have canonical on-chain or DA records; any third party can permissionless re-verify without exchange authorization; anchored history cannot be unilaterally rewritten; and risk windows between full PoR cycles are compressed via high-frequency anchoring, event-triggered supplementary issuance, or near-real-time proofs.

**Hard thresholds (all must be met)**:


| # | Requirement | Determination Criteria |
| - | ----------- | ---------------------- |
| S2-1 | Meets Stage 1 | All Stage 1 thresholds |
| S2-2 | Business-consistent constraints | Proof constraints cover the exchange’s actual business risk and can self-attest no fake users with negative net worth were inserted and no negative balances offset real liabilities; businesses involving lending, margin, portfolio margin, collateral haircuts, price parameters, or risk limits must be constrained in the circuit or equivalent public proof |
| S2-3 | On-chain anchoring | root / proof root / commitment on-chain, with non-deletable historical commitment chain |
| S2-4 | Data Availability | proof packages, vk, config, wallet_address_list bundle long-term available on stable public layer (on-chain, IPFS, third-party mirror/DA) |
| S2-5 | Low-barrier user verification | inclusion proof exportable, locally reproducible, with Web/WASM/GUI one-click verification and clear error messages; verification results bound to public root / global_proof / summary |
| S2-6 | Permissionless | Third parties can re-run core verification without login or API key application; if ZK is used, prefer on-chain verifier or equivalent public verification path |
| S2-7 | Parameter and version anchoring | vk/config/verifier version hashes or commitments fixed on-chain/DA; changes have audit trail |
| S2-8 | Anti-timing frequency | For top CEXs: at least one of “weekly full PoR + daily root/commitment anchoring + event-triggered supplementary issuance”; long-term monthly snapshots alone insufficient to maintain Stage 2 |
| S2-9 | Review independence | Conflicts of interest between reviewer and technology provider disclosed; no “athlete and referee” without explanation |


**Stage 2 equivalent statement**:

> Any user can low-barrier verify their own inclusion proof; any external observer can independently verify “whether the proof for this period’s canonical root holds”; and that proof can self-attest key security properties and constrain the exchange’s actual business risk; the exchange cannot unilaterally rewrite anchored history or arbitrarily extend unproven risk windows.

**Strong bonus items**:

- On-chain verifier contract can directly verify ZK proof.
- proof generation / verifier versions are traceable.
- Supports continuous proof or high-frequency snapshots (daily or higher).
- Complete historical commitment chain; proof/vk/root changes have audit trail.
- Standardized proof schema for lending, margin, derivatives netting, collateral haircuts, and price parameters.

**Current state**: Mainstream exchanges have no complete Stage 2 yet. Third-party archivers can put roots of captured artifacts on-chain, but this only proves “this third party has seen and archived this artifact set”—it **cannot** substitute for the exchange’s official Stage 2.

**Stage 2 trust assumption**: Users primarily trust cryptographic proofs, public data, and on-chain canonical records.

Stage 2 is **trust-minimized**, not trustless. It still relies on correct proof system implementation, truthful PoR coverage scope, honest disclosure of off-chain assets and off-balance-sheet liabilities, and review scope not being misread. If an exchange actually operates lending, margin, or derivatives business but PoR only proves spot balances or simple total liabilities, it cannot be considered Stage 2.

### 3.6 Stage Confidence Adjustments

Stage is a maturity label, not a security score. Confidence markers can be appended without changing Stage:


| Marker | Meaning |
| ------ | ------- |
| high confidence | Complete bilateral public disclosure, independent review, low implementation risk |
| medium confidence | Complete bilateral public disclosure, but off-chain distribution, user accessibility, or historical archiving insufficient |
| low confidence | sample-only, login walls, known HIGH defects, or extremely narrow attestation scope |


**Shortcoming handling rules**:


| Situation | Handling |
| --------- | -------- |
| Sample proof only, not current production | No production Stage 1; mark sample-only |
| Proof package requires login, registration, cookies, API key, or manual application | If liability side cannot be publicly reproduced, max Stage 0 |
| inclusion proof not exportable, not locally reproducible, or only black-box verifiable | Max Stage 0 |
| inclusion proof reproducible but high-friction verification flow | Max Stage 1 |
| No wallet_ownership_proof | Max Stage 0 |
| Trusted setup opaque, or no verifiable transcript | Max Stage 0 |
| No independent third-party review | Max Stage 0 |
| proof/root/vk server-hosted only, no archiving/anchoring | Max Stage 1 |
| Major unfixed implementation defects affecting proof semantics | Stage unchanged, confidence downgraded |
| Undisclosed conflicts of interest between reviewer and technology provider | Stage unchanged, confidence downgraded |


---

## 4. Auxiliary Lenses for Explaining Stage

### 4.1 Technology Lens: Technology Generation (Gen)

Gen describes what proof technology the exchange uses. It helps explain Stage but **does not directly determine Stage**.


| Gen | Name | Minimum Requirement | Examples |
| --- | ---- | ------------------- | -------- |
| **Gen 0** | No PoR | No ongoing PoR, or only traditional financial reports/marketing claims | Coinbase (exchange-level), KuCoin (stopped) |
| **Gen 1** | Merkle inclusion | Users can verify their own inclusion; no ZK solvency constraint | Bitget, MEXC, Bybit, Kraken |
| **Gen 2** | Merkle + ZK | Has global ZK / zk-SNARK / zk-STARK solvency proof | Binance, OKX, Gate.io, HTX |
| **Gen 3** | ZK + onchain + DA | Public proof verifiable on-chain; root/proof has tamper-resistant publication layer | None in industry yet |


**How Gen explains Stage**:

- Gen 0 typically corresponds to Pre-Stage; does not enter PoR Stage.
- Gen 1 typically reaches only Stage 0, unless it also has public wallet_address_list, wallet_ownership_proof, global_proof, review, and cross-verification mechanisms.
- Gen 2 indicates potentially stronger liability-side proof; Stage 1 foundation exists only when wallet_address_list, wallet_ownership_proof, global_proof, and trusted setup / parameter sources are all publicly verifiable.
- Gen 3 approaches Stage 2 technical form, but still must meet bilateral public disclosure, low-barrier inclusion proof, parameter transparency, DA, permissionless, frequency, and other thresholds.

**Why Gen does not directly determine Stage**:

- ZK proof can be advanced, but if proof packages are not public, trusted setup is not verifiable, or root is not anchored, users still must trust the exchange distribution layer or setup honesty.
- Merkle technology is more basic, but if asset side, wallet_ownership_proof, review, and historical records are more public, Stage may exceed some “off-chain ZK showcase.”
- Third-party review does not change Gen; it affects Stage thresholds and confidence.

### 4.2 Evidence Disclosure Lens: Evidence Level (E)

E describes how much evidence external observers can obtain and reproduce. It is evidence input for Stage determination, **not the final transparency rating**. E primarily describes “how far evidence is disclosed,” while also observing whether user-side verification paths are usable and low-barrier; but it does not automatically judge whether ownership, trusted setup, business-consistent constraints, or publication frequency meet Stage thresholds.


| E | Name | One-line Definition | Common Stage Interpretation |
| - | ---- | ------------------- | ----------------------------- |
| **E0** | No usable evidence | No usable PoR artifacts, or PoR has stopped | Pre-Stage |
| **E1** | Limited disclosure | Has summary, reserve ratio, user inclusion, or basic verifier, but cannot reproduce global relationship; verification path may still be complex | Stage 0 |
| **E2** | Public verifiable artifacts | Has permissionless public artifacts; asset side, liability side, or key parts reproducible | Stage 0 or Stage 1, depending on whether hard thresholds are complete |
| **E3** | Anchored permissionless evidence | E2 + low-barrier inclusion proof + on-chain anchoring / DA / permissionless verification / high-frequency updates | Stage 2 candidate |


**How E supports Stage determination**:

- E0 corresponds to Pre-Stage; does not enter PoR Stage.
- E1 indicates the exchange has minimum PoR disclosure, but usually insufficient to exceed Stage 0; if inclusion proof verification requires complex CLI or manual parameter assembly, user verifiability score should be reduced.
- E2 indicates publicly reproducible artifacts exist, but wallet_address_list, wallet_ownership_proof, global_proof, parameter sources, business-consistent constraints, and review independence must still be checked to determine Stage 1.
- E3 is the evidence form for Stage 2, but must simultaneously meet Stage 2 low-barrier user verification, business-consistent constraints, permissionless, DA, parameter and version anchoring, frequency, and independence requirements.

**E Level detailed determination points**:

- **E0**: No snapshot, summary, root, proof, wallet_address_list, verifier, or PoR has stopped.
- **E1**: Has snapshot, summary, reserve ratio, user inclusion proof, or basic verifier; but third parties cannot reproduce global asset–liability relationship. If user verification requires complex commands, environment setup, or manual proof handling, mark as high-friction E1.
- **E2**: Among core artifacts such as wallet_address_list, wallet_ownership_proof, global_proof, root/vk/config, proof schema, trusted setup transcript, some are publicly reproducible; if key items are missing, may still be only Stage 0.
- **E3**: On E2 basis, inclusion proof has low-barrier user verification path; root/proof/commitment has canonical anchoring; artifacts have stable DA; third parties can permissionless replay verification; and has higher-frequency or event-triggered updates.

### 4.3 Gen / E to Stage Explanation Mapping


| Combination | Stage Interpretation |
| ----------- | -------------------- |
| Gen 0 / E0 | Pre-Stage: no minimally usable PoR |
| Gen 1 + E1 | Typically Stage 0: has summary or user self-verification, but cannot reproduce global solvency |
| Gen 2 + E1 | Typically Stage 0: may have ZK technically, but production proof restricted from public or key artifacts incomplete |
| Gen 2 + E2 | Stage 0 or Stage 1: depends on whether wallet_ownership_proof, global_proof, parameter sources are complete |
| Gen 2/3 + E3 | Stage 2 candidate: still must check low-barrier inclusion proof, business-consistent constraints, permissionless, DA, frequency, parameter and version anchoring |


**Counterexamples**:

- Binance: Gen 2 + E2, wallet_address_list and global_proof both public, but missing wallet_ownership_proof and public trusted setup transcript; therefore does not enter Stage 1.
- OKX: Gen 2 + E2, wallet_address_list, wallet_ownership_proof, and global_proof all public; therefore can enter Stage 1; but proof still off-chain distributed, therefore not Stage 2.
- Bybit: Gen 1 + E1, has user Merkle and third-party institutional verification, but no public wallet_address_list / ZK; therefore still Stage 0.

---

## 5. Within-Stage Scoring

Stage is threshold-based; within the same Stage, a 100-point scale can be used for horizontal ranking:


| Category | Weight | Scoring Notes |
| -------- | ------ | ------------- |
| Public data availability | 15 | summary, history, wallet_address_list, global_proof, URL stability, machine readability |
| Liability-side proof strength and scope | 20 | Merkle / Sum Tree / ZK / per-user solvency / fake negative-net-worth user defense / off-balance-sheet liability boundary |
| Reserve-side verifiability and quality | 20 | wallet_address_list, block height, on-chain verifiability, clean assets, self-issued tokens, staking/custody classification |
| Tamper resistance, time anchoring, and frequency | 15 | On-chain root, third-party mirrors, DA, historical commitment, weekly/daily/event-triggered disclosure |
| Review independence and assurance level | 10 | Technical verification, AUP, limited assurance, reasonable assurance, reviewer qualifications and conflicts of interest |
| User verifiability | 10 | User proof, export capability, Web/WASM/GUI one-click verification, open-source verifier, clear error messages, real participation threshold and verification rate |
| Proof system and implementation risk | 10 | trusted setup, transparent setup, hash truncation, proof constraint sufficiency, parameter consistency, code maintenance, testing, audit finding remediation status |


**Floor score rules** (Stage unchanged, but total score capped):


| Condition | Total Score Cap |
| --------- | --------------- |
| No user liability-side data | 40 |
| No public asset-side data | 65 |
| No wallet_address_list and no ZK | 55 |
| Hash output <128 bit | 45 |
| No recent PoR | 20 |
| proof/root/vk server-hosted only with no archiving/anchoring | 80 |
| No wallet_ownership_proof | 60 |
| Trusted setup opaque or no verifiable transcript | 60 |
| Major unfixed implementation defects affecting proof semantics | 70 |
| Undisclosed conflicts of interest between reviewer and technology provider | 75 |
| Has lending/margin business but proof does not constrain risk-control logic | 70 |
| Cannot prove no fake users with negative net worth or negative balance offsetting real liabilities | 70 |


---

## 6. Exchange Mapping


| Exchange | **PoR Stage** | Gen | E Level | Stage Rationale |
| -------- | ------------- | --- | ------- | --------------- |
| Coinbase (exchange-level) | Pre-Stage | Gen 0 | E0 | Has public company financial disclosure, but no exchange-level PoR artifacts |
| KuCoin | Pre-Stage | Gen 0 | E0 | PoR stopped |
| Binance | Stage 0 | Gen 2 | E2 | wallet_address_list and global_proof public; but no wallet_ownership_proof, no public trusted setup transcript; Stage 1 blocked |
| OKX | **Stage 1** | Gen 2 | E2 | wallet_address_list, wallet_ownership_proof, and global_proof public; proof off-chain distributed |
| Gate.io | Stage 0 | Gen 2 | E1 | summary public; ZK package behind login; no wallet_address_list |
| HTX | Stage 0 | Gen 2 | E1 | Sample ZK package verifiable; production data and wallet_address_list gaps significant |
| Bybit | Stage 0 | Gen 1 | E1 | User Merkle + third-party institutional verification; no wallet_address_list / ZK |
| Bitget | Stage 0 | Gen 1 | weak E1 | 64-bit truncated hash; no independent attestation |


---

## 7. Practical Assessment Checklist

Before rating an exchange’s Stage, answer in order:

### 7.1 Pre-Stage checklist

If any one is met, mark as `Pre-Stage — No usable PoR`:

1. No public PoR in the last 12 months?
2. Has PoR stopped?
3. Marketing claims only, with no snapshot, summary, root, proof, or review report?
4. Traditional financial reports only, with no PoR artifacts mappable to user asset/liability side?
5. On-chain balance display only, with no user liability-side summary / inclusion / proof?

### 7.2 Stage 0 checklist

1. Has PoR been published continuously in the last 12 months?
2. Is there clear snapshot time and covered assets?
3. Is there summary or user inclusion proof?
4. If user proof is provided: are rules public, proof exportable, can third parties locally re-verify?
5. Does inclusion proof provide a low-barrier verification path (e.g., Web/WASM/GUI one-click verification), avoiding ordinary users having to run complex CLI, configure environments, or manually assemble parameters?
6. Is there an open-source verifier or sufficient algorithm description?

### 7.3 Stage 1 checklist

1. Is asset-side wallet_address_list public, login-free, and on-chain verifiable?
2. Is wallet_ownership_proof public and batch-verifiable?
3. Is liability-side global_proof downloadable and independently reproducible without login, registration, cookies, API key, or manual application?
4. Is there independent third-party review? Is the type disclosed (technical verification / AUP / limited assurance / reasonable assurance)?
5. Can both summary ↔ wallet_address_list aggregate and summary ↔ proof root / commitment be bound?
6. If trusted setup is used, is verifiable transcript / MPC ceremony public? If not using trusted setup, is the proof system transparent setup?
7. Are historical snapshot artifacts archivable?
8. Is user-side inclusion proof verification experience sufficiently low-barrier and bound to public global_proof / root / summary?
9. Is proof boundary clear (which liabilities covered, which off-balance-sheet items not covered)?

### 7.4 Stage 2 checklist

1. Are Stage 1 wallet_address_list + wallet_ownership_proof + global_proof + parameter transparency thresholds met?
2. Do proof constraints align with the exchange’s actual business?
3. Can it self-attest no fake users with negative net worth were inserted, no negative balances offset real liabilities, and user balance commitments match global liability aggregation?
4. If lending, margin, portfolio margin, collateral haircuts, negative balances, price parameters, or risk limits exist, are they constrained in the circuit or equivalent public proof?
5. Is inclusion proof exportable, locally reproducible, with Web/WASM/GUI one-click verification, clear error messages, and bound to public global_proof / root / summary?
6. Are root / proof / commitment on-chain anchored with a historical commitment chain?
7. Are proof packages, vk, config long-term available on stable DA layer?
8. Can third parties permissionless re-run core verification (no login/API key)?
9. Are vk/config/verifier versions fixed on-chain/DA with audit trail for changes?
10. Is publication frequency sufficient (weekly full PoR / daily anchoring / event-triggered supplementary issuance)?
11. Are conflicts of interest between reviewer and technology provider disclosed?

### 7.5 General risk items

1. Is reserve quality classified (clean / staked / off-chain / self-issued tokens)?
2. Are there major unfixed implementation defects?
3. Does the report clearly state what cannot be proven?

---

## 8. Recommendations for Regulators

PoR Stage can serve as tiered language for regulators setting exchange reserve transparency requirements, but it is a different dimension from traditional financial audit. Regulatory text should clarify: PoR proves asset–liability relationships under a specific snapshot and disclosure scope; traditional audit covers financial statements, internal controls, corporate governance, off-chain assets, and going concern at the company level. The two should be used complementarily.

### 8.1 Minimum Effective PoR Requirements

Regulators should not accept `Pre-Stage` or **Stage 0** as effective PoR. If an exchange custodies customer assets for the public, the minimum requirement for effective PoR should be at least **Stage 1**; Stage 0 can only serve as transitional disclosure, remediation observation, or risk notice—not as basis for meeting regulatory PoR requirements.


| Requirement | Regulatory Handling Recommendation |
| ----------- | ---------------------------------- |
| Pre-Stage | Should not be stated as “has PoR”; can only be marked as no minimally usable PoR |
| Stage 0 | Should not be recognized as effective PoR; can only serve as transitional disclosure or remediation status, and must be explicitly marked as still primarily relying on exchange self-reporting |
| Stage 1 | Can serve as minimum regulatory standard for effective PoR: wallet_address_list, wallet_ownership_proof, global_proof, parameter sources all publicly verifiable |
| Stage 2 | Can serve as long-term best practice: low-barrier inclusion proof, on-chain anchoring, DA, permissionless verification, continuous traceability |


### 8.2 Report Types That Should Not Be Conflated

Regulators should require exchanges to clearly distinguish in external disclosure:


| Type | Meaning |
| ---- | ------- |
| Technical verification | Checks whether technical artifacts such as Merkle/ZK/wallet_address_list/wallet_ownership_proof/global_proof are reproducible |
| AUP / agreed-upon procedures | Reviewer executes agreed check steps, but typically does not express overall assurance opinion |
| Limited assurance | Provides limited degree of assurance for specific scope |
| Reasonable assurance / full financial audit | Closer to traditional audit opinion, but must still state whether PoR liability scope is covered |


“Has audit” and “has PoR” answer different questions. Traditional financial audit provides company-level assurance; PoR provides publicly reproducible evidence for users and third parties. Regulatory and user disclosures should state each one’s coverage scope.

### 8.3 Differences and Complementarity of Traditional Audit and PoR

Traditional audit and PoR answer different questions; regulatory requirements should retain both and state their respective boundaries.


| Dimension | Traditional Audit | PoR |
| --------- | ----------------- | --- |
| Core question | Whether company financial statements are fairly presented in all material respects | Whether user liabilities under a snapshot are included in proof, and whether reserves suffice to cover disclosed scope |
| Time semantics | Typically quarterly, annual, or specific period | Typically a snapshot; can extend to high-frequency or continuous proof |
| Asset verification | Covers cash, bank deposits, receivables, off-chain assets, investments, etc. | Better suited for on-chain assets, wallet_ownership_proof, on-chain balances |
| Liability verification | Confirms liabilities per accounting scope, or via sampling/procedures | Can make reproducible commitments on full user balance commitments, Merkle/ZK inclusion, global_proof |
| Verification actors | Auditors, regulators, investors | Users, third-party researchers, regulators, automated verifiers |
| Output form | Audit opinion, AUP, assurance report, financial statement notes | root, proof, wallet_address_list, wallet_ownership_proof, vk/config, verifier output |
| Main blind spots | Does not necessarily provide user-level reproducible evidence; may not cover real-time custody risk | Cannot automatically prove off-chain assets, off-balance-sheet liabilities, internal controls, or operating capacity |


Complementary use:

1. **Traditional audit covers off-chain and company-level risk**: Bank deposits, receivables, related-party transactions, off-balance-sheet arrangements, internal controls, going concern—areas traditional audit handles better.
2. **PoR covers user-level and on-chain reproducible risk**: User inclusion, global_proof, wallet_ownership_proof, on-chain balances, public reproducibility of proof/root/commitment—parts traditional audit typically does not give users direct verification of.
3. **Accounting scope and PoR scope must align**: Regulators should require disclosure of differences between “what PoR covers” and “what financial statements recognize,” especially lending, margin, derivatives, related parties, and off-balance-sheet items.
4. **Auditors should not only verify exchange-generated conclusions**: If review scope only covers exchange-provided summary without verifying wallet_address_list, wallet_ownership_proof, global_proof, parameter sources, and historical artifacts, it cannot be stated as PoR Stage 1-level review.
5. **Regulators should adopt dual-track requirements**: Traditional audit for financial integrity and off-chain risk; PoR Stage for publicly reproducible reserve transparency. Both together better address user concerns about “redeemability.”

### 8.4 Fields Regulators Should Mandate for Disclosure

Regulatory reports should at minimum require machine-readable disclosure of:

1. Snapshot time, covered assets, covered chains, covered user liability scope.
2. reserve summary, liability summary, reserve ratio.
3. wallet_address_list, block height, balances, networks, asset classification.
4. wallet_ownership_proof and verification method.
5. global_proof, root, vk/config, proof schema.
6. Correspondence between proof constraints and business model and key security properties: whether it can self-attest no fake users with negative net worth, no negative balance offsetting real liabilities; whether lending, margin, portfolio margin, collateral haircuts, price parameters, and risk limits are constrained.
7. trusted setup transcript / MPC ceremony, or transparent setup proof system description.
8. Third-party review institution, review type, scope limitations, conflict-of-interest declarations.
9. Archival location, hashes, historical versions, and anchoring records for proof/root/commitment.
10. Publication frequency, event-triggered supplementary issuance rules, historical snapshot retention policy.
11. Clear statement of assets, liabilities, and off-balance-sheet items PoR does not cover.

### 8.5 Frequency and Event Triggers

For large CEXs, monthly PoR can only serve as a transitional baseline. Regulators should consider tiered higher-frequency requirements by risk:


| Scenario | Recommendation |
| -------- | -------------- |
| Normal ongoing disclosure | At least monthly full PoR; historical artifacts archivable |
| High-volume or high-leverage platforms | Weekly full PoR, or at least weekly liability-side updates |
| Extreme market volatility, bank runs, major security incidents | Event-triggered supplementary PoR |
| Stage 2 target | Low-barrier inclusion proof + daily root/commitment anchoring, progressing toward daily full PoR or near-real-time proof |


### 8.6 Enforcement and User Presentation

Regulators should require exchanges to display in a unified format:

```text
PoR Stage: Stage 1
Technology Generation: Gen 2
Evidence Level: E2
Scope: spot + margin liabilities, excludes off-balance-sheet exposures
Snapshot: 2026-06-01T00:00:00Z
```

If an exchange lacks wallet_address_list, wallet_ownership_proof, trusted setup transcript, public global_proof, key proof constraint documentation, or historical artifacts, the user interface should explicitly mark `UNVERIFIABLE`; marketing phrases such as “100% backed,” “audited,” or “ZK verified” must not cover critical gaps.

---

## 9. Recommendations for the Industry

The long-term goal of PoR is not for exchanges to publish more PDFs or marketing pages, but for users, third-party researchers, review institutions, and regulators to all be able to re-verify the same set of core facts: which assets the exchange controls, which liabilities it commits to, whether proof constraints cover real business risk, and whether this evidence can be traced long-term.

### 9.1 Recommendations for Exchanges

Exchanges should treat **Stage 1** as the near-term minimum target and **Stage 2** as the long-term build direction:

1. **Do not stop at Stage 0**: summary, reserve ratio, or user Merkle proof only demonstrate minimum disclosure; they do not prove global solvency.
2. **Publish bilateral core artifacts**: wallet_address_list, wallet_ownership_proof, global_proof, root/vk/config, proof schema, and trusted setup transcript should be obtainable without login, registration, cookies, API key, or manual application.
3. **Disclose proof coverage boundaries**: Clearly state which assets, liabilities, product lines, and risk scenarios PoR covers, and which off-chain assets, off-balance-sheet liabilities, lending, margin, derivatives, or related-party items it does not.
4. **Align proof constraints with business model**: PoR should not only prove static balances; it should self-attest no fake users with negative net worth were inserted and no negative balances offset real liabilities. If the exchange has lending, margin, portfolio margin, collateral haircuts, or price parameters, explain how these risk-control logics are constrained or why they are not covered.
5. **Lower user verification barriers**: inclusion proof should not be CLI-only for developers; exchanges should provide Web/WASM/GUI one-click verification, clear error messages, proof export, and local re-verification paths, and disclose verification entry points and participation rate metrics.
6. **Prioritize machine-readable and archivable**: JSON/CSV/ZIP/schema/API, open-source verifiers, and historical artifacts should take priority over unparseable PDFs; each snapshot period should have stable URL, hash, and version records.
7. **Increase publication frequency**: Monthly PoR can only serve as a transitional baseline; high-volume or high-leverage platforms should progress toward weekly, daily anchoring, or event-triggered supplementary issuance.

### 9.2 Recommendations for Third-Party Review Institutions

Third-party institutions should not merely restate exchange-generated conclusions; they should clearly state which artifacts they verified and which boundaries they did not:

1. Distinguish technical verification, AUP / agreed-upon procedures, limited assurance, and reasonable assurance; do not use “audited” to blur scope differences.
2. At minimum check consistency of wallet_address_list, wallet_ownership_proof, global_proof, parameter sources, proof schema, proof constraints, verifier output, and historical artifacts.
3. Disclose review scope, sampling methods, limitations, conflicts of interest, and unresolved issues.
4. If only summary or user inclusion is verified, conclusions cannot be stated as Stage 1-level PoR.

### 9.3 Recommendations for Tool and Infrastructure Builders

The industry needs shared verification tools and data infrastructure, not each exchange defining non-reusable formats:

1. Establish standardized PoR schema covering snapshot, wallet_address_list, wallet_ownership_proof, global_proof, vk/config, setup transcript, proof constraints, scope, and exclusion.
2. Build open-source verifiers, batch signature verification tools, proof replay tools, Web/WASM one-click verification components, and reproducible test sets.
3. Provide independent archiving and mirroring services so proof/root/commitment, wallet_address_list, and review reports do not depend on a single exchange server.
4. Promote on-chain anchoring, DA, IPFS/Arweave/object storage hash indexing, and other long-term traceable publication layers.
5. External display should show Stage, Gen, E Level, scope, snapshot, and key gaps—not reserve ratio alone.

### 9.4 Industry Migration Path

The industry can advance along the following path:


| Phase | Goal | Key Actions |
| ----- | ---- | ----------- |
| From nothing to something | Pre-Stage -> Stage 0 | Establish ongoing snapshot, summary, user proof, and basic verifier |
| From self-report to reproducible | Stage 0 -> Stage 1 | Publish wallet_address_list, wallet_ownership_proof, global_proof, parameter sources, and review scope |
| From disclosure to trust minimization | Stage 1 -> Stage 2 | Introduce low-barrier inclusion proof, canonical anchoring, DA, permissionless verification, key security property self-attestation, and business-consistent constraints |
| From single point to ecosystem collaboration | Single exchange disclosure -> industry public verification infrastructure | Standard schema, open-source verifiers, independent archiving, cross-institution re-verification |

# ardmere Service Audience

`ardmere` is not just a PoR verification tool for retail users. It is PoR transparency infrastructure for regulators, institutional allocators, independent reviewers, and exchanges.

One-line positioning:

> ardmere helps external observers determine whether an exchange PoR is genuinely verifiable, whether it meets an effective PoR standard, and which key evidence is missing.

---

## 1. Core audiences

### 1.1 Regulators / policymakers

Regulators are one of the most important audiences.

They do not need to run proofs themselves. They need answers to:

1. Which exchanges have effective PoR?
2. Which are only Stage 0 and must not be treated as effective PoR?
3. Which key artifacts are missing?
4. Are `wallet_address_list`, `wallet_ownership_proof`, and `global_proof` published?
5. Does the proof system rely on an opaque trusted setup?
6. Are there risks such as negative-net-worth fake users, negative balances offsetting real liabilities, off-chain distribution, or replaceable history?
7. How should PoR requirements be written into rules or licensing conditions?

Value for regulators:

- An actionable minimum PoR standard.
- Clear separation of `Pre-Stage`, `Stage 0`, `Stage 1`, and `Stage 2`.
- Explicit guidance that Stage 0 must not be accepted as effective PoR.
- A clear boundary between PoR and traditional audit, supporting a dual-track regulatory model.

### 1.2 Institutional investors / market makers / custody clients

Institutional allocators care about exchange transparency and custody risk.

They need:

1. Cross-exchange comparison of PoR transparency.
2. Fact-checking of claims such as “100% backed”, “audited”, or “ZK verified”.
3. Judgment on whether an exchange is Stage 0, Stage 1, or a Stage 2 candidate.
4. Integration of PoR risk signals into exchange onboarding, limits, market-making risk, and custody risk models.
5. Tracking of whether PoR is published continuously, historically archived, and whether standards have regressed.

Value for institutions:

- Exchange PoR risk ratings and transparency intelligence.
- Traceable evidence links, gap lists, and risk flags.
- Translation of technical verification into institution-readable ratings and reports.

### 1.3 Independent reviewers / security firms / accounting firms

Third-party reviewers can use ardmere as PoR review infrastructure.

They need:

1. Automated checks that PoR artifacts are complete.
2. Re-run proof / verifier workflows.
3. Review of `wallet_ownership_proof`, `global_proof`, `vk/config`, and trusted setup transcripts.
4. Checklists for AUP / technical verification / limited assurance work.
5. Explicit marking of unverifiable scope and evidence gaps.

Value for reviewers:

- A tooling layer for PoR technical assessment.
- Lower cost of rebuilding verifiers, fetchers, archives, and reconciliation tools.
- Protection against repeating exchange summaries while ignoring key artifacts.

### 1.4 Exchanges

Exchanges are also a potential audience, but independence must be preserved.

They need:

1. A clear view of what is missing for Stage 1 / Stage 2.
2. Pre-flight checks that PoR artifacts are complete.
3. Validation of proof, verifier, disclosure format, historical archive, and user verification UX.
4. A transparency improvement roadmap for regulators, users, and institutional clients.

Value for exchanges:

- PoR readiness assessment.
- Stage upgrade gap analysis.
- Artifact schemas, verifier output, and disclosure checklists.

If serving exchanges, clearly separate:

- independent public rating
- paid technical assessment
- remediation consulting

Avoid “grading your own homework” trust conflicts.

### 1.5 Advanced users / researchers / media

These users rarely run proofs daily, but they cite ratings and reports.

They need:

1. A concise, credible exchange transparency ranking.
2. Traceable evidence links.
3. Fact-checking of exchange marketing language.
4. Readable explanations of PoR technical differences, review scope, and risk gaps.

Value for them:

- Public transparency reports.
- Exchange PoR evidence indexes.
- Citable Stage conclusions and risk explanations.

---

## 2. Who are the first-phase direct customers?

Retail users are ultimate beneficiaries, but not necessarily first-phase direct customers.

First-phase priority audiences:

1. **Regulators and policy researchers**: need to define the minimum standard for effective PoR.
2. **Institutional allocators**: need to integrate PoR transparency into risk controls.
3. **Independent reviewers**: need a toolized, standardized PoR review workflow.

Exchanges can be a second-phase audience, with clear independence boundaries.

---

## 3. Product surfaces

### 3.1 Public Transparency Dashboard

For the public, researchers, media, and institutional screening.

Core content:

- Exchange PoR Stage.
- Technology Generation (Gen).
- Evidence Level (E).
- Missing artifacts.
- Risk flags.
- Historical snapshots.
- Last verification time.

### 3.2 Regulator / Institution Report

For regulators and institutional clients.

Core content:

- Stage conclusion per exchange.
- Whether the minimum effective PoR standard is met.
- Stage 0 / Stage 1 / Stage 2 blocked reasons.
- Detailed checks on proof, wallet, setup, DA, frequency, and audit taxonomy.
- How PoR complements traditional audit.

### 3.3 Verification API

For reviewers, institutional clients, and internal systems.

Core capabilities:

- Input PoR artifacts.
- Output verifier results.
- Output missing artifacts.
- Output risk flags.
- Output suggested Stage / E Level.
- Output machine-readable reports.

### 3.4 Exchange Readiness Assessment

For exchanges, with independence preserved.

Core content:

- Current Stage.
- Gaps to reach Stage 1 / Stage 2.
- Artifact schema recommendations.
- Inclusion proof user verification UX recommendations.
- Improvement roadmap for trusted setup, wallet ownership, DA, anchoring, and frequency.

---

## 4. Independence principles

To build long-term credibility, ardmere should enforce these boundaries:

1. Public ratings cannot be purchased by exchanges.
2. Paid technical assessments must disclose scope and conflicts of interest.
3. Remediation consulting must not be equated with independent certification.
4. Every rating must bind to artifacts, hashes, URLs, verification output, and timestamps.
5. Missing data must be marked `UNVERIFIABLE`; it must not default to pass.

---

## 5. Recommended rollout order

### Phase 1: Public Methodology + First Reports

Goal: establish credible methodology and sample reports.

Suggested subjects:

- OKX: Stage 1 reference sample.
- Binance: Gen 2 / E2 but Stage 1 blocked reference sample.
- Bybit or Bitget: Stage 0 / Merkle inclusion sample.

Outputs:

- public dashboard v0
- exchange transparency report v0
- methodology page

### Phase 2: Institution / Regulator Package

Goal: translate ratings into regulatory and institutional risk language.

Outputs:

- exchange PoR risk matrix
- minimum effective PoR standard brief
- policy guidance that Stage 0 must not count as effective PoR
- PoR and traditional audit complementarity brief

### Phase 3: API and Tooling

Goal: productize verification services.

Outputs:

- artifacts upload / fetch API
- verification API
- report API
- rule engine
- historical archive

### Phase 4: Exchange Readiness / Remediation

Goal: help exchanges improve PoR without sacrificing independence.

Outputs:

- readiness assessment
- Stage upgrade checklist
- artifacts schema
- user verification UX checklist

---

## 6. Conclusion

ardmere’s direct customers should not be defined first as retail traders. They should be defined as:

> regulators, institutional allocators, and independent reviewers.

Retail users are ultimate beneficiaries; exchanges are subjects of assessment and later service recipients.

The most robust product path is:

1. Build public, independent, traceable PoR Stage ratings first.
2. Serve regulators and institutional clients next.
3. Serve exchange improvement last, under explicit independence boundaries.

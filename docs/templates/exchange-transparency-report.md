# {{exchange}} PoR Transparency Report

> Template version: `ardmere/exchange-transparency-report@1`
> Methodology: `docs/por-transparency-framework.md`
> Assessment JSON: `{{assessment_json_path}}`

## 1. Summary

| Field | Value |
| --- | --- |
| Exchange | `{{exchange}}` |
| Snapshot | `{{snapshot_id}}` |
| Period | `{{period_seq}}` |
| Snapshot Time | `{{snapshot_time}}` |
| PoR Stage | `{{por_stage}} — {{stage_name}}` |
| Gen / Evidence | `{{gen}} / {{evidence_level}}` |
| Confidence | `{{confidence}}` |
| Effective PoR | `{{effective_por}}` |

{{headline}}

{{rationale}}

Under the ardmere framework, this snapshot should **not** be described as `{{forbidden_stage_or_claims}}` unless the corresponding evidence is public, reproducible, and listed in this report.

## 2. Stage Decision

### {{por_stage}}, {{effective_por_summary}}

{{stage_rationale}}

{{stage_block_intro}}

| Missing / Blocked Evidence | Risk Flag | Stage Effect | Why It Matters |
| --- | --- | --- | --- |
| `{{missing_or_blocked_evidence}}` | `{{risk_flag}}` | {{stage_effect}} | {{why_it_matters}} |

{{stage2_note}}

## 3. Frequency and Freshness

| Field | Value |
| --- | --- |
| Latest Snapshot | `{{latest_snapshot_time}}` |
| Previous Snapshot | `{{previous_snapshot_time}}` |
| Observed Cadence | `{{observed_cadence}}` |
| History Available | `{{history_available}}` |
| Event-triggered Updates | `{{event_triggered_updates}}` |
| Daily Root / Commitment Anchor | `{{daily_anchor}}` |
| Stage Impact | {{frequency_stage_impact}} |

{{frequency_notes}}

## 4. Evidence and Bundles

### Public Artifacts

| Artifact | SHA-256 | Source |
| --- | --- | --- |
| `{{artifact_kind}}` | `{{sha256}}` | `{{source_url_or_path}}` |

### Bundle References

| Bundle | Value |
| --- | --- |
| Artifact bundle root | `{{artifact_bundle_root}}` |
| Verification bundle root | `{{verification_bundle_root}}` |
| Artifact bundle SHA-256 | `{{artifact_bundle_sha256}}` |
| Verification bundle SHA-256 | `{{verification_bundle_sha256}}` |

## 5. Verifier Evidence

| Verifier | Version | Verdict | Coverage | Summary |
| --- | --- | --- | --- | --- |
| `{{verifier_id}}` | `{{version}}` | `{{verdict}}` | `{{coverage}}` | {{verifier_summary}} |

{{verifier_interpretation}}

## 6. Minimal Checklist

| Requirement | Status | Notes |
| --- | --- | --- |
| {{requirement}} | {{status}} | {{notes}} |

## 7. Recommendations

| Priority | Recommendation | Related Risk |
| --- | --- | --- |
| `{{priority}}` | {{recommendation}} | `{{related_risk}}` |

## 8. Boundary

This report evaluates public PoR artifacts and reproducibility. It does not evaluate the exchange's overall financial health, corporate governance, internal controls, off-chain assets, legal compliance, or complete off-balance-sheet liabilities. Traditional audit and PoR should be treated as complementary evidence tracks.

## 9. References

- Assessment JSON: `{{assessment_json_path}}`
- Artifact bundle: `{{artifact_bundle_path}}`
- Verification bundle: `{{verification_bundle_path}}`
- Supporting report: `{{supporting_report_path}}`

Missing data must be marked as `UNVERIFIABLE`, not treated as `PASS`.

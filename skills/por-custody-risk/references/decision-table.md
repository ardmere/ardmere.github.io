# Custody risk score

The skill scores custody risk. It does not decide whether to trade. There is no `decision` field and no allow / reduce / block output. `por-check.sh` reads only this table, so a second run has to land on the same integers.

`riskScore` is an integer from 0 to 100. Sum the applied `factors[].delta` values. If that sum is above 100, reduce deltas from the last factor backwards until the sum is 100, and say so in the note. `riskScore` is that sum.

`riskLevel` starts from the score, then a floor can raise it:

| Score | Level |
| --- | --- |
| 0–24 | `low` |
| 25–49 | `medium` |
| 50–79 | `high` |
| 80–100 | `critical` |

Floors, applied after the band: a blocking flag raises the level to at least `high`. Stage below effective PoR raises the level to at least `medium`. A file-level failure is `critical`.

`conclusive` is a separate flag. It is `true` whenever the table produced a score, and `false` only when the anchor record could not be read (see below), in which case `riskScore` is `null` rather than `0`. A `0` would be read as "low risk", which is a claim this skill is not entitled to make when it could not see the chain.

## Order

1. The file: missing, unreadable, or read. Stop on the first two.
2. The anchored record, when it was read and holds this snapshot: compare roots. Stop on a disagreement.
3. An anchor read that failed: inconclusive, no score.
4. The contract answered not-found: `ANCHOR_MISSING`, not a reserve judgement.
5. Blocking flags, in `riskFlags` order.
6. Stage factor, if the stage is below effective PoR.
7. Non-blocking flags, in `riskFlags` order, capped at 20.

## File failures

One factor, delta **100**, level `critical`, `conclusive: true`. Flag points are not added. The factor id is the first failure in this order:

| Id | When |
| --- | --- |
| `FILE_MISSING` | The assessment path does not exist |
| `FILE_UNREADABLE` | The assessment file is not valid JSON |

`FILE_MISSING` matters more than its name suggests. If the contract *does* hold an anchored record for the snapshotId you asked about, the note says so: the anchor exists and the local bytes to compare it against do not. That is a different statement from "there is no anchor", and it is the case the anchor read exists to distinguish.

## The anchored record

`scripts/por-check.sh` reads `getSnapshot(keccak256(auditId))` from the `ArdmerePoRAnchor` proxy ([`docs/anchor-query-api.md`](../../../docs/anchor-query-api.md)). This is a live `eth_call`, not a replay. With `--anchor-fixture` it replays a recorded answer instead, and every such output carries `anchor.live: false`.

Three fields are compared, each only when both sides carry a value: `artifactBundleRoot`, `verificationBundleRoot`, and `snapshot.exchangeMerkleRoot`. A zero `bytes32` on chain (the contract's "unknown") and an empty string in the assessment are not compared — the output marks them `compared: false`, `matches: null`.

| Id | When | Effect |
| --- | --- | --- |
| `ANCHOR_MISMATCH` | The contract holds this snapshot and a compared root disagrees | One factor, delta **100**, `critical` |
| `ANCHOR_UNAVAILABLE` | No RPC endpoint answered | `conclusive: false`, `riskScore: null` |
| `ANCHOR_MISSING` | The contract answered `SnapshotNotFound` for this snapshotId | One factor, delta **100**, plus `notAReserveJudgement: true` |

`ANCHOR_MISMATCH` means *these bytes are not the bytes that were anchored*. Anything that regenerates a bundle produces the same output as tampering does, so read it as "this file is not the anchored bundle", not as proof of misconduct.

It fires on real files. Of the 23 assessments in `docs/reports/` whose snapshotId has a v3 record on Sepolia, 12 disagree with the contract on `verificationBundleRoot` and 4 on `artifactBundleRoot`. Both numbers come from reading the live contract on 2026-11-07, not from a recorded fixture. A plausible explanation is that the verification bundle was rebuilt after the anchor transaction was sent, but **that is a hypothesis, not a finding** — nothing in this table can tell the two apart, and neither can `por verify` on its own. What the score can say is that the file on disk and the contract attest different bytes. Re-run the verifier and re-anchor to clear it; never edit the assessment to match the chain.

That rate is also why the mismatch is a short-circuit rather than one factor among many. A file whose bundle root is not the anchored one has no corroborated verification result, so adding flag points to it would decorate a number rather than measure anything.

`ANCHOR_MISSING` is **not** "reserves are bad". The live contract answering not-found means ardmere has not anchored this snapshot yet. A brand-new assessment reaches this state before its anchor transaction lands. Report it as "not anchored yet" and stop.

`ANCHOR_UNAVAILABLE` is **not** a mismatch either. If the RPC is down, the block explorer is rate-limiting, or `POR_RPC_URL` points at a chain without the proxy, nothing can be concluded. The reference skill scored an unavailable view as `critical`; that reads as a risk judgement about the exchange, which the evidence does not support. This skill returns `riskScore: null` instead, and exits non-zero so a caller cannot mistake the run for a score of zero.

## Blocking flags

These ids are the public framework's high flags that block a stage. Each occurrence in `riskFlags` adds **40**.

| Id | Delta |
| --- | --- |
| `UNVERIFIABLE` | 40 |
| `LOGIN_WALL` | 40 |
| `SAMPLE_ONLY` | 40 |
| `BLACK_BOX_INCLUSION_PROOF` | 40 |

`BLACK_BOX_INCLUSION_PROOF` and `LOGIN_WALL` describe the exchange's artifact, not the user's own file. Verifying a user's inclusion proof locally does not remove either flag: the flag says the public cannot check the tree, and one user's private proof does not change that. Report a pass as "your leaf was in the published tree", never as "this flag is resolved".

## Stage below effective PoR

If `porStage` is `Stage 0` or `Pre-Stage` (the exact value, or that value followed by a space, hyphen, or em dash), add one factor:

| Id | Delta |
| --- | --- |
| `STAGE_BELOW_EFFECTIVE_POR` | 35 |

## Non-blocking flags

Each remaining `riskFlags` entry adds the delta below. The sum of these deltas is capped at **20**. That cap is the base 0–20 for a file with `effectivePoR` true and no blocking flag. Points come from the flag id. Severity is only the fallback for an id that is not in this table.

| Id | Delta |
| --- | --- |
| `STAGE_BLOCKED` | 8 |
| `NO_WALLET_OWNERSHIP_PROOF` | 8 |
| `OPAQUE_TRUSTED_SETUP` | 8 |
| `BUSINESS_CONSTRAINT_GAP` | 8 |
| `NEGATIVE_NET_WORTH_DUMMY_USER_RISK` | 8 |
| `IMPLEMENTATION_SEMANTICS_RISK` | 8 |
| `HIGH_FRICTION_INCLUSION_PROOF` | 5 |
| `NO_INDEPENDENT_REVIEW` | 5 |
| `NO_CANONICAL_ANCHOR` | 5 |
| `CONFLICT_OF_INTEREST_UNDISCLOSED` | 5 |
| `HIGH_FREQUENCY_GAP` | 5 |

Fallback when the id is not listed above and is not blocking:

| Severity on the flag | Delta |
| --- | --- |
| `critical` | 10 |
| `high` | 8 |
| `medium` | 5 |
| `low` | 2 |
| `info` | 1 |
| missing or anything else | 5 |

The note on a flag factor is the assessment's `message` when it has one.

## Scores of files in this repository

Computed with `skills/por-custody-risk/scripts/por-check.sh`, not copied from anywhere. `sha` is the first 16 hex characters of the file's sha256, so you can confirm you are looking at the same bytes.

| File | sha256 (prefix) | Anchor | Flags applied | Score | Level |
| --- | --- | --- | --- | --- | --- |
| `docs/reports/okx/508756245-assessment.json` | `708b8b8aa6a8178d` | not anchored (`ANCHOR_MISSING` when read against a live RPC; the row below is the file-only score) | `NO_CANONICAL_ANCHOR` +5, `HIGH_FREQUENCY_GAP` +5. Stage 1, `effectivePoR` true | 10 | `low` |
| `docs/reports/kucoin/20261031-assessment.json` | `53290a7f5021032f` | anchored, 2 of 3 roots compared and agree | `UNVERIFIABLE` +40, `STAGE_BELOW_EFFECTIVE_POR` +35. Stage 0, `effectivePoR` false | 75 | `high` |
| `docs/reports/bitget/20261021-assessment.json` | `f33c596ad188d7b3` | anchored, 3 of 3 roots agree | `UNVERIFIABLE` +40, `STAGE_BELOW_EFFECTIVE_POR` +35. Stage 0 | 75 | `high` |
| `docs/reports/mexc/202610-assessment.json` | `a8631544e337795d` | anchored, 2 of 3 roots compared and agree | `UNVERIFIABLE` +40, stage +35, then `NO_PUBLIC_MERKLE_ROOT` +8, `WAF_GATED` +8, `LOGIN_GATED_USER_PROOF` +4 (capped at 20) | 95 | `critical` |
| `docs/reports/binance/PR01JUN26-assessment.json` | `250a2a76d87e9956` | **mismatch** when read against the live anchor: `verificationBundleRoot` and `artifactBundleRoot` disagree with the contract | `UNVERIFIABLE` +40, stage +35, `NO_WALLET_OWNERSHIP_PROOF` +8, `OPAQUE_TRUSTED_SETUP` +8, so 91 file-only; `ANCHOR_MISMATCH` +100 when the anchor is read | 91 / 100 | `critical` |
| `docs/reports/okx/506872725-assessment.json` | `8cbf3abbd6c4ae44` | **mismatch** when read against the live anchor: `verificationBundleRoot` disagrees | otherwise 10 (`NO_CANONICAL_ANCHOR` +5, `HIGH_FREQUENCY_GAP` +5); `ANCHOR_MISMATCH` +100 when the anchor is read | 10 / 100 | `low` / `critical` |

The last two rows are the reason the anchor read exists at all. Read file-only they look like ordinary Stage 0 / Stage 1 snapshots. Read against the chain, the bytes on disk are not the bytes that were anchored, and the score says so before any flag arithmetic runs.

Two rows show the same `75` score for different disclosures: KuCoin publishes no wallet list at all, Bitget publishes a Merkle root but keeps the tree behind a login. The score is a summary of the flagged evidence, not a ranking — read `factors[]` before quoting a number.

## What this table does not decide

* Whether to trade, deposit, or withdraw on an exchange. There is no `decision` field for that reason.
* Whether reserves are sufficient. That needs the wallet bundle and the full verifier, not this table.
* Whether a user's own inclusion proof is valid. That is `por-inclusion-proof`, which never adds or removes a `factor` here. A failing proof does not raise `riskScore`; it means the proof the user holds does not match the root it names.
* Whether an exchange is dishonest when a root disagrees. `ANCHOR_MISMATCH` says the file is not the anchored bundle. Re-run the verifier and re-anchor; investigate before accusing.


## On-chain coverage fields (schema v4)

Prefer `walletCoverageBps` and `reserveCoverageBps` from the anchor record when present. Deprecated `coverageBps` is an opaque max(verifier Coverage)×10000 and must not be treated as either ratio. v3 records decode with the new fields as `0` / `false` and `schemaVersion=3`.

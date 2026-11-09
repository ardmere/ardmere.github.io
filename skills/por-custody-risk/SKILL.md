---
name: por-custody-risk
description: Score PoR transparency risk (0–100, higher = worse disclosure posture) from a public assessment JSON and its ArdmerePoRAnchor record on Ethereum Sepolia. Use when an agent needs custody/transparency risk integers, or whether the assessment on disk matches the anchored bundle. Not reserve adequacy, solvency, or trading advice.
---

# PoR custody risk (transparency risk score)

This skill scores **PoR transparency risk** from evidence that is already public, and checks that evidence against a chain. **`riskScore` is not “how safe the exchange is.”** It measures **missing or unverifiable disclosure** (artifacts, proofs, anchors, verifier gaps). **Higher score = higher transparency risk** (worse posture in this table). Lower is better. It does not decide whether to trade.

Run the script and read the JSON. Use `scoreSemantics`, `riskScore`, `riskLevel`, `conclusive`, and `factors`. There is no `decision` field. Do not invent allow, reduce, or block. When presenting results to humans, always label the number as **transparency risk (higher = worse)**. The points live in [`references/decision-table.md`](references/decision-table.md); a second run has to land on the same integers.

## What this does NOT decide

* **Not whether to trade, deposit, or withdraw.** No `decision` field exists, on purpose.
* **Not whether reserves are sufficient.** That needs the wallet bundle and the full verifier, not this table.
* **Not whether an exchange is honest** when a root disagrees. `ANCHOR_MISMATCH` says the file on disk is not the anchored bundle.
* **Not whether a user's own inclusion proof is valid.** That is the [`por-inclusion-proof`](../por-inclusion-proof/SKILL.md) skill, which never adds or removes a `factor` here.
* **Not a verdict about reserves** when the contract answers not-found. `ANCHOR_MISSING` means "not anchored yet".
* **Not a score at all** when the chain could not be read. `conclusive: false`, `riskScore: null`.

## Read the record from Sepolia

The score is bound to the on-chain record, not to a recording. The read is a `getSnapshot` `eth_call` against the `ArdmerePoRAnchor` proxy (schema v5; v3/v4 rows remain readable):

```bash
skills/por-custody-risk/scripts/por-check.sh \
  docs/reports/kucoin/20261031-assessment.json
```

The script takes the audit id and exchange name from the file, computes the v6 composite `snapshotId = keccak256(abi.encodePacked(exchange, 0x00, auditId))` (and falls back to legacy `keccak256(bytes(auditId))` for v3–v5 rows), and calls the proxy at the address in `skills/lib/por_anchor.py` (`DEFAULT_CONTRACT`), on Sepolia (the chain id is in `docs/deployments.md`).
Both are the values in [`docs/deployments.md`](../../docs/deployments.md); override them with the environment variables below rather than editing the module. See [`docs/anchor-query-api.md`](../../docs/anchor-query-api.md) for the ABI, the `exchangeTag`, and the error codes.

| Environment variable | Meaning |
| --- | --- |
| `POR_ANCHOR_CONTRACT` | Override the proxy address (a fork, a redeploy, a mock) |
| `POR_RPC_URL` | Override the RPC endpoint. Unset means the public Sepolia endpoints in `skills/lib/por_anchor.py`, tried in order |

Read the `anchor` block to see where the record came from:

| Field | Meaning |
| --- | --- |
| `anchor.live` | `true` only when an RPC answered in this run |
| `anchor.source` | `rpc` for a live read, the fixture path for a replay, `none` when no read was attempted |
| `anchor.rpc` | the endpoint that answered |
| `anchor.status` | `found`, `not_found`, `unavailable`, or `skipped` |
| `anchor.contract` | the proxy the call was addressed to |
| `anchor.periodSeq`, `anchor.anchoredAt`, `anchor.schemaVersion` | what the contract says, not what the file says |
| `anchor.roots[]` | per-field comparison: `field`, `local`, `onchain`, `compared`, `matches` |
| `anchor.interpretation` | one sentence on what the status means here |

Three statuses each mean something different, and conflating them is the main way to get this wrong:

| Status | Meaning | Report it as |
| --- | --- | --- |
| `found`, roots agree | the bytes on disk are the bundle the contract attests | the score below is corroborated |
| `found`, a root disagrees | these bytes are not the anchored bytes | `ANCHOR_MISMATCH`, 100 / `critical`, and say which field |
| `not_found` | the contract holds no record for this snapshotId | "not anchored yet", never "reserves are bad" |
| `unavailable` | no endpoint answered | "cannot conclude", never a score of 0 |

Anchored-but-no-local-file and local-file-but-unanchored are the two asymmetries worth spelling out:

* **Anchored, no local file** — `FILE_MISSING` with a note saying the contract does hold a record. The anchor exists; the bytes to compare it against do not.
* **Local file, not anchored** — `ANCHOR_MISSING` with `notAReserveJudgement: true`. A fresh assessment is in this state until its anchor transaction lands.

## Offline and replay runs

`--anchor-fixture fixtures/anchor/sepolia-v3-records.json` replays a recorded `eth_call` answer. Every such output has `anchor.live: false` and says it is a replay. Say so if you quote a score from it: a replay is evidence that the code path works, not evidence about the chain now.

`--no-anchor` skips the read entirely and scores the file's own disclosure. The output marks `anchor.status: skipped` and the interpretation says the score describes the file alone. Use it when the assessment is not anchored yet, and never present it as a corroborated score.

```bash
# offline, recorded read
python3 skills/por-custody-risk/scripts/risk_score.py \
  --anchor-fixture fixtures/anchor/sepolia-v3-records.json \
  docs/reports/kucoin/20261031-assessment.json

# file only, no chain claim
python3 skills/por-custody-risk/scripts/risk_score.py --no-anchor \
  docs/reports/okx/508756245-assessment.json
```

## Example run

Quoted from a run against the recorded reads in `fixtures/anchor/sepolia-v3-records.json`. This is a replay, so `anchor.live` is `false`:

```bash
$ python3 skills/por-custody-risk/scripts/risk_score.py \
    --anchor-fixture fixtures/anchor/sepolia-v3-records.json \
    docs/reports/kucoin/20261031-assessment.json
{
  "exchange": "kucoin",
  "snapshot": "20261031",
  "riskScore": 99,
  "riskLevel": "critical",
  "conclusive": true,
  "factors": [
    {
      "id": "UNVERIFIABLE",
      "delta": 40,
      "note": "No public wallet_address_list."
    },
    {
      "id": "STAGE_BELOW_EFFECTIVE_POR",
      "delta": 35,
      "note": "porStage is Stage 0, below effective PoR."
    },
    {
      "id": "NO_PUBLIC_WALLET_ADDRESS_LIST",
      "delta": 16,
      "note": "No public wallet_address_list."
    },
    {
      "id": "NO_WALLET_OWNERSHIP_PROOF",
      "delta": 8,
      "note": "No public download channel for wallet ownership signatures / proofs"
    }
  ],
  "porStage": "Stage 0",
  "effectivePoR": false,
  "riskFlags": ["UNVERIFIABLE"],
  "assessmentSha256": "53290a7f5021032f580e538ebe1c3fa1ced0d2503a65984da93c6b53a0b0ef79",
  "anchor": {
    "status": "found",
    "live": false,
    "source": "fixtures/anchor/sepolia-v3-records.json",
    "periodSeq": 1,
    "schemaVersion": 3,
    "interpretation": "anchored; 2 of 3 roots the assessment declares were compared against the on-chain record and agree"
  }
}
```

Abridged: the full output also carries `anchor.contract`, `anchor.rpc`, the three bundle roots, and the `anchor.roots[]` rows. `2 of 3` is honest — the third field is `exchangeMerkleRoot`, which this exchange leaves empty, so there is nothing to compare.

The same file with `--no-anchor` scores identically, but the interpretation says the score describes the file alone. The score is the same because the anchor agreed; when it disagrees, the anchor wins and the score becomes 100.

## Output fields

| Field | Meaning |
| --- | --- |
| `scoreSemantics` | Always present: `kind: transparency_risk`, `scale` explains higher = worse |
| `riskScore` | 0–100 **transparency risk** (higher = worse), or `null` when inconclusive |
| `riskLevel` | `low`, `medium`, `high`, `critical`, or `null` when inconclusive |
| `conclusive` | `false` only when the anchor could not be read |
| `factors[]` | `id`, `delta`, `note` — the arithmetic, in application order |
| `notAReserveJudgement` | present and `true` on `ANCHOR_MISSING` only |
| `porStage`, `effectivePoR`, `riskFlags` | echoed from the assessment |
| `anchor.porStage`, `anchor.walletCoverageBps`, `anchor.reserveCoverageBps` | schema v4+ on-chain fields (v3 rows decode as 0/false); prefer these over deprecated `coverageBps` |
| `anchor.custodyRiskScore` | schema v5 field (`0–100` or `255`); cross-check via `custodyRiskScoreMatchesComputed` when recomputing |
| `assessmentSha256` | the file's sha256, so a quoted score can be tied to bytes |

Exit status is `0` when a score was produced and `1` when the run was inconclusive, so a caller cannot silently treat an unreadable chain as a score of zero.

## Hash and manifest

Unlike the reference skill, this repository keeps no `fixtures/digests.json` registry of assessment hashes, so there is no local manifest to disagree with. The anchor record is the integrity check that does exist: `ANCHOR_MISMATCH` fires when a compared root disagrees with the contract. `assessmentSha256` is printed so a reader can tell two files with the same name apart.

## Who publishes

This skill only reads. Anchoring a new snapshot is a separate, owner-signed step:

```bash
go build -o por ./cmd/por
./por anchor -exchange kucoin           # fetch → verify → bundles → anchor calldata
```

`por anchor` does not require a private key and does not broadcast: it prints the `cast send` for the operator (including `custodyRiskScore` from `risk_score.py --no-anchor`), and the transaction is sent from the permissioned anchor signer (`0xf11AcEcBB54bf72Db69f6BaC4f16FC6491cC670F` on Sepolia). Nothing in this skill writes to the chain, holds a key, or asks for one. See [`docs/por-cli.md`](../../docs/por-cli.md) for the flags.

Re-running the verifier and re-anchoring is the fix for `ANCHOR_MISMATCH` — never edit an assessment to match the chain.

## Verifying the user's own inclusion proof

A custody score describes an exchange. It says nothing about whether *this* user's balance was in the snapshot. That is what the user's own proof is for, and it belongs in the other skill:

```bash
skills/por-inclusion-proof/scripts/inclusion_check.py --list-schemes
skills/por-inclusion-proof/scripts/inclusion_check.py \
  --exchange bybit \
  --assessment docs/reports/bybit/2025061709-assessment.json \
  ~/Downloads/myProof.json
```

Read the proof from the user's own machine. Never ask them to paste it into chat, commit it, or upload it anywhere: it contains their account id, nonce, and balances. That skill prints verdicts only and writes no output file.

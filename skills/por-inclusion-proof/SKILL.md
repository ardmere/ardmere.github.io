---
name: por-inclusion-proof
description: Verify a user's own exchange inclusion proof locally, with dependency-free python3, and bind the recomputed root to the snapshot's published exchange merkle root from the ArdmerePoRAnchor record or an assessment file. Use when a person has downloaded their proof from an exchange and wants to know whether their leaf is in the tree that produced the published root. Does not judge the exchange, does not prove solvency, and never uploads the proof.
---

# PoR inclusion proof

This skill verifies one user's proof against one root. It answers a narrow question: **is my leaf in the tree that produced this root?** It does not answer whether the exchange holds enough reserves for everyone, and it does not change the custody score in [`por-custody-risk`](../por-custody-risk/SKILL.md).

Every scheme below is a port of the Go implementation this repository ships, and the tests check the Python against the same fixture roots the Go tests assert. If a scheme could not be reproduced faithfully it would not be listed.

## Privacy contract

The proof file is the one the exchange gives a logged-in user. It contains the account id, the nonce, and the balances. Treat it as the user's data, not as a report:

* **Read it from the user's own machine.** Never ask them to paste it into chat, commit it to a repository, attach it to an issue, or upload it anywhere.
* **The script prints verdicts only.** It never echoes the account id, the nonce, or any balance. `tests/test_inclusion.py` asserts this.
* **It writes no output file.** There is no `--out`, no cache, and nothing is created next to the proof. Redirect the stdout yourself if you need to keep it, on the user's machine only.
* **It reads nothing else on disk** except the proof you name and, if you ask, the assessment JSON you name.

If a user pastes a proof into chat anyway, do not repeat its contents back. Verify the file on their machine instead.

## Run it

```bash
skills/por-inclusion-proof/scripts/inclusion_check.py --list-schemes

skills/por-inclusion-proof/scripts/inclusion_check.py \
  --exchange bybit \
  ~/Downloads/myProof.json
```

`--list-schemes` lists exactly what is implemented and tested:

```text
$ skills/por-inclusion-proof/scripts/inclusion_check.py --list-schemes
bingx
bitget
bitunix
bybit
coinex
mexc
ourbit
toobit
```

All eight are ports of `internal/exchanges/<exchange>/merkle` and `internal/verifier/user_merkle_*.go`. Kraken, Gate.io, OKX, KuCoin, Poloniex, HTX and the other exchanges in this repository have public artifacts but no user proof scheme in Go, so they are deliberately absent rather than stubbed.

## Example run

A run against the repository fixture that the Go test uses:

```bash
$ python3 skills/por-inclusion-proof/scripts/inclusion_check.py \
    --exchange bitget fixtures/bitget/merkel_tree_bg.json
{
  "exchange": "bitget",
  "verdict": "PASS",
  "rootWidthBytes": 8,
  "recomputedRoot": "ca89456bb711c913",
  "detail": "path recomputes to the proof root and the folded balances match",
  "findings": [
    {
      "id": "TRUNCATED_HASH",
      "note": "Bitget truncates SHA-256 to 16 hex chars (64 bits); a pass is weak evidence"
    }
  ],
  "auditId": "202605",
  "scope": "Proves this account's leaf is in the tree that produced the proof root. It does not prove the exchange holds enough reserves for everyone, and it says nothing about any other account."
}
```

Exit status is `0` on `PASS`, `1` on `FAIL`, and `2` for an exchange this skill does not implement.

## Output fields

| Field | Meaning |
| --- | --- |
| `verdict` | `PASS` or `FAIL` |
| `rootWidthBytes` | how wide this scheme's root is: 8 for Bitget, 16 for MEXC and Ourbit, 32 for the rest |
| `recomputedRoot` | the root the proof's own path produces |
| `detail` | one sentence on what matched or did not |
| `findings[]` | `SELF_HASH_MISMATCH`, `PATH_MISMATCH`, `BALANCE_MISMATCH`, `HEIGHT_MISMATCH`, `TRUNCATED_HASH`, `ROOT_WIDTH_MISMATCH`, `ROOT_NOT_BOUND` |
| `binding` | `publishedRoot`, `publishedRoots`, `matches`, `widthsDiffer`, `source`, `anchor`, `note` |
| `auditId` / `proofId` | the snapshot the proof names, when the scheme carries one |
| `scope` | what a pass does and does not mean — quote it |

`binding.matches` is `null` when nothing could be bound. Do not report `null` as a match, and do not treat a `PASS` with `matches: null` as proof of anything beyond the user's own leaf.

## Binding to a published root

The root a proof names means nothing on its own — anyone can build a tree and hand you a proof for it. Bind it to what the snapshot actually published, in one of two ways.

From the on-chain anchor record, which is the stronger source because it is the contract's own answer:

```bash
skills/por-inclusion-proof/scripts/inclusion_check.py \
  --exchange bybit \
  --onchain 2025061709 \
  ~/Downloads/myProof.json
```

This reuses the same reader as the custody skill: `getSnapshot(keccak256(auditId))` against the proxy in `skills/lib/por_anchor.py` (`DEFAULT_CONTRACT`) on Sepolia, overriding with `POR_ANCHOR_CONTRACT` / `POR_RPC_URL`. `--anchor-fixture <file>` replays a recorded read instead, and then `binding.anchor.live` is `false`.

From an assessment file, when the snapshot is not anchored yet or you are working offline:

```bash
skills/por-inclusion-proof/scripts/inclusion_check.py \
  --exchange bitget \
  --assessment docs/reports/bitget/202605-assessment.json \
  ~/Downloads/merkel_tree_bg.json
```

Both flags can be given at once. If the assessment and the chain publish different roots, the output says `publishedRootsDisagree` and does not claim a match either way.

| Binding state | Meaning |
| --- | --- |
| `matches: true` | the proof's root equals the published root |
| `matches: false` | `ROOT_NOT_BOUND`; the verdict becomes `FAIL` |
| `matches: null`, no published root | the snapshot publishes nothing to bind to (`binding.note` says so) |
| `matches: null`, `widthsDiffer: true` | `ROOT_WIDTH_MISMATCH`; the widths differ, so they cannot be compared |
| `matches: null`, `anchor.status: not_found` | the snapshot is not anchored yet |
| `matches: null`, `anchor.status: unavailable` | the chain could not be read; nothing is concluded |

### Root widths

Two schemes truncate SHA-256, so their roots are not 256-bit:

| Scheme | Root bits | Note |
| --- | --- | --- |
| `bitget` | 64 | SHA-256 truncated to 16 hex chars; a pass is weak evidence |
| `mexc` | 128 | SHA-256 truncated to 32 hex chars |
| `ourbit` | 128 | same truncation class as MEXC |
| `bingx`, `bitunix`, `bybit`, `coinex`, `toobit` | 256 | full SHA-256 |

Bitget's assessment publishes a 256-bit `exchangeMerkleRoot` with its 64-bit root right-aligned in it, while its user proofs produce 64-bit roots. The two cannot be compared as strings, so this skill reports `ROOT_WIDTH_MISMATCH` rather than blaming the user's file. Treat such a binding as unknown, not as a failure. MEXC and Ourbit behave the same way: their assessments publish no root at all today, so a MEXC or Ourbit proof verifies but does not bind.

`TRUNCATED_HASH` matters when quoting a result. A Bitget pass rests on 64 bits; a collision is not a practical attack on one account's proof, but the evidence is weaker than a full-width pass and the output says so.

## What a pass proves

A pass proves the user's leaf is in the tree that produced the root the proof names. With `matches: true` it also proves that root is the one the snapshot published.

It does not prove the exchange holds enough reserves for everyone, that other accounts are included, or that the snapshot is honest. It also does not resolve `BLACK_BOX_INCLUSION_PROOF` or `LOGIN_WALL` in the custody score: those flags say the *public* cannot check the tree, and one user's private proof does not change that.

A `FAIL` means the file and the root it names disagree. Causes worth naming to the user before anything else: the proof was downloaded for a different snapshot than the assessment they are binding against, the file was edited or re-saved by a tool, or the exchange rotated its tree after the user downloaded the proof. Only after ruling those out is a mismatch evidence about the exchange.

## Who publishes

This skill only reads. The roots it binds against come from two owner-signed steps:

* **The anchor** — `por anchor` computes the bundle roots and prints the `cast send` for the permissioned signer to broadcast. Nothing here holds a key.
* **The assessment** — `por verify` writes `docs/reports/<exchange>/<snapshotId>-assessment.json`, including `snapshot.exchangeMerkleRoot` when the exchange publishes one.

If a proof will not bind and the snapshot is anchored, the next step is `por verify` on the exchange's published root, not an edit to the assessment.

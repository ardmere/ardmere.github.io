#!/usr/bin/env python3
"""Score exchange reserve risk from a public assessment plus its on-chain anchor record.

Reads one assessment JSON and the ArdmerePoRAnchor record for the same audit id,
then applies the scoring table in references/decision-table.md.

What it will not do:
  * emit a `decision` field (no allow / reduce / block)
  * treat "not anchored yet" as "reserves are bad"
  * invent a score when the chain could not be read

Standard library only. The record read lives in ../lib/por_anchor.py.
"""

import argparse
import hashlib
import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent.parent / "lib"))

BLOCKING = {
    "UNVERIFIABLE": 40,
    "LOGIN_WALL": 40,
    "SAMPLE_ONLY": 40,
    "BLACK_BOX_INCLUSION_PROOF": 40,
}
NONBLOCKING = {
    "STAGE_BLOCKED": 8,
    "NO_WALLET_OWNERSHIP_PROOF": 8,
    "OPAQUE_TRUSTED_SETUP": 8,
    "BUSINESS_CONSTRAINT_GAP": 8,
    "NEGATIVE_NET_WORTH_DUMMY_USER_RISK": 8,
    "IMPLEMENTATION_SEMANTICS_RISK": 8,
    "HIGH_FRICTION_INCLUSION_PROOF": 5,
    "NO_INDEPENDENT_REVIEW": 5,
    "NO_CANONICAL_ANCHOR": 5,
    "CONFLICT_OF_INTEREST_UNDISCLOSED": 5,
    "HIGH_FREQUENCY_GAP": 5,
}
SEVERITY_FALLBACK = {
    "critical": 10,
    "high": 8,
    "medium": 5,
    "low": 2,
    "info": 1,
}
NONBLOCKING_CAP = 20
DISCLOSURE_GAP = 16  # NO_PUBLIC_WALLET_ADDRESS_LIST; not subject to non-blocking cap
SCORE_CAP = 100
LIST_BLOCK_RULE_IDS = frozenset({
    "floor_no_public_wallet_address_list",
    "floor_no_snapshot_wallet_address_list",
    "floor_no_machine_readable_wallet_list",
})
RANK = {"low": 0, "medium": 1, "high": 2, "critical": 3}
ZERO_ROOT = "0" * 64


def band(score):
    if score >= 80:
        return "critical"
    if score >= 50:
        return "high"
    if score >= 25:
        return "medium"
    return "low"


def at_least(level, floor):
    return level if RANK[level] >= RANK[floor] else floor


def normalize_hex(value):
    if not isinstance(value, str):
        return None
    text = value.strip().lower()
    if text.startswith("0x"):
        text = text[2:]
    if len(text) != 64 or any(c not in "0123456789abcdef" for c in text):
        return None
    return text


def stage_below(stage):
    if not isinstance(stage, str):
        return False
    text = stage.strip()
    lowered = text.lower()
    if lowered == "pre-stage" or lowered.startswith(("pre-stage ", "pre-stage-", "pre-stage—")):
        return True
    return text == "Stage 0" or text.startswith(("Stage 0 ", "Stage 0-", "Stage 0—"))


def cap_sum(factors, limit):
    """Trim deltas from the last factor backwards until the total fits."""
    total = sum(item["delta"] for item in factors)
    if total <= limit:
        return
    overflow = total - limit
    for item in reversed(factors):
        if overflow <= 0:
            break
        cut = min(item["delta"], overflow)
        if cut:
            item["delta"] -= cut
            overflow -= cut
            item["note"] = (item["note"] + " Score capped at 100.").strip()


def anchor_view(doc):
    """Summarise an anchor read without copying its record wholesale.

    `None` means no read was attempted at all (the file-only mode); that is
    reported as `skipped`, which is honest about the scope of the score.
    """
    if doc is None:
        return {
            "status": "skipped",
            "live": False,
            "source": "none",
            "contract": None,
            "chainId": None,
            "rpc": None,
            "snapshotId": None,
            "note": "no anchor read was attempted; the score describes the file alone",
        }
    record = doc.get("record") or {}
    return {
        "status": doc.get("status"),
        "live": doc.get("live"),
        "source": doc.get("source"),
        "contract": doc.get("contract"),
        "chainId": doc.get("chainId"),
        "rpc": doc.get("rpc"),
        "snapshotId": doc.get("snapshotId"),
        "exchangeTag": doc.get("exchangeTag"),
        "periodSeq": record.get("periodSeq"),
        "snapshotTime": record.get("snapshotTime"),
        "anchoredAt": record.get("anchoredAt"),
        "schemaVersion": record.get("schemaVersion"),
        "porStage": record.get("porStage"),
        "effectivePoR": bool(record.get("effectivePoR")) if record.get("effectivePoR") is not None else None,
        "walletCoverageBps": record.get("walletCoverageBps"),
        "reserveCoverageBps": record.get("reserveCoverageBps"),
        "coverageBps": record.get("coverageBps"),  # deprecated
        "artifactBundleRoot": record.get("artifactBundleRoot"),
        "verificationBundleRoot": record.get("verificationBundleRoot"),
        "exchangeMerkleRoot": record.get("exchangeMerkleRoot"),
        "note": doc.get("note"),
        "error": doc.get("error"),
    }


def compare_roots(assessment, record):
    """Which of the roots the assessment declares disagree with the chain."""
    artifacts = (assessment.get("artifacts") or {}) if isinstance(assessment, dict) else {}
    snapshot_obj = (assessment.get("snapshot") or {}) if isinstance(assessment, dict) else {}
    pairs = (
        ("artifactBundleRoot", artifacts.get("artifactBundleRoot"), record.get("artifactBundleRoot")),
        ("verificationBundleRoot", artifacts.get("verificationBundleRoot"), record.get("verificationBundleRoot")),
        ("exchangeMerkleRoot", snapshot_obj.get("exchangeMerkleRoot"), record.get("exchangeMerkleRoot")),
    )
    compared = []
    disagreements = []
    for field, local_value, chain_value in pairs:
        chain = normalize_hex(chain_value)
        local = normalize_hex(local_value)
        row = {
            "field": field,
            "local": local,
            "onchain": chain,
            "compared": False,
            "matches": None,
        }
        if chain is None or chain == ZERO_ROOT:
            # The contract holds no value here; nothing to compare against.
            row["note"] = "the contract holds no value for this field"
        elif local is None:
            row["note"] = "the assessment publishes no value for this field"
        else:
            row["compared"] = True
            row["matches"] = local == chain
            if not row["matches"]:
                disagreements.append(field)
        compared.append(row)
    return compared, disagreements


def load_assessment(path):
    """(file_sha256, parsed doc or None)."""
    file_sha = hashlib.sha256(path.read_bytes()).hexdigest()
    try:
        return file_sha, json.loads(path.read_text())
    except json.JSONDecodeError:
        return file_sha, None


def _checklist_item(assessment, item_id):
    for row in assessment.get("checklist") or []:
        if isinstance(row, dict) and row.get("id") == item_id:
            return row
    return None


def has_public_wallet_zip(assessment):
    """True when the assessment's public artifact set includes a wallet zip."""
    if not isinstance(assessment, dict):
        return False
    present = (assessment.get("artifacts") or {}).get("present") or []
    for art in present:
        if not isinstance(art, dict):
            continue
        kind = str(art.get("kind") or "").lower()
        if "walletzip" in kind.replace("_", ""):
            return True
    return False


def derive_disclosure_factors(assessment, flag_ids):
    """Supplemental factors from stageBlockedReasons, checklist, and verifierResults."""
    extra = []
    if not isinstance(assessment, dict):
        return extra

    if not has_public_wallet_zip(assessment):
        list_note = None
        for row in assessment.get("stageBlockedReasons") or []:
            if not isinstance(row, dict):
                continue
            if row.get("ruleId") in LIST_BLOCK_RULE_IDS:
                list_note = row.get("reason") or "No public wallet_address_list."
                break
        if list_note is None:
            s1_1 = _checklist_item(assessment, "S1-1")
            if s1_1 and s1_1.get("status") in ("fail", "unverifiable"):
                list_note = s1_1.get("notes") or "wallet_address_list not public or on-chain verifiable."
        if list_note is not None:
            extra.append({
                "id": "NO_PUBLIC_WALLET_ADDRESS_LIST",
                "delta": DISCLOSURE_GAP,
                "note": list_note,
            })

    if "NO_WALLET_OWNERSHIP_PROOF" not in flag_ids:
        own_note = None
        for v in assessment.get("verifierResults") or []:
            if not isinstance(v, dict):
                continue
            vid = str(v.get("verifierId") or "")
            if not vid.startswith("address-ownership"):
                continue
            if v.get("verdict") in ("FAIL", "UNVERIFIABLE"):
                own_note = v.get("reason") or "address-ownership verifier did not PASS."
                break
        if own_note is None:
            s1_2 = _checklist_item(assessment, "S1-2")
            if s1_2 and s1_2.get("status") in ("fail", "unverifiable"):
                own_note = s1_2.get("notes") or "wallet_ownership_proof not batch-verifiable."
        if own_note is not None:
            extra.append({
                "id": "NO_WALLET_OWNERSHIP_PROOF",
                "delta": NONBLOCKING["NO_WALLET_OWNERSHIP_PROOF"],
                "note": own_note,
            })

    return extra


def assessment_fields(assessment):
    snapshot_obj = assessment.get("snapshot") if isinstance(assessment, dict) else None
    snapshot_obj = snapshot_obj if isinstance(snapshot_obj, dict) else {}
    snapshot = snapshot_obj.get("id")
    raw_flags = assessment.get("riskFlags") if isinstance(assessment, dict) else None
    flag_rows = [f for f in (raw_flags or []) if isinstance(f, dict)]
    return {
        "exchange": assessment.get("exchange") if isinstance(assessment, dict) else None,
        "snapshot": str(snapshot) if snapshot is not None else None,
        "stage": assessment.get("porStage") if isinstance(assessment, dict) else None,
        "effective": assessment.get("effectivePoR") if isinstance(assessment, dict) else None,
        "flag_rows": flag_rows,
        "flag_ids": [f["id"] for f in flag_rows if isinstance(f.get("id"), str)],
    }


def refuse(factor_id, note, fields, anchor, file_sha, extra=None):
    """A file we cannot score at all. One factor, 100, critical."""
    doc = {
        "exchange": fields.get("exchange"),
        "snapshot": fields.get("snapshot"),
        "riskScore": SCORE_CAP,
        "riskLevel": "critical",
        "conclusive": True,
        "factors": [{"id": factor_id, "delta": SCORE_CAP, "note": note}],
        "porStage": fields.get("stage"),
        "effectivePoR": fields.get("effective"),
        "riskFlags": fields.get("flag_ids"),
        "assessmentSha256": file_sha,
        "anchor": anchor,
    }
    if extra:
        doc.update(extra)
    return doc


def inconclusive(factor_id, note, fields, anchor, file_sha, extra=None):
    """Nothing can be concluded. Score is null: a 0 would read as 'low risk'."""
    doc = {
        "exchange": fields.get("exchange"),
        "snapshot": fields.get("snapshot"),
        "riskScore": None,
        "riskLevel": None,
        "conclusive": False,
        "factors": [{"id": factor_id, "delta": None, "note": note}],
        "porStage": fields.get("stage"),
        "effectivePoR": fields.get("effective"),
        "riskFlags": fields.get("flag_ids"),
        "assessmentSha256": file_sha,
        "anchor": anchor,
    }
    if extra:
        doc.update(extra)
    return doc


def score(assessment_path, anchor_doc=None):
    """Apply the decision table. `anchor_doc` is a por_anchor read result."""
    path = Path(assessment_path)
    anchor = anchor_view(anchor_doc)
    file_sha = None
    assessment = None
    fields = {"exchange": None, "snapshot": None, "stage": None, "effective": None,
              "flag_rows": [], "flag_ids": []}

    if path.is_file():
        file_sha, assessment = load_assessment(path)
        if isinstance(assessment, dict):
            fields = assessment_fields(assessment)

    # 1. The file itself.
    if not path.is_file():
        note = "The assessment file does not exist."
        if anchor.get("status") == "found":
            note += (" The contract does hold an anchored record for this snapshotId, so the anchor "
                     "exists but there are no local bytes to compare it against.")
        return refuse("FILE_MISSING", note, fields, anchor, file_sha)
    if assessment is None:
        return refuse("FILE_UNREADABLE",
                      "The assessment file is not valid JSON, so its stage and flags cannot be read.",
                      fields, anchor, file_sha)

    # 2. The anchor decides whether these bytes are the ones that were anchored. Only a
    #    positive contradiction stops the table; absence and unreadability are separate.
    roots, disagreements = [], []
    if anchor.get("status") == "found":
        roots, disagreements = compare_roots(assessment, anchor_doc.get("record") or {})
        anchor["roots"] = roots

    if disagreements:
        anchor["interpretation"] = (
            "the anchored record and the local file disagree, so the bytes scored here are not "
            "byte-identical to what the contract attested"
        )
        return refuse(
            "ANCHOR_MISMATCH",
            "The anchored record disagrees with this assessment on: "
            + ", ".join(disagreements)
            + ". The local file is not the bundle that was anchored. A regenerated assessment "
              "does this too, so read it as 'these bytes are not the anchored bytes', not as proof "
              "of tampering.",
            fields, anchor, file_sha,
        )
    if anchor.get("status") == "unavailable":
        anchor["interpretation"] = "no anchor read, so nothing about this snapshot is corroborated"
        note = anchor.get("error") or "No RPC endpoint answered."
        return inconclusive(
            "ANCHOR_UNAVAILABLE",
            f"The anchor record could not be read, so the file cannot be corroborated either way: {note}",
            fields, anchor, file_sha,
        )
    if anchor.get("status") == "not_found":
        anchor["interpretation"] = (
            "not anchored yet: the live contract answered not-found for this snapshotId. "
            "That is a statement about ardmere's anchor, not about the exchange's reserves."
        )
        return refuse(
            "ANCHOR_MISSING",
            "The live contract holds no record for this snapshotId, so nothing in this file is "
            "anchored yet. This is 'not anchored', not 'reserves are bad'.",
            fields, anchor, file_sha,
            extra={"notAReserveJudgement": True},
        )
    if anchor.get("status") == "skipped":
        anchor["interpretation"] = "no anchor was read; the score below describes the file's own disclosure only"
    elif anchor.get("status") == "found":
        declared = sum(1 for row in roots if row["compared"])
        anchor["interpretation"] = (
            f"anchored; {declared} of {len(roots)} roots the assessment declares were compared "
            "against the on-chain record and agree"
        )

    factors = []
    saw_blocking = False
    for flag in fields["flag_rows"]:
        flag_id = flag.get("id")
        if flag_id not in BLOCKING:
            continue
        saw_blocking = True
        factors.append({
            "id": flag_id,
            "delta": BLOCKING[flag_id],
            "note": flag.get("message") or "Blocking high flag.",
        })
    if stage_below(fields["stage"]):
        factors.append({
            "id": "STAGE_BELOW_EFFECTIVE_POR",
            "delta": 35,
            "note": f"porStage is {fields['stage']}, below effective PoR.",
        })
    derived_ids = set()
    for item in derive_disclosure_factors(assessment, fields["flag_ids"]):
        factors.append(item)
        derived_ids.add(item["id"])
    nonblocking_added = 0
    for flag in fields["flag_rows"]:
        flag_id = flag.get("id")
        if not isinstance(flag_id, str) or flag_id in BLOCKING:
            continue
        if flag_id in derived_ids:
            continue
        full = NONBLOCKING.get(flag_id)
        if full is None:
            full = SEVERITY_FALLBACK.get(str(flag.get("severity") or "").lower(), 5)
        room = NONBLOCKING_CAP - nonblocking_added
        applied = min(full, max(room, 0))
        note = flag.get("message") or "Non-blocking flag."
        if applied < full:
            note = f"{note} Non-blocking contribution capped at 20."
        factors.append({"id": flag_id, "delta": applied, "note": note})
        nonblocking_added += applied

    cap_sum(factors, SCORE_CAP)
    risk_score = sum(item["delta"] for item in factors)
    level = band(risk_score)
    if saw_blocking:
        level = at_least(level, "high")
    if stage_below(fields["stage"]):
        level = at_least(level, "medium")

    return {
        "exchange": fields["exchange"],
        "snapshot": fields["snapshot"],
        "riskScore": risk_score,
        "riskLevel": level,
        "conclusive": True,
        "factors": factors,
        "porStage": fields["stage"],
        "effectivePoR": fields["effective"],
        "riskFlags": fields["flag_ids"],
        "assessmentSha256": file_sha,
        "anchor": anchor,
    }


def main(argv):
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument("assessment", help="docs/reports/<exchange>/<snapshotId>-assessment.json")
    parser.add_argument("--audit-id", help="audit id to query (default: snapshot.id in the file)")
    parser.add_argument("--exchange", help="exchange name for the exchangeTag (default: from the file)")
    parser.add_argument("--anchor-contract", help="override the ArdmerePoRAnchor proxy address")
    parser.add_argument("--rpc-url", help="override the Sepolia RPC endpoint")
    parser.add_argument("--anchor-fixture", help="recorded anchor read to use instead of an RPC call")
    parser.add_argument("--no-anchor", action="store_true",
                        help="skip the chain read and score the file alone (says so in the output)")
    opts = parser.parse_args(argv)

    anchor_doc = None
    if opts.no_anchor:
        anchor_doc = {"status": "skipped", "live": False, "source": "none",
                      "note": "--no-anchor: the file was scored without a chain read"}
    else:
        file_sha, assessment = (None, None)
        path = Path(opts.assessment)
        if path.is_file():
            file_sha, assessment = load_assessment(path)
        audit_id = opts.audit_id
        exchange = opts.exchange
        if isinstance(assessment, dict):
            snapshot_obj = assessment.get("snapshot") or {}
            audit_id = audit_id or (str(snapshot_obj.get("id")) if snapshot_obj.get("id") is not None else None)
            exchange = exchange or assessment.get("exchange")
        if not audit_id:
            anchor_doc = {"status": "unavailable", "live": False, "source": "none",
                          "note": "no audit id could be determined, so no anchor was read"}
        else:
            import por_anchor

            anchor_doc = por_anchor.resolve_record(
                audit_id,
                exchange=exchange,
                contract=opts.anchor_contract,
                rpc=opts.rpc_url,
                fixture=opts.anchor_fixture,
            )

    doc = score(opts.assessment, anchor_doc)
    if "decision" in doc:
        raise SystemExit("refusing to print a decision field")
    json.dump(doc, sys.stdout, indent=2)
    sys.stdout.write("\n")
    return 0 if doc.get("conclusive") else 1


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))

#!/usr/bin/env python3
"""Verify a proof-of-reserves user inclusion proof on the user's own machine.

The proof file is the one the exchange gives a logged-in user (myProof.json,
merkel_tree_bg.json and friends). It carries the account id, nonce and balances,
so this script prints verdicts only: it never echoes those values, never writes
an output file, and never reads anything outside the file you name.

Every scheme below is a port of the Go implementation ardmere ships in
internal/exchanges/<exchange>/merkle and internal/verifier/user_merkle_*.go.
Parity is checked in tests/test_inclusion.py against the same fixtures the Go
tests use.

    inclusion_check.py --list-schemes
    inclusion_check.py --exchange bybit <proof.json>
    inclusion_check.py --exchange bitget --assessment <assessment.json> <proof.json>
"""

import argparse
import hashlib
import json
import sys
from decimal import Decimal, InvalidOperation
from pathlib import Path

SCHEMES = {}
# Root width in bytes each scheme produces, used to compare against a published root.
ROOT_WIDTH = {}

# The published root can also come from the on-chain anchor record.
sys.path.insert(0, str(Path(__file__).resolve().parent.parent.parent / "lib"))


class ProofError(ValueError):
    """The file is not a proof this scheme can read; not a verdict about the user."""


def sha256_hex(text):
    return hashlib.sha256(text.encode()).hexdigest()


def truncated(text, hex_len):
    return sha256_hex(text)[:hex_len]


def dec(value):
    """Parse a balance the way shopspring/decimal does, treating junk as zero."""
    if value is None:
        return Decimal(0)
    text = str(value).strip()
    if not text:
        return Decimal(0)
    try:
        return Decimal(text)
    except InvalidOperation:
        return Decimal(0)


def strict_dec(value):
    """decimal.NewFromString semantics: None when the text is not a number."""
    text = str(value).strip()
    if not text:
        return None
    try:
        return Decimal(text)
    except InvalidOperation:
        return None


def parse_dec_strict(value):
    text = str(value).strip()
    if not text:
        return Decimal(0)
    return Decimal(text)


def plain(value):
    """Format like shopspring Decimal.String(): no exponent, no trailing zeros."""
    amount = dec(value)
    if amount == 0:
        return "0"
    text = format(amount.normalize(), "f")
    if "." in text:
        text = text.rstrip("0").rstrip(".")
    return text


def fixed(value, scale):
    """Render a balance at a fixed scale, like Decimal.StringFixed (never exponent)."""
    quantum = Decimal(1).scaleb(-scale)
    return format(dec(value).quantize(quantum, rounding="ROUND_HALF_UP"), "f")


def round_down(value, scale):
    """Decimal.RoundDown: truncate toward zero at `scale` places."""
    quantum = Decimal(1).scaleb(-scale)
    return (dec(value) / quantum).to_integral_value(rounding="ROUND_DOWN") * quantum


def shopspring_str(value):
    """shopspring Decimal.String(): plain notation, trailing zeros trimmed."""
    amount = dec(value)
    if amount == 0:
        return "0"
    text = format(amount.normalize(), "f")
    if "." in text:
        text = text.rstrip("0").rstrip(".")
    return text or "0"


def shopspring_fixed(value, scale):
    """shopspring Decimal.StringFixed(scale): round half away from zero, keep the scale."""
    quantum = Decimal(1).scaleb(-scale)
    rounded = dec(value).quantize(quantum, rounding="ROUND_HALF_UP")
    return format(rounded, "f")


def scheme(*names):
    def register(fn):
        for name in names:
            SCHEMES[name] = fn
        return fn
    return register


def scheme_width(name, width):
    ROOT_WIDTH[name] = width


def result(exchange, verdict, root, detail, findings, extra=None):
    doc = {
        "exchange": exchange,
        "verdict": verdict,
        "rootWidthBytes": ROOT_WIDTH.get(exchange),
        "recomputedRoot": root,
        "detail": detail,
        "findings": findings,
    }
    if extra:
        doc.update(extra)
    return doc


def normalise_hex(value):
    if not isinstance(value, str):
        return None
    text = value.strip().lower()
    if text.startswith("0x"):
        text = text[2:]
    return text or None


def json_balances(balances):
    """Compact JSON with sorted keys and string values, as encoding/json emits it."""
    keys = sorted(balances)
    parts = [f'"{k}":"{balanced_value(balances[k])}"' for k in keys]
    return "{" + ",".join(parts) + "}"


def balanced_value(value):
    return value if isinstance(value, str) else plain(value)


# --- bitget: SHA-256 truncated to 64 bits, balances folded along the path ------

BITGET_HEX = 16


def bitget_balances(balances):
    keys = sorted(balances)
    parts = [f'"{k}":"{shopspring_str(balances[k])}"' for k in keys]
    return "{" + ",".join(parts) + "}"


def bitget_leaf(encrypt_uid, nonce, balances):
    return truncated(f"{encrypt_uid},{nonce},{bitget_balances(balances)}", BITGET_HEX)


def bitget_inner(left, right, balances, level):
    return truncated(f"{left}{right},{bitget_balances(balances)},{level}", BITGET_HEX)


def bitget_add(a, b):
    out = dict(a)
    for key, value in b.items():
        out[key] = shopspring_str(dec(out.get(key, 0)) + dec(value))
    return out


def bitget_normalise(balances):
    """decimal.NewFromString on every value; anything else rejects the file."""
    out = {}
    for key, value in balances.items():
        parsed = strict_dec(value)
        if parsed is None:
            raise ProofError(f"balance for {key} is not a decimal")
        out[key] = shopspring_str(parsed)
    return out


@scheme("bitget")
def verify_bitget(proof):
    scheme_width("bitget", BITGET_HEX // 2)
    findings = []
    uid = proof.get("encryptUid") or ""
    nonce = proof.get("nonce") or ""
    raw_balances = proof.get("balances") or {}
    path = proof.get("merkelPath") or proof.get("merklePath") or []
    root_node = proof.get("root") or {}
    root = (root_node.get("hash") or "").strip().lower()

    if not uid or not nonce or not raw_balances:
        raise ProofError("proof is missing encryptUid, nonce or balances")
    if not path or not root:
        raise ProofError("proof is missing merkelPath or root.hash")

    balances = bitget_normalise(raw_balances)
    current = bitget_leaf(uid, nonce, balances)
    declared = (proof.get("hash") or "").strip().lower()
    if declared and declared != current:
        return result("bitget", "FAIL", None,
                      "self hash does not match encryptUid+nonce+balances",
                      findings + [{"id": "SELF_HASH_MISMATCH",
                                   "note": "recomputed leaf differs from proof.hash"}])

    running = dict(balances)
    for sibling in path:
        sib_balances = bitget_normalise(sibling.get("balances") or {})
        sib_hash = (sibling.get("hash") or "").strip().lower()
        level = sibling.get("level") or 0
        if not sib_hash or level <= 0:
            raise ProofError("path node is missing hash or level")
        running = bitget_add(running, sib_balances)
        position = str(sibling.get("position") or "").strip().lower()
        if position in ("left", "l"):
            current = bitget_inner(sib_hash, current, running, level)
        else:
            current = bitget_inner(current, sib_hash, running, level)

    if current != root:
        return result("bitget", "FAIL", current, "path does not recompute to root.hash",
                      findings + [{"id": "PATH_MISMATCH", "note": "the path does not reach the root the proof names"}])
    root_balances = bitget_normalise(root_node.get("balances") or {})
    if {k: shopspring_str(v) for k, v in running.items()} != root_balances:
        return result("bitget", "FAIL", current, "folded balances differ from the root node balances",
                      findings + [{"id": "BALANCE_MISMATCH",
                                   "note": "the balances summed along the path do not equal the root balances"}])

    findings.append({"id": "TRUNCATED_HASH",
                     "note": f"Bitget truncates SHA-256 to {BITGET_HEX} hex chars (64 bits); a pass is weak evidence"})
    return result("bitget", "PASS", current,
                  "path recomputes to the proof root and the folded balances match",
                  findings, {"auditId": proof.get("auditId")})


# --- bingx: plain SHA-256 concatenation, no balance fold ----------------------


def bingx_path_entries(path_text):
    if isinstance(path_text, list):
        path_text = ",".join(str(p) for p in path_text)
    for part in str(path_text).split(","):
        part = part.strip()
        if not part:
            continue
        chunks = [c.strip().lower() for c in part.split(":")]
        if len(chunks) != 2:
            raise ProofError(f"path entry must be side:hash or hash:side, got {len(chunks)} chunks")
        first, second = chunks
        if first in ("l", "r"):
            yield first, second
        elif second in ("l", "r"):
            yield second, first
        else:
            raise ProofError("a path entry carries no side marker")


@scheme("bingx")
def verify_bingx(proof):
    scheme_width("bingx", 32)
    leaf = (proof.get("merkleLeaf") or "").strip().lower()
    if not leaf and proof.get("assetLine"):
        leaf = sha256_hex(proof["assetLine"])
    if not leaf:
        raise ProofError("proof is missing merkleLeaf (or assetLine to hash)")
    root = normalise_hex(proof.get("merkleRoot"))
    if not root:
        raise ProofError("proof is missing merkleRoot to compare against")

    current = leaf
    for side, sibling in bingx_path_entries(proof.get("merklePath") or ""):
        current = sha256_hex(sibling + current) if side == "l" else sha256_hex(current + sibling)

    if current != root:
        return result("bingx", "FAIL", current, "path does not recompute to merkleRoot",
                      [{"id": "PATH_MISMATCH", "note": "the path does not reach the root the proof names"}])
    return result("bingx", "PASS", current, "leaf and path recompute to merkleRoot",
                  [], {"auditId": proof.get("auditId")})


# --- bitunix: full SHA-256, height-descending parent, per-asset precision ------

BITUNIX_LEFT, BITUNIX_RIGHT, BITUNIX_ROOT = 1, 2, 3


def bitunix_balances_json(balances):
    keys = sorted(balances)
    parts = [f'"{k}":"{balanced_value(balances[k])}"' for k in keys]
    return "{" + ",".join(parts) + "}"


def bitunix_node_leaf(node):
    content = f"{node['encryptUid']},{node['nonce']},{bitunix_balances_json(node['balances'])}"
    return sha256_hex(content)


def bitunix_parent_leaf(left, right, parent):
    content = f"{left['merkleLeaf']}{right['merkleLeaf']},{bitunix_balances_json(parent['balances'])},{parent['height']}"
    return sha256_hex(content)


def bitunix_to_node(row):
    return {
        "encryptUid": row.get("encryptUid") or "",
        "merkleLeaf": (row.get("merkleLeaf") or "").strip().lower(),
        "nonce": row.get("nonce") or "",
        "balances": dict(row.get("balances") or {}),
        "height": row.get("height") or 0,
        "type": row.get("type") or 0,
        "precisions": dict(row.get("precisions") or {}),
    }


def bitunix_total_balances(left, right, precisions):
    out = {}
    for asset, value in left.items():
        scale = precisions.get(asset, 0)
        total = round_down(parse_dec_strict(value) + parse_dec_strict(right.get(asset, 0)), scale)
        out[asset] = fixed(total, scale)
    return out


def bitunix_parent(sibling, self_node, precisions):
    if sibling["type"] == BITUNIX_LEFT:
        left, right = sibling, self_node
    else:
        left, right = self_node, sibling
    parent = {
        "balances": bitunix_total_balances(left["balances"], right["balances"], precisions),
        "height": left["height"] - 1,
    }
    parent["merkleLeaf"] = bitunix_parent_leaf(left, right, parent)
    return parent


def bitunix_validate_self(self_node):
    if not self_node["balances"] or not self_node["precisions"] or self_node["height"] == 0:
        return False
    return bitunix_node_leaf(self_node) == self_node["merkleLeaf"]


def bitunix_validate_path_node(node):
    if not node["merkleLeaf"] or node["height"] == 0:
        return False
    if not BITUNIX_LEFT <= node["type"] <= BITUNIX_ROOT:
        return False
    if not node["balances"] or any(v == "" for v in node["balances"].values()):
        return False
    return True


@scheme("bitunix")
def verify_bitunix(proof):
    scheme_width("bitunix", 32)
    path = [bitunix_to_node(row) for row in (proof.get("path") or [])]
    self_node = bitunix_to_node(proof.get("self") or {})
    if len(path) < 2:
        raise ProofError("a Bitunix path needs at least two nodes")
    if not self_node["merkleLeaf"] or not self_node["encryptUid"] or not self_node["nonce"]:
        raise ProofError("self node is missing merkleLeaf, encryptUid or nonce")

    root_node = path[-1]
    if not bitunix_validate_self(self_node):
        return result("bitunix", "FAIL", None, "self leaf hash does not match encryptUid+nonce+balances",
                      [{"id": "SELF_HASH_MISMATCH", "note": "recomputed leaf differs from self.merkleLeaf"}])
    if not bitunix_validate_path_node(path[0]):
        raise ProofError("first path node is missing hash, height, type or balances")
    if path[0]["type"] == self_node["type"]:
        raise ProofError("the first path node sits on the same side as self")

    precisions = self_node["precisions"]
    parent = bitunix_parent(path[0], self_node, precisions)
    current = self_node
    for row in path[1:-1]:
        current = parent
        parent = bitunix_parent(row, current, precisions)

    if parent["merkleLeaf"] != root_node["merkleLeaf"]:
        return result("bitunix", "FAIL", parent["merkleLeaf"], "path does not recompute to the root leaf",
                      [{"id": "PATH_MISMATCH", "note": "the path does not reach the root the proof names"}])
    return result("bitunix", "PASS", parent["merkleLeaf"],
                  "leaf and path recompute to the root leaf under full SHA-256", [])


# --- bybit: 40-asset vector, SHA-256, height folded into the preimage ----------

BYBIT_ASSETS = [
    "AGLA", "APT", "ATOM", "AVAX", "MNT", "BTC", "DOGE", "DOT", "DYDX", "A",
    "ETH", "S", "LTC", "POL", "SOL", "USDC", "USDT", "XRP", "SHIB", "OP",
    "SAND", "LDO", "LINK", "COMP", "UNI", "MANA", "CRV", "BLUR", "SUSHI", "PEPE",
    "WLD", "BEAM", "AGI", "FET", "IMX", "RENDER", "GALA", "APEX", "SHRAP", "USDE",
]
BYBIT_SCALE = 8


def bybit_vector(balances):
    """The 40-asset vector exactly as the file holds it; concat uses these strings."""
    missing = [asset for asset in BYBIT_ASSETS if asset not in balances]
    if missing:
        raise ProofError(f"balance map is missing {len(missing)} of the 40 assets, first: {missing[0]}")
    for asset in BYBIT_ASSETS:
        if strict_dec(balances[asset]) is None:
            raise ProofError(f"balance for {asset} is not a decimal")
    return {asset: str(balances[asset]) for asset in BYBIT_ASSETS}


def bybit_concat(vector):
    return "".join(vector[asset] for asset in BYBIT_ASSETS)


def bybit_add(a, b):
    return {asset: shopspring_fixed(round_down(dec(a[asset]) + dec(b[asset]), BYBIT_SCALE), BYBIT_SCALE)
            for asset in BYBIT_ASSETS}


def bybit_combine(left_hash, right_hash, balances, height):
    return sha256_hex(left_hash + right_hash + bybit_concat(balances) + str(height))


def bybit_sibling_type(sibling):
    if sibling.get("type") == 3:
        return 3
    if sibling.get("type") == 1:
        return 2
    return 1


@scheme("bybit")
def verify_bybit(proof):
    scheme_width("bybit", 32)
    self_node = proof.get("self") or {}
    path = proof.get("path") or []
    if not (self_node.get("userHash") or "") or not path:
        raise ProofError("proof is missing self.userHash or path")

    self_balances = bybit_vector(self_node.get("balance") or {})
    expected_self = sha256_hex(self_node["userHash"] + bybit_concat(self_balances))
    if expected_self != (self_node.get("hash") or "").strip().lower():
        return result("bybit", "FAIL", None, "self hash does not match userHash+balances",
                      [{"id": "SELF_HASH_MISMATCH", "note": "recomputed leaf differs from self.hash"}])

    first = path[0]
    if first.get("type") == self_node.get("type"):
        raise ProofError("the first path node shares the self node's side")
    first_balances = bybit_vector(first.get("balance") or {})
    self_hash = (self_node["hash"] if isinstance(self_node.get("hash"), str) else "").strip().lower()
    first_hash = (first.get("hash") or "").strip().lower()

    if self_node.get("type") == 1:
        left, right = self_hash, first_hash
    else:
        left, right = first_hash, self_hash

    node_hash = bybit_combine(left, right, bybit_add(self_balances, first_balances), 2)
    node_balances = bybit_add(self_balances, first_balances)
    height = 1

    for index in range(1, len(path) - 1):
        height += 1
        sibling = path[index]
        sib_balances = bybit_vector(sibling.get("balance") or {})
        sib_hash = (sibling.get("hash") or "").strip().lower()
        if sibling.get("type") == 1:
            left, right = sib_hash, node_hash
        else:
            left, right = node_hash, sib_hash
        node_balances = bybit_add(node_balances, sib_balances)
        node_hash = bybit_combine(left, right, node_balances, height + 1)

    root = path[-1]
    root_balances = bybit_vector(root.get("balance") or {})
    root_hash = (root.get("hash") or "").strip().lower()
    if node_hash != root_hash:
        return result("bybit", "FAIL", node_hash, "path does not recompute to the root hash",
                      [{"id": "PATH_MISMATCH", "note": "the path does not reach the root the proof names"}])
    if node_balances != root_balances:
        return result("bybit", "FAIL", node_hash, "folded balances differ from the root node",
                      [{"id": "BALANCE_MISMATCH",
                        "note": "the balances summed along the path do not equal the root balances"}])
    if height + 1 != root.get("height"):
        return result("bybit", "FAIL", node_hash, "path height does not reach the root height",
                      [{"id": "HEIGHT_MISMATCH", "note": "the path is shorter or longer than the root height says"}])

    return result("bybit", "PASS", node_hash,
                  "path recomputes to the root hash and the folded balances add up",
                  [], {"auditId": self_node.get("auditId")})


# --- coinex: full SHA-256 with JSON balance fold, "left"/anything-else sides ---


def coinex_balances(balances):
    """encoding/json map output: sorted keys, compact, string values."""
    keys = sorted(balances)
    parts = [f'"{k}":"{balanced_value(balances[k])}"' for k in keys]
    return "{" + ",".join(parts) + "}"


def coinex_add(a, b):
    out = dict(a)
    for key, value in b.items():
        out[key] = plain(dec(out.get(key, 0)) + dec(value))
    return out


@scheme("coinex")
def verify_coinex(proof):
    scheme_width("coinex", 32)
    self_node = proof.get("self") or {}
    path = proof.get("path") or []
    root_node = proof.get("root") or {}
    root_hash = normalise_hex(root_node.get("hash"))
    if not root_hash or not path:
        raise ProofError("proof is missing root.hash or path")
    nonce = str(self_node.get("nonce") or "")
    if not nonce.strip():
        raise ProofError("proof is missing self.nonce")

    balances = {k: plain(v) for k, v in (self_node.get("balances") or {}).items()}
    current = sha256_hex(nonce + coinex_balances(balances))
    running = balances
    for index, node in enumerate(path):
        sib_hash = (node.get("hash") or "").strip().lower()
        if not sib_hash:
            # Padding node: sibling hash is self/current hash and balances are zero.
            current = sha256_hex(current + current + coinex_balances(running))
            continue
        sibling_balances = {k: plain(v) for k, v in (node.get("balances") or {}).items()}
        running = coinex_add(running, sibling_balances)
        position = str(node.get("pos") or "").strip().lower()
        if position == "left":
            current = sha256_hex(sib_hash + current + coinex_balances(running))
        else:
            current = sha256_hex(current + sib_hash + coinex_balances(running))

    if current != root_hash:
        return result("coinex", "FAIL", current, "path does not recompute to root.hash",
                      [{"id": "PATH_MISMATCH", "note": "the path does not reach the root the proof names"}])
    root_balances = {k: plain(v) for k, v in (root_node.get("balances") or {}).items()}
    if {k: plain(v) for k, v in running.items()} != root_balances:
        return result("coinex", "FAIL", current, "folded balances differ from the root node balances",
                      [{"id": "BALANCE_MISMATCH",
                        "note": "the balances summed along the path do not equal the root balances"}])
    return result("coinex", "PASS", current,
                  "path recomputes to the root hash and the folded balances match", [])


# --- mexc: SHA-256 truncated to 128 bits, level-indexed path ------------------

MEXC_HEX = 32
MEXC_ASSETS = ("BTC", "USDT", "ETH", "USDC")
MEXC_LEFT, MEXC_RIGHT, MEXC_ROOT = 1, 2, 3


def mexc_raw_balances(balances):
    return ",".join(f"{k}:{balanced_value(balances[k])}" for k in sorted(balances))


def mexc_plain_balances(balances):
    return ",".join(f"{k}:{plain(v)}" for k, v in sorted(balances.items()))


def mexc_inner(left_hash, right_hash, balances):
    return truncated(f"{left_hash},{right_hash},{mexc_plain_balances(balances)}", MEXC_HEX)


def mexc_merge(a, b):
    out = dict(a)
    for key, value in b.items():
        out[key] = plain(dec(out.get(key, 0)) + dec(value))
    return out


@scheme("mexc")
def verify_mexc(proof):
    scheme_width("mexc", MEXC_HEX // 2)
    path = proof.get("pathList") or []
    self_node = proof.get("self") or {}
    if not path or not (self_node.get("hash") or ""):
        raise ProofError("proof is missing pathList or self.hash")
    if not self_node.get("proofId") or not self_node.get("accountId"):
        raise ProofError("proof is missing self.proofId or self.accountId")

    balances = self_node.get("balances") or {}
    content = f"{self_node['proofId']},{self_node['accountId']},{mexc_raw_balances(balances)}"
    if truncated(content, MEXC_HEX) != self_node["hash"]:
        return result("mexc", "FAIL", None, "self hash does not match proofId+accountId+balances",
                      [{"id": "SELF_HASH_MISMATCH", "note": "recomputed leaf differs from self.hash"}])

    root_node = path[-1]
    root_hash = (root_node.get("hash") or "").strip().lower()
    if not root_hash:
        raise ProofError("proof is missing the root hash")
    depth = root_node.get("level") or 0

    level_map = {node.get("level"): node for node in path}
    current = {
        "level": self_node.get("level") or 0,
        "type": self_node.get("type") or 1,
        "hash": self_node["hash"],
        "balances": dict(balances),
    }
    left = right = None
    for level in range(depth):
        sibling = level_map.get(level)
        if sibling is None:
            sib_type = MEXC_RIGHT if current["type"] == MEXC_LEFT else MEXC_LEFT
            sib = {"type": sib_type, "hash": "", "balances": {}}
        else:
            sib = {
                "type": sibling.get("type") or 1,
                "hash": (sibling.get("hash") or "").strip().lower(),
                "balances": dict(sibling.get("balances") or {}),
            }
        if current["type"] == MEXC_LEFT:
            left, right = current, sib
        else:
            left, right = sib, current

        parent_level = level + 1
        parent = level_map.get(parent_level)
        if parent is None:
            parent_type = MEXC_LEFT
        else:
            parent_type = MEXC_RIGHT if (parent.get("type") or 1) == MEXC_LEFT else MEXC_LEFT

        folded = mexc_merge(left["balances"], {})
        if str(right.get("hash") or "").strip():
            for asset in MEXC_ASSETS:
                folded.setdefault(asset, "0")
            folded = mexc_merge(folded, right["balances"])
        current = {
            "level": parent_level,
            "type": parent_type,
            "hash": mexc_inner(left["hash"], right["hash"], folded),
            "balances": folded,
        }

    if left is None:
        raise ProofError("path is empty")
    computed = current["hash"]
    if computed != root_hash:
        return result("mexc", "FAIL", computed, "path does not recompute to the root hash",
                      [{"id": "PATH_MISMATCH", "note": "the path does not reach the root the proof names"}])
    findings = [{"id": "TRUNCATED_HASH",
                 "note": f"MEXC truncates SHA-256 to {MEXC_HEX} hex chars (128 bits); a pass is weak evidence"}]
    return result("mexc", "PASS", computed, "leaf and path recompute to the root hash",
                  findings, {"proofId": self_node.get("proofId")})


# --- ourbit: SHA-256 truncated to 128 bits, currency:market balance keys ------

OURBIT_HEX = 32
OURBIT_EMPTY, OURBIT_LEFT, OURBIT_RIGHT, OURBIT_ROOT = 0, 1, 2, 3


def ourbit_plain(amount):
    text = plain(amount)
    return "0" if text in ("", "-0") else text


def ourbit_balance_map(balances, is_leaf, enabled):
    out = {currency: Decimal(0) for currency in enabled}
    for key, value in balances.items():
        currency = str(key).split(":")[0]
        out[currency] = out.get(currency, Decimal(0)) + dec(value)
    if is_leaf:
        out = {k: v for k, v in out.items() if v != 0}
    return out


def ourbit_balance_string(balances):
    return ",".join(f"{k}:{ourbit_plain(balances[k])}" for k in sorted(balances))


def ourbit_node_hash(node, enabled):
    is_leaf = node["level"] == 0
    balances = ourbit_balance_map(node["balances"], is_leaf, enabled)
    if node["uid"]:
        content = f"PR-{node['snapshotDate']},{node['uid']},{ourbit_balance_string(balances)}"
    else:
        content = f"{node['leftChildHashId']},{node['rightChildHashId']},{ourbit_balance_string(balances)}"
    return truncated(content, OURBIT_HEX)


def ourbit_to_node(row):
    balances = row.get("balances") or {}
    if not isinstance(balances, dict):
        balances = {}
    return {
        "snapshotDate": row.get("snapshotDate") or "",
        "balances": {k: dec(v) for k, v in balances.items()},
        "hash": (row.get("hash") or "").strip().lower(),
        "level": row.get("level") or 0,
        "role": row.get("role") or 0,
        "uid": row.get("uid") or "",
        "leftChildHashId": "",
        "rightChildHashId": "",
    }


def ourbit_enabled(proof, self_node, path):
    declared = proof.get("allCurrencies") or []
    if declared:
        return sorted(str(c) for c in declared)
    seen = set()
    for node in [self_node, *path]:
        for key in node["balances"]:
            seen.add(str(key).split(":")[0])
    return sorted(seen)


def ourbit_construct(left, right, enabled):
    parent = {"balances": {}, "level": left["level"] + 1, "snapshotDate": left["snapshotDate"], "uid": "", "role": 0,
              "hash": "", "leftChildHashId": left["hash"], "rightChildHashId": right["hash"]}
    for key, value in list(left["balances"].items()) + list(right["balances"].items()):
        parent["balances"][key] = parent["balances"].get(key, Decimal(0)) + value
    parent["hash"] = ourbit_node_hash(parent, enabled)
    return parent


def ourbit_parent(friend, self_node, enabled):
    if friend["role"] == OURBIT_LEFT:
        return ourbit_construct(friend, self_node, enabled)
    return ourbit_construct(self_node, friend, enabled)


def ourbit_balances_equal(a, b):
    for key, value in a.items():
        other = b.get(key)
        if other is None or other != value:
            return False
    return True


@scheme("ourbit")
def verify_ourbit(proof):
    scheme_width("ourbit", OURBIT_HEX // 2)
    self_node = ourbit_to_node(proof.get("self") or {})
    path = [ourbit_to_node(row) for row in (proof.get("path") or [])]
    if not self_node["hash"] or not path:
        raise ProofError("proof is missing self.hash or path")
    enabled = ourbit_enabled(proof, self_node, path)

    if ourbit_node_hash(self_node, enabled) != self_node["hash"]:
        return result("ourbit", "FAIL", None, "self hash does not match PR-date+uid+balances",
                      [{"id": "SELF_HASH_MISMATCH", "note": "recomputed leaf differs from self.hash"}])
    first = path[0]
    if not first["hash"] or len(first["balances"]) == 0:
        raise ProofError("first path node is missing hash or balances")
    if first["role"] == self_node["role"]:
        raise ProofError("the first path node shares the self node's role")
    if len(path) <= 1:
        # The Go builder returns an empty node for a path this short, so nothing
        # can be recomputed and the check fails there too.
        return result("ourbit", "FAIL", None, "path is too short to recompute a root",
                      [{"id": "PATH_MISMATCH", "note": "a single-node path cannot reach a root"}])

    parent = ourbit_parent(first, self_node, enabled)
    current = self_node
    for row in path[1:-1]:
        current = parent
        parent = ourbit_parent(row, current, enabled)
    parent["role"] = OURBIT_ROOT

    root_node = path[-1]
    if parent["hash"] != root_node["hash"]:
        return result("ourbit", "FAIL", parent["hash"], "path does not recompute to the root hash",
                      [{"id": "PATH_MISMATCH", "note": "the path does not reach the root the proof names"}])
    if parent["level"] != root_node["level"]:
        return result("ourbit", "FAIL", parent["hash"], "path level does not reach the root level",
                      [{"id": "HEIGHT_MISMATCH", "note": "the path is shorter or longer than the root level says"}])
    if not ourbit_balances_equal(parent["balances"], root_node["balances"]):
        return result("ourbit", "FAIL", parent["hash"], "folded balances differ from the root node balances",
                      [{"id": "BALANCE_MISMATCH",
                        "note": "the balances summed along the path do not equal the root balances"}])

    findings = [{"id": "TRUNCATED_HASH",
                 "note": f"Ourbit truncates SHA-256 to {OURBIT_HEX} hex chars (128 bits); a pass is weak evidence"}]
    proof_id = f"PR-{self_node['snapshotDate']}" if self_node["snapshotDate"] else None
    return result("ourbit", "PASS", parent["hash"],
                  "path recomputes to the root hash and the folded balances add up",
                  findings, {"proofId": proof_id})


# --- toobit: SHA256-SMT with domain separators, integer balances -------------

TOOBIT_LEAF = "SMT_LEAF"
TOOBIT_NODE = "SMT_NODE"


def toobit_parse_balances(raw):
    out = {}
    for asset, amount in (raw or {}).items():
        text = str(amount).strip()
        try:
            out[asset] = int(text, 10)
        except ValueError:
            raise ProofError(f"invalid integer amount for {asset}") from None
    return out


def toobit_leaf_hash(nonce, balances):
    parts = [TOOBIT_LEAF, nonce]
    for asset in sorted(balances):
        parts.append(f"{asset}:{balances[asset]}")
    return sha256_hex("|".join(parts))


def toobit_parent_hash(left_hash, right_hash, balances):
    parts = [TOOBIT_NODE, left_hash, right_hash]
    for asset in sorted(balances):
        parts.append(f"{asset}:{balances[asset]}")
    return sha256_hex("|".join(parts))


def toobit_is_hash(value):
    text = str(value or "").strip()
    return len(text) == 64 and all(c in "0123456789abcdefABCDEF" for c in text)


@scheme("toobit")
def verify_toobit(proof):
    scheme_width("toobit", 32)
    self_node = proof.get("self") or {}
    root_node = proof.get("root") or {}
    path = proof.get("path") or []
    nonce = str(self_node.get("nonce") or "")
    leaf_hash = normalise_hex(self_node.get("leaf_hash"))
    root_hash = normalise_hex(root_node.get("hash"))
    if not nonce or not leaf_hash or not root_hash:
        raise ProofError("proof is missing self.nonce, self.leaf_hash or root.hash")
    if not toobit_is_hash(leaf_hash) or not toobit_is_hash(root_hash):
        raise ProofError("leaf or root hash is not 32-byte hex")

    balances = toobit_parse_balances(self_node.get("balances"))
    if toobit_leaf_hash(nonce, balances) != leaf_hash:
        return result("toobit", "FAIL", None, "self leaf hash does not match nonce+balances",
                      [{"id": "SELF_HASH_MISMATCH", "note": "recomputed leaf differs from self.leaf_hash"}])

    current_hash = leaf_hash
    running = dict(balances)
    sum_ok = True
    for index, sibling in enumerate(path):
        sib_hash = normalise_hex(sibling.get("hash"))
        if not toobit_is_hash(sib_hash or ""):
            raise ProofError(f"path node {index} is not a 32-byte hash")
        sib_balances = toobit_parse_balances(sibling.get("balances"))
        if any(amount < 0 for amount in sib_balances.values()):
            sum_ok = False
        folded = dict(running)
        for asset, amount in sib_balances.items():
            folded[asset] = folded.get(asset, 0) + amount
        position = str(sibling.get("position") or "").strip().lower()
        if position == "left":
            current_hash = toobit_parent_hash(sib_hash, current_hash, folded)
        elif position == "right":
            current_hash = toobit_parent_hash(current_hash, sib_hash, folded)
        else:
            raise ProofError(f"path node {index} has no left/right position")
        running = folded

    if current_hash != root_hash:
        return result("toobit", "FAIL", current_hash, "path does not recompute to the root hash",
                      [{"id": "PATH_MISMATCH", "note": "the path does not reach the root the proof names"}])

    totals = toobit_parse_balances(root_node.get("total_balances"))
    if totals:
        for asset, expected in totals.items():
            if running.get(asset) != expected:
                sum_ok = False
    if not sum_ok:
        return result("toobit", "FAIL", current_hash, "the totals along the path do not match the root totals",
                      [{"id": "BALANCE_MISMATCH",
                        "note": "summed balances or the sign check disagree with the root totals"}])
    return result("toobit", "PASS", current_hash,
                  "leaf and path recompute to the root hash and the totals add up", [],
                  {"auditId": proof.get("audit_id")})


# --- binding to a published root ---------------------------------------------


def normalise_root(value):
    return normalise_hex(value) if isinstance(value, str) else None


def load_published_roots(assessment_path, exchange, onchain_doc=None):
    """Collect the roots the snapshot publishes: assessment file, then the anchor."""
    binding = {
        "source": [],
        "publishedRoot": None,
        "publishedRoots": {},
        "matches": None,
        "widthsDiffer": False,
        "note": None,
        "snapshotId": None,
    }
    path = Path(assessment_path) if assessment_path else None
    if path is not None and path.is_file():
        try:
            assessment = json.loads(path.read_text())
        except (OSError, json.JSONDecodeError):
            binding["note"] = "the assessment file is unreadable; the proof was verified but not bound"
        else:
            if assessment.get("exchange") != exchange:
                binding["note"] = (f"the assessment is for {assessment.get('exchange')!r}, not {exchange!r}; "
                                   "the proof was verified but not bound")
            else:
                snapshot = assessment.get("snapshot") or {}
                binding["snapshotId"] = str(snapshot.get("id")) if snapshot.get("id") is not None else None
                root = normalise_root(snapshot.get("exchangeMerkleRoot"))
                if root:
                    binding["publishedRoots"]["assessment"] = root
                binding["source"].append(str(path))
    elif path is not None:
        binding["note"] = "the assessment file was not found; the proof was verified but not bound"

    if onchain_doc:
        record = onchain_doc.get("record") or {}
        root = normalise_root(record.get("exchangeMerkleRoot"))
        anchor_state = {
            "status": onchain_doc.get("status"),
            "live": onchain_doc.get("live"),
            "contract": onchain_doc.get("contract"),
            "chainId": onchain_doc.get("chainId"),
            "snapshotId": onchain_doc.get("snapshotId"),
            "rpc": onchain_doc.get("rpc"),
        }
        if onchain_doc.get("status") == "found" and root and root != "0" * 64:
            binding["publishedRoots"]["anchor"] = root
            binding["source"].append("on-chain anchor")
            binding["anchor"] = anchor_state
            if binding["snapshotId"] is None and onchain_doc.get("auditId") is not None:
                binding["snapshotId"] = str(onchain_doc.get("auditId"))
        elif onchain_doc.get("status") == "found":
            # Anchored, but the contract holds no exchange merkle root for it.
            anchor_state["note"] = "the anchored record holds no exchange merkle root (zero bytes32)"
            binding["anchor"] = anchor_state
        elif onchain_doc.get("status") == "not_found":
            anchor_state["note"] = "the contract holds no record for this snapshotId; not anchored yet"
            binding["anchor"] = anchor_state
        elif onchain_doc.get("status") == "unavailable":
            anchor_state["error"] = onchain_doc.get("error")
            anchor_state["note"] = "no RPC endpoint answered; the anchor was not consulted"
            binding["anchor"] = anchor_state
        elif onchain_doc.get("status") == "skipped":
            anchor_state["note"] = onchain_doc.get("note")
            binding["anchor"] = anchor_state

    roots = binding["publishedRoots"]
    if roots:
        # When both sources publish a root and they disagree, say so; binding is only
        # meaningful against a single published value.
        distinct = set(roots.values())
        if len(distinct) > 1:
            binding["note"] = "the assessment file and the on-chain record publish different merkle roots"
            binding["publishedRootsDisagree"] = True
        binding["publishedRoot"] = roots.get("anchor") or roots.get("assessment")
    return binding


def bind(binding, doc):
    """Fill in matches / widthsDiffer / note for one verified proof."""
    roots = binding.get("publishedRoots") or {}
    width = doc.get("rootWidthBytes")
    if not roots:
        if binding.get("note") is None:
            anchor = binding.get("anchor") or {}
            if anchor.get("status") == "not_found":
                binding["note"] = ("the snapshot is not anchored yet and publishes no exchange merkle "
                                   "root; the proof is not bound to any snapshot")
            elif anchor.get("status") == "unavailable":
                binding["note"] = ("the anchor could not be read and no assessment root was given; "
                                   "the proof is not bound to any snapshot")
            elif anchor.get("note"):
                binding["note"] = f"no published root to bind against: {anchor['note']}"
            else:
                binding["note"] = ("this snapshot publishes no exchange merkle root; the proof is not "
                                   "bound to the snapshot")
        return binding
    if binding.get("publishedRootsDisagree"):
        return binding

    results = {}
    width_problem = None
    recomputed = normalise_root(doc.get("recomputedRoot"))
    for label, root in roots.items():
        if isinstance(width, int) and width * 2 != len(root):
            width_problem = (f"the {label} publishes a {len(root) * 4}-bit root but this proof scheme "
                             f"produces {width * 8}-bit roots; they cannot be compared directly")
            results[label] = None
            # Padding is not a binding, but it is worth reporting: an exchange that keeps a
            # truncated root in a bytes32 word zero-pads it, so the bytes can still line up.
            trimmed = root.lstrip("0")
            if isinstance(width, int) and trimmed and len(trimmed) == width * 2:
                binding["rightAligned"] = {
                    "publishedWidthBits": len(root) * 4,
                    "schemeWidthBits": width * 8,
                    "equal": recomputed == trimmed,
                    "note": ("the published root is zero-padded to 32 bytes; a right-aligned comparison "
                             "is informational only and is not a binding"),
                }
        else:
            results[label] = recomputed == root
    binding["matches"] = results.get("anchor") if "anchor" in results else results.get("assessment")
    if width_problem:
        binding["widthsDiffer"] = True
        binding["note"] = width_problem
        binding["matches"] = None
    if binding.get("matches") is False and not width_problem:
        binding["note"] = "the proof root differs from the merkle root the snapshot publishes"
    return binding


# --- main ---------------------------------------------------------------------


def load_onchain(args, exchange):
    if not (args.onchain or args.anchor_contract or args.rpc_url):
        return None
    audit_id = args.audit_id or args.onchain
    if not audit_id:
        return {"status": "skipped", "note": "--anchor-contract needs --audit-id or --onchain <auditId>"}
    import por_anchor

    doc = por_anchor.resolve_record(
        audit_id,
        exchange=exchange,
        contract=args.anchor_contract,
        rpc=args.rpc_url,
        fixture=args.anchor_fixture,
    )
    return doc


def main(argv):
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument("proof", nargs="?", help="path to the proof JSON the exchange gave you")
    parser.add_argument("--exchange", help="one of: " + ", ".join(sorted(SCHEMES)))
    parser.add_argument("--assessment", help="assessment JSON to bind the proof root against")
    parser.add_argument("--onchain", metavar="AUDIT_ID",
                        help="bind against the ArdmerePoRAnchor record for this audit id")
    parser.add_argument("--anchor-contract", help="override the anchor proxy address")
    parser.add_argument("--rpc-url", help="override the Sepolia RPC endpoint")
    parser.add_argument("--anchor-fixture", help="recorded anchor read to use instead of an RPC call")
    parser.add_argument("--audit-id", help="audit id to look up when --anchor-contract is given alone")
    parser.add_argument("--list-schemes", action="store_true", help="print the supported schemes")
    opts = parser.parse_args(argv)

    if opts.list_schemes:
        print("\n".join(sorted(SCHEMES)))
        return 0
    if not opts.proof or not opts.exchange:
        parser.error("a proof path and --exchange are required")

    exchange = opts.exchange.strip().lower()
    if exchange not in SCHEMES:
        print(f"inclusion_check: unsupported exchange {exchange!r}; try: {', '.join(sorted(SCHEMES))}",
              file=sys.stderr)
        return 2
    try:
        proof = json.loads(Path(opts.proof).read_text())
    except (OSError, json.JSONDecodeError) as exc:
        print(f"inclusion_check: cannot read the proof: {exc}", file=sys.stderr)
        return 1
    if not isinstance(proof, dict):
        print("inclusion_check: the proof must be a JSON object", file=sys.stderr)
        return 1

    try:
        doc = SCHEMES[exchange](proof)
    except ProofError as exc:
        print(f"inclusion_check: {exc}", file=sys.stderr)
        return 1

    onchain = load_onchain(opts, exchange)
    if opts.assessment or onchain:
        binding = load_published_roots(opts.assessment, exchange, onchain)
        bind(binding, doc)
        doc["binding"] = binding
        if binding.get("widthsDiffer"):
            doc["findings"].append({"id": "ROOT_WIDTH_MISMATCH", "note": binding["note"]})
        elif binding.get("matches") is False:
            doc["verdict"] = "FAIL"
            doc["detail"] = "the proof root differs from the merkle root the snapshot publishes"
            doc["findings"].append({"id": "ROOT_NOT_BOUND",
                                    "note": "the snapshot publishes a different merkle root"})

    doc["scope"] = (
        "Proves this account's leaf is in the tree that produced the proof root. "
        "It does not prove the exchange holds enough reserves for everyone, and it says nothing "
        "about any other account."
    )
    json.dump(doc, sys.stdout, indent=2)
    sys.stdout.write("\n")
    return 0 if doc["verdict"] == "PASS" else 1


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))

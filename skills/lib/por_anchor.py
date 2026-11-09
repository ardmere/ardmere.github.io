"""Read an ArdmerePoRAnchor snapshot record over plain JSON-RPC.

Standard library only: Keccak-256 and the SnapshotRecord ABI decode are
implemented here rather than pulled from web3, so this runs on a bare
python3 with no pip install. When Foundry's `cast` is on PATH the caller can
use it instead; both paths must produce the same record.

Addresses, field order and error names follow docs/anchor-query-api.md
(schema v4). This module is a reader: it has no write path.
"""

import json
import os
import shutil
import subprocess
import urllib.error
import urllib.request

DEFAULT_CONTRACT = "0x0A5eB9f6c173429DBb418826EDFDf7fFe11433f7"  # pragma: allowlist secret
DEFAULT_RPC_URL = "https://ethereum-sepolia-rpc.publicnode.com"  # pragma: allowlist secret
FALLBACK_RPC_URLS = (
    "https://ethereum-sepolia.publicnode.com",  # pragma: allowlist secret
    "https://sepolia.gateway.tenderly.co",
    "https://1rpc.io/sepolia",
)
CHAIN_ID = 11155111  # pragma: allowlist secret
SCHEMA_VERSION = 5
CUSTODY_RISK_UNAVAILABLE = 255

# SnapshotRecord, in ABI order (docs/anchor-query-api.md §4).
RECORD_FIELDS = (
    ("snapshotId", "bytes32"),
    ("exchangeTag", "bytes32"),
    ("exchange", "string"),
    ("periodSeq", "uint32"),
    ("snapshotTime", "uint64"),
    ("btcBlockHeight", "uint32"),
    ("exchangeMerkleRoot", "bytes32"),
    ("artifactBundleRoot", "bytes32"),
    ("verificationBundleRoot", "bytes32"),
    ("verdictSummary", "uint8"),
    ("coverageBps", "uint16"),  # deprecated opaque max coverage; prefer wallet/reserve
    ("schemaVersion", "uint8"),
    ("anchoredAt", "uint256"),
    ("porStage", "uint8"),
    ("effectivePoR", "bool"),
    ("walletCoverageBps", "uint16"),
    ("reserveCoverageBps", "uint16"),
    ("custodyRiskScore", "uint8"),
)
RECORD_TUPLE = "(" + ",".join(kind for _name, kind in RECORD_FIELDS) + ")"

VIEW_SIGNATURES = {
    "getSnapshot": f"getSnapshot(bytes32){RECORD_TUPLE}",
    "latest": f"getLatestSnapshot(bytes32){RECORD_TUPLE}",
    "byPeriod": f"getSnapshotByPeriod(bytes32,uint32){RECORD_TUPLE}",
    "exists": "snapshotExists(bytes32)(bool)",
    "count": "getSnapshotCount(bytes32)(uint256)",
    "idAt": "getSnapshotIdAt(bytes32,uint256)(bytes32)",
    "schemaVersion": "SCHEMA_VERSION()(uint8)",
}


class AnchorError(RuntimeError):
    """The RPC or the contract could not answer."""


# --- keccak-256 (original Keccak padding, as Solidity uses) ---------------------

_RC = (
    0x0000000000000001, 0x0000000000008082, 0x800000000000808A, 0x8000000080008000,
    0x000000000000808B, 0x0000000080000001, 0x8000000080008081, 0x8000000000008009,
    0x000000000000008A, 0x0000000000000088, 0x0000000080008009, 0x000000008000000A,
    0x000000008000808B, 0x800000000000008B, 0x8000000000008089, 0x8000000000008003,
    0x8000000000008002, 0x8000000000000080, 0x000000000000800A, 0x800000008000000A,
    0x8000000080008081, 0x8000000000008080, 0x0000000080000001, 0x8000000080008008,
)
_ROT = {
    (0, 0): 0, (1, 0): 1, (2, 0): 62, (3, 0): 28, (4, 0): 27,
    (0, 1): 36, (1, 1): 44, (2, 1): 6, (3, 1): 55, (4, 1): 20,
    (0, 2): 3, (1, 2): 10, (2, 2): 43, (3, 2): 25, (4, 2): 39,
    (0, 3): 41, (1, 3): 45, (2, 3): 15, (3, 3): 21, (4, 3): 8,
    (0, 4): 18, (1, 4): 2, (2, 4): 61, (3, 4): 56, (4, 4): 14,
}
_MASK = (1 << 64) - 1


def _rotl(value, shift):
    if not shift:
        return value
    return ((value << shift) | (value >> (64 - shift))) & _MASK


def _keccak_f(state):
    for rnd in range(24):
        c = [state[x] ^ state[x + 5] ^ state[x + 10] ^ state[x + 15] ^ state[x + 20] for x in range(5)]
        d = [c[(x - 1) % 5] ^ _rotl(c[(x + 1) % 5], 1) for x in range(5)]
        for x in range(5):
            for y in range(5):
                state[x + 5 * y] ^= d[x]
        b = [0] * 25
        for x in range(5):
            for y in range(5):
                b[y + 5 * ((2 * x + 3 * y) % 5)] = _rotl(state[x + 5 * y], _ROT[(x, y)])
        for x in range(5):
            for y in range(5):
                state[x + 5 * y] = (
                    b[x + 5 * y]
                    ^ ((~b[(x + 1) % 5 + 5 * y]) & _MASK & b[(x + 2) % 5 + 5 * y])
                )
        state[0] ^= _RC[rnd]
    return state


def keccak256(data):
    """bytes -> 32-byte digest, original Keccak (0x01 pad), not SHA3 (0x06 pad)."""
    rate = 136
    block = bytearray(data) + b"\x01"
    while len(block) % rate:
        block.append(0)
    block[-1] ^= 0x80
    state = [0] * 25
    for offset in range(0, len(block), rate):
        chunk = block[offset:offset + rate]
        for i in range(rate // 8):
            state[i] ^= int.from_bytes(chunk[i * 8:i * 8 + 8], "little")
        _keccak_f(state)
    return b"".join(state[i].to_bytes(8, "little") for i in range(4))


def keccak_hex(data):
    return "0x" + keccak256(data).hex()


def exchange_tag(exchange):
    """keccak256(bytes(exchange)) — the bytes32 key the contract indexes on."""
    return keccak_hex(exchange.encode())


def snapshot_id(audit_id):
    """keccak256(bytes(auditId)) — the bytes32 key getSnapshot takes."""
    return keccak_hex(str(audit_id).encode())


def selector(signature):
    return "0x" + keccak256(signature.encode())[:4].hex()


# --- ABI decode ----------------------------------------------------------------


def _word(data, index):
    start = index * 32
    chunk = data[start:start + 32]
    if len(chunk) != 32:
        raise AnchorError("return data is shorter than the ABI head")
    return int.from_bytes(chunk, "big")


def decode_record(payload):
    """Decode a SnapshotRecord tuple from eth_call return data.

    Accepts both schema-v3 and schema-v4 ABI lengths. When the return is a v3
    payload (ending at anchoredAt), the v4 fields default to 0 / false.
    """
    data = bytes.fromhex(payload[2:] if payload.startswith("0x") else payload)
    if len(data) >= 32 and int.from_bytes(data[:32], "big") == 0x20:
        # Dynamic struct return: a head offset word precedes the tuple.
        data = data[32:]
    # Minimum = fields through anchoredAt (schema v3). Indices 0..12 inclusive.
    min_fields = 13
    if len(data) < 32 * min_fields:
        raise AnchorError("return data does not hold a full SnapshotRecord")
    record = {}
    for index, (name, kind) in enumerate(RECORD_FIELDS):
        if index * 32 + 32 > len(data) and kind != "string":
            # Trailing v4 fields absent from a v3-length return.
            if kind == "bool":
                record[name] = False
            elif kind.startswith("uint") or kind == "uint8" or kind == "uint16":
                record[name] = 0
            else:
                record[name] = 0
            continue
        if kind in ("bytes32",):
            start = index * 32
            record[name] = "0x" + data[start:start + 32].hex()
        elif kind == "string":
            offset = _word(data, index)
            length = _word(data, offset // 32)
            raw = data[offset + 32:offset + 32 + length]
            record[name] = raw.decode("utf-8", "replace")
        elif kind == "bool":
            if index * 32 + 32 > len(data):
                record[name] = False
            else:
                record[name] = _word(data, index) != 0
        else:
            if index * 32 + 32 > len(data):
                record[name] = 0
            else:
                record[name] = _word(data, index)
    if record["anchoredAt"] == 0:
        # Zeroed struct means "no record", whatever the RPC returned.
        return None
    return record


def decode_bool(payload):
    return _word(bytes.fromhex(payload[2:]), 0) != 0


def decode_uint(payload):
    return _word(bytes.fromhex(payload[2:]), 0)


def decode_error(payload):
    """Map revert data to a contract error name plus its args."""
    if not payload or payload in ("0x", "0x0"):
        return None
    data = bytes.fromhex(payload[2:])
    if len(data) < 4:
        return None
    head = "0x" + data[:4].hex()
    for error, args in (
        ("SnapshotNotFound", ["bytes32"]),
        ("SnapshotIndexOutOfBounds", ["bytes32", "uint256"]),
        ("Unauthorized", ["address"]),
        ("InvalidRoot", []),
        ("EmptyExchange", []),
        ("SnapshotAlreadyExists", ["bytes32"]),
        ("PeriodAlreadyAnchored", ["bytes32", "uint32"]),
    ):
        signature = error + "(" + ",".join(args) + ")"
        if selector(signature) == head:
            decoded = []
            for index, kind in enumerate(args):
                word = data[4 + index * 32:4 + (index + 1) * 32]
                if kind == "bytes32":
                    decoded.append("0x" + word.hex())
                else:
                    decoded.append(int.from_bytes(word, "big"))
            return {"name": error, "args": decoded}
    return {"name": "unrecognised revert", "args": ["0x" + data.hex()]}


# --- JSON-RPC ------------------------------------------------------------------


def _post(url, body, timeout):
    request = urllib.request.Request(
        url,
        data=json.dumps(body).encode(),
        headers={"Content-Type": "application/json", "User-Agent": "ardmere-por-skills/1.0"},
    )
    with urllib.request.urlopen(request, timeout=timeout) as response:
        return json.loads(response.read().decode())


def rpc_urls(explicit=None):
    """Configured RPC first, then the public fallbacks (POR_RPC_URL wins).

    `explicit` may be one URL or a sequence of them, so a caller can pin a single
    endpoint (and get a clean 'unavailable' when it is down) or pass a list to try.
    """
    if explicit:
        if isinstance(explicit, str):
            return [explicit]
        return [str(url) for url in explicit]
    configured = os.environ.get("POR_RPC_URL")
    if configured:
        return [configured]
    return [DEFAULT_RPC_URL, *FALLBACK_RPC_URLS]


def rpc_call(url, to, data, timeout=20):
    """eth_call at latest; returns the raw hex result or raises AnchorError."""
    body = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": "eth_call",
        "params": [{"to": to, "data": data}, "latest"],
    }
    try:
        payload = _post(url, body, timeout)
    except (urllib.error.URLError, OSError, ValueError) as exc:
        raise AnchorError(f"{url}: {exc}") from exc
    if "error" in payload:
        error = payload["error"] or {}
        revert = error.get("data") or {}
        if isinstance(revert, dict):
            revert = revert.get("data") or revert.get("result") or ""
        decoded = decode_error(revert) if isinstance(revert, str) else None
        if decoded:
            raise _Revert(decoded, error.get("message"))
        raise AnchorError(error.get("message") or "rpc error")
    return payload.get("result") or "0x"


class _Revert(AnchorError):
    def __init__(self, decoded, message=None):
        self.decoded = decoded
        super().__init__(message or f"reverted with {decoded.get('name')}")


def chain_id(url, timeout=20):
    payload = _post(url, {"jsonrpc": "2.0", "id": 1, "method": "eth_chainId", "params": []}, timeout)
    if "error" in payload:
        raise AnchorError(payload["error"].get("message") or "chain id unavailable")
    return int(payload["result"], 16)


def cast_available():
    return shutil.which("cast") is not None


def cast_call(url, to, signature, *args, timeout=25):
    command = ["cast", "call", to, signature, *[str(a) for a in args], "--rpc-url", url]
    proc = subprocess.run(command, capture_output=True, text=True, timeout=timeout)
    if proc.returncode != 0:
        raise AnchorError((proc.stderr or "cast call failed").strip()[:400])
    return proc.stdout.strip()


# --- the reader ----------------------------------------------------------------


def contract_address(explicit=None):
    return explicit or os.environ.get("POR_ANCHOR_CONTRACT") or DEFAULT_CONTRACT


def read_snapshot(audit_id, exchange=None, contract=None, rpc=None, prefer_cast=False, timeout=20):
    """Read the v4 record for one auditId.

    Returns a dict whose `status` is one of:
      found      — the contract returned a populated SnapshotRecord
      not_found  — the contract answered SnapshotNotFound / anchoredAt 0
      unavailable— no RPC answered; nothing can be concluded

    `live` is true only when an RPC answered in this process.
    """
    contract = contract_address(contract)
    sid = snapshot_id(audit_id)
    result = {
        "contract": contract,
        "chainId": CHAIN_ID,
        "network": "ethereum-sepolia",
        "auditId": str(audit_id),
        "snapshotId": sid,
        "exchangeTag": exchange_tag(exchange) if exchange else None,
        "live": False,
        "rpc": None,
        "status": "unavailable",
        "record": None,
        "error": None,
        "note": None,
    }
    errors = []
    for url in rpc_urls(rpc):
        try:
            if prefer_cast and cast_available():
                payload = cast_call(url, contract, VIEW_SIGNATURES["getSnapshot"], sid)
                if payload in ("", "0x"):
                    record = None
                else:
                    record = decode_record(payload)
            else:
                payload = rpc_call(url, contract, selector("getSnapshot(bytes32)") + sid[2:], timeout)
                record = decode_record(payload)
        except _Revert as exc:
            if (exc.decoded or {}).get("name") == "SnapshotNotFound":
                result.update(live=True, rpc=url, status="not_found", record=None,
                              note="the contract answered SnapshotNotFound for this snapshotId")
                return result
            errors.append(str(exc))
            continue
        except AnchorError as exc:
            errors.append(str(exc))
            continue
        result.update(live=True, rpc=url, status="found" if record else "not_found", record=record)
        if record is None:
            result["note"] = "the contract returned an empty SnapshotRecord (anchoredAt = 0)"
        return result
    result["error"] = "; ".join(errors[:3]) or "no RPC endpoint answered"
    result["note"] = "every RPC endpoint failed; this is not evidence about the snapshot"
    result["rpc"] = None
    return result


def read_latest(exchange, contract=None, rpc=None, timeout=20):
    """Most recent v4 anchor for an exchange tag."""
    contract = contract_address(contract)
    tag = exchange_tag(exchange)
    result = {
        "contract": contract,
        "chainId": CHAIN_ID,
        "exchange": exchange,
        "exchangeTag": tag,
        "live": False,
        "rpc": None,
        "status": "unavailable",
        "record": None,
        "error": None,
        "note": None,
    }
    errors = []
    for url in rpc_urls(rpc):
        try:
            payload = rpc_call(url, contract, selector("getLatestSnapshot(bytes32)") + tag[2:], timeout)
            record = decode_record(payload)
        except _Revert as exc:
            if (exc.decoded or {}).get("name") == "SnapshotNotFound":
                result.update(live=True, rpc=url, status="not_found",
                              note="this exchange has no v4 anchor yet")
                return result
            errors.append(str(exc))
            continue
        except AnchorError as exc:
            errors.append(str(exc))
            continue
        result.update(live=True, rpc=url, status="found" if record else "not_found", record=record)
        return result
    result["error"] = "; ".join(errors[:3]) or "no RPC endpoint answered"
    return result


def read_exists(audit_id, contract=None, rpc=None, timeout=20):
    contract = contract_address(contract)
    sid = snapshot_id(audit_id)
    for url in rpc_urls(rpc):
        try:
            payload = rpc_call(url, contract, selector("snapshotExists(bytes32)") + sid[2:], timeout)
        except AnchorError:
            continue
        return {"live": True, "rpc": url, "snapshotId": sid, "exists": decode_bool(payload)}
    return {"live": False, "rpc": None, "snapshotId": sid, "exists": None}


def load_fixture(path):
    """A recorded reader result, kept for offline runs. Never claims to be live.

    A fixture may omit `contract` / `chainId` / `rpc`; the defaults in this module
    are then reported, because that is what the recording was taken against. A
    fixture that does carry them (a fork, a redeploy) wins.
    """
    with open(path) as handle:
        doc = json.load(handle)
    doc["live"] = False
    doc.setdefault("contract", DEFAULT_CONTRACT)
    doc.setdefault("chainId", CHAIN_ID)
    doc.setdefault("network", "ethereum-sepolia")
    doc.setdefault("note", "recorded fixture, not a chain read in this run")
    return doc


def resolve_record(audit_id, exchange=None, contract=None, rpc=None, fixture=None, prefer_cast=False):
    """Live read, or the recorded fixture when the caller asked for one."""
    if fixture:
        doc = load_fixture(fixture)
        rows = doc.get("records") or [doc]
        for row in rows:
            if str(row.get("auditId")) == str(audit_id) or row.get("snapshotId") == snapshot_id(audit_id):
                merged = dict(doc)
                merged.update({k: v for k, v in row.items() if k != "records"})
                merged["live"] = False
                # The file's own `source` describes how it was recorded, not where this
                # run read it from, so name the file that was actually consulted.
                merged["source"] = str(fixture)
                merged.setdefault("status", "found" if merged.get("record") else "not_found")
                return merged
        return {
            "live": False, "status": "not_found", "record": None, "rpc": None,
            "contract": doc.get("contract"), "chainId": doc.get("chainId", CHAIN_ID),
            "auditId": str(audit_id), "snapshotId": snapshot_id(audit_id),
            "source": str(fixture),
            "note": "the recorded fixture holds no record for this snapshotId",
        }
    doc = read_snapshot(audit_id, exchange=exchange, contract=contract, rpc=rpc, prefer_cast=prefer_cast)
    doc["source"] = "rpc"
    return doc

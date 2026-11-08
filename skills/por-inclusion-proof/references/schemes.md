# Inclusion schemes

Each section is one scheme, ported from the Go package named at the top. The tests in `tests/test_inclusion.py` check every row against the same fixture the corresponding Go test uses, so a disagreement between Python and Go is a test failure rather than a silent divergence.

| Scheme | Go source | Root | Balances folded | Fixture | Root asserted |
| --- | --- | --- | --- | --- | --- |
| `bitget` | `internal/exchanges/bitget/merkle` | SHA-256 → 64 bits | yes, JSON object at every level | `fixtures/bitget/merkel_tree_bg.json` | `ca89456bb711c913` (Go test) |
| `bingx` | `internal/exchanges/bingx/merkle` | SHA-256 → 256 bits | no | `fixtures/bingx/mock_user_proof.json` | `e12dba68…4347`, the fixture's own `merkleRoot`; the Go test asserts validity only |
| `bitunix` | `internal/exchanges/bitunix/merkle` | SHA-256 → 256 bits | yes, at per-asset precision | `fixtures/bitunix/merkle_tree_bitunix.json` | `7d409a80…679d` (Go test) |
| `bybit` | `internal/exchanges/bybit/merkle` | SHA-256 → 256 bits | yes, fixed 40-asset vector | `fixtures/bybit/mock_user_merkle_tree_path_40_v5.json` | `c91d5de0…bc03` (Go test) |
| `coinex` | `internal/exchanges/coinex/merkle` | SHA-256 → 256 bits | yes, JSON object | `fixtures/coinex/mock_user_proof.json` | `8d44e525…f05f`, the fixture's own `root.hash`; the Go test asserts validity only |
| `mexc` | `internal/exchanges/mexc/merkle` | SHA-256 → 128 bits | yes, level-indexed | `fixtures/mexc/merkle_proof.json` | `35d01ffe…b0cc` (Go test) |
| `ourbit` | `internal/exchanges/ourbit/merkle` | SHA-256 → 128 bits | yes, `currency:market` keys | `fixtures/ourbit/merkel_tree_ourbit.json` | `df1d3451…ad8e` (Go test) |
| `toobit` | `internal/exchanges/toobit/merkle` | SHA-256 → 256 bits | yes, integer amounts | `fixtures/toobit/user_1_demo_proof.json` | `3586fdd2…f727` (Go test) |

The two rows marked "validity only" still fail the Python suite if the fold is wrong: the fixture carries its own declared root, and the test asserts the recomputed root equals it. What Go does not pin for those two is the *value*, only that the path validates — which is why the Python test asserts the fixture value explicitly rather than trusting the Go test's silence.

Full roots, for copy-paste:

```text
bitget  ca89456bb711c913
bingx   e12dba68a58da352b4fd8901f3d989e361cba37ce1fc4b28319b285c42244347
bitunix 7d409a80ebbbf4cef17f3f9c5f56f6e26040f1ca285395897cf695b21391679d
bybit   c91d5de0554f97244e4d9f8056fad70fa0cb2cdb23c290ec597a042645dcbc03
coinex  8d44e525e0ae827da3d27f8f27b50d4c5bc8d94b03b7a29bbf7219a8ba9cf05f
mexc    35d01ffef2beec3ff0166c03f3d9b0cc
ourbit  df1d34512abd50712b9f5713b51fad8e
toobit  3586fdd20930df70e214ca0157c68bb6226289824f05f92d5e0465eb5297f727
```

## `bitget`

Leaf preimage is `encryptUid,nonce,{"ASSET":"amount",...}` with keys sorted and amounts in shopspring `Decimal.String()` form (no exponent, no trailing zeros). Every inner node is

```text
sha256(left || right || "," || balances_json || "," || level)[:16]
```

where `balances_json` is the running sum of the balances along the path from the leaf up. `position` may be `left`/`l`; anything else is treated as `right`, matching the Go package. On a pass the leaf's folded balances must also equal `root.balances`, so a proof with a valid hash but a doctored balance map still fails.

64-bit roots are weak evidence. The verifier in Go marks them `VerdictWarn`; this skill emits `TRUNCATED_HASH` on every pass.

## `bingx`

The simplest scheme: a leaf hash and a flat path of `side:hash` pairs.

```text
leaf  = sha256(assetLine)              # or merkleLeaf, already hashed
node  = side == "l" ? sha256(sibling || node) : sha256(node || sibling)
```

No balance fold, no level, no height. Path entries accept `l:hash` and `hash:l` because the official README and its `por.py` disagree on the order; the Go package accepts both and so does this port.

## `bitunix`

A `self` node plus a height-descending path. The leaf is

```text
sha256(encryptUid || "," || nonce || "," || balances_json)
```

and each parent is

```text
sha256(left.merkleLeaf || right.merkleLeaf || "," || balances_json || "," || parent.height)
```

Balances are summed per asset and rounded **down** to that asset's `precisions[asset]`, then rendered at that fixed scale. The sibling's `type` decides the order: `1` is left, anything else places it right. The parent's height is `left.height - 1`, which is why the path is ordered from the leaf upward and the root is the last entry.

## `bybit`

The Gen 1.5 scheme over a fixed 40-asset vector (`Asset40V5`, in enum order, not alphabetical). The leaf is

```text
sha256(userHash || concat(balance vector))
```

and each combine is

```text
sha256(left || right || concat(balances) || height)
```

with balances summed and rounded down at scale 8. Three things are checked beyond the hash: the folded balances must equal the root node's, the path height must reach `root.height`, and the first path node must not sit on the same side as `self`. A proof whose last node is the root carries the root in `path[-1]`.

The 40-asset requirement is enforced: a balance map missing any asset is rejected, because the vector order would otherwise shift and the hash would be meaningless.

## `coinex`

Leaf and inner nodes both fold a JSON balance object:

```text
leaf  = sha256(self.nonce || balances_json)
node  = sha256(leftHash? || hash || balances_json)      # pos == "left" puts the sibling first
```

The official CLI concatenates self first for any position other than `left`, so `"right"` and unknown values behave the same. A path entry with an empty `hash` is a padding node: the sibling hash is the running hash itself and the balances are unchanged, which is how the tree fills an odd level.

## `mexc`

SHA-256 truncated to 32 hex characters. The path is indexed by `level` rather than ordered, and a missing level means a synthetic empty sibling of the opposite type.

```text
leaf = sha256(proofId || "," || accountId || "," || "ASSET:raw,...")[:32]
node = sha256(left || "," || right || "," || "ASSET:plain,...")[:32]
```

Note the difference between the two: the leaf uses the raw amount strings from the file, the inner nodes use `Decimal.String()`. When the right sibling is non-empty, the four assets `BTC`, `USDT`, `ETH`, `USDC` are seeded at `"0"` before folding, because the official verifier does that and the root hash depends on it.

128 bits is weaker than full SHA-256 but stronger than Bitget's 64. `TRUNCATED_HASH` still fires.

## `ourbit`

The closest scheme to MEXC, with three differences worth knowing:

* Balances are keyed `currency:market` (`BTC:spot`, `BTC:swap`) and summed per currency before hashing.
* `allCurrencies` decides which currencies appear in the preimage; when it is absent the list is derived from the keys actually present.
* A leaf drops zero currencies; an inner node keeps them, so `BTC:0` appears in a parent preimage but never in a leaf.

The leaf preimage is `PR-<snapshotDate>,<uid>,<balances>` and a parent is `<leftHash>,<rightHash>,<balances>`, each truncated to 32 hex characters. `role` decides the order: `1` is left. The final node must match the root's hash, level, **and** folded balances.

## `toobit`

A Summation Merkle Tree with domain separators, and the only scheme whose balances are integers rather than decimals:

```text
leaf = sha256("SMT_LEAF|"  || nonce   || "|ASSET:amount"... )
node = sha256("SMT_NODE|" || leftHash || "|" || rightHash || "|ASSET:sum"... )
```

Assets are sorted by name and amounts are base-10 integers, so no decimal rounding is involved. The root's `total_balances` are checked against the sums along the path when the file carries them, and a negative sibling amount fails the sum check.

---

## What is not here

These exchanges have public assessments in this repository but no user proof scheme, so they are not listed by `--list-schemes` rather than shipped as stubs:

| Exchange | Why not |
| --- | --- |
| `binance` | Public PoR is a wallet bundle, not a per-user Merkle proof |
| `okx` | Public Merkle bundle, no per-user proof file in the verifier |
| `kucoin` | Wallet list only; no inclusion proof artifact |
| `gateio` | Account-level artifact, no per-user proof scheme in Go |
| `poloniex`, `htx`, `phemex`, `woox`, `kraken`, `coinw` | No user proof scheme in this repository |
| `btcc` | Truncated Merkle root, no user proof scheme in Go |

If one of these gains a proof scheme, port the Go package, add it to `CASES` in `tests/test_inclusion.py` with the root its Go test asserts, and only then let it appear in `--list-schemes`.

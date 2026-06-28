# Bitget Proof-of-Reserves — data guide

Bitget uses a **Gen 1** pure Merkle Tree scheme. Its public Java verifier truncates SHA-256 to 16 hex characters (64 bits), so ardmere treats user inclusion as **weak evidence**.

## Public data planes

| Plane | Source | Public? | ardmere artifact |
|-------|--------|---------|------------------|
| Reserve ratio summary | PoR page / browser API | Browser capture | `summarySnapshot` |
| User Merkle path | Logged-in user downloads `merkel_tree_bg.json` | Login | `userMerkleProof` |
| Wallet addresses | Page claims wallet snapshots, but no stable machine-readable signed ZIP in ardmere | Not integrated | — |
| Global ZK | Not used | — | — |
| Third-party attestation | No clearly identified independent attestation | — | — |

## Official references

| Resource | URL |
|----------|-----|
| PoR page | https://www.bitget.com/proof-of-reserves |
| May 2026 announcement | https://www.bitget.com/support/articles/12560603884418 |
| Verification guide | https://www.bitget.com/academy/how-to-verify-your-assets-in-bitget |
| Open-source verifier | https://github.com/BitgetLimited/proof-of-reserves |

## Fetch

The PoR page is browser-oriented. Use browser DevTools to capture summary JSON, or import a saved summary bundle:

```bash
por fetch bitget -info-file ./info.json -coins-file ./coins.json

por fetch bitget \
  -summary-path ./fixtures/bitget/202605-summary.json \
  -user-proof ./fixtures/bitget/merkel_tree_bg.json
```

One-shot archive + verify + bundle:

```bash
por anchor -exchange bitget \
  -summary-path ./fixtures/bitget/202605-summary.json \
  -user-proof ./fixtures/bitget/merkel_tree_bg.json \
  -skip-rpc
```

## Active verifiers

| Verifier | Meaning |
|----------|---------|
| `artifact-integrity@1` | Archived summary/proof SHA256 matches bundle |
| `solvency-claim@1` | Self-reported total/per-coin reserve ratios >= 100% |
| `user-merkle-proof@bitget-1` | Recomputes `merkel_tree_bg.json` under Bitget's 64-bit truncated SHA-256 scheme |

## Stubbed dimensions

| Verifier | Reason |
|----------|--------|
| `internal-consistency@0` | No public wallet CSV to reconcile against summary |
| `onchain-balance-*@0` | No stable signed address bundle integrated |
| `address-ownership@0` | No public ownership signatures in ardmere |
| `global-zk-proof@0` | Bitget has no ZK proof |
| `third-party-attestation@0` | No clearly identified machine-readable independent attestation |

## Merkle proof semantics

Bitget's verifier uses:

```text
leaf = SHA256(encryptUid + "," + nonce + "," + JSON(balances))[0:16]
node = SHA256(leftHash + rightHash + "," + JSON(aggregatedBalances) + "," + level)[0:16]
```

`[0:16]` means the first 16 hex characters, i.e. 64-bit output. Collision resistance is only about 32 bits.

ardmere normalizes the 16-hex root to a left-padded `bytes32` for bundle/anchor compatibility, while the verifier compares the original 16-hex user proof root after trimming leading zeros.

## Boundaries

- `user-merkle-proof@bitget-1 PASS` proves only that the supplied `merkel_tree_bg.json` is consistent with Bitget's weak truncated-hash scheme.
- It does **not** prove exchange-wide solvency, wallet ownership, on-chain reserves, or full liability coverage.
- The official repository has been stale since 2023 and contains no visible test suite; production claims may have drifted from the open-source code.

# Deployments — ArdmerePoRAnchor

On-chain anchor contract deployments and first anchor transactions.

## Base Sepolia (testnet, chain id `84532`)

Bootstrap deployment using the dev signer from `PRIVATE_KEY` / `ANCHOR_SIGNER` in `~/.zshrc` (testnet only — **never use on mainnet**).

| Item | Value |
|---|---|
| **Contract** | [`0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9`](https://sepolia.basescan.org/address/0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9) |
| **Anchor signer** | `0xf2674A2b11b4a6CedC94ab57b22c86Df1fF36209` |
| **Schema version** | `2` (merged single-tx anchor) |
| **Deploy tx** | [`0x3e844a1099f932638f695ce0d3045f78e5cc4e5def63d0edfb40bb603bb2464c`](https://sepolia.basescan.org/tx/0x3e844a1099f932638f695ce0d3045f78e5cc4e5def63d0edfb40bb603bb2464c) |
| **Deploy block** | `42535778` |
| **Deploy gas price** | `0.011 gwei` |
| **Source verified** | [Blockscout](https://base-sepolia.blockscout.com/address/0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9?tab=contract) ✓ · [Sourcify (perfect match)](https://repo.sourcify.dev/contracts/full_match/84532/0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9/) ✓ |

### First anchor — `PR01JUN26` (binance, period 43)

| Item | Value |
|---|---|
| **Anchor tx** | [`0x3ce248b76d7638ea3326b93b6ef731fa40eb07f52c8397ab00633079614932bb`](https://sepolia.basescan.org/tx/0x3ce248b76d7638ea3326b93b6ef731fa40eb07f52c8397ab00633079614932bb) |
| **Block** | `42535860` |
| **Gas price** | `0.011 gwei` |
| **snapshotId** | `PR01JUN26` |
| **periodSeq** | `43` |
| **snapshotTime** | `2026-01-06T00:00:00Z` (unix `1767657600`) |
| **btcBlockHeight** | `951913` |
| **exchangeMerkleRoot** | `0x250e0c4ab441780d7276bd63a7e2bfb098213dd6bd9de89265923d8e3c11c2d1` |
| **artifactBundleRoot** | `0xf452e47dd22ed63dc4a905fe79da6c6f7a6975cc0d775f50d01879e97616671f` |
| **verificationBundleRoot** | `0x7be6a4a685aae652c4072a88130f3f21621854a1888175c80300ad040eb2d249` |
| **verdictSummary** | `0x09` |
| **coverageBps** | `10000` |
| **Local bundles** | `artifacts/PR01JUN26.{artifact-bundle,verification-bundle,anchor}.json` |

Verifier outcomes at anchor time:

- `internal-consistency@1.0` → **PASS**
- `onchain-balance-hot@1.0` → **FAIL** (coverage 4.04%)
- 4 stub verifiers → UNVERIFIABLE

### Gate.io public-data anchor — `20260316`

| Item | Value |
|---|---|
| **Anchor tx** | [`0xd3e47684a01c9c23b1086da459ec718f069567f763b76944997b337fd2b3b0c4`](https://sepolia.basescan.org/tx/0xd3e47684a01c9c23b1086da459ec718f069567f763b76944997b337fd2b3b0c4) |
| **Block** | `42816341` |
| **Gas price** | `0.006 gwei` |
| **snapshotId** | `20260316` |
| **periodSeq** | `0` |
| **snapshotTime** | `2026-03-16T00:00:00Z` (unix `1773619200`) |
| **btcBlockHeight** | `0` |
| **exchangeMerkleRoot** | `0x22299d113e6a4336509f35a9404025ba2a5a2274dc3414d69d37aa3198e843fe` |
| **artifactBundleRoot** | `0x34bd742f73b49f32902d047f3847ccd9fa1eb306ee8d9469cb8aaad9974fa65a` |
| **verificationBundleRoot** | `0x504fcc492ccb9fa6c41ebcb376f929be63e1c769488b432d1f0f3faaaea88505` |
| **verdictSummary** | `0x04` |
| **coverageBps** | `10000` |
| **Local bundles** | `artifacts/gateio/20260316/bundles/20260316.{artifact-bundle,verification-bundle,anchor}.json` |

Verifier outcomes at anchor time:

- `artifact-integrity@1.0` → **PASS**
- `solvency-claim@1.0` → **PASS** (self-reported public reserve ratios only)
- Gate.io no-login reserve-side / zk / user-inclusion verifiers → **UNVERIFIABLE**

This anchor records a Gate.io public-data assessment. It does not independently
prove Gate.io on-chain reserves or user liability inclusion.

### Explorer quick links

- **Contract (Basescan)**: https://sepolia.basescan.org/address/0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9
- **Contract (Blockscout, verified source)**: https://base-sepolia.blockscout.com/address/0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9?tab=contract
- **Deploy tx**: https://sepolia.basescan.org/tx/0x3e844a1099f932638f695ce0d3045f78e5cc4e5def63d0edfb40bb603bb2464c
- **First anchor tx**: https://sepolia.basescan.org/tx/0x3ce248b76d7638ea3326b93b6ef731fa40eb07f52c8397ab00633079614932bb
- **SnapshotAnchored event logs**: https://sepolia.basescan.org/address/0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9#events

### Basescan source verify

Blockscout + Sourcify verification succeeded (links above). The [Basescan contract page](https://sepolia.basescan.org/address/0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9#code) may still show bytecode-only until a Basescan API submission completes — use **Cross-chain Verify** on that page (Sourcify match is already indexed), or run:

```bash
export BASESCAN_API_KEY=<your key>   # https://sepolia.basescan.org/myapikey
export ANCHOR_SIGNER=0xf2674A2b11b4a6CedC94ab57b22c86Df1fF36209

cd contracts
forge verify-contract 0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9 \
  src/ArdmerePoRAnchor.sol:ArdmerePoRAnchor \
  --chain base-sepolia \
  --constructor-args $(cast abi-encode "constructor(address)" $ANCHOR_SIGNER) \
  --watch
```

## Base mainnet (production)

Not deployed yet. See [ADR-002](./decisions.md#adr-002-primary-anchor-chain-base).

# `contracts/` — ArdmerePoRAnchor

Foundry project for the [ArdmerePoRAnchor](./src/ArdmerePoRAnchor.sol) contract.

## Design (one paragraph)

A non-upgradable, stateless event emitter. Holds one immutable `signer`
address and exposes a single permissioned function `anchorSnapshot(...)` which
emits one `SnapshotAnchored` event per exchange snapshot period, carrying
**both** the artifact bundle root and the verification bundle root in a
single transaction. No admin, no proxies, no storage writes. This is the
"minimum trusted kernel" that backs ardmere's PoR claim provenance — see
[`docs/verifier-architecture.md` §6](../docs/verifier-architecture.md).

## Deployed (Base Sepolia)

| | |
|---|---|
| **Contract** | [`0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9`](https://sepolia.basescan.org/address/0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9) |
| **Deploy tx** | [`0x3e844a1099f932638f695ce0d3045f78e5cc4e5def63d0edfb40bb603bb2464c`](https://sepolia.basescan.org/tx/0x3e844a1099f932638f695ce0d3045f78e5cc4e5def63d0edfb40bb603bb2464c) |
| **First anchor tx** | [`0x3ce248b76d7638ea3326b93b6ef731fa40eb07f52c8397ab00633079614932bb`](https://sepolia.basescan.org/tx/0x3ce248b76d7638ea3326b93b6ef731fa40eb07f52c8397ab00633079614932bb) |
| **Verified source** | [Blockscout](https://base-sepolia.blockscout.com/address/0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9?tab=contract) · [Sourcify](https://repo.sourcify.dev/contracts/full_match/84532/0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9/) |

Full deployment record: [`docs/deployments.md`](../docs/deployments.md).

## Build and test

```bash
forge build       # compile (Solc 0.8.26)
forge test -vv    # unit tests
```

## Deploy to Base Sepolia

```bash
# 1) Configure
#    - Secrets: export PRIVATE_KEY and ANCHOR_SIGNER in ~/.zshrc (testnet only)
#    - Project: cp ../.env.example ../.env for RPC / contract addresses
cp ../.env.example ../.env
ln -sf ../.env .env   # Foundry loads .env from this directory

set -a && source .env && set +a
# Required in shell: PRIVATE_KEY, ANCHOR_SIGNER
# Required in .env: BASE_SEPOLIA_RPC_URL

# 2) Fund the address with Base Sepolia ETH from a faucet
#    e.g. https://www.alchemy.com/faucets/base-sepolia or
#         https://faucets.chain.link/base-sepolia

# 3) Deploy (example: 0.011 gwei gas price)
forge script script/Deploy.s.sol:Deploy \
  --rpc-url base_sepolia \
  --broadcast \
  --legacy \
  --with-gas-price 11000000 \
  -vvvv
```

## Verify on Basescan

Source verified on Blockscout + Sourcify (see links above). To also verify on Basescan:

```bash
export BASESCAN_API_KEY=<your key>
export ANCHOR_SIGNER=0xf2674A2b11b4a6CedC94ab57b22c86Df1fF36209

forge verify-contract 0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9 \
  src/ArdmerePoRAnchor.sol:ArdmerePoRAnchor \
  --chain base-sepolia \
  --constructor-args $(cast abi-encode "constructor(address)" $ANCHOR_SIGNER) \
  --watch
```

Or verify via Blockscout (no API key):

```bash
forge verify-contract 0xFb55EEaAd312C2564B64002Bd3DC9D922Bb7eeF9 \
  src/ArdmerePoRAnchor.sol:ArdmerePoRAnchor \
  --chain-id 84532 \
  --verifier blockscout \
  --verifier-url https://base-sepolia.blockscout.com/api \
  --etherscan-api-key '' \
  --constructor-args $(cast abi-encode "constructor(address)" $ANCHOR_SIGNER) \
  --watch
```

## Send an anchor

Run `por anchor` first (see [`../docs/por-cli.md`](../docs/por-cli.md)), then broadcast the printed `cast send` at your chosen gas price.

## Mainnet (later)

When ready for production, repeat with `--rpc-url base` and a separately-managed
production signer (see [`docs/decisions.md` ADR-002](../docs/decisions.md)).

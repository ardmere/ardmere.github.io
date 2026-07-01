# `contracts/` — ArdmerePoRAnchor

Foundry project for the [ArdmerePoRAnchor](./src/ArdmerePoRAnchor.sol) contract.

## Design (one paragraph)

A UUPS-upgradeable anchor contract deployed behind an `ERC1967Proxy`. Holds a
permissioned `signer` address and exposes `anchorSnapshot(...)` which stores an
immutable `SnapshotRecord` and emits one `SnapshotAnchored` event per exchange
snapshot period, carrying **both** the artifact bundle root and the verification
bundle root in a single transaction. View functions (`getSnapshot`, `getLatestSnapshot`, …)
enable O(1) reads without log scanning (schema v3). The deployer owns upgrade
rights; `setSigner` allows signer rotation without redeploying the proxy. This is
the "minimum trusted kernel" that backs ardmere's PoR claim provenance — see
[`docs/verifier-architecture.md` §6](../docs/verifier-architecture.md).

## License

Ardmere-authored Solidity in this directory (`src/`, `script/`, `test/`) is released under the [MIT License](../LICENSE) at the repository root (`SPDX-License-Identifier: MIT` in each file). Dependencies under `lib/` are vendored via Foundry (`forge install`) and retain their upstream licenses.

## Public source

Canonical source: [github.com/ardmere/ardmere.github.io](https://github.com/ardmere/ardmere.github.io) — `contracts/src/ArdmerePoRAnchor.sol`.

Full deployment record: [`docs/deployments.md`](../docs/deployments.md).

## Build and test

```bash
forge build       # compile (Solc 0.8.26)
forge test -vv    # unit tests
```

## Deploy to Ethereum Sepolia

```bash
# 1) Configure
#    - Secrets: export PRIVATE_KEY and ANCHOR_SIGNER in ~/.zshenv (testnet only)
#    - Project: cp ../.env.example ../.env for RPC / contract addresses
cp ../.env.example ../.env
ln -sf ../.env .env   # Foundry loads .env from this directory

set -a && source .env && set +a
source ~/.zshenv
# Required in shell: PRIVATE_KEY, ANCHOR_SIGNER
# Required in .env: SEPOLIA_RPC_URL

# 2) Fund the deployer address with Sepolia ETH from a faucet
#    e.g. https://cloud.google.com/application/web3/faucet/ethereum/sepolia

# 3) Deploy implementation + ERC1967 proxy
forge script script/Deploy.s.sol:Deploy \
  --rpc-url sepolia \
  --broadcast \
  -vvvv
```

Set `ANCHOR_CONTRACT` in `.env` to the **proxy** address logged by the script.

## Upgrade existing proxy (v2 → v3)

When a new implementation is released (e.g. on-chain query storage), upgrade the deployed proxy without changing its address:

```bash
# Required: PRIVATE_KEY (upgrade owner), ANCHOR_CONTRACT (proxy)
source ~/.zshenv
set -a && source .env && set +a

forge script script/Upgrade.s.sol:Upgrade \
  --rpc-url sepolia \
  --broadcast \
  -vvvv
```

Verify the new implementation on Etherscan (replace `<NEW_IMPLEMENTATION>`):

```bash
forge verify-contract <NEW_IMPLEMENTATION> \
  src/ArdmerePoRAnchor.sol:ArdmerePoRAnchor \
  --chain sepolia \
  --watch
```

Pre-upgrade anchors remain in event logs only; anchors after upgrade populate on-chain storage. See [anchor-query-api.md](../docs/anchor-query-api.md).

## Verify on Etherscan (Sepolia)

```bash
export ETHERSCAN_API_KEY=<your key>
export ANCHOR_SIGNER=<signer address>

# Implementation (no constructor args beyond disabled initializers)
forge verify-contract <IMPLEMENTATION_ADDRESS> \
  src/ArdmerePoRAnchor.sol:ArdmerePoRAnchor \
  --chain sepolia \
  --watch

# Proxy
forge verify-contract <PROXY_ADDRESS> \
  lib/openzeppelin-contracts/contracts/proxy/ERC1967/ERC1967Proxy.sol:ERC1967Proxy \
  --chain sepolia \
  --constructor-args $(cast abi-encode "constructor(address,bytes)" <IMPLEMENTATION_ADDRESS> $(cast calldata "initialize(address)" $ANCHOR_SIGNER)) \
  --watch
```

## Verify on Sourcify (optional)

```bash
# Implementation
ETHERSCAN_API_KEY= BASESCAN_API_KEY= forge verify-contract <IMPLEMENTATION_ADDRESS> \
  src/ArdmerePoRAnchor.sol:ArdmerePoRAnchor \
  --chain sepolia \
  --verifier sourcify \
  --verifier-url https://sourcify.dev/server/ \
  --watch
```

The proxy is usually verified on Etherscan as `ERC1967Proxy` (see commands above); Sourcify may show a runtime match for the proxy address.

## Send an anchor

Run `por anchor` first (see [`../docs/por-cli.md`](../docs/por-cli.md)), then broadcast the printed `cast send` against the **proxy** address.

## Mainnet (later)

When ready for production, repeat with `--rpc-url base` and a separately-managed
production signer (see [`docs/decisions.md` ADR-002](../docs/decisions.md)).

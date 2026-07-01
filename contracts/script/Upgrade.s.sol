// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script, console2} from "forge-std/Script.sol";
import {ArdmerePoRAnchor} from "../src/ArdmerePoRAnchor.sol";

/// @notice UUPS upgrade: deploy new implementation and point existing proxy at it.
///         Required env: PRIVATE_KEY, ANCHOR_CONTRACT (proxy address).
///         Run example:
///         forge script script/Upgrade.s.sol:Upgrade --rpc-url sepolia --broadcast -vvvv
contract Upgrade is Script {
    function run() external {
        uint256 pk = vm.envUint("PRIVATE_KEY");
        address proxy = vm.envAddress("ANCHOR_CONTRACT");

        vm.startBroadcast(pk);
        ArdmerePoRAnchor implementation = new ArdmerePoRAnchor();
        ArdmerePoRAnchor(proxy).upgradeToAndCall(address(implementation), "");
        vm.stopBroadcast();

        console2.log("New implementation:", address(implementation));
        console2.log("Proxy upgraded:", proxy);
        console2.log("SCHEMA_VERSION:", ArdmerePoRAnchor(proxy).SCHEMA_VERSION());
        console2.log("STORAGE_VERSION:", ArdmerePoRAnchor(proxy).STORAGE_VERSION());
    }
}

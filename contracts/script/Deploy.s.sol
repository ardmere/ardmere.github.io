// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script, console2} from "forge-std/Script.sol";
import {ArdmerePoRAnchor} from "../src/ArdmerePoRAnchor.sol";

/// @notice Deploys ArdmerePoRAnchor.
///         Required env: PRIVATE_KEY, ANCHOR_SIGNER (= address derived from PRIVATE_KEY).
///         Run example:
///         forge script script/Deploy.s.sol:Deploy --rpc-url base_sepolia --broadcast -vvvv
contract Deploy is Script {
    function run() external returns (ArdmerePoRAnchor deployed) {
        uint256 pk = vm.envUint("PRIVATE_KEY");
        address signer = vm.envAddress("ANCHOR_SIGNER");

        vm.startBroadcast(pk);
        deployed = new ArdmerePoRAnchor(signer);
        vm.stopBroadcast();

        console2.log("ArdmerePoRAnchor deployed at:", address(deployed));
        console2.log("Anchor signer:", deployed.signer());
    }
}

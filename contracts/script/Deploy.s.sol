// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script, console2} from "forge-std/Script.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {ArdmerePoRAnchor} from "../src/ArdmerePoRAnchor.sol";

/// @notice Deploys ArdmerePoRAnchor (UUPS implementation + ERC1967 proxy).
///         Required env: PRIVATE_KEY, ANCHOR_SIGNER (= address derived from PRIVATE_KEY).
///         Run example:
///         forge script script/Deploy.s.sol:Deploy --rpc-url sepolia --broadcast -vvvv
contract Deploy is Script {
    function run() external returns (ArdmerePoRAnchor proxy) {
        uint256 pk = vm.envUint("PRIVATE_KEY");
        address signer = vm.envAddress("ANCHOR_SIGNER");

        vm.startBroadcast(pk);
        ArdmerePoRAnchor implementation = new ArdmerePoRAnchor();
        bytes memory initData = abi.encodeCall(ArdmerePoRAnchor.initialize, (signer));
        ERC1967Proxy proxyContract = new ERC1967Proxy(address(implementation), initData);
        proxy = ArdmerePoRAnchor(address(proxyContract));
        vm.stopBroadcast();

        console2.log("Implementation:", address(implementation));
        console2.log("ArdmerePoRAnchor proxy (ANCHOR_CONTRACT):", address(proxy));
        console2.log("Anchor signer:", proxy.signer());
        console2.log("Upgrade owner:", proxy.owner());
    }
}

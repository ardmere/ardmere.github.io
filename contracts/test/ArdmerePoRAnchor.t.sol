// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ArdmerePoRAnchor} from "../src/ArdmerePoRAnchor.sol";

contract ArdmerePoRAnchorTest is Test {
    ArdmerePoRAnchor internal anchorContract;
    address internal constant SIGNER = address(0xBEEF);
    address internal constant STRANGER = address(0xCAFE);

    bytes32 internal constant SNAPSHOT_PR01JUN26 = keccak256("PR01JUN26");
    bytes32 internal constant EXCHANGE_TAG = keccak256("binance");
    bytes32 internal constant EXCHANGE_ROOT =
        0x250e0c4ab441780d7276bd63a7e2bfb098213dd6bd9de89265923d8e3c11c2d1;
    bytes32 internal constant ARTIFACT_ROOT =
        0xc5dacb50add81e87351b7fe38c18be4b5e0974c578c63c4a1670ecd2ea803c7e;
    bytes32 internal constant VERIFICATION_ROOT =
        0x7d124bc76311cb7f5021cafb696eeb1e6ec7d19ba35232644f41f05a196b81b5;

    event SnapshotAnchored(
        bytes32 indexed snapshotId,
        bytes32 indexed exchangeTag,
        string exchange,
        uint32 periodSeq,
        uint64 snapshotTime,
        uint32 btcBlockHeight,
        bytes32 exchangeMerkleRoot,
        bytes32 artifactBundleRoot,
        bytes32 verificationBundleRoot,
        uint8 verdictSummary,
        uint16 coverageBps,
        uint8 schemaVersion,
        uint256 anchoredAt
    );

    function setUp() public {
        anchorContract = new ArdmerePoRAnchor(SIGNER);
    }

    function test_constructor_setsSignerAndSchema() public view {
        assertEq(anchorContract.signer(), SIGNER);
        assertEq(anchorContract.SCHEMA_VERSION(), 2);
    }

    function test_anchorSnapshot_signerCanPublish() public {
        vm.warp(1_767_657_600);

        vm.expectEmit(true, true, false, true, address(anchorContract));
        emit SnapshotAnchored(
            SNAPSHOT_PR01JUN26,
            EXCHANGE_TAG,
            "binance",
            43,
            1_767_657_600,
            951_913,
            EXCHANGE_ROOT,
            ARTIFACT_ROOT,
            VERIFICATION_ROOT,
            0x01,
            17,
            2,
            1_767_657_600
        );

        vm.prank(SIGNER);
        anchorContract.anchorSnapshot(
            "binance",
            SNAPSHOT_PR01JUN26,
            43,
            1_767_657_600,
            951_913,
            EXCHANGE_ROOT,
            ARTIFACT_ROOT,
            VERIFICATION_ROOT,
            0x01,
            17
        );
    }

    function test_anchorSnapshot_revertsForStranger() public {
        vm.prank(STRANGER);
        vm.expectRevert(abi.encodeWithSelector(ArdmerePoRAnchor.Unauthorized.selector, STRANGER));
        anchorContract.anchorSnapshot(
            "binance",
            SNAPSHOT_PR01JUN26,
            1,
            0,
            0,
            bytes32(0),
            ARTIFACT_ROOT,
            VERIFICATION_ROOT,
            0,
            0
        );
    }

    function test_anchorSnapshot_revertsOnZeroArtifactRoot() public {
        vm.prank(SIGNER);
        vm.expectRevert(ArdmerePoRAnchor.InvalidRoot.selector);
        anchorContract.anchorSnapshot(
            "binance", SNAPSHOT_PR01JUN26, 1, 0, 0, bytes32(0), bytes32(0), VERIFICATION_ROOT, 0, 0
        );
    }

    function test_anchorSnapshot_revertsOnZeroVerificationRoot() public {
        vm.prank(SIGNER);
        vm.expectRevert(ArdmerePoRAnchor.InvalidRoot.selector);
        anchorContract.anchorSnapshot(
            "binance", SNAPSHOT_PR01JUN26, 1, 0, 0, bytes32(0), ARTIFACT_ROOT, bytes32(0), 0, 0
        );
    }

    function test_anchorSnapshot_revertsOnEmptyExchange() public {
        vm.prank(SIGNER);
        vm.expectRevert(ArdmerePoRAnchor.EmptyExchange.selector);
        anchorContract.anchorSnapshot(
            "", SNAPSHOT_PR01JUN26, 1, 0, 0, bytes32(0), ARTIFACT_ROOT, VERIFICATION_ROOT, 0, 0
        );
    }
}

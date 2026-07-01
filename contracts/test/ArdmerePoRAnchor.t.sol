// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {ArdmerePoRAnchor} from "../src/ArdmerePoRAnchor.sol";

contract ArdmerePoRAnchorTest is Test {
    ArdmerePoRAnchor internal anchorContract;
    ArdmerePoRAnchor internal implementation;
    address internal constant SIGNER = address(0xBEEF);
    address internal constant STRANGER = address(0xCAFE);
    address internal owner;

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

    event SignerUpdated(address indexed previousSigner, address indexed newSigner);

    function setUp() public {
        owner = address(this);
        implementation = new ArdmerePoRAnchor();
        bytes memory initData = abi.encodeCall(ArdmerePoRAnchor.initialize, (SIGNER));
        ERC1967Proxy proxy = new ERC1967Proxy(address(implementation), initData);
        anchorContract = ArdmerePoRAnchor(address(proxy));
    }

    function test_initialize_setsSignerSchemaAndOwner() public view {
        assertEq(anchorContract.signer(), SIGNER);
        assertEq(anchorContract.SCHEMA_VERSION(), 3);
        assertEq(anchorContract.STORAGE_VERSION(), 1);
        assertEq(anchorContract.owner(), owner);
    }

    function test_initialize_revertsOnZeroSigner() public {
        ArdmerePoRAnchor freshImpl = new ArdmerePoRAnchor();
        bytes memory initData = abi.encodeCall(ArdmerePoRAnchor.initialize, (address(0)));
        vm.expectRevert(ArdmerePoRAnchor.ZeroAddress.selector);
        new ERC1967Proxy(address(freshImpl), initData);
    }

    function test_setSigner_onlyOwner() public {
        address newSigner = address(0xDEAD);

        vm.expectEmit(true, true, false, true, address(anchorContract));
        emit SignerUpdated(SIGNER, newSigner);
        anchorContract.setSigner(newSigner);

        assertEq(anchorContract.signer(), newSigner);
    }

    function test_setSigner_revertsForNonOwner() public {
        vm.prank(STRANGER);
        vm.expectRevert();
        anchorContract.setSigner(address(0xDEAD));
    }

    function test_setSigner_revertsOnZeroAddress() public {
        vm.expectRevert(ArdmerePoRAnchor.ZeroAddress.selector);
        anchorContract.setSigner(address(0));
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
            3,
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

    function _anchorBinance(
        bytes32 snapshotId,
        uint32 periodSeq,
        uint64 snapshotTime,
        uint32 btcBlockHeight
    ) internal {
        vm.prank(SIGNER);
        anchorContract.anchorSnapshot(
            "binance",
            snapshotId,
            periodSeq,
            snapshotTime,
            btcBlockHeight,
            EXCHANGE_ROOT,
            ARTIFACT_ROOT,
            VERIFICATION_ROOT,
            0x01,
            17
        );
    }

    function test_anchorSnapshot_writesStorageAndQueries() public {
        vm.warp(1_767_657_600);
        _anchorBinance(SNAPSHOT_PR01JUN26, 43, 1_767_657_600, 951_913);

        assertTrue(anchorContract.snapshotExists(SNAPSHOT_PR01JUN26));

        ArdmerePoRAnchor.SnapshotRecord memory record = anchorContract.getSnapshot(SNAPSHOT_PR01JUN26);
        assertEq(record.snapshotId, SNAPSHOT_PR01JUN26);
        assertEq(record.exchangeTag, EXCHANGE_TAG);
        assertEq(record.exchange, "binance");
        assertEq(record.periodSeq, 43);
        assertEq(record.snapshotTime, 1_767_657_600);
        assertEq(record.btcBlockHeight, 951_913);
        assertEq(record.exchangeMerkleRoot, EXCHANGE_ROOT);
        assertEq(record.artifactBundleRoot, ARTIFACT_ROOT);
        assertEq(record.verificationBundleRoot, VERIFICATION_ROOT);
        assertEq(record.verdictSummary, 0x01);
        assertEq(record.coverageBps, 17);
        assertEq(record.schemaVersion, 3);
        assertEq(record.anchoredAt, 1_767_657_600);

        ArdmerePoRAnchor.SnapshotRecord memory latest = anchorContract.getLatestSnapshot(EXCHANGE_TAG);
        assertEq(latest.snapshotId, SNAPSHOT_PR01JUN26);

        ArdmerePoRAnchor.SnapshotRecord memory byPeriod =
            anchorContract.getSnapshotByPeriod(EXCHANGE_TAG, 43);
        assertEq(byPeriod.snapshotId, SNAPSHOT_PR01JUN26);

        assertEq(anchorContract.getSnapshotCount(EXCHANGE_TAG), 1);
        assertEq(anchorContract.getSnapshotIdAt(EXCHANGE_TAG, 0), SNAPSHOT_PR01JUN26);
    }

    function test_anchorSnapshot_latestSnapshotUpdates() public {
        bytes32 first = keccak256("PR01MAY26");
        bytes32 second = keccak256("PR01JUN26");

        vm.warp(1_767_657_600);
        _anchorBinance(first, 42, 1_767_657_600, 951_000);

        vm.warp(1_768_032_000);
        _anchorBinance(second, 43, 1_768_032_000, 952_000);

        assertEq(anchorContract.getLatestSnapshot(EXCHANGE_TAG).snapshotId, second);
        assertEq(anchorContract.getSnapshotCount(EXCHANGE_TAG), 2);
        assertEq(anchorContract.getSnapshotIdAt(EXCHANGE_TAG, 0), first);
        assertEq(anchorContract.getSnapshotIdAt(EXCHANGE_TAG, 1), second);
    }

    function test_anchorSnapshot_revertsOnDuplicateSnapshotId() public {
        vm.warp(1_767_657_600);
        _anchorBinance(SNAPSHOT_PR01JUN26, 43, 1_767_657_600, 951_913);

        vm.prank(SIGNER);
        vm.expectRevert(abi.encodeWithSelector(ArdmerePoRAnchor.SnapshotAlreadyExists.selector, SNAPSHOT_PR01JUN26));
        anchorContract.anchorSnapshot(
            "binance",
            SNAPSHOT_PR01JUN26,
            44,
            1_768_032_000,
            952_000,
            EXCHANGE_ROOT,
            ARTIFACT_ROOT,
            VERIFICATION_ROOT,
            0x01,
            17
        );
    }

    function test_getSnapshot_revertsWhenMissing() public {
        vm.expectRevert(abi.encodeWithSelector(ArdmerePoRAnchor.SnapshotNotFound.selector, SNAPSHOT_PR01JUN26));
        anchorContract.getSnapshot(SNAPSHOT_PR01JUN26);
    }

    function test_getSnapshotIdAt_revertsOnOutOfBounds() public {
        vm.expectRevert(
            abi.encodeWithSelector(ArdmerePoRAnchor.SnapshotIndexOutOfBounds.selector, EXCHANGE_TAG, 0)
        );
        anchorContract.getSnapshotIdAt(EXCHANGE_TAG, 0);
    }

    function test_uupsUpgrade_preservesSignerAndEnablesStorage() public {
        vm.warp(1_767_657_600);
        _anchorBinance(SNAPSHOT_PR01JUN26, 43, 1_767_657_600, 951_913);

        ArdmerePoRAnchor newImplementation = new ArdmerePoRAnchor();
        anchorContract.upgradeToAndCall(address(newImplementation), "");

        assertEq(anchorContract.signer(), SIGNER);
        assertEq(anchorContract.SCHEMA_VERSION(), 3);
        assertTrue(anchorContract.snapshotExists(SNAPSHOT_PR01JUN26));
        assertEq(anchorContract.getLatestSnapshot(EXCHANGE_TAG).snapshotId, SNAPSHOT_PR01JUN26);
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

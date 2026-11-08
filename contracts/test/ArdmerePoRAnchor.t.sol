// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {ArdmerePoRAnchor} from "../src/ArdmerePoRAnchor.sol";

/// @dev Minimal v3-shaped implementation used only to seed storage before a v4 upgrade.
contract ArdmerePoRAnchorV3Seed {
    uint8 public constant SCHEMA_VERSION = 3;
    uint8 public constant STORAGE_VERSION = 1;

    address public signer;

    struct SnapshotRecord {
        bytes32 snapshotId;
        bytes32 exchangeTag;
        string exchange;
        uint32 periodSeq;
        uint64 snapshotTime;
        uint32 btcBlockHeight;
        bytes32 exchangeMerkleRoot;
        bytes32 artifactBundleRoot;
        bytes32 verificationBundleRoot;
        uint8 verdictSummary;
        uint16 coverageBps;
        uint8 schemaVersion;
        uint256 anchoredAt;
    }

    /// @custom:storage-location erc7201:ardmere.storage.ArdmerePoRAnchor
    struct ArdmerePoRAnchorStorage {
        mapping(bytes32 snapshotId => SnapshotRecord) snapshots;
        mapping(bytes32 exchangeTag => bytes32[]) snapshotIds;
        mapping(bytes32 exchangeTag => mapping(uint32 periodSeq => bytes32)) snapshotByPeriod;
        mapping(bytes32 exchangeTag => bytes32) latestSnapshotId;
        mapping(bytes32 exchangeTag => uint256) snapshotCount;
    }

    bytes32 private constant ArdmerePoRAnchorStorageLocation =
        0x168b5695e07d5cc77c8fba711207364d2cbe159a40632dd678b45ab21ab90100;

    error Unauthorized(address caller);
    error InvalidRoot();
    error EmptyExchange();
    error ZeroAddress();
    error SnapshotAlreadyExists(bytes32 snapshotId);
    error PeriodAlreadyAnchored(bytes32 exchangeTag, uint32 periodSeq);

    function _getStorage() private pure returns (ArdmerePoRAnchorStorage storage $) {
        assembly {
            $.slot := ArdmerePoRAnchorStorageLocation
        }
    }

    function initialize(address _signer) external {
        if (_signer == address(0)) revert ZeroAddress();
        signer = _signer;
    }

    function anchorSnapshot(
        string calldata exchange,
        bytes32 snapshotId,
        uint32 periodSeq,
        uint64 snapshotTime,
        uint32 btcBlockHeight,
        bytes32 exchangeMerkleRoot,
        bytes32 artifactBundleRoot,
        bytes32 verificationBundleRoot,
        uint8 verdictSummary,
        uint16 coverageBps
    ) external {
        if (msg.sender != signer) revert Unauthorized(msg.sender);
        if (bytes(exchange).length == 0) revert EmptyExchange();
        if (artifactBundleRoot == bytes32(0) || verificationBundleRoot == bytes32(0)) {
            revert InvalidRoot();
        }
        ArdmerePoRAnchorStorage storage $ = _getStorage();
        if ($.snapshots[snapshotId].anchoredAt != 0) revert SnapshotAlreadyExists(snapshotId);
        bytes32 exchangeTag = keccak256(bytes(exchange));
        if ($.snapshotByPeriod[exchangeTag][periodSeq] != bytes32(0)) {
            revert PeriodAlreadyAnchored(exchangeTag, periodSeq);
        }
        SnapshotRecord storage record = $.snapshots[snapshotId];
        record.snapshotId = snapshotId;
        record.exchangeTag = exchangeTag;
        record.exchange = exchange;
        record.periodSeq = periodSeq;
        record.snapshotTime = snapshotTime;
        record.btcBlockHeight = btcBlockHeight;
        record.exchangeMerkleRoot = exchangeMerkleRoot;
        record.artifactBundleRoot = artifactBundleRoot;
        record.verificationBundleRoot = verificationBundleRoot;
        record.verdictSummary = verdictSummary;
        record.coverageBps = coverageBps;
        record.schemaVersion = SCHEMA_VERSION;
        record.anchoredAt = block.timestamp;
        $.snapshotIds[exchangeTag].push(snapshotId);
        $.snapshotByPeriod[exchangeTag][periodSeq] = snapshotId;
        $.latestSnapshotId[exchangeTag] = snapshotId;
        $.snapshotCount[exchangeTag] += 1;
    }

    function upgradeToAndCall(address newImplementation, bytes memory data) external {
        // Minimal UUPS stand-in for tests: write ERC-1967 implementation slot.
        bytes32 slot = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
        assembly {
            sstore(slot, newImplementation)
        }
        if (data.length > 0) {
            (bool ok,) = newImplementation.delegatecall(data);
            require(ok);
        }
    }
}

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
        uint256 anchoredAt,
        uint8 porStage,
        bool effectivePoR,
        uint16 walletCoverageBps,
        uint16 reserveCoverageBps
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
        assertEq(anchorContract.SCHEMA_VERSION(), 4);
        assertEq(anchorContract.STORAGE_VERSION(), 2);
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
            4,
            1_767_657_600,
            1, // Stage 0
            false,
            2500,
            8000
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
            17,
            1,
            false,
            2500,
            8000
        );
    }

    function _anchorBinance(
        bytes32 snapshotId,
        uint32 periodSeq,
        uint64 snapshotTime,
        uint32 btcBlockHeight,
        uint8 porStage,
        bool effectivePoR,
        uint16 walletCoverageBps,
        uint16 reserveCoverageBps
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
            17,
            porStage,
            effectivePoR,
            walletCoverageBps,
            reserveCoverageBps
        );
    }

    function test_anchorSnapshot_writesStorageAndQueries() public {
        vm.warp(1_767_657_600);
        _anchorBinance(SNAPSHOT_PR01JUN26, 43, 1_767_657_600, 951_913, 2, true, 5000, 9000);

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
        assertEq(record.schemaVersion, 4);
        assertEq(record.anchoredAt, 1_767_657_600);
        assertEq(record.porStage, 2);
        assertTrue(record.effectivePoR);
        assertEq(record.walletCoverageBps, 5000);
        assertEq(record.reserveCoverageBps, 9000);

        ArdmerePoRAnchor.SnapshotRecord memory latest = anchorContract.getLatestSnapshot(EXCHANGE_TAG);
        assertEq(latest.snapshotId, SNAPSHOT_PR01JUN26);

        ArdmerePoRAnchor.SnapshotRecord memory byPeriod =
            anchorContract.getSnapshotByPeriod(EXCHANGE_TAG, 43);
        assertEq(byPeriod.snapshotId, SNAPSHOT_PR01JUN26);

        assertEq(anchorContract.getSnapshotCount(EXCHANGE_TAG), 1);
        assertEq(anchorContract.getSnapshotIdAt(EXCHANGE_TAG, 0), SNAPSHOT_PR01JUN26);
    }

    function test_anchorSnapshot_abiEncodeDecodeRoundTrip() public {
        vm.warp(1_767_657_600);
        _anchorBinance(SNAPSHOT_PR01JUN26, 43, 1_767_657_600, 951_913, 3, true, 10000, 7500);

        bytes memory encoded = abi.encode(anchorContract.getSnapshot(SNAPSHOT_PR01JUN26));
        ArdmerePoRAnchor.SnapshotRecord memory decoded =
            abi.decode(encoded, (ArdmerePoRAnchor.SnapshotRecord));

        assertEq(decoded.snapshotId, SNAPSHOT_PR01JUN26);
        assertEq(decoded.coverageBps, 17);
        assertEq(decoded.schemaVersion, 4);
        assertEq(decoded.porStage, 3);
        assertTrue(decoded.effectivePoR);
        assertEq(decoded.walletCoverageBps, 10000);
        assertEq(decoded.reserveCoverageBps, 7500);
        assertEq(decoded.anchoredAt, 1_767_657_600);
    }

    function test_anchorSnapshot_latestSnapshotUpdates() public {
        bytes32 first = keccak256("PR01MAY26");
        bytes32 second = keccak256("PR01JUN26");

        vm.warp(1_767_657_600);
        _anchorBinance(first, 42, 1_767_657_600, 951_000, 1, false, 0, 0);

        vm.warp(1_768_032_000);
        _anchorBinance(second, 43, 1_768_032_000, 952_000, 2, true, 100, 200);

        assertEq(anchorContract.getLatestSnapshot(EXCHANGE_TAG).snapshotId, second);
        assertEq(anchorContract.getSnapshotCount(EXCHANGE_TAG), 2);
        assertEq(anchorContract.getSnapshotIdAt(EXCHANGE_TAG, 0), first);
        assertEq(anchorContract.getSnapshotIdAt(EXCHANGE_TAG, 1), second);
    }

    function test_anchorSnapshot_revertsOnDuplicateSnapshotId() public {
        vm.warp(1_767_657_600);
        _anchorBinance(SNAPSHOT_PR01JUN26, 43, 1_767_657_600, 951_913, 1, false, 0, 0);

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
            17,
            1,
            false,
            0,
            0
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
        _anchorBinance(SNAPSHOT_PR01JUN26, 43, 1_767_657_600, 951_913, 2, true, 1111, 2222);

        ArdmerePoRAnchor newImplementation = new ArdmerePoRAnchor();
        anchorContract.upgradeToAndCall(address(newImplementation), "");

        assertEq(anchorContract.signer(), SIGNER);
        assertEq(anchorContract.SCHEMA_VERSION(), 4);
        assertTrue(anchorContract.snapshotExists(SNAPSHOT_PR01JUN26));
        assertEq(anchorContract.getLatestSnapshot(EXCHANGE_TAG).snapshotId, SNAPSHOT_PR01JUN26);
        assertEq(anchorContract.getSnapshot(SNAPSHOT_PR01JUN26).walletCoverageBps, 1111);
    }

    function test_upgradeFromV3Storage_oldRecordReadableWithZeroNewFields() public {
        // Deploy a v3-shaped seed behind a proxy, write one record, then point at v4.
        ArdmerePoRAnchorV3Seed seedImpl = new ArdmerePoRAnchorV3Seed();
        // Minimal proxy: store implementation and call initialize via delegatecall manually.
        // Use ERC1967Proxy with empty init then set signer via raw storage is heavy —
        // instead deploy seed as the first implementation with a custom initializer payload.
        bytes memory initData = abi.encodeCall(ArdmerePoRAnchorV3Seed.initialize, (SIGNER));
        // ERC1967Proxy expects an initializer that may revert Ownable — seed has no Ownable.
        // Use a plain proxy pattern: deploy seed logic at proxy by constructing ERC1967 with init.
        ERC1967Proxy proxy = new ERC1967Proxy(address(seedImpl), initData);
        ArdmerePoRAnchorV3Seed seed = ArdmerePoRAnchorV3Seed(address(proxy));

        vm.warp(1_700_000_000);
        vm.prank(SIGNER);
        seed.anchorSnapshot(
            "binance",
            SNAPSHOT_PR01JUN26,
            43,
            1_700_000_000,
            900_000,
            EXCHANGE_ROOT,
            ARTIFACT_ROOT,
            VERIFICATION_ROOT,
            0x09,
            8500
        );

        ArdmerePoRAnchor v4 = new ArdmerePoRAnchor();
        // Replace implementation slot (seed's upgradeToAndCall).
        seed.upgradeToAndCall(address(v4), "");

        ArdmerePoRAnchor upgraded = ArdmerePoRAnchor(address(proxy));
        assertEq(upgraded.SCHEMA_VERSION(), 4);
        assertEq(upgraded.STORAGE_VERSION(), 2);
        assertTrue(upgraded.snapshotExists(SNAPSHOT_PR01JUN26));

        ArdmerePoRAnchor.SnapshotRecord memory record = upgraded.getSnapshot(SNAPSHOT_PR01JUN26);
        assertEq(record.snapshotId, SNAPSHOT_PR01JUN26);
        assertEq(record.exchange, "binance");
        assertEq(record.coverageBps, 8500);
        assertEq(record.schemaVersion, 3); // written under v3
        assertEq(record.anchoredAt, 1_700_000_000);
        assertEq(record.porStage, 0);
        assertFalse(record.effectivePoR);
        assertEq(record.walletCoverageBps, 0);
        assertEq(record.reserveCoverageBps, 0);
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
            0,
            0,
            false,
            0,
            0
        );
    }

    function test_anchorSnapshot_revertsOnZeroArtifactRoot() public {
        vm.prank(SIGNER);
        vm.expectRevert(ArdmerePoRAnchor.InvalidRoot.selector);
        anchorContract.anchorSnapshot(
            "binance",
            SNAPSHOT_PR01JUN26,
            1,
            0,
            0,
            bytes32(0),
            bytes32(0),
            VERIFICATION_ROOT,
            0,
            0,
            0,
            false,
            0,
            0
        );
    }

    function test_anchorSnapshot_revertsOnZeroVerificationRoot() public {
        vm.prank(SIGNER);
        vm.expectRevert(ArdmerePoRAnchor.InvalidRoot.selector);
        anchorContract.anchorSnapshot(
            "binance",
            SNAPSHOT_PR01JUN26,
            1,
            0,
            0,
            bytes32(0),
            ARTIFACT_ROOT,
            bytes32(0),
            0,
            0,
            0,
            false,
            0,
            0
        );
    }

    function test_anchorSnapshot_revertsOnEmptyExchange() public {
        vm.prank(SIGNER);
        vm.expectRevert(ArdmerePoRAnchor.EmptyExchange.selector);
        anchorContract.anchorSnapshot(
            "",
            SNAPSHOT_PR01JUN26,
            1,
            0,
            0,
            bytes32(0),
            ARTIFACT_ROOT,
            VERIFICATION_ROOT,
            0,
            0,
            0,
            false,
            0,
            0
        );
    }
}

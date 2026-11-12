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

/// @dev Schema v6 row (composite snapshotId, custody score, no USD totals) for backfill tests.
contract ArdmerePoRAnchorV6Seed {
    uint8 public constant SCHEMA_VERSION = 6;

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
        uint8 porStage;
        bool effectivePoR;
        uint16 walletCoverageBps;
        uint16 reserveCoverageBps;
        uint8 custodyRiskScore;
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

    function _getStorage() private pure returns (ArdmerePoRAnchorStorage storage $) {
        assembly {
            $.slot := ArdmerePoRAnchorStorageLocation
        }
    }

    function initialize(address _signer) external {
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
        uint16 coverageBps,
        uint8 porStage,
        bool effectivePoR,
        uint16 walletCoverageBps,
        uint16 reserveCoverageBps,
        uint8 custodyRiskScore
    ) external {
        ArdmerePoRAnchorStorage storage $ = _getStorage();
        bytes32 exchangeTag = keccak256(bytes(exchange));
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
        record.porStage = porStage;
        record.effectivePoR = effectivePoR;
        record.walletCoverageBps = walletCoverageBps;
        record.reserveCoverageBps = reserveCoverageBps;
        record.custodyRiskScore = custodyRiskScore;
        $.snapshotIds[exchangeTag].push(snapshotId);
        $.snapshotByPeriod[exchangeTag][periodSeq] = snapshotId;
        $.latestSnapshotId[exchangeTag] = snapshotId;
        $.snapshotCount[exchangeTag] += 1;
    }

    function upgradeToAndCall(address newImplementation, bytes memory data) external {
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
        uint16 reserveCoverageBps,
        uint8 custodyRiskScore,
        uint128 totalLiabilityUsd,
        uint128 totalReserveUsd,
        uint16 btcReserveRatioBps,
        uint16 ethReserveRatioBps,
        uint16 usdtReserveRatioBps,
        uint16 usdcReserveRatioBps
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
        assertEq(anchorContract.SCHEMA_VERSION(), 7);
        assertEq(anchorContract.STORAGE_VERSION(), 4);
        assertEq(anchorContract.USD_UNAVAILABLE(), 0);
        assertEq(anchorContract.RESERVE_RATIO_UNAVAILABLE(), 65535);
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
        bytes32 composite = anchorContract.expectedSnapshotId("binance", "PR01JUN26");

        vm.expectEmit(true, true, false, true, address(anchorContract));
        emit SnapshotAnchored(
            composite,
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
            7,
            1_767_657_600,
            1, // Stage 0
            false,
            2500,
            8000,
            42,
            1_000_000,
            1_100_000,
            10000,
            11000,
            65535,
            10150
        );

        vm.prank(SIGNER);
        anchorContract.anchorSnapshot(
            "binance",
            "PR01JUN26",
            composite,
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
            8000,
            42,
            1_000_000,
            1_100_000,
            10000,
            11000,
            65535,
            10150
        );
    }

    function _anchorBinance(
        string memory auditId,
        uint32 periodSeq,
        uint64 snapshotTime,
        uint32 btcBlockHeight,
        uint8 porStage,
        bool effectivePoR,
        uint16 walletCoverageBps,
        uint16 reserveCoverageBps,
        uint8 custodyRiskScore
    ) internal {
        bytes32 snapshotId = anchorContract.expectedSnapshotId("binance", auditId);
        vm.prank(SIGNER);
        anchorContract.anchorSnapshot(
            "binance",
            auditId,
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
            reserveCoverageBps,
            custodyRiskScore,
            0,
            0,
            65535,
            65535,
            65535,
            65535
        );
    }

    function test_anchorSnapshot_writesStorageAndQueries() public {
        vm.warp(1_767_657_600);
        bytes32 composite = anchorContract.expectedSnapshotId("binance", "PR01JUN26");
        _anchorBinance("PR01JUN26", 43, 1_767_657_600, 951_913, 2, true, 5000, 9000, 75);

        assertTrue(anchorContract.snapshotExists(composite));

        ArdmerePoRAnchor.SnapshotRecord memory record = anchorContract.getSnapshot(composite);
        assertEq(record.snapshotId, composite);
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
        assertEq(record.schemaVersion, 7);
        assertEq(record.totalLiabilityUsd, 0);
        assertEq(record.totalReserveUsd, 0);
        assertEq(record.btcReserveRatioBps, 65535);
        assertEq(record.custodyRiskScore, 75);
        assertEq(record.anchoredAt, 1_767_657_600);
        assertEq(record.porStage, 2);
        assertTrue(record.effectivePoR);
        assertEq(record.walletCoverageBps, 5000);
        assertEq(record.reserveCoverageBps, 9000);

        ArdmerePoRAnchor.SnapshotRecord memory latest = anchorContract.getLatestSnapshot(EXCHANGE_TAG);
        assertEq(latest.snapshotId, composite);

        ArdmerePoRAnchor.SnapshotRecord memory byPeriod =
            anchorContract.getSnapshotByPeriod(EXCHANGE_TAG, 43);
        assertEq(byPeriod.snapshotId, composite);

        assertEq(anchorContract.getSnapshotCount(EXCHANGE_TAG), 1);
        assertEq(anchorContract.getSnapshotIdAt(EXCHANGE_TAG, 0), composite);
    }

    function test_anchorSnapshot_abiEncodeDecodeRoundTrip() public {
        vm.warp(1_767_657_600);
        bytes32 composite = anchorContract.expectedSnapshotId("binance", "PR01JUN26");
        _anchorBinance("PR01JUN26", 43, 1_767_657_600, 951_913, 3, true, 10000, 7500, 10);

        bytes memory encoded = abi.encode(anchorContract.getSnapshot(composite));
        ArdmerePoRAnchor.SnapshotRecord memory decoded =
            abi.decode(encoded, (ArdmerePoRAnchor.SnapshotRecord));

        assertEq(decoded.snapshotId, composite);
        assertEq(decoded.coverageBps, 17);
        assertEq(decoded.schemaVersion, 7);
        assertEq(decoded.btcReserveRatioBps, 65535);
        assertEq(decoded.custodyRiskScore, 10);
        assertEq(decoded.porStage, 3);
        assertTrue(decoded.effectivePoR);
        assertEq(decoded.walletCoverageBps, 10000);
        assertEq(decoded.reserveCoverageBps, 7500);
        assertEq(decoded.anchoredAt, 1_767_657_600);
    }

    function test_anchorSnapshot_latestSnapshotUpdates() public {
        bytes32 first = anchorContract.expectedSnapshotId("binance", "PR01MAY26");
        bytes32 second = anchorContract.expectedSnapshotId("binance", "PR01JUN26");

        vm.warp(1_767_657_600);
        _anchorBinance("PR01MAY26", 42, 1_767_657_600, 951_000, 1, false, 0, 0, 255);

        vm.warp(1_768_032_000);
        _anchorBinance("PR01JUN26", 43, 1_768_032_000, 952_000, 2, true, 100, 200, 50);

        assertEq(anchorContract.getLatestSnapshot(EXCHANGE_TAG).snapshotId, second);
        assertEq(anchorContract.getSnapshotCount(EXCHANGE_TAG), 2);
        assertEq(anchorContract.getSnapshotIdAt(EXCHANGE_TAG, 0), first);
        assertEq(anchorContract.getSnapshotIdAt(EXCHANGE_TAG, 1), second);
    }

    function test_anchorSnapshot_revertsOnDuplicateSnapshotId() public {
        vm.warp(1_767_657_600);
        bytes32 composite = anchorContract.expectedSnapshotId("binance", "PR01JUN26");
        _anchorBinance("PR01JUN26", 43, 1_767_657_600, 951_913, 1, false, 0, 0, 0);

        vm.prank(SIGNER);
        vm.expectRevert(abi.encodeWithSelector(ArdmerePoRAnchor.SnapshotAlreadyExists.selector, composite));
        anchorContract.anchorSnapshot(
            "binance",
            "PR01JUN26",
            composite,
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
            0,
            0,
            0,
            0,
            65535,
            65535,
            65535,
            65535
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
        bytes32 composite = anchorContract.expectedSnapshotId("binance", "PR01JUN26");
        _anchorBinance("PR01JUN26", 43, 1_767_657_600, 951_913, 2, true, 1111, 2222, 88);

        ArdmerePoRAnchor newImplementation = new ArdmerePoRAnchor();
        anchorContract.upgradeToAndCall(address(newImplementation), "");

        assertEq(anchorContract.signer(), SIGNER);
        assertEq(anchorContract.SCHEMA_VERSION(), 7);
        assertTrue(anchorContract.snapshotExists(composite));
        assertEq(anchorContract.getLatestSnapshot(EXCHANGE_TAG).snapshotId, composite);
        assertEq(anchorContract.getSnapshot(composite).walletCoverageBps, 1111);
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
        assertEq(upgraded.SCHEMA_VERSION(), 7);
        assertEq(upgraded.STORAGE_VERSION(), 4);
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
        assertEq(record.custodyRiskScore, 0);
        assertEq(record.totalLiabilityUsd, 0);
        assertEq(record.totalReserveUsd, 0);
        assertEq(record.btcReserveRatioBps, 0);
        assertEq(record.usdcReserveRatioBps, 0);
    }

    function test_anchorSnapshot_storesReserveTotalsAndRatioSentinel() public {
        vm.warp(1_767_657_600);
        bytes32 composite = anchorContract.expectedSnapshotId("binance", "TOTALS");
        vm.prank(SIGNER);
        anchorContract.anchorSnapshot(
            "binance",
            "TOTALS",
            composite,
            7,
            1_767_657_600,
            1,
            EXCHANGE_ROOT,
            ARTIFACT_ROOT,
            VERIFICATION_ROOT,
            0x01,
            1,
            1,
            false,
            0,
            0,
            10,
            250_000_000_000,
            260_000_000_000,
            10200,
            65535,
            10000,
            65534
        );
        ArdmerePoRAnchor.SnapshotRecord memory record = anchorContract.getSnapshot(composite);
        assertEq(record.schemaVersion, 7);
        assertEq(record.totalLiabilityUsd, 250_000_000_000);
        assertEq(record.totalReserveUsd, 260_000_000_000);
        assertEq(record.btcReserveRatioBps, 10200);
        assertEq(record.ethReserveRatioBps, 65535);
        assertEq(record.usdtReserveRatioBps, 10000);
        assertEq(record.usdcReserveRatioBps, 65534);
    }

    function test_anchorSnapshot_revertsOnInvalidCustodyRiskScore() public {
        bytes32 composite = anchorContract.expectedSnapshotId("binance", "BAD_SCORE");
        vm.prank(SIGNER);
        vm.expectRevert(abi.encodeWithSelector(ArdmerePoRAnchor.InvalidCustodyRiskScore.selector, 101));
        anchorContract.anchorSnapshot(
            "binance",
            "BAD_SCORE",
            composite,
            99,
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
            0,
            101,
            0,
            0,
            65535,
            65535,
            65535,
            65535
        );
    }

    function test_custodyRiskUnavailable_constant() public view {
        assertEq(anchorContract.CUSTODY_RISK_UNAVAILABLE(), 255);
    }

    function test_anchorSnapshot_revertsForStranger() public {
        bytes32 composite = anchorContract.expectedSnapshotId("binance", "PR01JUN26");
        vm.prank(STRANGER);
        vm.expectRevert(abi.encodeWithSelector(ArdmerePoRAnchor.Unauthorized.selector, STRANGER));
        anchorContract.anchorSnapshot(
            "binance",
            "PR01JUN26",
            composite,
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
            0,
            255,
            0,
            0,
            65535,
            65535,
            65535,
            65535
        );
    }

    function test_anchorSnapshot_revertsOnZeroArtifactRoot() public {
        bytes32 composite = anchorContract.expectedSnapshotId("binance", "PR01JUN26");
        vm.prank(SIGNER);
        vm.expectRevert(ArdmerePoRAnchor.InvalidRoot.selector);
        anchorContract.anchorSnapshot(
            "binance",
            "PR01JUN26",
            composite,
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
            0,
            255,
            0,
            0,
            65535,
            65535,
            65535,
            65535
        );
    }

    function test_anchorSnapshot_revertsOnZeroVerificationRoot() public {
        bytes32 composite = anchorContract.expectedSnapshotId("binance", "PR01JUN26");
        vm.prank(SIGNER);
        vm.expectRevert(ArdmerePoRAnchor.InvalidRoot.selector);
        anchorContract.anchorSnapshot(
            "binance",
            "PR01JUN26",
            composite,
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
            0,
            255,
            0,
            0,
            65535,
            65535,
            65535,
            65535
        );
    }

    function test_anchorSnapshot_revertsOnEmptyExchange() public {
        bytes32 composite = anchorContract.expectedSnapshotId("binance", "PR01JUN26");
        vm.prank(SIGNER);
        vm.expectRevert(ArdmerePoRAnchor.EmptyExchange.selector);
        anchorContract.anchorSnapshot(
            "",
            "PR01JUN26",
            composite,
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
            0,
            255,
            0,
            0,
            65535,
            65535,
            65535,
            65535
        );
    }

    function test_backfillReserveFigures_keepsOldFieldsAndRejectsSecond() public {
        ArdmerePoRAnchorV6Seed seedImpl = new ArdmerePoRAnchorV6Seed();
        bytes memory initData = abi.encodeCall(ArdmerePoRAnchorV6Seed.initialize, (SIGNER));
        ERC1967Proxy proxy = new ERC1967Proxy(address(seedImpl), initData);
        ArdmerePoRAnchorV6Seed seed = ArdmerePoRAnchorV6Seed(address(proxy));

        bytes32 composite = keccak256(abi.encodePacked(bytes("binance"), bytes1(0), bytes("PR04NOV26")));
        vm.warp(1_793_750_400);
        seed.anchorSnapshot(
            "binance",
            composite,
            6,
            1_793_750_400,
            974_457,
            EXCHANGE_ROOT,
            ARTIFACT_ROOT,
            VERIFICATION_ROOT,
            0x0d,
            10000,
            1,
            false,
            0,
            0,
            91
        );

        ArdmerePoRAnchor upgradedImpl = new ArdmerePoRAnchor();
        seed.upgradeToAndCall(address(upgradedImpl), "");
        // Seed did not initialize Ownable; the backfill owner is this test contract.
        vm.store(
            address(proxy),
            0x9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300,
            bytes32(uint256(uint160(address(this))))
        );

        ArdmerePoRAnchor upgraded = ArdmerePoRAnchor(address(proxy));
        assertEq(upgraded.expectedSnapshotId("binance", "PR04NOV26"), composite);

        upgraded.backfillReserveFigures("binance", "PR04NOV26", 132_692_864_216, 133_789_440_153, 10003, 10002, 10155, 10312);

        ArdmerePoRAnchor.SnapshotRecord memory record = upgraded.getSnapshot(composite);
        assertEq(record.snapshotId, composite);
        assertEq(record.periodSeq, 6);
        assertEq(record.snapshotTime, 1_793_750_400);
        assertEq(record.btcBlockHeight, 974_457);
        assertEq(record.exchangeMerkleRoot, EXCHANGE_ROOT);
        assertEq(record.artifactBundleRoot, ARTIFACT_ROOT);
        assertEq(record.verificationBundleRoot, VERIFICATION_ROOT);
        assertEq(record.coverageBps, 10000);
        assertEq(record.custodyRiskScore, 91);
        assertEq(record.porStage, 1);
        assertFalse(record.effectivePoR);
        assertEq(record.schemaVersion, 7);
        assertEq(record.totalLiabilityUsd, 132_692_864_216);
        assertEq(record.totalReserveUsd, 133_789_440_153);
        assertEq(record.btcReserveRatioBps, 10003);
        assertEq(record.ethReserveRatioBps, 10002);
        assertEq(record.usdtReserveRatioBps, 10155);
        assertEq(record.usdcReserveRatioBps, 10312);
        assertEq(upgraded.getSnapshotCount(keccak256("binance")), 1);

        vm.expectRevert(abi.encodeWithSelector(ArdmerePoRAnchor.ReserveBackfillAlreadyApplied.selector, composite));
        upgraded.backfillReserveFigures("binance", "PR04NOV26", 1, 1, 1, 1, 1, 1);
        assertEq(upgraded.getSnapshot(composite).totalLiabilityUsd, 132_692_864_216);
        assertEq(upgraded.getSnapshot(composite).custodyRiskScore, 91);
    }

    function test_backfillReserveFigures_revertsWhenMissingAlreadyV7OrNotOwner() public {
        bytes32 missing = anchorContract.expectedSnapshotId("binance", "PR04NOV26");
        vm.expectRevert(abi.encodeWithSelector(ArdmerePoRAnchor.SnapshotNotFound.selector, missing));
        anchorContract.backfillReserveFigures("binance", "PR04NOV26", 1, 1, 1, 1, 1, 1);

        bytes32 composite = anchorContract.expectedSnapshotId("binance", "TOTALS");
        vm.prank(SIGNER);
        anchorContract.anchorSnapshot(
            "binance",
            "TOTALS",
            composite,
            7,
            1,
            1,
            EXCHANGE_ROOT,
            ARTIFACT_ROOT,
            VERIFICATION_ROOT,
            1,
            1,
            1,
            false,
            0,
            0,
            10,
            5,
            6,
            7,
            8,
            9,
            10
        );
        vm.expectRevert(abi.encodeWithSelector(ArdmerePoRAnchor.ReserveBackfillAlreadyApplied.selector, composite));
        anchorContract.backfillReserveFigures("binance", "TOTALS", 9, 9, 9, 9, 9, 9);
        assertEq(anchorContract.getSnapshot(composite).totalLiabilityUsd, 5);
        assertEq(anchorContract.getSnapshot(composite).custodyRiskScore, 10);

        vm.prank(STRANGER);
        vm.expectRevert(abi.encodeWithSignature("OwnableUnauthorizedAccount(address)", STRANGER));
        anchorContract.backfillReserveFigures("binance", "TOTALS", 1, 1, 1, 1, 1, 1);
    }

    function test_expectedSnapshotId_distinctForSharedAuditId() public view {
        bytes32 phemex = anchorContract.expectedSnapshotId("phemex", "20261101");
        bytes32 htx = anchorContract.expectedSnapshotId("htx", "20261101");
        bytes32 ourbit = anchorContract.expectedSnapshotId("ourbit", "20261101");
        assertTrue(phemex != htx);
        assertTrue(htx != ourbit);
        assertTrue(phemex != ourbit);
        assertTrue(phemex != keccak256("20261101"));
    }

    function test_legacySnapshotId_matchesKeccakBytes() public view {
        assertEq(anchorContract.legacySnapshotId("PR01JUN26"), SNAPSHOT_PR01JUN26);
    }

    function test_anchorSnapshot_revertsOnInvalidSnapshotId() public {
        bytes32 expected = anchorContract.expectedSnapshotId("binance", "PR01JUN26");
        vm.prank(SIGNER);
        vm.expectRevert(abi.encodeWithSelector(ArdmerePoRAnchor.InvalidSnapshotId.selector, expected, SNAPSHOT_PR01JUN26));
        anchorContract.anchorSnapshot(
            "binance",
            "PR01JUN26",
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
            0,
            255,
            0,
            0,
            65535,
            65535,
            65535,
            65535
        );
    }
}

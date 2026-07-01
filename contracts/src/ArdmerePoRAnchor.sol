// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

/// @title  ArdmerePoRAnchor — minimal trusted-kernel anchor contract (UUPS upgradeable)
/// @notice One on-chain anchor per exchange snapshot: both artifact and
///         verification bundle roots in a single event / transaction, with
///         queryable on-chain storage (schema v3).
///         See docs/verifier-architecture.md §6.
contract ArdmerePoRAnchor is Initializable, OwnableUpgradeable, UUPSUpgradeable {
    /// @dev Event schema version (v3 = merged single-tx anchor + on-chain storage).
    uint8 public constant SCHEMA_VERSION = 3;

    /// @dev On-chain snapshot storage layout version (independent of event schema).
    uint8 public constant STORAGE_VERSION = 1;

    address public signer;

    /// @notice Immutable on-chain record for one anchored snapshot.
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

    // keccak256(abi.encode(uint256(keccak256("ardmere.storage.ArdmerePoRAnchor")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant ArdmerePoRAnchorStorageLocation =
        0x168b5695e07d5cc77c8fba711207364d2cbe159a40632dd678b45ab21ab90100;

    /// @notice Emitted once per exchange snapshot period.
    /// @param  exchangeTag             keccak256(bytes(exchange)) — indexed for filtering
    /// @param  exchange                Plain exchange name, e.g. "binance"
    /// @param  snapshotId              keccak256(native audit id), e.g. keccak256("PR01JUN26")
    /// @param  periodSeq               1-based sequence number for this exchange
    /// @param  snapshotTime            Exchange snapshot UTC unix timestamp
    /// @param  btcBlockHeight          BTC block height time anchor (0 if N/A)
    /// @param  exchangeMerkleRoot      Exchange self-reported Merkle root (0 if unknown)
    /// @param  artifactBundleRoot      Merkle root over fetched raw artifacts
    /// @param  verificationBundleRoot  Merkle root over verifier results
    /// @param  verdictSummary          Compact verifier outcome bitfield
    /// @param  coverageBps             On-chain audit coverage × 10_000
    /// @param  schemaVersion           Event schema version (= SCHEMA_VERSION)
    /// @param  anchoredAt              block.timestamp
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

    error Unauthorized(address caller);
    error InvalidRoot();
    error EmptyExchange();
    error ZeroAddress();
    error SnapshotNotFound(bytes32 snapshotId);
    error SnapshotAlreadyExists(bytes32 snapshotId);
    error PeriodAlreadyAnchored(bytes32 exchangeTag, uint32 periodSeq);
    error SnapshotIndexOutOfBounds(bytes32 exchangeTag, uint256 index);

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function _getArdmerePoRAnchorStorage() private pure returns (ArdmerePoRAnchorStorage storage $) {
        assembly {
            $.slot := ArdmerePoRAnchorStorageLocation
        }
    }

    /// @notice One-time proxy initialization. Deployer becomes upgrade owner.
    function initialize(address _signer) external initializer {
        if (_signer == address(0)) revert ZeroAddress();
        __Ownable_init(msg.sender);
        signer = _signer;
    }

    /// @notice Rotate the permissioned anchor signer (owner only).
    function setSigner(address _signer) external onlyOwner {
        if (_signer == address(0)) revert ZeroAddress();
        address previous = signer;
        signer = _signer;
        emit SignerUpdated(previous, _signer);
    }

    /// @notice Anchor one complete PoR snapshot (artifact + verification roots).
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

        _writeSnapshot(
            exchange,
            snapshotId,
            periodSeq,
            snapshotTime,
            btcBlockHeight,
            exchangeMerkleRoot,
            artifactBundleRoot,
            verificationBundleRoot,
            verdictSummary,
            coverageBps
        );
        _emitSnapshotAnchored(snapshotId);
    }

    function _writeSnapshot(
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
    ) private {
        ArdmerePoRAnchorStorage storage $ = _getArdmerePoRAnchorStorage();
        if ($.snapshots[snapshotId].anchoredAt != 0) revert SnapshotAlreadyExists(snapshotId);

        bytes32 exchangeTag = keccak256(bytes(exchange));
        if ($.snapshotByPeriod[exchangeTag][periodSeq] != bytes32(0)) {
            revert PeriodAlreadyAnchored(exchangeTag, periodSeq);
        }

        uint256 anchoredAt = block.timestamp;
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
        record.anchoredAt = anchoredAt;

        $.snapshotIds[exchangeTag].push(snapshotId);
        $.snapshotByPeriod[exchangeTag][periodSeq] = snapshotId;
        $.latestSnapshotId[exchangeTag] = snapshotId;
        $.snapshotCount[exchangeTag] += 1;
    }

    function _emitSnapshotAnchored(bytes32 snapshotId) private {
        SnapshotRecord storage record = _getArdmerePoRAnchorStorage().snapshots[snapshotId];
        emit SnapshotAnchored(
            record.snapshotId,
            record.exchangeTag,
            record.exchange,
            record.periodSeq,
            record.snapshotTime,
            record.btcBlockHeight,
            record.exchangeMerkleRoot,
            record.artifactBundleRoot,
            record.verificationBundleRoot,
            record.verdictSummary,
            record.coverageBps,
            record.schemaVersion,
            record.anchoredAt
        );
    }

    /// @notice Whether a snapshot record exists in on-chain storage.
    function snapshotExists(bytes32 snapshotId) external view returns (bool) {
        return _getArdmerePoRAnchorStorage().snapshots[snapshotId].anchoredAt != 0;
    }

    /// @notice Fetch one snapshot by id.
    function getSnapshot(bytes32 snapshotId) external view returns (SnapshotRecord memory) {
        SnapshotRecord memory record = _getArdmerePoRAnchorStorage().snapshots[snapshotId];
        if (record.anchoredAt == 0) revert SnapshotNotFound(snapshotId);
        return record;
    }

    /// @notice Fetch the most recently anchored snapshot for an exchange tag.
    function getLatestSnapshot(bytes32 exchangeTag) external view returns (SnapshotRecord memory) {
        bytes32 snapshotId = _getArdmerePoRAnchorStorage().latestSnapshotId[exchangeTag];
        if (snapshotId == bytes32(0)) revert SnapshotNotFound(snapshotId);
        return _getArdmerePoRAnchorStorage().snapshots[snapshotId];
    }

    /// @notice Fetch a snapshot by exchange tag and period sequence.
    function getSnapshotByPeriod(bytes32 exchangeTag, uint32 periodSeq)
        external
        view
        returns (SnapshotRecord memory)
    {
        bytes32 snapshotId = _getArdmerePoRAnchorStorage().snapshotByPeriod[exchangeTag][periodSeq];
        if (snapshotId == bytes32(0)) revert SnapshotNotFound(snapshotId);
        return _getArdmerePoRAnchorStorage().snapshots[snapshotId];
    }

    /// @notice Number of snapshots anchored for an exchange tag.
    function getSnapshotCount(bytes32 exchangeTag) external view returns (uint256) {
        return _getArdmerePoRAnchorStorage().snapshotCount[exchangeTag];
    }

    /// @notice Snapshot id at chronological index (0 = earliest) for an exchange tag.
    function getSnapshotIdAt(bytes32 exchangeTag, uint256 index) external view returns (bytes32) {
        ArdmerePoRAnchorStorage storage $ = _getArdmerePoRAnchorStorage();
        bytes32[] storage ids = $.snapshotIds[exchangeTag];
        if (index >= ids.length) revert SnapshotIndexOutOfBounds(exchangeTag, index);
        return ids[index];
    }

    function _authorizeUpgrade(address) internal override onlyOwner {}
}

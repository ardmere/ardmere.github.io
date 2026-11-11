// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

/// @title  ArdmerePoRAnchor — minimal trusted-kernel anchor contract (UUPS upgradeable)
/// @notice One on-chain anchor per exchange snapshot: both artifact and
///         verification bundle roots in a single event / transaction, with
///         queryable on-chain storage (schema v7; code-ready, proxy not yet upgraded).
///         See docs/verifier-architecture.md §6.
contract ArdmerePoRAnchor is Initializable, OwnableUpgradeable, UUPSUpgradeable {
    /// @dev Event/storage schema version (v7 = self-reported USD totals + major-coin reserve ratios).
    ///      v6 = composite snapshotId enforcement. Live Sepolia proxy stays on v6 until upgraded.
    uint8 public constant SCHEMA_VERSION = 7;

    /// @dev On-chain snapshot storage layout version (independent of event schema).
    ///      v2 = SnapshotRecord appends porStage/effectivePoR/walletCoverageBps/reserveCoverageBps
    ///      after anchoredAt (v3 records remain readable; new fields decode as 0/false).
    ///      v3 = appends custodyRiskScore after reserveCoverageBps (v4 rows decode as 0).
    ///      v4 = appends totalLiabilityUsd, totalReserveUsd, and four major-coin ratio bps.
    ///      Pre-v7 rows keep schemaVersion < 7; the new slots read as 0 and are unavailable.
    uint8 public constant STORAGE_VERSION = 4;

    /// @dev custodyRiskScore sentinel when the anchor pipeline could not produce a score (uint8).
    ///      Matches NEAR sibling semantics and skills/por-custody-risk (conclusive: false offline).
    uint8 public constant CUSTODY_RISK_UNAVAILABLE = 255;

    /// @dev Whole-USD totals. 0 = the snapshot did not publish amounts, or every coin was unpriced.
    uint128 public constant USD_UNAVAILABLE = 0;

    /// @dev Major-coin reserve ratio bps. 10000 = 100%. 65535 = coin omitted or ratio unavailable.
    uint16 public constant RESERVE_RATIO_UNAVAILABLE = 65535;

    address public signer;

    /// @notice Immutable on-chain record for one anchored snapshot.
    /// @dev New v4 fields are appended AFTER `anchoredAt` so v3 ABI/storage readers
    ///      that stop at `anchoredAt` keep working; trailing slots are zero for old rows.
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
        /// @dev Deprecated in v4: opaque max(verifier.Coverage)×10000. Prefer wallet/reserve fields.
        uint16 coverageBps;
        uint8 schemaVersion;
        uint256 anchoredAt;
        /// @dev 0=Pre-Stage, 1=Stage 0, 2=Stage 1, 3=Stage 2 (matches assessment.porStage).
        uint8 porStage;
        bool effectivePoR;
        /// @dev verifiable_addrs / published_addrs × 10_000
        uint16 walletCoverageBps;
        /// @dev verifiable_usd / claimed_reserve_usd × 10_000 (CoinGecko at snapshotTime)
        uint16 reserveCoverageBps;
        /// @dev 0–100 inclusive custody risk at anchor time; 255 = unavailable / inconclusive.
        uint8 custodyRiskScore;
        /// @dev Self-reported customer liabilities in whole USD. 0 = unavailable. Schema v7.
        uint128 totalLiabilityUsd;
        /// @dev Self-reported reserves (self-custody + third-party) in whole USD. 0 = unavailable.
        uint128 totalReserveUsd;
        /// @dev BTC reserve/liability × 10_000. 65535 = unavailable. Not a solvency opinion.
        uint16 btcReserveRatioBps;
        uint16 ethReserveRatioBps;
        uint16 usdtReserveRatioBps;
        uint16 usdcReserveRatioBps;
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
    /// @param  snapshotId              v6+: keccak256(abi.encodePacked(exchange,0,auditId)); legacy rows use keccak256(bytes(auditId))
    /// @param  periodSeq               1-based sequence number for this exchange
    /// @param  snapshotTime            Exchange snapshot UTC unix timestamp
    /// @param  btcBlockHeight          BTC block height time anchor (0 if N/A)
    /// @param  exchangeMerkleRoot      Exchange self-reported Merkle root (0 if unknown)
    /// @param  artifactBundleRoot      Merkle root over fetched raw artifacts
    /// @param  verificationBundleRoot  Merkle root over verifier results
    /// @param  verdictSummary          Compact verifier outcome bitfield
    /// @param  coverageBps             Deprecated opaque coverage × 10_000 (kept for readers)
    /// @param  schemaVersion           Event schema version (= SCHEMA_VERSION)
    /// @param  anchoredAt              block.timestamp
    /// @param  porStage                0=Pre-Stage .. 3=Stage 2
    /// @param  effectivePoR            Same effective-PoR flag as the assessment
    /// @param  walletCoverageBps       Verifiable address count ratio × 10_000
    /// @param  reserveCoverageBps      Verifiable reserve USD ratio × 10_000
    /// @param  custodyRiskScore        0–100 custody risk at anchor; 255 if unavailable
    /// @param  totalLiabilityUsd       Self-reported liabilities, whole USD; 0 = unavailable (v7)
    /// @param  totalReserveUsd         Self-reported reserves, whole USD; 0 = unavailable (v7)
    /// @param  btcReserveRatioBps      BTC reserve ratio × 10_000; 65535 = unavailable
    /// @param  ethReserveRatioBps      ETH reserve ratio × 10_000; 65535 = unavailable
    /// @param  usdtReserveRatioBps     USDT reserve ratio × 10_000; 65535 = unavailable
    /// @param  usdcReserveRatioBps     USDC reserve ratio × 10_000; 65535 = unavailable
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

    error Unauthorized(address caller);
    error InvalidRoot();
    error EmptyExchange();
    error ZeroAddress();
    error SnapshotNotFound(bytes32 snapshotId);
    error SnapshotAlreadyExists(bytes32 snapshotId);
    error PeriodAlreadyAnchored(bytes32 exchangeTag, uint32 periodSeq);
    error SnapshotIndexOutOfBounds(bytes32 exchangeTag, uint256 index);
    error InvalidCustodyRiskScore(uint8 value);
    error InvalidSnapshotId(bytes32 expected, bytes32 provided);

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

    /// @notice keccak256(abi.encodePacked(normalizedExchange, bytes1(0), bytes(auditId))) — v6 storage key.
    function expectedSnapshotId(string calldata exchange, string calldata auditId) public pure returns (bytes32) {
        return keccak256(abi.encodePacked(_normalizedExchange(exchange), bytes1(0), bytes(auditId)));
    }

    /// @notice keccak256(bytes(auditId)) — v3–v5 legacy key; remains readable forever.
    function legacySnapshotId(string calldata auditId) public pure returns (bytes32) {
        return keccak256(bytes(auditId));
    }

    function _normalizedExchange(string memory exchange) private pure returns (bytes memory) {
        bytes memory out = bytes(exchange);
        for (uint256 i = 0; i < out.length; ++i) {
            uint8 c = uint8(out[i]);
            if (c >= 0x41 && c <= 0x5A) {
                out[i] = bytes1(c + 32);
            }
        }
        return out;
    }

    /// @notice Anchor one complete PoR snapshot (artifact + verification roots).
    function anchorSnapshot(
        string calldata exchange,
        string calldata auditId,
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
        uint8 custodyRiskScore,
        uint128 totalLiabilityUsd,
        uint128 totalReserveUsd,
        uint16 btcReserveRatioBps,
        uint16 ethReserveRatioBps,
        uint16 usdtReserveRatioBps,
        uint16 usdcReserveRatioBps
    ) external {
        if (msg.sender != signer) revert Unauthorized(msg.sender);
        if (bytes(exchange).length == 0) revert EmptyExchange();
        if (artifactBundleRoot == bytes32(0) || verificationBundleRoot == bytes32(0)) {
            revert InvalidRoot();
        }
        bytes32 expected = expectedSnapshotId(exchange, auditId);
        if (snapshotId != expected) revert InvalidSnapshotId(expected, snapshotId);

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
            coverageBps,
            porStage,
            effectivePoR,
            walletCoverageBps,
            reserveCoverageBps,
            custodyRiskScore,
            totalLiabilityUsd,
            totalReserveUsd,
            btcReserveRatioBps,
            ethReserveRatioBps,
            usdtReserveRatioBps,
            usdcReserveRatioBps
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
        uint16 coverageBps,
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
    ) private {
        if (custodyRiskScore > 100 && custodyRiskScore != CUSTODY_RISK_UNAVAILABLE) {
            revert InvalidCustodyRiskScore(custodyRiskScore);
        }
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
        record.porStage = porStage;
        record.effectivePoR = effectivePoR;
        record.walletCoverageBps = walletCoverageBps;
        record.reserveCoverageBps = reserveCoverageBps;
        record.custodyRiskScore = custodyRiskScore;
        record.totalLiabilityUsd = totalLiabilityUsd;
        record.totalReserveUsd = totalReserveUsd;
        record.btcReserveRatioBps = btcReserveRatioBps;
        record.ethReserveRatioBps = ethReserveRatioBps;
        record.usdtReserveRatioBps = usdtReserveRatioBps;
        record.usdcReserveRatioBps = usdcReserveRatioBps;

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
            record.anchoredAt,
            record.porStage,
            record.effectivePoR,
            record.walletCoverageBps,
            record.reserveCoverageBps,
            record.custodyRiskScore,
            record.totalLiabilityUsd,
            record.totalReserveUsd,
            record.btcReserveRatioBps,
            record.ethReserveRatioBps,
            record.usdtReserveRatioBps,
            record.usdcReserveRatioBps
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

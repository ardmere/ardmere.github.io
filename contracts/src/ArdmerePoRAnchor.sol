// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

/// @title  ArdmerePoRAnchor — minimal trusted-kernel anchor contract
/// @notice One on-chain anchor per exchange snapshot: both artifact and
///         verification bundle roots in a single event / transaction.
///         See docs/verifier-architecture.md §6.
contract ArdmerePoRAnchor {
    /// @dev Event schema version (v2 = merged single-tx anchor).
    uint8 public constant SCHEMA_VERSION = 2;

    address public immutable signer;

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

    error Unauthorized(address caller);
    error InvalidRoot();
    error EmptyExchange();

    constructor(address _signer) {
        require(_signer != address(0), "signer=0");
        signer = _signer;
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

        emit SnapshotAnchored(
            snapshotId,
            keccak256(bytes(exchange)),
            exchange,
            periodSeq,
            snapshotTime,
            btcBlockHeight,
            exchangeMerkleRoot,
            artifactBundleRoot,
            verificationBundleRoot,
            verdictSummary,
            coverageBps,
            SCHEMA_VERSION,
            block.timestamp
        );
    }
}

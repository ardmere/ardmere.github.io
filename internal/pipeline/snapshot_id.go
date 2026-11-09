package pipeline

import (
	"strings"

	"github.com/ardmere/ardmere/internal/keccak"
)

// NormalizeExchange lowercases the exchange slug (assessments store ASCII slugs).
func NormalizeExchange(exchange string) string {
	return strings.ToLower(strings.TrimSpace(exchange))
}

// LegacySnapshotID is keccak256(bytes(auditId)) — v3–v5 on-chain keys.
func LegacySnapshotID(auditID string) [32]byte {
	return keccak.Sum256([]byte(auditID))
}

// CompositeSnapshotID is keccak256(abi.encodePacked(normalizeExchange, 0x00, auditId)).
func CompositeSnapshotID(exchange, auditID string) [32]byte {
	ex := []byte(NormalizeExchange(exchange))
	packed := make([]byte, 0, len(ex)+1+len(auditID))
	packed = append(packed, ex...)
	packed = append(packed, 0)
	packed = append(packed, []byte(auditID)...)
	return keccak.Sum256(packed)
}

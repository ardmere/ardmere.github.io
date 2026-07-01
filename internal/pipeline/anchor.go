package pipeline

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/ardmere/ardmere/internal/keccak"
	"github.com/ardmere/ardmere/internal/por"
)

const anchorSelector = "0x3ab606f4"

// AnchorParams holds fields for on-chain anchorSnapshot calldata.
type AnchorParams struct {
	Exchange               string
	SnapshotID             string
	PeriodSeq              uint32
	SnapshotTimeUnix       uint64
	BTCBlockHeight         uint32
	ExchangeMerkleRoot     [32]byte
	ArtifactBundleRoot     string
	VerificationBundleRoot string
	VerdictSummary         uint8
	CoverageBps            uint16
}

// AnchorParamsFrom builds anchor params from a completed verify run.
func AnchorParamsFrom(snap por.Snapshot, exchangeRoot [32]byte, artRoot, verRoot string, verdictSummary uint8, coverageBps uint16) AnchorParams {
	return AnchorParams{
		Exchange:               snap.Exchange,
		SnapshotID:             snap.ID,
		PeriodSeq:              snap.PeriodSeq,
		SnapshotTimeUnix:       uint64(snap.SnapshotTime.Unix()),
		BTCBlockHeight:         BTCBlockHeight(snap),
		ExchangeMerkleRoot:     exchangeRoot,
		ArtifactBundleRoot:     artRoot,
		VerificationBundleRoot: verRoot,
		VerdictSummary:         verdictSummary,
		CoverageBps:            coverageBps,
	}
}

// ParseHex32 decodes a 32-byte hex merkle root.
func ParseHex32(s string) [32]byte {
	s = strings.TrimPrefix(strings.TrimSpace(s), "0x")
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != 32 {
		return [32]byte{}
	}
	var out [32]byte
	copy(out[:], b)
	return out
}

// PrintAnchorCalldata prints cast send instructions for anchorSnapshot().
func PrintAnchorCalldata(p AnchorParams) {
	snapHash := keccak.Sum256([]byte(p.SnapshotID))
	exchangeTag := keccak.Sum256([]byte(p.Exchange))
	contract := os.Getenv("ANCHOR_CONTRACT")
	rpcURL := os.Getenv("RPC_URL")
	pk := os.Getenv("PRIVATE_KEY")

	fmt.Println()
	fmt.Println("==================================================================")
	fmt.Printf("ANCHOR PAYLOAD — exchange=%s snapshot=%s period=%d (1 tx)\n",
		p.Exchange, p.SnapshotID, p.PeriodSeq)
	fmt.Println("==================================================================")
	fmt.Printf("exchange              : %s\n", p.Exchange)
	fmt.Printf("exchangeTag   (b32)   : 0x%s\n", hex.EncodeToString(exchangeTag[:]))
	fmt.Printf("snapshotId    (utf8)  : %s\n", p.SnapshotID)
	fmt.Printf("snapshotId    (b32)   : 0x%s\n", hex.EncodeToString(snapHash[:]))
	fmt.Printf("periodSeq             : %d\n", p.PeriodSeq)
	fmt.Printf("snapshotTime (unix)   : %d\n", p.SnapshotTimeUnix)
	fmt.Printf("btcBlockHeight        : %d\n", p.BTCBlockHeight)
	fmt.Printf("exchangeMerkleRoot    : 0x%s\n", hex.EncodeToString(p.ExchangeMerkleRoot[:]))
	fmt.Printf("artifactBundleRoot    : %s\n", p.ArtifactBundleRoot)
	fmt.Printf("verificationBundleRoot: %s\n", p.VerificationBundleRoot)
	fmt.Printf("verdictSummary        : 0x%02x\n", p.VerdictSummary)
	fmt.Printf("coverageBps           : %d\n", p.CoverageBps)
	fmt.Println()

	sig := "anchorSnapshot(string,bytes32,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16)"
	args := fmt.Sprintf(`"%s" 0x%s %d %d %d 0x%s %s %s 0x%02x %d`,
		p.Exchange,
		hex.EncodeToString(snapHash[:]),
		p.PeriodSeq,
		p.SnapshotTimeUnix,
		p.BTCBlockHeight,
		hex.EncodeToString(p.ExchangeMerkleRoot[:]),
		p.ArtifactBundleRoot,
		p.VerificationBundleRoot,
		p.VerdictSummary,
		p.CoverageBps,
	)

	if contract == "" || rpcURL == "" || pk == "" {
		fmt.Println("# set ANCHOR_CONTRACT, RPC_URL, PRIVATE_KEY then use cast directly:")
		fmt.Printf("cast send <ANCHOR_CONTRACT> '%s' \\\n", sig)
		fmt.Printf("  %s \\\n", args)
		fmt.Println("  --rpc-url <RPC_URL> --private-key <PRIVATE_KEY>")
	} else {
		fmt.Printf("cast send %s '%s' \\\n", contract, sig)
		fmt.Printf("  %s \\\n", args)
		fmt.Printf("  --rpc-url %s --private-key %s\n", rpcURL, pk)
	}
}

// ValidateAnchorSelector panics if keccak impl drifts from the deployed contract.
func ValidateAnchorSelector() {
	sig := "anchorSnapshot(string,bytes32,uint32,uint64,uint32,bytes32,bytes32,bytes32,uint8,uint16)"
	got := keccak.Sum256([]byte(sig))
	want, _ := hex.DecodeString(strings.TrimPrefix(anchorSelector, "0x"))
	for i := 0; i < 4; i++ {
		if got[i] != want[i] {
			panic("anchorSnapshot selector mismatch (keccak impl bug?)")
		}
	}
}

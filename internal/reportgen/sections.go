package reportgen

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ardmere/ardmere/internal/artifacts"
	"github.com/ardmere/ardmere/internal/bundle"
)

func writeSection3Frequency(b *strings.Builder, freq frequencyInfo) {
	b.WriteString("## 3. Frequency and Freshness\n\n")
	b.WriteString("| Field | Value |\n| --- | --- |\n")
	fmt.Fprintf(b, "| Latest snapshot (evaluation set) | `%s` |\n", freq.LatestSnapshot)
	fmt.Fprintf(b, "| Previous snapshot | `%s` |\n", freq.PreviousSnapshot)
	fmt.Fprintf(b, "| Observed cadence | `%s` |\n", freq.ObservedCadence)
	fmt.Fprintf(b, "| History available | `%s` |\n", freq.HistoryAvailable)
	fmt.Fprintf(b, "| Event-triggered updates | `%s` |\n", freq.EventTriggered)
	fmt.Fprintf(b, "| Daily root / commitment anchor | `%s` |\n", freq.DailyAnchor)
	fmt.Fprintf(b, "| Stage impact | %s |\n\n", freq.StageImpact)
	if freq.Notes != "" {
		b.WriteString(freq.Notes + "\n\n")
	}
}

func writeSection4Evidence(b *strings.Builder, exchangeID, snapshotID, snapshotDir string, art bundle.ArtifactBundle, ver bundle.VerificationBundle) {
	b.WriteString("## 4. Evidence and Bundles\n\n")
	b.WriteString("### Public Artifacts\n\n")
	b.WriteString("| Artifact | SHA-256 | Source |\n| --- | --- | --- |\n")
	for _, a := range art.Artifacts {
		source := a.URL
		if source == "" {
			source = a.LocalPath
		}
		fmt.Fprintf(b, "| `%s` | `%s` | %s |\n", a.Kind, a.SHA256, escapeCell(source))
	}
	b.WriteString("\n### Bundle References\n\n")
	b.WriteString("| Bundle | Value |\n| --- | --- |\n")
	fmt.Fprintf(b, "| Artifact bundle root | `%s` |\n", art.MerkleRoot)
	fmt.Fprintf(b, "| Verification bundle root | `%s` |\n", ver.MerkleRoot)

	artPath := artifacts.ResolveBundlePath(snapshotDir, snapshotID, ".artifact-bundle.json")
	verPath := artifacts.ResolveBundlePath(snapshotDir, snapshotID, ".verification-bundle.v2.json")
	if _, err := os.Stat(verPath); err != nil {
		verPath = artifacts.ResolveBundlePath(snapshotDir, snapshotID, ".verification-bundle.json")
	}

	artSHA := bundleFileSHA256(artPath)
	verSHA := bundleFileSHA256(verPath)
	fmt.Fprintf(b, "| Artifact bundle SHA-256 | `%s` |\n", artSHA)
	fmt.Fprintf(b, "| Verification bundle SHA-256 | `%s` |\n\n", verSHA)

	relBase := fmt.Sprintf("../../../artifacts/%s/%s/bundles", exchangeID, snapshotID)
	fmt.Fprintf(b, "Local bundle paths: [%s.artifact-bundle.json](%s/%s.artifact-bundle.json), [%s.verification-bundle.v2.json](%s/%s.verification-bundle.v2.json)\n\n",
		snapshotID, relBase, snapshotID,
		snapshotID, relBase, snapshotID)
}

func writeSection7Checklist(b *strings.Builder, items []checklistItem) {
	b.WriteString("## 7. Minimal Checklist\n\n")
	b.WriteString("| Requirement | Status | Notes |\n| --- | --- | --- |\n")
	for _, item := range items {
		fmt.Fprintf(b, "| %s (%s) | `%s` | %s |\n",
			escapeCell(item.Question), item.ID, item.Status, escapeCell(item.Notes))
	}
	b.WriteString("\n")
}

func writeSection8Recommendations(b *strings.Builder, recs []recommendationItem) {
	b.WriteString("## 8. Recommendations\n\n")
	b.WriteString("| Priority | Recommendation | Related Risk |\n| --- | --- | --- |\n")
	for _, rec := range recs {
		flags := strings.Join(rec.RelatedRiskFlags, ", ")
		if flags == "" {
			flags = "—"
		} else {
			flags = "`" + flags + "`"
		}
		fmt.Fprintf(b, "| `%s` | %s | %s |\n", rec.Priority, escapeCell(rec.Text), flags)
	}
	b.WriteString("\n")
}

func bundleFileSHA256(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return "UNVERIFIABLE"
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "UNVERIFIABLE"
	}
	return hex.EncodeToString(h.Sum(nil))
}

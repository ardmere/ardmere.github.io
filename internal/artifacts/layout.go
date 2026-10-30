// Package artifacts manages the on-disk layout under artifacts/<exchange>/<snapshot>/.
package artifacts

import (
	"path/filepath"
)

// ExchangeBase returns artifactsBase/<exchangeID>.
func ExchangeBase(artifactsBase, exchangeID string) string {
	return filepath.Join(artifactsBase, exchangeID)
}

// SnapshotRoot returns artifactsBase/<exchangeID>/<snapshotID>.
func SnapshotRoot(artifactsBase, exchangeID, snapshotID string) string {
	return filepath.Join(artifactsBase, exchangeID, snapshotID)
}

// SnapshotRootFromDir normalizes a path that may be the snapshot root or a
// child (raw/ or bundles/) back to the snapshot root. Idempotent.
func SnapshotRootFromDir(dir string) string {
	if dir == "" {
		return dir
	}
	base := filepath.Base(dir)
	if base == "raw" || base == "bundles" {
		return filepath.Dir(dir)
	}
	return dir
}

// ResolveBundlePath returns the path to a named bundle file under the snapshot.
// dir may be the snapshot root or the bundles/ directory.
// suffix is e.g. ".artifact-bundle.json" → <snapshotID>.artifact-bundle.json.
func ResolveBundlePath(dir, snapshotID, suffix string) string {
	name := snapshotID + suffix
	if filepath.Base(dir) == "bundles" {
		return filepath.Join(dir, name)
	}
	return filepath.Join(dir, "bundles", name)
}

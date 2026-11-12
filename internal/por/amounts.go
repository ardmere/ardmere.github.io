package por

// Extra keys for CoinSummary.Extra / Snapshot.Extra amount publication.
const (
	ExtraReserveAmount     = "reserveAmount"
	ExtraLiabilityAmount   = "liabilityAmount"
	ExtraReserveRate       = "reserveRate"
	ExtraUnpublished       = "unpublished"
	ExtraAmountsPublished  = "amountsPublished"
	ExtraAmountsFalse      = "false"
)

// MarkAmountsUnpublished records that reserve and liability quantities were not
// published for this coin (distinct from a true published zero).
func MarkAmountsUnpublished(extra map[string]string) map[string]string {
	if extra == nil {
		extra = map[string]string{}
	}
	extra[ExtraReserveAmount] = ExtraUnpublished
	extra[ExtraLiabilityAmount] = ExtraUnpublished
	return extra
}

// SnapshotAmountsNotPublished sets the snapshot-level flag used by reserve_figures
// to skip inventing USD totals from missing amounts.
func SnapshotAmountsNotPublished(extra map[string]string) map[string]string {
	if extra == nil {
		extra = map[string]string{}
	}
	extra[ExtraAmountsPublished] = ExtraAmountsFalse
	return extra
}

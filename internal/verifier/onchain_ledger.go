package verifier

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/verifier/onchainconfig"
	"github.com/ardmere/ardmere/internal/walletzip"
)

// OnchainBalanceLedger verifies non-EVM ledger rows (UTXO, XRPL, Solana).
type OnchainBalanceLedger struct {
	Exchange        string
	WalletCSV       walletzip.File
	WalletZipPath   string
	WalletZipSha256 string
	BapiSha256      string
	SnapshotID      string
	Ledger          *rpc.LedgerClient
	Concurrency     int
}

func (v OnchainBalanceLedger) Run(ctx context.Context) Verification {
	if v.Concurrency <= 0 {
		v.Concurrency = 3
	}
	if v.Ledger == nil {
		v.Ledger = rpc.NewLedger()
	}

	out := Verification{
		VerifierID:     "onchain-balance-ledger",
		Version:        "1.4",
		SnapshotID:     v.SnapshotID,
		VerifiedAt:     time.Now().UTC(),
		InputArtifacts: []string{v.BapiSha256, v.WalletZipSha256},
	}

	var allRows []walletzip.Row
	csvFile := v.WalletCSV
	if csvFile == 0 {
		csvFile = walletzip.WalletFileForExchange(v.Exchange)
	}
	_, err := walletzip.ForEachRow(v.WalletZipPath, csvFile, func(r walletzip.Row) error {
		allRows = append(allRows, r)
		return nil
	})
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("read HotCold csv: %v", err)
		return out
	}

	supported := loadLedgerSupportedFor(v.Exchange)
	var rows []walletzip.Row
	unsupportedPairs := map[string]int{}
	for _, r := range allRows {
		key := r.Coin + "|" + r.Network
		if _, ok := supported[key]; !ok {
			if shouldReportUnsupportedLedger(v.Exchange, r.Coin, r.Network) {
				unsupportedPairs[key]++
			}
			continue
		}
		rows = append(rows, r)
	}

	if len(rows) == 0 {
		out.Verdict = VerdictPartial
		out.Reason = "no HotCold rows in supported ledger (coin,network) set"
		out.Findings = unsupportedFindings(unsupportedPairs)
		return out
	}

	var sampleNote string
	rows, sampleNote = maybeSubsampleOnchainRows(v.Exchange, rows)
	if sampleNote != "" {
		out.Findings = append(out.Findings, Finding{
			Subject: "onchain-sample",
			Field:   "ledger_rows",
			Status:  VerdictPass,
			Note:    sampleNote,
		})
	}

	type result struct {
		row        walletzip.Row
		actual     decimal.Decimal
		used       string
		mode       string
		components map[string]string
		err        error
		unv        bool
		unvNote    string
	}
	jobs := make(chan walletzip.Row, len(rows))
	results := make(chan result, len(rows))
	var wg sync.WaitGroup
	for i := 0; i < v.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range jobs {
				spec := supported[r.Coin+"|"+r.Network]
				actual, used, mode, components, unv, unvNote, err := v.balanceAtHeight(ctx, r, spec)
				results <- result{row: r, actual: actual, used: used, mode: mode, components: components, err: err, unv: unv, unvNote: unvNote}
			}
		}()
	}
	for _, r := range rows {
		jobs <- r
	}
	close(jobs)
	go func() { wg.Wait(); close(results) }()

	var passCount, warnCount, failCount, unverifiableCount, errCount int
	for res := range results {
		key := res.row.Coin + "|" + res.row.Network
		spec := supported[key]
		if res.err != nil {
			errCount++
			out.Findings = append(out.Findings, Finding{
				Subject: res.row.Address,
				Field:   key,
				Status:  VerdictWarn,
				Note:    fmt.Sprintf("ledger rpc error (provider=%s): %v", res.used, res.err),
			})
			continue
		}
		if res.unv {
			unverifiableCount++
			out.Findings = append(out.Findings, Finding{
				Subject: res.row.Address,
				Field:   key + "#" + fmt.Sprintf("%d", res.row.Height),
				Claim:   res.row.Balance.String(),
				Actual:  res.actual.String(),
				Status:  VerdictUnverifiable,
				Note:    res.unvNote,
				Components: map[string]string{"mode": res.mode},
			})
			continue
		}

		if unv, note := btcNativeCustodyUnverifiable(res.row, res.actual); unv {
			unverifiableCount++
			out.Findings = append(out.Findings, Finding{
				Subject:    res.row.Address,
				Field:      key + "#" + fmt.Sprintf("%d", res.row.Height),
				Claim:      res.row.Balance.String(),
				Actual:     res.actual.String(),
				Status:     VerdictUnverifiable,
				Note:       note,
				Components: map[string]string{"mode": res.mode},
			})
			continue
		}

		claim := res.row.Balance
		diff := res.actual.Sub(claim)

		if res.row.Coin == "APT" && res.row.Network == "APT" {
			status, aptNote := aptosLedgerMismatch(res.actual, claim, res.components)
			if status == VerdictPass {
				passCount++
				continue
			}
			note := fmt.Sprintf("on-chain balance != csv claim by %s (provider=%s)", diff.Abs().String(), res.used)
			if aptNote != "" {
				note = aptNote + "; " + note
			}
			components := map[string]string{"mode": res.mode}
			for k, v := range res.components {
				components[k] = v
			}
			switch status {
			case VerdictUnverifiable:
				unverifiableCount++
			case VerdictWarn:
				warnCount++
			default:
				failCount++
			}
			out.Findings = append(out.Findings, Finding{
				Subject:    res.row.Address,
				Field:      key + "#" + fmt.Sprintf("%d", res.row.Height),
				Claim:      claim.String(),
				Actual:     res.actual.String(),
				Status:     status,
				Note:       note,
				Components: components,
			})
			continue
		}

		within, surplus := classifyBalanceMismatch(res.actual, claim)
		if within {
			passCount++
			continue
		}
		if surplus {
			passCount++
			continue
		}

		note := fmt.Sprintf("on-chain balance != csv claim by %s (provider=%s)", diff.Abs().String(), res.used)
		if snap := ledgerLiveSnapshotNote(spec, res.components); snap != "" {
			note = snap + "; " + note
		}
		if tronSnap := tronSnapshotMismatchNote(res.row.Coin, res.row.Network); tronSnap != "" {
			note = tronSnap + "; " + note
		}
		status := VerdictFail
		if ledgerLiveSnapshotNote(spec, res.components) != "" || tronSnapshotMismatchNote(res.row.Coin, res.row.Network) != "" {
			status = VerdictWarn
		}
		if status == VerdictWarn {
			warnCount++
		} else {
			failCount++
		}
		components := map[string]string{"mode": res.mode}
		for k, v := range res.components {
			components[k] = v
		}
		out.Findings = append(out.Findings, Finding{
			Subject:    res.row.Address,
			Field:      key + "#" + fmt.Sprintf("%d", res.row.Height),
			Claim:      claim.String(),
			Actual:     res.actual.String(),
			Status:     status,
			Note:       note,
			Components: components,
		})
	}

	out.Coverage = float64(passCount+failCount) / float64(len(allRows))
	out.Findings = append(out.Findings, unsupportedFindings(unsupportedPairs)...)
	out.Findings = append(out.Findings, Finding{
		Subject: "summary",
		Field:   "verified_pairs",
		Actual: fmt.Sprintf("%d pass / %d warn / %d fail / %d unverifiable / %d rpc-error -- of %d supported ledger rows / %d total HotCold rows",
			passCount, warnCount, failCount, unverifiableCount, errCount, len(rows), len(allRows)),
		Status: VerdictPass,
	})

	switch {
	case failCount > 0:
		out.Verdict = VerdictFail
	case unverifiableCount > 0:
		out.Verdict = VerdictWarn
		out.Reason = "some ledger rows cannot be verified without archive API (e.g. high-volume DOGE)"
	case warnCount > 0:
		out.Verdict = VerdictWarn
		out.Reason = "some rows use live Solana balance or other snapshot API limits vs CSV"
	case errCount > len(rows)/4:
		out.Verdict = VerdictWarn
		out.Reason = "more than 25% of supported ledger rows failed to query"
	case len(rows) < len(allRows):
		out.Verdict = VerdictPartial
	default:
		out.Verdict = VerdictPass
	}
	return out
}

func (v OnchainBalanceLedger) balanceAtHeight(ctx context.Context, r walletzip.Row, spec onchainconfig.LedgerSpec) (
	actual decimal.Decimal, used, mode string, components map[string]string, unverifiable bool, unvNote string, err error,
) {
	components = map[string]string{}
	switch spec.Kind {
	case onchainconfig.LedgerEsplora:
		raw, u, e := v.Ledger.EsploraBalanceAtHeight(ctx, spec.EsploraBase, r.Address, r.Height)
		if e != nil && spec.Alchemy != "" {
			if bal, u2, mode, comps, e2 := v.balanceAlchemy(ctx, r, spec); e2 == nil {
				return bal, u2, mode, comps, false, "", nil
			}
		}
		if e != nil && spec.Blockcypher != "" {
			return v.balanceBlockcypher(ctx, r, spec)
		}
		if e != nil && spec.Blockchair != "" {
			raw, u, live, e2 := v.Ledger.BlockchairBalanceAtHeight(ctx, spec.Blockchair, r.Address, r.Height)
			if e2 == nil {
				mode := "blockchair_utxo"
				if live {
					mode = "blockchair_live"
					components["live_snapshot"] = "true"
				}
				return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, mode, components, false, "", nil
			}
		}
		if e != nil {
			return decimal.Zero, u, "esplora", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "esplora_utxo", nil, false, "", nil

	case onchainconfig.LedgerAlchemy:
		bal, u, mode, comps, err := v.balanceAlchemy(ctx, r, spec)
		if err == nil {
			return bal, u, mode, comps, false, "", nil
		}
		if spec.Blockchair != "" {
			raw, u, live, e := v.Ledger.BlockchairBalanceAtHeight(ctx, spec.Blockchair, r.Address, r.Height)
			if e == nil {
				mode := "blockchair_utxo"
				if live {
					mode = "blockchair_live"
					components["live_snapshot"] = "true"
				}
				return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, mode, components, false, "", nil
			}
		}
		return decimal.Zero, u, "alchemy", nil, false, "", err

	case onchainconfig.LedgerBlockchair:
		raw, u, live, e := v.Ledger.BlockchairBalanceAtHeight(ctx, spec.Blockchair, r.Address, r.Height)
		if e != nil && spec.Blockchair == "zcash" {
			raw, u, e = v.Ledger.CipherScanZECBalance(ctx, r.Address)
			if e == nil {
				live = true
				components["live_snapshot"] = "true"
			}
		}
		if e != nil && spec.Alchemy != "" {
			if bal, u2, mode, comps, e2 := v.balanceAlchemy(ctx, r, spec); e2 == nil {
				return bal, u2, mode, comps, false, "", nil
			}
		}
		if e != nil && spec.Blockcypher != "" {
			return v.balanceBlockcypher(ctx, r, spec)
		}
		if e != nil {
			return decimal.Zero, u, "blockchair", nil, false, "", e
		}
		mode := "blockchair_utxo"
		if live {
			mode = "blockchair_live"
			components["live_snapshot"] = "true"
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, mode, components, false, "", nil

	case onchainconfig.LedgerBlockcypher:
		return v.balanceBlockcypher(ctx, r, spec)

	case onchainconfig.LedgerXRPL:
		raw, u, e := v.Ledger.XRPLBalanceAtLedger(ctx, r.Address, r.Height)
		if e != nil {
			return decimal.Zero, u, "xrpl", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "xrpl_account", nil, false, "", nil

	case onchainconfig.LedgerSolanaNative:
		slot := r.Height
		raw, u, historical, e := v.Ledger.SolanaNativeBalanceAtSlot(ctx, r.Address, slot)
		if e != nil {
			return decimal.Zero, u, "solana", components, false, "", e
		}
		mode := "solana_native"
		if historical {
			mode = "solana_native_slot"
			components["historical_slot"] = strconv.FormatInt(slot, 10)
			components["history_provider"] = rpc.SolanaHistoryProviderFromUsed(u)
		}
		return decimal.NewFromInt(int64(raw)).Shift(-int32(spec.Decimals)), u, mode, components, false, "", nil

	case onchainconfig.LedgerSolanaSPL:
		slot := r.Height
		raw, u, historical, e := v.Ledger.SolanaSPLBalanceAtSlot(ctx, r.Address, spec.Mint, slot)
		if e != nil {
			return decimal.Zero, u, "solana-spl", components, false, "", e
		}
		mode := "solana_spl"
		if historical {
			mode = "solana_spl_slot"
			components["historical_slot"] = strconv.FormatInt(slot, 10)
			components["history_provider"] = rpc.SolanaHistoryProviderFromUsed(u)
		}
		return decimal.NewFromInt(int64(raw)).Shift(-int32(spec.Decimals)), u, mode, components, false, "", nil

	case onchainconfig.LedgerAptos:
		raw, u, comps, e := v.Ledger.AptosAccountedAPTAtVersion(ctx, r.Address, r.Height)
		if e != nil {
			return decimal.Zero, u, "aptos", comps, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "aptos_accounted", comps, false, "", nil

	case onchainconfig.LedgerAptosFA:
		raw, u, e := v.Ledger.AptosFABalanceAtVersion(ctx, r.Address, spec.Mint, r.Height)
		if e != nil {
			return decimal.Zero, u, "aptos-fa", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "aptos_fa", nil, false, "", nil

	case onchainconfig.LedgerNear:
		raw, u, live, e := v.Ledger.NearNativeBalanceAtHeight(ctx, r.Address, r.Height)
		if e != nil {
			return decimal.Zero, u, "near", nil, false, "", e
		}
		mode := "near_native"
		if live {
			mode = "near_native_live"
		}
		return decimal.NewFromBigInt(raw, -int32(spec.Decimals)), u, mode, nil, false, "", nil

	case onchainconfig.LedgerNearFT:
		raw, u, live, e := v.Ledger.NearFTBalanceAtHeight(ctx, spec.Mint, r.Address, r.Height)
		if e != nil {
			return decimal.Zero, u, "near-ft", nil, false, "", e
		}
		mode := "near_ft"
		if live {
			mode = "near_ft_live"
		}
		return decimal.NewFromBigInt(raw, -int32(spec.Decimals)), u, mode, nil, false, "", nil

	case onchainconfig.LedgerSui:
		raw, u, e := v.Ledger.SuiBalanceAtCheckpoint(ctx, r.Address, "", r.Height)
		if e != nil {
			return decimal.Zero, u, "sui", nil, false, "", e
		}
		return decimal.NewFromInt(int64(raw)).Shift(-int32(spec.Decimals)), u, "sui_native", nil, false, "", nil

	case onchainconfig.LedgerSuiCoin:
		raw, u, e := v.Ledger.SuiBalanceAtCheckpoint(ctx, r.Address, spec.Mint, r.Height)
		if e != nil {
			return decimal.Zero, u, "sui-coin", nil, false, "", e
		}
		return decimal.NewFromInt(int64(raw)).Shift(-int32(spec.Decimals)), u, "sui_coin", nil, false, "", nil

	case onchainconfig.LedgerHbar:
		raw, u, e := v.Ledger.HbarBalanceAtBlock(ctx, r.Address, r.Height)
		if e != nil {
			return decimal.Zero, u, "hbar", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "hbar_native", nil, false, "", nil

	case onchainconfig.LedgerHbarHTS:
		raw, u, e := v.Ledger.HbarHTSBalanceAtBlock(ctx, r.Address, spec.Mint, r.Height)
		if e != nil {
			return decimal.Zero, u, "hbar-hts", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "hbar_hts", nil, false, "", nil

	case onchainconfig.LedgerSubstrate:
		raw, u, e := v.Ledger.SubstrateBalanceAtBlock(ctx, spec.RPCURL, r.Address, r.Height)
		if e != nil {
			return decimal.Zero, u, "substrate", nil, false, "", e
		}
		return decimal.NewFromBigInt(raw, -int32(spec.Decimals)), u, "substrate_free", nil, false, "", nil

	case onchainconfig.LedgerTonJetton:
		raw, u, e := v.Ledger.TonJettonBalance(ctx, r.Address, spec.Mint, r.Height)
		if e != nil {
			return decimal.Zero, u, "ton-jetton", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "ton_jetton_live", nil, false, "", nil

	case onchainconfig.LedgerAlgoASA:
		assetID := int64(spec.AssetID)
		if assetID == 0 && spec.Mint != "" {
			assetID, _ = strconv.ParseInt(spec.Mint, 10, 64)
		}
		raw, u, e := v.Ledger.AlgoASABalanceAtRound(ctx, r.Address, assetID, r.Height)
		if e != nil {
			return decimal.Zero, u, "algo-asa", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "algo_asa", nil, false, "", nil

	case onchainconfig.LedgerAptosCoin:
		raw, u, e := v.Ledger.AptosCoinBalanceByTypeAtVersion(ctx, r.Address, spec.Mint, r.Height)
		if e != nil {
			return decimal.Zero, u, "aptos-coin", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "aptos_coin", nil, false, "", nil

	case onchainconfig.LedgerSubstrateAsset:
		raw, u, e := v.Ledger.SubstrateAssetBalanceAtBlock(ctx, spec.RPCURL, r.Address, int64(spec.AssetID), r.Height)
		if e != nil {
			return decimal.Zero, u, "substrate-asset", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "substrate_asset", nil, false, "", nil

	case onchainconfig.LedgerStellarAsset:
		raw, u, e := v.Ledger.StellarAssetBalanceAtLedger(ctx, r.Address, spec.Mint, spec.Issuer, r.Height)
		if e != nil {
			return decimal.Zero, u, "stellar-asset", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "stellar_asset", nil, false, "", nil

	case onchainconfig.LedgerTezosFA:
		raw, u, e := v.Ledger.TezosFABalanceAtLevel(ctx, r.Address, spec.Mint, r.Height)
		if e != nil {
			return decimal.Zero, u, "tezos-fa", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "tezos_fa", nil, false, "", nil

	case onchainconfig.LedgerXRPLEmitted:
		raw, u, e := v.Ledger.XRPLEmittedBalanceAtLedger(ctx, r.Address, spec.Issuer, spec.Mint, r.Height, spec.Decimals)
		if e != nil {
			return decimal.Zero, u, "xrpl-issued", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "xrpl_issued", nil, false, "", nil

	case onchainconfig.LedgerChromia:
		raw, u, e := v.Ledger.ChromiaCHRBalanceAtHeight(ctx, spec.RPCURL, spec.Mint, r.Address, r.Height)
		if e != nil {
			return decimal.Zero, u, "chromia", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "chromia_ft4", nil, false, "", nil

	case onchainconfig.LedgerFilecoin:
		raw, u, e := v.Ledger.FilecoinActorBalance(ctx, r.Address, r.Height)
		if e != nil {
			return decimal.Zero, u, "filecoin", nil, false, "", e
		}
		return decimal.NewFromBigInt(raw, -int32(spec.Decimals)), u, "filecoin_actor", nil, false, "", nil

	case onchainconfig.LedgerTonNative:
		raw, u, e := v.Ledger.TonNativeBalance(ctx, r.Address, r.Height)
		if e != nil {
			return decimal.Zero, u, "ton-native", nil, false, "", e
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "ton_native", nil, false, "", nil

	default:
		return decimal.Zero, "", string(spec.Kind), nil, false, "", fmt.Errorf("unsupported ledger kind %q", spec.Kind)
	}
}

func (v OnchainBalanceLedger) balanceAlchemy(ctx context.Context, r walletzip.Row, spec onchainconfig.LedgerSpec) (
	decimal.Decimal, string, string, map[string]string, error,
) {
	if spec.Alchemy == "" {
		return decimal.Zero, "", "alchemy", nil, fmt.Errorf("alchemy chain not configured")
	}
	raw, u, err := v.Ledger.AlchemyBalanceAtHeight(ctx, spec.Alchemy, r.Address, r.Height)
	if err != nil {
		return decimal.Zero, u, "alchemy", nil, err
	}
	return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "alchemy_utxo", nil, nil
}

func (v OnchainBalanceLedger) balanceBlockcypher(ctx context.Context, r walletzip.Row, spec onchainconfig.LedgerSpec) (
	decimal.Decimal, string, string, map[string]string, bool, string, error,
) {
	raw, nTx, u, e := v.Ledger.BlockcypherBalanceBefore(ctx, spec.Blockcypher, r.Address, r.Height)
	if e != nil {
		return decimal.Zero, u, "blockcypher", nil, false, "", e
	}
	if spec.MaxTxCount > 0 && nTx > spec.MaxTxCount {
		if spec.Blockchair != "" {
			raw, u, live, e2 := v.Ledger.BlockchairBalanceAtHeight(ctx, spec.Blockchair, r.Address, r.Height)
			if e2 == nil {
				components := map[string]string{"blockcypher_tx_count": fmt.Sprintf("%d", nTx)}
				mode := "blockchair_utxo"
				if live {
					mode = "blockchair_live"
					components["live_snapshot"] = "true"
				}
				return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, mode, components, false, "", nil
			}
		}
		return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "blockcypher", nil,
			true, fmt.Sprintf("address has %d txs; exceeds public API limit %d — set BLOCKCHAIR_API_KEY or use archive provider", nTx, spec.MaxTxCount),
			nil
	}
	return decimal.NewFromInt(raw).Shift(-int32(spec.Decimals)), u, "blockcypher_utxo", nil, false, "", nil
}

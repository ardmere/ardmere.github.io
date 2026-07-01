package verifier

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/keccak"
	"github.com/ardmere/ardmere/internal/rpc"
	"github.com/ardmere/ardmere/internal/walletzip"
)

const (
	bscStakeHubAddress    = "0x0000000000000000000000000000000000002002"
	ethDepositContract    = "0x00000000219ab540356cbb839cbe05303d7705fa"
	evmNativeDecimals     = 18
	maxStakeHubValidators = 500
)

type accountedBalance struct {
	total      decimal.Decimal
	used       string
	components map[string]string
	note       string
	incomplete bool
}

type bscStakeHubCache struct {
	mu              sync.Mutex
	creditByHeight  map[int64][]string
	stakedByAddress map[string]stakeHubAccounted
}

type ethDepositCache struct {
	mu        sync.Mutex
	byAddress map[string]ethDepositCached
}

type ethDepositCached struct {
	accounted ETHDepositAccounted
	err       error
}

type stakeHubAccounted struct {
	staked   decimal.Decimal
	locked   decimal.Decimal
	used     string
	scanned  int
	scanNote string
}

func newBSCStakeHubCache() *bscStakeHubCache {
	return &bscStakeHubCache{
		creditByHeight:  map[int64][]string{},
		stakedByAddress: map[string]stakeHubAccounted{},
	}
}

func newETHDepositCache() *ethDepositCache {
	return &ethDepositCache{byAddress: map[string]ethDepositCached{}}
}

func (v OnchainBalanceHot) accountedBalance(ctx context.Context, r walletzip.Row, stakeCache *bscStakeHubCache, ethCache *ethDepositCache) (accountedBalance, error) {
	net := loadNativeSupported(v.Exchange)[r.Coin+"|"+r.Network]
	bi, used, err := v.RPC.GetBalance(ctx, net, r.Address, r.Height)
	if err != nil {
		return accountedBalance{used: used}, err
	}
	liquid := decimal.NewFromBigInt(bi, -evmNativeDecimals)
	out := accountedBalance{
		total: liquid,
		used:  used,
		components: map[string]string{
			"liquid":    liquid.String(),
			"staked":    "0",
			"unbonding": "0",
		},
	}

	switch r.Coin + "|" + r.Network {
	case "BNB|BSC":
		staked, err := v.bscStakeHubAccounted(ctx, r.Address, r.Height, stakeCache)
		if err != nil {
			out.incomplete = true
			out.note = fmt.Sprintf("stakehub lookup unavailable; liquid-only observation (provider=%s): %v", used, err)
			return out, nil
		}
		out.total = out.total.Add(staked.staked).Add(staked.locked)
		out.components["staked"] = staked.staked.String()
		out.components["unbonding"] = staked.locked.String()
		out.used = joinProviders(used, staked.used)
		out.note = fmt.Sprintf("liquid + StakeHub pooled/unbonding; scanned %d credit contracts", staked.scanned)
		if staked.scanNote != "" {
			out.note += "; " + staked.scanNote
		}
	case "ETH|ETH":
		gap := r.Balance.Sub(liquid)
		if likelyEthDepositGap(r.Balance, liquid, gap) {
			deposits, err := v.ethDepositAccounted(ctx, r.Address, r.Height, ethCache)
			if err != nil {
				out.incomplete = true
				out.components["unsupported"] = gap.String()
				out.note = fmt.Sprintf("likely ETH2 deposit balance; deposit indexer unavailable: %v", err)
				break
			}
			out.total = out.total.Add(deposits.Deposited)
			out.components["staked"] = deposits.Deposited.String()
			out.note = fmt.Sprintf("liquid + Eth2 deposits; %d deposit txs via %s; %s", deposits.TxCount, deposits.Source, deposits.Note)
		}
	}

	return out, nil
}

func (v OnchainBalanceHot) tryAccountedHeightWindow(ctx context.Context, r walletzip.Row, stakeCache *bscStakeHubCache, ethCache *ethDepositCache) (accountedBalance, int64, string, bool) {
	if r.Coin+"|"+r.Network != "ETH|ETH" {
		return accountedBalance{}, 0, "", false
	}
	for _, delta := range []int64{-1, 1} {
		h := r.Height + delta
		if h < 0 {
			continue
		}
		row := r
		row.Height = h
		ab, err := v.accountedBalance(ctx, row, stakeCache, ethCache)
		if err != nil {
			continue
		}
		if balanceWithinTolerance(ab.total, r.Balance) {
			note := fmt.Sprintf("height-boundary: csv block %d matches at block %d", r.Height, h)
			return ab, h, note, true
		}
	}
	return accountedBalance{}, 0, "", false
}

func (v OnchainBalanceHot) ethDepositAccounted(ctx context.Context, address string, height int64, cache *ethDepositCache) (ETHDepositAccounted, error) {
	key := strings.ToLower(address) + fmt.Sprintf("@%d", height)
	cache.mu.Lock()
	if cached, ok := cache.byAddress[key]; ok {
		cache.mu.Unlock()
		return cached.accounted, cached.err
	}
	cache.mu.Unlock()

	accounted, err := v.ETHDeposits.DepositedETH(ctx, address, height)
	cache.mu.Lock()
	cache.byAddress[key] = ethDepositCached{accounted: accounted, err: err}
	cache.mu.Unlock()
	return accounted, err
}

func (v OnchainBalanceHot) bscStakeHubAccounted(ctx context.Context, delegator string, height int64, cache *bscStakeHubCache) (stakeHubAccounted, error) {
	key := strings.ToLower(delegator) + fmt.Sprintf("@%d", height)
	cache.mu.Lock()
	if cached, ok := cache.stakedByAddress[key]; ok {
		cache.mu.Unlock()
		return cached, nil
	}
	cache.mu.Unlock()

	credits, used, err := v.stakeHubCreditContracts(ctx, height, cache)
	if err != nil {
		return stakeHubAccounted{}, err
	}
	pooledByCredit, lockedByCredit, batchUsed, err := stakeHubBatchBalances(ctx, v.RPC, delegator, credits, height)
	if err != nil {
		return stakeHubAccounted{}, err
	}
	var totalPooled, totalLocked decimal.Decimal
	usedProviders := joinProviders(used, batchUsed)
	lockedUnavailable := 0
	for i, credit := range credits {
		pooledWei, ok := pooledByCredit[credit]
		if !ok {
			return stakeHubAccounted{}, fmt.Errorf("missing getPooledBNB result for %s", credit)
		}
		lockedWei, ok := lockedByCredit[credit]
		if !ok {
			lockedUnavailable++
			lockedWei = new(big.Int)
		}
		totalPooled = totalPooled.Add(decimal.NewFromBigInt(pooledWei, -evmNativeDecimals))
		totalLocked = totalLocked.Add(decimal.NewFromBigInt(lockedWei, -evmNativeDecimals))
		_ = i
	}

	out := stakeHubAccounted{
		staked:  totalPooled,
		locked:  totalLocked,
		used:    usedProviders,
		scanned: len(credits),
	}
	if lockedUnavailable > 0 {
		out.scanNote = fmt.Sprintf("lockedBNBs unavailable on %d/%d credit contracts; treated as zero", lockedUnavailable, len(credits))
	}
	cache.mu.Lock()
	cache.stakedByAddress[key] = out
	cache.mu.Unlock()
	return out, nil
}

func (v OnchainBalanceHot) stakeHubCreditContracts(ctx context.Context, height int64, cache *bscStakeHubCache) ([]string, string, error) {
	cache.mu.Lock()
	if credits, ok := cache.creditByHeight[height]; ok {
		cache.mu.Unlock()
		return credits, "cache", nil
	}
	cache.mu.Unlock()

	data := append(selector("getValidators(uint256,uint256)"), encodeUint64(0)...)
	data = append(data, encodeUint64(maxStakeHubValidators)...)
	raw, used, err := v.RPC.CallContract(ctx, rpc.NetBSC, bscStakeHubAddress, data, height)
	if err != nil {
		return nil, used, fmt.Errorf("StakeHub getValidators: %w", err)
	}
	credits, err := decodeSecondAddressArray(raw)
	if err != nil {
		return nil, used, fmt.Errorf("decode getValidators credit contracts: %w", err)
	}
	if len(credits) == 0 {
		return nil, used, fmt.Errorf("StakeHub getValidators returned no credit contracts")
	}

	cache.mu.Lock()
	cache.creditByHeight[height] = credits
	cache.mu.Unlock()
	return credits, used, nil
}

func stakeHubBatchBalances(ctx context.Context, c *rpc.Client, delegator string, credits []string, height int64) (map[string]*big.Int, map[string]*big.Int, string, error) {
	pooledSig := selector("getPooledBNB(address)")
	lockedSig := selector("lockedBNBs(address,uint256)")
	delegatorArg := encodeAddress(delegator)
	zeroArg := encodeUint64(0)

	var calls []rpc.ContractCall
	pooledIDs := make(map[int]string, len(credits))
	lockedIDs := make(map[int]string, len(credits))
	nextID := 1
	for _, credit := range credits {
		pooledData := append(append([]byte(nil), pooledSig...), delegatorArg...)
		pid := nextID
		nextID++
		pooledIDs[pid] = credit
		calls = append(calls, rpc.ContractCall{ID: pid, To: credit, Data: pooledData})

		lockedData := append(append(append([]byte(nil), lockedSig...), delegatorArg...), zeroArg...)
		lid := nextID
		nextID++
		lockedIDs[lid] = credit
		calls = append(calls, rpc.ContractCall{ID: lid, To: credit, Data: lockedData})
	}

	rawByID, used, err := c.CallContractBatch(ctx, rpc.NetBSC, calls, height)
	if err != nil {
		return nil, nil, used, err
	}
	pooled := make(map[string]*big.Int, len(credits))
	locked := make(map[string]*big.Int, len(credits))
	for id, credit := range pooledIDs {
		raw, ok := rawByID[id]
		if !ok || len(raw) < 32 {
			return nil, nil, used, fmt.Errorf("short getPooledBNB response for %s", credit)
		}
		pooled[credit] = new(big.Int).SetBytes(raw[:32])
	}
	for id, credit := range lockedIDs {
		raw, ok := rawByID[id]
		if !ok {
			continue
		}
		if len(raw) < 32 {
			continue
		}
		locked[credit] = new(big.Int).SetBytes(raw[:32])
	}
	return pooled, locked, used, nil
}

func callUint256(ctx context.Context, c *rpc.Client, net rpc.Network, to, sig string, args [][]byte, height int64) (*big.Int, string, error) {
	data := selector(sig)
	for _, arg := range args {
		data = append(data, arg...)
	}
	raw, used, err := c.CallContract(ctx, net, to, data, height)
	if err != nil {
		return nil, used, err
	}
	if len(raw) < 32 {
		return nil, used, fmt.Errorf("short uint256 response: %d bytes", len(raw))
	}
	return new(big.Int).SetBytes(raw[:32]), used, nil
}

func selector(sig string) []byte {
	h := keccak.Sum256([]byte(sig))
	return append([]byte(nil), h[:4]...)
}

func encodeUint64(v uint64) []byte {
	return rpc.To32Bytes(new(big.Int).SetUint64(v).Bytes())
}

func encodeAddress(addr string) []byte {
	b, err := rpc.EncodeHexAddress(addr)
	if err != nil {
		return make([]byte, 32)
	}
	return rpc.To32Bytes(b)
}

func decodeSecondAddressArray(raw []byte) ([]string, error) {
	if len(raw) < 64 {
		return nil, fmt.Errorf("short response: %d bytes", len(raw))
	}
	offset := int(new(big.Int).SetBytes(raw[32:64]).Int64())
	return decodeAddressArrayAt(raw, offset)
}

func decodeAddressArrayAt(raw []byte, offset int) ([]string, error) {
	if offset < 0 || offset+32 > len(raw) {
		return nil, fmt.Errorf("offset %d outside response length %d", offset, len(raw))
	}
	n := int(new(big.Int).SetBytes(raw[offset : offset+32]).Int64())
	start := offset + 32
	if start+n*32 > len(raw) {
		return nil, fmt.Errorf("address array length %d exceeds response length %d", n, len(raw))
	}
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		word := raw[start+i*32 : start+(i+1)*32]
		out = append(out, "0x"+hex.EncodeToString(word[12:32]))
	}
	return out, nil
}

func likelyEthDepositGap(claim, liquid, gap decimal.Decimal) bool {
	if gap.LessThan(decimal.NewFromInt(32)) {
		return false
	}
	// Near-empty EOA with any material gap — try Eth2 deposit indexer (OKX 212/914 ETH rows).
	if liquid.LessThan(decimal.NewFromInt(1)) {
		return true
	}
	// Large claim with tiny liquid share (e.g. 1.7k / 312k ETH omnibus label).
	if claim.GreaterThan(decimal.NewFromInt(100)) {
		if liquid.Div(claim).LessThan(decimal.NewFromFloat(0.02)) {
			return true
		}
	}
	// Eth2 deposits are gwei-denominated and usually aggregate to exact 32 ETH
	// multiples. Keep this heuristic for mid-size gaps with some liquid dust.
	remainder := gap.Mod(decimal.NewFromInt(32)).Abs()
	closeToMultiple := remainder.LessThanOrEqual(decimal.NewFromFloat(1e-4)) ||
		decimal.NewFromInt(32).Sub(remainder).LessThanOrEqual(decimal.NewFromFloat(1e-4))
	return closeToMultiple || claim.GreaterThan(decimal.NewFromInt(1000))
}

// ethHotInternalCustodyLikely marks ETH|ETH rows where liquid (+ any indexed deposits)
// is far below CSV — omnibus / internal ledger, not a reserve shortfall signal.
func ethHotInternalCustodyLikely(claim, accounted decimal.Decimal) (bool, string) {
	if claim.LessThanOrEqual(decimal.Zero) || accounted.GreaterThan(claim) {
		return false, ""
	}
	if claim.LessThan(decimal.NewFromInt(100)) {
		return false, ""
	}
	ratio := accounted.Div(claim)
	if ratio.GreaterThanOrEqual(decimal.NewFromFloat(0.02)) {
		return false, ""
	}
	return true, "chain liquid+indexed deposits << CSV ETH allocation; likely omnibus/internal custody label — not single-address reserve"
}

func joinProviders(providers ...string) string {
	seen := map[string]bool{}
	var out []string
	for _, p := range providers {
		for _, part := range strings.Split(p, ",") {
			part = strings.TrimSpace(part)
			if part == "" || seen[part] {
				continue
			}
			seen[part] = true
			out = append(out, part)
		}
	}
	return strings.Join(out, ",")
}

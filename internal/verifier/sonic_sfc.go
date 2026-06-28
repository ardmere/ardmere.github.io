package verifier

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/ardmere/ardmere/internal/rpc"
)

const (
	sonicSFCAddress    = "0xFC00FACE00000000000000000000000000000000"
	maxSonicValidators = 64
)

type sonicSFCCache struct {
	mu              sync.Mutex
	lastIDByHeight  map[int64]uint64
	stakedByAddress map[string]sonicSFCAccounted
}

type sonicSFCAccounted struct {
	liquid    decimal.Decimal
	staked    decimal.Decimal
	scanned   int
	used      string
	byValidator map[string]string
}

func newSonicSFCCache() *sonicSFCCache {
	return &sonicSFCCache{
		lastIDByHeight:  map[int64]uint64{},
		stakedByAddress: map[string]sonicSFCAccounted{},
	}
}

func (v OnchainBalanceToken) sonicAccountedBalance(ctx context.Context, address string, height int64, cache *sonicSFCCache) (sonicSFCAccounted, error) {
	key := strings.ToLower(address) + fmt.Sprintf("@%d", height)
	cache.mu.Lock()
	if cached, ok := cache.stakedByAddress[key]; ok {
		cache.mu.Unlock()
		return cached, nil
	}
	cache.mu.Unlock()

	liquid, used, err := nativeBalanceDecimal(ctx, v.RPC, rpc.NetSonic, address, height, 18)
	if err != nil {
		return sonicSFCAccounted{}, err
	}

	lastID, used2, err := v.sonicLastValidatorID(ctx, height, cache)
	if err != nil {
		return sonicSFCAccounted{liquid: liquid, used: used}, fmt.Errorf("sonic lastValidatorID: %w", err)
	}
	used = joinProviders(used, used2)

	stakedWei, byValidator, batchUsed, err := sonicStakeBatch(ctx, v.RPC, address, lastID, height)
	if err != nil {
		return sonicSFCAccounted{liquid: liquid, used: used}, err
	}
	used = joinProviders(used, batchUsed)

	out := sonicSFCAccounted{
		liquid:      liquid,
		staked:      decimal.NewFromBigInt(stakedWei, -18),
		scanned:     int(lastID),
		used:        used,
		byValidator: byValidator,
	}
	cache.mu.Lock()
	cache.stakedByAddress[key] = out
	cache.mu.Unlock()
	return out, nil
}

func (v OnchainBalanceToken) sonicLastValidatorID(ctx context.Context, height int64, cache *sonicSFCCache) (uint64, string, error) {
	cache.mu.Lock()
	if id, ok := cache.lastIDByHeight[height]; ok {
		cache.mu.Unlock()
		return id, "", nil
	}
	cache.mu.Unlock()

	raw, used, err := callUint256(ctx, v.RPC, rpc.NetSonic, sonicSFCAddress, "lastValidatorID()", nil, height)
	if err != nil {
		return 0, used, err
	}
	id := raw.Uint64()
	if id == 0 || id > maxSonicValidators {
		id = maxSonicValidators
	}
	cache.mu.Lock()
	cache.lastIDByHeight[height] = id
	cache.mu.Unlock()
	return id, used, nil
}

func sonicStakeBatch(ctx context.Context, c *rpc.Client, delegator string, lastID uint64, height int64) (*big.Int, map[string]string, string, error) {
	if lastID == 0 {
		return new(big.Int), nil, "", nil
	}
	getStakeSig := selector("getStake(address,uint256)")
	delegatorArg := encodeAddress(delegator)
	var calls []rpc.ContractCall
	for vid := uint64(1); vid <= lastID; vid++ {
		data := append(append([]byte(nil), getStakeSig...), delegatorArg...)
		data = append(data, encodeUint64(vid)...)
		calls = append(calls, rpc.ContractCall{ID: int(vid), To: sonicSFCAddress, Data: data})
	}
	rawByID, used, err := c.CallContractBatch(ctx, rpc.NetSonic, calls, height)
	if err != nil {
		return nil, nil, used, err
	}
	total := new(big.Int)
	byValidator := map[string]string{}
	for vid := uint64(1); vid <= lastID; vid++ {
		raw, ok := rawByID[int(vid)]
		if !ok || len(raw) < 32 {
			continue
		}
		amt := new(big.Int).SetBytes(raw[:32])
		if amt.Sign() == 0 {
			continue
		}
		total.Add(total, amt)
		byValidator[fmt.Sprintf("validator_%d", vid)] = decimal.NewFromBigInt(amt, -18).String()
	}
	return total, byValidator, used, nil
}

// SonicSFCSummary is the liquid + SFC staked breakdown for one delegator.
type SonicSFCSummary struct {
	Liquid       decimal.Decimal
	Staked       decimal.Decimal
	Scanned      int
	Used         string
	ByValidator  map[string]string
}

// SonicSFCAccounted exposes Sonic SFC delegation aggregation for diagnostics.
func SonicSFCAccounted(ctx context.Context, rpcClient *rpc.Client, delegator string, height int64) (SonicSFCSummary, error) {
	v := OnchainBalanceToken{RPC: rpcClient, Concurrency: 1}
	out, err := v.sonicAccountedBalance(ctx, delegator, height, newSonicSFCCache())
	if err != nil {
		return SonicSFCSummary{}, err
	}
	return SonicSFCSummary{
		Liquid:      out.liquid,
		Staked:      out.staked,
		Scanned:     out.scanned,
		Used:        out.used,
		ByValidator: out.byValidator,
	}, nil
}

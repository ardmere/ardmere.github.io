package verifier

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	okxcommon "github.com/okx/proof-of-reserves/common"

	"github.com/ardmere/ardmere/internal/exchanges/okx/okxwallet"
)

// AddressOwnershipOKX verifies "I am an OKX address" signatures in the public wallet CSV.
type AddressOwnershipOKX struct {
	WalletZipPath   string
	WalletZipSha256 string
	SnapshotID      string
	// AptosRPC optional base URL (default: public mainnet fullnode).
	AptosRPC string
}

func (v AddressOwnershipOKX) Run() Verification {
	out := Verification{
		VerifierID:     "address-ownership",
		Version:        "okx-1.1",
		SnapshotID:     v.SnapshotID,
		VerifiedAt:     time.Now().UTC(),
		InputArtifacts: []string{v.WalletZipSha256},
	}

	state := &okxOwnershipState{
		aptosRPC:  strings.TrimRight(v.AptosRPC, "/"),
		authCache: map[string]string{},
	}
	if state.aptosRPC == "" {
		state.aptosRPC = defaultAptosFullnode
	}

	var checked, passed, failed int64
	var aptosResourceOK int64
	_, err := okxwallet.ForEachAddressRow(v.WalletZipPath, func(r okxwallet.AddressRow) error {
		if okxcommon.IsVerifyAddressBannedCoin(strings.ToUpper(r.Coin)) {
			return nil
		}
		checked++
		via, err := verifyOKXAddressRow(state, r)
		if err != nil {
			failed++
			if failed <= 5 {
				out.Findings = append(out.Findings, Finding{
					Subject: r.Address,
					Field:   "signature",
					Claim:   r.Message,
					Status:  VerdictFail,
					Note:    err.Error(),
				})
			}
			return nil
		}
		passed++
		if via == ownershipViaAptosAuthKey {
			aptosResourceOK++
		}
		return nil
	})
	if err != nil {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("read okx wallet csv: %v", err)
		return out
	}
	if checked == 0 {
		out.Verdict = VerdictUnverifiable
		out.Reason = "no verifiable address rows in wallet zip"
		return out
	}
	out.Coverage = float64(passed) / float64(checked)
	if failed == 0 {
		out.Verdict = VerdictPass
		out.Reason = fmt.Sprintf("verified %d/%d address signatures", passed, checked)
		if aptosResourceOK > 0 {
			out.Reason += fmt.Sprintf(" (%d Aptos resource accounts via authentication_key)", aptosResourceOK)
		}
	} else {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("%d/%d address signatures failed", failed, checked)
	}
	return out
}

type ownershipPath int

const (
	ownershipViaStandard ownershipPath = iota
	ownershipViaAptosAuthKey
)

type okxOwnershipState struct {
	aptosRPC  string
	authCache map[string]string
	http      *http.Client
}

func verifyOKXAddressRow(state *okxOwnershipState, r okxwallet.AddressRow) (ownershipPath, error) {
	coin := strings.ToUpper(r.Coin)
	if okxcommon.IsVerifyAddressBannedCoin(coin) {
		return ownershipViaStandard, nil
	}
	if r.Address == "" || r.Message == "" || r.Signature1 == "" {
		return ownershipViaStandard, fmt.Errorf("missing address/message/signature")
	}
	coinType, exists := okxcommon.PorCoinTypeMap[coin]
	if !exists {
		coinType = okxcommon.EcdsaCoinType
	}
	switch coinType {
	case okxcommon.EvmCoinTye:
		// Align with okx/proof-of-reserves common.verifyCSVLineInternal (EVM):
		// verify owner1+sign1 when present; only verify owner2 when both owner2 and sign2 are set.
		if okxFieldSet(r.EOA1) {
			if err := safeVerify(func() error {
				return okxcommon.VerifyEvmCoin(coin, r.EOA1, r.Message, r.Signature1)
			}); err != nil {
				return ownershipViaStandard, err
			}
			if okxFieldSet(r.EOA2) && okxFieldSet(r.Signature2) {
				return ownershipViaStandard, safeVerify(func() error {
					return okxcommon.VerifyEvmCoin(coin, r.EOA2, r.Message, r.Signature2)
				})
			}
			return ownershipViaStandard, nil
		}
		return ownershipViaStandard, safeVerify(func() error {
			return okxcommon.VerifyEvmCoin(coin, r.Address, r.Message, r.Signature1)
		})
	case okxcommon.EcdsaCoinType:
		if okxFieldSet(r.EOA1) {
			if err := safeVerify(func() error {
				return okxcommon.VerifyEcdsaCoin(coin, r.EOA1, r.Message, r.Signature1)
			}); err != nil {
				return ownershipViaStandard, err
			}
			if okxFieldSet(r.EOA2) && okxFieldSet(r.Signature2) {
				return ownershipViaStandard, safeVerify(func() error {
					return okxcommon.VerifyEcdsaCoin(coin, r.EOA2, r.Message, r.Signature2)
				})
			}
			return ownershipViaStandard, nil
		}
		return ownershipViaStandard, safeVerify(func() error {
			return okxcommon.VerifyEcdsaCoin(coin, r.Address, r.Message, r.Signature1)
		})
	case okxcommon.Ed25519CoinType:
		return verifyOKXEd25519Row(state, coin, r)
	case okxcommon.TrxCoinType:
		return ownershipViaStandard, safeVerify(func() error {
			return okxcommon.VerifyTRX(r.Address, r.Message, r.Signature1)
		})
	case okxcommon.BethCoinType:
		return ownershipViaStandard, safeVerify(func() error {
			return okxcommon.VerifyBETH(r.Address, r.Message, r.Signature1)
		})
	case okxcommon.UTXOCoinType:
		return ownershipViaStandard, safeVerify(func() error {
			return okxcommon.VerifyUtxoCoin(coin, r.Address, r.Message, r.Signature1, r.Signature2, r.Script)
		})
	case okxcommon.StarkCoinType:
		return ownershipViaStandard, safeVerify(func() error {
			return okxcommon.VerifyStarkCoin(coin, r.Address, r.Message, r.Signature1, r.Script)
		})
	default:
		return ownershipViaStandard, fmt.Errorf("unsupported coin type %s", coinType)
	}
}

func verifyOKXEd25519Row(state *okxOwnershipState, coin string, r okxwallet.AddressRow) (ownershipPath, error) {
	err := safeVerify(func() error {
		return okxcommon.VerifyEd25519Coin(coin, r.Address, r.Message, r.Signature1, r.Script)
	})
	if err == nil {
		return ownershipViaStandard, nil
	}
	if !strings.Contains(err.Error(), "recovery address not match") {
		return ownershipViaStandard, err
	}
	addrType, ok := okxcommon.PorCoinAddressTypeMap[coin]
	if !ok || addrType != "APTOS" {
		return ownershipViaStandard, err
	}
	recovers, derr := aptosRecoverAddrsFromPubKey(r.Script)
	if derr != nil {
		return ownershipViaStandard, fmt.Errorf("%v; derive recover: %w", err, derr)
	}
	auth, aerr := state.aptosAuthenticationKey(r.Address)
	if aerr != nil {
		return ownershipViaStandard, fmt.Errorf("%v; aptos auth_key: %w", err, aerr)
	}
	authN := normalizeAptosHexAddr(auth)
	for _, rec := range recovers {
		if normalizeAptosHexAddr(rec) == authN {
			return ownershipViaAptosAuthKey, nil
		}
	}
	return ownershipViaStandard, err
}

const defaultAptosFullnode = "https://fullnode.mainnet.aptoslabs.com/v1"

func (s *okxOwnershipState) aptosAuthenticationKey(addr string) (string, error) {
	key := normalizeAptosHexAddr(addr)
	if s.authCache != nil {
		if v, ok := s.authCache[key]; ok {
			return v, nil
		}
	}
	fetch := fetchAptosAuthKey
	if fetch == nil {
		fetch = defaultFetchAptosAuthKey
	}
	base := s.aptosRPC
	if base == "" {
		base = defaultAptosFullnode
	}
	auth, err := fetch(base, key, s.http)
	if err != nil {
		return "", err
	}
	if s.authCache != nil {
		s.authCache[key] = auth
	}
	return auth, nil
}

// fetchAptosAuthKey is overridable in tests.
var fetchAptosAuthKey = defaultFetchAptosAuthKey

func defaultFetchAptosAuthKey(base, addr string, client *http.Client) (string, error) {
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	u := strings.TrimRight(base, "/") + "/accounts/" + url.PathEscape(addr)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "ardmere-por/okx-address-ownership")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var out struct {
		AuthenticationKey string `json:"authentication_key"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", err
	}
	if out.AuthenticationKey == "" {
		return "", fmt.Errorf("empty authentication_key")
	}
	return out.AuthenticationKey, nil
}

func aptosRecoverAddrsFromPubKey(pubkeyHex string) ([]string, error) {
	pubkeyBytes, err := okxcommon.Decode(pubkeyHex)
	if err != nil {
		return nil, err
	}
	if len(pubkeyBytes) == 0 {
		return nil, fmt.Errorf("empty public key")
	}
	publicKey := append(append([]byte{}, pubkeyBytes...), 0x0)
	rAddr := "0x" + hex.EncodeToString(okxcommon.Sha256Hash(publicKey))
	re := regexp.MustCompile("^0x0*")
	short := re.ReplaceAllString(rAddr, "0x")
	return []string{short, rAddr}, nil
}

func normalizeAptosHexAddr(addr string) string {
	addr = strings.TrimSpace(strings.ToLower(addr))
	if addr == "" {
		return addr
	}
	if !strings.HasPrefix(addr, "0x") {
		addr = "0x" + addr
	}
	h := strings.TrimPrefix(addr, "0x")
	if len(h) < 64 {
		h = strings.Repeat("0", 64-len(h)) + h
	}
	return "0x" + h
}

// okxFieldSet mirrors OKX CSV null sentinels in verifyCSVLineInternal.
func okxFieldSet(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	switch strings.ToLower(s) {
	case "null", `\n`, `\N`:
		return false
	}
	return true
}

// safeVerify catches panics from upstream OKX helpers (e.g. MustDecode on empty hex).
func safeVerify(fn func() error) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("verify panic: %v", rec)
		}
	}()
	return fn()
}

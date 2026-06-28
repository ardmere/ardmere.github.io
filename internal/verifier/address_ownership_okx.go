package verifier

import (
	"fmt"
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
}

func (v AddressOwnershipOKX) Run() Verification {
	out := Verification{
		VerifierID:     "address-ownership",
		Version:        "okx-1",
		SnapshotID:     v.SnapshotID,
		VerifiedAt:     time.Now().UTC(),
		InputArtifacts: []string{v.WalletZipSha256},
	}

	var checked, passed, failed int64
	_, err := okxwallet.ForEachAddressRow(v.WalletZipPath, func(r okxwallet.AddressRow) error {
		if okxcommon.IsVerifyAddressBannedCoin(strings.ToUpper(r.Coin)) {
			return nil
		}
		checked++
		if err := verifyOKXAddressRow(r); err != nil {
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
	} else {
		out.Verdict = VerdictFail
		out.Reason = fmt.Sprintf("%d/%d address signatures failed", failed, checked)
	}
	return out
}

func verifyOKXAddressRow(r okxwallet.AddressRow) error {
	coin := strings.ToUpper(r.Coin)
	if okxcommon.IsVerifyAddressBannedCoin(coin) {
		return nil
	}
	if r.Address == "" || r.Message == "" || r.Signature1 == "" {
		return fmt.Errorf("missing address/message/signature")
	}
	coinType, exists := okxcommon.PorCoinTypeMap[coin]
	if !exists {
		coinType = okxcommon.EcdsaCoinType
	}
	switch coinType {
	case okxcommon.EvmCoinTye:
		if r.EOA1 != "" && r.EOA2 != "" {
			if err := okxcommon.VerifyEvmCoin(coin, r.EOA1, r.Message, r.Signature1); err != nil {
				return err
			}
			return okxcommon.VerifyEvmCoin(coin, r.EOA2, r.Message, r.Signature2)
		}
		if r.EOA1 != "" {
			return okxcommon.VerifyEvmCoin(coin, r.EOA1, r.Message, r.Signature1)
		}
		return okxcommon.VerifyEvmCoin(coin, r.Address, r.Message, r.Signature1)
	case okxcommon.EcdsaCoinType:
		if r.EOA1 != "" && r.EOA2 != "" {
			if err := okxcommon.VerifyEcdsaCoin(coin, r.EOA1, r.Message, r.Signature1); err != nil {
				return err
			}
			return okxcommon.VerifyEcdsaCoin(coin, r.EOA2, r.Message, r.Signature2)
		}
		if r.EOA1 != "" {
			return okxcommon.VerifyEcdsaCoin(coin, r.EOA1, r.Message, r.Signature1)
		}
		return okxcommon.VerifyEcdsaCoin(coin, r.Address, r.Message, r.Signature1)
	case okxcommon.Ed25519CoinType:
		return okxcommon.VerifyEd25519Coin(coin, r.Address, r.Message, r.Signature1, r.Script)
	case okxcommon.TrxCoinType:
		return okxcommon.VerifyTRX(r.Address, r.Message, r.Signature1)
	case okxcommon.BethCoinType:
		return okxcommon.VerifyBETH(r.Address, r.Message, r.Signature1)
	case okxcommon.UTXOCoinType:
		return okxcommon.VerifyUtxoCoin(coin, r.Address, r.Message, r.Signature1, r.Signature2, r.Script)
	case okxcommon.StarkCoinType:
		return okxcommon.VerifyStarkCoin(coin, r.Address, r.Message, r.Signature1, r.Script)
	default:
		return fmt.Errorf("unsupported coin type %s", coinType)
	}
}

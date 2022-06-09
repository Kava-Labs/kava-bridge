package session

import (
	tss_common "github.com/binance-chain/tss-lib/common"
	"github.com/binance-chain/tss-lib/tss"
)

// SigningSessionResult is the result of a signing session.
type SigningSessionResult struct {
	Signature *tss_common.SignatureData
	Err       *tss.Error
}

func (res *SigningSessionResult) HasSignature() bool {
	return res.Signature != nil
}

// NewSigningSessionResult returns a new SigningSessionResult.
func NewSigningSessionResult(
	signature *tss_common.SignatureData,
	err *tss.Error,
) SigningSessionResult {
	return SigningSessionResult{
		Signature: signature,
		Err:       err,
	}
}

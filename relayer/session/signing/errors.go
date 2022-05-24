package signing

import "errors"

var (
	ErrInvalidSessionType = errors.New("invalid session type")
	ErrMismatchedTxHash   = errors.New("mismatched tx hash in signing session")
)

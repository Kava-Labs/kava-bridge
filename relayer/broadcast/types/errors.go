package types

import "errors"

var (
	ErrMsgIDEmpty                = errors.New("message ID is empty")
	ErrMsgInsufficientRecipients = errors.New("not enough recipient peer IDs, requires at least 1")
	ErrMsgTTLTooShort            = errors.New("message TTL is too short, must be at least 1 second")
	ErrMsgExpired                = errors.New("message is expired")
)

package session_test

import (
	"testing"

	"github.com/binance-chain/tss-lib/common"
	"github.com/kava-labs/kava-bridge/relayer/session"
	"github.com/stretchr/testify/require"
)

func TestHasSignature(t *testing.T) {
	result := session.NewSigningSessionResult(nil, nil)
	require.False(t, result.HasSignature())

	result = session.NewSigningSessionResult(&common.SignatureData{}, nil)
	require.True(t, result.HasSignature())
}

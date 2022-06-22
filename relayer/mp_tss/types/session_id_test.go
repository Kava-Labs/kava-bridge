package types_test

import (
	"testing"

	"github.com/kava-labs/kava-bridge/relayer/mp_tss/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsPeerParticipant(t *testing.T) {
	sessionID := types.AggregateSigningSessionID(
		AppendSlices(
			types.SigningSessionIDPart{1}.Bytes(),
			types.SigningSessionIDPart{2}.Bytes(),
			types.SigningSessionIDPart{3}.Bytes(),
			types.SigningSessionIDPart{4}.Bytes(),
			types.SigningSessionIDPart{5}.Bytes(),
		),
	)

	require.True(t, sessionID.Validate())

	for i := 1; i <= 10; i++ {
		isParticipant := sessionID.IsPeerParticipant(types.SigningSessionIDPart{byte(i)})

		if i <= 5 {
			assert.True(t, isParticipant)
		} else {
			assert.False(t, isParticipant)
		}
	}
}

// AppendSlices concatenates all the given slices.
func AppendSlices(slices ...[]byte) []byte {
	var result []byte
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

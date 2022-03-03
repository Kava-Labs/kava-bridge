package relayer_test

import (
	"testing"
	"time"

	"github.com/kava-labs/kava-bridge/relayer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoordinator_KavaNotInitialized(t *testing.T) {
	// Due to sequencing retries, the coordinator must never process an ethereum block without first being
	// initialized with a kava block.
	//
	// All ethereum blocks must be rejected with an error and return a nil slice for the signer
	// to execute (this is defensive in case the consumer ignores the returned error).
	//
	// If we allowed ethereum blocks to be added before initializing the coordinator with kava blocks, then we could
	// end up non-deterministically sequencing new transfers with retries from failed kava transactions.
	coordinator := relayer.NewCoordinator()

	kavaActions, err := coordinator.AddEthBlock(relayer.NewEthBlock(1, time.Now()))

	assert.Empty(t, kavaActions, "expected kava actions to have length 0")
	assert.Equal(t, err, relayer.ErrKavaNotInitialized, "expected an error to be returned")
}

func TestCoordinator_KavaEnforcedBlockOrder(t *testing.T) {
	blockTime := time.Now()
	coordinator := relayer.NewCoordinator()

	kavaBlock := relayer.NewKavaBlock(1, blockTime)

	err := coordinator.AddKavaBlock(kavaBlock)
	require.NoError(t, err, "expected no error when adding first block")

	kavaBlock = relayer.NewKavaBlock(1, blockTime.Add(10*time.Second))
	err = coordinator.AddKavaBlock(kavaBlock)
	require.Equal(t, err, relayer.ErrInvalidBlockHeight, "expected same height to error")

	kavaBlock = relayer.NewKavaBlock(1, blockTime.Add(10*time.Second))
	err = coordinator.AddKavaBlock(kavaBlock)
	require.Equal(t, err, relayer.ErrInvalidBlockHeight, "expected previous height to error")

	kavaBlock = relayer.NewKavaBlock(0, blockTime.Add(10*time.Second))
	err = coordinator.AddKavaBlock(kavaBlock)
	require.Equal(t, err, relayer.ErrInvalidBlockHeight, "expected earlier height to error")

	kavaBlock = relayer.NewKavaBlock(3, blockTime.Add(10*time.Second))
	err = coordinator.AddKavaBlock(kavaBlock)
	require.Equal(t, err, relayer.ErrInvalidBlockHeight, "expected height gap to error")

	kavaBlock = relayer.NewKavaBlock(2, blockTime)
	err = coordinator.AddKavaBlock(kavaBlock)
	require.Equal(t, err, relayer.ErrInvalidBlockTime, "expected same timestamp to error")

	kavaBlock = relayer.NewKavaBlock(2, blockTime.Add(-1*time.Second))
	err = coordinator.AddKavaBlock(kavaBlock)
	require.Equal(t, err, relayer.ErrInvalidBlockTime, "expected previous timestamp to error")
}

// TODO: this case is only relevant when an ethereum block has *not* been added yet
// it is valid to add a ethereum block after the kava block IF AND ONLY IF an ethereum block with one less
// height is added -- meaning
func TestCoordinator_EthTimePastKava_NoPreviousEthBlock(t *testing.T) {
	// Similar to the not initialized case, or being unable to process an ethereum block until at least one
	// kava block has been added, the coordinator must never process an ethereum block that has a timestamp
	// equal or later than the last added kava block.
	//
	// All ethereum blocks greater than the kava chain must be rejected with an error and return a nil slice for the signer
	// to execute (this is defensive in case the consumer ignores the returned error).
	//
	// If we allowed ethereum blocks to be added before initializing the coordinator with kava blocks, then we could
	// end up non-deterministically sequencing new transfers with retries from failed kava transactions.
	kavaBlockTime := time.Now()
	coordinator := relayer.NewCoordinator()
	err := coordinator.AddKavaBlock(relayer.NewKavaBlock(1, kavaBlockTime))
	require.NoError(t, err)

	// greater than kava block time
	kavaActions, err := coordinator.AddEthBlock(relayer.NewEthBlock(1, kavaBlockTime.Add(1*time.Second)))

	assert.Empty(t, kavaActions, "expected kava actions to have length 0")
	assert.Equal(t, err, relayer.ErrEthBlockAhead, "expected an error to be returned")

	// equal to kava block time
	kavaActions, err = coordinator.AddEthBlock(relayer.NewEthBlock(1, kavaBlockTime))

	assert.Empty(t, kavaActions, "expected kava actions to have length 0")
	assert.Equal(t, err, relayer.ErrEthBlockAhead, "expected an error to be returned")
}

func TestCoordinator_EthBlockBeforeKava(t *testing.T) {
	// In order to prioritize throughput, we can process ethereum blocks before a kava block with no error
	kavaBlockTime := time.Now()
	coordinator := relayer.NewCoordinator()
	err := coordinator.AddKavaBlock(relayer.NewKavaBlock(1, kavaBlockTime))
	require.NoError(t, err)

	// 100 seconds in the past
	kavaActions, err := coordinator.AddEthBlock(relayer.NewEthBlock(1, kavaBlockTime.Add(-100*time.Second)))
	assert.Empty(t, kavaActions, "expected kava actions to have length 0")
	assert.NoError(t, err, "expected no error")

	// 10 seconds in the past
	kavaActions, err = coordinator.AddEthBlock(relayer.NewEthBlock(2, kavaBlockTime.Add(-10*time.Second)))
	assert.Empty(t, kavaActions, "expected kava actions to have length 0")
	assert.NoError(t, err, "expected no error")

	// 1 second in the past
	kavaActions, err = coordinator.AddEthBlock(relayer.NewEthBlock(3, kavaBlockTime.Add(-1*time.Second)))
	assert.Empty(t, kavaActions, "expected kava actions to have length 0")
	assert.NoError(t, err, "expected no error")
}

func TestCoordinator_EthEnforcedBlockOrder(t *testing.T) {
	blockTime := time.Now()
	coordinator := relayer.NewCoordinator()

	err := coordinator.AddKavaBlock(relayer.NewKavaBlock(1, blockTime.Add(10*time.Minute)))
	ethBlock := relayer.NewEthBlock(1, blockTime)

	_, err = coordinator.AddEthBlock(ethBlock)
	require.NoError(t, err, "expected no error when adding first block")

	ethBlock = relayer.NewEthBlock(1, blockTime.Add(10*time.Second))
	_, err = coordinator.AddEthBlock(ethBlock)
	require.Equal(t, err, relayer.ErrInvalidBlockHeight, "expected same height to error")

	ethBlock = relayer.NewEthBlock(1, blockTime.Add(10*time.Second))
	_, err = coordinator.AddEthBlock(ethBlock)
	require.Equal(t, err, relayer.ErrInvalidBlockHeight, "expected previous height to error")

	ethBlock = relayer.NewEthBlock(0, blockTime.Add(10*time.Second))
	_, err = coordinator.AddEthBlock(ethBlock)
	require.Equal(t, err, relayer.ErrInvalidBlockHeight, "expected earlier height to error")

	ethBlock = relayer.NewEthBlock(3, blockTime.Add(10*time.Second))
	_, err = coordinator.AddEthBlock(ethBlock)
	require.Equal(t, err, relayer.ErrInvalidBlockHeight, "expected height gap to error")

	ethBlock = relayer.NewEthBlock(2, blockTime)
	_, err = coordinator.AddEthBlock(ethBlock)
	require.Equal(t, err, relayer.ErrInvalidBlockTime, "expected same timestamp to error")

	ethBlock = relayer.NewEthBlock(2, blockTime.Add(-1*time.Second))
	_, err = coordinator.AddEthBlock(ethBlock)
	require.Equal(t, err, relayer.ErrInvalidBlockTime, "expected previous timestamp to error")
}

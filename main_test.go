package main_test

import (
	"testing"
	"time"

	bridge "github.com/kava-labs/kava-bridge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRelayer_KavaNotInitialized(t *testing.T) {
	// Due to sequencing retries, the relayer must never process an ethereum block without first being
	// initialized with a kava block.
	//
	// All ethereum blocks must be rejected with an error and return a nil slice for the signer
	// to execute (this is defensive in case the consumer ignores the returned error).
	//
	// If we allowed ethereum blocks to be added before initializing the relayer with kava blocks, then we could
	// end up non-deterministically sequencing new transfers with retries from failed kava transactions.
	relayer := bridge.NewRelayer()

	kavaActions, err := relayer.AddEthBlock(bridge.NewEthBlock(1, time.Now()))

	assert.Empty(t, kavaActions, "expected kava actions to have length 0")
	assert.Equal(t, err, bridge.ErrKavaNotInitialized, "expected an error to be returned")
}

func TestRelayer_KavaEnforcedBlockOrder(t *testing.T) {
	blockTime := time.Now()
	relayer := bridge.NewRelayer()

	kavaBlock := bridge.NewKavaBlock(1, blockTime)

	err := relayer.AddKavaBlock(kavaBlock)
	require.NoError(t, err, "expected no error when adding first block")

	kavaBlock = bridge.NewKavaBlock(1, blockTime.Add(10*time.Second))
	err = relayer.AddKavaBlock(kavaBlock)
	require.Equal(t, err, bridge.ErrInvalidBlockHeight, "expected same height to error")

	kavaBlock = bridge.NewKavaBlock(1, blockTime.Add(10*time.Second))
	err = relayer.AddKavaBlock(kavaBlock)
	require.Equal(t, err, bridge.ErrInvalidBlockHeight, "expected previous height to error")

	kavaBlock = bridge.NewKavaBlock(0, blockTime.Add(10*time.Second))
	err = relayer.AddKavaBlock(kavaBlock)
	require.Equal(t, err, bridge.ErrInvalidBlockHeight, "expected earlier height to error")

	kavaBlock = bridge.NewKavaBlock(3, blockTime.Add(10*time.Second))
	err = relayer.AddKavaBlock(kavaBlock)
	require.Equal(t, err, bridge.ErrInvalidBlockHeight, "expected height gap to error")

	kavaBlock = bridge.NewKavaBlock(2, blockTime)
	err = relayer.AddKavaBlock(kavaBlock)
	require.Equal(t, err, bridge.ErrInvalidBlockTime, "expected same timestamp to error")

	kavaBlock = bridge.NewKavaBlock(2, blockTime.Add(-1*time.Second))
	err = relayer.AddKavaBlock(kavaBlock)
	require.Equal(t, err, bridge.ErrInvalidBlockTime, "expected previous timestamp to error")
}

// TODO: this case is only relevant when an ethereum block has *not* been added yet
// it is valid to add a ethereum block after the kava block IF AND ONLY IF an ethereum block with one less
// height is added -- meaning
func TestRelayer_EthTimePastKava_NoPreviousEthBlock(t *testing.T) {
	// Similar to the not initialized case, or being unable to process an ethereum block until at least one
	// kava block has been added, the relayer must never process an ethereum block that has a timestamp
	// equal or later than the last added kava block.
	//
	// All ethereum blocks greater than the kava chain must be rejected with an error and return a nil slice for the signer
	// to execute (this is defensive in case the consumer ignores the returned error).
	//
	// If we allowed ethereum blocks to be added before initializing the relayer with kava blocks, then we could
	// end up non-deterministically sequencing new transfers with retries from failed kava transactions.
	kavaBlockTime := time.Now()
	relayer := bridge.NewRelayer()
	err := relayer.AddKavaBlock(bridge.NewKavaBlock(1, kavaBlockTime))
	require.NoError(t, err)

	// greater than kava block time
	kavaActions, err := relayer.AddEthBlock(bridge.NewEthBlock(1, kavaBlockTime.Add(1*time.Second)))

	assert.Empty(t, kavaActions, "expected kava actions to have length 0")
	assert.Equal(t, err, bridge.ErrEthBlockAhead, "expected an error to be returned")

	// equal to kava block time
	kavaActions, err = relayer.AddEthBlock(bridge.NewEthBlock(1, kavaBlockTime))

	assert.Empty(t, kavaActions, "expected kava actions to have length 0")
	assert.Equal(t, err, bridge.ErrEthBlockAhead, "expected an error to be returned")
}

func TestRelayer_EthBlockBeforeKava(t *testing.T) {
	// In order to prioritize throughput, we can process ethereum blocks before a kava block with no error
	kavaBlockTime := time.Now()
	relayer := bridge.NewRelayer()
	err := relayer.AddKavaBlock(bridge.NewKavaBlock(1, kavaBlockTime))
	require.NoError(t, err)

	// 100 seconds in the past
	kavaActions, err := relayer.AddEthBlock(bridge.NewEthBlock(1, kavaBlockTime.Add(-100*time.Second)))
	assert.Empty(t, kavaActions, "expected kava actions to have length 0")
	assert.NoError(t, err, "expected no error")

	// 10 seconds in the past
	kavaActions, err = relayer.AddEthBlock(bridge.NewEthBlock(2, kavaBlockTime.Add(-10*time.Second)))
	assert.Empty(t, kavaActions, "expected kava actions to have length 0")
	assert.NoError(t, err, "expected no error")

	// 1 second in the past
	kavaActions, err = relayer.AddEthBlock(bridge.NewEthBlock(3, kavaBlockTime.Add(-1*time.Second)))
	assert.Empty(t, kavaActions, "expected kava actions to have length 0")
	assert.NoError(t, err, "expected no error")
}

func TestRelayer_EthEnforcedBlockOrder(t *testing.T) {
	blockTime := time.Now()
	relayer := bridge.NewRelayer()

	err := relayer.AddKavaBlock(bridge.NewKavaBlock(1, blockTime.Add(10*time.Minute)))
	ethBlock := bridge.NewEthBlock(1, blockTime)

	_, err = relayer.AddEthBlock(ethBlock)
	require.NoError(t, err, "expected no error when adding first block")

	ethBlock = bridge.NewEthBlock(1, blockTime.Add(10*time.Second))
	_, err = relayer.AddEthBlock(ethBlock)
	require.Equal(t, err, bridge.ErrInvalidBlockHeight, "expected same height to error")

	ethBlock = bridge.NewEthBlock(1, blockTime.Add(10*time.Second))
	_, err = relayer.AddEthBlock(ethBlock)
	require.Equal(t, err, bridge.ErrInvalidBlockHeight, "expected previous height to error")

	ethBlock = bridge.NewEthBlock(0, blockTime.Add(10*time.Second))
	_, err = relayer.AddEthBlock(ethBlock)
	require.Equal(t, err, bridge.ErrInvalidBlockHeight, "expected earlier height to error")

	ethBlock = bridge.NewEthBlock(3, blockTime.Add(10*time.Second))
	_, err = relayer.AddEthBlock(ethBlock)
	require.Equal(t, err, bridge.ErrInvalidBlockHeight, "expected height gap to error")

	ethBlock = bridge.NewEthBlock(2, blockTime)
	_, err = relayer.AddEthBlock(ethBlock)
	require.Equal(t, err, bridge.ErrInvalidBlockTime, "expected same timestamp to error")

	ethBlock = bridge.NewEthBlock(2, blockTime.Add(-1*time.Second))
	_, err = relayer.AddEthBlock(ethBlock)
	require.Equal(t, err, bridge.ErrInvalidBlockTime, "expected previous timestamp to error")
}

package main

import (
	"errors"
	"time"
)

var ErrKavaNotInitialized = errors.New("kava blocks must be added before adding eth blocks")
var ErrEthBlockAhead = errors.New("eth block must have a timestamp less than the last kava block")
var ErrInvalidBlockHeight = errors.New("block height must be 1 greater than last block")
var ErrInvalidBlockTime = errors.New("timestamp must be greater than last timestamp")

// ethBlock represents an ethereum block with incoming transfers at specific height and time
type ethBlock struct {
	height    uint64
	blockTime time.Time
}

func NewEthBlock(height uint64, blockTime time.Time) ethBlock {
	return ethBlock{
		height:    height,
		blockTime: blockTime,
	}
}

// kavaBlock represents an kavaereum block with incoming transfers at specific height and time
type kavaBlock struct {
	height    uint64
	blockTime time.Time
}

func NewKavaBlock(height uint64, blockTime time.Time) kavaBlock {
	return kavaBlock{
		height:    height,
		blockTime: blockTime,
	}
}

// kavaAction represents a mint to finish a transfer from ethereum to kava
type kavaAction struct {
}

type kavaActions []kavaAction

func NewKavaActions(actions ...kavaAction) kavaActions {
	return actions
}

// Relayer provides an interface to add source and destination blocks for
// a chain.  Each addition, emits actions identified by the source sequence.
// The nounce is unique to every transfer.
//
// Two algorithms -- fixed time, or next block / block ahead
//
// Throughput over retry latency
//
// we only support block ahead for eth -> kava transfers, and interface
// is designed only around that terminology
//
//
// Rejects eth blocks that are ahead of kava chain blocks
// this is OK since added eth blocks are behind the tip of the kava chain due to confirmations.
//
// Retries are sequenced among transfers
// Actions -- token, amount, address, event nounce, signer sequence
//
// Main interface AddEthBlock, AddKavaBlock -> actions
//
// We must sequence transactions from ethereum to kava --
//
// When we see a failure in a kava block, then we must pick moment in the ethereum block
// history to sequence the transaction that all relayers can agree on.  In order to do this,
// we emit these actions in the first ethereum block after the kava block at expereinced the failures.
//
// This results in a delay of about ~confirmation time~ until the transaction is retried again.
type relayer struct {
	lastKavaBlockHeight uint64
	lastKavaBlockTime   time.Time
	lastEthBlockHeight  uint64
	lastEthBlockTime    time.Time
}

// NewRelayer returns a new relayer
func NewRelayer() *relayer {
	return &relayer{}
}

// AddEthBlock adds an ethereum block to the relayer, returning any actions required
// to be send to the kava chain
func (r *relayer) AddEthBlock(blk ethBlock) (kavaActions, error) {
	if r.lastKavaBlockTime.IsZero() {
		return nil, ErrKavaNotInitialized
	}

	if !blk.blockTime.Before(r.lastKavaBlockTime) {
		return nil, ErrEthBlockAhead
	}

	if blk.height != r.lastEthBlockHeight+1 {
		return nil, ErrInvalidBlockHeight
	}

	if !blk.blockTime.After(r.lastEthBlockTime) {
		return nil, ErrInvalidBlockTime
	}

	r.lastEthBlockHeight = blk.height
	r.lastEthBlockTime = blk.blockTime

	return nil, nil
}

// AddKavaBlock adds kava block state to the relayer, never returning any actions
// we must ensure that kava blocks stay ahead of ethereum blocks
func (r *relayer) AddKavaBlock(blk kavaBlock) error {
	if blk.height != r.lastKavaBlockHeight+1 {
		return ErrInvalidBlockHeight
	}

	if !blk.blockTime.After(r.lastKavaBlockTime) {
		return ErrInvalidBlockTime
	}

	r.lastKavaBlockHeight = blk.height
	r.lastKavaBlockTime = blk.blockTime

	return nil
}

func main() {
}

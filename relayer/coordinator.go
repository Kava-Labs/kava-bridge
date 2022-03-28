package relayer

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrSourceBlockAhead is returned when the source blocks are ahead of the destination blocks
	ErrSourceBlockAhead = errors.New("source block must have a timestamp less than the last destination block")
	// ErrInvalidBlockHeight is returned when an invalid block height is added
	ErrInvalidBlockHeight = errors.New("block height must be 1 greater than last block")
	// ErrInvalidBlockTime is returned when an invalid block time is added
	ErrInvalidBlockTime = errors.New("timestamp must be greater than last timestamp")
	// ErrUnkownBlockOrigin is returned when in invalid block origin is added
	ErrUnknownBlockOrigin = errors.New("unknown block type")
)

// BlockOrigin determines if a block is a Source or Desitation block.  The coordinator uses
// this type to determine how to process an incoming block.
type BlockOrigin uint8

const (
	// Source represents a block origin of the source chain.
	Source BlockOrigin = iota + 1
	// Destination represents a block origin of the destinatino chain.
	Destination
)

// Payload represents any data that needs to be transferred cross-chain.
type Payload interface{}

// Block represents a block from a source or destiatnion chain with height, time, and zero or more payloads
type Block struct {
	Origin   BlockOrigin
	Height   uint64
	Time     time.Time
	Payloads []Payload
}

// NewBlock instantiates a new block with a origin, heigh, time, and payloads
func NewBlock(origin BlockOrigin, height uint64, time time.Time, payloads []Payload) Block {
	return Block{
		Origin:   origin,
		Height:   height,
		Time:     time,
		Payloads: payloads,
	}
}

// SigningOutput maps a payload in a block to an account nonce on the destination chain.  The
// coordinator is responsible for determining which nonce to assign to each payload.
type SigningOutput struct {
	nonce   uint64
	payload Payload
}

// NewSigningOutput returns a new signing output for a nonce and payload
func NewSigningOutput(nonce uint64, payload Payload) SigningOutput {
	return SigningOutput{
		nonce:   nonce,
		payload: payload,
	}
}

// Nonce returns the nonce of the signing output
func (s SigningOutput) Nonce() uint64 {
	return s.nonce
}

// Payload returns the payload of the signing output
func (s SigningOutput) Payload() Payload {
	return s.payload
}

// Coordinator sequences source and destination blocks between chains
// and outputs deterministicly ordered signing outputs
type Coordinator struct {
	// record the last height and time for the source chain
	lastSourceHeight uint64
	lastSourceTime   time.Time

	// record the last block height and time for the destination chain
	lastDestinationHeight uint64
	lastDestinationTime   time.Time

	nonce uint64
	// TODO: add bounds/backpressure, prevent OOM on syncing
	pending []Block

	outputs chan []SigningOutput
}

// NewCoordinator returns a new coordinator
// TODO: nonce and block initialization
func NewCoordinator() *Coordinator {
	return &Coordinator{
		nonce:   0,
		outputs: make(chan []SigningOutput),
	}
}

// AddBlock process a new source or destination block.  This is not go routine safe, as only
// one go routine should call this method.
//
// When adding a source block, it may block if there  are no avaible routines to read from
// the signing output channel.  If an error is returned, no signing outputs are generated
// and blocks may be safely retried.
func (c *Coordinator) AddBlock(ctx context.Context, block Block) (err error) {
	switch block.Origin {
	case Source:
		err = c.addSourceBlock(ctx, block)
	case Destination:
		err = c.addDestinationBlock(block)
	default:
		err = ErrUnknownBlockOrigin
	}

	return err
}

// SigningOutputs returns the output channel for all signing outputs.  This channel must be read from
// in order to process blocks in AddBlock.
func (c *Coordinator) SigningOutputs() <-chan []SigningOutput {
	return c.outputs
}

// Close closes the signing outputs channel.  This should used the upstream process adding blocks to
// signal downstream output readers to stop.  Calling AddBlock after Close will result in an error.
// TODO: add internal state to prevent writing to closed channel.
func (c *Coordinator) Close() {
	close(c.outputs)
}

func (c *Coordinator) addSourceBlock(ctx context.Context, block Block) error {
	if block.Height != c.lastSourceHeight+1 {
		return ErrInvalidBlockHeight
	}
	if !block.Time.After(c.lastSourceTime) {
		return ErrInvalidBlockTime
	}

	// TODO: add fixed time implementation for kava -> eth transfers
	// This only works for eth -> kava
	if !block.Time.Before(c.lastDestinationTime) {
		return ErrSourceBlockAhead
	}

	outputs := []SigningOutput{}
	for _, payload := range block.Payloads {
		c.nonce = c.nonce + 1
		outputs = append(outputs, NewSigningOutput(c.nonce, payload))
	}
	for _, p := range c.pending {
		if p.Time.Before(block.Time) {
			for _, payload := range p.Payloads {
				c.nonce = c.nonce + 1
				outputs = append(outputs, NewSigningOutput(c.nonce, payload))
			}

			c.pending = c.pending[1:]
		}
	}

	// Write outputs out atomically with cancellation
	if len(outputs) > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case c.outputs <- outputs:
		}
	}

	c.lastSourceHeight = block.Height
	c.lastSourceTime = block.Time

	return nil
}

func (c *Coordinator) addDestinationBlock(block Block) error {
	if block.Height != c.lastDestinationHeight+1 {
		return ErrInvalidBlockHeight
	}
	if !block.Time.After(c.lastDestinationTime) {
		return ErrInvalidBlockTime
	}

	// Any destination block with payloads must be added
	// to the pending queue
	if len(block.Payloads) > 0 {
		c.pending = append(c.pending, block)
	}

	c.lastDestinationHeight = block.Height
	c.lastDestinationTime = block.Time

	return nil
}

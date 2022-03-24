package relayer_test

import (
	"context"
	"time"

	"github.com/kava-labs/kava-bridge/relayer"
)

// payload is a payload only used for testing
type payload struct {
	sequence uint64
}

// newPayload creates a new test payload
func newPayload(sequence uint64) payload {
	return payload{
		sequence: sequence,
	}
}

// Sequence returns the sequence of the payload
func (p payload) Sequence() uint64 {
	return p.sequence
}

// addSecond is a test helper to add second durations to a time
func addSecond(t time.Time, d time.Duration) time.Time {
	return t.Add(d * time.Second)
}

// blocks is a test helper for creating a empty or populated slice of blocks
func blocks(blocks ...relayer.Block) []relayer.Block {
	return blocks
}

// eth is a test helper for returning a new eth block
func eth(height uint64, blockTime time.Time, payloads ...relayer.Payload) relayer.Block {
	return relayer.NewBlock(relayer.Source, height, blockTime, payloads)
}

// kava is a test helper for returning a new kava block
func kava(height uint64, blockTime time.Time, payloads ...relayer.Payload) relayer.Block {
	return relayer.NewBlock(relayer.Destination, height, blockTime, payloads)
}

// expectedOutput represents an expected signging output from a coordinator
type expectedOutput struct {
	nonce   uint64
	payload relayer.Payload
}

// outputs is a test helper for creating a empty or populated slice of signing outputs
func outputs(outputs ...[]expectedOutput) [][]expectedOutput {
	return outputs
}

// output is a test helper for creating a empty or populated slice of signing outputs
func output(outputs ...expectedOutput) []expectedOutput {
	return outputs
}

func newOutput(nonce uint64, payload relayer.Payload) expectedOutput {
	return expectedOutput{
		nonce:   nonce,
		payload: payload,
	}
}

// newCoordinator is a test helper for creating a new default eth to kava coordinator
func newCoordinator() *relayer.Coordinator {
	return relayer.NewCoordinator()
}

// orderedBlockPermutations is a recursive function that takes a slice pointer to store results in,
// an accumulating array for each permutation, and two sets of blocks.
//
// The number of permutations stored is (len(a) + len(b)) choose len(a).  Every combination
// of the two given block sequences is created and individual ordering of each sequence is kept.
func orderedBlockPermutations(permutations *[][]relayer.Block, acc, a, b []relayer.Block) {
	if len(a) == 0 {
		*permutations = append(*permutations, append(acc, b...))
		return
	}

	if len(b) == 0 {
		*permutations = append(*permutations, append(acc, a...))
		return
	}

	orderedBlockPermutations(permutations, append(acc, a[0]), a[1:], b)
	orderedBlockPermutations(permutations, append(acc, b[0]), b[1:], a)
}

// writeBlocks continually calls AddBlock, skipping to future blocks when adding a block results in an error.
// If all blocks in the slice error, then it returns an error; otherwise it collects all outputs generated from the
// coordinator.
func writeBlocks(ctx context.Context, c *relayer.Coordinator, blocks []relayer.Block) ([][]relayer.SigningOutput, error) {
	writeErr := make(chan error)

	go func() {
		addedBlockIndexes := make([]bool, len(blocks))
		defer c.Close()

		// add all blocks prioritizing the order
		for currentIndex, block := range blocks {
			// skip any blocks already added
			if addedBlockIndexes[currentIndex] {
				continue
			}

			// TODO: we are only allowed to add the next block of the next chain
			// TODO: account for added block index in failure case (currently attempts to add them again)
			for {
				// attempt to add block
				origErr := c.AddBlock(ctx, block)

				// if added break inner loop and continue
				if origErr == nil {
					addedBlockIndexes[currentIndex] = true
					break
				}

				// if failed, try next blocks in sequence to get unstuck
				for i := currentIndex + 1; i < len(blocks); i++ {
					err := c.AddBlock(ctx, blocks[i])

					// if added then break
					if err == nil {
						addedBlockIndexes[i] = true
						break
					}

					if i == len(blocks)-1 {
						writeErr <- origErr
					}
				}
			}
		}
	}()

	var outputs [][]relayer.SigningOutput

	for {
		select {
		case err := <-writeErr:
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		case p, more := <-c.SigningOutputs():
			if more {
				outputs = append(outputs, p)
			} else {
				return outputs, nil
			}
		}
	}
}

// startNullReader starts a null reader and returns a stop func
func startNullReader(c *relayer.Coordinator) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-c.SigningOutputs():
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancel
}

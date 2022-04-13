package p2p

import (
	"fmt"
	"time"
)

func retry(attempts int, sleep time.Duration, f func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			time.Sleep(sleep)
		}
		err = f()
		if err == nil {
			return nil
		}

		log.Debugw("retry", "attempt", i, "error", err)
	}

	return fmt.Errorf("failed after %d attempts: %w", attempts, err)
}

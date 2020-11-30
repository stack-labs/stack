package client

import (
	"context"
	"math"
	"time"
)

type BackoffFunc func(ctx context.Context, req Request, attempts int) (time.Duration, error)

// exponential backoff is a function x^e multiplied by a factor of 0.1 second.
// Result is limited to 2 minute.
func exponentialBackoff(ctx context.Context, req Request, attempts int) (time.Duration, error) {
	if attempts > 13 {
		return 2 * time.Minute, nil
	}
	return time.Duration(math.Pow(float64(attempts), math.E)) * time.Millisecond * 100, nil
}

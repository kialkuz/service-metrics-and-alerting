package postgresql

import (
	"context"
	"time"
)

const retryAttempts = 3

func WithRetry[T any](
	ctx context.Context,
	cfg RetryConfig,
	fn func(ctx context.Context) (T, error),
) (T, error) {
	var zero T

	res, err := fn(ctx)
	if err == nil {
		return res, nil
	}

	for i := 0; i < cfg.Attempts; i++ {
		timer := time.NewTimer(retryAttempts * time.Millisecond)

		select {
		case <-ctx.Done():
			timer.Stop()
			return zero, ctx.Err()

		case <-timer.C:
			res, err := fn(ctx)
			if err == nil {
				return res, nil
			}
		}
	}

	return zero, err
}

package postgresql

import (
	"context"
	"kialkuz/service-metrics-and-alerting/pkg/pgerrors"
	"time"
)

var additionalReplaySendInterval = map[int]int{
	1: 1,
	2: 3,
	3: 5,
}

func WithRetry[T any](
	ctx context.Context,
	cfg RetryConfig,
	classifier *pgerrors.PostgresErrorClassifier,
	fn func(ctx context.Context) (T, error),
) (T, error) {
	var zero T

	res, err := fn(ctx)

	if err != nil {
		classification := classifier.Classify(err)
		if classification == pgerrors.Retriable {
			for i := 1; i <= cfg.Attempts; i++ {
				timer := time.NewTimer(time.Duration(additionalReplaySendInterval[i]) * time.Second)

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
		}

		return zero, err
	}

	return res, nil
}

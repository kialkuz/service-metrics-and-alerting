package retrier

import (
	"context"
	"kialkuz/service-metrics-and-alerting/pkg/errorclassifier"
	"time"
)

var additionalReplaySendInterval = map[int]int{
	1: 1,
	2: 3,
	3: 5,
}

//go:generate go run go.uber.org/mock/mockgen -source=retrier.go -destination=mocks/retrier_mock.go -package=mocks -typed
type Retrier interface {
	WithRetry(ctx context.Context, fn func(ctx context.Context) error) error
}

type RetrierService struct {
	cfg        RetryConfig
	classifier errorclassifier.ErrorClassifier
}

func NewRetrierService(cfg RetryConfig, classifier errorclassifier.ErrorClassifier) Retrier {
	return RetrierService{
		cfg:        cfg,
		classifier: classifier,
	}
}

func (r RetrierService) WithRetry(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	err := fn(ctx)

	if err != nil {
		classification := r.classifier.Classify(err)
		if classification == errorclassifier.Retriable {
			for i := 1; i <= r.cfg.Attempts; i++ {
				timer := time.NewTimer(time.Duration(additionalReplaySendInterval[i]) * time.Second)

				select {
				case <-ctx.Done():
					timer.Stop()
					return ctx.Err()

				case <-timer.C:
					err := fn(ctx)
					if err == nil {
						return nil
					}
				}
			}
		}

		return err
	}

	return nil
}

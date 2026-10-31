package collectalgo

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/storage"
	"kialkuz/service-metrics-and-alerting/internal/model"
	service "kialkuz/service-metrics-and-alerting/internal/service/agent"
	"kialkuz/service-metrics-and-alerting/pkg/errorclassifier"

	"golang.org/x/sync/errgroup"
)

const CollectSingle = "collect_single"

type Single struct {
	metricsService service.MetricsAgentService
	collectConfig  CollectConfig
}

func NewSingle(
	metricsService service.MetricsAgentService,
	collectConfig CollectConfig,
) Single {
	return Single{
		metricsService: metricsService,
		collectConfig:  collectConfig,
	}
}

func (s Single) CollectAndSend(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g := new(errgroup.Group)

	metrics := make(chan dto.Metrics, s.collectConfig.RateLimit)
	var err error

	g.Go(func() error {
		metrics, err = s.collect(metrics)
		if err != nil {
			return err
		}

		return nil
	})

	for w := 1; w <= s.collectConfig.RateLimit; w++ {
		g.Go(func() error {
			return s.send(ctx, metrics)
		})
	}

	if err := g.Wait(); err != nil {
		cancel()
	}

	return nil
}

func (s Single) collect(metrics chan dto.Metrics) (chan dto.Metrics, error) {
	defer close(metrics)

	for fieldName, fieldValue := range s.metricsService.CollectCounter() {
		metrics <- dto.Metrics{
			ID:    fieldName,
			MType: model.Counter,
			Delta: &fieldValue,
		}
	}

	gaugeMetrics, err := s.metricsService.CollectGauge()
	if err != nil {
		return nil, err
	}
	for fieldName, fieldValue := range gaugeMetrics {
		metrics <- dto.Metrics{
			ID:    fieldName,
			MType: model.Gauge,
			Value: &fieldValue,
		}
	}

	return metrics, nil
}

func (s Single) send(ctx context.Context, metrics <-chan dto.Metrics) error {
	for metric := range metrics {
		_, err := storage.WithRetry(
			ctx,
			storage.RetryConfig{Attempts: retryAttempts},
			errorclassifier.NewNetworkErrorClassifier(),
			func(ctx context.Context) (struct{}, error) {
				_, err := s.metricsService.SendSingleMetric(ctx, metric)
				return struct{}{}, err
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

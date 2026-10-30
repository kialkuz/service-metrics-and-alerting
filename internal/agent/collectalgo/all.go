package collectalgo

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/storage"
	"kialkuz/service-metrics-and-alerting/internal/model"
	service "kialkuz/service-metrics-and-alerting/internal/service/agent"
	"kialkuz/service-metrics-and-alerting/pkg/errorclassifier"
)

const (
	CollectAll    = "collect_all"
	retryAttempts = 3
)

type All struct {
	metricsService service.MetricsAgentService
	collectConfig  CollectConfig
}

func NewAll(
	metricsService service.MetricsAgentService,
	collectConfig CollectConfig,
) All {
	return All{
		metricsService: metricsService,
		collectConfig:  collectConfig,
	}
}

func (a All) CollectAndSend(ctx context.Context) error {
	metrics, err := a.collect()
	if err != nil {
		return err
	}

	err = a.send(ctx, metrics)
	if err != nil {
		return err
	}

	return nil
}

func (a All) collect() ([]dto.Metrics, error) {
	metricsGauge, err := a.metricsService.CollectGauge()
	if err != nil {
		return nil, err
	}

	metricsCounter := a.metricsService.CollectCounter()

	metrics := make([]dto.Metrics, 0, len(metricsCounter)+len(metricsGauge))

	for fieldName, fieldValue := range metricsCounter {
		metrics = append(metrics, dto.Metrics{
			ID:    fieldName,
			MType: model.Counter,
			Delta: &fieldValue,
		})
	}

	for fieldName, fieldValue := range metricsGauge {
		metrics = append(metrics, dto.Metrics{
			ID:    fieldName,
			MType: model.Gauge,
			Value: &fieldValue,
		})
	}

	return metrics, nil
}

func (a All) send(ctx context.Context, metrics []dto.Metrics) error {
	_, err := storage.WithRetry(
		ctx,
		storage.RetryConfig{Attempts: retryAttempts},
		errorclassifier.NewNetworkErrorClassifier(),
		func(ctx context.Context) (struct{}, error) {
			_, err := a.metricsService.SendListMetrics(ctx, metrics)
			return struct{}{}, err
		},
	)
	if err != nil {
		return err
	}

	return nil
}

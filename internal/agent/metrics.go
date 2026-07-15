package agent

import (
	"context"
	"errors"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/model"
	service "kialkuz/service-metrics-and-alerting/internal/service/agent"
	"time"

	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"

	"golang.org/x/sync/errgroup"
)

const defaultCountAdditionalReplaySendMetrics = 3

var additionalReplaySendInterval = map[int]int{
	1: 1,
	2: 3,
	3: 5,
}

type MetricsAgent struct {
	metricsService service.MetricsAgentService
}

func NewMetricsAgent(metricsService service.MetricsAgentService) *MetricsAgent {
	return &MetricsAgent{
		metricsService: metricsService,
	}
}

func (a *MetricsAgent) CollectAndSend(reportInterval, pollInterval, rateLimit int) error {
	var err error

	now := time.Now()

	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g := new(errgroup.Group)

	for {
		select {
		case <-ctx.Done():
			return err
		case <-ticker.C:
			if time.Now().After(now.Add(time.Duration(reportInterval) * time.Second)) {
				metrics := make(chan dto.Metrics, rateLimit)

				g.Go(func() error {
					metrics, err = a.collect(metrics)
					if err != nil {
						return err
					}

					return nil
				})

				for w := 1; w <= rateLimit; w++ {
					g.Go(func() error {
						return a.send(metrics)
					})
				}

				if err := g.Wait(); err != nil {
					cancel()
				}

				now = time.Now()
			}
		}
	}
}

func (a *MetricsAgent) collect(metrics chan dto.Metrics) (chan dto.Metrics, error) {
	defer close(metrics)

	for fieldName, fieldValue := range a.metricsService.CollectCounter() {
		metrics <- dto.Metrics{
			ID:    fieldName,
			MType: model.Counter,
			Delta: &fieldValue,
		}
	}

	gaugeMetrics, err := a.metricsService.CollectGauge()
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

func (a *MetricsAgent) send(metrics <-chan dto.Metrics) error {
	for metric := range metrics {
		_, err := a.metricsService.SendSingleMetric(metric)
		if errors.Is(err, pkgErrors.ErrSendMetrics) {
			for i := 1; i <= defaultCountAdditionalReplaySendMetrics; i++ {
				timer := time.NewTimer(time.Duration(additionalReplaySendInterval[i]) * time.Second)
				<-timer.C

				if _, err = a.metricsService.SendSingleMetric(metric); err == nil {
					break
				}
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

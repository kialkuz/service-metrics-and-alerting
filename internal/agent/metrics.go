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
	metrics        chan dto.Metrics
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
				a.metrics = make(chan dto.Metrics, rateLimit)

				go a.collect(a.metrics)

				for w := 1; w <= rateLimit; w++ {
					g.Go(a.send)
				}

				if err := g.Wait(); err != nil {
					cancel()
				}

				now = time.Now()
			}
		}
	}
}

func (a *MetricsAgent) collect(metrics chan<- dto.Metrics) {
	defer close(metrics)

	for fieldName, fieldValue := range a.metricsService.CollectCounter() {
		metrics <- dto.Metrics{
			ID:    fieldName,
			MType: model.Counter,
			Delta: &fieldValue,
		}
	}

	for fieldName, fieldValue := range a.metricsService.CollectGauge() {
		metrics <- dto.Metrics{
			ID:    fieldName,
			MType: model.Gauge,
			Value: &fieldValue,
		}
	}
}

func (a *MetricsAgent) send() error {
	for metric := range a.metrics {
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

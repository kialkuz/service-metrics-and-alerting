package agent

import (
	"context"
	"errors"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/model"
	service "kialkuz/service-metrics-and-alerting/internal/service/agent"
	"time"

	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"
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

func (a *MetricsAgent) Collect(reportInterval, pollInterval int) error {
	var err error

	now := time.Now()

	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return err
		case <-ticker.C:
			if time.Now().After(now.Add(time.Duration(reportInterval) * time.Second)) {
				err = a.send()

				if errors.Is(err, pkgErrors.ErrSendMetrics) {
					for i := 1; i <= defaultCountAdditionalReplaySendMetrics; i++ {
						timer := time.NewTimer(time.Duration(additionalReplaySendInterval[i]) * time.Second)
						<-timer.C

						if err = a.send(); err == nil {
							break
						}
					}

					if err != nil {
						cancel()
					}
				}

				now = time.Now()
			}
		}
	}
}

func (a *MetricsAgent) send() error {
	var metricsForSend []dto.Metrics

	for fieldName, fieldValue := range a.metricsService.CollectCounter() {
		metricsForSend = append(metricsForSend, dto.Metrics{
			ID:    fieldName,
			MType: model.Counter,
			Delta: &fieldValue,
		})
	}

	for fieldName, fieldValue := range a.metricsService.CollectGauge() {
		metricsForSend = append(metricsForSend, dto.Metrics{
			ID:    fieldName,
			MType: model.Gauge,
			Value: &fieldValue,
		})
	}

	if len(metricsForSend) > 0 {
		_, err := a.metricsService.SendListMetrics(metricsForSend)
		if err != nil {
			return fmt.Errorf("%w: %w", pkgErrors.ErrSendMetrics, err)
		}
	}

	return nil
}

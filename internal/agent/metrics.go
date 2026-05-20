package agent

import (
	"context"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/model"
	service "kialkuz/service-metrics-and-alerting/internal/service/agent"
	"log"
	"time"
)

type MetricsAgent struct {
	metricsService service.MetricsAgentService
}

func NewMetricsAgent(metricsService service.MetricsAgentService) *MetricsAgent {
	return &MetricsAgent{
		metricsService: metricsService,
	}
}

func (a *MetricsAgent) Collect(reportInterval, pollInterval int) {
	now := time.Now()

	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if time.Now().After(now.Add(time.Duration(reportInterval) * time.Second)) {
				a.send()

				now = time.Now()
			}
		}
	}
}

func (a *MetricsAgent) send() {
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
			log.Println(fmt.Errorf("error agent send metric: %s", err))
		}
	}
}
